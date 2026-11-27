package ui

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// indexAt returns the rune the caret would take at a point: the click of
// a press lays the caret on the character the point is over.
func TestTextAreaIndexAt(t *testing.T) {
	s := "alpha\nbeta\ngamma"
	tt := coreNewTester(func(c *context) { coreTextArea(c, &s).Fill() }, 400, 300)
	tt.Press(20, 15)
	tt.Release(20, 15)
	a := textAreaState(tt).editor.area

	// The click above put the caret on the first paragraph.
	if got := a.indexAt(textAreaState(tt).editor, 20, 15); got >= 6 {
		t.Errorf("indexAt on the first line gave %d, want a rune of \"alpha\"", got)
	}

	// A point on the second line gives a rune past the first newline.
	if got := a.indexAt(textAreaState(tt).editor, 20, 15+18); got < 6 || got > 10 {
		t.Errorf("indexAt on the second line gave %d, want a rune of \"beta\"", got)
	}
}

// forget drops the layouts of the paragraphs far from view, keeping their
// heights, so a long document does not keep every layout it showed.
func TestTextAreaForgetsLayouts(t *testing.T) {
	var b strings.Builder
	for i := range 2000 {
		b.WriteString("line ")
		b.WriteString(strings.Repeat("x", i%7))
		b.WriteByte('\n')
	}
	s := b.String()
	tt := coreNewTester(func(c *context) { coreTextArea(c, &s).Fill() }, 400, 300)
	tt.Press(20, 15)
	tt.Release(20, 15)
	tt.Key(Ctrl, KeyHome)
	st := textAreaState(tt)
	a := st.editor.area

	// Scroll in steps, as a wheel does, so that the paragraphs the view
	// passed keep their layouts until there are too many of them.
	for range 120 {
		tt.Scroll(100, 100, 0, 100)
		tt.Frame()
	}
	if a.first == 0 {
		t.Fatal("the view did not scroll into the text")
	}

	// The count of laid paragraphs stays within the window and its
	// margin; layouts beyond it were dropped.
	if a.laid > maxLaid+a.last-a.first {
		t.Errorf("%d paragraphs keep their layouts, want at most the window and its margin", a.laid)
	}
	// Most paragraphs far above the view lost theirs (the caret's keeps
	// its, by reveal, and with it a window around it): only a few may
	// remain of the hundreds the view passed.
	kept := 0
	for i := 1; i < a.first-maxLaid/4 && i < len(st.editor.buf.paras); i++ {
		if st.editor.buf.paras[i].layout != nil {
			kept++
		}
	}
	if kept > 3*(maxLaid/4) {
		t.Errorf("%d paragraphs far above the view %d..%d keep their layouts", kept, a.first, a.last)
	}
	// A paragraph in view keeps its layout.
	if st.editor.buf.paras[a.last].layout == nil {
		t.Errorf("paragraph %d, in view, lost its layout", a.last)
	}
	// The paragraphs in view keep theirs.
	for i := a.first; i <= a.last && i < len(st.editor.buf.paras); i++ {
		if st.editor.buf.paras[i].layout == nil {
			t.Errorf("paragraph %d, in view %d..%d, lost its layout", i, a.first, a.last)
			break
		}
	}

	// The text is unchanged after the layouts came and went.
	if got := utf8.RuneCountInString(s); st.editor.buf.n != got {
		t.Errorf("the buffer holds %d runes, want %d", st.editor.buf.n, got)
	}
}
