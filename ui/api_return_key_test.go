package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// ReturnKey tells the text input which return key to show; ReturnNext also
// advances focus after submitting.
func TestReturnKeyAPI(t *testing.T) {
	value := ""
	kind := ReturnSearch
	var got platform.ReturnKeyKind
	tt := NewTester(func(c *Context) {
		in := TextInput(c, &value).ReturnKey(kind)
		if !in.Valid() {
			t.Fatal("TextInput is invalid")
		}
		got = in.node().st.traits.ReturnKey
	}, 200, 80)
	defer tt.rt.close()
	for _, tc := range []struct {
		kind ReturnKeyType
		want platform.ReturnKeyKind
	}{
		{ReturnDone, platform.ReturnDone},
		{ReturnGo, platform.ReturnGo},
		{ReturnSearch, platform.ReturnSearch},
		{ReturnSend, platform.ReturnSend},
		{ReturnDefault, platform.ReturnDefault},
	} {
		kind = tc.kind
		tt.Frame()
		if got != tc.want {
			t.Errorf("ReturnKey(%v) set %v, want %v", tc.kind, got, tc.want)
		}
	}
}

func TestReturnKeyNextFocusesNextAPI(t *testing.T) {
	value := ""
	tt := NewTester(func(c *Context) {
		Column(c).Gap(8).Children(func() {
			TextInput(c, &value).ReturnKey(ReturnNext)
			Button(c, "After")
		})
	}, 200, 100)
	defer tt.rt.close()
	tt.Frame()
	tt.Key(0, KeyTab)
	active, _ := tt.TextInputActive()
	if !active {
		t.Fatal("the text input did not receive focus")
	}
	tt.Key(0, KeyEnter)
	tt.Frame()
	if tt.rt.focused == 0 || tt.Focused("After") == false {
		t.Fatalf("ReturnNext did not move focus to the next button: focused=%d", tt.rt.focused)
	}
}

func TestReturnKeyAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.ReturnKey(ReturnSearch).Valid() {
		t.Error("an absent element set a return key")
	}
}
