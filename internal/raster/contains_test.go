package raster

import (
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

func TestContainsInterior(t *testing.T) {
	r := scene.Rect{X: 0, Y: 0, W: 100, H: 100}
	// A rounded corner of 20 keeps the middle and the straight sides, and
	// takes the corner's arc away.
	sharp := [4]float32{20, 0, 0, 0}
	for _, c := range []struct {
		x, y float32
		want bool
	}{
		{50, 50, true},
		{30, 5, true},
		{5, 30, true},
		{5, 5, false},
		{2, 2, false},
		{200, 50, false},
		{-1, 50, false},
	} {
		if got := Contains(r, sharp, false, c.x, c.y); got != c.want {
			t.Errorf("Contains(sharp corner, %v, %v) = %v, want %v", c.x, c.y, got, c.want)
		}
		if got := Contains(r, sharp, true, c.x, c.y); got != c.want {
			t.Errorf("Contains(continuous corner, %v, %v) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
	// Radii larger than the sides fit the rect: a corner of 100 in a
	// 100×100 rect centers at the far corner.
	if Contains(r, [4]float32{200, 0, 0, 0}, false, 30, 30) != true {
		t.Error("a point under a fitted corner arc is outside")
	}
	if Contains(r, [4]float32{200, 0, 0, 0}, false, 0, 0) != false {
		t.Error("the very corner of a fully rounded rect is outside")
	}
	// No corners: the whole interior counts.
	if !Contains(r, [4]float32{}, false, 99, 99) {
		t.Error("a square clip does not contain its interior corner")
	}
	if Contains(r, [4]float32{}, false, 100, 100) {
		t.Error("a square clip contains its corner point")
	}
}
