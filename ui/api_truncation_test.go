package ui

import (
	"image"
	"testing"
)

// truncBox shows a long line of text in a narrow box, cut as style says,
// and returns the box it takes and the pixels it draws.
func truncBox(t *testing.T, style func(e Element) Element) (Rect, *image.RGBA) {
	t.Helper()
	tt := NewTester(func(c *Context) {
		Box(c).Width(80).Children(func() {
			e := Text(c, "src/components/picker.go").SingleLine()
			if style != nil {
				e = style(e)
			}
		})
	}, 200, 100)
	tt.Frame()
	r, ok := tt.Find("src/components/picker.go")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", tt.Texts())
	}
	return r, tt.Image()
}

// sameBox reports whether two pictures look the same in a box.
func sameBox(a, b *image.RGBA, r Rect) bool {
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if a.RGBAAt(x, y) != b.RGBAAt(x, y) {
				return false
			}
		}
	}
	return true
}

// The public API of how a text that does not fit is cut: at its end, in
// its middle, at its start, and with a mark of the app's choosing.
func TestTruncationAPI(t *testing.T) {
	end, endImg := truncBox(t, nil)

	// The middle keeps the start and the end of the line.
	middle, middleImg := truncBox(t, func(e Element) Element { return e.TruncateMiddle() })
	if sameBox(endImg, middleImg, end) {
		t.Error("the middle cut looks the end cut")
	}

	// The start keeps the end of the line.
	_, startImg := truncBox(t, func(e Element) Element { return e.TruncateStart() })
	if sameBox(endImg, startImg, end) {
		t.Error("the start cut looks the end cut")
	}
	if sameBox(middleImg, startImg, middle) {
		t.Error("the start cut looks the middle cut")
	}

	// The cut ends with a mark of the app's choosing.
	_, markImg := truncBox(t, func(e Element) Element { return e.Ellipsis("~") })
	if sameBox(endImg, markImg, end) {
		t.Error("the mark of the cut changed nothing")
	}

	// Outside a frame the view cuts no text.
	if (Element{}).TruncateMiddle().Valid() {
		t.Error("a text was cut outside a frame")
	}
	if (Element{}).TruncateStart().Valid() {
		t.Error("a text was cut outside a frame")
	}
	if (Element{}).Ellipsis("~").Valid() {
		t.Error("a text was cut outside a frame")
	}
}
