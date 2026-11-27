package ui

import (
	"testing"
	"time"

	"github.com/egoist/mygo/internal/platform"
)

// A tap on the status bar flings what is under it back to its top, and a
// low-memory note from the system lets the next frame make again what it
// let go of.
func TestSurfaceScrollToTop(t *testing.T) {
	var s ScrollState
	tt := coreNewTester(func(c *context) {
		coreScroll(c).Fill().TrackScroll(&s).Children(func() {
			coreBox(c).Height(1000).Shrink(0)
		})
	}, 200, 200)
	// The insets are what a phone's bars cover the screen with: they
	// enable the touch handling the system's notes come through.
	tt.SetInsets(0, 0, 0, 0)
	s.Y = 300
	tt.Frame()
	if s.Y != 300 {
		t.Fatalf("scrolled to %v", s.Y)
	}
	tt.rt.systemEvent(platform.SurfaceEvent{Kind: platform.SurfaceScrollToTop})
	// The fling takes it back to the top over the frames of its decay: the
	// first builds with it on, and the rest follow its refresh.
	tt.Frame()
	for range 40 {
		tt.Advance(100 * time.Millisecond)
		if s.Y == 0 {
			break
		}
	}
	if s.Y != 0 {
		t.Errorf("the scroll to the top left it at %v", s.Y)
	}
}

// A low-memory note lets go of what the next frame makes again.
func TestSurfaceLowMemory(t *testing.T) {
	tt := coreNewTester(func(c *context) {
		coreText(c, "Kept")
	}, 100, 40)
	tt.Frame()
	tt.rt.systemEvent(platform.SurfaceEvent{Kind: platform.SurfaceLowMemory})
	if !tt.asked() {
		t.Error("the note asks for no frame")
	}
	tt.Frame()
	if !tt.HasText("Kept") {
		t.Error("the text is gone after a low-memory note")
	}
}

// A screen the system pops itself leaves the router at the page of the
// last, as the system's stack has it.
func TestSurfaceStackPopped(t *testing.T) {
	enableTouch()
	s := &stackHost{can: true, moves: true}
	r := NewRouter("/")
	r.SystemStack = true
	tt := stackTester(s, r, nil)
	tt.Frame()
	r.Push("/notes")
	tt.Frame()
	tt.rt.systemEvent(platform.SurfaceEvent{Kind: platform.SurfaceStackPopped, To: 1})
	tt.Frame()
	if r.Path() != "/" {
		t.Errorf("the popped stack leaves the router at %s", r.Path())
	}
}

// The system's gesture of going back takes the page, follows the finger
// and, cancelled, lets it go.
func TestSurfaceBackGesture(t *testing.T) {
	enableTouch()
	s := &stackHost{can: true}
	r := NewRouter("/")
	r.SystemStack = true
	tt := stackTester(s, r, nil)
	tt.Frame()
	r.Push("/notes")
	tt.Frame()
	if tt.rt.edgeBack != r {
		t.Fatal("the router does not take the system's back gesture")
	}
	tt.rt.systemEvent(platform.SurfaceEvent{Kind: platform.SurfaceBackGesture, Text: "started"})
	tt.Frame()
	if !r.swiping {
		t.Fatal("the gesture does not take the page")
	}
	tt.rt.systemEvent(platform.SurfaceEvent{Kind: platform.SurfaceBackGesture, Text: "progress", Scale: 1})
	tt.Frame()
	if r.swipe != 0.35 {
		t.Errorf("the page went %v of the way, want 0.35", r.swipe)
	}
	tt.rt.systemEvent(platform.SurfaceEvent{Kind: platform.SurfaceBackGesture, Text: "cancelled"})
	tt.Frame()
	if r.swiping {
		t.Error("the cancelled gesture still takes the page")
	}
}
