package ui

import (
	"fmt"
	"testing"
)

// The public API of a list's focus: the view gives the list the keyboard
// focus, and tells whether the list, or one of its rows, has it.
func TestListFocusAPI(t *testing.T) {
	var s ListState
	var (
		show     = true
		focus    bool
		focused  bool
		within   bool
		after    bool
		hasFocus bool
	)
	tt := NewTester(func(c *Context) {
		if !show {
			// A list that is not built in this pass has no focus.
			hasFocus = s.Focus(c)
			hasFocus = hasFocus || s.Focused(c)
			hasFocus = hasFocus || s.FocusWithin(c)
			return
		}
		List(c, &s, 3, func(i int) {
			Text(c, fmt.Sprintf("Row %d", i))
			focus = s.Focus(c)
			focused = s.Focused(c)
			within = s.FocusWithin(c)
		})
		after = s.Focused(c)
	}, 300, 300)
	tt.Frame()
	if !focus {
		t.Error("the list was not built in this pass")
	}
	if !focused {
		t.Error("the list has not the focus")
	}
	if !within {
		t.Error("the list has not the focus within it")
	}
	if !after {
		t.Error("the list has not the focus after it is built")
	}

	// A list that is not built in this pass has no focus.
	show = false
	tt.Frame()
	if hasFocus {
		t.Error("a list that is not built has the focus")
	}

	// Outside a frame the focus is nowhere.
	p := tt.rt.public
	if s.Focus(p) || s.Focused(p) || s.FocusWithin(p) {
		t.Error("the list has the focus outside a frame")
	}
}
