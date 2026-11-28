package ui

import "testing"

// InputMethod controls whether the focused text input composes text with
// the system's input method.
func TestInputMethodAPI(t *testing.T) {
	value := ""
	on := false
	tt := NewTester(func(c *Context) {
		TextInput(c, &value).InputMethod(on).Focus()
	}, 200, 80)
	defer tt.rt.close()
	for _, wantNoIME := range []bool{true, false} {
		on = !wantNoIME
		tt.Frame()
		active, noIME := tt.TextInputActive()
		if !active {
			t.Fatal("the focused TextInput did not activate text input")
		}
		if noIME != wantNoIME {
			t.Errorf("InputMethod(%v) gave NoInputMethod=%v, want %v", on, noIME, wantNoIME)
		}
	}
}

func TestInputMethodAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.InputMethod(false).Valid() {
		t.Error("an absent element changed its input method")
	}
}
