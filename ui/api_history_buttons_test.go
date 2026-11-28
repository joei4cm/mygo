package ui

import (
	"testing"
)

// The public API of the back and the forward buttons: they follow the
// pages of a router's history, and sit disabled at its ends.
func TestHistoryButtonsAPI(t *testing.T) {
	var page string
	r := NewRouter("/")
	tt := NewTester(func(c *Context) {
		r.View(c, func(rt *Route) {
			switch {
			case rt.Match("/"):
				page = "Home"
				Text(c, "Home")
			case rt.Match("/one"):
				page = "One"
				Text(c, "One")
			case rt.Match("/two"):
				page = "Two"
				Text(c, "Two")
			}
		})
		Row(c).Children(func() {
			BackButton(c, r)
			ForwardButton(c, r)
		})
	}, 400, 300)
	tt.Frame()
	if page != "Home" {
		t.Errorf("the page is %q", page)
	}
	back, ok := tt.Access("Back")
	if !ok || !back.Disabled {
		t.Errorf("back: %+v, %v", back, ok)
	}
	forward, ok := tt.Access("Forward")
	if !ok || !forward.Disabled {
		t.Errorf("forward: %+v, %v", forward, ok)
	}

	// Forward through two pages, and back sits ready.
	r.Push("/one")
	tt.Frame()
	r.Push("/two")
	tt.Frame()
	if page != "Two" {
		t.Errorf("the page is %q", page)
	}
	if back, _ = tt.Access("Back"); back.Disabled {
		t.Error("back sits disabled in the middle of the history")
	}
	if forward, _ = tt.Access("Forward"); !forward.Disabled {
		t.Error("forward sits enabled at the end of the history")
	}

	// Back goes to the page before, forward goes back to it.
	if err := tt.Click("Back"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if page != "One" {
		t.Errorf("the page is %q", page)
	}
	if err := tt.Click("Forward"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if page != "Two" {
		t.Errorf("the page is %q", page)
	}

	// Outside a frame the view builds no button.
	p := tt.rt.public
	if BackButton(p, r) != (Element{}) || ForwardButton(p, r) != (Element{}) {
		t.Error("a button was built outside a frame")
	}
}
