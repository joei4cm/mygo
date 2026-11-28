package ui

import (
	"testing"
	"time"
)

// The public API of what a press of the pointer did: it went down inside
// the element or outside it, it was given up, and a finger held the
// element still for half a second. Each answer holds for the one frame
// it is asked in, so the test counts the frames that said so.
func TestPressAPI(t *testing.T) {
	inside, outside, canceled, long, clicks := 0, 0, 0, 0, 0
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			e := Box(c).Size(100, 100)
			if e.Clicked() {
				clicks++
			}
			if e.PressedInside() {
				inside++
			}
			if e.PressedOutside() {
				outside++
			}
			if e.PressCanceled() {
				canceled++
			}
			if e.LongPressed() {
				long++
			}
		})
	}, 300, 200)
	tt.Frame()
	if inside != 0 || outside != 0 || canceled != 0 || long != 0 {
		t.Error("a quiet frame reports a press")
	}

	// A press inside the element says so for the frame it went down in.
	tt.Press(50, 50)
	if inside != 1 || outside != 0 {
		t.Errorf("a press inside: %d frames said inside, %d outside", inside, outside)
	}

	// The frame after it is quiet again, and the release clicks once.
	tt.Release(50, 50)
	tt.Frame()
	if inside != 1 {
		t.Errorf("a frame after the down said inside %d times", inside)
	}

	// A press outside the element says so.
	tt.Press(250, 150)
	if outside != 1 || inside != 1 {
		t.Errorf("a press outside: %d frames said outside, %d inside", outside, inside)
	}
	tt.Release(250, 150)

	// A press given up says so, and clicks nothing.
	inside, outside, clicks = 0, 0, 0
	tt.Press(50, 50)
	tt.CancelPress()
	tt.Frame()
	if canceled != 1 {
		t.Errorf("a given-up press: %d frames said canceled", canceled)
	}
	if clicks != 0 {
		t.Errorf("the given-up press clicked %d times", clicks)
	}
	tt.Release(50, 50)

	// A finger held the element still for half a second says so, clicks
	// nothing as it lifts, and no frame after it says so again.
	inside, outside, canceled, clicks = 0, 0, 0, 0
	tt.FingerDown(50, 50)
	tt.Advance(400 * time.Millisecond)
	if long != 0 {
		t.Fatal("the finger held still for less than half a second long pressed")
	}
	tt.Advance(200 * time.Millisecond)
	tt.Frame()
	if long != 1 {
		t.Errorf("a finger held still for half a second: %d frames said long pressed", long)
	}
	tt.FingerUp(50, 50)
	tt.Frame()
	if clicks != 0 {
		t.Errorf("the long press clicked %d times", clicks)
	}
	if canceled != 0 || inside != 1 || outside != 0 {
		t.Errorf("a frame after the lift reports a press: %d canceled, %d inside, %d outside; the down alone said inside", canceled, inside, outside)
	}

	// Outside a frame no press is reported.
	if (Element{}).PressedInside() {
		t.Error("a press was reported outside a frame")
	}
	if (Element{}).PressedOutside() {
		t.Error("a press was reported outside a frame")
	}
	if (Element{}).PressCanceled() {
		t.Error("a press was reported outside a frame")
	}
	if (Element{}).LongPressed() {
		t.Error("a press was reported outside a frame")
	}
}
