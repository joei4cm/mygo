package ui

import "testing"

// The public API of what stays while the rest scrolls: a heading that
// sticks to the top of the container around it holds there while its
// section scrolls under it, until the section itself goes.
func TestStickyAPI(t *testing.T) {
	red := RGB(220, 30, 30)
	sticky := false
	tt := NewTester(func(c *Context) {
		Scroll(c).Size(120, 100).Children(func() {
			Column(c).Children(func() {
				h := Box(c).Size(120, 20).Background(red)
				if sticky {
					h = h.Sticky().Top(0)
				}
				Box(c).Size(120, 60)
				Box(c).Size(120, 60)
			})
		})
	}, 240, 100)
	defer tt.rt.close()
	tt.Frame()
	reds := func() int {
		n := 0
		for y := 0; y < 20; y++ {
			for x := 0; x < 120; x++ {
				if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
					n++
				}
			}
		}
		return n
	}
	if reds() < 2000 {
		t.Fatal("the heading does not show at the top before scrolling")
	}
	// Scrolled, a heading that does not stick goes with its section; one
	// that does holds at the top while its section is still there.
	tt.Scroll(60, 50, 0, 100)
	tt.Frame()
	if n := reds(); sticky && n < 2000 {
		t.Errorf("the heading does not hold at the top: %d red px", n)
	} else if !sticky && n != 0 {
		t.Errorf("the heading went with its section, but %d red px hold", n)
	}
	// Asked to stick, the heading holds at the top at the same scroll.
	sticky = true
	tt.Frame()
	if n := reds(); n < 2000 {
		t.Errorf("the heading does not hold at the top: %d red px", n)
	}
	// Its room in the flow is still its own: scrolled back, it is back.
	tt.Scroll(60, 50, 0, -100)
	tt.Frame()
	if n := reds(); n < 2000 {
		t.Errorf("scrolled back, the heading shows %d red px", n)
	}
}

// Outside a frame nothing sticks.
func TestStickyAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Sticky().Valid() {
		t.Error("an absent element sticks")
	}
}
