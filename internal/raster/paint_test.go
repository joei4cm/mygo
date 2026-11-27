package raster

import (
	"math"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

var (
	redOpaque    = scene.Color{R: 255, A: 255}
	blueOpaque   = scene.Color{B: 255, A: 255}
	redBlueStops = []scene.Stop{{Color: redOpaque, At: 0}, {Color: blueOpaque, At: 1}}
)

func aboutColors(a, b [4]float32, tol float32) bool {
	for i := range a {
		if math.Abs(float64(a[i]-b[i])) > float64(tol) {
			return false
		}
	}
	return true
}

func TestNewPainter(t *testing.T) {
	c := scene.Color{R: 255, G: 128, B: 64, A: 128}
	op := &scene.Op{
		Paint:    scene.PaintSolid,
		Color:    c,
		Color2:   scene.Color{B: 255, A: 255},
		Gradient: [4]float32{1, 2, 3, 4},
		Rect:     scene.Rect{X: 7, Y: 8, W: 10, H: 10},
	}
	p := newPainter(op, 0.5, nil)
	if want := c.Premul(0.5); p.c1 != want {
		t.Errorf("c1 = %v, want %v", p.c1, want)
	}
	if want := op.Color2.Premul(0.5); p.c2 != want {
		t.Errorf("c2 = %v, want %v", p.c2, want)
	}
	if p.g != op.Gradient || p.origin != [2]float32{7, 8} {
		t.Errorf("gradient %v origin %v, want %v at (7, 8)", p.g, p.origin, op.Gradient)
	}

	// An Oklab mix keeps its colors' Oklab, times the premultiplied alpha.
	op.Paint = scene.PaintOklab
	p = newPainter(op, 1, nil)
	want := scene.Oklab(c)
	for i := range want {
		want[i] *= float32(c.A) / 255
	}
	if !aboutColors([4]float32{p.lab1[0], p.lab1[1], p.lab1[2], 0}, [4]float32{want[0], want[1], want[2], 0}, 1e-4) {
		t.Errorf("lab1 = %v, want %v", p.lab1, want)
	}
}

func TestNewPainterStops(t *testing.T) {
	stops := []scene.Stop{
		{Color: scene.Color{R: 255, A: 255}, At: 0},
		{Color: scene.Color{G: 255, A: 255}, At: 0.25},
		{Color: scene.Color{B: 255, A: 255}, At: 0.5},
		{Color: scene.Color{R: 255, G: 255, A: 255}, At: 1},
	}
	op := &scene.Op{Paint: scene.PaintLinearStops, Stops: [2]int32{1, 3}}
	p := newPainter(op, 0.5, stops)
	if len(p.stops) != 2 {
		t.Fatalf("the painter took %d stops, want the two of 1 to 3", len(p.stops))
	}
	want := stops[1].Color.Premul(0.5)
	if p.stops[0].c != want || p.stops[0].at != 0.25 {
		t.Errorf("the first stop is %v at %v, want %v at 0.25", p.stops[0].c, p.stops[0].at, want)
	}
	wantLab := scene.Oklab(stops[1].Color)
	for i := range wantLab {
		wantLab[i] *= 0.5
	}
	if p.stops[0].lab != wantLab {
		t.Errorf("the first stop's lab is %v, want %v", p.stops[0].lab, wantLab)
	}

	// An empty or backwards range takes no stops; one that runs past the
	// scene's stops ends at the last, and one longer than MaxStops is cut.
	for _, c := range []struct {
		name     string
		lo, hi   int32
		all      []scene.Stop
		minStops int
	}{
		{"empty", 2, 2, stops, 0},
		{"backwards", 3, 1, stops, 0},
		{"before the first", -1, 2, stops, 0},
		{"past the scene's stops", 2, 99, stops, 0},
		{"more than fit", 0, 34, make([]scene.Stop, scene.MaxStops+2), scene.MaxStops},
	} {
		op := &scene.Op{Paint: scene.PaintLinearStops, Stops: [2]int32{c.lo, c.hi}}
		p := newPainter(op, 1, c.all)
		if len(p.stops) != c.minStops {
			t.Errorf("%s: the painter took %d stops, want %d", c.name, len(p.stops), c.minStops)
		}
	}
}

func TestPainterVisible(t *testing.T) {
	op := &scene.Op{Paint: scene.PaintLinearStops, Stops: [2]int32{0, 2}}
	all := []scene.Stop{{Color: scene.Color{A: 255}, At: 0}, {Color: scene.Color{A: 255}, At: 1}}
	if p := newPainter(op, 1, all); !p.visible() {
		t.Error("a gradient of opaque stops shows nowhere")
	}
	transparent := []scene.Stop{{Color: scene.Color{}, At: 0}, {Color: scene.Color{}, At: 1}}
	if p := newPainter(op, 1, transparent); p.visible() {
		t.Error("a gradient of transparent stops shows")
	}

	op.Paint = scene.PaintSolid
	op.Color = redOpaque
	if p := newPainter(op, 1, nil); !p.visible() {
		t.Error("an opaque solid shows nowhere")
	}
	op.Color = scene.Color{}
	if p := newPainter(op, 1, nil); p.visible() {
		t.Error("a transparent solid shows")
	}

	// A mix shows when only its second color does.
	op.Paint = scene.PaintLinear
	op.Color = scene.Color{}
	op.Color2 = blueOpaque
	if p := newPainter(op, 1, nil); !p.visible() {
		t.Error("a mix of a transparent first color shows nowhere")
	}
	op.Color2 = scene.Color{}
	if p := newPainter(op, 1, nil); p.visible() {
		t.Error("a mix of two transparent colors shows")
	}
}

func TestPainterAtSolidAndStripes(t *testing.T) {
	op := &scene.Op{Paint: scene.PaintSolid, Color: redOpaque}
	p := newPainter(op, 1, nil)
	if c := p.at(12, 34); c != redOpaque.Premul(1) {
		t.Errorf("a solid's fill is %v, want its color", c)
	}

	// Stripes along x of a stripe three wide, ten apart: inside a stripe
	// the first color, in the gap the second, wrapping around the period.
	op.Paint = scene.PaintStripes
	op.Color = redOpaque
	op.Color2 = blueOpaque
	op.Gradient = [4]float32{1, 0, 3, 10}
	p = newPainter(op, 1, nil)
	if c := p.at(1, 0); c != [4]float32{1, 0, 0, 1} {
		t.Errorf("in a stripe the fill is %v, want the first color", c)
	}
	if c := p.at(6, 0); c != [4]float32{0, 0, 1, 1} {
		t.Errorf("in a gap the fill is %v, want the second color", c)
	}
	if c := p.at(-1, 0); c != [4]float32{0, 0, 1, 1} {
		t.Errorf("past the period the fill wraps to %v, want the second color", c)
	}
}

func TestPainterAtStops(t *testing.T) {
	op := &scene.Op{Paint: scene.PaintLinearStops, Gradient: [4]float32{0, 0, 10, 0}, Stops: [2]int32{0, 2}}
	p := newPainter(op, 1, redBlueStops)
	if c := p.at(2.5, 0); c != [4]float32{0.75, 0, 0.25, 1} {
		t.Errorf("a quarter along the gradient is %v, want three quarters red", c)
	}
	if c := p.at(-5, 0); c != [4]float32{1, 0, 0, 1} {
		t.Errorf("before the first stop the fill is %v, want the first stop's color", c)
	}
	if c := p.at(15, 0); c != [4]float32{0, 0, 1, 1} {
		t.Errorf("after the last stop the fill is %v, want the last stop's color", c)
	}

	// A conic takes its place around the center: the top of the circle is
	// the first stop, the bottom halfway to the second.
	op.Paint = scene.PaintConic
	op.Gradient = [4]float32{0, 0, 0, 0}
	p = newPainter(op, 1, redBlueStops)
	if c := p.at(0, -1); c != [4]float32{1, 0, 0, 1} {
		t.Errorf("at the start angle the fill is %v, want the first stop's color", c)
	}
	if c := p.at(0, 1); !aboutColors(c, [4]float32{0.5, 0, 0.5, 1}, 1e-6) {
		t.Errorf("half a turn later the fill is %v, want the stops mixed", c)
	}

	// A radial takes its place by the distance from the center.
	op.Paint = scene.PaintRadial
	op.Gradient = [4]float32{0, 0, 10, 0}
	p = newPainter(op, 1, redBlueStops)
	if c := p.at(3, 4); !aboutColors(c, [4]float32{0.5, 0, 0.5, 1}, 1e-6) {
		t.Errorf("halfway out of the circle the fill is %v, want the stops mixed", c)
	}
}

func TestPainterAtLinear(t *testing.T) {
	op := &scene.Op{Paint: scene.PaintLinear, Color: redOpaque, Color2: blueOpaque, Gradient: [4]float32{0, 0, 10, 0}}
	p := newPainter(op, 1, nil)
	if c := p.at(2.5, 0); c != [4]float32{0.75, 0, 0.25, 1} {
		t.Errorf("a quarter along the gradient is %v, want three quarters red", c)
	}
	// A gradient of two ends at the same place takes its first color.
	op.Gradient = [4]float32{5, 5, 5, 5}
	p = newPainter(op, 1, nil)
	if c := p.at(12, 34); c != [4]float32{1, 0, 0, 1} {
		t.Errorf("a gradient of nowhere is %v, want the first color", c)
	}
}

func TestPainterAtOklab(t *testing.T) {
	op := &scene.Op{Paint: scene.PaintOklab, Color: redOpaque, Color2: blueOpaque, Gradient: [4]float32{0, 0, 10, 0}}
	p := newPainter(op, 1, nil)

	// At a stop the color is that color, back from Oklab unchanged.
	if c := p.at(10, 0); !aboutColors(c, [4]float32{0, 0, 1, 1}, 1e-4) {
		t.Errorf("at the second end the fill is %v, want the second color", c)
	}

	// Halfway the Oklab of the two labs mixed; the alpha comes premultiplied.
	half := scene.Oklab(redOpaque)
	halfL2 := scene.Oklab(blueOpaque)
	for i := range half {
		half[i] = (half[i] + halfL2[i]) / 2
	}
	rgb := fromOklab(half)
	want := [4]float32{clamp01(rgb[0]), clamp01(rgb[1]), clamp01(rgb[2]), 1}
	if c := p.at(5, 0); !aboutColors(c, want, 1e-4) {
		t.Errorf("halfway the fill is %v, want the Oklab mix %v", c, want)
	}

	// Two transparent colors fill nothing, however far along.
	op.Color = scene.Color{}
	op.Color2 = scene.Color{}
	p = newPainter(op, 1, nil)
	if c := p.at(5, 0); c != [4]float32{} {
		t.Errorf("a mix of two transparent colors fills %v, want nothing", c)
	}
}

func TestStopsAt(t *testing.T) {
	// No stops, nothing to fill with.
	p := painter{paint: scene.PaintLinearStops}
	if c := p.stopsAt(0.5); c != [4]float32{} {
		t.Errorf("a gradient of no stops fills %v, want nothing", c)
	}

	one := []paintStop{{c: [4]float32{1, 0, 0, 1}, at: 0.5}}
	p = painter{paint: scene.PaintLinearStops, stops: one}
	if c := p.stopsAt(0); c != one[0].c {
		t.Errorf("before the one stop the fill is %v, want its color", c)
	}
	if c := p.stopsAt(1); c != one[0].c {
		t.Errorf("after the one stop the fill is %v, want its color", c)
	}

	paint := func(mode scene.Paint) painter {
		return painter{paint: mode, stops: []paintStop{
			{c: [4]float32{1, 0, 0, 1}, lab: scene.Oklab(redOpaque), at: 0},
			{c: [4]float32{0, 0, 1, 1}, lab: scene.Oklab(blueOpaque), at: 1},
		}}
	}
	srgb := paint(scene.PaintLinearStops)
	if c := srgb.stopsAt(0.5); c != [4]float32{0.5, 0, 0.5, 1} {
		t.Errorf("halfway in sRGB the fill is %v, want the stops mixed", c)
	}
	conic := paint(scene.PaintConic)
	if c := conic.stopsAt(0.25); c != [4]float32{0.75, 0, 0.25, 1} {
		t.Errorf("a quarter in sRGB the fill is %v, want the stops mixed", c)
	}

	// In Oklab halfway is the Oklab of the two labs mixed, premultiplied.
	ok := paint(scene.PaintLinearStopsOklab)
	lab := scene.Oklab(redOpaque)
	lab2 := scene.Oklab(blueOpaque)
	for i := range lab {
		lab[i] = (lab[i] + lab2[i]) / 2
	}
	rgb := fromOklab(lab)
	want := [4]float32{clamp01(rgb[0]), clamp01(rgb[1]), clamp01(rgb[2]), 1}
	if c := ok.stopsAt(0.5); !aboutColors(c, want, 1e-4) {
		t.Errorf("halfway in Oklab the fill is %v, want %v", c, want)
	}

	// Two transparent stops fill nothing in Oklab, however far along.
	ok.stops = []paintStop{{lab: scene.Oklab(redOpaque), at: 0}, {lab: scene.Oklab(blueOpaque), at: 1}}
	if c := ok.stopsAt(0.5); c != [4]float32{} {
		t.Errorf("a mix of two transparent stops fills %v, want nothing", c)
	}

	// A painter takes at most MaxStops of the scene's, so a point past the
	// last kept stop holds that stop's color.
	many := make([]paintStop, scene.MaxStops)
	for i := range many {
		many[i] = paintStop{c: [4]float32{float32(i) / (scene.MaxStops - 1), 0, 0, 1}, at: float32(i) / (scene.MaxStops - 1)}
	}
	p = painter{paint: scene.PaintLinearStops, stops: many}
	if c := p.stopsAt(1); c != many[scene.MaxStops-1].c {
		t.Errorf("past the kept stops the fill is %v, want the last's", c)
	}
}
