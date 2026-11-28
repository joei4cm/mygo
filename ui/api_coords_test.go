package ui

import "testing"

// The public API of reading an element's layout box before transforms and
// mapping points between the element and the window.
func TestCoordsAPI(t *testing.T) {
	var layout, visual Rect
	var wx, wy, lx, ly float32
	var ok, singularOK bool

	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Padding(30, 0, 0, 40).Children(func() {
			moved := Box(c).Size(80, 20).Translate(10, 5)
			layout = moved.LayoutBounds()
			visual = moved.Bounds()
			wx, wy = moved.LocalToWindow(5, 4)
			lx, ly, ok = moved.WindowToLocal(wx, wy)

			flat := Box(c).Size(40, 20).Scale(0, 0)
			_, _, singularOK = flat.WindowToLocal(50, 60)
		})
	}, 300, 200)
	tt.Frame()

	// LayoutBounds keeps the box before the translation, where Bounds moves with it.
	if layout != (Rect{X: 40, Y: 30, W: 80, H: 20}) {
		t.Errorf("LayoutBounds gave %v, want {40 30 80 20}", layout)
	}
	if visual != (Rect{X: 50, Y: 35, W: 80, H: 20}) {
		t.Errorf("Bounds gave %v, want {50 35 80 20}", visual)
	}

	// LocalToWindow adds the element's position and its transform; WindowToLocal inverts it.
	if wx != 55 || wy != 39 {
		t.Errorf("LocalToWindow(5, 4) gave (%v, %v), want (55, 39)", wx, wy)
	}
	if !ok || lx != 5 || ly != 4 {
		t.Errorf("WindowToLocal(%v, %v) gave (%v, %v, %v), want (5, 4, true)", wx, wy, lx, ly, ok)
	}

	// A singular transform has no inverse.
	if singularOK {
		t.Error("WindowToLocal succeeded through a singular Scale(0) transform")
	}

	// Outside a frame the element reports empty bounds and maps nothing.
	if got := (Element{}).LayoutBounds(); got != (Rect{}) {
		t.Errorf("LayoutBounds outside a frame gave %v", got)
	}
	if x, y := (Element{}).LocalToWindow(5, 4); x != 0 || y != 0 {
		t.Errorf("LocalToWindow outside a frame gave (%v, %v)", x, y)
	}
	if x, y, valid := (Element{}).WindowToLocal(55, 39); x != 0 || y != 0 || valid {
		t.Errorf("WindowToLocal outside a frame gave (%v, %v, %v)", x, y, valid)
	}
}
