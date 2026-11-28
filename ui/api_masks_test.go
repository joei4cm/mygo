package ui

import (
	"image"
	"testing"
)

// maskTint returns how much the pixel at x, y of the image is red, from
// 0 (the white of the background) to 255 (the red of the box).
func maskTint(img *image.RGBA, x, y int) int {
	c := img.RGBAAt(x, y)
	return 255 - int(c.G)
}

// maskRed counts the red pixels of the image in the rectangle r.
func maskRed(img *image.RGBA, r Rect) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			c := img.RGBAAt(x, y)
			if c.A > 200 && c.R > 200 && c.G < 60 && c.B < 60 {
				n++
			}
		}
	}
	return n
}

// maskInk counts the pixels of the image that are not white, and those
// that keep a color of their own.
func maskInk(img *image.RGBA) (int, int) {
	ink, colored := 0, 0
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			c := img.RGBAAt(x, y)
			if c.R == 255 && c.G == 255 && c.B == 255 {
				continue
			}
			ink++
			if c.R != c.G || c.G != c.B {
				colored++
			}
		}
	}
	return ink, colored
}

// The public API of what an element paints and of what it hides.
func TestMasksAPI(t *testing.T) {
	t.Run("Gradient", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Box(c).Size(100, 20).Gradient(RGB(255, 0, 0), RGB(0, 0, 255), 90)
		}, 200, 100)
		tt.Frame()

		left := tt.Image().RGBAAt(4, 10)
		right := tt.Image().RGBAAt(96, 10)
		if left.A == 0 || right.A == 0 {
			t.Fatalf("the gradient painted nothing: left %v, right %v", left, right)
		}
		if left.R <= left.B {
			t.Errorf("the left of the gradient is %v, want the red of its start", left)
		}
		if right.B <= right.R {
			t.Errorf("the right of the gradient is %v, want the blue of its end", right)
		}
	})

	t.Run("MaskLinearGradient", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Box(c).Size(100, 20).Background(RGB(255, 0, 0)).
				MaskLinearGradient(LinearGradient{From: RGB(0, 0, 0), To: Transparent, Angle: 90})
		}, 200, 100)
		tt.Frame()

		if tint := maskTint(tt.Image(), 4, 10); tint < 150 {
			t.Errorf("the masked box shows little at its left, tint %d", tint)
		}
		if tint := maskTint(tt.Image(), 96, 10); tint > 40 {
			t.Errorf("the masked box shows at its right, tint %d", tint)
		}
	})

	t.Run("MaskRadialGradient", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Box(c).Size(100, 100).Background(RGB(255, 0, 0)).
				MaskRadialGradient(RadialGradient{
					Stops:  []GradientStop{{RGB(0, 0, 0), 0}, {Transparent, 1}},
					Radius: 30,
				})
		}, 200, 140)
		tt.Frame()

		if tint := maskTint(tt.Image(), 50, 50); tint < 150 {
			t.Errorf("the masked box shows little at its middle, tint %d", tint)
		}
		if tint := maskTint(tt.Image(), 2, 2); tint > 40 {
			t.Errorf("the masked box shows at its corner, tint %d", tint)
		}
	})

	t.Run("MaskConicGradient", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Box(c).Size(100, 100).Background(RGB(255, 0, 0)).
				MaskConicGradient(ConicGradient{
					Stops: []GradientStop{{RGB(0, 0, 0), 0}, {Transparent, 0.5}},
				})
		}, 200, 140)
		tt.Frame()

		n := maskRed(tt.Image(), Rect{X: 0, Y: 0, W: 100, H: 100})
		if n == 0 {
			t.Error("the conic mask hides the box whole")
		}
		if n >= 100*100 {
			t.Errorf("the conic mask shows the box whole, %d pixels", n)
		}
	})

	t.Run("Masks", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			solid := LinearGradient{From: RGB(0, 0, 0), To: RGB(0, 0, 0)}
			Box(c).Size(100, 40).Padding(10).Background(RGB(255, 0, 0)).Masks(
				Mask{Image: solid, Box: MaskContentBox, Composite: MaskExclude},
				Mask{Image: solid},
			)
		}, 200, 100)
		tt.Frame()

		if tint := maskTint(tt.Image(), 4, 20); tint < 150 {
			t.Errorf("the ring of the mask shows little, tint %d", tint)
		}
		if tint := maskTint(tt.Image(), 50, 20); tint > 40 {
			t.Errorf("the hole of the mask shows, tint %d", tint)
		}

		one := NewTester(func(c *Context) {
			solid := LinearGradient{From: RGB(0, 0, 0), To: RGB(0, 0, 0)}
			Box(c).Size(100, 40).Background(RGB(255, 0, 0)).Masks(Mask{Image: solid})
		}, 200, 100)
		one.Frame()
		if tint := maskTint(one.Image(), 50, 20); tint < 150 {
			t.Errorf("the one layer mask hides the box, tint %d", tint)
		}
	})

	t.Run("ClipInset", func(t *testing.T) {
		tt := NewTester(func(c *Context) {
			Box(c).Size(40, 40).Background(RGB(255, 0, 0)).ClipInset(10, 10, 10, 10, 0)
		}, 100, 100)
		tt.Frame()

		if tint := maskTint(tt.Image(), 2, 2); tint > 40 {
			t.Errorf("the clipped box shows at its edge, tint %d", tint)
		}
		if tint := maskTint(tt.Image(), 20, 20); tint < 150 {
			t.Errorf("the clipped box shows little inside, tint %d", tint)
		}
	})

	t.Run("ClipY", func(t *testing.T) {
		build := func(clip bool) *Tester {
			return NewTester(func(c *Context) {
				col := Column(c).Size(40, 40)
				if clip {
					col = col.ClipY()
				}
				col.Children(func() {
					Box(c).Size(40, 80).Background(RGB(255, 0, 0))
				})
			}, 100, 120)
		}
		plain := build(false)
		plain.Frame()
		if n := maskRed(plain.Image(), Rect{X: 0, Y: 40, W: 100, H: 80}); n == 0 {
			t.Error("the child of the column is clipped without ClipY")
		}

		clipped := build(true)
		clipped.Frame()
		if n := maskRed(clipped.Image(), Rect{X: 0, Y: 40, W: 100, H: 80}); n != 0 {
			t.Errorf("the child of the column draws below it, %d pixels", n)
		}
		if n := maskRed(clipped.Image(), Rect{X: 0, Y: 0, W: 100, H: 40}); n == 0 {
			t.Error("the child of the column draws nothing inside it")
		}
	})

	t.Run("Grayscale", func(t *testing.T) {
		pix := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for i := 0; i < len(pix.Pix); i += 4 {
			pix.Pix[i], pix.Pix[i+3] = 255, 255
		}
		bm := NewBitmap(pix)
		build := func(gray bool) *Tester {
			return NewTester(func(c *Context) {
				img := Image(c, bm)
				if gray {
					img = img.Grayscale()
				}
			}, 60, 60)
		}
		colored := build(false)
		colored.Frame()
		ink, red := maskInk(colored.Image())
		if ink == 0 || red == 0 {
			t.Errorf("the picture draws %d pixels, %d of them red", ink, red)
		}

		gray := build(true)
		gray.Frame()
		ink, red = maskInk(gray.Image())
		if ink == 0 {
			t.Error("the grayscaled picture draws nothing")
		}
		if red != 0 {
			t.Errorf("the grayscaled picture keeps a color, %d pixels", red)
		}
	})

	t.Run("Outside a frame", func(t *testing.T) {
		var e Element
		for name, got := range map[string]Element{
			"Gradient":           e.Gradient(RGB(0, 0, 0), RGB(255, 255, 255), 90),
			"ClipInset":          e.ClipInset(1, 1, 1, 1, 0),
			"ClipY":              e.ClipY(),
			"Grayscale":          e.Grayscale(),
			"MaskLinearGradient": e.MaskLinearGradient(LinearGradient{}),
			"MaskRadialGradient": e.MaskRadialGradient(RadialGradient{}),
			"MaskConicGradient":  e.MaskConicGradient(ConicGradient{}),
			"Masks":              e.Masks(),
		} {
			if got != (Element{}) {
				t.Errorf("%s outside a frame gave %v, want nothing", name, got)
			}
		}
	})
}
