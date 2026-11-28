package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// The public API of describing an element to assistive technology, and of
// keeping it out of what assistive technology reads.
func TestAccessAPIDescribe(t *testing.T) {
	var filled bool
	tt := NewTester(func(c *Context) {
		Box(c).Size(120, 40).Role(RoleButton).Label("Save").Access(func(n *platform.AccessNode) {
			filled = true
			n.Description = "Writes the file"
		})
	}, 200, 100)
	tt.Frame()
	if filled {
		t.Fatal("the description ran while no assistive technology read the window")
	}
	tt.send(platform.SurfaceEvent{Kind: platform.AccessibilityOn})
	tt.Frame()
	if !filled {
		t.Fatal("the description never ran")
	}
	info, ok := tt.Access("Save")
	if !ok {
		t.Fatal("no button named Save")
	}
	if info.Description != "Writes the file" {
		t.Errorf("the button's description is %q", info.Description)
	}
}

// An element hidden from assistive technology still shows, and neither it
// nor what it holds is in the tree.
func TestAccessAPIHidden(t *testing.T) {
	hide := true
	tt := NewTester(func(c *Context) {
		Column(c).Gap(8).Children(func() {
			Box(c).Width(120).Height(30).Children(func() {
				Text(c, "Secret")
			}).AccessHidden(hide)
			Text(c, "Shown")
		})
	}, 200, 100)
	tt.Frame()
	if !tt.HasText("Secret") {
		t.Fatal("the hidden element did not show")
	}
	tt.send(platform.SurfaceEvent{Kind: platform.AccessibilityOn})
	tt.Frame()
	if _, ok := tt.Access("Secret"); ok {
		t.Error("the hidden element is in the tree")
	}
	if _, ok := tt.Access("Shown"); !ok {
		t.Error("the element beside the hidden one is not in the tree")
	}
	// Shown again, it is read once more.
	hide = false
	tt.Frame()
	if _, ok := tt.Access("Secret"); !ok {
		t.Error("the element shown again is not in the tree")
	}
}

// Outside a frame no element is described, nor kept out of the tree.
func TestAccessAPIOutsideFrame(t *testing.T) {
	var e Element
	if e.Access(func(n *platform.AccessNode) { n.Label = "x" }).Valid() {
		t.Error("an absent element was described")
	}
	if e.AccessHidden(true).Valid() {
		t.Error("an absent element was hidden")
	}
}
