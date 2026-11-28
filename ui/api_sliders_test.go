package ui

import (
	"testing"
)

// The public API of a slider: a press on it sets the value by where the
// press sits, and the keys move it.
func TestSliderBaseAPI(t *testing.T) {
	value := 0.0
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Padding(10).Children(func() {
			SliderBase(c, &value, 0, 100).Size(120, 20).PaddingX(10).Label("volume")
		})
	}, 200, 100)
	r, _ := tt.Find("volume")
	tt.Press(r.X+60, r.Y+10)
	tt.Release(r.X+60, r.Y+10)
	if value != 50 {
		t.Errorf("pressing in the middle sets %v", value)
	}
	tt.Press(r.X+110, r.Y+10)
	tt.Release(r.X+110, r.Y+10)
	if value != 100 {
		t.Errorf("pressing at the end sets %v", value)
	}
	tt.Key(0, KeyHome)
	if value != 0 {
		t.Errorf("Home sets %v", value)
	}
	tt.Key(0, KeyEnd)
	if value != 100 {
		t.Errorf("End sets %v", value)
	}

	// Outside a frame the view builds no slider.
	p := tt.rt.public
	if SliderBase(p, &value, 0, 100) != (Element{}) {
		t.Error("a slider was built outside a frame")
	}
}

// The public API of a stepped slider: a press on it snaps the value to a
// step, and the arrows step it.
func TestStepSliderAPI(t *testing.T) {
	value := 0.0
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Padding(10).Children(func() {
			StepSlider(c, &value, 0, 100, 25).Size(120, 20).PaddingX(10).Label("zoom")
		})
	}, 200, 100)
	r, _ := tt.Find("zoom")
	tt.Press(r.X+60, r.Y+10)
	tt.Release(r.X+60, r.Y+10)
	if value != 50 {
		t.Errorf("pressing in the middle sets %v", value)
	}
	tt.Press(r.X+20, r.Y+10)
	tt.Release(r.X+20, r.Y+10)
	if value != 0 {
		t.Errorf("pressing near the start snaps to %v", value)
	}
	tt.Key(0, KeyRight)
	if value != 25 {
		t.Errorf("Right sets %v", value)
	}

	// Outside a frame the view builds no slider.
	p := tt.rt.public
	if StepSlider(p, &value, 0, 100, 25) != (Element{}) {
		t.Error("a slider was built outside a frame")
	}
}

// The public API of a range slider: a press on it moves the knob nearest
// to it, which snaps to a step.
func TestRangeSliderAPI(t *testing.T) {
	low, high := 25.0, 75.0
	tt := NewTester(func(c *Context) {
		Column(c).AlignItems(Start).Padding(10).Children(func() {
			RangeSlider(c, &low, &high, 0, 100, 25).Size(120, 20).PaddingX(10).Label("price")
		})
	}, 200, 100)
	r, _ := tt.Find("price")
	tt.Press(r.X+20, r.Y+10)
	tt.Release(r.X+20, r.Y+10)
	if low != 0 {
		t.Errorf("pressing at the start sets the low knob to %v", low)
	}
	tt.Press(r.X+100, r.Y+10)
	tt.Release(r.X+100, r.Y+10)
	if high != 100 {
		t.Errorf("pressing at the end sets the high knob to %v", high)
	}

	// Outside a frame the view builds no slider.
	p := tt.rt.public
	if RangeSlider(p, &low, &high, 0, 100, 25) != (Element{}) {
		t.Error("a slider was built outside a frame")
	}
}

// The public API of a spinner: assistive technology sees a progress
// indicator of no value.
func TestSpinnerAPI(t *testing.T) {
	var spinnerRole Role
	tt := NewTester(func(c *Context) {
		s := Spinner(c)
		spinnerRole = c.rt.states[s.node().ID()].role
	}, 200, 100)
	tt.Frame()
	if spinnerRole != RoleProgress {
		t.Errorf("the role is %v", spinnerRole)
	}

	// Outside a frame the view builds no spinner.
	p := tt.rt.public
	if Spinner(p) != (Element{}) {
		t.Error("a spinner was built outside a frame")
	}
}
