package ui

import (
	"testing"
)

// The public API of a router: the view builds the page it shows, tells
// the element holding it, and whether the system shows the pages as its
// own stack.
func TestRouterAPI(t *testing.T) {
	r := NewRouter("/")
	var (
		page    Element
		route   *Route
		inStack bool
		built   bool
	)
	tt := NewTester(func(c *Context) {
		r.View(c, func(rt *Route) {
			switch {
			case rt.Match("/"):
				rt.Title("Home")
				route = rt
				page = rt.Page()
				inStack = r.InSystemStack(c)
				built = true
			}
			Text(c, "Body")
		})
	}, 400, 300)
	tt.h.stack = &stackHost{can: true}
	tt.Frame()
	if !built {
		t.Fatal("the page did not build")
	}
	if page == (Element{}) {
		t.Error("the page has no element")
	}
	if !tt.HasText("Body") {
		t.Errorf("the texts are %q", tt.Texts())
	}

	// The system shows the pages as its own stack only when the router
	// asks for it.
	r.SystemStack = true
	tt.Frame()
	if !inStack {
		t.Error("the pages are not in the system's stack")
	}
	r.SystemStack = false
	tt.Frame()
	if inStack {
		t.Error("the pages are in the system's stack")
	}

	// Outside a frame the router tells nothing.
	p := tt.rt.public
	if r.InSystemStack(p) {
		t.Error("the router answers outside a frame")
	}
	if route.Page() != (Element{}) {
		t.Error("the route answers outside a frame")
	}
	if route.View(p, func(*Route) {}) != (Element{}) {
		t.Error("the route answers outside a frame")
	}
	if r.View(p, func(*Route) {}) != (Element{}) {
		t.Error("the router answers outside a frame")
	}
}

// The public API of a route's pages: the view builds the pages the rest
// of the path holds.
func TestRouteViewAPI(t *testing.T) {
	r := NewRouter("/files/{name...}")
	var (
		name  string
		built bool
	)
	tt := NewTester(func(c *Context) {
		r.View(c, func(rt *Route) {
			if rt.Match("/files/{name...}") {
				rt.View(c, func(rt *Route) {
					name = rt.Param("name")
					built = true
					Text(c, "Name: "+name)
				})
			}
		})
	}, 400, 300)
	tt.Frame()
	r.Push("/files/notes")
	tt.Frame()
	if !built {
		t.Fatal("the pages did not build")
	}
	if name != "notes" {
		t.Errorf("the name is %q", name)
	}
	if !tt.HasText("Name: notes") {
		t.Errorf("the texts are %q", tt.Texts())
	}
}
