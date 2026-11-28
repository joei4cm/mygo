package ui

import "testing"

// The public API of one element attached to another: a panel follows its
// anchor, keeps its distance from the window's edges, and tells which way
// it went.
func TestAttachAPI(t *testing.T) {
	var edge float32 = -1
	left, top := float32(2), float32(10)
	var viewAt, viewSelf, drawAt, drawSelf Anchor
	var panel *node
	tt := NewTester(func(c *Context) {
		b := ButtonBase(c).Size(40, 20).Label("anchor").Absolute().Left(left).Top(top)
		Overlay(c, func() {
			p := Box(c).Size(60, 30).Label("panel").AttachTo(b, AnchorBottomStart, AnchorTopStart).AttachedOrigin()
			if edge >= 0 {
				p.AttachEdgeMargin(edge)
			}
			viewAt, viewSelf = p.AttachedAnchors()
			p.Draw(func(*Painter, Rect) { drawAt, drawSelf = p.AttachedAnchors() })
			panel = p.node()
		})
	}, 400, 300)
	defer tt.rt.close()
	find := func() Rect {
		t.Helper()
		r, ok := tt.Find("panel")
		if !ok {
			t.Fatal("no panel")
		}
		return r
	}
	// Below, its left edge on the anchor's, 2 from the window's, it stands
	// 4 from the edge by default, and its origin is its top left.
	tt.Frame()
	if r := find(); r.X != 4 || r.Y != 30 {
		t.Fatalf("below: x %v, y %v, want 4, 30", r.X, r.Y)
	}
	if drawAt != AnchorBottomStart || drawSelf != AnchorTopStart || viewAt != AnchorBottomStart || viewSelf != AnchorTopStart {
		t.Fatalf("below: drawn %v %v, view %v %v", drawAt, drawSelf, viewAt, viewSelf)
	}
	if panel.originX != 0 || panel.originY != 0 {
		t.Errorf("below: origin %v %v, want 0 0", panel.originX, panel.originY)
	}
	// An edge margin of 12 keeps it 12 from the edge; 0, on the anchor.
	edge = 12
	tt.Frame()
	if r := find(); r.X != 12 {
		t.Errorf("edge margin 12: x %v, want 12", r.X)
	}
	edge = 0
	tt.Frame()
	if r := find(); r.X != 2 {
		t.Errorf("edge margin 0: x %v, want 2, on the anchor", r.X)
	}
	// Without room below, it goes above: drawn and its origin in this
	// frame, read by the view in the next.
	top = 260
	tt.Frame()
	if drawAt != AnchorTopStart || drawSelf != AnchorBottomStart {
		t.Errorf("above, drawn: %v %v", drawAt, drawSelf)
	}
	if viewAt != AnchorBottomStart || viewSelf != AnchorTopStart {
		t.Errorf("above, the view of the frame it moved: %v %v, still those asked for", viewAt, viewSelf)
	}
	if panel.originX != 0 || panel.originY != 1 {
		t.Errorf("above: origin %v %v, want 0 1", panel.originX, panel.originY)
	}
	tt.Frame()
	if viewAt != AnchorTopStart || viewSelf != AnchorBottomStart {
		t.Errorf("above, in the view: %v %v", viewAt, viewSelf)
	}
	// Near the right edge too, it aligns its end with the anchor's.
	left = 350
	tt.Frame()
	tt.Frame()
	if viewAt != AnchorTopEnd || viewSelf != AnchorBottomEnd {
		t.Errorf("above, near the right edge: %v %v", viewAt, viewSelf)
	}
	if panel.originX != 1 || panel.originY != 1 {
		t.Errorf("near the right edge: origin %v %v, want 1 1", panel.originX, panel.originY)
	}
}

// Outside a frame these attachments change nothing.
func TestAttachAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.AttachEdgeMargin(12).Valid() {
		t.Error("an absent element took an edge margin")
	}
	if absent.AttachTo(absent, AnchorBottomStart, AnchorTopStart).Valid() {
		t.Error("an absent element attached to another")
	}
	if absent.AttachedOrigin().Valid() {
		t.Error("an absent element took an origin")
	}
	at, self := absent.AttachedAnchors()
	if at != 0 || self != 0 {
		t.Errorf("an absent element reported anchors %v %v", at, self)
	}
}
