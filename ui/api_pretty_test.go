package ui

import "testing"

// Pretty marks a text element to keep a word from standing alone on its
// last line; the layout behavior is covered by TestPrettyText.
func TestPrettyAPI(t *testing.T) {
	var marked bool
	tt := NewTester(func(c *Context) {
		text := Text(c, "A paragraph of text").Pretty()
		if !text.Valid() {
			t.Fatal("Pretty returned an invalid element")
		}
		marked = text.node().pretty
	}, 200, 80)
	defer tt.rt.close()
	tt.Frame()
	if !marked {
		t.Fatal("Pretty did not mark the text for pretty wrapping")
	}
}

func TestPrettyAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Pretty().Valid() {
		t.Error("an absent element enabled pretty wrapping")
	}
}
