package ui

import "testing"

// The public API of which way a box lays its children out: along the
// lines or down the page.
func TestDirectionAPI(t *testing.T) {
	// A column laid out as a row lays its children out from left to
	// right.
	var a, b Rect
	tt := NewTester(func(c *Context) {
		Column(c).Size(200, 20).Row().Children(func() {
			a = Box(c).Size(40, 10).Bounds()
			b = Box(c).Size(40, 10).Bounds()
		})
	}, 300, 200)
	tt.Frame()
	if a != (Rect{X: 0, Y: 5, W: 40, H: 10}) || b != (Rect{X: 40, Y: 5, W: 40, H: 10}) {
		t.Errorf("the children of the row stand at %v and %v, want side by side", a, b)
	}

	// A row laid out as a column lays its children out from top to
	// bottom.
	tt = NewTester(func(c *Context) {
		Row(c).Size(20, 200).Column().Children(func() {
			a = Box(c).Size(10, 40).Bounds()
			b = Box(c).Size(10, 40).Bounds()
		})
	}, 300, 200)
	tt.Frame()
	if a != (Rect{X: 5, Y: 0, W: 10, H: 40}) || b != (Rect{X: 5, Y: 40, W: 10, H: 40}) {
		t.Errorf("the children of the column stand at %v and %v, want one below the other", a, b)
	}

	// Outside a frame the view lays nothing out.
	if (Element{}).Row().Valid() {
		t.Error("a direction was set outside a frame")
	}
	if (Element{}).Column().Valid() {
		t.Error("a direction was set outside a frame")
	}
}
