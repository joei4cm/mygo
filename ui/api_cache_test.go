package ui

import "testing"

// The public API of a layer kept by Cache: it paints like the uncached
// element, stays in a bitmap across frames that only rebuild its sibling,
// and redraws that bitmap when its own content changes.
func TestCacheAPI(t *testing.T) {
	accent := RGB(255, 255, 255)
	var cached, plain *state
	makeView := func(keep bool, out **state) func(*Context) {
		return func(c *Context) {
			Row(c).Size(220, 60).Background(RGB(255, 255, 255)).Children(func() {
				card := Box(c).Size(120, 50).Background(RGB(30, 60, 200)).Padding(8).Children(func() {
					Text(c, "Daily mix").TextColor(accent)
					Box(c).Absolute().Right(0).Top(0).Size(30, 50).Gradient(RGB(30, 60, 200), RGB(200, 40, 120), 30)
				})
				if keep {
					card = card.Cache()
				}
				*out = card.node().st
				Box(c).Size(20, 20).Draw(func(p *Painter, r Rect) {
					p.AnimationFrame()
					p.Fill(Rect{r.X, r.Y, 10, r.H}, RGB(0, 0, 0), 0)
				})
			})
		}
	}
	ct := NewTester(makeView(true, &cached), 240, 80)
	defer ct.rt.close()
	pt := NewTester(makeView(false, &plain), 240, 80)
	defer pt.rt.close()
	ct.Frame()
	pt.Frame()
	if cached == nil || cached.effectLayer == nil || cached.effectLayer.img == nil {
		t.Fatal("Cache did not keep the layer")
	}
	layer := cached.effectLayer.img
	if layer.Transient {
		t.Fatal("Cache's bitmap is transient")
	}
	if diff := sameImage(ct.Image(), pt.Image()); diff != "" {
		t.Fatalf("cached layer differs from painting it: %s", diff)
	}
	// A frame that rebuilds and paints both views leaves the still layer
	// alone; changing its text paints the layer again, like the plain view.
	version := layer.Version()
	ct.Frame()
	pt.Frame()
	if layer.Version() != version {
		t.Errorf("the still layer was drawn again: %d -> %d", version, layer.Version())
	}
	accent = RGB(255, 220, 0)
	ct.Frame()
	pt.Frame()
	if layer.Version() == version {
		t.Error("changed content did not redraw the kept layer")
	}
	if diff := sameImage(ct.Image(), pt.Image()); diff != "" {
		t.Errorf("redrawn cached layer differs from painting it: %s", diff)
	}
}

// Outside a frame Cache cannot keep an absent element.
func TestCacheAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.Cache().Valid() {
		t.Error("an absent element kept a layer")
	}
}
