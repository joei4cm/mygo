package raster

import (
	"image"
	"math"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

// quads lays out quads of video, InstanceWords words each: the quad's
// rectangle, its size, and four words of data.
func quads(rects ...scene.Rect) (int32, []uint32) {
	var in []uint32
	for _, r := range rects {
		in = append(in,
			math.Float32bits(r.X), math.Float32bits(r.Y), math.Float32bits(r.W), math.Float32bits(r.H),
			math.Float32bits(r.W), math.Float32bits(r.H), 1, 2, 3, 4)
	}
	return 10, in
}

func TestSplitInstances(t *testing.T) {
	words, in := quads(scene.Rect{X: 10, Y: 10, W: 30, H: 20}, scene.Rect{X: 60, Y: 10, W: 30, H: 20})
	s := &scene.Scene{
		Ops: []scene.Op{
			{Kind: scene.OpFill, Rect: scene.Rect{X: 0, Y: 0, W: 10, H: 10}},
			{Kind: scene.OpEffect, Rect: scene.Rect{X: 0, Y: 0, W: 100, H: 100}, Start: 0, Transform: scene.Translation(10, 20)},
			{Kind: scene.OpFill, Rect: scene.Rect{X: 0, Y: 0, W: 5, H: 5}},
		},
		Effects: []scene.EffectOp{{InstanceWords: words, Instances: in}},
	}
	bounds := []image.Rectangle{
		image.Rect(0, 0, 10, 10),
		image.Rect(0, 0, 60, 60), // the second quad falls past this
		image.Rect(0, 0, 5, 5),
	}

	// Before an app paints quads, nothing splits: the scene is as it is.
	var inr instancer
	out, ob := inr.Split(s, bounds)
	if len(out.Ops) != len(s.Ops) || out != s {
		// Split returns s itself, not a copy, while the code is off.
		t.Errorf("Split before EnableInstances changed the scene: %d ops", len(out.Ops))
	}
	if len(ob) != len(bounds) {
		t.Errorf("Split before EnableInstances changed the bounds: %d", len(ob))
	}

	EnableInstances()
	EnableInstances() // twice is fine
	out, ob = inr.Split(s, bounds)

	// The two fills keep, and the effect takes its quads that the op's
	// bounds hold: the first quad of three, at 10,10 after the op's
	// translation of 10,20; the second at 60,10 → 70,30 falls past the
	// op's bounds of 20×20 and is left out.
	if len(out.Ops) != 3 || out.Ops[0].Kind != scene.OpFill || out.Ops[1].Kind != scene.OpEffect || out.Ops[2].Kind != scene.OpFill {
		t.Fatalf("the split scene holds %d ops: %+v", len(out.Ops), out.Ops)
	}
	if out.Ops[1].Rect != (scene.Rect{X: 10, Y: 10, W: 30, H: 20}) {
		t.Errorf("the quad's op draws %+v, want the quad", out.Ops[1].Rect)
	}
	// The quad's op points at the copy of the effect that holds it alone,
	// which the scene carries after the frame's own effects.
	if out.Ops[1].Start != 1 || out.Effects[1].Instance != 0 || len(out.Effects[1].Instances) != 10 {
		t.Errorf("the first quad's effect = index %d, instance %d, %d words", out.Ops[1].Start, out.Effects[1].Instance, len(out.Effects[1].Instances))
	}
	if got := out.Effects[1].Instances[0:4]; got[0] != math.Float32bits(float32(10)) || got[3] != math.Float32bits(float32(20)) {
		t.Error("the first quad's copy does not hold its own rectangle")
	}
	// The bounds of a quad's op are its bounds under the op's map, a few
	// pixels to spare, clipped to the op's.
	if ob[1] != (image.Rect(20-4, 30-4, 50+4, 50+4)) {
		t.Errorf("the quad's bounds = %v, want the mapped quad with room to spare", ob[1])
	}
	if ob[2] != bounds[2] || ob[0] != bounds[0] {
		t.Error("the fills' bounds moved")
	}
	if out.Effects[0].InstanceCount() != 2 {
		t.Error("the frame's own effect lost its quads")
	}

	// A quad that is empty takes nothing.
	words, in = quads(scene.Rect{X: 10, Y: 10, W: 30, H: 20}, scene.Rect{X: 60, Y: 10, W: 0, H: 20})
	s2 := &scene.Scene{
		Ops:     []scene.Op{{Kind: scene.OpEffect, Rect: scene.Rect{X: 0, Y: 0, W: 100, H: 100}, Start: 0}},
		Effects: []scene.EffectOp{{InstanceWords: words, Instances: in}},
	}
	all := []image.Rectangle{image.Rect(0, 0, 200, 200)}
	var inr2 instancer
	out2, _ := inr2.Split(s2, all)
	if len(out2.Ops) != 1 || out2.Ops[0].Rect != (scene.Rect{X: 10, Y: 10, W: 30, H: 20}) {
		t.Errorf("an empty quad took an op: %d ops, first at %+v", len(out2.Ops), out2.Ops[0].Rect)
	}

	// A scene without quads comes back as it is.
	s3 := &scene.Scene{Ops: []scene.Op{{Kind: scene.OpFill, Rect: scene.Rect{X: 0, Y: 0, W: 1, H: 1}}}}
	out3, ob3 := inr2.Split(s3, []image.Rectangle{image.Rect(0, 0, 1, 1)})
	if out3 != s3 || len(ob3) != 1 {
		t.Error("a scene without quads was rebuilt")
	}
}

func TestInstanced(t *testing.T) {
	words, in := quads(scene.Rect{X: 0, Y: 0, W: 1, H: 1})
	fx := scene.EffectOp{InstanceWords: words, Instances: in}
	s := &scene.Scene{Ops: []scene.Op{{Kind: scene.OpEffect, Start: 0}}, Effects: []scene.EffectOp{fx}}
	if !instanced(s, &s.Ops[0]) {
		t.Error("an effect with quads is not instanced")
	}
	for _, op := range []scene.Op{
		{Kind: scene.OpFill, Start: 0},
		{Kind: scene.OpEffect, Start: -1},
		{Kind: scene.OpEffect, Start: 1},
		{Kind: scene.OpEffect, Start: 0},
	} {
		if instanced(s, &op) && op.Kind != scene.OpEffect {
			t.Errorf("a %v op is instanced", op.Kind)
		}
	}
	// An effect of no quads is not instanced, wherever it points.
	empty := scene.EffectOp{InstanceWords: 5, Instances: make([]uint32, 20)}
	s2 := &scene.Scene{Ops: []scene.Op{{Kind: scene.OpEffect, Start: 0}}, Effects: []scene.EffectOp{empty}}
	if instanced(s2, &s2.Ops[0]) {
		t.Error("an effect with no quads is instanced")
	}
}

func TestQuadBounds(t *testing.T) {
	quad := scene.Rect{X: 10.5, Y: 20.5, W: 30, H: 40}
	if got := quadBounds(scene.Affine{}, quad); got != (image.Rect(10-4, 20-4, 41+4, 61+4)) {
		t.Errorf("quadBounds of an identity = %v, want the quad with room to spare", got)
	}
	if got := quadBounds(scene.Translation(10, 20), quad); got != (image.Rect(20-4, 40-4, 51+4, 81+4)) {
		t.Errorf("quadBounds under a translation = %v, want the mapped quad with room to spare", got)
	}
}
