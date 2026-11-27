package ui

import (
	"testing"
	"time"
)

// The public Context's animations: a view that asks for an animation frame
// keeps getting frames, and one that disables the animations sees its
// changes at once.
func TestContextAnimations(t *testing.T) {
	var (
		frames  int
		instant bool
		moving  bool
		disable bool
	)
	tt := NewTester(func(c *Context) {
		frames++
		if disable {
			c.DisableAnimations()
		}
		instant = c.Instant()
		if moving {
			c.AnimationFrame()
		}
	}, 200, 200)

	// A view that asks for nothing gets no frame.
	moving = false
	tt.Frame()
	frames = 0
	tt.Advance(100 * time.Millisecond)
	if frames != 0 {
		t.Errorf("the view got %d frames", frames)
	}

	// A view that asks for an animation frame keeps getting them.
	moving = true
	tt.Frame()
	frames = 0
	tt.Advance(100 * time.Millisecond)
	if frames < 2 {
		t.Errorf("the view asked for %d frames", frames)
	}
	moving = false

	// A frame that disables the animations shows its changes at once, and
	// the frames after it do not.
	disable = true
	tt.Frame()
	if !instant {
		t.Error("the frame does not show its changes at once")
	}
	disable = false
	tt.Frame()
	if instant {
		t.Error("the frame shows its changes at once")
	}

	// Outside a frame the Context asks for nothing.
	p := tt.rt.public
	p.AnimationFrame()
	p.DisableAnimations()
	if p.Instant() {
		t.Error("the Context answers outside a frame")
	}
}
