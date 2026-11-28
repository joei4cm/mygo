package ui

import (
	"testing"
	"time"
)

// The public API of a finger pulling a container past its top: a short
// pull does nothing, one far enough fires once and keeps refreshing
// until the app says it is done.
func TestPullToRefreshAPI(t *testing.T) {
	defer func(k overscroll) { overscrollKind = k }(overscrollKind)
	overscrollKind = overscrollBounce
	var r RefreshState
	fired := 0
	tt := NewTester(func(c *Context) {
		e := Scroll(c).Fill()
		if e.PullToRefresh(&r) {
			fired++
		}
		e.Children(func() { Box(c).Size(200, 50).Shrink(0) })
	}, 200, 200)
	defer tt.rt.close()
	tt.Frame()
	// A pull that does not reach refreshes nothing.
	tt.FingerDown(100, 20)
	tt.Advance(time.Millisecond)
	tt.FingerMove(100, 30)
	tt.Advance(16 * time.Millisecond)
	tt.FingerMove(100, 60)
	tt.FingerUp(100, 60)
	for range 100 {
		tt.Advance(16 * time.Millisecond)
	}
	if fired != 0 || r.Refreshing {
		t.Fatalf("a short pull fired %d times, refreshing %v", fired, r.Refreshing)
	}
	// One that reaches fires once and rests refreshing.
	tt.FingerDown(100, 20)
	tt.Advance(time.Millisecond)
	tt.FingerMove(100, 30)
	tt.Advance(16 * time.Millisecond)
	tt.FingerMove(100, 190)
	tt.FingerUp(100, 190)
	tt.Frame()
	if fired != 1 || !r.Refreshing {
		t.Fatalf("a long pull: fired %d, refreshing %v", fired, r.Refreshing)
	}
	// The app finishing refreshes it.
	r.Refreshing = false
	tt.Frame()
	if fired != 1 {
		t.Errorf("finishing fired %d times", fired)
	}
}

// Outside a frame nothing pulls.
func TestPullToRefreshAPIOutsideFrame(t *testing.T) {
	var absent Element
	var r RefreshState
	if absent.PullToRefresh(&r) {
		t.Error("an absent element refreshed")
	}
}
