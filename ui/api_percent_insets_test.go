package ui

import "testing"

// The public API of positioning an element by a percentage of its parent:
// physical edges and logical start and end edges.
func TestPercentInsetsAPI(t *testing.T) {
	var tl, br, ltrStart, rtlStart, ltrEnd, rtlEnd Rect
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Box(c).Size(200, 100).Children(func() {
				tl = Box(c).Absolute().TopPercent(20).LeftPercent(10).Size(40, 20).Bounds()
				br = Box(c).Absolute().BottomPercent(20).RightPercent(10).Size(40, 20).Bounds()
			})
			Box(c).Size(200, 100).Direction(LeftToRight).Children(func() {
				ltrStart = Box(c).Absolute().Top(0).StartPercent(10).Size(40, 20).Bounds()
				ltrEnd = Box(c).Absolute().Top(30).EndPercent(10).Size(40, 20).Bounds()
			})
			Box(c).Size(200, 100).Direction(RightToLeft).Children(func() {
				rtlStart = Box(c).Absolute().Top(0).StartPercent(10).Size(40, 20).Bounds()
				rtlEnd = Box(c).Absolute().Top(30).EndPercent(10).Size(40, 20).Bounds()
			})
		})
	}, 300, 400)
	tt.Frame()

	// TopPercent(20) on 100px = y 20; LeftPercent(10) on 200px = x 20.
	if tl.X != 20 || tl.Y != 20 {
		t.Errorf("TopPercent/LeftPercent gave %v, want at (20,20)", tl)
	}
	// BottomPercent(20) = 100 - 20 - 20 = 60; RightPercent(10) = 200 - 20 - 40 = 140.
	if br.X != 140 || br.Y != 60 {
		t.Errorf("BottomPercent/RightPercent gave %v, want at (140,60)", br)
	}
	// Logical start/end follow direction: LTR puts start at left (20) and end at right (140).
	if ltrStart.X != 20 || ltrEnd.X != 140 {
		t.Errorf("LTR StartPercent/EndPercent gave %v / %v, want x=20 / x=140", ltrStart, ltrEnd)
	}
	// RTL flips start to right (140) and end to left (20).
	if rtlStart.X != 140 || rtlEnd.X != 20 {
		t.Errorf("RTL StartPercent/EndPercent gave %v / %v, want x=140 / x=20", rtlStart, rtlEnd)
	}

	// Outside a frame the view positions nothing.
	if (Element{}).TopPercent(20).Valid() {
		t.Error("TopPercent built outside a frame")
	}
	if (Element{}).RightPercent(10).Valid() {
		t.Error("RightPercent built outside a frame")
	}
	if (Element{}).BottomPercent(20).Valid() {
		t.Error("BottomPercent built outside a frame")
	}
	if (Element{}).LeftPercent(10).Valid() {
		t.Error("LeftPercent built outside a frame")
	}
	if (Element{}).StartPercent(10).Valid() {
		t.Error("StartPercent built outside a frame")
	}
	if (Element{}).EndPercent(10).Valid() {
		t.Error("EndPercent built outside a frame")
	}
}
