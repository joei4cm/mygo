package ui

import "testing"

// The public API of where the lines of a wrapping row or column stand:
// at its start, its middle or its end, and the other way for lines that
// wrap the other way.
func TestWrapAPI(t *testing.T) {
	// Four children of 100 by 10 wrap onto two lines of a row of 200 by
	// 60: the first two on one line, the last two on the next.
	boxes := func(style func(Element) Element) [4]Rect {
		var rs [4]Rect
		tt := NewTester(func(c *Context) {
			Column(c).AlignItems(Start).Children(func() {
				row := Row(c).Size(200, 60).Wrap()
				if style != nil {
					row = style(row)
				}
				row.Children(func() {
					rs[0] = Box(c).Size(100, 10).Bounds()
					rs[1] = Box(c).Size(100, 10).Bounds()
					rs[2] = Box(c).Size(100, 10).Bounds()
					rs[3] = Box(c).Size(100, 10).Bounds()
				})
			})
		}, 300, 200)
		tt.Frame()
		return rs
	}

	// The lines stand at the top.
	if rs := boxes(nil); rs[0].Y != 0 || rs[2].Y != 10 {
		t.Errorf("the lines of the row stand at %v and %v, want at 0 and 10", rs[0].Y, rs[2].Y)
	}

	// They stand in its middle.
	if rs := boxes(func(e Element) Element { return e.AlignContent(Center) }); rs[0].Y != 20 || rs[2].Y != 30 {
		t.Errorf("the lines of the row stand at %v and %v, want at 20 and 30", rs[0].Y, rs[2].Y)
	}

	// They stand at its bottom.
	if rs := boxes(func(e Element) Element { return e.AlignContent(End) }); rs[0].Y != 40 || rs[2].Y != 50 {
		t.Errorf("the lines of the row stand at %v and %v, want at 40 and 50", rs[0].Y, rs[2].Y)
	}

	// Wrapping the other way, the first line stands at the bottom.
	if rs := boxes(func(e Element) Element { return e.WrapReverse() }); rs[0].Y != 50 || rs[2].Y != 40 {
		t.Errorf("the lines of the row stand at %v and %v, want the first at 50 and the next at 40", rs[0].Y, rs[2].Y)
	}

	// In a column the lines are columns: wrapping the other way, the
	// first stands at the right.
	var cs [4]Rect
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Column(c).Size(60, 200).Wrap().WrapReverse().Children(func() {
				cs[0] = Box(c).Size(10, 100).Bounds()
				cs[1] = Box(c).Size(10, 100).Bounds()
				cs[2] = Box(c).Size(10, 100).Bounds()
				cs[3] = Box(c).Size(10, 100).Bounds()
			})
		})
	}, 200, 300)
	tt.Frame()
	if cs[0].X != 50 || cs[2].X != 40 {
		t.Errorf("the columns of the column stand at %v and %v, want the first at 50 and the next at 40", cs[0].X, cs[2].X)
	}

	// Outside a frame the view places no line.
	if (Element{}).WrapReverse().Valid() {
		t.Error("a container wrapped outside a frame")
	}
	if (Element{}).AlignContent(End).Valid() {
		t.Error("lines were placed outside a frame")
	}
}
