package raster

import (
	"image"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

func TestDash(t *testing.T) {
	// A dashed border of width 10 around a 90×90 rect: the top side takes
	// three segments of 30, dash, gap, dash.
	r := scene.Rect{X: 0, Y: 0, W: 90, H: 90}
	w := [4]float32{10, 10, 10, 10}
	// The center of the first dash, the edge between the dash and the
	// gap, the center of the gap, and the ramp past the edge.
	for _, c := range []struct {
		px, py float32
		want   float32
	}{
		{15, 0.5, 1},
		{30, 0.5, 0.5},
		{45, 0.5, 0},
		{31, 0.5, 0},
		{29, 0.5, 1},
	} {
		if got := dash(r, w, c.px, c.py); got != c.want {
			t.Errorf("dash at (%v, %v) = %v, want %v", c.px, c.py, got, c.want)
		}
	}
	// A short side is one dash across.
	if got := dash(scene.Rect{X: 0, Y: 0, W: 20, H: 90}, w, 10, 0.5); got != 1 {
		t.Errorf("a side that is one dash = %v, want 1", got)
	}
	// A pixel belongs to the side it is nearest: at the bottom edge the
	// pattern runs along the bottom, whose border is three times as wide
	// and takes the whole side as one dash, where the top of the same
	// width takes a gap there.
	tall := scene.Rect{X: 0, Y: 0, W: 90, H: 90}
	wide2 := [4]float32{10, 10, 30, 10}
	if got := dash(tall, wide2, 45, 89.5); got != 1 {
		t.Errorf("the bottom's own dash = %v, want 1", got)
	}
	if got := dash(tall, wide2, 45, 0.5); got != 0 {
		t.Errorf("the top's gap at the same x = %v, want 0", got)
	}
	// A side of no width never takes the pixel.
	noTop := [4]float32{0, 10, 10, 10}
	if got := dash(r, noTop, 0.5, 0.5); got != 1 {
		t.Errorf("a corner pixel with no top border = %v, want the left side's dash 1", got)
	}
}

func TestErase(t *testing.T) {
	p := []byte{255, 128, 0, 255}
	erase(p, 1)
	if p[0] != 0 || p[1] != 0 || p[2] != 0 || p[3] != 0 {
		t.Errorf("a full erase = %v, want nothing", p)
	}
	p = []byte{255, 128, 0, 255}
	erase(p, 0.5)
	if p[0] != 128 || p[1] != 64 || p[2] != 0 || p[3] != 128 {
		t.Errorf("a half erase = %v, want half of each", p)
	}
	p = []byte{255, 255, 255, 255}
	erase(p, 0)
	if p[0] != 255 || p[3] != 255 {
		t.Errorf("no erase = %v, want the pixel as it is", p)
	}
}

func TestHole(t *testing.T) {
	dst := NewImage(8, 8)
	for i := range dst.Pix {
		dst.Pix[i] = 255
	}
	s := &scene.Scene{Width: 8, Height: 8}
	r := renderer{dst: dst, s: s, area: image.Rect(0, 0, 8, 8)}
	r.updateBounds()
	r.hole(&scene.Op{Kind: scene.OpHole, Rect: scene.Rect{X: 2, Y: 2, W: 4, H: 4}})
	for y := range 8 {
		for x := range 8 {
			inside := x >= 2 && x < 6 && y >= 2 && y < 6
			if got := dst.Pix[4*(y*8+x)]; (inside && got != 0) || (!inside && got != 255) {
				t.Fatalf("hole at (%d, %d) = %d, want %d", x, y, got, map[bool]byte{true: 0, false: 255}[inside])
			}
		}
	}

	// A hole of half opacity takes half of each pixel.
	for i := range dst.Pix {
		dst.Pix[i] = 255
	}
	r.hole(&scene.Op{Kind: scene.OpHole, Rect: scene.Rect{X: 2, Y: 2, W: 4, H: 4}, Opacity: 0.5})
	if got := dst.Pix[4*(3*8+3)]; got != 128 {
		t.Errorf("a half hole = %d, want 128", got)
	}
	if got := dst.Pix[0]; got != 255 {
		t.Errorf("outside the hole = %d, want 255", got)
	}

	// A rounded hole takes what the corner's arc leaves.
	for i := range dst.Pix {
		dst.Pix[i] = 255
	}
	radii := [4]float32{4, 0, 0, 0}
	r.hole(&scene.Op{Kind: scene.OpHole, Rect: scene.Rect{X: 2, Y: 2, W: 4, H: 4}, Radii: radii})
	if got := dst.Pix[4*(2*8+2)]; got != 255 {
		t.Errorf("the pixel past the arc = %d, want 255", got)
	}
	if got := dst.Pix[4*(4*8+4)]; got != 0 {
		t.Errorf("inside the arc = %d, want 0", got)
	}
	if got := dst.Pix[4*(2*8+3)]; got == 255 || got == 0 {
		t.Errorf("under the arc = %d, want part of the pixel", got)
	}
}
