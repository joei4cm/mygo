package ui

import "testing"

// The public API of a text the user edits in place: a double click shows
// an input of it, and Enter keeps what was typed.
func TestEditableTextAPI(t *testing.T) {
	value := "notes.txt"
	var changed bool
	tt := NewTester(func(c *Context) {
		if EditableText(c, &value).Changed() {
			changed = true
		}
	}, 300, 200)

	// The text shows as it is.
	if !tt.HasText("notes.txt") {
		t.Fatalf("the texts are %q", tt.Texts())
	}

	// A double click shows an input of it, with the text before its
	// extension selected.
	r, ok := tt.Find("notes.txt")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", tt.Texts())
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if tt.HasText("notes.txt") {
		t.Fatalf("the text is still a text among %q", tt.Texts())
	}

	// What is typed replaces the text before the extension.
	tt.Type("memo")
	tt.Key(0, KeyEnter)
	if value != "memo.txt" {
		t.Fatalf("the value is %q", value)
	}
	if !changed {
		t.Fatal("the change was not reported")
	}
	if !tt.HasText("memo.txt") {
		t.Fatalf("the texts are %q", tt.Texts())
	}

	// Outside a frame the view builds no text.
	p := tt.rt.public
	if EditableText(p, &value).Valid() {
		t.Error("a text was built outside a frame")
	}
}

// The public API of a text the user edits in place: Escape goes back to
// the text as it was.
func TestEditableTextEscapeAPI(t *testing.T) {
	value := "notes.txt"
	var changed bool
	tt := NewTester(func(c *Context) {
		if EditableText(c, &value).Changed() {
			changed = true
		}
	}, 300, 200)

	r, ok := tt.Find("notes.txt")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", tt.Texts())
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if tt.HasText("notes.txt") {
		t.Fatalf("the text is still a text among %q", tt.Texts())
	}
	tt.Type("other")
	tt.Key(0, KeyEscape)
	if value != "notes.txt" {
		t.Fatalf("the value is %q", value)
	}
	if changed {
		t.Fatal("the change was reported")
	}
	if !tt.HasText("notes.txt") {
		t.Fatalf("the texts are %q", tt.Texts())
	}
}
