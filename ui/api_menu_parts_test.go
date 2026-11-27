package ui

import (
	"slices"
	"testing"
)

// The public API of a menu: the menu shows while it is open, its items
// report their choice, and it groups them and holds submenus.
func TestMenuPartsAPI(t *testing.T) {
	var open bool
	var (
		shown     bool
		chosen    []string
		itemRole  Role
		checkRole Role
		radRole   Role
	)
	tt := NewTester(func(c *Context) {
		m := MenuBase(c, &open)
		m.Trigger.Children(func() { Text(c, "Actions") })
		m.Popup(func(panel Element) {
			item := m.Item()
			itemRole = c.rt.states[item.node().ID()].role
			item.Children(func() { Text(c, "Open") })
			if item.Clicked() {
				chosen = append(chosen, "Open")
			}
			check := m.CheckItem(false)
			checkRole = c.rt.states[check.node().ID()].role
			check.Children(func() { Text(c, "Pinned") })
			if check.Clicked() {
				chosen = append(chosen, "Pinned")
			}
			radio := m.RadioItem(true)
			radRole = c.rt.states[radio.node().ID()].role
			radio.Children(func() { Text(c, "Name") })
			if radio.Clicked() {
				chosen = append(chosen, "Name")
			}
			sortBy := m.SubmenuItem()
			sortBy.Children(func() { Text(c, "Sort by") })
			m.Submenu(sortBy, func(panel Element) {
				panel.Children(func() {
					by := m.RadioItem(false)
					by.Children(func() { Text(c, "Date") })
					if by.Clicked() {
						chosen = append(chosen, "Date")
					}
				})
			})
			m.Group(func(group Element) {
				group.Children(func() { Text(c, "Grouped") })
			})
		})
		shown = m.Open()
	}, 400, 400)
	tt.Frame()
	if shown {
		t.Error("the menu shows before it is opened")
	}

	// A click on the trigger opens the menu, and its items show.
	tt.Click("Actions")
	tt.Frame()
	if !shown {
		t.Error("the menu does not show")
	}
	for _, text := range []string{"Open", "Pinned", "Name", "Sort by", "Grouped"} {
		if !tt.HasText(text) {
			t.Errorf("the item %q shows nowhere among %q", text, tt.Texts())
		}
	}

	// The items are what assistive technology sees them as.
	if itemRole != RoleMenuItem || checkRole != RoleMenuItemCheckBox || radRole != RoleMenuItemRadio {
		t.Errorf("the roles are %v, %v, %v", itemRole, checkRole, radRole)
	}

	// A click on an item chooses it, and the menu closes.
	tt.Click("Open")
	tt.Frame()
	if !slices.Equal(chosen, []string{"Open"}) {
		t.Errorf("the choices are %q", chosen)
	}
	if shown {
		t.Error("the menu shows after an item was chosen")
	}

	// The pointer resting on an item with a submenu opens it.
	tt.Click("Actions")
	tt.Frame()
	sortBy, ok := tt.Find("Sort by")
	if !ok {
		t.Fatalf("the item shows nowhere among %q", tt.Texts())
	}
	tt.Move(sortBy.X+sortBy.W/2, sortBy.Y+sortBy.H/2)
	tt.Advance(menuAimGrace * 2)
	if !tt.HasText("Date") {
		t.Errorf("the submenu shows nowhere among %q", tt.Texts())
	}

	// A part that is not the view's tells nothing.
	var none MenuParts
	none.Submenu(Element{}, func(Element) {})
	if none.Open() {
		t.Error("a part that is not the view's answers")
	}
	if none.Item() != (Element{}) || none.CheckItem(false) != (Element{}) ||
		none.RadioItem(false) != (Element{}) || none.SubmenuItem() != (Element{}) {
		t.Error("a part that is not the view's answers")
	}
	if none.Group(func(Element) {}) != (Element{}) || none.Popup(func(Element) {}) != (Element{}) {
		t.Error("a part that is not the view's answers")
	}
}
