package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// The public API of which way a slider runs: a slider of one's own that
// goes up takes its value from how far down the press stands, moves with
// the up and down keys, and reads vertical to assistive technology.
func TestVerticalAPI(t *testing.T) {
	v := 0.0
	tt := NewTester(func(c *Context) {
		SliderBase(c, &v, 0, 100).Size(20, 100).Label("Up").Vertical().Step(10)
	}, 200, 200)
	defer tt.rt.close()
	r, ok := tt.Find("Up")
	if !ok {
		t.Fatal("no slider")
	}
	// A press 30% of the way down the track stands at 70.
	tt.Press(r.X+10, r.Y+r.H*0.3)
	tt.Release(r.X+10, r.Y+r.H*0.3)
	if v != 70 {
		t.Errorf("pressing 30%% down a vertical slider: %v", v)
	}
	// The keys move it up and down.
	tt.Key(0, KeyUp)
	if v != 80 {
		t.Errorf("Up: %v", v)
	}
	tt.Key(0, KeyDown)
	if v != 70 {
		t.Errorf("Down: %v", v)
	}
	tt.send(platform.SurfaceEvent{Kind: platform.AccessibilityOn})
	if n := accessNode(t, tt.h.access, platform.RoleSlider, "Up"); n.States&platform.AccessVertical == 0 {
		t.Error("the vertical slider is not vertical to assistive technology")
	}
}

// Outside a frame no slider turns, and a context that is not building
// makes no slider.
func TestVerticalAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Vertical().Valid() {
		t.Error("an absent element turned")
	}
	var noContext Context
	v := 0.0
	if SliderBase(&noContext, &v, 0, 100).Valid() {
		t.Error("a context that is not building made a slider")
	}
}
