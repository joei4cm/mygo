package ui

import (
	"image/color"
	"testing"
)

// PreviousBoundary and NextBoundary step between whole-grapheme UTF-16
// offsets, snapping interior offsets and clamping at the ends.
func TestTextLayoutBoundaries(t *testing.T) {
	// "á" is one grapheme of 2 UTF-16 units (0..2),
	// "😀" is one grapheme of 2 UTF-16 units (2..4),
	// "b" is one grapheme of 1 UTF-16 unit (4..5).
	l := ShapeText("a\u0301\U0001F600b", Font{Size: 16}, 0)

	for _, c := range []struct {
		index, want int
	}{
		{-2, 0},
		{0, 0},
		{1, 0}, // inside a + combining acute snaps to 0
		{2, 0},
		{3, 0}, // inside surrogate pair snaps to 2, then steps to 0
		{4, 2},
		{5, 4},
		{99, 4}, // past the end clamps to 5, then steps to 4
	} {
		if got := l.PreviousBoundary(c.index); got != c.want {
			t.Errorf("PreviousBoundary(%d) = %d, want %d", c.index, got, c.want)
		}
	}

	for _, c := range []struct {
		index, want int
	}{
		{-2, 2},
		{0, 2},
		{1, 2}, // inside a + combining acute snaps to 0, then steps to 2
		{2, 4},
		{3, 4}, // inside surrogate pair snaps to 2, then steps to 4
		{4, 5},
		{5, 5},
		{99, 5},
	} {
		if got := l.NextBoundary(c.index); got != c.want {
			t.Errorf("NextBoundary(%d) = %d, want %d", c.index, got, c.want)
		}
	}

	// An empty layout keeps every boundary at 0.
	empty := ShapeText("", Font{Size: 16}, 0)
	if got := empty.PreviousBoundary(0); got != 0 {
		t.Errorf("empty PreviousBoundary(0) = %d, want 0", got)
	}
	if got := empty.NextBoundary(0); got != 0 {
		t.Errorf("empty NextBoundary(0) = %d, want 0", got)
	}
}

// Painter.TextLayout paints a retained layout with its default color and
// the colors of its spans, and ignores a nil layout.
func TestPainterTextLayout(t *testing.T) {
	green := RGB(0, 180, 0)
	red := RGB(220, 0, 0)
	layout := ShapeRichText([]Span{
		{Text: "██", Color: green, colorSet: true},
		{Text: "██"},
	}, Font{Size: 20}, 0)
	bg := color.RGBA{255, 255, 255, 255}
	white := RGB(255, 255, 255)

	tt := coreNewTester(func(c *context) {
		coreBox(c).Size(120, 40).Background(white).Draw(func(p *Painter, r Rect) {
			p.TextLayout(nil, r.X, r.Y, red)
			p.TextLayout(layout, r.X+4, r.Y+4, red)
		})
	}, 120, 40)

	// The first span paints in its own green, the second in the red the
	// painter was given.
	img := tt.Image()
	var greens, reds int
	for y := range 40 {
		for x := range 120 {
			switch img.RGBAAt(x, y) {
			case color.RGBA{R: green.R, G: green.G, B: green.B, A: 255}:
				greens++
			case color.RGBA{R: red.R, G: red.G, B: red.B, A: 255}:
				reds++
			}
		}
	}
	if greens == 0 {
		t.Error("Painter.TextLayout did not paint the green span")
	}
	if reds == 0 {
		t.Error("Painter.TextLayout did not paint the default red color")
	}

	// Drawing only a nil layout leaves the background untouched.
	ttNil := coreNewTester(func(c *context) {
		coreBox(c).Size(40, 20).Background(white).Draw(func(p *Painter, r Rect) {
			p.TextLayout(nil, r.X, r.Y, red)
		})
	}, 40, 20)
	imgNil := ttNil.Image()
	for y := range 20 {
		for x := range 40 {
			if got := imgNil.RGBAAt(x, y); got != bg {
				t.Fatalf("nil TextLayout painted %v at (%d,%d), want background", got, x, y)
			}
		}
	}
}
