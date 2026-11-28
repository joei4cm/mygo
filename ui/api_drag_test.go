package ui

import (
	"slices"
	"testing"
)

// The public API of dragging a value within a window, of an element the
// pointer drags about, and of files dragged in from another app.
func TestDragAPI(t *testing.T) {
	t.Run("Dragging", func(t *testing.T) {
		var dragging, missed bool
		var at Point
		tt := NewTester(func(c *Context) {
			src := Box(c).Size(100, 40).Drag("hello")
			dragging = src.Dragging()
			if p, ok := src.DragMissed(); ok {
				missed, at = true, p
			}
		}, 240, 240)

		tt.Press(50, 20)
		tt.Move(90, 60)
		tt.Frame()
		if !dragging {
			t.Fatal("the source does not report the drag it started")
		}

		// Let go over nothing: the drag is missed, where it ended.
		tt.Release(200, 200)
		tt.Frame()
		if dragging {
			t.Error("the source still reports the drag it let go")
		}
		if !missed || at != (Point{X: 200, Y: 200}) {
			t.Errorf("the drag missed %v at %v, want true at (200, 200)", missed, at)
		}

		// A drag given up with Escape is not missed.
		missed, at = false, Point{}
		tt.Press(50, 20)
		tt.Move(90, 60)
		tt.Frame()
		if !dragging {
			t.Fatal("the source does not report the drag it started again")
		}
		tt.Key(0, KeyEscape)
		tt.Release(200, 200)
		tt.Frame()
		if missed {
			t.Error("the drag given up with Escape is missed")
		}
	})

	t.Run("DragOffset", func(t *testing.T) {
		var off Point
		var vx, vy float32
		var shown Point
		tt := NewTester(func(c *Context) {
			el := Box(c).Size(80, 40).Draggable(&off, DraggableOptions{})
			vx, vy = el.DragVelocity()
			shown = el.DragOffset()
		}, 240, 240)

		tt.Press(40, 20)
		tt.Advance(16 * 1e6)
		tt.Move(80, 60)
		tt.Frame()
		if shown != (Point{X: 40, Y: 40}) {
			t.Errorf("the element shows %v from its place, want (40, 40)", shown)
		}
		if vx == 0 || vy == 0 {
			t.Errorf("the pointer moves the element at (%v, %v), want both non-zero", vx, vy)
		}

		tt.Release(80, 60)
		tt.Frame()
		if off != (Point{X: 40, Y: 40}) {
			t.Errorf("the element rests at %v, where the pointer let it go", off)
		}
	})

	t.Run("DroppedFiles", func(t *testing.T) {
		var over, after bool
		var dropped [][]string
		tt := NewTester(func(c *Context) {
			target := Box(c).Size(200, 120)
			over = target.FileDragOver()
			dropped = append(dropped, target.DroppedFiles())
		}, 240, 240)
		tt.Frame()

		if ok := tt.DragFiles(50, 30); !ok {
			t.Fatal("no element takes the files dragged over it")
		}
		tt.Frame()
		if !over {
			t.Error("the target does not see the files over it")
		}

		if ok := tt.DropFiles(50, 30, "/tmp/report.pdf"); !ok {
			t.Fatal("the files dropped on the target are not taken")
		}
		tt.Frame()
		var took []string
		for _, files := range dropped {
			if slices.Equal(files, []string{"/tmp/report.pdf"}) {
				took = files
			}
		}
		if took == nil {
			t.Errorf("the target takes none of the files dropped on it: %q", dropped)
		}

		tt.DragFilesOut()
		tt.Frame()
		if after {
			t.Error("the target still sees the files dragged out of the window")
		}
	})

	t.Run("Outside a frame", func(t *testing.T) {
		var e Element
		if got := e.Draggable(nil, DraggableOptions{}); got != (Element{}) {
			t.Errorf("Draggable outside a frame gave %v, want nothing", got)
		}
		if p, ok := e.DragMissed(); ok || p != (Point{}) {
			t.Errorf("DragMissed outside a frame gave (%v, %v)", p, ok)
		}
		if e.Dragging() {
			t.Error("Dragging outside a frame is true")
		}
		if p := e.DragOffset(); p != (Point{}) {
			t.Errorf("DragOffset outside a frame gave %v", p)
		}
		if vx, vy := e.DragVelocity(); vx != 0 || vy != 0 {
			t.Errorf("DragVelocity outside a frame gave (%v, %v)", vx, vy)
		}
		if f := e.DroppedFiles(); f != nil {
			t.Errorf("DroppedFiles outside a frame gave %q", f)
		}
		if e.FileDragOver() {
			t.Error("FileDragOver outside a frame is true")
		}
	})
}
