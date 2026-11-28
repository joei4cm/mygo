package ui

import (
	"testing"
)

// The public API of an outline: roots at the top, the children of an open
// item below it, which its arrow opens and closes.
func TestOutlineAPI(t *testing.T) {
	tree := map[string][]string{
		"/": {"a", "c"},
		"a": {"a/1"},
	}
	var s OutlineState[string]
	tt := NewTester(func(c *Context) {
		Outline(c, &s, []string{"/"}, func(dir string) []string { return tree[dir] }, func(path string) {
			Text(c, path)
		})
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("/") || tt.HasText("a") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Opening the root shows its children.
	s.Open.Add("/")
	tt.Frame()
	if !tt.HasText("a") || !tt.HasText("c") || tt.HasText("a/1") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// The arrow of an item opens and closes it.
	a, ok := tt.Find("a")
	if !ok {
		t.Fatalf("the item shows nowhere among %q", tt.Texts())
	}
	tt.ClickAt(a.X-10, a.Y+a.H/2)
	tt.Frame()
	if !tt.HasText("a/1") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	tt.ClickAt(a.X-10, a.Y+a.H/2)
	tt.Frame()
	if tt.HasText("a/1") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Closing the root hides what it holds.
	s.Open.Remove("/")
	tt.Frame()
	if tt.HasText("a/1") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Outside a frame the view builds no outline.
	p := tt.rt.public
	if Outline(p, &s, nil, func(dir string) []string { return nil }, func(item string) {}) != (Element{}) {
		t.Error("an outline was built outside a frame")
	}
}
