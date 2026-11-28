package ui

import (
	"testing"
)

// The public API of a group of check boxes: the box over them is checked
// while they all are, and a click on it checks them all, or none.
func TestCheckboxGroupAPI(t *testing.T) {
	mail, calendar := true, false
	tt := NewTester(func(c *Context) {
		CheckboxGroup(c, "Notifications", func() {
			Checkbox(c, &mail, "Mail")
			Checkbox(c, &calendar, "Calendar")
		})
	}, 400, 300)
	tt.Frame()
	if !mail || calendar {
		t.Errorf("the boxes are %v and %v", mail, calendar)
	}
	if !tt.HasText("Notifications") || !tt.HasText("Mail") || !tt.HasText("Calendar") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A click on the box over them checks them all.
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !mail || !calendar {
		t.Errorf("the boxes are %v and %v", mail, calendar)
	}

	// Another click checks none of them.
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if mail || calendar {
		t.Errorf("the boxes are %v and %v", mail, calendar)
	}

	// Outside a frame the view builds no group.
	p := tt.rt.public
	if CheckboxGroup(p, "Notifications", func() {}) != (Element{}) {
		t.Error("a group was built outside a frame")
	}
}

// The public API of a check box that can be mixed: a click checks it,
// another unchecks it.
func TestCheckboxStateAPI(t *testing.T) {
	state := CheckMixed
	tt := NewTester(func(c *Context) {
		CheckboxState(c, &state, "Select all")
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("Select all") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if err := tt.Click("Select all"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state != CheckOn {
		t.Errorf("the state is %v", state)
	}
	if err := tt.Click("Select all"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state != CheckOff {
		t.Errorf("the state is %v", state)
	}

	// Outside a frame the view builds no check box.
	p := tt.rt.public
	if CheckboxState(p, &state, "Select all") != (Element{}) {
		t.Error("a check box was built outside a frame")
	}
}

// The public API of a check box that can be mixed, without a look: a
// click checks it, another unchecks it.
func TestCheckboxStateBaseAPI(t *testing.T) {
	state := CheckOff
	tt := NewTester(func(c *Context) {
		cb := CheckboxStateBase(c, &state).Label("Select all")
		cb.Children(func() { Text(c, "All") })
	}, 400, 300)
	tt.Frame()
	if err := tt.Click("All"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state != CheckOn {
		t.Errorf("the state is %v", state)
	}
	if err := tt.Click("All"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state != CheckOff {
		t.Errorf("the state is %v", state)
	}

	// Outside a frame the view builds no check box.
	p := tt.rt.public
	if CheckboxStateBase(p, &state) != (Element{}) {
		t.Error("a check box was built outside a frame")
	}
}
