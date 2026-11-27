package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

func TestDecodeBitmapScaled(t *testing.T) {
	// An 8×4 PNG: red on the left half, blue on the right half.
	img := image.NewRGBA(image.Rect(0, 0, 8, 4))
	for y := range 4 {
		for x := range 8 {
			c := color.RGBA{R: 255, A: 255}
			if x >= 4 {
				c = color.RGBA{B: 255, A: 255}
			}
			img.SetRGBA(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()

	// Shrinking preserves the aspect ratio and averages each region.
	bm, err := DecodeBitmapScaled(data, 4)
	if err != nil {
		t.Fatalf("DecodeBitmapScaled(4): %v", err)
	}
	if w, h := bm.Size(); w != 4 || h != 2 {
		t.Fatalf("size = %d×%d, want 4×2", w, h)
	}
	if got := bm.img.Pix[0:4]; !bytes.Equal(got, []byte{255, 0, 0, 255}) {
		t.Errorf("left pixel = %v, want red", got)
	}
	if got := bm.img.Pix[2*4 : 3*4]; !bytes.Equal(got, []byte{0, 0, 255, 255}) {
		t.Errorf("right pixel = %v, want blue", got)
	}

	// Shrinking to 1×1 averages both halves and rounds to the nearest byte.
	one, err := DecodeBitmapScaled(data, 1)
	if err != nil {
		t.Fatalf("DecodeBitmapScaled(1): %v", err)
	}
	if w, h := one.Size(); w != 1 || h != 1 {
		t.Fatalf("1×1 size = %d×%d, want 1×1", w, h)
	}
	if got := one.img.Pix[:4]; !bytes.Equal(got, []byte{128, 0, 128, 255}) {
		t.Errorf("averaged pixel = %v, want [128 0 128 255]", got)
	}

	// A larger maxSide or a non-positive one leaves the size unchanged.
	for _, maxSide := range []int{16, 8, 0, -4} {
		keep, err := DecodeBitmapScaled(data, maxSide)
		if err != nil {
			t.Fatalf("DecodeBitmapScaled(%d): %v", maxSide, err)
		}
		if w, h := keep.Size(); w != 8 || h != 4 {
			t.Errorf("maxSide %d: size = %d×%d, want 8×4", maxSide, w, h)
		}
	}

	// Unlike DecodeBitmap, DecodeBitmapScaled does not keep the encoded bytes,
	// so ReleasePixels leaves the pixels in place.
	bm.ReleasePixels()
	if len(bm.img.Pix) == 0 {
		t.Error("ReleasePixels dropped pixels of a scaled bitmap with no remake")
	}

	// Undecodable bytes fail with an error.
	if _, err := DecodeBitmapScaled([]byte("not an image"), 4); err == nil {
		t.Error("DecodeBitmapScaled accepted invalid bytes")
	}
}

func TestShrinkImage(t *testing.T) {
	// Already fits or non-positive maxSide returns the same image untouched.
	small := scene.NewImageRGBA(4, 3, make([]byte, 4*3*4))
	if out := shrinkImage(small, 4); out != small {
		t.Error("shrinkImage copied an image that already fits")
	}
	if out := shrinkImage(small, 0); out != small {
		t.Error("shrinkImage modified an image when maxSide is 0")
	}

	// Tall image (4×8): green on the top half, white on the bottom half.
	tallPix := make([]byte, 4*8*4)
	for y := range 8 {
		for x := range 4 {
			off := (y*4 + x) * 4
			if y < 4 {
				copy(tallPix[off:], []byte{0, 255, 0, 255})
			} else {
				copy(tallPix[off:], []byte{255, 255, 255, 255})
			}
		}
	}
	tall := shrinkImage(scene.NewImageRGBA(4, 8, tallPix), 4)
	if tall.W != 2 || tall.H != 4 {
		t.Fatalf("tall shrink size = %d×%d, want 2×4", tall.W, tall.H)
	}
	if got := tall.Pix[:4]; !bytes.Equal(got, []byte{0, 255, 0, 255}) {
		t.Errorf("top pixel = %v, want green", got)
	}
	if got := tall.Pix[3*2*4 : 3*2*4+4]; !bytes.Equal(got, []byte{255, 255, 255, 255}) {
		t.Errorf("bottom pixel = %v, want white", got)
	}

	// Extreme aspect ratios clamp the shorter side to 1 instead of 0.
	wide := shrinkImage(scene.NewImageRGBA(10, 1, make([]byte, 10*4)), 3)
	if wide.W != 3 || wide.H != 1 {
		t.Errorf("wide strip size = %d×%d, want 3×1", wide.W, wide.H)
	}
	narrow := shrinkImage(scene.NewImageRGBA(1, 10, make([]byte, 10*4)), 3)
	if narrow.W != 1 || narrow.H != 3 {
		t.Errorf("narrow strip size = %d×%d, want 1×3", narrow.W, narrow.H)
	}

	// Half-up rounding: sums [6, 10, 14, 1020] over 4 pixels round to [2, 3, 4, 255].
	pix := []byte{
		0, 1, 2, 255,
		1, 2, 3, 255,
		2, 3, 4, 255,
		3, 4, 5, 255,
	}
	avg := shrinkImage(scene.NewImageRGBA(2, 2, pix), 1)
	if got := avg.Pix[:4]; !bytes.Equal(got, []byte{2, 3, 4, 255}) {
		t.Errorf("rounded average = %v, want [2 3 4 255]", got)
	}
}
