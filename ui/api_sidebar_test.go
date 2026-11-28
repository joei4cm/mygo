package ui

import "testing"

// The public API of a section of a sidebar: it holds its items while it
// is open, and a click on its title opens and closes it.
func TestSidebarSectionAPI(t *testing.T) {
	open := false
	tt := NewTester(func(c *Context) {
		SidebarSection(c, "Favorites", &open, func() {
			Text(c, "Air")
		})
	}, 300, 300)

	// A closed section shows its title, not its items.
	if !tt.HasText("Favorites") {
		t.Fatal("the section shows no title")
	}
	if tt.HasText("Air") {
		t.Fatal("the closed section shows its items")
	}

	// A click on the title opens it, and shows what it holds.
	if err := tt.Click("Favorites"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.HasText("Air") {
		t.Fatalf("the click did not open the section (%v)", open)
	}

	// A click again closes it, and hides what it held.
	if err := tt.Click("Favorites"); err != nil {
		t.Fatal(err)
	}
	if open || tt.HasText("Air") {
		t.Fatalf("the click did not close the section (%v)", open)
	}
}

// The public API of a section of a sidebar that keeps no state of its own:
// it stays open, and a click on its title does nothing.
func TestSidebarSectionWithoutStateAPI(t *testing.T) {
	tt := NewTester(func(c *Context) {
		SidebarSection(c, "Places", nil, func() {
			Text(c, "Macintosh HD")
		})
	}, 300, 300)

	if !tt.HasText("Macintosh HD") {
		t.Fatal("the section hides its items")
	}
	if err := tt.Click("Places"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Macintosh HD") {
		t.Fatal("the click closed a section that keeps no state")
	}

	// Outside a frame the view builds no section.
	p := tt.rt.public
	if SidebarSection(p, "Places", nil, func() {}).Valid() {
		t.Error("a section was built outside a frame")
	}
}
