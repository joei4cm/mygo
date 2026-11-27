package ui

import "testing"

// A page reads its query parameters from its route, absent ones as "",
// and the context gives the page the router it is in.
func TestRouterQuery(t *testing.T) {
	r := NewRouter("/settings")
	r.Transition = TransitionNone
	var tab, sort, who string
	var pageRouter *Router
	tt := coreNewTester(func(c *context) {
		r.coreView(c, func(rt *Route) {
			switch {
			case rt.Match("/settings"):
				tab, sort = rt.Query("tab"), rt.Query("sort")
				pageRouter = c.Router()
			case rt.Match("/about"):
				who = rt.Query("who")
			}
		})
	}, 300, 200)
	tt.Frame()
	if tab != "" || sort != "" {
		t.Errorf("a page without a query reads %q, %q", tab, sort)
	}
	if pageRouter != r {
		t.Error("the page's context gave another router")
	}

	// The parameters of the location come out unescaped, the ones of
	// other names left alone.
	r.Push("/settings?tab=go%20ui&sort=a")
	tt.Frame()
	if tab != "go ui" || sort != "a" {
		t.Errorf("the query reads %q, %q", tab, sort)
	}

	// A page without the parameter reads "" again.
	r.Push("/about")
	tt.Frame()
	if who != "" {
		t.Errorf("about reads %q, want \"\"", who)
	}
}
