package ui

import (
	"testing"
)

// The public API of a tree: its items sit at their depth, and the items
// inside one show below it while it is open.
func TestTreePartsAPI(t *testing.T) {
	srcOpen := true
	var (
		depths  []int
		grouped bool
	)
	tt := NewTester(func(c *Context) {
		tree := TreeBase(c)
		tree.Root.Children(func() {
			item := tree.Item(&srcOpen)
			depths = append(depths[:0], item.Depth())
			item.Row.Children(func() { Text(c, "src") })
			grouped = false
			if item.Group(func() {
				leaf := tree.Item(nil)
				depths = append(depths, leaf.Depth())
				leaf.Row.Children(func() { Text(c, "main.go") })
			}) != (Element{}) {
				grouped = true
			}
		})
	}, 400, 300)
	tt.Frame()
	if len(depths) != 2 || depths[0] != 0 || depths[1] != 1 {
		t.Errorf("the depths are %v", depths)
	}
	if !tt.HasText("main.go") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Closed, the item holds nothing.
	srcOpen = false
	tt.Frame()
	if grouped {
		t.Error("the item holds something while it is closed")
	}
	if tt.HasText("main.go") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A part that is not the view's tells nothing.
	var none TreeParts
	var noneItem TreeItemParts
	if none.Item(nil).Depth() != 0 || noneItem.Depth() != 0 || noneItem.Group(func() {}) != (Element{}) {
		t.Error("a part that is not the view's answers")
	}
}

// The public API of a segmented control: its segments choose the one
// that is chosen.
func TestSegmentedPartsAPI(t *testing.T) {
	selected := 0
	tt := NewTester(func(c *Context) {
		s := SegmentedBase(c, &selected, 2)
		s.Track.Children(func() {
			for i, name := range []string{"List", "Grid"} {
				segment := s.Segment(i)
				segment.Children(func() { Text(c, name) })
			}
		})
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("List") || !tt.HasText("Grid") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	tt.Click("Grid")
	tt.Frame()
	if selected != 1 {
		t.Errorf("the chosen segment is %d", selected)
	}

	// A part that is not the view's tells nothing.
	var none SegmentedParts
	if none.Segment(0) != (Element{}) {
		t.Error("a part that is not the view's answers")
	}
}
