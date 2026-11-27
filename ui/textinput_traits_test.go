package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// The keyboard traits of a text input reach its state, and a stateless
// element takes none of them.
func TestTextInputTraits(t *testing.T) {
	value := ""
	var n *node
	tt := NewTester(func(c *Context) {
		el := TextInput(c, &value).Keyboard(KeyboardNumber).ReturnKey(ReturnNext).Capitalize(CapitalizeWords).Autocorrect(false).Autofill(AutofillEmail)
		n = el.node()
	}, 100, 40)
	tt.Frame()
	w := n.st.traits
	if w.Keyboard != KeyboardNumber || w.ReturnKey != ReturnNext || w.Capitalize != CapitalizeWords || w.Autocorrect != platform.AutocorrectOff || w.Autofill != AutofillEmail {
		t.Errorf("the traits are %+v", w)
	}
	// Back on, the keyboard corrects again.
	tt = NewTester(func(c *Context) {
		el := TextInput(c, &value).Autocorrect(true)
		n = el.node()
	}, 100, 40)
	tt.Frame()
	if n.st.traits.Autocorrect != platform.AutocorrectOn {
		t.Errorf("Autocorrect(true) made %v", n.st.traits.Autocorrect)
	}
	// A stateless element takes none: the calls return it for chaining,
	// and nothing reads a state it has not.
	var plain Element
	var plainValid bool
	tt = NewTester(func(c *Context) {
		plain = Box(c).Keyboard(KeyboardNumber).ReturnKey(ReturnNext).Capitalize(CapitalizeWords).Autocorrect(true).Autofill(AutofillEmail)
		plainValid = plain.Valid()
		// The traits of no element bind to nothing, and chain.
		if none := (Element{}).Keyboard(KeyboardNumber).ReturnKey(ReturnNext).Capitalize(CapitalizeWords).Autocorrect(true).Autofill(AutofillEmail); none.Valid() {
			t.Error("the traits of no element made one")
		}
	}, 100, 40)
	tt.Frame()
	if !plainValid {
		t.Error("the builders of a stateless element did not return it")
	}
}
