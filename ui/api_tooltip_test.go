package ui

import (
	"testing"
	"time"
)

// The public API of a tooltip: it shows once the pointer has rested on
// its anchor, and goes as the pointer leaves.
func TestTooltipBaseAPI(t *testing.T) {
	var tip Element
	tt, now := clockTester(func(c *context) {
		pc := makeContext(c)
		Column(pc).Children(func() {
			save := Button(pc, "Save")
			tip = TooltipBase(pc, save, func(tip Element) {
				tip.Children(func() { Text(pc, "Saves the file") })
			})
		})
	}, 400, 300)
	tick := ticker(tt, now)
	s, _ := tt.Find("Save")
	tt.Move(s.X+s.W/2, s.Y+s.H/2)
	tick(500 * time.Millisecond)
	if tt.HasText("Saves the file") || tip != (Element{}) {
		t.Fatal("the tip shows before the pointer rested 0.6 seconds")
	}
	tick(150 * time.Millisecond)
	if !tt.HasText("Saves the file") || tip == (Element{}) {
		t.Fatal("no tip once the pointer rested")
	}

	// A click on its anchor hides the tip.
	tt.ClickAt(s.X+s.W/2, s.Y+s.H/2)
	tick(16 * time.Millisecond)
	if tt.HasText("Saves the file") || tip != (Element{}) {
		t.Error("the tip shows after a click on its anchor")
	}

	// Back and rested, it shows again.
	tt.Move(390, 290)
	tt.Frame()
	tt.Move(s.X+s.W/2, s.Y+s.H/2)
	tick(time.Second)
	if !tt.HasText("Saves the file") || tip == (Element{}) {
		t.Fatal("no tip once the pointer rested again")
	}

	// The tip goes as the pointer leaves its anchor.
	tt.Move(390, 290)
	tt.Frame()
	if tt.HasText("Saves the file") || tip != (Element{}) {
		t.Error("the tip shows after the pointer left")
	}

	// Outside a frame the view builds no tip.
	p := tt.rt.public
	if TooltipBase(p, Element{}, func(tip Element) {}) != (Element{}) {
		t.Error("a tip was built outside a frame")
	}
}
