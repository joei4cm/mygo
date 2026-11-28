package ui

import (
	"image"
	"testing"
)

// redPixels counts the red pixels in a box of the picture.
func redPixels(img *image.RGBA, r Rect) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			p := img.RGBAAt(x, y)
			if p.R == 255 && p.G == 0 && p.B == 0 && p.A == 255 {
				n++
			}
		}
	}
	return n
}

// The public API of drawing a border: its width, whether it is dashed,
// and the width at the end side, which follows the direction.
func TestBorderAPI(t *testing.T) {
	red := RGB(255, 0, 0)

	// A four-pixel solid border draws more of the color than a one-pixel
	// border around the same box.
	borderPixels := func(width float32, style BorderStyle) int {
		tt := NewTester(func(c *Context) {
			Column(c).AlignItems(Start).Children(func() {
				Box(c).Size(60, 40).BorderWidth(width).BorderColor(red).BorderStyle(style)
			})
		}, 100, 80)
		tt.Frame()
		return redPixels(tt.Image(), Rect{W: 100, H: 80})
	}
	thin := borderPixels(1, BorderSolid)
	thick := borderPixels(4, BorderSolid)
	if thin <= 0 || thick <= thin {
		t.Errorf("solid borders of widths 1 and 4 draw %d and %d red pixels", thin, thick)
	}

	// A dashed border draws less of the color than the same solid border,
	// but does draw some.
	dashed := borderPixels(4, BorderDashed)
	if dashed <= 0 || dashed >= thick {
		t.Errorf("the dashed border draws %d red pixels; the solid draws %d", dashed, thick)
	}

	// BorderEndWidth takes room on the right in LTR and on the left in
	// RTL, leaving the same 90-pixel content width.
	var ltr, rtl Rect
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Box(c).Size(100, 40).Direction(LeftToRight).BorderEndWidth(10).Children(func() {
				ltr = Box(c).Fill().Bounds()
			})
			Box(c).Size(100, 40).Direction(RightToLeft).BorderEndWidth(10).Children(func() {
				rtl = Box(c).Fill().Bounds()
			})
		})
	}, 200, 120)
	tt.Frame()
	if ltr != (Rect{X: 0, Y: 0, W: 90, H: 40}) {
		t.Errorf("LTR border-end leaves content %v, want 90 wide at the left", ltr)
	}
	if rtl != (Rect{X: 10, Y: 40, W: 90, H: 40}) {
		t.Errorf("RTL border-end leaves content %v, want 90 wide at the right", rtl)
	}

	// Outside a frame the view draws no border.
	if (Element{}).BorderWidth(4).Valid() {
		t.Error("a border width was set outside a frame")
	}
	if (Element{}).BorderStyle(BorderDashed).Valid() {
		t.Error("a border style was set outside a frame")
	}
	if (Element{}).BorderEndWidth(10).Valid() {
		t.Error("a border end width was set outside a frame")
	}
}
