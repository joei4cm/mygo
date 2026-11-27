package ui

import (
	"testing"
	"time"
)

// The public API of the toasts: the view builds them in a viewport of its
// own, counts down the time a toast shows still, and closes them with
// their buttons.
func TestToastPartsAPI(t *testing.T) {
	var (
		added    bool
		acted    bool
		viewport bool
		titles   []string
		left     time.Duration
		closed   bool
		undo     bool
	)
	tt := NewTester(func(c *Context) {
		if !added {
			added = true
			c.AddToast(Toast{ID: "upload", Title: "Uploading", Timeout: 10 * time.Second, OnClose: func() { closed = true }})
		}
		if !acted {
			acted = true
			c.AddToast(Toast{ID: "note", Title: "Note deleted", Timeout: -1, Action: "Undo", OnAction: func() { undo = true }})
		}
		ToastViewportBase(c, func(v Element, toasts []Toast) {
			viewport = true
			titles = titles[:0]
			for _, t := range toasts {
				titles = append(titles, t.Title)
			}
			for _, t := range toasts {
				toast := ToastBase(c, t)
				if t.ID == "upload" {
					left = toast.Left()
				}
				toast.Root.Children(func() {
					Text(c, t.Title)
					if t.ID == "note" {
						toast.ActionButton().Children(func() { Text(c, "Undo") })
					} else {
						toast.CloseButton().Children(func() { Text(c, "Close") })
					}
				})
			}
			// The time of the toasts moves, so the frames come.
			c.AnimationFrame()
		})
	}, 400, 300)

	// The viewport holds the toasts the window shows, and a toast counts
	// down the time it shows still.
	tt.Frame()
	if !viewport {
		t.Fatal("the viewport did not build")
	}
	if len(titles) != 2 || titles[0] != "Uploading" || titles[1] != "Note deleted" {
		t.Errorf("the toasts are %q", titles)
	}
	if left != 10*time.Second {
		t.Errorf("the toast shows still for %v", left)
	}

	// The time of a toast that runs moves on.
	tt.Advance(2 * time.Second)
	tt.Frame()
	if left != 8*time.Second {
		t.Errorf("the toast shows still for %v", left)
	}

	// The close button of a toast closes it.
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !closed {
		t.Error("the toast did not close")
	}
	if tt.HasText("Uploading") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if !tt.HasText("Note deleted") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// The action button of a toast runs its action and closes it.
	if err := tt.Click("Undo"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !undo {
		t.Error("the button of the toast did not run its action")
	}
	if tt.HasText("Note deleted") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A part that is not the view's tells nothing.
	var none ToastParts
	if none.Left() != 0 {
		t.Errorf("a part that is not the view's says %v", none.Left())
	}
	if none.ActionButton() != (Element{}) || none.CloseButton() != (Element{}) {
		t.Error("a part that is not the view's answers")
	}

	// Outside a frame the view builds no toast.
	p := tt.rt.public
	built := false
	if v := ToastViewportBase(p, func(v Element, toasts []Toast) { built = true }); v != (Element{}) {
		t.Error("a viewport was built outside a frame")
	}
	if built {
		t.Error("the toasts were built outside a frame")
	}
	if tb := ToastBase(p, Toast{Title: "Late"}); tb.Root != (Element{}) {
		t.Error("a toast was built outside a frame")
	}
}
