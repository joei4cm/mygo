package ui

import "testing"

// SelectColumns chooses every row of a Table choosing cells and the
// columns from "A" to "C" as the block's anchor and cursor.
func TestSelectColumns(t *testing.T) {
	tt, s, sel := cellTable(t, 5, 4, 0)
	tt.Click("r1c2") // pick a column first, SelectColumns replaces the whole block
	wantBlock(t, s, sel, "click first", []int{1}, ColumnRange{Cursor: "B", Anchor: "B"}, []string{"B"})

	s.SelectColumns("A", "C")
	wantBlock(t, s, sel, "SelectColumns A to C", []int{0, 1, 2, 3, 4}, ColumnRange{Cursor: "C", Anchor: "A"}, []string{"A", "B", "C"})
}

// SelectColumns is a no-op for a table that does not choose cells.
func TestSelectColumnsWithoutCellSelect(t *testing.T) {
	s := &ListState{CellSelect: false}
	sel := &Selection[int]{}
	s.Selection = sel
	s.SelectColumns("A", "C")
	if s.Cells != (ColumnRange{}) || len(chosenRows(sel)) != 0 {
		t.Errorf("SelectColumns mutated state without CellSelect: cells %+v rows %v", s.Cells, chosenRows(sel))
	}
}

// SelectColumns is also a no-op before a frame exists.
func TestSelectColumnsWithoutFrame(t *testing.T) {
	s := &ListState{CellSelect: true}
	sel := &Selection[int]{}
	s.Selection = sel
	s.SelectColumns("A", "C")
	if s.Cells != (ColumnRange{}) || len(chosenRows(sel)) != 0 {
		t.Errorf("SelectColumns mutated state without a frame: cells %+v rows %v", s.Cells, chosenRows(sel))
	}
}

// SelectColumns updates a repeated call in the same view pass, choosing the
// new range instead of the old one.
func TestSelectColumnsIdempotent(t *testing.T) {
	_, s, sel := cellTable(t, 3, 4, 0)
	s.SelectColumns("A", "B")
	s.SelectColumns("C", "D")
	wantBlock(t, s, sel, "SelectColumns twice", []int{0, 1, 2}, ColumnRange{Cursor: "D", Anchor: "C"}, []string{"C", "D"})
}
