package ui

import "testing"

// The public API of the space around and inside an element: the margins
// on both sides of an axis, the margin at the end of the lines, the
// padding on both sides of the vertical axis, and the least height.
func TestMarginsAPI(t *testing.T) {
	// The margins on both sides of the horizontal axis make room for
	// them: the next child of the row stands 20 further.
	var a, b Rect
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Row(c).Children(func() {
				a = Box(c).Size(40, 20).MarginX(10).Bounds()
				b = Box(c).Size(40, 20).Bounds()
			})
		})
	}, 300, 200)
	tt.Frame()
	if a.X != 10 || b.X != 60 {
		t.Errorf("the first child stands at %v and the next at %v, want 10 and 60", a.X, b.X)
	}

	// The margins on both sides of the vertical axis make room for them
	// in a column too.
	var ay, by Rect
	tt = NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			ay = Box(c).Size(40, 20).MarginY(10).Bounds()
			by = Box(c).Size(40, 20).Bounds()
		})
	}, 300, 200)
	tt.Frame()
	if ay.Y != 10 || by.Y != 40 {
		t.Errorf("the first child stands at %v and the next at %v, want 10 and 40", ay.Y, by.Y)
	}

	// The margin at the end of the lines sits on the right of a row that
	// runs left to right, and on the left of one that runs right to left.
	ltrA, ltrB, rtlA, rtlB := Rect{}, Rect{}, Rect{}, Rect{}
	tt = NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Row(c).Size(200, 20).Direction(LeftToRight).Children(func() {
				ltrA = Box(c).Size(40, 20).MarginEnd(10).Bounds()
				ltrB = Box(c).Size(40, 20).Bounds()
			})
			Row(c).Size(200, 20).Direction(RightToLeft).Children(func() {
				rtlA = Box(c).Size(40, 20).MarginEnd(10).Bounds()
				rtlB = Box(c).Size(40, 20).Bounds()
			})
		})
	}, 300, 200)
	tt.Frame()
	if ltrA.X != 0 || ltrB.X != 50 {
		t.Errorf("left to right, the first child stands at %v and the next at %v, want 0 and 50", ltrA.X, ltrB.X)
	}
	if rtlA.X != 160 || rtlB.X != 110 {
		t.Errorf("right to left, the first child stands at %v and the next at %v, want 160 and 110", rtlA.X, rtlB.X)
	}

	// The padding on both sides of the vertical axis holds a child away
	// from the top and the bottom, and gives the box its height.
	var parent, child Rect
	tt = NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			p := Box(c).PaddingY(10)
			p.Children(func() {
				child = Box(c).Size(40, 20).Bounds()
			})
			parent = p.Bounds()
		})
	}, 300, 200)
	tt.Frame()
	if child.Y != 10 {
		t.Errorf("the child of the padded box stands at %v, want 10", child.Y)
	}
	if parent.H != 40 {
		t.Errorf("the padded box is %v high, want 40", parent.H)
	}

	// The least height holds a smaller box up to it.
	var small Rect
	tt = NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			small = Box(c).Size(40, 30).MinHeight(80).Bounds()
		})
	}, 300, 200)
	tt.Frame()
	if small.H != 80 {
		t.Errorf("the box of the least height 80 is %v high", small.H)
	}

	// Outside a frame the view spaces nothing.
	if (Element{}).MarginX(10).Valid() {
		t.Error("a margin was set outside a frame")
	}
	if (Element{}).MarginY(10).Valid() {
		t.Error("a margin was set outside a frame")
	}
	if (Element{}).MarginEnd(10).Valid() {
		t.Error("a margin was set outside a frame")
	}
	if (Element{}).PaddingY(10).Valid() {
		t.Error("a padding was set outside a frame")
	}
	if (Element{}).MinHeight(80).Valid() {
		t.Error("a height was set outside a frame")
	}
}
