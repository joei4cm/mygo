package ui

import "testing"

// The setters that take a side by its line's start or end place their
// value on the right or the left edge as the element's direction says.
func TestEndSettersLeftToRight(t *testing.T) {
	e := &node{}
	e.MarginEnd(5)
	if e.margin[1] != 5 || e.margin[3] != 0 {
		t.Errorf("MarginEnd left the margins at %v, want the right one 5", e.margin)
	}
	if e.logical&logMarginEnd == 0 {
		t.Error("MarginEnd did not note the end margin")
	}

	e = &node{}
	e.BorderEndWidth(2)
	if e.border[1] != 2 || e.border[3] != 0 {
		t.Errorf("BorderEndWidth left the borders at %v, want the right one 2", e.border)
	}
	if e.logical&logBorderEnd == 0 {
		t.Error("BorderEndWidth did not note the end border")
	}

	e = &node{}
	e.EndPercent(25)
	if e.inset[1] != percent(25) || e.inset[3] != (length{}) {
		t.Errorf("EndPercent left the insets at %v, want the right one 25%%", e.inset)
	}
	if e.logical&logInsetEnd == 0 {
		t.Error("EndPercent did not note the end inset")
	}
}

func TestEndSettersRightToLeft(t *testing.T) {
	e := &node{dir: dirRTL}
	e.MarginEnd(5)
	if e.margin[3] != 5 || e.margin[1] != 0 {
		t.Errorf("right to left, MarginEnd left the margins at %v, want the left one 5", e.margin)
	}
	if e.logical&logMarginEnd == 0 {
		t.Error("MarginEnd did not note the end margin")
	}

	e = &node{dir: dirRTL}
	e.BorderEndWidth(2)
	if e.border[3] != 2 || e.border[1] != 0 {
		t.Errorf("right to left, BorderEndWidth left the borders at %v, want the left one 2", e.border)
	}
	if e.logical&logBorderEnd == 0 {
		t.Error("BorderEndWidth did not note the end border")
	}

	e = &node{dir: dirRTL}
	e.EndPercent(25)
	if e.inset[3] != percent(25) || e.inset[1] != (length{}) {
		t.Errorf("right to left, EndPercent left the insets at %v, want the left one 25%%", e.inset)
	}
	if e.logical&logInsetEnd == 0 {
		t.Error("EndPercent did not note the end inset")
	}
}
