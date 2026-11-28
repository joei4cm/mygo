package ui

import "testing"

// Effect draws the element's rendered content through fx; nil removes
// the effect and draws the element as it is.
func TestEffectAPI(t *testing.T) {
	fx := &layerSpy{ids: map[uint64]bool{}}
	var effect Effect = fx
	var applied Effect
	tt := NewTester(func(c *Context) {
		box := Box(c).Width(80).Height(40).Background(RGB(220, 20, 30)).Effect(effect)
		applied = box.node().effect
	}, 160, 100)
	defer tt.rt.close()

	if applied != fx {
		t.Fatalf("element has effect %T, want the supplied effect", applied)
	}
	if len(fx.ids) == 0 {
		t.Fatal("the supplied effect did not paint")
	}

	effect = nil
	tt.Frame()
	applied = tt.rt.c.root.first.effect
	if applied != nil {
		t.Fatalf("nil did not remove the effect: %T", applied)
	}
	if c := tt.Image().RGBAAt(20, 20); c.R != 220 || c.G != 20 {
		t.Errorf("without the effect, the element does not draw as it is: %v", c)
	}
}

func TestEffectAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Effect(&layerSpy{ids: map[uint64]bool{}}).Valid() {
		t.Error("an absent element took an effect")
	}
}
