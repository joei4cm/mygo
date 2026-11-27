package ui

import "testing"

// Depth tells how deep the item of a row sits: the items at the top at 0,
// their open children a level under them, and so on down.
func TestOutlineDepth(t *testing.T) {
	var s OutlineState[string]
	tt := coreNewTester(func(c *context) {
		coreOutline(c, &s, []string{"src", "docs", "README.md"}, children, func(path string) {
			coreText(c, baseName(path))
		}).Grow(1)
	}, 400, 400)
	tt.SetPreferences(Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if s.Rows() != 3 {
		t.Fatalf("%d rows closed", s.Rows())
	}
	for row, want := range map[int]int{0: 0, 1: 0, 2: 0} {
		if got := s.Depth(row); got != want {
			t.Errorf("closed: the depth of row %d is %d, want %d", row, got, want)
		}
	}
	// Open src and the ui inside it: their children sit a level deeper
	// each, the roots at the top stay at 0.
	src, _ := tt.Find("src")
	tt.ClickAt(src.X-10, src.Y+src.H/2)
	ui, _ := tt.Find("ui")
	tt.ClickAt(ui.X-10, ui.Y+ui.H/2)
	tt.Frame()
	if s.Rows() != 7 {
		t.Fatalf("%d rows open", s.Rows())
	}
	for _, want := range []struct {
		item  string
		depth int
	}{
		{"src", 0},
		{"src/ui", 1},
		{"src/ui/list.go", 2},
		{"src/ui/tree.go", 2},
		{"src/main.go", 1},
		{"docs", 0},
		{"README.md", 0},
	} {
		depth, at := -1, -1
		for row := range s.Rows() {
			if s.Item(row) == want.item {
				depth, at = s.Depth(row), row
				break
			}
		}
		if at < 0 {
			t.Errorf("%q shows nowhere", want.item)
		} else if depth != want.depth {
			t.Errorf("%q sits at depth %d, want %d", want.item, depth, want.depth)
		}
	}
}
