package ui

import "testing"

// nearSplit reports whether a and b are the same, within half a DIP.
func nearSplit(a, b float32) bool {
	return a > b-0.5 && a < b+0.5
}

// The public API of a split with the first pane above the second, taking
// a percentage of the room: dragging the divider moves it.
func TestSplitVerticalPercentAPI(t *testing.T) {
	pct := float32(50)
	changes := 0
	var first, second Rect
	tt := NewTester(func(c *Context) {
		e := SplitVerticalPercent(c, &pct, func() {
			first = Box(c).Fill().Bounds()
		}, func() {
			second = Box(c).Fill().Bounds()
		}).Grow(1).FillWidth()
		if e.Changed() {
			changes++
		}
	}, 200, 200)
	tt.Frame()

	// The panes share the room, as the percentage says.
	if !nearSplit(first.H, 99.5) || !nearSplit(second.H, 99.5) {
		t.Fatalf("the panes are %v and %v", first, second)
	}
	if !nearSplit(pct, 50) {
		t.Fatalf("the percentage is %v", pct)
	}

	// Dragging the divider down gives the room to the first pane, which
	// is a change.
	tt.Move(100, 100)
	tt.Press(100, 100)
	tt.Move(100, 140)
	tt.Release(100, 140)
	tt.Frame()
	if !nearSplit(first.H, 139.5) || !nearSplit(second.H, 59.5) {
		t.Fatalf("the panes are %v and %v", first, second)
	}
	if pct < 69 || pct > 71 {
		t.Fatalf("the percentage is %v", pct)
	}
	if changes != 1 {
		t.Fatalf("the move was reported %d times", changes)
	}

	// Dragging it back up takes the room away again.
	tt.Move(100, 140)
	tt.Press(100, 140)
	tt.Move(100, 120)
	tt.Release(100, 120)
	tt.Frame()
	if !nearSplit(first.H, 119.5) {
		t.Fatalf("the panes are %v and %v", first, second)
	}
	if pct < 59 || pct > 61 {
		t.Fatalf("the percentage is %v", pct)
	}

	// Outside a frame the view builds no split.
	p := tt.rt.public
	if SplitVerticalPercent(p, &pct, func() {}, func() {}).Valid() {
		t.Error("a split was built outside a frame")
	}
}
