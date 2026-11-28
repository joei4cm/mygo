package ui

import "testing"

// The public API of laying children out in a grid: a box becomes a grid,
// its first rows are tracks of the app's choosing or n rows of equal
// height, and a child goes to the row it names.
func TestGridAPI(t *testing.T) {
	// A grid of two rows: the first 30 high, the second the rest of the
	// 100, which a child without a size of its own fills, and a child of
	// its own choosing goes to the second row.
	var a, b, d Rect
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Row(c).Size(200, 100).Grid().Columns(2).RowTracks(Fixed(30), Fr(1)).Children(func() {
				a = Box(c).Size(10, 10).Bounds()
				b = Box(c).Size(10, 10).RowStart(2).Bounds()
				d = Box(c).RowStart(2).ColumnStart(2).Bounds()
			})
		})
	}, 300, 300)
	tt.Frame()
	if a.Y != 0 || a.H != 10 {
		t.Errorf("the first child takes %v, want its 10 high in the first row", a)
	}
	if b.Y != 30 || b.H != 10 {
		t.Errorf("the child of the second row takes %v, want its 10 high at 30", b)
	}
	if d != (Rect{X: 100, Y: 30, W: 100, H: 70}) {
		t.Errorf("the child without a size takes %v, want the second row's 70", d)
	}

	// GridRows gives the grid rows of equal height, which a child without
	// a size of its own fills.
	var c1, c2 Rect
	tt = NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Row(c).Size(200, 100).Grid().Columns(1).GridRows(2).Children(func() {
				c1 = Box(c).Size(10, 10).Bounds()
				c2 = Box(c).Bounds()
			})
		})
	}, 300, 300)
	tt.Frame()
	if c1.Y != 0 {
		t.Errorf("the first child of the even rows stands at %v, want 0", c1.Y)
	}
	if c2 != (Rect{X: 0, Y: 50, W: 200, H: 50}) {
		t.Errorf("the second child of the even rows takes %v, want the bottom half", c2)
	}

	// Outside a frame the view builds no grid.
	if (Element{}).Grid().Valid() {
		t.Error("a grid was built outside a frame")
	}
	if (Element{}).RowTracks(Fixed(30)).Valid() {
		t.Error("rows were set outside a frame")
	}
	if (Element{}).GridRows(2).Valid() {
		t.Error("rows were set outside a frame")
	}
	if (Element{}).RowStart(2).Valid() {
		t.Error("a row was chosen outside a frame")
	}
}
