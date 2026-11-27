package ui

import (
	"testing"
	"time"
)

func TestEditorEmacsKey(t *testing.T) {
	now := func() time.Time { return time.Unix(0, 0) }
	newEd := func(s string, caret int) *editor {
		ed := newEditor(now, nil)
		ed.multiline = true
		ed.setText(s)
		ed.caret, ed.anchor = caret, caret
		ed.hasDesired = true
		return ed
	}

	// Control-A moves to the start of the current paragraph.
	ed := newEd("alpha\nbeta\ngamma", 8) // middle of "beta"
	if !ed.emacsKey(KeyA, false) {
		t.Fatal("KeyA was not handled")
	}
	if ed.caret != 6 || ed.anchor != 6 {
		t.Errorf("KeyA caret/anchor = %d/%d, want 6/6", ed.caret, ed.anchor)
	}
	if ed.hasDesired {
		t.Error("KeyA left hasDesired set")
	}

	// Shift+Control-E extends the selection to the end of the current paragraph.
	ed.hasDesired = true
	if !ed.emacsKey(KeyE, true) {
		t.Fatal("Shift+KeyE was not handled")
	}
	if ed.caret != 10 || ed.anchor != 6 {
		t.Errorf("Shift+KeyE caret/anchor = %d/%d, want 10/6", ed.caret, ed.anchor)
	}
	if ed.hasDesired {
		t.Error("Shift+KeyE left hasDesired set")
	}

	// Shift+Control-A extends back to the paragraph start; plain Control-E collapses at the end.
	ed = newEd("alpha\nbeta\ngamma", 8)
	if !ed.emacsKey(KeyA, true) {
		t.Fatal("Shift+KeyA was not handled")
	}
	if ed.caret != 6 || ed.anchor != 8 {
		t.Errorf("Shift+KeyA caret/anchor = %d/%d, want 6/8", ed.caret, ed.anchor)
	}
	if !ed.emacsKey(KeyE, false) {
		t.Fatal("KeyE was not handled")
	}
	if ed.caret != 10 || ed.anchor != 10 {
		t.Errorf("KeyE caret/anchor = %d/%d, want 10/10", ed.caret, ed.anchor)
	}

	// Control-K deletes from the caret to the end of the paragraph, ignoring any prior anchor.
	ed = newEd("alpha\nbeta\ngamma", 8)
	ed.anchor = 6
	if !ed.emacsKey(KeyK, false) {
		t.Fatal("KeyK in paragraph was not handled")
	}
	if got := ed.String(); got != "alpha\nbe\ngamma" {
		t.Errorf("KeyK in paragraph gave %q, want %q", got, "alpha\nbe\ngamma")
	}
	if ed.caret != 8 || ed.anchor != 8 {
		t.Errorf("KeyK caret/anchor = %d/%d, want 8/8", ed.caret, ed.anchor)
	}
	if ed.hasDesired {
		t.Error("KeyK left hasDesired set")
	}

	// At the end of a non-last paragraph, Control-K removes the line break and joins the next one.
	ed.hasDesired = true
	if !ed.emacsKey(KeyK, false) {
		t.Fatal("KeyK at paragraph end was not handled")
	}
	if got := ed.String(); got != "alpha\nbegamma" {
		t.Errorf("KeyK at paragraph end gave %q, want %q", got, "alpha\nbegamma")
	}
	if ed.caret != 8 || ed.anchor != 8 {
		t.Errorf("KeyK join caret/anchor = %d/%d, want 8/8", ed.caret, ed.anchor)
	}
	if ed.hasDesired {
		t.Error("KeyK join left hasDesired set")
	}

	// At the end of the final paragraph, Control-K is still handled and leaves the text intact.
	ed = newEd("alpha\nbeta", 10)
	if !ed.emacsKey(KeyK, false) {
		t.Fatal("KeyK at buffer end was not handled")
	}
	if got := ed.String(); got != "alpha\nbeta" {
		t.Errorf("KeyK at buffer end changed text to %q", got)
	}
	if ed.hasDesired {
		t.Error("KeyK at buffer end left hasDesired set")
	}

	// Unrelated keys are ignored and leave editor state untouched.
	ed = newEd("alpha\nbeta", 8)
	ed.anchor = 6
	if ed.emacsKey(KeyZ, false) {
		t.Error("KeyZ was unexpectedly handled")
	}
	if ed.caret != 8 || ed.anchor != 6 || !ed.hasDesired || ed.String() != "alpha\nbeta" {
		t.Errorf("KeyZ mutated editor state: caret=%d anchor=%d hasDesired=%v text=%q", ed.caret, ed.anchor, ed.hasDesired, ed.String())
	}
}
