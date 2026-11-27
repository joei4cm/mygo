package ui

import "testing"

// The buffer's version counts every edit and every snapshot adopted, and
// a snapshot taken before them keeps its own.
func TestTextBufferVersion(t *testing.T) {
	b := NewTextBuffer("abc")
	v0 := b.Version()
	snap := b.Snapshot()
	if v := snap.Version(); v != v0 {
		t.Fatalf("the snapshot's version is %d, want the buffer's %d", v, v0)
	}

	b.Replace(1, 2, "X")
	if v := b.Version(); v != v0+1 {
		t.Errorf("after a replace the version is %d, want %d", v, v0+1)
	}
	if v := snap.Version(); v != v0 {
		t.Errorf("the old snapshot's version moved to %d, want %d", v, v0)
	}

	b.Set("hello")
	if v := b.Version(); v != v0+2 {
		t.Errorf("after Set the version is %d, want %d", v, v0+2)
	}

	b.Restore(snap)
	if v := b.Version(); v != v0+3 {
		t.Errorf("after Restore the version is %d, want %d", v, v0+3)
	}
	if got := b.String(); got != "abc" {
		t.Errorf("Restore gave %q, want the snapshot's text", got)
	}
}

// UTF16Offset and RuneOffset map between rune indexes and UTF-16 units,
// counting a rune outside the BMP as two units.
func TestTextBufferUTF16Offsets(t *testing.T) {
	// "a𐐀b": rune 1 is a surrogate pair, two units long.
	b := NewTextBuffer("a𐐀b")
	if l := b.UTF16Len(); l != 4 {
		t.Fatalf("the buffer's UTF-16 length is %d, want 4", l)
	}
	for _, c := range []struct {
		index, want int
	}{
		{0, 0},
		{1, 1}, // past "a"
		{2, 3}, // past the pair
		{3, 4}, // past "b"
	} {
		if got := b.UTF16Offset(c.index); got != c.want {
			t.Errorf("UTF16Offset(%d) = %d, want %d", c.index, got, c.want)
		}
	}
	for _, c := range []struct {
		unit, want int
	}{
		{0, 0},
		{1, 1}, // inside the pair rounds down to its rune
		{2, 1},
		{3, 2}, // past the pair
		{4, 3},
		{9, 3}, // past the end clamps
	} {
		if got := b.RuneOffset(c.unit); got != c.want {
			t.Errorf("RuneOffset(%d) = %d, want %d", c.unit, got, c.want)
		}
	}
}
