package ui

import "testing"

// SetChanged reports the value of an element as changed: in the pass that
// calls it, and, called between frames, in the next build; the passes
// after read it unchanged again.
func TestSetChanged(t *testing.T) {
	value := ""
	var n *node
	pass := 0
	var again, later bool
	tt := coreNewTester(func(c *context) {
		pass++
		n = coreTextInputBase(c, &value)
		if pass == 1 {
			n.SetChanged()
			if !n.Changed() {
				t.Error("SetChanged did not report the value changed in its pass")
			}
		} else if pass == 2 {
			again = n.Changed()
		} else if pass >= 3 && pass <= 6 {
			later = later || n.Changed()
		}
	}, 100, 40)
	tt.Frame()
	if again {
		t.Error("the pass after still reads the value as changed")
	}
	// Called between frames, the change reports in the next build.
	n.SetChanged()
	tt.Frame()
	if !later {
		t.Error("a change made between frames missed the next build")
	}
}
