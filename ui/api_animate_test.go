package ui

import (
	"testing"
	"time"
)

// The public API of the values an element animates: springs and colors
// that follow a target, sequences of keyframes, a progress that loops,
// a timeline of tweens, a pose the element moves to, and the questions
// about them.
func TestAnimateAPI(t *testing.T) {
	t.Run("Spring", func(t *testing.T) {
		var v, velocity, jumped float32
		var settled, arrived, jumpSettled bool
		target := float32(0)
		jumpTo := float32(0)
		tt := NewTester(func(c *Context) {
			box := Box(c).Size(40, 20)
			v = box.Spring("x", target)
			velocity = box.Velocity("x")
			settled = box.Settled("x")
			if box.Arrived("x") {
				arrived = true
			}
			jumped = box.Jump("j", jumpTo)
			jumpSettled = box.Settled("j")
		}, 120, 60)
		tt.Frame()
		if v != 0 || !settled {
			t.Errorf("the value at rest is %v, settled %v, want 0 and true", v, settled)
		}

		target = 100
		tt.Frame()
		if settled {
			t.Error("the value follows its new target at once")
		}

		tt.Advance(100 * time.Millisecond)
		tt.Frame()
		if velocity == 0 {
			t.Error("the moving value has no velocity")
		}

		tt.Advance(2 * time.Second)
		tt.Frame()
		if v != 100 || !settled {
			t.Errorf("the value ends at %v, settled %v, want 100 and true", v, settled)
		}
		if !arrived {
			t.Error("the value never reports coming to rest")
		}

		// Jump puts the value where it says at once.
		jumpTo = 50
		tt.Frame()
		if jumped != 50 {
			t.Errorf("Jump gave %v, want 50", jumped)
		}
		if !jumpSettled {
			t.Error("the value jumped to does not rest there")
		}
	})

	t.Run("AnimateColor", func(t *testing.T) {
		var col Color
		var settled bool
		target := RGB(255, 0, 0)
		tt := NewTester(func(c *Context) {
			box := Box(c).Size(40, 20)
			col = box.AnimateColor("c", target, Animation{Duration: 10 * time.Second})
			settled = box.Settled("c")
		}, 120, 60)
		tt.Frame()
		if col != RGB(255, 0, 0) || !settled {
			t.Errorf("the color at rest is %v, settled %v, want red and true", col, settled)
		}

		target = RGB(0, 0, 255)
		tt.Advance(5 * time.Second)
		tt.Frame()
		if settled {
			t.Error("the color follows its new target at once")
		}

		tt.Advance(15 * time.Second)
		tt.Frame()
		if col != RGB(0, 0, 255) || !settled {
			t.Errorf("the color ends at %v, settled %v, want blue and true", col, settled)
		}
	})

	t.Run("Keyframes", func(t *testing.T) {
		var v float32
		tt := NewTester(func(c *Context) {
			box := Box(c).Size(40, 20)
			v = box.Keyframes("k", []Keyframe{
				{At: 0, Value: 0},
				{At: time.Second, Value: 100},
			})
		}, 120, 60)
		tt.Frame()
		if v != 0 {
			t.Errorf("the sequence starts at %v, want 0", v)
		}

		tt.Advance(250 * time.Millisecond)
		tt.Frame()
		if v != 25 {
			t.Errorf("a quarter through, the sequence is at %v, want 25", v)
		}

		tt.Advance(2 * time.Second)
		tt.Frame()
		if v != 100 {
			t.Errorf("the sequence holds its last value at %v, want 100", v)
		}
	})

	t.Run("Loop", func(t *testing.T) {
		var v float32
		tt := NewTester(func(c *Context) {
			box := Box(c).Size(40, 20)
			v = box.Loop("l", time.Second, Linear)
		}, 120, 60)
		tt.Frame()
		if v != 0 {
			t.Errorf("the loop starts at %v, want 0", v)
		}

		tt.Advance(250 * time.Millisecond)
		tt.Frame()
		if !near32(v, 0.25) {
			t.Errorf("a quarter of the way round, the loop is at %v, want about 0.25", v)
		}

		tt.Advance(2 * time.Second)
		tt.Frame()
		if !near32(v, 0.25) {
			t.Errorf("two and a quarter turns in, the loop is at %v, want about 0.25", v)
		}
	})

	t.Run("Timeline", func(t *testing.T) {
		var x float32
		run := uint64(1)
		tt := NewTester(func(c *Context) {
			box := Box(c).Size(40, 20)
			tl := box.Timeline("intro", TimelineOptions{Run: run})
			x = tl.Tween(Tween{From: 0, To: 100, Duration: time.Second})
		}, 120, 60)
		tt.Frame()
		if x != 0 {
			t.Errorf("the tween starts at %v, want 0", x)
		}

		tt.Advance(500 * time.Millisecond)
		tt.Frame()
		if x <= 0 || x >= 100 {
			t.Errorf("halfway through, the tween is at %v, want between 0 and 100", x)
		}

		tt.Advance(2 * time.Second)
		tt.Frame()
		if x != 100 {
			t.Errorf("the tween ends at %v, want 100", x)
		}

		// Playing the timeline over with a new Run starts it again.
		run = 2
		tt.Frame()
		if x >= 100 {
			t.Errorf("the tween played again starts at %v, want before 100", x)
		}
	})

	t.Run("AnimatePose", func(t *testing.T) {
		var trans bool
		var y float32
		target := Pose{}
		tt := NewTester(func(c *Context) {
			box := Box(c).Size(40, 20)
			box.AnimatePose(target, Animation{Duration: 100 * time.Millisecond})
			trans = box.Transitioning()
			y = box.Bounds().Y
		}, 120, 200)
		tt.Frame()
		if trans || y != 0 {
			t.Errorf("the pose at rest shows at %v, transitioning %v, want 0 and false", y, trans)
		}

		target = Pose{Y: 40}
		tt.Frame()
		if !trans {
			t.Error("the element does not report moving to its new pose")
		}

		tt.Advance(500 * time.Millisecond)
		tt.Frame()
		if y != 40 {
			t.Errorf("the pose ends at %v, want 40", y)
		}
		if trans {
			t.Error("the element still reports moving once it is there")
		}
	})

	t.Run("Outside a frame", func(t *testing.T) {
		var e Element
		if got := e.AnimateColor("c", RGB(255, 0, 0), Animation{}); got != (Color{}) {
			t.Errorf("AnimateColor outside a frame gave %v, want nothing", got)
		}
		if got := e.AnimatePose(Pose{}, Animation{}); got != (Element{}) {
			t.Errorf("AnimatePose outside a frame gave %v, want nothing", got)
		}
		if got := e.Keyframes("k", []Keyframe{{At: 0, Value: 7}, {At: time.Second, Value: 9}}); got != 0 {
			t.Errorf("Keyframes outside a frame gave %v, want nothing", got)
		}
		if got := e.Spring("x", 42); got != 0 {
			t.Errorf("Spring outside a frame gave %v, want nothing", got)
		}
		if got := e.Jump("x", 42); got != 0 {
			t.Errorf("Jump outside a frame gave %v, want nothing", got)
		}
		if got := e.Loop("l", time.Second, Linear); got != 0 {
			t.Errorf("Loop outside a frame gave %v, want 0", got)
		}
		if got := e.Timeline("t", TimelineOptions{}); got != nil {
			t.Errorf("Timeline outside a frame gave %v, want none", got)
		}
		if e.Settled("x") {
			t.Error("Settled outside a frame is true")
		}
		if e.Arrived("x") {
			t.Error("Arrived outside a frame is true")
		}
		if got := e.Velocity("x"); got != 0 {
			t.Errorf("Velocity outside a frame gave %v, want nothing", got)
		}
		if e.Transitioning() {
			t.Error("Transitioning outside a frame is true")
		}
	})
}

// near32 reports whether v is within a twentieth of want.
func near32(v, want float32) bool { return v > want-0.05 && v < want+0.05 }
