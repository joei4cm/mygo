package ui

import (
	"testing"

	transfer "github.com/egoist/mygo/transfer"
)

// The public API of the Tab order: a container that takes the focus to
// hold its keys stands in it where it is asked to — Tab passes over one
// that is no tab stop, and stops at one that is.
func TestTabStopAPI(t *testing.T) {
	stop := true
	tt := NewTester(func(c *Context) {
		Column(c).Children(func() {
			Button(c, "First")
			panel := Box(c).Size(100, 30).Label("Panel")
			if stop {
				panel = panel.TabStop(true)
			} else {
				panel = panel.TabStop(false)
			}
			panel.Children(func() {
				Button(c, "Inside")
			})
			Button(c, "After")
		})
	}, 200, 140)
	defer tt.rt.close()
	tt.Frame()
	tt.Key(0, KeyTab)
	if !tt.Focused("First") {
		t.Fatal("Tab did not stop at the first button")
	}
	tt.Key(0, KeyTab)
	if stop && !tt.Focused("Panel") {
		t.Error("Tab did not stop at the panel")
	}
	if !stop && !tt.Focused("Inside") {
		t.Errorf("Tab did not pass the panel: %v", tt.Focused("Inside"))
	}
	tt.Key(0, KeyTab)
	if !tt.Focused("After") && !tt.Focused("Inside") {
		t.Error("Tab left the column")
	}
	if stop && !tt.Focused("Inside") {
		t.Error("Tab skipped the panel's button")
	}
}

// The public API of what Copy puts on the clipboard: a selectable
// container copies the parts of its selection as the data its fn returns.
func TestCopyDataAPI(t *testing.T) {
	tt := NewTester(func(c *Context) {
		Column(c).Selectable().Padding(12).Gap(8).CopyData(func(parts []SelectedText) transfer.Data {
			s := ""
			for _, p := range parts {
				s += p.Text[p.Start:p.End] + ";"
			}
			return transfer.TextData(s)
		}).Children(func() {
			Text(c, "Alpha beta")
			Text(c, "Gamma")
		})
	}, 300, 120)
	defer tt.rt.close()
	tt.Frame()
	selectBetween(t, tt, "Alpha beta", 6, "Gamma", 3)
	tt.Frame()
	tt.Key(Cmd, KeyC)
	if got := tt.Clipboard(); got != "beta;Gam;" {
		t.Errorf("the copy put %q on the clipboard", got)
	}
}

// Outside a frame neither changes anything.
func TestTabStopAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.TabStop(false).Valid() {
		t.Error("an absent element took a place in the Tab order")
	}
	if absent.CopyData(func([]SelectedText) transfer.Data { return transfer.TextData("") }).Valid() {
		t.Error("an absent element took copy data")
	}
}
