package raster

import (
	"math"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

// nv12Frame builds a 4×4 NV12 frame of video: a luma plane of the row
// values, and a chroma plane of the Cb, Cr pairs of its 2×2 samples,
// each with the padding its stride leaves unused.
func nv12Frame(t *testing.T, luma [4][4]byte, chroma [2][2][2]byte) (*scene.Image, *scene.Surface) {
	t.Helper()
	lp := make([]byte, 0, 4*8)
	for _, row := range luma {
		lp = append(lp, row[0], row[1], row[2], row[3], 0xAA, 0xAA, 0xAA, 0xAA)
	}
	cp := make([]byte, 0, 2*8)
	for _, row := range chroma {
		cp = append(cp, row[0][0], row[0][1], row[1][0], row[1][1], 0xAA, 0xAA, 0xAA, 0xAA)
	}
	img := scene.NewImageRGBA(4, 4, nil)
	s := &scene.Surface{
		Format: scene.SurfaceNV12,
		Matrix: scene.BT601 + scene.FullRange,
		Planes: [2]scene.Plane{{Pix: lp, Stride: 8}, {Pix: cp, Stride: 8}},
	}
	return img, s
}

func about(a, b float32) bool { return math.Abs(float64(a-b)) < 1e-4 }

func TestSampleNV12(t *testing.T) {
	// Luma columns 0 and 2 hold white, 1 and 3 black; the chroma is
	// neutral everywhere, so the color follows the luma.
	luma := [4][4]byte{{255, 0, 255, 0}, {255, 0, 255, 0}, {255, 0, 255, 0}, {255, 0, 255, 0}}
	chroma := [2][2][2]byte{{{128, 128}, {128, 128}}, {{128, 128}, {128, 128}}}
	img, s := nv12Frame(t, luma, chroma)
	all := scene.Rect{X: 0, Y: 0, W: 4, H: 4}

	// A texel's center is its sample.
	if c := sampleNV12(img, s, 0.5, 1.5, all); !about(c[0], 1) || !about(c[1], 1) || !about(c[2], 1) || c[3] != 1 {
		t.Errorf("the center of a white texel = (%v, %v, %v, %v)", c[0], c[1], c[2], c[3])
	}
	if c := sampleNV12(img, s, 1.5, 1.5, all); !about(c[0], 0) || !about(c[1], 0) || !about(c[2], 0) {
		t.Errorf("the center of a black texel = (%v, %v, %v, %v)", c[0], c[1], c[2], c[3])
	}
	// Halfway between a white and a black column is a half gray.
	if c := sampleNV12(img, s, 1.0, 1.5, all); !about(c[0], 0.5) || !about(c[1], 0.5) || !about(c[2], 0.5) {
		t.Errorf("between the columns = (%v, %v, %v), want half gray", c[0], c[1], c[2])
	}
	// A sample past a side holds the edge texel.
	if c := sampleNV12(img, s, -5, 1.5, all); !about(c[0], 1) {
		t.Errorf("left of the frame = %v, want the left edge's white", c[0])
	}
	if c := sampleNV12(img, s, 9, 1.5, all); !about(c[0], 0) {
		t.Errorf("right of the frame = %v, want the right edge's black", c[0])
	}

	// A frame's source rect keeps the sampling inside it. This frame's
	// columns 2 and 3 are black and white — the opposite of the frame
	// above at its own edges — so a sample clamped to the rect tells the
	// rect's edges from the frame's.
	luma2 := [4][4]byte{{255, 255, 0, 255}, {255, 255, 0, 255}, {255, 255, 0, 255}, {255, 255, 0, 255}}
	img2, s2 := nv12Frame(t, luma2, chroma)
	src := scene.Rect{X: 2, Y: 0, W: 2, H: 4}
	if c := sampleNV12(img2, s2, -5, 1.5, src); !about(c[0], 0) {
		t.Errorf("left of the source rect = %v, want its left edge's black", c[0])
	}
	if c := sampleNV12(img2, s2, 9, 1.5, src); !about(c[0], 1) {
		t.Errorf("right of the source rect = %v, want its right edge's white", c[0])
	}
}

func TestSampleNV12Chroma(t *testing.T) {
	// The luma is white everywhere; the chroma's left column holds Cb 0
	// and the right one Cb 255, so a sample's blue names the chroma texel
	// at half its position: Cb 0 takes the blue below the luma, Cb 255
	// holds it at the clamp, and half of each is in between.
	luma := [4][4]byte{{255, 255, 255, 255}, {255, 255, 255, 255}, {255, 255, 255, 255}, {255, 255, 255, 255}}
	chroma := [2][2][2]byte{{{0, 128}, {255, 128}}, {{0, 128}, {255, 128}}}
	img, s := nv12Frame(t, luma, chroma)
	all := scene.Rect{X: 0, Y: 0, W: 4, H: 4}

	// 1 + 2·(1−Kb)·(Cb − 1/2), at the chroma texels' own centers 0.5 and
	// 2.5 in luma, and halfway between them.
	const (
		kb         = 0.114
		noBlue     = 1 - 2*(1-kb)*(128.0/255) // Cb 0: about 0.1175
		allBlue    = 1 + 2*(1-kb)*(127.0/255) // Cb 255: past the clamp
		halfBlueCb = 63.75/255 - 128.0/255    // the two, averaged
		halfBlue   = 1 + 2*(1-kb)*halfBlueCb  // about 0.5535
	)
	if c := sampleNV12(img, s, 0.5, 0.5, all); !about(c[2], noBlue) || !about(c[0], 1) || !about(c[1], 1) {
		t.Errorf("the left chroma texel = (%v, %v, %v), want (%v, 1, %v)", c[0], c[1], c[2], 1, noBlue)
	}
	if c := sampleNV12(img, s, 2.5, 0.5, all); !about(c[2], 1) {
		t.Errorf("the right chroma texel = (%v, %v, %v), want blue held at the clamp", c[0], c[1], c[2])
	}
	if c := sampleNV12(img, s, 1.5, 0.5, all); !about(c[2], halfBlue) {
		t.Errorf("between the chroma texels, blue = %v, want %v", c[2], halfBlue)
	}
}
