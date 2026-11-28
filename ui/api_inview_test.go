package ui

import "testing"

// The public API of what shows in the window: a card in a scroll
// container tells whether it showed at all, how much of it showed, and
// how near it is to showing.
func TestInViewAPI(t *testing.T) {
	var inView, near, half, quarter, full, deep, noAmount bool
	tt := NewTester(func(c *Context) {
		Scroll(c).Size(120, 100).Children(func() {
			Box(c).Size(120, 40)
			Box(c).Size(120, 70)
			card := Box(c).Key("card").Size(120, 80)
			Box(c).Size(120, 60)
			inView = card.InView()
			near = card.InViewMargin(200)
			half = card.InViewAmount(0.5)
			quarter = card.InViewAmount(0.25)
			full = card.InViewAmount(1)
			deep = card.InViewMargin(-40)
			noAmount = card.InViewAmount(0)
		})
	}, 240, 200)
	defer tt.rt.close()
	tt.Frame()
	if inView || !near || half || quarter || full || deep {
		t.Fatalf("below the fold: in view %v, within 200 %v, half %v, a quarter %v, all %v, within -40 %v",
			inView, near, half, quarter, full, deep)
	}
	if noAmount != inView {
		t.Errorf("an amount of 0 answered %v, the view %v", noAmount, inView)
	}
	// Scrolled by 40, 30 of the 80 DIPs of the card show: a third of it.
	tt.Scroll(60, 50, 0, 40)
	tt.Frame()
	if !inView || half || !quarter || full || deep {
		t.Fatalf("a third in view: in view %v, half %v, a quarter %v, all %v, within -40 %v",
			inView, half, quarter, full, deep)
	}
	// Scrolled by 110, the card shows whole.
	tt.Scroll(60, 50, 0, 70)
	tt.Frame()
	if !inView || !full || !deep {
		t.Fatalf("all in view: in view %v, all %v, within -40 %v", inView, full, deep)
	}
	if noAmount != inView {
		t.Errorf("an amount of 0 answered %v, the view %v", noAmount, inView)
	}
}

// Outside a frame these views know nothing.
func TestInViewAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.InView() {
		t.Error("an absent element is in view")
	}
	if absent.InViewMargin(10) {
		t.Error("an absent element is within a margin")
	}
	if absent.InViewAmount(0.5) {
		t.Error("an absent element shows a part of itself")
	}
}
