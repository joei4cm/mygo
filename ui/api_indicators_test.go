package ui

import (
	"testing"
)

// The public API of the small indicators: the meter shows its value, the
// avatar its initials, the badge its text, and the divider a line.
func TestIndicatorsAPI(t *testing.T) {
	var meterRole Role
	tt := NewTester(func(c *Context) {
		Column(c).Children(func() {
			m := Meter(c, 9, 0, 10, &MeterLevels{Warning: 5, Critical: 8}).Label("Disk")
			meterRole = c.rt.states[m.node().ID()].role
			Avatar(c, "Ada Lovelace", nil)
			Badge(c, "3")
			Divider(c)
			Text(c, "Body")
		})
	}, 400, 300)
	tt.Frame()
	if meterRole != RoleMeter {
		t.Errorf("the role is %v", meterRole)
	}
	if _, ok := tt.Access("Disk"); !ok {
		t.Error("the meter has no label")
	}
	if !tt.HasText("AL") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if !tt.HasText("3") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if _, ok := tt.Access("Ada Lovelace"); !ok {
		t.Error("the avatar has no label")
	}

	// Outside a frame the view builds none of them.
	p := tt.rt.public
	if Meter(p, 9, 0, 10, nil) != (Element{}) || Avatar(p, "Ada", nil) != (Element{}) ||
		Badge(p, "3") != (Element{}) || Divider(p) != (Element{}) {
		t.Error("something was built outside a frame")
	}
}
