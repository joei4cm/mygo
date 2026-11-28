package ui

import (
	"image/color"
	"testing"
)

// fakePlatformView stands for a view of the platform, whose handle is h.
type fakePlatformView struct{ h uintptr }

func (f fakePlatformView) NativeView() uintptr { return f.h }

// rgbaOfColor is the color as the pixels of the window hold it.
func rgbaOfColor(c Color) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, c.A}
}

// The public API of a view of the platform: one with no implementation is
// a plain box, and one with an implementation takes the focus.
func TestNativeViewAPI(t *testing.T) {
	red := RGB(255, 0, 0)

	// A view of the platform with no implementation shows in the place of
	// a plain box: it paints its background and holds its children. A
	// closed view is no view at all, and shows as one.
	for _, v := range []PlatformView{nil, fakePlatformView{h: 0}} {
		tt := NewTester(func(c *Context) {
			NativeView(c, v).Size(100, 60).Background(red).Children(func() {
				Text(c, "Inside")
			})
		}, 200, 200)
		if !tt.HasText("Inside") {
			t.Fatalf("the texts are %q", tt.Texts())
		}
		if got := tt.Image().RGBAAt(5, 5); got != rgbaOfColor(red) {
			t.Fatalf("the box painted %v", got)
		}
		if got := tt.Image().RGBAAt(150, 150); got != (color.RGBA{255, 255, 255, 255}) {
			t.Fatalf("the box painted past its own box: %v", got)
		}
	}

	// A view of the platform has no size of its own: it takes the room of
	// its row.
	var bounds Rect
	row := NewTester(func(c *Context) {
		Row(c).Fill().Children(func() {
			Box(c).Width(50).Height(80)
			bounds = NativeView(c, nil).Bounds()
		})
	}, 200, 200)
	row.Frame()
	if bounds.H != 200 {
		t.Fatalf("the view of the platform takes %v", bounds)
	}

	// A view of the platform shows in the place of the box, and takes the
	// focus as any view of the platform does.
	view := NewTester(func(c *Context) {
		NativeView(c, fakePlatformView{h: 1}).Size(100, 60).Background(red).Children(func() {
			Text(c, "Inside")
		})
	}, 200, 200)
	view.Key(0, KeyTab)
	if !view.Focused("Inside") {
		t.Fatal("the view of the platform took no focus")
	}
	if got := view.Image().RGBAAt(5, 5); got == rgbaOfColor(red) {
		t.Fatal("the box painted its background over the view of the platform")
	}

	// Outside a frame the view builds no view of the platform.
	p := view.rt.public
	if NativeView(p, nil).Valid() {
		t.Error("a view of the platform was built outside a frame")
	}
}
