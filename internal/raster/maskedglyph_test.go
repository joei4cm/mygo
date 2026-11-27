package raster

import (
	"image"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

func TestAtlasSample(t *testing.T) {
	// A mask atlas of 8×8 coverage bytes, whose cell holds 100 at its
	// bottom left texel and 200 at its bottom right one.
	a := &scene.Atlas{BPP: 1, W: 8, H: 8, Pix: make([]byte, 64)}
	a.Pix[3*8+2] = 100
	a.Pix[3*8+3] = 200
	// A texel's center is its value.
	if c := atlasSample(a, 3.5, 3.5, 2, 2, 2, 2); c[0] != 200.0/255 {
		t.Errorf("the cell's bottom right texel = %v, want %v", c[0], 200.0/255)
	}
	// Halfway across the cell is the two texels' average.
	if c := atlasSample(a, 3.0, 3.5, 2, 2, 2, 2); c[0] != 150.0/255 {
		t.Errorf("between the texels = %v, want %v", c[0], 150.0/255)
	}
	// A sample beside the cell holds the cell's edge texel.
	if c := atlasSample(a, 1.0, 3.5, 2, 2, 2, 2); c[0] != 100.0/255 {
		t.Errorf("past the cell = %v, want the edge texel's %v", c[0], 100.0/255)
	}
	// A colored atlas of four bytes a pixel takes all of them.
	b := &scene.Atlas{BPP: 4, W: 2, H: 2, Pix: []byte{10, 20, 30, 40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	if c := atlasSample(b, 0.5, 0.5, 0, 0, 2, 2); c[0] != 10.0/255 || c[1] != 20.0/255 || c[2] != 30.0/255 || c[3] != 40.0/255 {
		t.Errorf("the colored atlas texel = %v, want its bytes", c)
	}
}

type fakeEffectPixels struct{ c [4]float32 }

func (f *fakeEffectPixels) Begin(*scene.EffectOp, scene.Rect, [4]float32) {}
func (f *fakeEffectPixels) Color(x, y float32, b *scene.BackdropImage) [4]float32 {
	return f.c
}

func TestMaskedEffectGlyphs(t *testing.T) {
	// A mask atlas whose cell holds full coverage at the one texel a
	// glyph of three at (1, 1) reads at its device pixel (3, 3).
	atlas := &scene.Atlas{BPP: 1, W: 8, H: 8, Pix: make([]byte, 64)}
	atlas.Pix[4*8+4] = 255
	s := &scene.Scene{
		Width:     8,
		Height:    8,
		MaskAtlas: atlas,
		Glyphs: []scene.Glyph{
			{X: 1, Y: 1, W: 3, H: 3, U: 2, V: 2, UW: 3, VH: 3, Color: scene.Color{R: 1, G: 0, B: 0, A: 1}},
			{X: 5, Y: 5, W: 2, H: 2, U: 6, V: 6, UW: 2, VH: 2, Color: scene.Color{R: 0, G: 1, B: 0, A: 1}, Colored: true},
		},
	}
	dst := NewImage(8, 8)
	r := renderer{dst: dst, s: s, area: image.Rect(0, 0, 8, 8)}
	r.updateBounds()
	fx := &scene.EffectOp{Glyphs: [2]int32{0, 2}}
	px := &fakeEffectPixels{c: [4]float32{1, 0, 0, 1}}
	r.maskedEffect(&scene.Op{Kind: scene.OpEffect}, fx, px)

	p := dst.Pix[4*(3*8+3):][:4]
	if p[0] != 0 || p[1] != 0 || p[2] < 200 || p[3] < 200 {
		t.Errorf("the glyph's pixel = %v, want about %v", p, []byte{0, 0, 255, 255})
	}
	// The colored glyph took nothing: masked glyphs come from the mask
	// atlas alone, and colored ones are someone else's.
	if q := dst.Pix[4*(5*8+5):][:4]; q[0] != 0 || q[1] != 0 || q[2] != 0 || q[3] != 0 {
		t.Errorf("a colored glyph = %v, want untouched", q)
	}
	if q := dst.Pix[4*(1*8+1):][:4]; q[0] != 0 || q[3] != 0 {
		t.Errorf("a pixel of no coverage = %v, want untouched", q)
	}

	// Half opacity takes half of the pixel.
	for i := range dst.Pix {
		dst.Pix[i] = 0
	}
	r.maskedEffect(&scene.Op{Kind: scene.OpEffect, Opacity: 0.5}, fx, px)
	if got := dst.Pix[4*(3*8+3)+3]; got < 110 || got > 145 {
		t.Errorf("a half glyph's alpha = %d, want about 128", got)
	}
}
