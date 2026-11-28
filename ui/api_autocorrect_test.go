package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// Autocorrect sets the trait sent with the text input to an on-screen
// keyboard, in either direction.
func TestAutocorrectAPI(t *testing.T) {
	value := "hello"
	on := false
	var got platform.AutocorrectKind
	tt := NewTester(func(c *Context) {
		in := TextInput(c, &value).Autocorrect(on)
		if !in.Valid() {
			t.Fatal("TextInput is invalid")
		}
		got = in.node().st.traits.Autocorrect
	}, 200, 80)
	defer tt.rt.close()
	tt.Frame()
	if got != platform.AutocorrectOff {
		t.Fatalf("Autocorrect(false) set %v, want off", got)
	}
	on = true
	tt.Frame()
	if got != platform.AutocorrectOn {
		t.Fatalf("Autocorrect(true) set %v, want on", got)
	}
}

func TestAutocorrectAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Autocorrect(true).Valid() {
		t.Error("an absent element enabled autocorrection")
	}
}
