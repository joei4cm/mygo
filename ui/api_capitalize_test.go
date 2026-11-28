package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// Capitalize tells the on-screen keyboard what to capitalize as the user types.
func TestCapitalizeAPI(t *testing.T) {
	value := ""
	kind := CapitalizeWords
	var got platform.CapitalizeKind
	tt := NewTester(func(c *Context) {
		in := TextInput(c, &value).Capitalize(kind)
		if !in.Valid() {
			t.Fatal("TextInput is invalid")
		}
		got = in.node().st.traits.Capitalize
	}, 200, 80)
	defer tt.rt.close()
	for _, tc := range []struct {
		kind platform.CapitalizeKind
		want platform.CapitalizeKind
	}{
		{platform.CapitalizeWords, platform.CapitalizeWords},
		{platform.CapitalizeSentences, platform.CapitalizeSentences},
		{platform.CapitalizeCharacters, platform.CapitalizeCharacters},
		{platform.CapitalizeNone, platform.CapitalizeNone},
		{platform.CapitalizeDefault, platform.CapitalizeDefault},
	} {
		kind = tc.kind
		tt.Frame()
		if got != tc.want {
			t.Errorf("Capitalize(%v) set %v, want %v", tc.kind, got, tc.want)
		}
	}
}

func TestCapitalizeAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Capitalize(CapitalizeWords).Valid() {
		t.Error("an absent element set capitalization")
	}
}
