package ui

import "testing"

// The public API of the content of a window: it holds the app's view, and
// gives it the window's own context.
func TestViewAPI(t *testing.T) {
	tt := NewTester(func(c *Context) {}, 100, 100)

	var got *Context
	runs := 0
	content := View(func(c *Context) { runs++; got = c })
	if content == nil {
		t.Fatal("the content is nothing")
	}
	content.view(&tt.rt.c)
	if got != tt.rt.public {
		t.Fatal("the content did not give the window's context to the view")
	}
	if runs != 1 {
		t.Fatalf("the view ran %d times", runs)
	}
}
