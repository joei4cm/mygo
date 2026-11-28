package ui

import (
	"testing"
)

// The public API of an accordion: a click on a section's header opens and
// closes it, and the arrows move between the headers.
func TestAccordionAPI(t *testing.T) {
	var (
		general, privacy bool
		changes          int
	)
	tt := NewTester(func(c *Context) {
		Accordion(c, func() {
			if AccordionItem(c, "General", &general, func() { Text(c, "General settings") }).Changed() {
				changes++
			}
			if AccordionItem(c, "Privacy", &privacy, func() { Text(c, "Privacy settings") }).Changed() {
				changes++
			}
		})
	}, 400, 300)
	tt.Frame()
	if general || privacy {
		t.Errorf("the sections are %v and %v", general, privacy)
	}
	if tt.HasText("General settings") || tt.HasText("Privacy settings") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A click on a header opens its section.
	if err := tt.Click("General"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !general || privacy {
		t.Errorf("the sections are %v and %v", general, privacy)
	}
	if !tt.HasText("General settings") || tt.HasText("Privacy settings") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	if changes == 0 {
		t.Error("no change was reported")
	}

	// A click on another header opens it beside the first.
	if err := tt.Click("Privacy"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !general || !privacy {
		t.Errorf("the sections are %v and %v", general, privacy)
	}
	if !tt.HasText("General settings") || !tt.HasText("Privacy settings") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// The arrows move between the headers, and Space closes the one with
	// the focus.
	tt.Key(0, KeyUp)
	tt.Key(0, KeySpace)
	tt.Frame()
	if general || !privacy {
		t.Errorf("the sections are %v and %v", general, privacy)
	}
	if tt.HasText("General settings") || !tt.HasText("Privacy settings") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Outside a frame the view builds no accordion.
	p := tt.rt.public
	if Accordion(p, func() {}) != (Element{}) ||
		AccordionItem(p, "General", &general, func() {}) != (Element{}) {
		t.Error("an accordion was built outside a frame")
	}
}
