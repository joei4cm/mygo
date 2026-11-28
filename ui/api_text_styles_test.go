package ui

import (
	"image"
	"testing"
)

// textBox shows a text styled by style, and returns the box it takes and
// the pixels it draws.
func textBox(t *testing.T, style func(e Element) Element) (Rect, *image.RGBA) {
	t.Helper()
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			e := Text(c, "Hello")
			if style != nil {
				e = style(e)
			}
		})
	}, 300, 200)
	r, ok := tt.Find("Hello")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", tt.Texts())
	}
	return r, tt.Image()
}

// ink counts the dark pixels of a box of a picture.
func ink(img *image.RGBA, r Rect) int {
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if c := img.RGBAAt(x, y); c.R < 128 && c.A > 0 {
				n++
			}
		}
	}
	return n
}

// The public API of the styles of a text: its weight, the height of its
// lines, the space between its letters, and the line under it.
func TestTextStylesAPI(t *testing.T) {
	plain, plainImg := textBox(t, nil)

	// A heavier weight takes more room.
	bold, _ := textBox(t, func(e Element) Element { return e.FontWeight(700) })
	if bold.W <= plain.W {
		t.Errorf("the bold text is %v wide, the plain %v", bold.W, plain.W)
	}

	// The height of the lines takes more room.
	tall, _ := textBox(t, func(e Element) Element { return e.LineHeight(2) })
	if tall.H <= plain.H {
		t.Errorf("the text of taller lines is %v high, the plain %v", tall.H, plain.H)
	}

	// The space between the letters takes more room.
	spaced, _ := textBox(t, func(e Element) Element { return e.LetterSpacing(2) })
	if spaced.W <= plain.W {
		t.Errorf("the spaced text is %v wide, the plain %v", spaced.W, plain.W)
	}

	// The underline draws a line under the text.
	under, underImg := textBox(t, func(e Element) Element { return e.Underline() })
	if ink(underImg, under) <= ink(plainImg, plain) {
		t.Error("the underline draws no line")
	}

	// An italic font draws the letters leaning.
	italic, italicImg := textBox(t, func(e Element) Element { return e.Italic() })
	if same(italicImg, plainImg, italic) {
		t.Error("the italic text looks the plain one")
	}

	// Outside a frame the view styles no text.
	if (Element{}).FontWeight(700).Valid() {
		t.Error("a text was styled outside a frame")
	}
	if (Element{}).Italic().Valid() {
		t.Error("a text was styled outside a frame")
	}
	if (Element{}).LetterSpacing(2).Valid() {
		t.Error("a text was styled outside a frame")
	}
	if (Element{}).LineHeight(2).Valid() {
		t.Error("a text was styled outside a frame")
	}
	if (Element{}).Underline().Valid() {
		t.Error("a text was styled outside a frame")
	}
}

// same reports whether two pictures look the same in a box.
func same(a, b *image.RGBA, r Rect) bool {
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if a.RGBAAt(x, y) != b.RGBAAt(x, y) {
				return false
			}
		}
	}
	return true
}
