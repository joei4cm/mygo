package ui

import "testing"

// A handle resolves to its control of this pass, and its focus queries
// answer from the committed frame: the focused control, and those around it.
func TestRefResolveFocused(t *testing.T) {
	var parent, child, orphan Handle
	var parentResolved, childResolved, orphanResolved bool
	var parentRect, childRect Rect
	var parentFocused, parentWithin, childFocused, childWithin bool
	var outside Handle
	var ctxFocused, ctxWithin, outsideWithin bool
	tt := NewTester(func(c *Context) {
		Box(c).Size(80, 60).Bind(&parent).Label("boundParent").Children(func() {
			Box(c).Size(40, 20).Bind(&child).Label("boundChild")
		})
		Box(c).Size(40, 20).Bind(&outside)
		outsideWithin = outside.FocusWithin(c)
		resolvedParent, resolvedChild := c.Resolve(parent), c.Resolve(child)
		parentResolved, childResolved = resolvedParent.Valid(), resolvedChild.Valid()
		orphanResolved = c.Resolve(orphan).Valid()
		parentRect, childRect = parent.Bounds(c), child.Bounds(c)
		parentFocused, parentWithin = parent.Focused(c), parent.FocusWithin(c)
		childFocused, childWithin = child.Focused(c), child.FocusWithin(c)
		ctxFocused, ctxWithin = c.Focused(child), c.FocusWithin(parent)
	}, 100, 100)
	tt.Frame()
	if orphanResolved {
		t.Error("Resolve returned an element for an unbound handle")
	}
	if !parentResolved || !childResolved {
		t.Fatalf("Resolve returned nothing for bound controls: parent %v, child %v", parentResolved, childResolved)
	}
	if p, ok := tt.Find("boundParent"); !ok || p != parentRect {
		t.Errorf("the parent resolves to %v, the box shows at %v", parentRect, p)
	}
	if r, ok := tt.Find("boundChild"); !ok || r != childRect {
		t.Errorf("the child resolves to %v, the box shows at %v", childRect, r)
	}
	if parentFocused || childFocused || parentWithin || childWithin {
		t.Errorf("nothing is focused yet: parent %v/%v, child %v/%v", parentFocused, parentWithin, childFocused, childWithin)
	}
	// Focus the child: it answers Focused and FocusWithin both (the
	// focused control is within itself), its parent only FocusWithin.
	child.Focus()
	tt.Frame()
	tt.Frame()
	if !childFocused {
		t.Error("the focused child does not answer Focused")
	}
	if !childWithin {
		t.Error("the focused child does not answer FocusWithin")
	}
	if !parentWithin {
		t.Error("the parent of the focused child does not answer FocusWithin")
	}
	if parentFocused {
		t.Error("the parent answers Focused")
	}
	if !ctxFocused || !ctxWithin {
		t.Errorf("the context forms answer Focused %v, FocusWithin %v", ctxFocused, ctxWithin)
	}
	if outsideWithin {
		t.Error("FocusWithin answers for a control outside the focus")
	}
}

// CancelFocus withdraws a request, one that waits for a binding and one
// that a binding carries, in one window or everywhere; focus works again
// after one.
func TestRefCancelFocus(t *testing.T) {
	var handle Handle
	var focused bool
	var win *Context
	var bound bool
	tt := NewTester(func(c *Context) {
		win = c
		if bound {
			Box(c).Size(40, 20).Bind(&handle)
		}
		focused = handle.Focused(c)
	}, 100, 100)

	var nilHandle *Handle
	nilHandle.CancelFocus()

	// Requested before its first binding, cancelled before it binds.
	handle.Focus()
	handle.CancelFocus()
	bound = true
	tt.Frame()
	tt.Frame()
	if focused {
		t.Error("a pending request, cancelled, still focused the control")
	}

	// Requested on a binding, cancelled in its window before it lands.
	var unbound Handle
	unbound.Focus()
	unbound.CancelFocus(win)
	handle.Focus()
	handle.CancelFocus(win)
	tt.Frame()
	tt.Frame()
	if focused {
		t.Error("a request, cancelled in its window, still focused the control")
	}

	// And a request afterwards takes; RequestFocus is its alias.
	handle.RequestFocus()
	tt.Frame()
	tt.Frame()
	if !focused {
		t.Error("focus does not answer after a cancel")
	}

	// Cancelling in more than one window panics.
	if func() (panicked bool) {
		defer func() { panicked = recover() != nil }()
		handle.CancelFocus(win, win)
		return
	}() == false {
		t.Error("CancelFocus in two windows does not panic")
	}
}

// Binding a handle to two controls of one build panics.
func TestRefDoubleBindPanics(t *testing.T) {
	var handle Handle
	var panicked bool
	tt := NewTester(func(c *Context) {
		defer func() { panicked = recover() != nil }()
		Box(c).Bind(nil) // no handle, nothing binds
		Box(c).Bind(&handle)
		Box(c).Bind(&handle)
	}, 100, 100)
	tt.Frame()
	if !panicked {
		t.Error("binding a handle to two controls does not panic")
	}
}
