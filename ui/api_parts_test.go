package ui

import (
	"testing"
)

// The combobox's popup opens and closes, and filters the options with what
// the user typed.
func TestComboboxPartsAPI(t *testing.T) {
	var (
		text      string
		open      bool
		filtering bool
		closeIt   bool
	)
	tt := NewTester(func(c *Context) {
		p := ComboboxBase(c, &text)
		p.Input.Label("Combo")
		if closeIt {
			p.SetOpen(false)
		}
		open = p.Open()
		filtering = p.Filtering()
		p.Popup(func(panel Element) {
			panel.Children(func() {
				for _, value := range []string{"Alpha", "Beta"} {
					p.Item(value).Children(func() { Text(c, value) })
				}
			})
		})
	}, 300, 300)

	// A click on the input opens the popup, which shows the options.
	tt.Click("Combo")
	tt.Frame()
	if !open {
		t.Error("the popup does not show")
	}
	if filtering {
		t.Error("the popup filters what was not typed")
	}
	if !tt.HasText("Alpha") || !tt.HasText("Beta") {
		t.Errorf("the options are %q", tt.Texts())
	}

	// A part that is not the view's tells nothing.
	var none ComboboxParts
	none.SetOpen(false)
	if none.Open() || none.Filtering() {
		t.Error("a part that is not the view's answers")
	}

	// What the user types filters them.
	tt.Type("al")
	tt.Frame()
	if !filtering {
		t.Error("the popup does not filter what was typed")
	}

	// The view closes the popup itself.
	closeIt = true
	tt.Frame()
	if open {
		t.Error("the popup shows though the view closed it")
	}
	if tt.HasText("Alpha") {
		t.Errorf("the options are %q", tt.Texts())
	}
}

// The accordion's item opens and closes, and shows what is below it while
// it is open.
func TestAccordionPartsAPI(t *testing.T) {
	open := false
	var (
		progress float32
		panel    Element
	)
	tt := NewTesterWith(func(c *Context) {
		a := AccordionBase(c)
		item := a.Item(&open)
		item.Trigger.Children(func() { Text(c, "Section") })
		panel = item.Panel(func() { Text(c, "Under it") })
		progress = item.Progress()
	}, 300, 200, TesterOptions{Preferences: Preferences{ReduceMotion: true}})
	tt.Frame()
	if progress != 0 {
		t.Errorf("the progress is %v", progress)
	}
	if panel != (Element{}) {
		t.Error("the panel is there while the item is closed")
	}
	if tt.HasText("Under it") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A click on the header opens the item, and what is below it shows.
	tt.Click("Section")
	tt.Frame()
	if !open {
		t.Error("the item did not open")
	}
	if progress != 1 {
		t.Errorf("the progress is %v", progress)
	}
	if panel == (Element{}) {
		t.Error("the panel is gone while the item is open")
	}
	if !tt.HasText("Under it") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A part that is not the view's tells nothing: the accordion's item
	// is one, too.
	var none AccordionParts
	noneItem := none.Item(&open)
	var noneParts CollapsibleParts
	if noneParts.Progress() != 0 || noneParts.Panel(func() {}) != (Element{}) {
		t.Error("a part that is not the view's answers")
	}
	if noneItem.Progress() != 0 || noneItem.Panel(func() {}) != (Element{}) {
		t.Error("a part that is not the view's answers")
	}
}
