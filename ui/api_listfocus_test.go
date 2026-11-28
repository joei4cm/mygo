package ui

import (
	"fmt"
	"testing"
)

// The public API of a list's focus: asking the list for the focus, and
// asking it whether it or one of its rows holds it.
func TestListFocusAPI(t *testing.T) {
	var s ListState
	var focused, within bool
	var rowWithin bool
	var ask, took bool
	show := true
	tt := NewTester(func(c *Context) {
		if show {
			List(c, &s, 3, func(i int) {
				Box(c).Height(24).Children(func() {
					Button(c, fmt.Sprintf("Row%d", i))
					rowWithin = s.FocusWithin(c)
				})
			})
			if ask {
				took = s.Focus(c)
				ask = false
			}
		}
		focused = s.Focused(c)
		within = s.FocusWithin(c)
	}, 200, 150)
	tt.Frame()
	if focused || within {
		t.Fatalf("a list without the focus: focused %v, within %v", focused, within)
	}
	// The list takes the focus, and says so in the same pass.
	ask = true
	tt.Frame()
	if ask {
		t.Fatal("the list did not take the focus")
	}
	if !took {
		t.Fatal("the list did not report that it took the focus")
	}
	if !focused || !within {
		t.Fatalf("the focused list: focused %v, within %v", focused, within)
	}
	if !rowWithin {
		t.Error("a row builder did not see the focus within the list")
	}
	// The focus stays with the list in the next pass.
	tt.Frame()
	if !focused {
		t.Error("the list lost the focus")
	}
	// A row of the list takes the focus: the list is within it, not focused.
	if err := tt.Click("Row1"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if focused {
		t.Error("the list was focused while one of its rows was")
	}
	if !within {
		t.Error("the list did not see the focus within it")
	}
	if !tt.Focused("Row1") {
		t.Error("the row did not take the focus")
	}
	// A list that is not built in this pass has no focus to give or to report.
	show = false
	tt.Frame()
	if focused || within {
		t.Fatalf("a list that is not built: focused %v, within %v", focused, within)
	}
	if s.Focus(nil) {
		t.Error("a list outside a frame took the focus")
	}
	if s.Focused(nil) || s.FocusWithin(nil) {
		t.Error("a list outside a frame reported the focus")
	}
}
