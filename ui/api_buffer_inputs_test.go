package ui

import (
	"testing"
)

// The public API of a text input over a buffer: typing into it writes to
// the buffer, and so does typing into the one without a look.
func TestTextInputBufferAPI(t *testing.T) {
	buf := NewTextBuffer("")
	tt := NewTester(func(c *Context) {
		TextInputBuffer(c, buf).Label("Name")
	}, 400, 300)
	tt.Frame()
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Ada")
	if buf.String() != "Ada" {
		t.Errorf("the buffer is %q", buf.String())
	}
	// Enter in a single line does not type a new line.
	tt.Key(0, KeyEnter)
	if buf.String() != "Ada" {
		t.Errorf("the buffer is %q", buf.String())
	}

	// The same without a look.
	buf2 := NewTextBuffer("")
	tt2 := NewTester(func(c *Context) {
		TextInputBufferBase(c, buf2).Label("Name")
	}, 400, 300)
	tt2.Frame()
	if err := tt2.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt2.Type("Bob")
	if buf2.String() != "Bob" {
		t.Errorf("the buffer is %q", buf2.String())
	}

	// Outside a frame the view builds no input.
	p := tt.rt.public
	if TextInputBuffer(p, buf) != (Element{}) || TextInputBufferBase(p, buf) != (Element{}) {
		t.Error("an input was built outside a frame")
	}
}

// The public API of a text area over a buffer: typing into it writes to
// the buffer, and Enter goes to the next line.
func TestTextAreaBufferBaseAPI(t *testing.T) {
	buf := NewTextBuffer("")
	tt := NewTester(func(c *Context) {
		TextAreaBufferBase(c, buf).Label("Notes")
	}, 400, 300)
	tt.Frame()
	if err := tt.Click("Notes"); err != nil {
		t.Fatal(err)
	}
	tt.Type("first")
	tt.Key(0, KeyEnter)
	tt.Type("second")
	if buf.String() != "first\nsecond" {
		t.Errorf("the buffer is %q", buf.String())
	}

	// Outside a frame the view builds no area.
	p := tt.rt.public
	if TextAreaBufferBase(p, buf) != (Element{}) {
		t.Error("an area was built outside a frame")
	}
}
