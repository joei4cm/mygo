package ui

import "testing"

// The public API of a select of your own: a trigger that opens a popup of
// options, which the app opens and closes, highlights, and reads the
// highlight of.
func TestSelectOwnAPI(t *testing.T) {
	size := "M"
	sizes := []string{"S", "M", "L", "XL"}
	highlighted := map[string]bool{}
	var open bool
	var wantOpen bool
	var highlight string
	tt := NewTester(func(c *Context) {
		sel := SelectBase(c, &size)
		sel.Trigger.Label("size").Padding(4).Children(func() { Text(c, "Size: "+size) })
		clear(highlighted)
		open = sel.Open()
		sel.Popup(func(panel Element) {
			panel.Label("popup")
			for _, s := range sizes {
				item := sel.Item(s).Padding(4).Label("item " + s)
				if item.Highlighted() {
					highlighted[s] = true
				}
				item.Children(func() { Text(c, s) })
			}
		})
		if highlight != "" {
			sel.Highlight(highlight)
		}
		if open != wantOpen {
			sel.SetOpen(wantOpen)
		}
	}, 300, 300)
	defer tt.rt.close()
	tt.Frame()
	if open || tt.HasText("XL") {
		t.Fatalf("the popup shows before it is asked: open %v", open)
	}
	// Opened, it shows the options with the highlight on the choice.
	wantOpen = true
	tt.Frame()
	if !open || !tt.HasText("XL") || !highlighted["M"] || len(highlighted) != 1 {
		t.Fatalf("open %v, text %v, highlights %v; want open, the options, M only", open, tt.HasText("XL"), highlighted)
	}
	// The app moves the highlight as the user types.
	highlight = "L"
	tt.Frame()
	if !highlighted["L"] || len(highlighted) != 1 {
		t.Errorf("the highlight asked for L, on %v", highlighted)
	}
	// With the focus on the trigger, Enter chooses the highlighted one
	// and closes the popup.
	tt.Key(0, KeyTab)
	wantOpen = false
	tt.Key(0, KeyEnter)
	if size != "L" || open || tt.HasText("XL") {
		t.Errorf("Enter chose %q, open %v, text %v; want L, closed", size, open, tt.HasText("XL"))
	}
}

// Outside a frame a select without its trigger does nothing.
func TestSelectOwnAPIOutsideFrame(t *testing.T) {
	var absent SelectParts[string]
	absent.Highlight("S")
	absent.SetOpen(true)
	if absent.Open() {
		t.Error("an absent select is open")
	}
	if absent.Item("S").Valid() {
		t.Error("an absent select made an option")
	}
	var noElement Element
	if noElement.Highlighted() {
		t.Error("an absent element is highlighted")
	}
}
