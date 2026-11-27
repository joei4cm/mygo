package ui

import (
	"fmt"
	"testing"
)

// ScrollTo brings a grid's item into the view where its align says, and
// ScrollIntoView as little as shows it whole.
func TestGridViewScrollTo(t *testing.T) {
	const n = 2000
	s := GridState{Label: func(i int) string { return fmt.Sprintf("Photo %d", i) }}
	tt := coreNewTester(func(c *context) {
		coreGridView(c, &s, n, 100, 80, func(i int) {
			coreTextf(c, "Photo %d", i)
		})
	}, 448, 400)
	tt.Frame()
	// 4 columns of 106: item 95 is the 23rd row.
	if s.cols != 4 {
		t.Fatalf("%d columns", s.cols)
	}
	row := func() (float32, bool) {
		r, ok := tt.Find("Photo 95")
		return r.Y, ok && r.Y >= 0 && r.Y+r.H <= 400
	}
	s.ScrollTo(95, Start)
	tt.Frame()
	yStart, in := row()
	if !in {
		t.Fatalf("ScrollTo(95, Start) leaves it at %v", yStart)
	}
	// Each alignment scrolls to its own place, Start at the top of the
	// view and End at its bottom.
	s.ScrollTo(95, Center)
	tt.Frame()
	yCenter, in := row()
	if !in || yCenter <= yStart {
		t.Errorf("ScrollTo(95, Center) leaves it at %v, Start at %v", yCenter, yStart)
	}
	s.ScrollTo(95, End)
	tt.Frame()
	yEnd, in := row()
	if !in || yEnd <= yCenter {
		t.Errorf("ScrollTo(95, End) leaves it at %v, Center at %v", yEnd, yCenter)
	}
	// Away from it and back, ScrollIntoView shows it whole, scrolling
	// as little as that takes.
	s.ScrollTo(0, Start)
	tt.Frame()
	s.ScrollIntoView(95)
	tt.Frame()
	yNear, in := row()
	if !in {
		t.Fatalf("ScrollIntoView(95) leaves it at %v", yNear)
	}
}
