package ui

import (
	"testing"
)

// The public API of a switch, a toggle and a radio button: a click turns
// a switch and a toggle on and off, a radio button picks its value, and
// Changed reports it.
func TestSwitchToggleRadioAPI(t *testing.T) {
	var (
		on       bool
		toggled  bool
		pick     string
		switchCh bool
		toggleCh bool
		radioCh  bool
		swRole   Role
		tgRole   Role
		rdRole   Role
	)
	tt := NewTester(func(c *Context) {
		sw := SwitchBase(c, &on)
		swRole = c.rt.states[sw.node().ID()].role
		sw.Children(func() { Text(c, "Wi-Fi") })
		if sw.Changed() {
			switchCh = true
		}
		tg := ToggleBase(c, &toggled)
		tgRole = c.rt.states[tg.node().ID()].role
		tg.Children(func() { Text(c, "Bold") })
		if tg.Changed() {
			toggleCh = true
		}
		a := RadioBase(c, &pick, "Light")
		rdRole = c.rt.states[a.node().ID()].role
		a.Children(func() { Text(c, "Light") })
		if a.Changed() {
			radioCh = true
		}
		b := RadioBase(c, &pick, "Dark")
		b.Children(func() { Text(c, "Dark") })
	}, 400, 300)

	// Nothing changes before a click.
	tt.Frame()
	if on || toggled || pick != "" {
		t.Errorf("the switch is %v, the toggle %v, the radio %q", on, toggled, pick)
	}
	if switchCh || toggleCh || radioCh {
		t.Error("something changed before a click")
	}

	// Assistive technology sees a switch, a toggle button and a radio
	// button.
	if swRole != RoleSwitch || tgRole != RoleToggleButton || rdRole != RoleRadio {
		t.Errorf("the roles are %v, %v, %v", swRole, tgRole, rdRole)
	}

	// A click turns the switch on, and another turns it off.
	if err := tt.Click("Wi-Fi"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !on {
		t.Error("the switch did not turn on")
	}
	if err := tt.Click("Wi-Fi"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if on {
		t.Error("the switch did not turn off")
	}

	// A click turns the toggle on, and another turns it off.
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !toggled {
		t.Error("the toggle did not turn on")
	}
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if toggled {
		t.Error("the toggle did not turn off")
	}

	// A click picks a radio button, and another pick replaces it.
	if err := tt.Click("Dark"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if pick != "Dark" {
		t.Errorf("the picked value is %q", pick)
	}
	if err := tt.Click("Light"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if pick != "Light" {
		t.Errorf("the picked value is %q", pick)
	}

	// Changed reported every change.
	if !switchCh || !toggleCh || !radioCh {
		t.Errorf("changed: %v, %v, %v", switchCh, toggleCh, radioCh)
	}

	// Outside a frame the view builds none of them.
	p := tt.rt.public
	if SwitchBase(p, &on) != (Element{}) || ToggleBase(p, &toggled) != (Element{}) ||
		RadioBase(p, &pick, "Light") != (Element{}) {
		t.Error("something was built outside a frame")
	}
}
