package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// Autofill tells the on-screen keyboard what kind of text the input holds.
func TestAutofillAPI(t *testing.T) {
	value := ""
	kind := AutofillEmail
	var got platform.AutofillKind
	tt := NewTester(func(c *Context) {
		in := TextInput(c, &value).Autofill(kind)
		if !in.Valid() {
			t.Fatal("TextInput is invalid")
		}
		got = in.node().st.traits.Autofill
	}, 200, 80)
	defer tt.rt.close()
	tt.Frame()
	if got != platform.AutofillEmail {
		t.Fatalf("Autofill(Email) set %v, want Email", got)
	}
	kind = AutofillOneTimeCode
	tt.Frame()
	if got != platform.AutofillOneTimeCode {
		t.Fatalf("Autofill(OneTimeCode) set %v, want OneTimeCode", got)
	}
	kind = AutofillNone
	tt.Frame()
	if got != platform.AutofillNone {
		t.Fatalf("Autofill(None) set %v, want None", got)
	}
}

func TestAutofillAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Autofill(AutofillEmail).Valid() {
		t.Error("an absent element enabled autofill")
	}
}
