package ui

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

// The public API of an image: it shows the bitmap at its own size, and at
// another size scaled to fit.
func TestImageAPI(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 100, 40))
	draw.Draw(src, src.Bounds(), image.NewUniform(color.RGBA{R: 255, A: 255}), image.Point{}, draw.Src)
	b := NewBitmap(src)

	// At its own size, 100 by 40.
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() { Image(c, b) })
	}, 400, 300)
	tt.Frame()
	px := tt.Image()
	if got := px.RGBAAt(5, 5); got.R < 0xf0 || got.G > 0x10 || got.B > 0x10 {
		t.Errorf("at (5, 5) the color is %v", got)
	}
	if got := px.RGBAAt(95, 35); got.R < 0xf0 || got.G > 0x10 || got.B > 0x10 {
		t.Errorf("at (95, 35) the color is %v", got)
	}
	if got := px.RGBAAt(150, 10); got.R > 0xf0 && got.G < 0x10 && got.B < 0x10 {
		t.Errorf("the picture reaches past its own size: %v", got)
	}

	// At half, 50 by 20.
	tt2 := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() { Image(c, b).Size(50, 20) })
	}, 400, 300)
	tt2.Frame()
	px2 := tt2.Image()
	if got := px2.RGBAAt(45, 10); got.R < 0xf0 || got.G > 0x10 || got.B > 0x10 {
		t.Errorf("at (45, 10) the color is %v", got)
	}
	if got := px2.RGBAAt(70, 10); got.R > 0xf0 && got.G < 0x10 && got.B < 0x10 {
		t.Errorf("the picture reaches past its half size: %v", got)
	}

	// Outside a frame the view builds no image.
	p := tt.rt.public
	if Image(p, b) != (Element{}) {
		t.Error("an image was built outside a frame")
	}
}
