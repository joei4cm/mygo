package ui

import "testing"

// Skew slants an element after the transform it has, around the origin it
// has.
func TestElementSkew(t *testing.T) {
	var skewed, composed *node
	tt := coreNewTester(func(c *context) {
		skewed = coreBox(c).Absolute().Left(10).Top(10).Size(20, 20).Background(RGB(220, 20, 30)).Skew(45, 0)
		composed = coreBox(c).Absolute().Left(40).Top(10).Size(20, 20).Background(RGB(220, 20, 30)).Rotate(90).Skew(30, 0)
	}, 100, 100)
	tt.Frame()
	if skewed.transform.flat() != Skewing(45, 0).flat() {
		t.Errorf("Skew(45, 0) alone made %v", skewed.transform.flat())
	}
	// Skew comes after the transform the element has.
	if composed.transform.flat() != Rotation(90).Then(Skewing(30, 0)).flat() {
		t.Errorf("Rotate(90).Skew(30, 0) made %v", composed.transform.flat())
	}

	isRed := func(tt *Tester, x, y int) bool {
		c := tt.Image().RGBAAt(x, y)
		return c.R > 200 && c.G < 60 && c.B < 60
	}
	// Skewed 45 degrees around its middle, the box slants: rows above
	// the middle slide left, rows under it slide right, a pixel's worth
	// a row.
	tt = coreNewTester(func(c *context) {
		coreBox(c).Absolute().Left(10).Top(10).Size(20, 20).Background(RGB(220, 20, 30)).Skew(45, 0)
	}, 100, 100)
	tt.Frame()
	if !isRed(tt, 5, 11) {
		t.Error("the box skewed about its middle did not slide its top rows left")
	}
	if isRed(tt, 10, 29) {
		t.Error("the box skewed about its middle did not slide its bottom rows right")
	}
	// A box without a skew covers the one and not the other.
	tt = coreNewTester(func(c *context) {
		coreBox(c).Absolute().Left(10).Top(10).Size(20, 20).Background(RGB(220, 20, 30))
	}, 100, 100)
	tt.Frame()
	if !isRed(tt, 10, 29) || isRed(tt, 5, 11) {
		t.Error("the box without a skew is not where it should be")
	}
	// Skewed around its top left corner instead, the slant keeps that
	// corner's row where it is.
	tt = coreNewTester(func(c *context) {
		coreBox(c).Absolute().Left(10).Top(10).Size(20, 20).Background(RGB(220, 20, 30)).TransformOrigin(0, 0).Skew(45, 0)
	}, 100, 100)
	tt.Frame()
	if isRed(tt, 5, 11) || !isRed(tt, 12, 11) {
		t.Error("the box skewed around its corner did not slide from it alone")
	}
	if !isRed(tt, 29, 10) {
		t.Error("the box skewed around its corner moved the corner's row")
	}
}
