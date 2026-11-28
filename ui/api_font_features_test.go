package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/text"
)

// FontFeatures turns on OpenType features of the font, by tag or as
// tag=value. The text of the element and of the elements inside it takes
// them; a text that sets its own takes its own instead.
func TestFontFeaturesAPI(t *testing.T) {
	var params, own text.Params
	tt := NewTester(func(c *Context) {
		Row(c).FontFeatures("tnum", "liga=0").Children(func() {
			Text(c, "12:30").Label("inherited")
			params = Text(c, "12:30").node().textParams(0)
			own = Text(c, "12:30").FontFeatures("smcp").node().textParams(0)
		})
	}, 400, 200)
	defer tt.rt.close()
	if s := params.Style; s.Features != "tnum,liga=0" {
		t.Errorf("the text inside takes %q, not %q", s.Features, "tnum,liga=0")
	}
	if s := own.Style; s.Features != "smcp" {
		t.Errorf("the text with its own takes %q, not %q", s.Features, "smcp")
	}
	if _, ok := tt.Find("inherited"); !ok {
		t.Fatal("the row shows nothing")
	}
}

func TestFontFeaturesAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.FontFeatures("tnum").Valid() {
		t.Error("an absent element turned a font feature on")
	}
}
