package ui

import (
	"image"
	"testing"
)

// ringBlue counts the pixels of the image that lean blue, as the theme's
// focus ring does over a light background.
func ringBlue(img *image.RGBA) int {
	n := 0
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.B > 200 && c.B-c.R > 40 && c.B-c.G > 30 {
				n++
			}
		}
	}
	return n
}

// ringGreen counts the pixels of the image that lean green.
func ringGreen(img *image.RGBA) int {
	n := 0
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.G > 150 && c.G > c.R && c.G > c.B {
				n++
			}
		}
	}
	return n
}

// ringRow counts the blue pixels of the image along the row y, between
// x0 and x1.
func ringRow(img *image.RGBA, x0, x1, y int) int {
	n := 0
	for x := x0; x < x1; x++ {
		c := img.RGBAAt(x, y)
		if c.B > 200 && c.B-c.R > 40 && c.B-c.G > 30 {
			n++
		}
	}
	return n
}

// The public API of how an element shows the keyboard focus.
func TestFocusAPI(t *testing.T) {
	t.Run("FocusRing", func(t *testing.T) {
		build := func(show bool) (*Tester, *bool) {
			visible := new(bool)
			tt := NewTester(func(c *Context) {
				Column(c).Padding(30, 0, 0, 40).Children(func() {
					b := Button(c, "Go").Size(80, 20).FocusRing(show)
					*visible = b.FocusVisible()
				})
			}, 180, 90)
			return tt, visible
		}
		ringed, visible := build(true)
		ringed.Key(0, KeyTab)
		ringed.Frame()
		if !*visible {
			t.Error("the button does not report the focus from the keyboard")
		}
		if n := ringBlue(ringed.Image()); n == 0 {
			t.Error("the focused button draws no ring")
		}

		own, _ := build(false)
		own.Key(0, KeyTab)
		own.Frame()
		if n := ringBlue(own.Image()); n != 0 {
			t.Errorf("the button draws its own ring after FocusRing(false), %d pixels", n)
		}
	})

	t.Run("FocusRingColor", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Column(c).Padding(30, 0, 0, 40).Children(func() {
				Button(c, "Go").Size(80, 20).FocusRingColor(RGB(0, 160, 0))
			})
		}, 180, 90)
		tt.Key(0, KeyTab)
		tt.Frame()

		if n := ringGreen(tt.Image()); n == 0 {
			t.Error("the ring of the button is not of the color FocusRingColor set")
		}
		if n := ringBlue(tt.Image()); n != 0 {
			t.Errorf("the ring keeps the theme's blue, %d pixels", n)
		}
	})

	t.Run("FocusRingOffset", func(t *testing.T) {
		build := func(offset bool) *Tester {
			return NewTester(func(c *Context) {
				Column(c).Padding(30, 0, 0, 40).Children(func() {
					b := Button(c, "Go").Size(80, 20)
					if offset {
						b = b.FocusRingOffset(5)
					}
				})
			}, 180, 90)
		}
		plain := build(false)
		plain.Key(0, KeyTab)
		plain.Frame()
		// The box of the button, and the rows its ring lies on with the
		// theme's offset of 1: just inside its top and bottom.
		r, ok := plain.Find("Go")
		if !ok {
			t.Fatal("the tester does not find the button")
		}
		if n := ringRow(plain.Image(), int(r.X)-6, int(r.X+r.W)+6, int(r.Y)-2); n == 0 {
			t.Error("the ring of the theme is not just outside the top of the button")
		}

		out := build(true)
		out.Key(0, KeyTab)
		out.Frame()
		rr, ok := out.Find("Go")
		if !ok {
			t.Fatal("the tester does not find the button")
		}
		if n := ringRow(out.Image(), int(rr.X)-8, int(rr.X+rr.W)+8, int(rr.Y)-6); n == 0 {
			t.Error("the ring does not stand 5 outside the button")
		}
		if n := ringRow(out.Image(), int(rr.X)-8, int(rr.X+rr.W)+8, int(rr.Y)-2); n != 0 {
			t.Error("the ring of the offset 5 still stands at the theme's offset")
		}
	})

	t.Run("FocusGroup", func(t *testing.T) {
		var inA, inB, visible, after bool
		tt := NewTester(func(c *Context) {
			Column(c).Padding(30, 0, 0, 40).Gap(8).Children(func() {
				Row(c).Gap(8).FocusGroup(Horizontal).Children(func() {
					a := Button(c, "first").Size(70, 20)
					b := Button(c, "second").Size(70, 20)
					inA, inB, visible = a.Focused(), b.Focused(), a.FocusVisible()
				})
				out := Button(c, "after").Size(70, 20)
				after = out.Focused()
			})
		}, 220, 110)

		tt.Key(0, KeyTab)
		tt.Frame()
		if !inA || inB {
			t.Errorf("Tab enters the group at (%v, %v), want the first", inA, inB)
		}
		if !visible {
			t.Error("the first button does not show the focus")
		}

		tt.Key(0, KeyRight)
		tt.Frame()
		if !inB || inA {
			t.Errorf("the right arrow moves within the group to (%v, %v), want the second", inA, inB)
		}

		tt.Key(0, KeyTab)
		tt.Frame()
		if inA || inB || !after {
			t.Errorf("Tab leaves the group at (%v, %v, after %v), want the button after it", inA, inB, after)
		}
	})

	t.Run("KeepFocus", func(t *testing.T) {
		build := func(keep bool) (*Tester, *bool) {
			focused := new(bool)
			tt := NewTester(func(c *Context) {
				Column(c).Padding(30, 0, 0, 40).Gap(10).Children(func() {
					inp := TextInput(c, new(string))
					*focused = inp.Focused()
					btn := Button(c, "go")
					if keep {
						btn = btn.KeepFocus()
					}
					_ = btn
				})
			}, 200, 120)
			return tt, focused
		}
		for _, keep := range []bool{true, false} {
			tt, focused := build(keep)
			tt.Key(0, KeyTab)
			tt.Frame()
			if !*focused {
				t.Fatalf("Tab does not focus the input (keep %v)", keep)
			}
			r, ok := tt.Find("go")
			if !ok {
				t.Fatalf("the tester does not find the button (keep %v)", keep)
			}
			tt.Press(r.X+r.W/2, r.Y+r.H/2)
			tt.Frame()
			if keep != *focused {
				t.Errorf("with KeepFocus %v the input is focused %v after a press on the button", keep, *focused)
			}
		}
	})

	t.Run("Outside a frame", func(t *testing.T) {
		var e Element
		for name, got := range map[string]Element{
			"FocusGroup":      e.FocusGroup(Horizontal),
			"FocusRing":       e.FocusRing(false),
			"FocusRingColor":  e.FocusRingColor(RGB(0, 0, 0)),
			"FocusRingOffset": e.FocusRingOffset(2),
			"KeepFocus":       e.KeepFocus(),
		} {
			if got != (Element{}) {
				t.Errorf("%s outside a frame gave %v, want nothing", name, got)
			}
		}
		if e.FocusVisible() {
			t.Error("FocusVisible outside a frame is true")
		}
	})
}
