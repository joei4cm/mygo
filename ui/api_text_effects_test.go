package ui

import (
	"image"
	"testing"
)

// fxRed counts the red pixels of the image.
func fxRed(img *image.RGBA) int {
	n := 0
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.A > 100 && c.R > 150 && c.G < 100 && c.B < 100 {
				n++
			}
		}
	}
	return n
}

// fxEdge returns the rightmost column of the image that has ink on it.
func fxEdge(img *image.RGBA) int {
	last := 0
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.A > 100 && (c.R < 220 || c.G < 220) {
				if x > last {
					last = x
				}
			}
		}
	}
	return last
}

// fxInk counts the pixels of the image that are not the background.
func fxInk(img *image.RGBA) int {
	n := 0
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.A > 100 && (c.R < 220 || c.G < 220) {
				n++
			}
		}
	}
	return n
}

// The public API of how a text is stroked, shadowed, tabbed, spaced and
// tagged with a language.
func TestTextEffectsAPI(t *testing.T) {
	t.Run("TextStroke", func(t *testing.T) {
		build := func(stroke bool) *Tester {
			return NewTester(func(c *Context) {
				txt := Text(c, "Hello").TextColor(RGB(0, 0, 0)).FontSize(24)
				if stroke {
					txt = txt.TextStroke(3, RGB(255, 0, 0))
				}
			}, 200, 60)
		}
		plain := build(false)
		plain.Frame()
		if n := fxRed(plain.Image()); n != 0 {
			t.Errorf("the plain text has %d red pixels", n)
		}

		stroked := build(true)
		stroked.Frame()
		if n := fxRed(stroked.Image()); n == 0 {
			t.Error("the stroke of the text paints none of its color")
		}
	})

	t.Run("TextShadow", func(t *testing.T) {
		build := func(shadow bool) *Tester {
			return NewTester(func(c *Context) {
				txt := Text(c, "Hello").TextColor(RGB(0, 0, 0)).FontSize(24)
				if shadow {
					txt = txt.TextShadow(4, 4, 0, RGB(255, 0, 0))
				}
			}, 200, 60)
		}
		plain := build(false)
		plain.Frame()
		if n := fxRed(plain.Image()); n != 0 {
			t.Errorf("the plain text casts %d red pixels", n)
		}

		cast := build(true)
		cast.Frame()
		if n := fxRed(cast.Image()); n == 0 {
			t.Error("the shadow of the text paints none of its color")
		}
	})

	t.Run("TabSize", func(t *testing.T) {
		var set2, set8 float32
		build := func(n float32) *Tester {
			return NewTester(func(c *Context) {
				txt := Text(c, "a\tb").FontSize(16).TabSize(n)
				if n == 2 {
					set2 = txt.node().ts.tabSize
				} else {
					set8 = txt.node().ts.tabSize
				}
			}, 200, 60)
		}
		wide := build(8)
		wide.Frame()
		narrow := build(2)
		narrow.Frame()
		if set2 != 2 || set8 != 8 {
			t.Errorf("TabSize put (%v, %v) on the texts, want (2, 8)", set2, set8)
		}
		if edge := fxEdge(wide.Image()); edge <= fxEdge(narrow.Image()) {
			t.Errorf("a tab of 8 spaces runs no further than one of 2 (%d <= %d)",
				edge, fxEdge(narrow.Image()))
		}
	})

	t.Run("PreserveSpaces", func(t *testing.T) {
		var kept bool
		build := func(keep bool) *Tester {
			return NewTester(func(c *Context) {
				// Whitespace ending a line takes no room unless kept, and
				// an end-aligned line shows the difference.
				txt := Text(c, "hi      ").FontSize(16).Size(60, 20).TextAlign(End)
				if keep {
					txt = txt.PreserveSpaces()
					kept = txt.node().keepSpaces
				}
			}, 90, 50)
		}
		folded := build(false)
		folded.Frame()

		spaces := build(true)
		spaces.Frame()
		if !kept {
			t.Error("PreserveSpaces kept no spaces")
		}
		if edge := fxEdge(folded.Image()); edge-fxEdge(spaces.Image()) < 10 {
			t.Errorf("the kept spaces take no room (%d against %d)",
				fxEdge(spaces.Image()), edge)
		}
	})

	t.Run("Lang", func(t *testing.T) {
		var tag string
		tt := NewTester(func(c *Context) {
			txt := Text(c, "こんにちは").FontSize(16).Lang("ja")
			tag = txt.node().ts.lang
		}, 200, 60)
		tt.Frame()
		if tag != "ja" {
			t.Errorf("Lang put %q on the text, want \"ja\"", tag)
		}
		if n := fxInk(tt.Image()); n == 0 {
			t.Error("the text of the other language draws nothing")
		}
	})

	t.Run("Outside a frame", func(t *testing.T) {
		var e Element
		for name, got := range map[string]Element{
			"TabSize":        e.TabSize(2),
			"PreserveSpaces": e.PreserveSpaces(),
		} {
			if got != (Element{}) {
				t.Errorf("%s outside a frame gave %v, want nothing", name, got)
			}
		}
		e.TextStroke(3, RGB(255, 0, 0))
		e.TextShadow(4, 4, 0, RGB(255, 0, 0))
		e.Lang("ja")
	})
}
