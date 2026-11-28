package richtext

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDocumentFormat(t *testing.T) {
	d := NewDocument(NewText("hello world"))
	d.SetSelection(0, 5)
	red := ui.RGB(200, 30, 30)
	d.Format(func(a *Attributes) { a.Color = red })

	if got := d.Text().AttributesAt(2).Color; got != red {
		t.Fatalf("selection color = %v, want %v", got, red)
	}
	if got := d.Text().AttributesAt(7).Color; got == red {
		t.Fatalf("unselected text changed color: %v", got)
	}
	if a, c := d.Selection(); a != 0 || c != 5 {
		t.Fatalf("Format moved selection to %d, %d", a, c)
	}
	if !d.CanUndo() {
		t.Fatal("Format did not make an undo step")
	}
	d.Undo()
	if got := d.Text().AttributesAt(2).Color; got == red {
		t.Fatal("undo did not revert Format")
	}

	// At a caret, Format changes typing attributes and strips any Attachment.
	d.SetSelection(5, 5)
	d.Format(func(a *Attributes) {
		a.Color = red
		a.Attachment = &Attachment{Name: "ignored"}
	})
	ta := d.TypingAttributes()
	if ta.Color != red {
		t.Fatalf("typing color = %v, want %v", ta.Color, red)
	}
	if ta.Attachment != nil {
		t.Fatalf("typing attributes kept Attachment: %+v", ta.Attachment)
	}
	d.InsertText("!")
	if got := d.Text().AttributesAt(5).Color; got != red {
		t.Fatalf("inserted rune color = %v, want %v", got, red)
	}
}

func TestDocumentToggleStrikethroughAndCode(t *testing.T) {
	d := NewDocument(NewText("abcdef"))
	d.SetSelection(0, 3)
	d.ToggleStrikethrough()
	if !d.Text().AttributesAt(1).Strikethrough || d.Text().AttributesAt(4).Strikethrough {
		t.Fatal("ToggleStrikethrough did not turn on for [0,3)")
	}
	// Mixed selection turns all of it on.
	d.SetSelection(0, 6)
	d.ToggleStrikethrough()
	if !d.Text().AttributesAt(4).Strikethrough {
		t.Fatal("ToggleStrikethrough on mixed selection did not turn all on")
	}
	// All on turns it off.
	d.ToggleStrikethrough()
	if d.Text().AttributesAt(1).Strikethrough || d.Text().AttributesAt(4).Strikethrough {
		t.Fatal("ToggleStrikethrough on fully struck selection did not turn off")
	}

	d.SetSelection(0, 3)
	d.ToggleCode()
	if !d.Text().AttributesAt(1).Code || d.Text().AttributesAt(4).Code {
		t.Fatal("ToggleCode did not turn on for [0,3)")
	}
	d.SetSelection(0, 6)
	d.ToggleCode()
	if !d.Text().AttributesAt(4).Code {
		t.Fatal("ToggleCode on mixed selection did not turn all on")
	}
	d.ToggleCode()
	if d.Text().AttributesAt(1).Code || d.Text().AttributesAt(4).Code {
		t.Fatal("ToggleCode on fully code selection did not turn off")
	}
}

func TestDocumentToggleCodeBlock(t *testing.T) {
	d := NewDocument(TextOf(
		Paragraph{
			Style: ParagraphStyle{Heading: 2, List: Tasks, Checked: true},
			Runs:  []Run{{Text: "first"}},
		},
		Paragraph{
			Runs: []Run{{Text: "second"}},
		},
	))
	// First paragraph becomes a code block, clearing heading, list and checked.
	d.SetSelection(0, 5)
	d.ToggleCodeBlock()
	p0 := d.Text().ParagraphStyleAt(0)
	if !p0.CodeBlock || p0.Heading != 0 || p0.List != NoList || p0.Checked {
		t.Fatalf("first paragraph after ToggleCodeBlock = %+v", p0)
	}
	if d.Text().ParagraphStyleAt(6).CodeBlock {
		t.Fatal("second paragraph unexpectedly became a code block")
	}

	// Selecting across a code block and a plain paragraph turns both into code blocks.
	d.SetSelection(0, 12)
	d.ToggleCodeBlock()
	if !d.Text().ParagraphStyleAt(0).CodeBlock || !d.Text().ParagraphStyleAt(6).CodeBlock {
		t.Fatal("ToggleCodeBlock on mixed paragraphs did not turn all on")
	}

	// Selecting all code blocks turns them off.
	d.ToggleCodeBlock()
	if d.Text().ParagraphStyleAt(0).CodeBlock || d.Text().ParagraphStyleAt(6).CodeBlock {
		t.Fatal("ToggleCodeBlock on all code blocks did not turn off")
	}
}

func TestDocumentInsertAttachment(t *testing.T) {
	d := NewDocument(NewText("hi there"))
	d.SetSelection(0, 8)
	d.ToggleBold()
	d.SetLink("https://example.com")
	d.SetSelection(2, 8)

	att := &Attachment{Name: "photo.png", URL: "file:///photo.png", Width: 64, Height: 48}
	d.InsertAttachment(att)

	if got := d.Text().String(); got != "hi"+string(ObjectReplacement) {
		t.Fatalf("text = %q, want %q", got, "hi"+string(ObjectReplacement))
	}
	if a, c := d.Selection(); a != 3 || c != 3 {
		t.Fatalf("selection = %d, %d, want 3, 3", a, c)
	}
	attrs := d.Text().AttributesAt(2)
	if attrs.Attachment != att {
		t.Fatalf("attachment = %+v, want %+v", attrs.Attachment, att)
	}
	if !attrs.Bold() {
		t.Error("attachment did not keep extending typing attribute (Bold)")
	}
	if attrs.Link != "" {
		t.Errorf("attachment unexpectedly extended Link: %q", attrs.Link)
	}
	if !d.CanUndo() {
		t.Fatal("InsertAttachment did not make an undo step")
	}
	d.Undo()
	if got := d.Text().String(); got != "hi there" {
		t.Fatalf("after undo, text = %q, want %q", got, "hi there")
	}
}
