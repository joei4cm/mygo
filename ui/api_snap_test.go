package ui

import (
	"testing"
	"time"
)

// The public API of where a scroll container comes to rest: a strip of
// cards that snaps rests showing a card, however far a notch of the wheel
// goes.
func TestSnapAPI(t *testing.T) {
	var s ScrollState
	tt := NewTester(func(c *Context) {
		ScrollHorizontal(c).Size(400, 100).ScrollSnap(SnapMandatory).TrackScroll(&s).Children(func() {
			for i := range 10 {
				Box(c).Key(i).Size(150, 100).Shrink(0).SnapAlign(Start)
			}
		})
	}, 400, 200)
	defer tt.rt.close()
	tt.Move(200, 50)
	// A notch rests at the next card, not where its short way ends.
	tt.Wheel(200, 50, 40, 0)
	tt.Advance(2 * time.Second)
	if s.X != 150 || tt.FramePending() {
		t.Errorf("after a notch: at %v, frames pending %v; want 150, none", s.X, tt.FramePending())
	}
	// Notches on the way go a card further each.
	tt.Wheel(200, 50, 40, 0)
	tt.Wheel(200, 50, 40, 0)
	tt.Advance(2 * time.Second)
	if s.X != 450 {
		t.Errorf("after two notches: at %v, want 450", s.X)
	}
}

// A column that snaps vertically rests showing a card, while a mark that
// snaps sideways only leaves its resting free.
func TestSnapAlignYAPI(t *testing.T) {
	var s ScrollState
	tt := NewTester(func(c *Context) {
		Scroll(c).Size(200, 80).ScrollSnap(SnapMandatory).TrackScroll(&s).Children(func() {
			for i := range 8 {
				Box(c).Key(i).Size(200, 60).Shrink(0).SnapAlignY(Start)
			}
		})
	}, 200, 200)
	defer tt.rt.close()
	tt.Move(100, 40)
	tt.Wheel(100, 40, 0, 40)
	tt.Advance(2 * time.Second)
	if s.Y != 60 {
		t.Errorf("after a notch down: at %v, want 60", s.Y)
	}
}

// What marks a place sideways only does not hold a column that snaps.
func TestSnapAlignXAPI(t *testing.T) {
	var s ScrollState
	tt := NewTester(func(c *Context) {
		Scroll(c).Size(200, 80).ScrollSnap(SnapMandatory).TrackScroll(&s).Children(func() {
			for i := range 8 {
				Box(c).Key(i).Size(200, 60).Shrink(0).SnapAlignX(Start)
			}
		})
	}, 200, 200)
	defer tt.rt.close()
	tt.Move(100, 40)
	tt.Wheel(100, 40, 0, 40)
	tt.Advance(2 * time.Second)
	if s.Y == 60 {
		t.Errorf("at %v, want anywhere but a place marked sideways only", s.Y)
	}
	if s.Y <= 0 {
		t.Errorf("the wheel moved nothing: at %v", s.Y)
	}
}

// Outside a frame these marks change nothing.
func TestSnapAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.ScrollSnap(SnapMandatory).Valid() {
		t.Error("an absent element snapped")
	}
	if absent.SnapAlign(Start).Valid() {
		t.Error("an absent element marked a place")
	}
	if absent.SnapAlignX(Start).Valid() {
		t.Error("an absent element marked a place sideways")
	}
	if absent.SnapAlignY(Start).Valid() {
		t.Error("an absent element marked a place vertically")
	}
}
