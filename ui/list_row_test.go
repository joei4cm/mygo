package ui

import "testing"

// A row reads the choice of its list, and whether the list holds the
// focus, as the row's own element would.
func TestListRowSelected(t *testing.T) {
	var s ListState
	sel := -1
	s.Selected = &sel
	var selected, focused [6]bool
	tt := NewTester(func(c *Context) {
		Column(c).Fill().Children(func() {
			List(c, &s, 6).Grow(1).Rows(func(r ListRow) {
				selected[r.Index], focused[r.Index] = r.Selected(), r.ListFocused()
				Textf(r.Context, "Row %d", r.Index).Height(24)
			})
			Button(c, "Elsewhere")
		})
	}, 300, 300)
	tt.Frame()
	for i := range selected {
		if selected[i] || focused[i] {
			t.Fatalf("nothing is chosen yet, row %d reads selected %v, list focus %v", i, selected[i], focused[i])
		}
	}

	// A click chooses the row, and the list takes the focus.
	tt.Click("Row 2")
	tt.Frame()
	for i := range selected {
		if selected[i] != (i == 2) {
			t.Errorf("after choosing Row 2, row %d reads selected %v", i, selected[i])
		}
		if !focused[i] {
			t.Errorf("after choosing Row 2, row %d reads no list focus", i)
		}
	}

	// A click outside leaves the rows without the list's focus.
	tt.Click("Elsewhere")
	tt.Frame()
	for i := range focused {
		if focused[i] {
			t.Errorf("a click outside, and row %d still reads a list focus", i)
		}
	}
	for i := range selected {
		if selected[i] != (i == 2) {
			t.Errorf("a click outside, and row %d reads selected %v", i, selected[i])
		}
	}
}
