package raster

import (
	"math"
	"testing"
)

func TestBlendSubpixel(t *testing.T) {
	// A pixel starts black, transparent, and takes a color's channels
	// each by its own coverage; the pixel is BGRA, the color RGBA.
	p := []byte{0, 0, 0, 0}
	blendSubpixel(p, [3]float32{0, 0, 1}, [3]float32{1, 1, 1}, 1)
	if p[0] != 255 || p[1] != 0 || p[2] != 0 || p[3] != 255 {
		t.Fatalf("a blue over nothing = %v, want opaque blue", p)
	}
	p = []byte{0, 0, 0, 0}
	blendSubpixel(p, [3]float32{1, 0, 0}, [3]float32{0.5, 0.5, 0.5}, 1)
	if p[0] != 0 || p[1] != 0 || p[2] != 128 || p[3] != 128 {
		t.Fatalf("a half-covered red = %v, want %v", p, []byte{0, 0, 128, 128})
	}
	// Over an opaque black pixel, a full red keeps it opaque and turns it
	// red; the alpha takes the channels' mean.
	p = []byte{0, 0, 0, 255}
	blendSubpixel(p, [3]float32{1, 0, 0}, [3]float32{1, 1, 1}, 1)
	if p[0] != 0 || p[1] != 0 || p[2] != 255 || p[3] != 255 {
		t.Fatalf("a red over black = %v, want opaque red", p)
	}
	// Half of a green over an opaque blue greens the blue channel by half
	// and keeps the rest of it.
	p = []byte{255, 0, 0, 255}
	blendSubpixel(p, [3]float32{0, 1, 0}, [3]float32{1, 1, 1}, 0.5)
	if p[0] != 128 || p[1] != 128 || p[2] != 0 || p[3] != 255 {
		t.Fatalf("a half green over blue = %v, want %v", p, []byte{128, 128, 0, 255})
	}
}

func TestUnpremul(t *testing.T) {
	if got := unpremul([4]float32{0.25, 0.5, 1, 0.5}); got != [3]float32{0.5, 1, 2} {
		t.Errorf("unpremul = %v, want {0.5 1 2}", got)
	}
	if got := unpremul([4]float32{0.3, 0.6, 0.9, 0}); got != ([3]float32{}) {
		t.Errorf("unpremul of nothing = %v, want zeros", got)
	}
	if got := unpremul([4]float32{0.3, 0.6, 0.9, -1}); got != ([3]float32{}) {
		t.Errorf("unpremul of a negative alpha = %v, want zeros", got)
	}
}

func TestToLinear(t *testing.T) {
	for _, c := range []struct{ v, want float32 }{
		{0, 0},
		{1, 1},
		// The piecewise boundary and its value, and a middle gray.
		{0.04045, 0.04045 / 12.92},
		{0.5, 0.21404114},
	} {
		if got := toLinear(c.v); math.Abs(float64(got-c.want)) > 1e-6 {
			t.Errorf("toLinear(%v) = %v, want %v", c.v, got, c.want)
		}
	}
	if got := toSRGB(toLinear(0.5)); math.Abs(float64(got-0.5)) > 1e-4 {
		t.Errorf("toSRGB(toLinear(0.5)) = %v, want 0.5", got)
	}
}
