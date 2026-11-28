package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// Keyboard sends the on-screen keyboard layout requested for a text input.
func TestKeyboardAPI(t *testing.T) {
	value := ""
	kind := KeyboardEmail
	var got platform.KeyboardKind
	tt := NewTester(func(c *Context) {
		in := TextInput(c, &value).Keyboard(kind)
		if !in.Valid() {
			t.Fatal("TextInput is invalid")
		}
		got = in.node().st.traits.Keyboard
	}, 200, 80)
	defer tt.rt.close()
	for _, tc := range []struct {
		kind platform.KeyboardKind
		want platform.KeyboardKind
	}{
		{platform.KeyboardNumber, platform.KeyboardNumber},
		{platform.KeyboardDecimal, platform.KeyboardDecimal},
		{platform.KeyboardPhone, platform.KeyboardPhone},
		{platform.KeyboardURL, platform.KeyboardURL},
		{platform.KeyboardSearch, platform.KeyboardSearch},
		{platform.KeyboardDefault, platform.KeyboardDefault},
	} {
		kind = tc.kind
		tt.Frame()
		if got != tc.want {
			t.Errorf("Keyboard(%v) set %v, want %v", tc.kind, got, tc.want)
		}
	}
}

func TestKeyboardAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Keyboard(KeyboardEmail).Valid() {
		t.Error("an absent element set a keyboard")
	}
}
