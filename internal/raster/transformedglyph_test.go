package raster

import (
	"image"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

func TestTransformedGlyphs(t *testing.T) {
	// A color atlas whose cell holds an opaque red texel at its top left,
	// and a mask atlas whose cell holds full coverage a texel over.
	colors := &scene.Atlas{BPP: 4, W: 4, H: 4, Pix: make([]byte, 64)}
	colors.Pix[0], colors.Pix[3] = 255, 255
	masks := &scene.Atlas{BPP: 1, W: 4, H: 4, Pix: make([]byte, 16)}
	masks.Pix[1*4+1] = 255

	// A colored glyph of four at the origin, drawn by an op translated by
	// half a pixel, takes the pixel between its atlas's first two texels
	// at half of each, and the rows under them too.
	colors.Pix[4], colors.Pix[5], colors.Pix[7] = 0, 255, 255
	s := &scene.Scene{
		Width:      4,
		Height:     4,
		ColorAtlas: colors,
		MaskAtlas:  masks,
		Glyphs: []scene.Glyph{
			{X: 0, Y: 0, W: 4, H: 4, U: 0, V: 0, UW: 4, VH: 4, Color: scene.Color{R: 255, G: 0, B: 0, A: 255}, Colored: true},
		},
	}
	dst := NewImage(4, 4)
	r := renderer{dst: dst, s: s, area: image.Rect(0, 0, 4, 4)}
	r.updateBounds()
	r.transformedGlyphs(&scene.Op{Kind: scene.OpGlyphs, Transform: scene.Translation(0.5, 0.5), Start: 0, End: 1})
	p := dst.Pix[4*(1*4+1):][:4]
	if p[0] != 0 || p[1] != 64 || p[2] != 64 || p[3] != 128 {
		t.Errorf("the colored glyph = %v, want a quarter of red and of green", p)
	}
	if q := dst.Pix[4*(3*4+1):][:4]; q[0] != 0 || q[3] != 0 {
		t.Errorf("past the glyph's texels = %v, want untouched", q)
	}

	// A masked glyph carries its own transform on top of its op's: one of
	// (1, 0) more lands it at (2, 1), where it blues by the coverage its
	// atlas cell holds there.
	s2 := &scene.Scene{
		Width:      4,
		Height:     4,
		ColorAtlas: colors,
		MaskAtlas:  masks,
		Glyphs: []scene.Glyph{
			{Transform: scene.Translation(1, 0), X: 0, Y: 0, W: 2, H: 2, U: 1, V: 1, UW: 2, VH: 2, Color: scene.Color{R: 0, G: 0, B: 255, A: 255}},
		},
	}
	dst2 := NewImage(4, 4)
	r2 := renderer{dst: dst2, s: s2, area: image.Rect(0, 0, 4, 4)}
	r2.updateBounds()
	r2.transformedGlyphs(&scene.Op{Kind: scene.OpGlyphs, Transform: scene.Translation(1, 1), Start: 0, End: 1})
	q := dst2.Pix[4*(1*4+2):][:4]
	if q[0] == 0 || q[3] == 0 {
		t.Errorf("the shifted masked glyph = %v, want its blue by its coverage", q)
	}
	if q := dst2.Pix[4*(1*4+1):][:4]; q[0] != 0 || q[3] != 0 {
		t.Errorf("where the glyph is not = %v, want untouched", q)
	}

	// A glyph whose transform has no inverse takes nothing.
	s3 := &scene.Scene{
		Width:     4,
		Height:    4,
		MaskAtlas: masks,
		Glyphs: []scene.Glyph{
			{X: 0, Y: 0, W: 2, H: 2, U: 1, V: 1, UW: 2, VH: 2, Color: scene.Color{R: 0, G: 0, B: 255, A: 255}},
		},
	}
	dst3 := NewImage(4, 4)
	r3 := renderer{dst: dst3, s: s3, area: image.Rect(0, 0, 4, 4)}
	r3.updateBounds()
	r3.transformedGlyphs(&scene.Op{Kind: scene.OpGlyphs, Transform: scene.Matrix(0, 0, 0, 0, 0, 0), Start: 0, End: 1})
	for i := range dst3.Pix {
		if dst3.Pix[i] != 0 {
			t.Fatalf("a glyph without an inverse drew at byte %d", i)
		}
	}
}
