package ui

import (
	"image"
	"testing"
)

// barGray counts the pixels of the image that are the gray of a scroll
// bar's thumb, darker than the background.
func barGray(img *image.RGBA, r Rect) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			c := img.RGBAAt(x, y)
			dg, db := int(c.R)-int(c.G), int(c.B)-int(c.R)
			if c.A == 255 && c.R < 200 && dg > -4 && dg < 4 && db > -4 && db < 4 {
				n++
			}
		}
	}
	return n
}

// barRed counts the red pixels of the image in the rectangle r.
func barRed(img *image.RGBA, r Rect) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			c := img.RGBAAt(x, y)
			if c.A > 200 && c.R > 200 && c.G < 80 && c.B < 80 {
				n++
			}
		}
	}
	return n
}

// The public API of a scroll container's scroll bars, and of scrolling an
// element into view.
func TestScrollbarsAPI(t *testing.T) {
	// scrollBox builds a window of a scroll container 140 wide and 100
	// high, of a box 200 by 300, whose scroll bars keep showing.
	scrollBox := func(set func(Element) Element) *Tester {
		return NewTester(func(c *Context) {
			sc := ScrollBoth(c).Size(140, 100).KeepScrollbars(true)
			if set != nil {
				sc = set(sc)
			}
			sc.Children(func() {
				Box(c).Size(200, 300)
			})
		}, 170, 130)
	}

	t.Run("Scrollbars", func(t *testing.T) {
		shown := scrollBox(nil)
		shown.Frame()
		band := Rect{X: 128, Y: 4, W: 16, H: 92}
		if n := barGray(shown.Image(), band); n == 0 {
			t.Error("the overflowing scroll container draws no scroll bar")
		}

		hidden := scrollBox(func(e Element) Element { return e.Scrollbars(false) })
		hidden.Frame()
		if n := barGray(hidden.Image(), band); n != 0 {
			t.Errorf("the scroll bar shows after Scrollbars(false), %d pixels", n)
		}
	})

	t.Run("ScrollbarAxes", func(t *testing.T) {
		vband := Rect{X: 128, Y: 4, W: 16, H: 92}
		hband := Rect{X: 4, Y: 88, W: 120, H: 16}
		both := scrollBox(func(e Element) Element { return e.ScrollbarAxes(true, true) })
		both.Frame()
		if n := barGray(both.Image(), vband); n == 0 {
			t.Error("the vertical bar is missing")
		}
		if n := barGray(both.Image(), hband); n == 0 {
			t.Error("the horizontal bar is missing")
		}

		vertical := scrollBox(func(e Element) Element {
			return e.ScrollbarAxes(false, true)
		})
		vertical.Frame()
		if n := barGray(vertical.Image(), vband); n == 0 {
			t.Error("the vertical bar is missing after ScrollbarAxes(false, true)")
		}
		if n := barGray(vertical.Image(), hband); n != 0 {
			t.Errorf("the horizontal bar shows after ScrollbarAxes(false, true), %d pixels", n)
		}
	})

	t.Run("ScrollbarInsets", func(t *testing.T) {
		inset := scrollBox(func(e Element) Element { return e.ScrollbarInsets(64, 0, 0) })
		inset.Frame()
		if n := barGray(inset.Image(), Rect{X: 128, Y: 4, W: 16, H: 58}); n != 0 {
			t.Errorf("the bar does not keep 64 below the top, %d pixels there", n)
		}
		if n := barGray(inset.Image(), Rect{X: 128, Y: 68, W: 16, H: 28}); n == 0 {
			t.Error("the bar draws nothing below the inset")
		}
	})

	t.Run("ScrollbarStyle", func(t *testing.T) {
		styled := scrollBox(func(e Element) Element {
			return e.ScrollbarStyle(ScrollbarStyle{
				ThumbWidth: 14,
				TrackWidth: 14,
				Square:     true,
				Flush:      true,
				Thumb:      RGB(255, 0, 0),
			})
		})
		styled.Frame()
		if n := barRed(styled.Image(), Rect{X: 118, Y: 4, W: 30, H: 92}); n == 0 {
			t.Error("the thumb is not of the color ScrollbarStyle set")
		}
		// The thumb is 14 DIPs wide against the right edge of the box.
		if n := barRed(styled.Image(), Rect{X: 118, Y: 10, W: 14, H: 1}); n == 0 {
			t.Errorf("the thumb does not reach 14 DIPs wide")
		}
		if n := barRed(styled.Image(), Rect{X: 110, Y: 10, W: 8, H: 1}); n != 0 {
			t.Errorf("the thumb is wider than ScrollbarStyle set")
		}
	})

	t.Run("ScrollIntoView", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Scroll(c).Size(100, 100).Children(func() {
				Column(c).Children(func() {
					Box(c).Size(100, 50)
					Box(c).Size(100, 100)
					far := Box(c).Size(100, 50).Background(RGB(255, 0, 0))
					Box(c).Size(100, 150)
					far.ScrollIntoView()
				})
			})
		}, 150, 130)
		tt.Frame()

		if n := barRed(tt.Image(), Rect{X: 0, Y: 50, W: 100, H: 50}); n == 0 {
			t.Error("ScrollIntoView does not bring the box into the container")
		}
	})

	t.Run("Outside a frame", func(t *testing.T) {
		var e Element
		for name, got := range map[string]Element{
			"ScrollIntoView":  e.ScrollIntoView(),
			"ScrollbarAxes":   e.ScrollbarAxes(true, true),
			"ScrollbarInsets": e.ScrollbarInsets(64, 0, 0),
			"ScrollbarStyle":  e.ScrollbarStyle(ScrollbarStyle{}),
			"Scrollbars":      e.Scrollbars(false),
			"KeepScrollbars":  e.KeepScrollbars(true),
		} {
			if got != (Element{}) {
				t.Errorf("%s outside a frame gave %v, want nothing", name, got)
			}
		}
	})
}
