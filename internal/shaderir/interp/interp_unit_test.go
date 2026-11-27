package interp

import (
	"math"
	"testing"

	"github.com/egoist/mygo/internal/shaderir"
)

func TestMemberOffset(t *testing.T) {
	s := &shaderir.Type{Main: shaderir.Struct, Sub: []shaderir.Type{
		{Main: shaderir.Float},
		{Main: shaderir.Vec2},
		{Main: shaderir.Int},
	}}
	for _, c := range []struct {
		index int
		want  int
	}{
		{0, 0},
		{1, 1}, // past the float
		{2, 3}, // past the float and the vec2
	} {
		if got := memberOffset(s, c.index); got != c.want {
			t.Errorf("memberOffset of member %d = %d, want %d", c.index, got, c.want)
		}
	}
}

func TestConstBits(t *testing.T) {
	c := &compiler{}
	o := c.constBits(shaderir.Type{Main: shaderir.Bool}, 1)
	if o.off != 0 || c.size != 1 {
		t.Fatalf("a bool constant took %d words at %d, want 1 at 0", c.size, o.off)
	}
	s := &state{r: make([]uint32, c.size)}
	eval(s, o)
	if s.r[0] != 1 {
		t.Errorf("the constant wrote %d, want 1", s.r[0])
	}
	// A vector takes as many words as its lanes, of which the constant's
	// run fills the first; the constructor fills the rest.
	c2 := &compiler{}
	o2 := c2.constBits(shaderir.Type{Main: shaderir.Vec2}, 7)
	if o2.off != 0 || c2.size != 2 {
		t.Fatalf("a vec2 constant took %d words, want 2", c2.size)
	}
	s2 := &state{r: make([]uint32, c2.size)}
	eval(s2, o2)
	if s2.r[0] != 7 || s2.r[1] != 0 {
		t.Errorf("the constant wrote %v, want 7 then nothing", s2.r)
	}
}

func TestDot(t *testing.T) {
	s := &state{r: []uint32{
		math.Float32bits(1.5), math.Float32bits(2.5), // a
		math.Float32bits(4), math.Float32bits(-2), // b
	}}
	a := operand{off: 0, t: shaderir.Type{Main: shaderir.Vec2}}
	b := operand{off: 2, t: shaderir.Type{Main: shaderir.Vec2}}
	if got := dot(s, a, b); got != 1.5*4+2.5*(-2) {
		t.Errorf("dot = %v, want %v", got, 1.5*4+2.5*(-2))
	}
}
