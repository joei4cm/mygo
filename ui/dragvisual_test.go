package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

// shiftVisual moves an element with a visual transform, and those inside
// it, on the screen alone: their laid-out places stay.
func TestShiftVisual(t *testing.T) {
	e := &node{world: scene.Rotation(90)}
	ch := &node{}
	e.first = ch
	e.x, e.y = 10, 20

	wx, wy := e.world.Point(e.x, e.y)
	shiftVisual(e, 30, 40)

	// The laid-out place stays where it was.
	if e.x != 10 || e.y != 20 {
		t.Errorf("shiftVisual moved the element to %v,%v, want it to stay", e.x, e.y)
	}
	// The world moved by the shift: the element draws 30,40 over where it
	// did, rotated as before.
	nx, ny := e.world.Point(e.x, e.y)
	if nx != wx+30 || ny != wy+40 {
		t.Errorf("the element's place on the screen moved to %v,%v, want %v,%v over", nx, ny, wx+30, wy+40)
	}
	// The child's world moved with its parent's.
	cx0, cy0 := scene.Rotation(90).Point(0, 0)
	cx, cy := ch.world.Point(0, 0)
	if cx != cx0+30 || cy != cy0+40 {
		t.Errorf("the child draws at %v,%v, want %v,%v over", cx, cy, cx0+30, cy0+40)
	}
}

// restoreVisual gives an element, and those inside it, the worlds their
// states had.
func TestRestoreVisual(t *testing.T) {
	e := &node{world: scene.Rotation(90), st: &state{world: scene.Rotation(90)}}
	ch := &node{st: &state{world: scene.Affine{}}}
	e.first = ch
	saved, savedChild := e.world, ch.world

	e.world = scene.Translation(5, 5).Mul(e.world)
	ch.world = scene.Translation(5, 5).Mul(ch.world)
	restoreVisual(e)

	if e.world != saved || ch.world != savedChild {
		t.Error("restoreVisual did not give the element and its child their worlds back")
	}
}
