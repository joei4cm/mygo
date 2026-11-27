package ui

import (
	"testing"
)

// The public API of an alert dialog: it shows while it is open, a click
// on a button chooses it and closes it, and Escape dismisses it.
func TestAlertDialogAPI(t *testing.T) {
	var (
		open    = true
		chosen  = -1
		alerted bool
	)
	tt := NewTester(func(c *Context) {
		if i := AlertDialog(c, &open, "Delete Notes?", "You can't undo this.", "Cancel", "Delete"); i != -1 {
			chosen = i
			alerted = true
		}
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("Delete Notes?") || !tt.HasText("You can't undo this.") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if alerted {
		t.Error("the alert was answered before a click")
	}

	// A click on a button chooses it, and closes the alert.
	if err := tt.Click("Delete"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if chosen != 1 {
		t.Errorf("the chosen button is %d", chosen)
	}
	if open {
		t.Error("the alert stays open")
	}
	if tt.HasText("Delete Notes?") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Escape dismisses the alert, which reports it.
	open = true
	tt.Frame()
	tt.KeyDown(KeyEscape)
	tt.Frame()
	if chosen != AlertDismissed {
		t.Errorf("the alert reported %d", chosen)
	}
	if open {
		t.Error("the alert stays open")
	}

	// A click on a button chooses it even where Escape follows, which the
	// frame after no longer reports.
	open = true
	tt.Frame()
	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	tt.KeyDown(KeyEscape)
	tt.Frame()
	if chosen != 0 {
		t.Errorf("the alert reported %d", chosen)
	}

	// Outside a frame the alert answers nothing.
	p := tt.rt.public
	if AlertDialog(p, &open, "Delete?", "", "OK") != 0 {
		t.Error("an alert was answered outside a frame")
	}
}

// The public API of a dialog: it shows while it is open, and a click on
// the backdrop dismisses it, which the panel reports.
func TestDialogBaseAPI(t *testing.T) {
	var (
		open      bool
		dismissed bool
	)
	tt := NewTester(func(c *Context) {
		DialogBase(c, &open, func(backdrop Element, panel Element) {
			panel.Children(func() { Text(c, "Panel") })
			if panel.Dismissed() {
				dismissed = true
			}
		})
	}, 400, 300)
	tt.Frame()
	if tt.HasText("Panel") {
		t.Errorf("the dialog shows before it is opened: %q", tt.Texts())
	}

	// Open, the dialog shows its panel.
	open = true
	tt.Frame()
	if !tt.HasText("Panel") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if dismissed {
		t.Error("the dialog was dismissed before a click")
	}

	// A click on the backdrop dismisses it, and the panel reports it.
	tt.ClickAt(20, 20)
	tt.Frame()
	if open {
		t.Error("the dialog stays open")
	}
	if !dismissed {
		t.Error("the panel did not report the dismissal")
	}

	// Outside a frame the view builds no dialog.
	p := tt.rt.public
	if DialogBase(p, &open, func(backdrop Element, panel Element) {}) != (Element{}) {
		t.Error("a dialog was built outside a frame")
	}
}
