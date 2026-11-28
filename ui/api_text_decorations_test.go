package ui

import (
	"image"
	"testing"
)

// decoBox shows a text decorated by style, and returns the box it takes
// and the pixels it draws.
func decoBox(t *testing.T, style func(e Element) Element) (Rect, *image.RGBA) {
	t.Helper()
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			e := Text(c, "Hello")
			if style != nil {
				e = style(e)
			}
		})
	}, 300, 200)
	r, ok := tt.Find("Hello")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", tt.Texts())
	}
	return r, tt.Image()
}

// inkOf counts the dark pixels of a box of a picture.
func inkOf(img *image.RGBA, r Rect) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if c := img.RGBAAt(x, y); c.R < 128 && c.A > 0 {
				n++
			}
		}
	}
	return n
}

// lastInk is the lowest row of dark pixels of a box of a picture.
func lastInk(img *image.RGBA, r Rect) int {
	last := -1
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if c := img.RGBAAt(x, y); c.R < 128 && c.A > 0 {
				last = y
			}
		}
	}
	return last
}

// sameInk reports whether two pictures look the same in a box.
func sameInk(a, b *image.RGBA, r Rect) bool {
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if a.RGBAAt(x, y) != b.RGBAAt(x, y) {
				return false
			}
		}
	}
	return true
}

// The public API of the decorations of a text: the line under it, in its
// shapes, colors and thickness, and how far below it goes.
func TestTextDecorationsAPI(t *testing.T) {
	plain, plainImg := decoBox(t, nil)

	// A wavy underline draws a line under the text.
	wavy, wavyImg := decoBox(t, func(e Element) Element { return e.WavyUnderline() })
	if inkOf(wavyImg, wavy) <= inkOf(plainImg, plain) {
		t.Error("the wavy underline draws no line")
	}

	// A wave is not a straight line.
	solid, solidImg := decoBox(t, func(e Element) Element { return e.Underline() })
	if sameInk(wavyImg, solidImg, wavy) {
		t.Error("the wavy underline draws a straight line")
	}

	// A dotted underline draws a line of dots, which is not a solid line.
	dotted, dottedImg := decoBox(t, func(e Element) Element { return e.DottedUnderline() })
	if inkOf(dottedImg, dotted) <= inkOf(plainImg, plain) {
		t.Error("the dotted underline draws no line")
	}
	if sameInk(dottedImg, solidImg, dotted) {
		t.Error("the dotted underline draws a solid line")
	}

	// The offset of the underline puts it lower.
	low, lowImg := decoBox(t, func(e Element) Element { return e.Underline().UnderlineOffset(4) })
	if lastInk(lowImg, low) <= lastInk(solidImg, solid) {
		t.Error("the offset moved the underline nowhere")
	}

	// The decoration has a color of its own.
	red, redImg := decoBox(t, func(e Element) Element { return e.Underline().DecorationColor(RGB(255, 0, 0)) })
	reds := 0
	for y := int(red.Y); y < int(red.Y+red.H); y++ {
		for x := int(red.X); x < int(red.X+red.W); x++ {
			if c := redImg.RGBAAt(x, y); c.R > 200 && c.G < 60 && c.A > 0 {
				reds++
			}
		}
	}
	if reds == 0 {
		t.Error("the underline took no color of its own")
	}

	// The decoration has a thickness of its own.
	thick, thickImg := decoBox(t, func(e Element) Element { return e.Underline().DecorationThickness(3) })
	if inkOf(thickImg, thick) <= inkOf(solidImg, solid) {
		t.Error("a thicker underline drew no more")
	}

	// Outside a frame the view decorates no text.
	if (Element{}).WavyUnderline().Valid() {
		t.Error("a text was decorated outside a frame")
	}
	if (Element{}).DottedUnderline().Valid() {
		t.Error("a text was decorated outside a frame")
	}
	if (Element{}).UnderlineOffset(4).Valid() {
		t.Error("a text was decorated outside a frame")
	}
	if (Element{}).DecorationColor(RGB(255, 0, 0)).Valid() {
		t.Error("a text was decorated outside a frame")
	}
	if (Element{}).DecorationThickness(3).Valid() {
		t.Error("a text was decorated outside a frame")
	}
}
