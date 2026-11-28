package ui

import "testing"

// The public API of the sizes a control takes: a box fills its parent's
// height, and a segmented control spreads its segments over the width it
// fills, each taking an equal share.
func TestFillAPI(t *testing.T) {
	red := RGB(220, 30, 30)
	tt := NewTester(func(c *Context) {
		Row(c).Size(200, 80).Children(func() {
			Box(c).Size(40, 10).FillHeight().Background(red)
		})
	}, 200, 80)
	defer tt.rt.close()
	tt.Frame()
	reds := 0
	for y := 0; y < 80; y++ {
		for x := 0; x < 60; x++ {
			if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
				reds++
			}
		}
	}
	if reds < 2000 {
		t.Errorf("the box fills %d px of its parent, want most of 40 by 80", reds)
	}
}

// A segmented control fills its parent's width only when it is asked to,
// and its segments share it equally.
func TestFillEquallyAPI(t *testing.T) {
	choice := 0
	equally := false
	tt := NewTester(func(c *Context) {
		seg := Segmented(c, &choice, "List", "Grid")
		if equally {
			seg.FillEqually()
		}
	}, 300, 60)
	defer tt.rt.close()
	tt.Frame()
	r, ok := tt.Find("Grid")
	if !ok {
		t.Fatal("no Grid segment")
	}
	if r.X+r.W/2 >= 150 {
		t.Errorf("without FillEqually the segment ends at %v, the control shrunk", r.X+r.W/2)
	}
	equally = true
	tt.Frame()
	r, ok = tt.Find("Grid")
	if !ok {
		t.Fatal("no Grid segment")
	}
	if r.X+r.W/2 <= 150 {
		t.Errorf("the second segment stands at %v..%v, want its half of the width", r.X, r.X+r.W)
	}
	// Each segment takes its half: a press past the middle hits the second.
	tt.ClickAt(225, r.Y+r.H/2)
	if choice != 1 {
		t.Errorf("a press on the second half chose %d", choice)
	}
}

// Outside a frame these sizes change nothing.
func TestFillAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.FillHeight().Valid() {
		t.Error("an absent element filled a height")
	}
	if absent.FillEqually().Valid() {
		t.Error("an absent element filled equally")
	}
}
