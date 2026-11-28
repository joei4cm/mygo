package ui

import (
	"image"
	"testing"
)

// extentOf returns the box the pixels of the color c fill in a picture,
// or nothing where there are none.
func extentOf(img *image.RGBA, c Color) Rect {
	lo, hi := image.Pt(1<<30, 1<<30), image.Pt(-1, -1)
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if p := img.RGBAAt(x, y); p.R == c.R && p.G == c.G && p.B == c.B && p.A == c.A {
				if x < lo.X {
					lo.X = x
				}
				if y < lo.Y {
					lo.Y = y
				}
				if x > hi.X {
					hi.X = x
				}
				if y > hi.Y {
					hi.Y = y
				}
			}
		}
	}
	if hi.X < 0 {
		return Rect{}
	}
	return Rect{X: float32(lo.X), Y: float32(lo.Y), W: float32(hi.X - lo.X + 1), H: float32(hi.Y - lo.Y + 1)}
}

// The public API of how an element is drawn askew: turned around, turned
// about the axes, moved by a share of its size, and gone from view when
// its back faces it.
func TestTransformAPI(t *testing.T) {
	red := RGB(255, 0, 0)
	// count lays out a red box of 60 by 20 at 70, 70, styled by style,
	// and reports how many red pixels it drew and the box they fill.
	count := func(style func(Element) Element) (int, Rect) {
		tt := NewTester(func(c *Context) {
			Column(c).Padding(70).Children(func() {
				e := Box(c).Size(60, 20).Background(red)
				if style != nil {
					e = style(e)
				}
			})
		}, 200, 200)
		tt.Frame()
		img := tt.Image()
		n := 0
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				if p := img.RGBAAt(x, y); p.R == red.R && p.G == red.G && p.B == red.B && p.A == red.A {
					n++
				}
			}
		}
		return n, extentOf(img, red)
	}

	// The plain box fills its 60 by 20.
	n, r := count(nil)
	if n != 1200 || r != (Rect{X: 70, Y: 70, W: 60, H: 20}) {
		t.Fatalf("the plain box drew %d pixels in %v", n, r)
	}

	// Turned a quarter way round, the box stands on its end.
	if _, r = count(func(e Element) Element { return e.Rotate(90) }); r != (Rect{X: 90, Y: 50, W: 20, H: 60}) {
		t.Errorf("the box turned a quarter round fills %v, want it on its end", r)
	}

	// The same turn, given as a transform of the app's own.
	if _, r = count(func(e Element) Element { return e.Transform(Rotation(90)) }); r != (Rect{X: 90, Y: 50, W: 20, H: 60}) {
		t.Errorf("the box of a quarter turn fills %v, want it on its end", r)
	}

	// Turned about the vertical or the horizontal axis a quarter way, the
	// box stands edge on, and draws nothing.
	if n, _ = count(func(e Element) Element { return e.RotateY(90) }); n != 0 {
		t.Errorf("the box turned about the vertical axis drew %d pixels", n)
	}
	if n, _ = count(func(e Element) Element { return e.RotateX(90) }); n != 0 {
		t.Errorf("the box turned about the horizontal axis drew %d pixels", n)
	}

	// Turned halfway about the vertical axis it still draws; hidden from
	// the back, it does not.
	if n, _ = count(func(e Element) Element { return e.RotateY(180) }); n != 1200 {
		t.Errorf("the box turned halfway round drew %d pixels", n)
	}
	if n, _ = count(func(e Element) Element { return e.RotateY(180).BackfaceHidden() }); n != 0 {
		t.Errorf("the box hidden from the back drew %d pixels", n)
	}

	// Moved by a half of its width, the box stands 30 further right.
	if _, r = count(func(e Element) Element { return e.TranslatePercent(50, 0) }); r != (Rect{X: 100, Y: 70, W: 60, H: 20}) {
		t.Errorf("the box moved by half its width fills %v, want it at 100", r)
	}

	// A box brought towards the viewer under a perspective draws larger.
	n, r = count(func(e Element) Element { return e.TranslateZ(100).Perspective(400) })
	if n <= 1200 || r.W <= 60 || r.H <= 20 {
		t.Errorf("the box brought towards the viewer drew %d pixels in %v, want it larger", n, r)
	}

	// Outside a frame the view draws nothing askew.
	if (Element{}).Rotate(90).Valid() {
		t.Error("a box was turned outside a frame")
	}
	if (Element{}).RotateX(90).Valid() || (Element{}).RotateY(90).Valid() {
		t.Error("a box was turned outside a frame")
	}
	if (Element{}).TranslatePercent(50, 0).Valid() {
		t.Error("a box was moved outside a frame")
	}
	if (Element{}).TranslateZ(100).Valid() || (Element{}).Perspective(400).Valid() {
		t.Error("a box was drawn in perspective outside a frame")
	}
	if (Element{}).BackfaceHidden().Valid() {
		t.Error("a box was hidden from the back outside a frame")
	}
	if (Element{}).Transform(Rotation(90)).Valid() {
		t.Error("a transform was set outside a frame")
	}
}
