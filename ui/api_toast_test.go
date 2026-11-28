package ui

import (
	"testing"
	"time"
)

// The public API of a toast: showing one, closing it, and building the
// toasts of a window with a look of its own.
func TestToastAPI(t *testing.T) {
	var shown, acted, closed, closing bool
	var id string
	tt := NewTester(func(c *Context) {
		if !shown {
			shown = true
			c.Toast("Saved")
			c.ToastAction("Note deleted", "Undo", func() { acted = true })
			id = c.AddToast(Toast{ID: "upload", Title: "Uploading", Timeout: -1, OnClose: func() { closed = true }})
		}
		if closing {
			closing = false
			c.CloseToast(id)
		}
	}, 240, 160)
	tt.Frame()
	if !tt.HasText("Saved") || !tt.HasText("Note deleted") || !tt.HasText("Uploading") {
		t.Fatalf("the toasts did not show: %v", tt.Texts())
	}
	if id != "upload" {
		t.Fatalf("AddToast gave the toast the ID %q", id)
	}
	// The button of a toast runs its action and closes it.
	if err := tt.Click("Undo"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !acted {
		t.Error("the toast's action did not run")
	}
	if tt.HasText("Note deleted") {
		t.Error("the toast stayed after its action")
	}
	// Closing a toast by its ID runs its OnClose, and leaves the others.
	if closed || !tt.HasText("Uploading") {
		t.Fatal("the toast without a timeout closed by itself")
	}
	closing = true
	tt.Frame()
	if !closed {
		t.Error("the closed toast did not run its OnClose")
	}
	if tt.HasText("Uploading") {
		t.Error("the closed toast stayed")
	}
	if !tt.HasText("Saved") {
		t.Error("closing one toast closed another")
	}
	// A toast closes when its time is over.
	tt.Advance(toastTime + time.Second)
	tt.Frame()
	if tt.HasText("Saved") {
		t.Error("the toast stayed after its time")
	}
}

// The public API of a window's own toasts: the viewport that holds them
// and the toast the view builds in it.
func TestToastViewportAPI(t *testing.T) {
	var added, closed, retried bool
	var left time.Duration
	tt := NewTester(func(c *Context) {
		if !added {
			added = true
			c.AddToast(Toast{ID: "one", Title: "Uploading", Timeout: 5 * time.Second, OnClose: func() { closed = true }})
			c.AddToast(Toast{ID: "two", Title: "Failed", Timeout: -1, Action: "Retry", OnAction: func() { retried = true }})
		}
		ToastViewportBase(c, func(viewport Element, toasts []Toast) {
			viewport.Padding(16).AlignItems(End).Gap(8)
			for _, t := range toasts {
				toast := ToastBase(c, t)
				if t.ID == "one" {
					left = toast.Left()
				}
				toast.Root.Row().Gap(8).Padding(8).Background(RGB(20, 20, 20)).Children(func() {
					Text(c, t.Title).TextColor(RGB(255, 255, 255))
					if t.Action != "" {
						toast.ActionButton().Children(func() { Text(c, t.Action) })
					} else {
						toast.CloseButton().Children(func() { Text(c, "Close") })
					}
				})
			}
		})
	}, 240, 160)
	tt.Frame()
	if !tt.HasText("Uploading") {
		t.Fatalf("the window's own toast did not show: %v", tt.Texts())
	}
	if left != 5*time.Second {
		t.Errorf("a new toast had %v left, want 5s", left)
	}
	tt.Advance(2 * time.Second)
	tt.Frame()
	if left != 3*time.Second {
		t.Errorf("the toast had %v left after 2s, want 3s", left)
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !closed {
		t.Error("the toast's close button did not close the toast")
	}
	if tt.HasText("Uploading") {
		t.Error("the closed toast stayed")
	}
	// The action button runs the toast's action and closes it.
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !retried {
		t.Error("the toast's action button did not run its action")
	}
	if tt.HasText("Failed") {
		t.Error("the toast stayed after its action")
	}
}

// Outside a frame no toast shows, closes, or builds.
func TestToastAPIOutsideFrame(t *testing.T) {
	var c *Context
	c.Toast("Saved")
	c.ToastAction("Note deleted", "Undo", func() {})
	c.CloseToast("")
	if c.AddToast(Toast{Title: "Uploading"}) != "" {
		t.Error("a context that is not building added a toast")
	}
	if ToastViewportBase(c, func(Element, []Toast) {}).Valid() {
		t.Error("a context that is not building built a viewport")
	}
	p := ToastBase(c, Toast{Title: "Uploading"})
	if p.Root.Valid() || p.ActionButton().Valid() || p.CloseButton().Valid() || p.Left() != 0 {
		t.Error("a context that is not building built a toast")
	}
}
