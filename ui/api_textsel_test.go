package ui

import "testing"

// The public API of what a selection looks like and holds: a container
// that selects its text highlights it in the color it sets, and reads the
// parts the selection holds.
func TestSelectionStyleAPI(t *testing.T) {
	red := RGB(220, 30, 30)
	var texts []SelectedText
	tt := NewTester(func(c *Context) {
		col := Column(c).Selectable().Padding(12).Gap(8).SelectionColor(red)
		texts = col.SelectedTexts()
		col.Children(func() {
			Text(c, "Alpha beta")
			Text(c, "Gamma")
		})
	}, 300, 150)
	defer tt.rt.close()
	tt.Frame()
	if len(texts) != 0 {
		t.Fatalf("a selection before there is one: %+v", texts)
	}
	selectBetween(t, tt, "Alpha beta", 6, "Gamma", 3)
	tt.Frame()
	if len(texts) != 2 ||
		texts[0].Text != "Alpha beta" || texts[0].Start != 6 || texts[0].End != 10 ||
		texts[1].Text != "Gamma" || texts[1].Start != 0 || texts[1].End != 3 {
		t.Fatalf("the selection holds %+v", texts)
	}
	reds := 0
	for y := 0; y < 150; y++ {
		for x := 0; x < 300; x++ {
			if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
				reds++
			}
		}
	}
	if reds < 150 {
		t.Errorf("the highlight is not the color set: %d red px", reds)
	}
	// A press away empties the selection.
	tt.ClickAt(150, 34)
	tt.Frame()
	if len(texts) != 0 {
		t.Errorf("the selection holds %+v after a click away", texts)
	}
}

// The public API of the looks of an input: runs of its text in the colors
// the frame styles them in, and its caret in the color it sets.
func TestTextInputStyleAPI(t *testing.T) {
	red := RGB(220, 30, 30)
	text := "hello world"
	styled := true
	tt := NewTester(func(c *Context) {
		in := TextInput(c, &text).Size(220, 24)
		if styled {
			in = in.TextRanges(TextRange{Start: 0, End: 5, Color: red})
		}
		_ = in
	}, 300, 100)
	defer tt.rt.close()
	tt.Frame()
	reds := 0
	for y := 0; y < 40; y++ {
		for x := 0; x < 150; x++ {
			if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
				reds++
			}
		}
	}
	if reds < 40 {
		t.Errorf("the run styled red shows %d red px", reds)
	}
	// Without the styles of the frame, the text is plain.
	styled = false
	tt.Frame()
	reds = 0
	for y := 0; y < 40; y++ {
		for x := 0; x < 150; x++ {
			if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
				reds++
			}
		}
	}
	if reds != 0 {
		t.Errorf("unstyled text shows %d red px", reds)
	}
}

// The caret takes the color the input sets, not the theme's accent.
func TestCaretColorAPI(t *testing.T) {
	red := RGB(220, 30, 30)
	text := "hello world"
	tt := NewTester(func(c *Context) {
		TextInput(c, &text).Size(220, 24).CaretColor(red)
	}, 300, 100)
	defer tt.rt.close()
	tt.Frame()
	for y := 0; y < 40; y++ {
		for x := 0; x < 300; x++ {
			if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
				t.Fatal("a red pixel before the field was edited")
			}
		}
	}
	// Editing puts the caret where the press went, in the color set.
	tt.ClickAt(210, 12)
	tt.Frame()
	reds := 0
	for y := 0; y < 40; y++ {
		for x := 0; x < 300; x++ {
			if c := tt.Image().RGBAAt(x, y); c.A > 200 && c.R > 200 && c.G < 90 && c.B < 90 {
				reds++
			}
		}
	}
	if reds < 8 {
		t.Errorf("the caret shows %d red px", reds)
	}
}

// Outside a frame these styles change nothing and read none.
func TestTextSelAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.SelectionColor(RGB(220, 30, 30)).Valid() {
		t.Error("an absent element took a selection color")
	}
	if absent.CaretColor(RGB(220, 30, 30)).Valid() {
		t.Error("an absent element took a caret color")
	}
	if absent.TextRanges(TextRange{Start: 0, End: 5}).Valid() {
		t.Error("an absent element styled runs")
	}
	if texts := absent.SelectedTexts(); texts != nil {
		t.Errorf("an absent element holds %+v", texts)
	}
}
