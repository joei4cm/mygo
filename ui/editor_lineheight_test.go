package ui

import (
	"testing"
	"time"

	"github.com/egoist/mygo/internal/text"
)

// lineHeight comes from the area's empty line while there is one, then
// from the first line of the layout, and is 18 without either.
func TestEditorLineHeight(t *testing.T) {
	sys := text.New()

	// No area, no layout: the default.
	ed := newEditor(func() time.Time { return time.Unix(0, 0) }, sys)
	if got := ed.lineHeight(); got != 18 {
		t.Errorf("lineHeight = %v, want the default 18", got)
	}

	// An area whose line has a height answers with it.
	ed.area = &area{}
	ed.area.line.Height = 21.5
	if got := ed.lineHeight(); got != 21.5 {
		t.Errorf("with an area, lineHeight = %v, want the area line's 21.5", got)
	}

	// An area whose line has none falls back to the layout's first line.
	ed.area.line.Height = 0
	l := sys.Shape(text.Params{Text: "two\nlines", Style: text.Style{Size: 14}})
	if n := len(l.Lines); n != 2 {
		t.Fatalf("the layout has %d lines, want 2", n)
	}
	ed.layout = l
	if got := ed.lineHeight(); got != l.Lines[0].Height {
		t.Errorf("with a layout, lineHeight = %v, want the first line's %v", got, l.Lines[0].Height)
	}
	if l.Lines[0].Height == 0 {
		t.Error("the layout's first line has no height; the test says nothing")
	}
}
