package ui

import (
	"testing"
)

// The public API of a collapsible: its trigger opens and closes it, and
// the panel holds what it shows while it is open.
func TestCollapsibleBaseAPI(t *testing.T) {
	var (
		open     bool
		panel    Element
		progress float32
	)
	tt := NewTesterWith(func(c *Context) {
		p := CollapsibleBase(c, &open)
		p.Trigger.Children(func() { Text(c, "Advanced") })
		panel = p.Panel(func() { Text(c, "Under it") })
		progress = p.Progress()
	}, 400, 300, TesterOptions{Preferences: Preferences{ReduceMotion: true}})
	tt.Frame()
	if progress != 0 {
		t.Errorf("the progress is %v", progress)
	}
	if panel != (Element{}) {
		t.Error("the panel is there while the collapsible is closed")
	}
	if tt.HasText("Under it") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A click on the trigger opens it, and the panel shows.
	if err := tt.Click("Advanced"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !open {
		t.Error("the collapsible did not open")
	}
	if progress != 1 {
		t.Errorf("the progress is %v", progress)
	}
	if panel == (Element{}) {
		t.Error("the panel is gone while the collapsible is open")
	}
	if !tt.HasText("Under it") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// Another click closes it.
	if err := tt.Click("Advanced"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if open {
		t.Error("the collapsible did not close")
	}
	if tt.HasText("Under it") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A part that is not the view's tells nothing.
	var none CollapsibleParts
	if none.Progress() != 0 || none.Panel(func() {}) != (Element{}) {
		t.Error("a part that is not the view's answers")
	}

	// Outside a frame the view builds no collapsible.
	p := tt.rt.public
	if CollapsibleBase(p, &open) != (CollapsibleParts{}) {
		t.Error("a collapsible was built outside a frame")
	}
}

// The public API of a theme's region: its children build under the theme
// it gives.
func TestThemedAPI(t *testing.T) {
	tt := NewTester(func(c *Context) {
		th := *c.Theme()
		th.FontSize *= 4
		Themed(c, &th).Children(func() { Text(c, "Inside") })
		Text(c, "Outside")
	}, 400, 300)
	tt.Frame()
	in, ok := tt.Find("Inside")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", tt.Texts())
	}
	out, ok := tt.Find("Outside")
	if !ok {
		t.Fatalf("the text shows nowhere among %q", tt.Texts())
	}
	if in.H <= out.H {
		t.Errorf("the text inside is %v tall, no taller than the %v outside", in.H, out.H)
	}

	// A nil theme is a column.
	tt2 := NewTester(func(c *Context) {
		Themed(c, nil).Children(func() { Text(c, "Inside") })
	}, 400, 300)
	tt2.Frame()
	if !tt2.HasText("Inside") {
		t.Errorf("the texts are %q", tt2.Texts())
	}

	// Outside a frame the view builds no region.
	p := tt.rt.public
	if Themed(p, nil) != (Element{}) {
		t.Error("a region was built outside a frame")
	}
}
