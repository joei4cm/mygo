package ui

import (
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

// The public API of the relations assistive technology reads: a field
// named by the text built before it, a control whose value is read beside
// its name, and a field that has the focus for the option of a list of
// its own that its keys are on.
func TestAccessRelationsAPI(t *testing.T) {
	mail, reply := "", ""
	var field Element
	var option Element
	tt := NewTester(func(c *Context) {
		Column(c).Gap(8).Children(func() {
			label := Text(c, "Email address")
			TextInput(c, &mail).Width(200).LabelledBy(label)
			Button(c, "Zoom").Value("100%")
			field = TextInput(c, &reply).Width(200).Label("Reply to")
			Column(c).Role(RoleList).Label("own-list").Children(func() {
				Row(c).Role(RoleListItem).Label("Alpha").Children(func() { Text(c, "Alpha") })
				option = Row(c).Role(RoleListItem).Label("Beta").Children(func() { Text(c, "Beta") })
				Row(c).Role(RoleListItem).Label("Gamma").Children(func() { Text(c, "Gamma") })
			})
			field.ActiveDescendant(option)
		})
	}, 300, 220)
	defer tt.rt.close()
	tt.Click("Reply to")
	tt.send(platform.SurfaceEvent{Kind: platform.AccessibilityOn})
	a, ok := tt.AccessByRole(RoleTextField, "Email address")
	if !ok {
		t.Fatal("the field is not named by the text built before it")
	}
	if a.Value != "" {
		t.Errorf("the field reads %q as its value", a.Value)
	}
	if b, ok := tt.AccessByRole(RoleButton, "Zoom"); !ok || b.Value != "100%" {
		t.Errorf("the button's value %q, want 100%%", b.Value)
	}
	// The option the field's keys are on has the focus for it.
	focusID := uint64(0)
	for i := range tt.h.access.Nodes {
		if n := &tt.h.access.Nodes[i]; n.Label == "Beta" && n.Role == platform.RoleListItem {
			focusID = n.ID
		}
	}
	if focusID == 0 || tt.h.access.Focus != focusID {
		t.Errorf("the focus is on %d, want the option Beta %d", tt.h.access.Focus, focusID)
	}
	// The field keeps the keyboard, but reads as having it for the option.
	if r, ok := tt.AccessByRole(RoleTextField, "Reply to"); !ok || r.Focused {
		t.Errorf("the reply field itself focused %v", ok)
	}
}

// Outside a frame these relations change nothing.
func TestAccessRelationsAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Value("3 of 10").Valid() {
		t.Error("an absent element took a value")
	}
	if absent.LabelledBy(absent).Valid() {
		t.Error("an absent element took a label")
	}
	if absent.ActiveDescendant(absent).Valid() {
		t.Error("an absent element took an active descendant")
	}
}
