package ui

import "testing"

// The public API of a container that scrolls its children both ways: a
// wheel over it moves what it holds.
func TestScrollBothAPI(t *testing.T) {
	var state ScrollState
	tt := NewTester(func(c *Context) {
		ScrollBoth(c).Size(200, 200).TrackScroll(&state).Children(func() {
			Box(c).Size(600, 600).Label("Big")
		})
	}, 300, 300)

	// The container scrolls as far as its content goes.
	if state.MaxX != 400 || state.MaxY != 400 {
		t.Fatalf("the container scrolls to %v", state)
	}

	// A wheel down scrolls the content up.
	tt.Scroll(100, 100, 0, 80)
	if state.Y != 80 {
		t.Fatalf("the content did not scroll down: %v", state)
	}
	if state.X != 0 {
		t.Fatalf("the content scrolled sideways: %v", state)
	}

	// A wheel right scrolls the content left.
	tt.Scroll(100, 100, 60, 0)
	if state.X != 60 {
		t.Fatalf("the content did not scroll sideways: %v", state)
	}
	if state.Y != 80 {
		t.Fatalf("the content scrolled down: %v", state)
	}

	// Outside a frame the view builds no container.
	p := tt.rt.public
	if ScrollBoth(p).Valid() {
		t.Error("a container was built outside a frame")
	}
}
