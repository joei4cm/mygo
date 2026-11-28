package ui

import (
	"testing"
	"time"
)

// apiChildExitLeaving returns the copies of the elements that left, under
// the parent they left from.
func apiChildExitLeaving(e *node) []*node {
	var out []*node
	for c := e.first; c != nil; c = c.next {
		if c.leaving != 0 {
			out = append(out, c)
		}
	}
	return out
}

// apiChildExitPage runs a view of a page that goes, and returns the tester
// and the copy of the page that left.
func apiChildExitPage(t *testing.T, exit *Motion, childExit *Motion) (*Tester, *node) {
	t.Helper()
	page := "One"
	tt := NewTester(func(c *Context) {
		box := Box(c).Fill()
		if childExit != nil {
			box = box.ChildExit(*childExit)
		}
		box.Children(func() {
			tr := ElementTransition{Duration: 200 * time.Millisecond, Ease: Linear}
			if exit != nil {
				tr.Exit = exit
			}
			Column(c).Key(page).Fill().Transition(tr).Children(func() {
				Text(c, page)
			})
		})
	}, 200, 200)
	tt.Frame()
	page = "Two"
	tt.Frame()
	leaving := apiChildExitLeaving(tt.rt.c.root.first)
	if len(leaving) != 1 {
		t.Fatalf("%d copies of the page gone, not 1", len(leaving))
	}
	return tt, leaving[0]
}

// A parent says where the children that leave go, in place of the Exit of
// their own transitions.
func TestChildExitAPI(t *testing.T) {
	tt, ghost := apiChildExitPage(t, &Motion{X: 10}, &Motion{X: -40})
	defer tt.rt.close()
	// The copy starts where the page was.
	if ghost.x != 0 {
		t.Errorf("the copy starts at %v, not 0", ghost.x)
	}
	// Halfway through the exit, it is on its way the way the parent said.
	tt.Advance(100 * time.Millisecond)
	tt.Frame()
	if ghost.x != -20 {
		t.Errorf("halfway, the copy is at %v, not -20", ghost.x)
	}
	// Once the exit is over, the copy is gone.
	tt.Advance(200 * time.Millisecond)
	tt.Frame()
	if len(apiChildExitLeaving(tt.rt.c.root.first)) != 0 {
		t.Error("the copy is still there once its exit ended")
	}
}

// Without a parent saying where they go, the children leaving go the way
// their own transition says.
func TestChildExitAPIWithoutParent(t *testing.T) {
	tt, ghost := apiChildExitPage(t, &Motion{X: 10}, nil)
	defer tt.rt.close()
	tt.Advance(100 * time.Millisecond)
	tt.Frame()
	if ghost.x != 5 {
		t.Errorf("halfway, the copy is at %v, not 5", ghost.x)
	}
}

// A parent saying twice says the last: the copy goes the way of the last.
func TestChildExitAPITwice(t *testing.T) {
	page := "One"
	tt := NewTester(func(c *Context) {
		Box(c).Fill().ChildExit(Motion{X: -40}).ChildExit(Motion{X: -80}).Children(func() {
			Column(c).Key(page).Fill().Transition(ElementTransition{
				Duration: 200 * time.Millisecond,
				Ease:     Linear,
				Exit:     &Motion{X: 10},
			}).Children(func() {
				Text(c, page)
			})
		})
	}, 200, 200)
	defer tt.rt.close()
	tt.Frame()
	page = "Two"
	tt.Frame()
	leaving := apiChildExitLeaving(tt.rt.c.root.first)
	if len(leaving) != 1 {
		t.Fatalf("%d copies of the page gone, not 1", len(leaving))
	}
	tt.Advance(100 * time.Millisecond)
	tt.Frame()
	if x := leaving[0].x; x != -40 {
		t.Errorf("halfway, the copy is at %v, not -40", x)
	}
}

func TestChildExitAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.ChildExit(Motion{X: -40}).Valid() {
		t.Error("an absent element sent its children away")
	}
}
