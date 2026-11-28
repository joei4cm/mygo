package ui

import (
	"slices"
	"testing"
)

// ObserveInput sees input in its descendants without taking it from them.
func TestObserveInputAPI(t *testing.T) {
	var value string
	var events []InputEvent
	tt := NewTester(func(c *Context) {
		Column(c).ObserveInput(func(ev InputEvent) { events = append(events, ev) }).Children(func() {
			TextInput(c, &value).Label("Field")
		})
	}, 220, 100)
	defer tt.rt.close()
	if err := tt.Click("Field"); err != nil {
		t.Fatal(err)
	}
	events = nil // discard the pointer events from focusing the field
	tt.TypeKey(0, KeyB, "b")
	if value != "b" {
		t.Fatalf("the observed input took the key from the field: value=%q", value)
	}
	var got []InputKind
	for _, ev := range events {
		got = append(got, ev.Kind)
	}
	want := []InputKind{InputKeyDown, InputText, InputKeyUp}
	if !slices.Equal(got, want) {
		t.Fatalf("observed event kinds %v, want %v", got, want)
	}
	if events[0].Key != KeyB || events[1].Text != "b" || events[2].Key != KeyB {
		t.Fatalf("observed events lost their key or text: %+v", events)
	}
}

func TestObserveInputAPIPointer(t *testing.T) {
	var events []InputEvent
	tt := NewTester(func(c *Context) {
		Box(c).Size(100, 80).ObserveInput(func(ev InputEvent) { events = append(events, ev) }).Children(func() {
			Button(c, "Inside")
		})
	}, 160, 100)
	defer tt.rt.close()
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	var got []InputKind
	for _, ev := range events {
		got = append(got, ev.Kind)
	}
	want := []InputKind{InputPointerDown, InputPointerUp}
	if !slices.Equal(got, want) {
		t.Fatalf("observed pointer event kinds %v, want %v", got, want)
	}
}

func TestObserveInputAPIOutsideFrame(t *testing.T) {
	var absent Element
	if absent.ObserveInput(func(InputEvent) {}).Valid() {
		t.Error("an absent element started observing input")
	}
}
