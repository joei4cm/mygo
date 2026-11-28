package ui

import (
	"strings"
	"testing"
)

// The public API of a field's error: the message shows below the control
// in the theme's danger color, the control reads invalid with the message
// as its description, and its border draws in that color; clearing the
// error leaves all of it as it was.
func TestFieldErrorAPI(t *testing.T) {
	mail := "not an email"
	bad := true
	tt := NewTester(func(c *Context) {
		fld := Field(c, "Email", func() {
			TextInput(c, &mail).Width(200)
		})
		if bad {
			fld.Error("Enter an email address")
		}
	}, 300, 150)
	defer tt.rt.close()
	tt.Frame()
	if !tt.HasText("Enter an email address") {
		t.Fatal("the message does not show")
	}
	a, ok := tt.AccessByRole(RoleTextField, "Email")
	if !ok {
		t.Fatal("the control is not in the tree")
	}
	if !a.Invalid {
		t.Error("the control does not read invalid")
	}
	if !strings.Contains(a.Description, "Enter an email address") {
		t.Errorf("the control reads %q as its description", a.Description)
	}
	if reds := countDanger(tt, 300, 150); reds < 30 {
		t.Errorf("the danger color shows %d px", reds)
	}
	// A valid value leaves nothing behind.
	bad = false
	tt.Frame()
	if tt.HasText("Enter an email address") {
		t.Error("the message stays")
	}
	if a, _ := tt.AccessByRole(RoleTextField, "Email"); a.Invalid {
		t.Error("the control still reads invalid")
	}
	if reds := countDanger(tt, 300, 150); reds != 0 {
		t.Errorf("the danger color stays: %d px", reds)
	}
}

// countDanger counts the strong pixels of the theme's danger color.
func countDanger(tt *Tester, w, h int) int {
	reds := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
				reds++
			}
		}
	}
	return reds
}

// Outside a frame an error on an absent element changes nothing.
func TestFieldErrorAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Error("bad").Valid() {
		t.Error("an absent element took an error")
	}
}
