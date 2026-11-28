package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"testing"

	"github.com/egoist/mygo/internal/wait"
)

// The public API of a video: it shows the frames of its player.
func TestVideoAPI(t *testing.T) {
	red := color.RGBA{255, 0, 0, 255}
	g := &gif.GIF{Config: image.Config{Width: 4, Height: 4}}
	g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, 4, 4), color.Palette{red}))
	g.Delay = append(g.Delay, 10)
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		t.Fatal(err)
	}
	d, err := NewGIFDecoder(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	p, err := NewVideoPlayer(d)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	tt := NewTester(func(c *Context) {
		Video(c, p).Size(16, 16).Fit(FillBox)
	}, 40, 40)
	wait.Until(t, "the video to decode its frame", func() bool {
		return p.frames != nil && len(p.frames) >= 1
	})
	tt.Frame()
	if got := tt.Image().RGBAAt(8, 8); got != red {
		t.Fatalf("the video shows %v", got)
	}

	// A video of no player shows nothing, and asks for nothing.
	empty := NewTester(func(c *Context) {
		Video(c, nil).Size(16, 16)
	}, 40, 40)
	empty.Frame()
	if got := empty.Image().RGBAAt(8, 8); got == red {
		t.Fatalf("the video of no player shows %v", got)
	}

	// Outside a frame the view builds no video.
	q := tt.rt.public
	if Video(q, p).Valid() {
		t.Error("a video was built outside a frame")
	}
}
