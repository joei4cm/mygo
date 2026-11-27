package ui

import (
	"math"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

// fx returns an effect with a shader that tints every pixel red, so the
// tests can tell an effect's op apart from a plain fill.
func fx() *scene.Effect {
	return &scene.Effect{Name: "test"}
}

func paramsOf(v [4]float32) [5][4]float32 {
	var ps [5][4]float32
	ps[0] = v
	return ps
}

func testPainter() (*Painter, *scene.Scene) {
	var s scene.Scene
	rt := &engine{host: &windowHost{onGPU: true}}
	return &Painter{rt: rt, s: &s, scale: 1, opacity: 1}, &s
}

func TestPainterEffect(t *testing.T) {
	p, s := testPainter()
	p.Effect(fx(), Rect{X: 10, Y: 20, W: 30, H: 40}, [4]float32{5, 5, 5, 5}, 3, paramsOf([4]float32{1, 2, 3, 4}))
	if len(s.Ops) != 1 || len(s.Effects) != 1 {
		t.Fatalf("Effect left %d ops and %d effects, want one of each", len(s.Ops), len(s.Effects))
	}
	op := s.Ops[0]
	if op.Kind != scene.OpEffect || op.Rect != (scene.Rect{X: 10, Y: 20, W: 30, H: 40}) {
		t.Errorf("the effect's op is %+v, want the rectangle as given", op)
	}
	if op.Start != 0 || op.Opacity != 1 {
		t.Errorf("the effect's op starts at %d at opacity %v, want 0 at 1", op.Start, op.Opacity)
	}
	if s.Effects[0].Effect == nil || s.Effects[0].Blur != 3 || s.Effects[0].Params[0] != [4]float32{1, 2, 3, 4} {
		t.Errorf("the effect holds %+v, want the shader, its blur and its param", s.Effects[0])
	}

	// No effect, or an empty rectangle, paints nothing.
	p, s = testPainter()
	p.Effect(nil, Rect{W: 30, H: 40}, [4]float32{}, 0, [5][4]float32{})
	p.Effect(fx(), Rect{}, [4]float32{}, 0, [5][4]float32{})
	p.Effect(fx(), Rect{W: 30, H: 0}, [4]float32{}, 0, [5][4]float32{})
	if len(s.Ops) != 0 || len(s.Effects) != 0 {
		t.Errorf("empty effects left %d ops and %d effects, want none", len(s.Ops), len(s.Effects))
	}
}

func TestPainterShaderEffect(t *testing.T) {
	p, s := testPainter()
	p.ShaderEffect(fx(), Rect{X: 1, Y: 2, W: 8, H: 4}, 2, ShaderInput{Uniforms: []uint32{9}})
	if len(s.Ops) != 1 || len(s.Effects) != 1 {
		t.Fatalf("ShaderEffect left %d ops and %d effects, want one of each", len(s.Ops), len(s.Effects))
	}
	if s.Ops[0].Kind != scene.OpEffect || s.Effects[0].Uniforms == nil {
		t.Error("ShaderEffect did not paint the effect with its uniforms")
	}

	// No effect or an empty rectangle paints nothing.
	p, s = testPainter()
	p.ShaderEffect(nil, Rect{W: 8, H: 4}, 2, ShaderInput{})
	p.ShaderEffect(fx(), Rect{}, 2, ShaderInput{})
	if len(s.Ops) != 0 || len(s.Effects) != 0 {
		t.Errorf("empty shader effects left %d ops and %d effects, want none", len(s.Ops), len(s.Effects))
	}
}

func TestPainterShaderInstances(t *testing.T) {
	p, s := testPainter()
	p.ShaderInstances(fx(), []Rect{{X: 0, Y: 0, W: 8, H: 8}, {X: 20, Y: 4, W: 6, H: 2}},
		[]float32{1, 2, 3, 4, 5, 6, 7, 8}, 2, 1, ShaderInput{})
	if len(s.Ops) != 1 || len(s.Effects) != 1 {
		t.Fatalf("ShaderInstances left %d ops and %d effects, want one of each", len(s.Ops), len(s.Effects))
	}
	op := s.Ops[0]
	// The op's rectangle is the quads' bounds; each quad's words hold its
	// box, its size in DIPs, and its own two floats of data.
	if op.Rect != (scene.Rect{X: 0, Y: 0, W: 26, H: 8}) {
		t.Errorf("the instances' bounds are %v, want the quads' union", op.Rect)
	}
	if len(s.Words) != 2*(scene.InstanceHead+2) {
		t.Fatalf("the instances took %d words, want %d", len(s.Words), 2*(scene.InstanceHead+2))
	}
	q1x := math.Float32frombits(s.Words[0])
	if q1x != 0 {
		t.Errorf("the first quad's x is %v, want 0", q1x)
	}
	// The second quad's words start after the first's: its box, its size,
	// then its two floats of data, which were 3 and 4.
	d2 := math.Float32frombits(s.Words[scene.InstanceHead+2+2+4])
	if d2 != 3 {
		t.Errorf("the second quad's first datum is %v, want 3", d2)
	}
	if d2b := math.Float32frombits(s.Words[scene.InstanceHead+2+2+5]); d2b != 4 {
		t.Errorf("the second quad's second datum is %v, want 4", d2b)
	}

	// Empty quads hold their words but take no room; a run of none paints nothing.
	p, s = testPainter()
	p.ShaderInstances(fx(), []Rect{{W: 4, H: 4}, {}}, []float32{0, 0, 0, 0}, 2, 1, ShaderInput{})
	if len(s.Ops) != 1 {
		t.Fatalf("an empty quad dropped the op: %d ops", len(s.Ops))
	}
	if got := (scene.Rect{X: math.Float32frombits(s.Words[4]), Y: math.Float32frombits(s.Words[5]),
		W: math.Float32frombits(s.Words[6]), H: math.Float32frombits(s.Words[7])}); !got.Empty() {
		t.Errorf("the empty quad holds %v, want nothing", got)
	}

	p, s = testPainter()
	p.ShaderInstances(fx(), []Rect{{}}, []float32{0}, 0, 1, ShaderInput{})
	if len(s.Ops) != 0 || len(s.Effects) != 0 || len(s.Words) != 0 {
		t.Errorf("quads of none left %d ops, %d effects and %d words, want none", len(s.Ops), len(s.Effects), len(s.Words))
	}

	// Bad inputs paint nothing: no effect, no quads, a stride past the
	// limit, or less data than the quads' strides hold.
	p, s = testPainter()
	p.ShaderInstances(nil, []Rect{{W: 4, H: 4}}, []float32{0}, 0, 1, ShaderInput{})
	p.ShaderInstances(fx(), nil, nil, 0, 1, ShaderInput{})
	p.ShaderInstances(fx(), []Rect{{W: 4, H: 4}}, []float32{0}, scene.MaxInstanceData+1, 1, ShaderInput{})
	p.ShaderInstances(fx(), []Rect{{W: 4, H: 4}}, []float32{}, 1, 1, ShaderInput{})
	if len(s.Ops) != 0 || len(s.Effects) != 0 || len(s.Words) != 0 {
		t.Errorf("bad instances left %d ops, %d effects and %d words, want none", len(s.Ops), len(s.Effects), len(s.Words))
	}
}

func TestPainterEffectWords(t *testing.T) {
	p, s := testPainter()
	w := p.EffectWords(4)
	if len(w) != 4 || len(s.Words) != 4 {
		t.Fatalf("EffectWords(4) gave %d words of %d, want 4 of 4", len(w), len(s.Words))
	}
	for _, v := range w {
		if v != 0 {
			t.Errorf("a fresh word holds %d, want 0", v)
		}
	}
	copy(w, []uint32{1, 2, 3, 4})
	if s.Words[2] != 3 {
		t.Error("the words are not the scene's own")
	}
	// A second call continues where the first ended, zeroed.
	w2 := p.EffectWords(2)
	if len(w2) != 2 || len(s.Words) != 6 || w2[0] != 0 || s.Words[4] != 0 || s.Words[3] != 4 {
		t.Errorf("the next words are %v of %v, want two zeroed after the first four", w2, s.Words)
	}
}

func TestShaderInputOp(t *testing.T) {
	bm := NewBitmap(stripedImage(2, 2, 1, 0))
	in := ShaderInput{Uniforms: []uint32{7}, WideUniforms: []uint32{8}, Images: [4]*Bitmap{bm}, Blur: 2}
	op := in.op(fx())
	if op.Effect == nil || op.Uniforms == nil || op.WideUniforms == nil || op.Blur != 2 {
		t.Errorf("the op holds %+v, want the effect, its uniforms and its blur", op)
	}
	if len(op.Images) != 4 || op.Images[0] == nil || op.Images[1] != nil {
		t.Error("the op does not keep the input's image first and nil for the rest")
	}
}

func TestMaterialOfNode(t *testing.T) {
	// A plain material paints as it is; a builder sees the element.
	e := &node{}
	e.Material(fakeMaterial{})
	if _, ok := e.material.(fakeMaterial); !ok || e.fill != fillMaterial {
		t.Errorf("Material left fill %v with material %v, want fillMaterial of the material", e.fill, e.material)
	}
	b := &fakeBuilder{built: nil}
	e = &node{}
	e.Material(b)
	if b.built != e {
		t.Error("Material did not build a builder's material from the element")
	}
}

type fakeMaterial struct{}

func (fakeMaterial) PaintMaterial(p *Painter, box Rect, radii [4]float32) {
	p.Fill(box, RGB(255, 0, 0), 0)
}

type fakeBuilder struct{ built *node }

func (b *fakeBuilder) BuildMaterial(e *node) Material {
	b.built = e
	return fakeMaterial{}
}

func (fakeBuilder) PaintMaterial(p *Painter, box Rect, radii [4]float32) {
	p.Fill(box, RGB(0, 0, 255), 0)
}
