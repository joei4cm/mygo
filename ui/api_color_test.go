package ui

import (
	"testing"
)

// The public API of a color picker: it shows the color in hex, a click on
// a swatch chooses the color, and the field edits it.
func TestColorPickerAPI(t *testing.T) {
	color := Hex("#3b82f6")
	var changed bool
	tt := NewTester(func(c *Context) {
		if ColorPicker(c, &color).Changed() {
			changed = true
		}
	}, 400, 600)
	tt.Frame()
	info, ok := tt.Access("Hex")
	if !ok || info.Value != "#3b82f6" {
		t.Errorf("the hex field: %+v, %v", info, ok)
	}

	// A click on a swatch chooses its color.
	if err := tt.Click("Red"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if color != Hex("#ef4444") {
		t.Errorf("the color is %v", color)
	}
	if info, ok = tt.Access("Hex"); !ok || info.Value != "#ef4444" {
		t.Errorf("the hex field: %+v, %v", info, ok)
	}
	if !changed {
		t.Error("no change was reported")
	}

	// The field edits the color in hex.
	if err := tt.Click("Hex"); err != nil {
		t.Fatal(err)
	}
	tt.Key(Ctrl, KeyA)
	tt.Type("#00ff00")
	tt.Frame()
	if color != Hex("#00ff00") {
		t.Errorf("the color is %v", color)
	}

	// Outside a frame the view builds no picker.
	p := tt.rt.public
	if ColorPicker(p, &color) != (Element{}) {
		t.Error("a picker was built outside a frame")
	}
}

// The public API of a color well: a click opens a picker below it, which
// choosing a color closes.
func TestColorWellAPI(t *testing.T) {
	color := RGB(255, 0, 0)
	tt := NewTester(func(c *Context) {
		ColorWell(c, &color).Label("Tint")
	}, 400, 600)
	tt.Frame()
	if _, ok := tt.Access("Hex"); ok {
		t.Error("the picker shows before it is opened")
	}
	info, ok := tt.Access("Tint")
	if !ok || info.Value != "#ff0000" {
		t.Errorf("the well: %+v, %v", info, ok)
	}

	// A click on the well opens its picker.
	if err := tt.Click("Tint"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if info, ok := tt.Access("Hex"); !ok || info.Value != "#ff0000" {
		t.Errorf("the hex field: %+v, %v", info, ok)
	}

	// Choosing a swatch sets the color.
	if err := tt.Click("Blue"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if color != Hex("#3b82f6") {
		t.Errorf("the color is %v", color)
	}

	// Outside a frame the view builds no well.
	p := tt.rt.public
	if ColorWell(p, &color) != (Element{}) {
		t.Error("a well was built outside a frame")
	}
}
