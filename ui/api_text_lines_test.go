package ui

import "testing"

// The public API of the lines of a text: the lines of a text area that an
// app paints beside, the balance of the lines of a paragraph, and the
// height of the lines of a text.
func TestTextLinesAPI(t *testing.T) {
	// A text area keeps where it laid out its lines.
	var lines TextLines
	value := "one\ntwo\nthree"
	tt := NewTester(func(c *Context) {
		TextArea(c, &value).TrackLines(&lines)
	}, 300, 200)
	tt.Frame()
	if lines.Count() != 3 {
		t.Errorf("the text area has %d lines, want 3", lines.Count())
	}
	if lines.Height() <= 0 {
		t.Errorf("the lines of the text area have the height %v", lines.Height())
	}
	if lines.Top(1) <= 0 {
		t.Errorf("the second line of the text area starts at %v", lines.Top(1))
	}

	// A text of a fixed line height is as tall as the height it is given.
	fixed := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Text(c, "Hello").FixedLineHeight(30)
		})
	}, 300, 200)
	fixed.Frame()
	fr, ok := fixed.Find("Hello")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", fixed.Texts())
	}
	if fr.H != 30 {
		t.Errorf("the text of lines 30 high is %v high", fr.H)
	}

	// A balanced text takes lines of about the same length: its longest
	// line is shorter than the longest line of the text without it, and
	// it breaks its lines at other words.
	text := "the quick brown fox jumps over a lazy dog and then some more words to fill the line here"
	var balanced, plain *node
	bt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			balanced = Text(c, text).Width(137).Balance().node()
			plain = Text(c, text).Width(137).node()
		})
	}, 300, 400)
	bt.Frame()
	if balanced == nil || balanced.tl == nil || plain == nil || plain.tl == nil {
		t.Fatal("the texts have no layout")
	}
	if len(balanced.tl.Lines) != len(plain.tl.Lines) {
		t.Fatalf("the balanced text takes %d lines, the plain %d", len(balanced.tl.Lines), len(plain.tl.Lines))
	}
	if balanced.tl.Width >= plain.tl.Width {
		t.Errorf("the longest line of the balanced text is %v wide, the plain %v", balanced.tl.Width, plain.tl.Width)
	}
	same := true
	for i, l := range balanced.tl.Lines {
		if l.Start != plain.tl.Lines[i].Start {
			same = false
			break
		}
	}
	if same {
		t.Error("the balanced text breaks its lines at the same words")
	}

	// Outside a frame the view lays out no line.
	if (Element{}).Balance().Valid() {
		t.Error("a text was balanced outside a frame")
	}
	if (Element{}).FixedLineHeight(30).Valid() {
		t.Error("a text was styled outside a frame")
	}
	if (Element{}).TrackLines(&lines).Valid() {
		t.Error("a text area tracked its lines outside a frame")
	}
}
