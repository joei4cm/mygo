package ui

import (
	"testing"
)

// The public Context's toasts: a toast the view adds shows, and closes
// when the view says so; a toast with a button runs its action.
func TestContextToasts(t *testing.T) {
	var (
		id      string
		add     bool
		added   bool
		closeID string
		closed  bool
		undo    bool
		action  bool
		acted   bool
		message bool
		told    bool
	)
	tt := NewTester(func(c *Context) {
		if add && !added {
			added = true
			id = c.AddToast(Toast{Title: "Uploading", Timeout: -1, OnClose: func() { closed = true }})
			c.AddToast(Toast{Title: "Saving", Timeout: -1})
		}
		if action && !acted {
			acted = true
			c.ToastAction("Note deleted", "Undo", func() { undo = true })
		}
		if message && !told {
			told = true
			c.Toast("Saved")
		}
		if closeID != "" {
			c.CloseToast(closeID)
		}
	}, 400, 300)

	// A toast the view adds shows, with the ID it returns.
	add = true
	tt.Frame()
	if id == "" {
		t.Fatal("the toast has no ID")
	}
	if !tt.HasText("Uploading") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// The view closes it, and its OnClose runs, while the other toast
	// stays.
	closeID = id
	tt.Frame()
	if tt.HasText("Uploading") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if !tt.HasText("Saving") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if !closed {
		t.Error("the toast did not close")
	}

	// A toast with a button shows the button, which runs its action and
	// closes the toast.
	action = true
	tt.Frame()
	if !tt.HasText("Note deleted") || !tt.HasText("Undo") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	tt.Click("Undo")
	tt.Frame()
	if !undo {
		t.Error("the button of the toast did not run its action")
	}
	if tt.HasText("Note deleted") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A toast with a message shows it.
	message = true
	tt.Frame()
	if !tt.HasText("Saved") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Outside a frame the view shows no toast.
	p := tt.rt.public
	if p.AddToast(Toast{Title: "Late"}) != "" {
		t.Error("a toast was added outside a frame")
	}
	p.CloseToast(id)
	p.Toast("Late")
	p.ToastAction("Late", "Undo", func() {})
	tt.Frame()
	if tt.HasText("Late") {
		t.Errorf("the texts are %q", tt.Texts())
	}
}
