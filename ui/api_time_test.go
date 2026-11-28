package ui

import (
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
)

// The public API of a field of a time: it shows the time in the region's
// own way, and the arrows step it.
func TestTimeInputAPI(t *testing.T) {
	alarm := time.Date(2026, 10, 4, 7, 30, 0, 0, time.UTC)
	changes := 0
	tt := NewTester(func(c *Context) {
		if TimeInput(c, &alarm).Label("Alarm").Changed() {
			changes++
		}
	}, 300, 100)
	tt.SetRegion(i18n.Region{Locale: "en-GB", FirstWeekday: time.Monday, HourCycle: 24})

	// The time shows in the region's own way.
	if !tt.HasText("07") || !tt.HasText("30") {
		t.Fatalf("the texts are %q", tt.Texts())
	}

	// The arrows step the hours, which is a change.
	tt.Click("07")
	tt.Key(0, KeyUp)
	if alarm.Hour() != 8 || changes != 1 {
		t.Fatalf("Up: %v, %d changes", alarm, changes)
	}
	tt.Key(0, KeyDown)
	if alarm.Hour() != 7 || changes != 2 {
		t.Fatalf("Down: %v, %d changes", alarm, changes)
	}

	// Digits typed set the hours, and go on to the minutes.
	tt.Key(0, Key1)
	tt.Key(0, Key5)
	if alarm.Hour() != 15 {
		t.Fatalf("1 and 5: %v", alarm)
	}
	tt.Key(0, Key4)
	tt.Key(0, Key5)
	if alarm.Minute() != 45 {
		t.Fatalf("4 and 5 in the minutes: %v", alarm)
	}
	if y, m, d := alarm.Date(); y != 2026 || m != 10 || d != 4 {
		t.Errorf("the date changed: %v", alarm)
	}

	// Outside a frame the view builds no field.
	p := tt.rt.public
	if TimeInput(p, &alarm).Valid() {
		t.Error("a field was built outside a frame")
	}
}
