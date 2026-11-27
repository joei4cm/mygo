package ui

import "testing"

// Rows gives the grid n rows of equal height, Fr(1) each, replacing
// whatever tracks it had.
func TestGridRows(t *testing.T) {
	e := &node{}
	e.RowTracks(Fixed(50), Fixed(60))
	if len(e.rows) != 2 || e.rows[0] != (Fixed(50)) || e.rows[1] != (Fixed(60)) {
		t.Fatalf("RowTracks left %v, want two fixed tracks", e.rows)
	}

	if got := e.Rows(3); got != e {
		t.Error("Rows did not return the node")
	}
	if len(e.rows) != 3 {
		t.Fatalf("Rows(3) left %d rows, want 3", len(e.rows))
	}
	for k, tr := range e.rows {
		if tr != (Fr(1)) {
			t.Errorf("row %d is %v, want Fr(1)", k, tr)
		}
	}

	// Rows(0) clears the rows.
	e.Rows(0)
	if len(e.rows) != 0 {
		t.Errorf("Rows(0) left %d rows, want none", len(e.rows))
	}
}

// BottomPercent places an absolute element by a percentage of its
// parent's height from the bottom.
func TestBottomPercent(t *testing.T) {
	e := &node{}
	if got := e.BottomPercent(25); got != e {
		t.Error("BottomPercent did not return the node")
	}
	if e.inset[2] != (percent(25)) {
		t.Errorf("BottomPercent left the bottom inset at %v, want a quarter", e.inset[2])
	}
	if !e.insetSet {
		t.Error("BottomPercent did not note an inset set")
	}
}

// KeepFocus keeps the keyboard focus on an element when another takes it.
func TestKeepFocus(t *testing.T) {
	e := &node{}
	if got := e.KeepFocus(); got != e {
		t.Error("KeepFocus did not return the node")
	}
	if e.flags&flagKeepFocus == 0 {
		t.Error("KeepFocus did not set its flag")
	}
}
