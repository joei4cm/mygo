package ui

import (
	"testing"
)

// The public API of a menu button: the button opens a menu with the items
// the app adds, and choosing one reports the choice.
func TestMenuButtonAPI(t *testing.T) {
	var chosen string
	tt := NewTester(func(c *Context) {
		MenuButton(c, "Sort by", func(m *Menu) {
			for _, by := range []string{"Name", "Date"} {
				if m.Item(by).Chosen() {
					chosen = by
				}
			}
		})
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("Sort by") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A click on the button opens its menu, with the items of the app.
	if err := tt.Click("Sort by"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	items := tt.Menu()
	if len(items) != 2 || items[0] != "Name" || items[1] != "Date" {
		t.Errorf("the menu holds %q", items)
	}

	// Choosing an item runs the app's function again, which reports the
	// choice.
	if err := tt.ChooseMenuItem("Date"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if chosen != "Date" {
		t.Errorf("the chosen item is %q", chosen)
	}

	// Outside a frame the view builds no button.
	p := tt.rt.public
	if MenuButton(p, "Sort by", func(m *Menu) {}) != (Element{}) {
		t.Error("a button was built outside a frame")
	}
}

// The public API of a menu button whose menu the window draws: its items
// show in the window, and a click on one chooses it.
func TestPopoverMenuButtonAPI(t *testing.T) {
	var chosen string
	tt := NewTester(func(c *Context) {
		PopoverMenuButton(c, "View", func(m *Menu) {
			if m.Item("Grid").Chosen() {
				chosen = "Grid"
			}
		})
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("View") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if tt.HasText("Grid") {
		t.Errorf("the menu shows before it is opened: %q", tt.Texts())
	}
	if err := tt.Click("View"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("Grid") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if err := tt.Click("Grid"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if chosen != "Grid" {
		t.Errorf("the chosen item is %q", chosen)
	}

	// Outside a frame the view builds no button.
	p := tt.rt.public
	if PopoverMenuButton(p, "View", func(m *Menu) {}) != (Element{}) {
		t.Error("a button was built outside a frame")
	}
}

// The public API of a context menu: a click with the secondary button on
// the target opens it, and choosing an item closes it.
func TestContextMenuBaseAPI(t *testing.T) {
	var (
		open   bool
		chosen string
	)
	tt := NewTester(func(c *Context) {
		row := Button(c, "Row")
		m := ContextMenuBase(c, row, &open)
		m.Popup(func(panel Element) {
			item := m.Item()
			item.Children(func() { Text(c, "Copy") })
			if item.Clicked() {
				chosen = "Copy"
			}
		})
	}, 400, 300)
	tt.Frame()
	if tt.HasText("Copy") {
		t.Errorf("the menu shows before it is opened: %q", tt.Texts())
	}

	// A click with the secondary button on the target opens the menu.
	if err := tt.RightClick("Row"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !open {
		t.Error("the menu does not open")
	}
	if !tt.HasText("Copy") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A click on an item chooses it, and the menu closes.
	if err := tt.Click("Copy"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if chosen != "Copy" {
		t.Errorf("the chosen item is %q", chosen)
	}
	if open {
		t.Error("the menu stays open")
	}

	// Outside a frame the view builds no menu.
	p := tt.rt.public
	if ContextMenuBase(p, Element{}, &open) != (MenuParts{}) {
		t.Error("a menu was built outside a frame")
	}
}
