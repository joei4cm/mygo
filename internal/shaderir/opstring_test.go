package shaderir_test

import (
	"go/token"
	"testing"

	shaderir "github.com/egoist/mygo/internal/shaderir"
)

func TestTypeString(t *testing.T) {
	cases := []struct {
		name shaderir.Type
		want string
	}{
		{shaderir.Type{}, "none"},
		{shaderir.Type{Main: shaderir.Bool}, "bool"},
		{shaderir.Type{Main: shaderir.Int}, "int"},
		{shaderir.Type{Main: shaderir.Float}, "float"},
		{shaderir.Type{Main: shaderir.Uint}, "uint"},
		{shaderir.Type{Main: shaderir.Vec4}, "vec4"},
		{shaderir.Type{Main: shaderir.IVec2}, "ivec2"},
		{shaderir.Type{Main: shaderir.UVec3}, "uvec3"},
		{shaderir.Type{Main: shaderir.Mat4}, "mat4"},
		{shaderir.Type{Main: shaderir.Texture}, "texture"},
		{shaderir.Type{Main: shaderir.Array, Length: 3, Sub: []shaderir.Type{{Main: shaderir.Float}}}, "[3]float"},
		{shaderir.Type{Main: shaderir.Struct, Sub: []shaderir.Type{{Main: shaderir.Int}, {Main: shaderir.Vec2}}}, "struct{int,vec2}"},
		{shaderir.Type{Main: shaderir.BasicType(99)}, "?(unknown type: 99)"},
	}
	for _, c := range cases {
		if got := c.name.String(); got != c.want {
			t.Errorf("String of %+v = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestOpFromToken(t *testing.T) {
	for _, c := range []struct {
		tok  token.Token
		want shaderir.Op
	}{
		{token.ADD, shaderir.Add},
		{token.SUB, shaderir.Sub},
		{token.NOT, shaderir.NotOp},
		{token.QUO, shaderir.Div},
		{token.QUO_ASSIGN, shaderir.Div},
		{token.REM, shaderir.ModOp},
		{token.SHL, shaderir.LeftShift},
		{token.SHR, shaderir.RightShift},
		{token.LSS, shaderir.LessThanOp},
		{token.LEQ, shaderir.LessThanEqualOp},
		{token.GTR, shaderir.GreaterThanOp},
		{token.GEQ, shaderir.GreaterThanEqualOp},
		{token.EQL, shaderir.EqualOp},
		{token.NEQ, shaderir.NotEqualOp},
		{token.AND, shaderir.And},
		{token.AND_NOT, shaderir.AndNot},
		{token.XOR, shaderir.Xor},
		{token.OR, shaderir.Or},
		{token.LAND, shaderir.AndAnd},
		{token.LOR, shaderir.OrOr},
	} {
		if op, ok := shaderir.OpFromToken(c.tok, shaderir.Type{}, shaderir.Type{}); !ok || op != c.want {
			t.Errorf("OpFromToken(%v) = (%v, %v), want (%v, true)", c.tok, op, ok, c.want)
		}
	}
	// A multiplication of matrices takes the matrix op, of anything else
	// the one that goes by parts; and vectors compare by parts too.
	mat := shaderir.Type{Main: shaderir.Mat2}
	if op, ok := shaderir.OpFromToken(token.MUL, mat, shaderir.Type{}); !ok || op != shaderir.MatrixMul {
		t.Errorf("OpFromToken(MUL, mat2, none) = (%v, %v), want the matrix multiply", op, ok)
	}
	if op, ok := shaderir.OpFromToken(token.MUL, shaderir.Type{}, mat); !ok || op != shaderir.MatrixMul {
		t.Errorf("OpFromToken(MUL, none, mat2) = (%v, %v), want the matrix multiply", op, ok)
	}
	if op, ok := shaderir.OpFromToken(token.MUL, shaderir.Type{Main: shaderir.Vec2}, shaderir.Type{}); !ok || op != shaderir.ComponentWiseMul {
		t.Errorf("OpFromToken(MUL, vec2, none) = (%v, %v), want the product by parts", op, ok)
	}
	if op, ok := shaderir.OpFromToken(token.EQL, shaderir.Type{Main: shaderir.Vec2}, shaderir.Type{}); !ok || op != shaderir.VectorEqualOp {
		t.Errorf("OpFromToken(EQL, vec2, none) = (%v, %v), want the vector equality", op, ok)
	}
	if op, ok := shaderir.OpFromToken(token.NEQ, shaderir.Type{}, shaderir.Type{Main: shaderir.Int}); !ok || op != shaderir.NotEqualOp {
		t.Errorf("OpFromToken(NEQ, none, int) = (%v, %v), want the scalar inequality", op, ok)
	}
	// A token that is none of the operations says so.
	if op, ok := shaderir.OpFromToken(token.ELLIPSIS, shaderir.Type{}, shaderir.Type{}); ok {
		t.Errorf("OpFromToken(…) = (%v, %v), want none", op, ok)
	}
}
