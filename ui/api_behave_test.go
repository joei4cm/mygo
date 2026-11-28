package ui

import (
	"image"
	"image/color"
	"testing"
)

// behaveRed counts the strong red pixels of the image.
func behaveRed(img *image.RGBA) int {
	n := 0
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.A > 200 && c.R > 200 && c.G < 60 && c.B < 60 {
				n++
			}
		}
	}
	return n
}

// The public API of how an element behaves under the pointer and among
// its siblings: the cursor it shows, how a disabled one dims, where it
// takes presses, what it lets pass through, the lines between its
// children, and where its picture stands in its box.
func TestBehaveAPI(t *testing.T) {
	t.Run("Cursor", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Box(c).Size(100, 50).Debug().Unselectable().Cursor(CursorPointer)
		}, 200, 100)
		tt.Move(50, 25)
		if got := tt.Cursor(); got != CursorPointer {
			t.Errorf("the pointer over the box is %v, want the one it sets", got)
		}
		tt.Move(150, 80)
		if got := tt.Cursor(); got != CursorDefault {
			t.Errorf("the pointer off the box is %v, want the default", got)
		}
	})

	t.Run("DimDisabled", func(t *testing.T) {
		build := func(dim bool) *Tester {
			return NewTester(func(c *Context) {
				b := Box(c).Size(100, 50).Background(RGB(255, 0, 0)).Disabled(true)
				if !dim {
					b = b.DimDisabled(false)
				}
			}, 200, 100)
		}
		dimmed := build(true)
		dimmed.Frame()
		if n := behaveRed(dimmed.Image()); n != 0 {
			t.Errorf("the disabled box keeps its color, %d red pixels", n)
		}

		whole := build(false)
		whole.Frame()
		if n := behaveRed(whole.Image()); n != 100*50 {
			t.Errorf("DimDisabled(false) draws %d red pixels, want the box whole", n)
		}
	})

	t.Run("HitTest", func(t *testing.T) {
		var left, right int
		tt := NewTester(func(c *Context) {
			Column(c).Children(func() {
				Box(c).Size(200, 50).OnClick(func() { left++ }).
					HitTest(func(x, y float32, r Rect) bool { return x < r.W/2 })
				Box(c).Size(200, 50).OnClick(func() { right++ })
			})
		}, 240, 140)
		tt.Frame()

		tt.ClickAt(30, 25)
		tt.ClickAt(150, 25)
		if left != 1 || right != 0 {
			t.Errorf("the halves took (%d, %d) of the presses, want the left one only", left, right)
		}
	})

	t.Run("PassThrough", func(t *testing.T) {
		build := func(pass bool) (*Tester, *int, *int) {
			over, under := new(int), new(int)
			tt := NewTester(func(c *Context) {
				Column(c).Children(func() {
					Box(c).Size(200, 50).Background(RGB(0, 128, 0)).OnClick(func() { *under++ })
				})
				Overlay(c, func() {
					top := Box(c).Absolute().Left(0).Top(0).Size(200, 50).OnClick(func() { *over++ })
					if pass {
						top = top.PassThrough()
					}
				})
			}, 240, 120)
			return tt, over, under
		}
		stopped, over, under := build(false)
		stopped.ClickAt(50, 25)
		if *over != 1 || *under != 0 {
			t.Errorf("the panel over it took the press (%d, %d), want it stopped", *over, *under)
		}

		through, over, under := build(true)
		through.ClickAt(50, 25)
		if *over != 0 || *under != 1 {
			t.Errorf("the press went (%d, %d), want it through to the box under", *over, *under)
		}
	})

	t.Run("Dividers", func(t *testing.T) {
		build := func(div bool) *Tester {
			return NewTester(func(c *Context) {
				row := Row(c).Size(200, 60).Gap(10)
				if div {
					row = row.Dividers(2, RGB(255, 0, 0))
				}
				row.Children(func() {
					Box(c).Size(40, 40)
					Box(c).Size(40, 40)
					Box(c).Size(40, 40)
				})
			}, 220, 100)
		}
		plain := build(false)
		plain.Frame()
		if n := behaveRed(plain.Image()); n != 0 {
			t.Errorf("the row without dividers draws %d red pixels", n)
		}

		divided := build(true)
		divided.Frame()
		// Two lines, 2 DIPs thick, across the 60-high row.
		if n := behaveRed(divided.Image()); n < 200 {
			t.Errorf("the dividers draw only %d red pixels, want two lines across the row", n)
		}
	})

	t.Run("ObjectPosition", func(t *testing.T) {
		// A picture of a red sky over a blue sea, taller than its box.
		pix := image.NewRGBA(image.Rect(0, 0, 20, 40))
		for y := 0; y < 40; y++ {
			for x := 0; x < 20; x++ {
				if y < 20 {
					pix.SetRGBA(x, y, color.RGBA{R: 255, A: 255})
				} else {
					pix.SetRGBA(x, y, color.RGBA{B: 255, A: 255})
				}
			}
		}
		bm := NewBitmap(pix)
		build := func(top bool) *Tester {
			return NewTester(func(c *Context) {
				img := Image(c, bm).Size(20, 20).Fit(Cover)
				if top {
					img = img.ObjectPosition(0.5, 0)
				}
			}, 90, 60)
		}
		middle := build(false)
		middle.Frame()
		if c := middle.Image().RGBAAt(8, 18); c.B < 200 {
			t.Errorf("the picture centers on the sea at its bottom, got %v", c)
		}

		sky := build(true)
		sky.Frame()
		if c := sky.Image().RGBAAt(8, 18); c.R < 200 {
			t.Errorf("ObjectPosition(0.5, 0) does not show the sky at the bottom, got %v", c)
		}
	})

	t.Run("Outside a frame", func(t *testing.T) {
		var e Element
		for name, got := range map[string]Element{
			"Cursor":         e.Cursor(CursorPointer),
			"DimDisabled":    e.DimDisabled(false),
			"HitTest":        e.HitTest(nil),
			"PassThrough":    e.PassThrough(),
			"Dividers":       e.Dividers(2, RGB(255, 0, 0)),
			"ObjectPosition": e.ObjectPosition(0.5, 0),
			"Debug":          e.Debug(),
			"Unselectable":   e.Unselectable(),
		} {
			if got != (Element{}) {
				t.Errorf("%s outside a frame gave %v, want nothing", name, got)
			}
		}
	})
}
