package ui

import "testing"

// The public API of sizing an element by a percentage of its parent: its
// width, height, bounds and flex basis.
func TestPercentSizesAPI(t *testing.T) {
	var sized, floor, ceil, basis Rect
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Children(func() {
			Box(c).Size(200, 100).AlignItems(Start).Children(func() {
				sized = Box(c).WidthPercent(50).HeightPercent(40).Bounds()
			})
			Box(c).Size(200, 100).AlignItems(Start).Children(func() {
				floor = Box(c).Size(20, 10).MinWidthPercent(50).MinHeightPercent(40).Bounds()
			})
			Box(c).Size(200, 100).AlignItems(Start).Children(func() {
				ceil = Box(c).Size(180, 90).MaxWidthPercent(50).MaxHeightPercent(40).Bounds()
			})
			Row(c).Size(200, 100).AlignItems(Start).Children(func() {
				basis = Box(c).BasisPercent(25).Height(20).Bounds()
			})
		})
	}, 300, 500)
	tt.Frame()

	if sized.W != 100 || sized.H != 40 {
		t.Errorf("WidthPercent/HeightPercent gave %v, want 100x40", sized)
	}
	if floor.W != 100 || floor.H != 40 {
		t.Errorf("MinWidthPercent/MinHeightPercent gave %v, want 100x40", floor)
	}
	if ceil.W != 100 || ceil.H != 40 {
		t.Errorf("MaxWidthPercent/MaxHeightPercent gave %v, want 100x40", ceil)
	}
	if basis.W != 50 {
		t.Errorf("BasisPercent gave %v, want width 50", basis)
	}

	// Outside a frame the view sizes nothing.
	if (Element{}).WidthPercent(50).Valid() {
		t.Error("WidthPercent built outside a frame")
	}
	if (Element{}).HeightPercent(40).Valid() {
		t.Error("HeightPercent built outside a frame")
	}
	if (Element{}).MinWidthPercent(50).Valid() {
		t.Error("MinWidthPercent built outside a frame")
	}
	if (Element{}).MinHeightPercent(40).Valid() {
		t.Error("MinHeightPercent built outside a frame")
	}
	if (Element{}).MaxWidthPercent(50).Valid() {
		t.Error("MaxWidthPercent built outside a frame")
	}
	if (Element{}).MaxHeightPercent(40).Valid() {
		t.Error("MaxHeightPercent built outside a frame")
	}
	if (Element{}).BasisPercent(25).Valid() {
		t.Error("BasisPercent built outside a frame")
	}
}
