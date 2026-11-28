package ui

import (
	"testing"
)

// The public API of a find bar: its field takes the focus as it opens,
// typing changes the query, Enter and its buttons move among the matches,
// and Escape and Done close it.
func TestFindBarAPI(t *testing.T) {
	open := true
	query := "find"
	current := 0
	tt := NewTester(func(c *Context) {
		FindBar(c, &open, &query, 5, &current)
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("1 of 5") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// The field takes the focus as the bar opens, the query chosen: typing
	// replaces it.
	tt.Type("z")
	if query != "z" {
		t.Errorf("the query is %q", query)
	}

	// Enter goes to the next match.
	tt.Key(0, KeyEnter)
	tt.Frame()
	if current != 1 || !tt.HasText("2 of 5") {
		t.Errorf("at %d among %q", current, tt.Texts())
	}

	// Its buttons go to the next and the previous, round the ends.
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if current != 2 || !tt.HasText("3 of 5") {
		t.Errorf("at %d among %q", current, tt.Texts())
	}
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if current != 1 {
		t.Errorf("at %d", current)
	}
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if current != 4 || !tt.HasText("5 of 5") {
		t.Errorf("at %d among %q", current, tt.Texts())
	}

	// Done closes the bar.
	if err := tt.Click("Done"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if open {
		t.Error("the bar stays open")
	}

	// Escape closes the bar, which opens anew, its field with the focus.
	open = true
	tt.Frame()
	tt.Key(0, KeyEscape)
	tt.Frame()
	if open {
		t.Error("the bar stays open")
	}

	// Outside a frame the view builds no bar.
	p := tt.rt.public
	if FindBar(p, &open, &query, 5, &current) != (Element{}) {
		t.Error("a bar was built outside a frame")
	}
}
