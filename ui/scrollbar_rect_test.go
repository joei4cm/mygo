package ui

import "testing"

// The track of a scroll bar runs along the edge of the view it scrolls,
// and the thumb down the middle of the track, keeping to the insets of
// the style at both ends of the scroll.
func TestScrollbarTrackRect(t *testing.T) {
	style := ScrollbarStyle{TrackWidth: 8, ThumbWidth: 4}
	// A vertical bar runs down the right edge of the view.
	v := Scrollbar{Box: Rect{X: 10, Y: 20, W: 200, H: 100}, View: 50, Content: 200, Style: style}
	if got := v.trackRect(style); got != (Rect{X: 202, Y: 20, W: 8, H: 100}) {
		t.Errorf("the vertical track is %v", got)
	}
	// A horizontal one along the bottom edge.
	h := Scrollbar{Box: Rect{X: 10, Y: 20, W: 200, H: 100}, Horizontal: true, View: 50, Content: 200, Style: style}
	if got := h.trackRect(style); got != (Rect{X: 10, Y: 112, W: 200, H: 8}) {
		t.Errorf("the horizontal track is %v", got)
	}

	// The thumb stays in the track, in its middle across, and at either
	// end of the scroll keeps to the inset of two DIPs.
	for _, off := range []float32{0, 75, 150} {
		v.Offset = off
		th, track := v.thumbRect(style), v.trackRect(style)
		if th.X != track.X+(track.W-th.W)/2 || th.W != style.ThumbWidth {
			t.Errorf("the thumb at %v sits at %v, not in the middle of the track", off, th)
		}
		if th.Y < track.Y+2 || th.Y+th.H > track.Y+track.H-2 {
			t.Errorf("the thumb at %v is %v, out of the track's insets", off, th)
		}
	}
	v.Offset = 0
	if th := v.thumbRect(style); th.Y != v.Box.Y+2 {
		t.Errorf("the thumb at the top of the scroll is at %v, want an inset from the track's top", th)
	}
	v.Offset = v.Content - v.View
	if th := v.thumbRect(style); th.Y+th.H != v.Box.Y+v.Box.H-2 {
		t.Errorf("the thumb at the end of the scroll ends at %v, want an inset from the track's end", th)
	}

	// A horizontal thumb runs across the middle of its track, an inset
	// from the track's start where it starts.
	h.Offset = 0
	if th := h.thumbRect(style); th.X != h.Box.X+2 || th.Y != 114 || th.H != style.ThumbWidth {
		t.Errorf("the horizontal thumb at the start of the scroll is %v", th)
	}
	// A flush style runs without the insets of its own.
	flush := ScrollbarStyle{TrackWidth: 8, ThumbWidth: 4, Flush: true}
	v.Style = flush
	v.Offset = 0
	if th := v.thumbRect(flush); th.Y != v.Box.Y {
		t.Errorf("the flush thumb at the start of the scroll is at %v, want the track's top", th)
	}
}
