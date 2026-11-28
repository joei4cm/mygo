package ui

import "testing"

// Lines sets the minimum and maximum line counts used to size a growing text area.
func TestLinesAPI(t *testing.T) {
	value := "one\ntwo"
	min, max := 2, 5
	var got [2]int
	tt := NewTester(func(c *Context) {
		in := TextArea(c, &value).Lines(min, max)
		if !in.Valid() {
			t.Fatal("TextArea is invalid")
		}
		got = in.node().st.editor.lines
	}, 220, 140)
	defer tt.rt.close()
	for _, tc := range [][2]int{{1, 4}, {3, 6}, {2, 5}} {
		min, max = tc[0], tc[1]
		tt.Frame()
		if got != tc {
			t.Errorf("Lines(%d, %d) set %v", tc[0], tc[1], got)
		}
	}
}

func TestLinesAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Lines(1, 4).Valid() {
		t.Error("an absent element set line limits")
	}
}
