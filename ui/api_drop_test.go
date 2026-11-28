package ui

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/internal/platform"
	"github.com/egoist/mygo/transfer"
)

// The public API of a drag and drop of the app's values: dragging a value
// over a target tells it, and dropping it gives it there.
func TestDragAndDropAPI(t *testing.T) {
	var dropped []string
	var over bool
	tt := NewTester(func(c *Context) {
		Column(c).Gap(20).Padding(20).Children(func() {
			Box(c).Size(100, 40).Label("Source").Drag("hello")
			target := Box(c).Size(200, 100).Label("Target")
			if v, ok := Drop[string](target); ok {
				dropped = append(dropped, v)
			}
			_, over = DragOver[string](target)
		})
	}, 400, 400)
	src, _ := tt.Find("Source")
	target, _ := tt.Find("Target")

	// Dragged over the target, it sees the value; dropped on it, it gets
	// it.
	tt.Press(src.X+10, src.Y+10)
	tt.Move(target.X+20, target.Y+20)
	if !over {
		t.Fatal("the target does not see the value over it")
	}
	tt.Release(target.X+20, target.Y+20)
	if !slices.Equal(dropped, []string{"hello"}) || over {
		t.Fatalf("dropped %q, over %v", dropped, over)
	}

	// Outside a frame the view asks nothing of a part that is not its.
	if _, ok := Drop[string](Element{}); ok {
		t.Error("a drop was asked of nothing")
	}
	if _, ok := DragOver[string](Element{}); ok {
		t.Error("a drag was asked of nothing")
	}
}

// The public API of a drag and drop of native data: the target negotiates
// what it takes as the drag hovers, and reads it as it drops.
func TestDropDataAPI(t *testing.T) {
	var got []transfer.Drop
	var results []transfer.Result
	var over bool
	custom := transfer.Format("application/vnd.mygo.test+json")
	data := transfer.New(transfer.NewItem(
		transfer.Bytes(transfer.Text, []byte("fallback")),
		transfer.Lazy(custom, func() ([]byte, error) { return []byte(`{"id":7}`), nil }),
	))
	source := NewTester(func(c *Context) {
		Box(c).Size(100, 50).DragData(data, transfer.DragOptions{Operations: transfer.Copy | transfer.Move, Done: func(r transfer.Result) { results = append(results, r) }})
	}, 200, 200)
	target := NewTester(func(c *Context) {
		zone := Box(c).Size(100, 100)
		opts := transfer.DropOptions{Formats: []transfer.Format{custom}, Operations: transfer.Copy | transfer.Move}
		_, over = DataDragOver(zone, opts)
		if d, ok := DropData(zone, opts); ok {
			got = append(got, d)
		}
	}, 200, 200)

	// A drag from another app hovers, and drops.
	source.Press(10, 10)
	source.Move(20, 20)
	d := &platform.DataDragEvent{Offer: transfer.Offer{Formats: data.Formats(), Operations: transfer.Copy | transfer.Move, Suggested: transfer.Move}}
	if !target.send2(platform.SurfaceEvent{Kind: platform.DataDragOver, X: 10, Y: 10, Drag: d}) {
		t.Fatal("the drag did not reach the target")
	}
	if !over {
		t.Fatal("the target does not see the drag over it")
	}
	d.Data = source.h.dragData
	if !target.send2(platform.SurfaceEvent{Kind: platform.DataDrop, X: 10, Y: 10, Drag: d}) {
		t.Fatal("drop was rejected")
	}
	source.h.finishDataDrag(transfer.Result{Operation: d.Operation})
	source.settle()
	if len(got) != 1 || got[0].Operation != transfer.Move || over || len(results) != 1 {
		t.Fatalf("drops %v, results %v, over %v", got, results, over)
	}
	// Only the format the target asked for is read; the fallback is left.
	if !slices.Equal(got[0].Data.Formats(), []transfer.Format{custom}) {
		t.Fatalf("the drop carries %v", got[0].Data.Formats())
	}

	// Outside a frame the view asks nothing of a part that is not its.
	opts := transfer.DropOptions{Formats: []transfer.Format{custom}}
	if _, ok := DataDragOver(Element{}, opts); ok {
		t.Error("a drag was asked of nothing")
	}
	if _, ok := DropData(Element{}, opts); ok {
		t.Error("a drop was asked of nothing")
	}
}

// The public API of a drag source that asks for its data when the gesture
// starts.
func TestDragDataFromAPI(t *testing.T) {
	var calls int
	data := transfer.New(transfer.NewItem(transfer.Bytes(transfer.Text, []byte("late"))))
	source := NewTester(func(c *Context) {
		Box(c).Size(100, 50).DragDataFrom(func() transfer.Data {
			calls++
			return transfer.New(transfer.NewItem(transfer.Bytes(transfer.Text, []byte("late"))))
		})
	}, 200, 200)

	// Before the drag, nothing is asked for.
	source.Frame()
	if calls != 0 {
		t.Fatalf("the data was asked for %d times before the drag", calls)
	}
	source.Press(10, 10)
	source.Move(20, 20)
	source.settle()
	if calls != 1 {
		t.Fatalf("the data was asked for %d times", calls)
	}
	if formats := source.h.dragData.Formats(); !slices.Equal(formats, []transfer.Format{transfer.Text}) {
		t.Fatalf("the drag carries %v", formats)
	}

	// An absent element drags nothing.
	if (Element{}).Drag("x").Valid() {
		t.Error("nothing dragged a value")
	}
	if (Element{}).DragData(data).Valid() {
		t.Error("nothing dragged data")
	}
	if (Element{}).DragDataFrom(func() transfer.Data { return data }).Valid() {
		t.Error("nothing dragged data from a function")
	}
}
