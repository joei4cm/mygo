package ui

import "testing"

// SetChanged lets a custom widget report its value changing through the
// same Changed API as built-in controls.
func TestSetChangedAPI(t *testing.T) {
	var changed, sawChange bool
	tt := NewTester(func(c *Context) {
		Column(c).Children(func() {
			widget := Box(c).Key("custom")
			if Button(c, "Change").Clicked() {
				widget.SetChanged()
			}
			changed = widget.Changed()
			sawChange = sawChange || changed
		})
	}, 200, 100)
	defer tt.rt.close()
	tt.Frame()
	if changed {
		t.Fatal("the custom widget reported a change before input")
	}
	if err := tt.Click("Change"); err != nil {
		t.Fatal(err)
	}
	if !sawChange {
		t.Fatal("SetChanged did not make Changed report the update")
	}
	tt.Frame()
	if changed {
		t.Fatal("the custom widget kept reporting the same change")
	}
}

func TestSetChangedAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.SetChanged().Valid() {
		t.Error("an absent element reported a change")
	}
	if absent.Changed() {
		t.Error("an absent element was changed")
	}
}
