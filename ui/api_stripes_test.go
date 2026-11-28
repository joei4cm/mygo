package ui

import (
	"image"
	"testing"
)

// paintedCount counts the pixels of the color c in a box of a picture.
func paintedCount(img *image.RGBA, r Rect, c Color) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if p := img.RGBAAt(x, y); p.R == c.R && p.G == c.G && p.B == c.B && p.A == c.A {
				n++
			}
		}
	}
	return n
}

// The public API of what a box paints on its background: stripes to mark
// what is unavailable, and what a text paints behind its lines.
func TestStripesAPI(t *testing.T) {
	red := RGB(255, 0, 0)

	// A box of stripes paints stripes of the color, where a plain box
	// paints none.
	striped := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Box(c).Size(60, 60).Stripes(red, 10, 10, 0)
		})
	}, 300, 200)
	striped.Frame()
	if paintedCount(striped.Image(), Rect{X: 0, Y: 0, W: 60, H: 60}, red) == 0 {
		t.Error("the stripes painted none of the color")
	}

	plain := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Box(c).Size(60, 60)
		})
	}, 300, 200)
	plain.Frame()
	if paintedCount(plain.Image(), Rect{X: 0, Y: 0, W: 60, H: 60}, red) != 0 {
		t.Error("the plain box painted the color")
	}

	// A text with a background paints it behind its lines, where a plain
	// text paints none.
	back := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Text(c, "Hello").TextBackground(red)
		})
	}, 300, 200)
	back.Frame()
	r, ok := back.Find("Hello")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", back.Texts())
	}
	if paintedCount(back.Image(), r, red) == 0 {
		t.Error("the background of the text painted none of the color")
	}

	// Outside a frame the view paints nothing.
	if (Element{}).Stripes(red, 10, 10, 0).Valid() {
		t.Error("stripes were set outside a frame")
	}
	if (Element{}).TextBackground(red).Valid() {
		t.Error("a background was set outside a frame")
	}
}
