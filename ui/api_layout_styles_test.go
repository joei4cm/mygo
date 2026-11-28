package ui

import "testing"

// The public API of the layout of a container: the space between its
// children, and the shape of an element.
func TestLayoutStylesAPI(t *testing.T) {
	// A horizontal gap between the children of a row.
	var first, second Rect
	row := NewTester(func(c *Context) {
		Row(c).GapX(20).Children(func() {
			first = Box(c).Size(50, 20).Bounds()
			second = Box(c).Size(50, 20).Bounds()
		})
	}, 200, 100)
	row.Frame()
	if first.X != 0 || second.X != 70 {
		t.Fatalf("the children are at %v and %v", first, second)
	}

	// A vertical gap between the children of a column.
	var top, bottom Rect
	column := NewTester(func(c *Context) {
		Column(c).GapY(20).Children(func() {
			top = Box(c).Size(50, 20).Bounds()
			bottom = Box(c).Size(50, 20).Bounds()
		})
	}, 200, 100)
	column.Frame()
	if top.Y != 0 || bottom.Y != 40 {
		t.Fatalf("the children are at %v and %v", top, bottom)
	}

	// An aspect ratio makes the height the width divided by it.
	var wide, tall Rect
	shape := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			wide = Box(c).Width(100).AspectRatio(2).Bounds()
			tall = Box(c).Height(40).AspectRatio(2).Bounds()
		})
	}, 200, 100)
	shape.Frame()
	if wide.W != 100 || wide.H != 50 {
		t.Fatalf("the element of a width is %v", wide)
	}
	if tall.W != 80 || tall.H != 40 {
		t.Fatalf("the element of a height is %v", tall)
	}

	// Outside a frame the view lays nothing out.
	if (Element{}).GapX(20).Valid() {
		t.Error("a gap was set outside a frame")
	}
	if (Element{}).GapY(20).Valid() {
		t.Error("a gap was set outside a frame")
	}
	if (Element{}).AspectRatio(2).Valid() {
		t.Error("a ratio was set outside a frame")
	}
}
