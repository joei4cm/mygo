package ui

import "testing"

// KeepOpen keeps a menu open while its item is chosen; a plain item closes it.
func TestKeepOpenAPI(t *testing.T) {
	var open bool
	var chosen []string
	theme := "Light"
	tt := NewTester(func(c *Context) {
		m := MenuBase(c, &open)
		m.Trigger.Children(func() { Text(c, "Account") })
		m.Popup(func(panel Element) {
			panel.Padding(4)
			item := m.Item()
			item.Children(func() { Text(c, "Profile") })
			if item.Clicked() {
				chosen = append(chosen, "Profile")
			}
			m.Group(func(row Element) {
				row.Gap(4)
				for _, th := range []string{"Light", "Dark", "System"} {
					item := m.RadioItem(theme == th).KeepOpen()
					item.Children(func() { Text(c, th) })
					if item.Clicked() {
						chosen = append(chosen, th)
						theme = th
					}
				}
			})
		})
	}, 300, 300)
	defer tt.rt.close()
	open = true
	tt.Frame()
	if !tt.HasText("Profile") {
		t.Fatal("the menu did not open")
	}
	// A radio item that keeps the menu open: the choice is made and the menu stays.
	if err := tt.Click("Dark"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if theme != "Dark" || !open || !tt.HasText("Profile") {
		t.Fatalf("choosing Dark: theme %q, open %v", theme, open)
	}
	// A plain item closes the menu.
	if err := tt.Click("Profile"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if open || tt.HasText("Profile") {
		t.Fatalf("choosing Profile: open %v", open)
	}
	if len(chosen) != 2 || chosen[0] != "Dark" || chosen[1] != "Profile" {
		t.Fatalf("chosen %q", chosen)
	}
}

func TestKeepOpenAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.KeepOpen().Valid() {
		t.Error("an absent element kept a menu open")
	}
}
