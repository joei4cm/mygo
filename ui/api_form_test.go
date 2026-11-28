package ui

import (
	"testing"
)

// The public API of a form: the labels of its fields line up as wide as
// the widest, and a click on a label focuses its field.
func TestFormAPI(t *testing.T) {
	name, mail := "", ""
	tt := NewTester(func(c *Context) {
		Form(c, func() {
			Field(c, "Name", func() { TextInput(c, &name).Label("The name") })
			Field(c, "Electronic mail", func() { TextInput(c, &mail).Label("The mail") })
		})
	}, 400, 300)
	tt.Frame()

	// The labels line up: the fields start at the same place.
	n, ok := tt.Access("The name")
	if !ok {
		t.Fatal("the field names nothing")
	}
	e, ok := tt.Access("The mail")
	if !ok {
		t.Fatal("the field names nothing")
	}
	if n.Bounds.X != e.Bounds.X {
		t.Errorf("the fields start at %v and %v", n.Bounds.X, e.Bounds.X)
	}

	// A click on a label focuses its field, which types.
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Ada")
	if name != "Ada" {
		t.Errorf("the name is %q", name)
	}

	// Outside a frame the view builds no form.
	p := tt.rt.public
	if Form(p, func() {}) != (Element{}) {
		t.Error("a form was built outside a frame")
	}
}

// The public API of a fieldset: its legend names the group for assistive
// technology, and shows above what it holds.
func TestFieldsetAPI(t *testing.T) {
	news := false
	tt := NewTester(func(c *Context) {
		Fieldset(c, "Notifications", func() {
			Checkbox(c, &news, "News")
		})
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("Notifications") || !tt.HasText("News") {
		t.Errorf("the texts are %q", tt.Texts())
	}
	info, ok := tt.Access("Notifications")
	if !ok {
		t.Fatal("the legend names nothing")
	}
	if info.Role != RoleGroup {
		t.Errorf("the role is %v", info.Role)
	}

	// Outside a frame the view builds no fieldset.
	p := tt.rt.public
	if Fieldset(p, "Notifications", func() {}) != (Element{}) {
		t.Error("a fieldset was built outside a frame")
	}
}

// The public API of breadcrumbs: every crumb but the last is a link, and
// a click on one reports where the user went.
func TestBreadcrumbsAPI(t *testing.T) {
	chosen := -1
	tt := NewTester(func(c *Context) {
		Breadcrumbs(c, []string{"Home", "Files", "notes.txt"}, &chosen)
	}, 400, 300)
	tt.Frame()
	if !tt.HasText("Home") || !tt.HasText("Files") || !tt.HasText("notes.txt") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// A click on a crumb reports it.
	if err := tt.Click("Files"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if chosen != 1 {
		t.Errorf("the chosen crumb is %d", chosen)
	}
	if err := tt.Click("Home"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if chosen != 0 {
		t.Errorf("the chosen crumb is %d", chosen)
	}

	// Outside a frame the view builds no breadcrumbs.
	p := tt.rt.public
	if Breadcrumbs(p, nil, &chosen) != (Element{}) {
		t.Error("breadcrumbs were built outside a frame")
	}
}
