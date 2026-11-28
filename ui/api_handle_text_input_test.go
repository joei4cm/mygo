package ui

import "testing"

// HandleTextInput connects an element to the system's text services: the
// pointer over it is the text cursor, a click gives it the focus and it
// takes typed text. A nil client, or a nil client of a pointer kind,
// disconnects it, and it connects again.
func TestHandleTextInputAPI(t *testing.T) {
	real := &primitiveClient{text: "hi", selection: TextInputSelection{Range: TextInputRange{Start: 2, End: 2}}}
	var client TextInputClient = real
	tt := NewTester(func(c *Context) {
		Box(c).Size(250, 80).HandleTextInput(client)
	}, 300, 100)
	defer tt.rt.close()

	tt.Move(20, 20)
	if got := tt.Cursor(); got != CursorText {
		t.Errorf("over the element, the cursor is %v, not CursorText", got)
	}
	tt.ClickAt(20, 20)
	if active, _ := tt.TextInputActive(); !active {
		t.Fatal("the element takes no text")
	}
	tt.Type("!")
	if real.text != "hi!" {
		t.Fatalf("typed text gave %q, want %q", real.text, "hi!")
	}
	if real.changes == 0 {
		t.Error("the client heard no change")
	}

	// A nil client disconnects it: it no longer takes text, and typing
	// goes nowhere.
	client = nil
	tt.Frame()
	if active, _ := tt.TextInputActive(); active {
		t.Error("a nil client left the element taking text")
	}
	tt.Type("?")
	if real.text != "hi!" {
		t.Errorf("a disconnected client heard typing: %q", real.text)
	}

	// A nil client of a pointer kind disconnects it the same.
	client = (*primitiveClient)(nil)
	tt.Frame()
	if active, _ := tt.TextInputActive(); active {
		t.Error("a typed nil client left the element taking text")
	}

	// A client connects it again.
	client = real
	tt.Frame()
	tt.ClickAt(20, 20)
	if active, _ := tt.TextInputActive(); !active {
		t.Error("the element takes no text with its client back")
	}
	tt.Type("~")
	if real.text != "hi!~" {
		t.Errorf("typed text gave %q, want %q", real.text, "hi!~")
	}
}

func TestHandleTextInputAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.HandleTextInput(&primitiveClient{}).Valid() {
		t.Error("an absent element connected to the text services")
	}
}
