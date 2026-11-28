package ui

import (
	"image"
	"image/color"
	"testing"
)

func apiFitBitmap() *Bitmap {
	// Four quadrants: red top-left, green top-right, blue bottom-left, yellow bottom-right.
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	colors := [2][2]color.RGBA{{{255, 0, 0, 255}, {0, 255, 0, 255}}, {{0, 0, 255, 255}, {255, 255, 0, 255}}}
	for y := range 32 {
		for x := range 32 {
			img.Set(x, y, colors[y/16][x/16])
		}
	}
	return NewBitmap(img)
}

func TestFitAPI(t *testing.T) {
	red, green, blue, yellow, white := color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{255, 255, 0, 255}, color.RGBA{255, 255, 255, 255}
	fits := []struct {
		name string
		fit  Fit
		want map[[2]int]color.RGBA
	}{
		{"contain", Contain, map[[2]int]color.RGBA{{20, 8}: red, {44, 8}: green, {20, 24}: blue, {44, 24}: yellow}},
		{"cover", Cover, map[[2]int]color.RGBA{{12, 8}: red, {52, 8}: green, {12, 24}: blue, {52, 24}: yellow}},
		{"fill", FillBox, map[[2]int]color.RGBA{{20, 8}: red, {44, 8}: green, {20, 24}: blue, {44, 24}: yellow}},
		{"natural", NaturalSize, map[[2]int]color.RGBA{{24, 16}: blue, {8, 4}: white, {56, 28}: white, {40, 16}: yellow}},
	}
	for _, tc := range fits {
		t.Run(tc.name, func(t *testing.T) {
			tt := NewTester(func(c *Context) {
				// A 64x32 box at the window's origin.
				Image(c, apiFitBitmap()).Size(64, 32).Fit(tc.fit)
			}, 80, 50)
			defer tt.rt.close()
			tt.Frame()
			for pt, want := range tc.want {
				if got := tt.Image().RGBAAt(pt[0], pt[1]); got != want {
					t.Errorf("pixel %v = %v, want %v", pt, got, want)
				}
			}
		})
	}
}

func TestFitAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Fit(Cover).Valid() {
		t.Error("an absent element changed its fit mode")
	}
}
