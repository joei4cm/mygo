package ui

import "testing"

// Skew slants an element around its origin, after any transform before
// it: pixels and hits move with the slant.
func TestSkewAPI(t *testing.T) {
	clicks := 0
	var sx, sy float32
	tt := NewTester(func(c *Context) {
		box := Box(c).Absolute().Left(60).Top(60).Size(40, 40).
			Background(RGB(220, 20, 30)).
			TransformOrigin(0, 0).
			Skew(sx, sy)
		if box.Clicked() {
			clicks++
		}
	}, 220, 180)
	defer tt.rt.close()

	// Unskewed: (65, 95) is inside the box; (110, 90) is outside.
	if c := tt.Image().RGBAAt(65, 95); c.R != 220 || c.G != 20 {
		t.Fatalf("unskewed inside (65,95): %v", c)
	}
	if c := tt.Image().RGBAAt(110, 90); c.R == 220 && c.G == 20 {
		t.Fatalf("unskewed outside (110,90): %v", c)
	}

	// Skew 45° along X around the top-left origin: x' = x + y.
	// Local (20, 30) moves to (60+50, 60+30) = (110, 90), while (65, 95)
	// leaves the box.
	sx, sy = 45, 0
	tt.Frame()
	if c := tt.Image().RGBAAt(110, 90); c.R != 220 || c.G != 20 {
		t.Fatalf("skewed X inside (110,90): %v", c)
	}
	if c := tt.Image().RGBAAt(65, 95); c.R == 220 && c.G == 20 {
		t.Fatalf("skewed X left (65,95) red: %v", c)
	}
	tt.ClickAt(110, 90)
	if clicks != 1 {
		t.Fatalf("click inside skewed box: %d clicks, want 1", clicks)
	}
	tt.ClickAt(65, 95)
	if clicks != 1 {
		t.Fatalf("click outside skewed box: %d clicks, want 1", clicks)
	}

	// Skew 45° along Y around the top-left origin: y' = y + x.
	// Local (30, 20) moves to (60+30, 60+50) = (90, 110).
	sx, sy = 0, 45
	tt.Frame()
	if c := tt.Image().RGBAAt(90, 110); c.R != 220 || c.G != 20 {
		t.Fatalf("skewed Y inside (90,110): %v", c)
	}
	if c := tt.Image().RGBAAt(95, 65); c.R == 220 && c.G == 20 {
		t.Fatalf("skewed Y left (95,65) red: %v", c)
	}
}

// Skew composes after an existing transform on the same element.
func TestSkewAPIComposes(t *testing.T) {
	tt := NewTester(func(c *Context) {
		Box(c).Absolute().Left(60).Top(60).Size(40, 40).
			Background(RGB(20, 200, 40)).
			TransformOrigin(0, 0).
			Translate(20, 0).
			Skew(45, 0)
	}, 220, 180)
	defer tt.rt.close()

	// Translate(20,0) then Skew(45,0): local (30, 30) -> (50, 30) -> (80, 30),
	// so screen (140, 90) is green while (90, 90) is outside.
	if c := tt.Image().RGBAAt(140, 90); c.G != 200 || c.R != 20 {
		t.Fatalf("composed inside (140,90): %v", c)
	}
	if c := tt.Image().RGBAAt(90, 90); c.G == 200 && c.R == 20 {
		t.Fatalf("composed outside (90,90) green: %v", c)
	}
}

func TestSkewAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Skew(30, 0).Valid() {
		t.Error("an absent element skewed")
	}
}
