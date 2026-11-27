package ui

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/internal/imagesrc"
	"github.com/egoist/mygo/internal/scene"
)

// FillRadii fills a rectangle whose corners have radii of their own,
// clockwise from the top left.
func TestFillRadii(t *testing.T) {
	white := [4]uint8{255, 255, 255, 255}
	tt := NewTester(func(c *Context) {
		Box(c).Size(100, 100).Draw(func(p *Painter, r Rect) {
			p.FillRadii(Rect{X: 10, Y: 10, W: 80, H: 80}, [4]float32{40, 0, 0, 0}, gRed)
			// A radius on every corner fills like Fill does with one.
			p.FillRadii(Rect{X: 10, Y: 10, W: 80, H: 80}, [4]float32{20, 20, 20, 20}, gBlue)
		})
	}, 100, 100)
	img := tt.Image()
	// The blue fill covers the red one; its rounded corners show.
	for _, c := range []struct {
		x, y int
		want [4]uint8
	}{
		{11, 11, white}, // outside the corner of 20
		{30, 11, [4]uint8{0, 0, 255, 255}},
		{70, 70, [4]uint8{0, 0, 255, 255}}, // inside the bottom right corner of 20
	} {
		if got := rgba(img, c.x, c.y); !nearRGBA8(got, c.want, 0) {
			t.Errorf("at (%d, %d) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
	// Alone, the red fill's corners of 40 on the top left and none
	// elsewhere show: white past the corner, red on the square ones.
	tt = NewTester(func(c *Context) {
		Box(c).Size(100, 100).Draw(func(p *Painter, r Rect) {
			p.FillRadii(Rect{X: 10, Y: 10, W: 80, H: 80}, [4]float32{40, 0, 0, 0}, gRed)
		})
	}, 100, 100)
	img = tt.Image()
	for _, c := range []struct {
		x, y int
		want [4]uint8
	}{
		{11, 11, white},                    // outside the corner of 40
		{50, 11, [4]uint8{255, 0, 0, 255}}, // past the corner, on the top edge
		{11, 89, [4]uint8{255, 0, 0, 255}}, // the bottom left corner is square
		{89, 11, [4]uint8{255, 0, 0, 255}}, // as the top right one
	} {
		if got := rgba(img, c.x, c.y); !nearRGBA8(got, c.want, 0) {
			t.Errorf("at (%d, %d) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
}

// FillConicGradient fills a rectangle with a conic gradient of its own,
// rounded as the radius says.
func TestFillConicGradient(t *testing.T) {
	stops := []GradientStop{{gRed, 0}, {gGreen, 0.25}, {gBlue, 0.5}, {gRed, 1}}
	tt := NewTester(func(c *Context) {
		Box(c).Size(100, 100).Draw(func(p *Painter, r Rect) {
			p.FillConicGradient(Rect{X: 10, Y: 10, W: 80, H: 80}, ConicGradient{Stops: stops, Angle: 90}, 0)
		})
	}, 100, 100)
	img := tt.Image()
	// From the right, a quarter in, at the bottom, is green; halfway, at
	// the left, is blue.
	for _, c := range []struct {
		x, y int
		want [4]uint8
	}{
		{89, 50, [4]uint8{255, 0, 0, 255}},
		{50, 89, [4]uint8{0, 255, 0, 255}},
		{10, 50, [4]uint8{0, 0, 255, 255}},
	} {
		if got := rgba(img, c.x, c.y); !nearRGBA8(got, c.want, 8) {
			t.Errorf("at (%d, %d) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
	// Rounded by 40, the corners show the background.
	tt = NewTester(func(c *Context) {
		Box(c).Size(100, 100).Draw(func(p *Painter, r Rect) {
			p.FillConicGradient(Rect{X: 10, Y: 10, W: 80, H: 80}, ConicGradient{Stops: stops}, 40)
		})
	}, 100, 100)
	if got := tt.Image().RGBAAt(11, 11); got != (color.RGBA{255, 255, 255, 255}) {
		t.Errorf("outside the corner of 40 = %v, want the background", got)
	}
}

// fakeTextureFrames is an imagesrc.Frames of one picture, the one a
// texture-backed Image shows.
type fakeTextureFrames struct {
	imagesrc.Mark
	img        *scene.Image
	size       [2]float32
	key        any
	gotRelease bool
	pending    bool // a frame the window has not asked for yet
}

func (f *fakeTextureFrames) Size() (w, h float32) { return f.size[0], f.size[1] }

func (f *fakeTextureFrames) Frame(key any, invalidate func()) *scene.Image {
	f.key, f.gotRelease = key, invalidate != nil
	if f.pending {
		f.pending = false
		invalidate()
	}
	return f.img
}

func solidPix(w, h int, r, g, b uint8) []byte {
	pix := make([]byte, w*h*4)
	for i := 0; i < len(pix); i += 4 {
		pix[i], pix[i+1], pix[i+2], pix[i+3] = r, g, b, 255
	}
	return pix
}

// An Image of a texture draws the frame it presented, and the frame after.
func TestDrawTexture(t *testing.T) {
	white := color.RGBA{255, 255, 255, 255}
	frames := &fakeTextureFrames{img: scene.NewImageRGBA(8, 8, solidPix(8, 8, 0, 255, 0)), size: [2]float32{8, 8}}
	tt := NewTester(func(c *Context) {
		Image(c, frames).Size(32, 32)
	}, 100, 100)
	img := tt.Image()
	if got := img.RGBAAt(16, 16); got != (color.RGBA{0, 255, 0, 255}) {
		t.Errorf("the middle of the texture %v, want green", got)
	}
	if got := img.RGBAAt(2, 2); got != (color.RGBA{0, 255, 0, 255}) {
		t.Errorf("inside the image %v, want green", got)
	}
	if got := img.RGBAAt(60, 16); got != white {
		t.Errorf("outside the image %v, want the background", got)
	}
	if frames.key == nil {
		t.Error("Frame got no key")
	}
	if !frames.gotRelease {
		t.Error("Frame got no invalidate")
	}
	// The next frame it presents is the one that shows; asking the
	// window for it is what its invalidate does.
	frames.img = scene.NewImageRGBA(8, 8, solidPix(8, 8, 255, 0, 0))
	frames.pending = true
	tt.Frame()
	if got := tt.Image().RGBAAt(16, 16); got != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("the middle of the next texture %v, want red", got)
	}

	// Without a frame yet, nothing shows.
	none := &fakeTextureFrames{size: [2]float32{8, 8}}
	tt = NewTester(func(c *Context) {
		Image(c, none).Size(32, 32)
	}, 100, 100)
	if got := tt.Image().RGBAAt(16, 16); got != white {
		t.Errorf("without a frame %v, want the background", got)
	}
	// Nor with an empty one.
	empty := &fakeTextureFrames{img: scene.NewImageRGBA(0, 0, nil)}
	tt = NewTester(func(c *Context) {
		Image(c, empty).Size(32, 32)
	}, 100, 100)
	if got := tt.Image().RGBAAt(16, 16); got != white {
		t.Errorf("with an empty frame %v, want the background", got)
	}
	// Nor where the image has no room.
	tt = NewTester(func(c *Context) {
		Image(c, frames).Size(0, 32)
	}, 100, 100)
	if got := tt.Image().RGBAAt(0, 16); got != white {
		t.Errorf("with no room %v, want the background", got)
	}
}
