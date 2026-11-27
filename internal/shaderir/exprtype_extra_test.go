package shaderir_test

import (
	"go/constant"
	"reflect"
	"testing"

	s "github.com/egoist/mygo/internal/shaderir"
)

type typeEnv struct{}

func (typeEnv) LocalType(i int) s.Type {
	return []s.Type{{Main: s.Vec3}, {Main: s.Mat3}, {Main: s.Array, Sub: []s.Type{{Main: s.Uint}}}, {Main: s.Float},
		{Main: s.Struct, Sub: []s.Type{{Main: s.Float}, {Main: s.IVec4}}}}[i]
}
func (typeEnv) UniformType(int) s.Type { return s.Type{Main: s.IVec2} }
func (typeEnv) FuncReturn(int) s.Type  { return s.Type{Main: s.Vec4} }

func TestExpressionTypes(t *testing.T) {
	num := s.Expr{Type: s.NumberExpr, Const: constant.MakeInt64(1)}
	local := func(i int) s.Expr { return s.Expr{Type: s.LocalVariable, Index: i} }
	cases := []struct {
		name string
		expr s.Expr
		want s.BasicType
	}{
		{"integer", num, s.Int},
		{"unsigned", s.Expr{Type: s.NumberExpr, Const: constant.MakeInt64(1), Unsigned: true}, s.Uint},
		{"float", s.Expr{Type: s.NumberExpr, Const: constant.MakeFloat64(1.5)}, s.Float},
		{"bool", s.Expr{Type: s.NumberExpr, Const: constant.MakeBool(true)}, s.Bool},
		{"uniform", s.Expr{Type: s.UniformVariable}, s.IVec2},
		{"texture", s.Expr{Type: s.TextureVariable}, s.Texture},
		{"not", s.Expr{Type: s.Unary, Op: s.NotOp, Exprs: []s.Expr{num}}, s.Bool},
		{"unary", s.Expr{Type: s.Unary, Exprs: []s.Expr{local(0)}}, s.Vec3},
		{"comparison", s.Expr{Type: s.Binary, Op: s.LessThanOp, Exprs: []s.Expr{num, num}}, s.Bool},
		{"vector scalar", s.Expr{Type: s.Binary, Op: s.Add, Exprs: []s.Expr{num, local(0)}}, s.Vec3},
		{"matrix vector", s.Expr{Type: s.Binary, Op: s.MatrixMul, Exprs: []s.Expr{local(1), local(0)}}, s.Vec3},
		{"shift", s.Expr{Type: s.Binary, Op: s.LeftShift, Exprs: []s.Expr{local(0), num}}, s.Vec3},
		{"right shift", s.Expr{Type: s.Binary, Op: s.RightShift, Exprs: []s.Expr{local(3), num}}, s.Float},
		{"logic", s.Expr{Type: s.Binary, Op: s.OrOr, Exprs: []s.Expr{local(3), num}}, s.Bool},
		{"vector op vector", s.Expr{Type: s.Binary, Op: s.ComponentWiseMul, Exprs: []s.Expr{local(0), local(0)}}, s.Vec3},
		{"scalar op vector", s.Expr{Type: s.Binary, Op: s.Add, Exprs: []s.Expr{local(3), num}}, s.Float},
		{"scalar op matrix", s.Expr{Type: s.Binary, Op: s.Add, Exprs: []s.Expr{local(1), num}}, s.Mat3},
		{"constant op matrix", s.Expr{Type: s.Binary, Op: s.Add, Exprs: []s.Expr{num, local(1)}}, s.Mat3},
		{"constants", s.Expr{Type: s.Binary, Op: s.Add, Exprs: []s.Expr{num, num}}, s.Int},
		{"selection", s.Expr{Type: s.Selection, Exprs: []s.Expr{num, num, local(0)}}, s.Vec3},
		{"selection all constants", s.Expr{Type: s.Selection, Exprs: []s.Expr{num, num}}, s.Int},
		{"selection without operands", s.Expr{Type: s.Selection, Exprs: []s.Expr{num}}, s.None},
		{"swizzle", s.Expr{Type: s.FieldSelector, Exprs: []s.Expr{local(0), {Type: s.SwizzlingExpr, Swizzling: "xy"}}}, s.Vec2},
		{"struct member", s.Expr{Type: s.FieldSelector, Exprs: []s.Expr{local(4), {Type: s.StructMember, Index: 1}}}, s.IVec4},
		{"array index", s.Expr{Type: s.Index, Exprs: []s.Expr{local(2), num}}, s.Uint},
		{"matrix index", s.Expr{Type: s.Index, Exprs: []s.Expr{local(1), num}}, s.Vec3},
		{"vector index", s.Expr{Type: s.Index, Exprs: []s.Expr{local(0), num}}, s.Float},
		{"function", s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.FunctionExpr}}}, s.Vec4},
		{"constructor", s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.BuiltinFuncExpr, BuiltinFunc: s.UVec2F}}}, s.UVec2},
		{"step", s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.BuiltinFuncExpr, BuiltinFunc: s.Step}, num, local(0)}}, s.Vec3},
		{"smoothstep", s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.BuiltinFuncExpr, BuiltinFunc: s.Smoothstep}, num, num, local(0)}}, s.Vec3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := s.TypeOf(&c.expr, typeEnv{}); got.Main != c.want {
				t.Errorf("TypeOf = %v, want %v", got.Main, c.want)
			}
		})
	}
}

func TestBuiltinCallTypes(t *testing.T) {
	rows := []struct {
		fn   s.BuiltinFunc
		want s.BasicType
	}{
		{s.BoolF, s.Bool}, {s.FrontFacing, s.Bool},
		{s.IntF, s.Int}, {s.Len, s.Int}, {s.Cap, s.Int},
		{s.UintF, s.Uint},
		{s.FloatF, s.Float}, {s.Length, s.Float}, {s.Distance, s.Float}, {s.Dot, s.Float},
		{s.Vec2F, s.Vec2}, {s.Vec3F, s.Vec3}, {s.Cross, s.Vec3},
		{s.Vec4F, s.Vec4}, {s.TexelAt, s.Vec4},
		{s.IVec2F, s.IVec2}, {s.IVec3F, s.IVec3}, {s.IVec4F, s.IVec4},
		{s.UVec2F, s.UVec2}, {s.UVec3F, s.UVec3}, {s.UVec4F, s.UVec4},
		{s.Mat2F, s.Mat2}, {s.Mat3F, s.Mat3}, {s.Mat4F, s.Mat4},
	}
	for _, r := range rows {
		e := s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.BuiltinFuncExpr, BuiltinFunc: r.fn}}}
		if got := s.TypeOf(&e, typeEnv{}); got.Main != r.want {
			t.Errorf("TypeOf(%s) = %v, want %v", r.fn, got.Main, r.want)
		}
	}
	// Min and Max take the first argument's type; a name that is none of the
	// built-ins takes the operands' shared type.
	fl := s.Expr{Type: s.LocalVariable, Index: 3}
	for _, c := range []struct {
		name string
		expr s.Expr
		want s.BasicType
	}{
		{"min", s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.BuiltinFuncExpr, BuiltinFunc: s.Min}, fl, {Type: s.NumberExpr, Const: constant.MakeFloat64(1)}}}, s.Float},
		{"max", s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.BuiltinFuncExpr, BuiltinFunc: s.Max}, fl, {Type: s.NumberExpr, Const: constant.MakeFloat64(1)}}}, s.Float},
		{"unknown builtin", s.Expr{Type: s.Call, Exprs: []s.Expr{{Type: s.BuiltinFuncExpr, BuiltinFunc: "nosuch"}, fl}}, s.Float},
	} {
		if got := s.TypeOf(&c.expr, typeEnv{}); got.Main != c.want {
			t.Errorf("TypeOf(%s) = %v, want %v", c.name, got.Main, c.want)
		}
	}
}

func TestBinaryOperandTypes(t *testing.T) {
	cases := []struct {
		name              string
		op                s.Op
		left, right, want s.BasicType
		ok                bool
	}{
		{"vector plus scalar", s.Add, s.Vec3, s.Float, s.Vec3, true},
		{"mixed scalars", s.Add, s.Int, s.Uint, s.None, false},
		{"bool arithmetic", s.Add, s.Bool, s.Bool, s.None, false},
		{"matching vector comparison", s.VectorEqualOp, s.UVec2, s.UVec2, s.Bool, true},
		{"mismatched vector comparison", s.VectorEqualOp, s.UVec2, s.UVec3, s.None, false},
		{"matrix comparison", s.EqualOp, s.Mat2, s.Mat2, s.None, false},
		{"integer vector and scalar", s.And, s.IVec3, s.Int, s.IVec3, true},
		{"mixed integer kinds", s.And, s.IVec3, s.Uint, s.None, false},
		{"scalar modulo vector", s.ModOp, s.Int, s.IVec3, s.None, false},
		{"shift different integer kinds", s.LeftShift, s.IVec3, s.UVec3, s.IVec3, true},
		{"shift different vector sizes", s.LeftShift, s.IVec3, s.UVec2, s.None, false},
		{"matrix times vector", s.MatrixMul, s.Mat3, s.Vec3, s.Vec3, true},
		{"matrix times wrong vector", s.MatrixMul, s.Mat3, s.Vec2, s.None, false},
		{"divide by matrix", s.Div, s.Mat3, s.Mat3, s.None, false},
		{"vector and bool", s.AndAnd, s.Bool, s.Bool, s.Bool, true},
		{"bool and float", s.AndAnd, s.Bool, s.Float, s.None, false},
		{"int equals int", s.EqualOp, s.Int, s.Int, s.Bool, true},
		{"array equals array", s.EqualOp, s.Array, s.Array, s.None, false},
		{"texture equals texture", s.EqualOp, s.Texture, s.Texture, s.None, false},
		{"int less than uint", s.LessThanOp, s.Int, s.Uint, s.None, false},
		{"modulo int vector", s.ModOp, s.IVec3, s.Int, s.IVec3, true},
		{"or vector by scalar", s.Or, s.UVec2, s.Uint, s.UVec2, true},
		{"scalar or vector", s.AndNot, s.Uint, s.UVec2, s.UVec2, true},
		{"scalar modulo vector", s.ModOp, s.Uint, s.UVec2, s.None, false},
		{"complement kind", s.Xor, s.Int, s.Uint, s.None, false},
		{"int shifted by int", s.LeftShift, s.Int, s.Int, s.Int, true},
		{"scalar shifted by vector", s.RightShift, s.Int, s.IVec3, s.None, false},
		{"uint vectors shifted", s.RightShift, s.UVec2, s.UVec2, s.UVec2, true},
		{"shift by float", s.LeftShift, s.Int, s.Float, s.None, false},
		{"matrix times float", s.MatrixMul, s.Mat2, s.Float, s.Mat2, true},
		{"float times matrix", s.MatrixMul, s.Float, s.Mat4, s.Mat4, true},
		{"vec2 times mat2", s.MatrixMul, s.Vec2, s.Mat2, s.Vec2, true},
		{"mat2 divided by float", s.Div, s.Mat2, s.Float, s.Mat2, true},
		{"float vector by float scalar", s.Div, s.Vec4, s.Float, s.Vec4, true},
		{"float scalar by float vector", s.Div, s.Float, s.Vec4, s.Vec4, true},
		{"int scalar by int vector", s.Sub, s.Int, s.IVec2, s.IVec2, true},
		{"uint scalar by uint vector", s.ComponentWiseMul, s.Uint, s.UVec4, s.UVec4, true},
		{"texture arithmetic", s.Add, s.Texture, s.Texture, s.None, false},
		{"array arithmetic", s.Sub, s.Array, s.Array, s.None, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := s.TypeFromBinaryOp(c.op, s.Type{Main: c.left}, s.Type{Main: c.right}, nil, nil)
			if ok != c.ok || got.Main != c.want {
				t.Errorf("got (%v,%v), want (%v,%v)", got.Main, ok, c.want, c.ok)
			}
		})
	}
}

func TestResolveBinaryConstants(t *testing.T) {
	i := constant.MakeInt64(2)
	f := constant.MakeFloat64(2.5)
	b := constant.MakeBool(true)
	cases := []struct {
		name   string
		op     s.Op
		l, r   constant.Value
		lt, rt s.BasicType
		lk, rk constant.Kind
		ok     bool
	}{
		{"promote left", s.Add, i, f, s.None, s.None, constant.Float, constant.Float, true},
		{"promote right", s.Add, f, i, s.None, s.None, constant.Float, constant.Float, true},
		{"untyped against vector", s.Add, i, nil, s.None, s.Vec3, constant.Float, constant.Unknown, true},
		{"untyped against uint", s.Add, nil, i, s.Uint, s.None, constant.Unknown, constant.Int, true},
		{"invalid shift", s.LeftShift, f, i, s.None, s.None, constant.Unknown, constant.Unknown, false},
		{"invalid boolean", s.Add, b, i, s.None, s.None, constant.Unknown, constant.Unknown, false},
		{"same kind", s.Add, i, constant.MakeInt64(3), s.None, s.None, constant.Int, constant.Int, true},
		{"shift of integers", s.RightShift, i, constant.MakeInt64(1), s.None, s.None, constant.Int, constant.Int, true},
		{"untyped against matrix", s.Add, i, nil, s.None, s.Mat3, constant.Float, constant.Unknown, true},
		{"untyped bool against bool", s.Add, b, nil, s.None, s.Bool, constant.Bool, constant.Unknown, true},
		{"untyped bool against float", s.Add, b, nil, s.None, s.Float, constant.Unknown, constant.Unknown, false},
		{"float against untyped", s.Add, nil, f, s.Vec2, s.None, constant.Unknown, constant.Float, true},
		{"bool against untyped bool", s.Add, nil, b, s.Bool, s.None, constant.Unknown, constant.Bool, true},
		{"typed against typed", s.Add, f, f, s.Float, s.Float, constant.Float, constant.Float, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l, r, ok := s.ResolveUntypedConstsForBinaryOp(c.op, c.l, c.r, s.Type{Main: c.lt}, s.Type{Main: c.rt})
			kind := func(v constant.Value) constant.Kind {
				if v == nil {
					return constant.Unknown
				}
				return v.Kind()
			}
			if ok != c.ok || kind(l) != c.lk || kind(r) != c.rk {
				t.Errorf("got (%v,%v,%v), want kinds (%v,%v,%v)", l, r, ok, c.lk, c.rk, c.ok)
			}
		})
	}
}

func TestUntypedBinaryOperandTypes(t *testing.T) {
	i, b := constant.MakeInt64(1), constant.MakeBool(true)
	untyped := s.Type{}
	cases := []struct {
		name string
		op   s.Op
		l, r constant.Value
		want s.BasicType
		ok   bool
	}{
		{"add", s.Add, i, constant.MakeInt64(2), s.None, true},
		{"bool and", s.AndAnd, b, constant.MakeBool(false), s.Bool, true},
		{"bool or", s.OrOr, b, constant.MakeBool(false), s.Bool, true},
		{"modulo", s.ModOp, i, constant.MakeInt64(2), s.None, true},
		{"float modulo", s.ModOp, constant.MakeFloat64(1), constant.MakeFloat64(2), s.None, false},
		{"bitwise", s.And, i, constant.MakeInt64(2), s.None, true},
		{"shift", s.LeftShift, i, constant.MakeInt64(2), s.None, true},
		{"equality", s.EqualOp, i, constant.MakeInt64(2), s.Bool, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := s.TypeFromBinaryOp(c.op, untyped, untyped, c.l, c.r)
			if ok != c.ok || got.Main != c.want {
				t.Errorf("got (%v,%v), want (%v,%v)", got.Main, ok, c.want, c.ok)
			}
		})
	}
}

func TestFragmentLibrary(t *testing.T) {
	empty := &s.Program{}
	lib, entry, funcs := empty.FragmentLibrary()
	if lib != empty || entry != nil || funcs != nil {
		t.Fatal("empty fragment changed")
	}
	block := &s.Block{Stmts: []s.Stmt{{Type: s.ExprStmt, Exprs: []s.Expr{{Type: s.FunctionExpr, Index: 4}}}}}
	p := &s.Program{Varyings: []s.Type{{Main: s.Vec2}}, Funcs: []s.Func{{Index: 4, Block: &s.Block{}}, {Index: 8, Block: &s.Block{}}}, VertexFunc: s.VertexFunc{Block: &s.Block{}}, FragmentFunc: s.FragmentFunc{Block: block}}
	lib, entry, funcs = p.FragmentLibrary()
	if lib == p || lib.VertexFunc.Block != nil || lib.FragmentFunc.Block != nil || entry.Index != 9 || entry.Block != block || entry.Return.Main != s.Vec4 || !reflect.DeepEqual(entry.InParams, []s.Type{{Main: s.Vec4}, {Main: s.Vec2}}) {
		t.Fatalf("unexpected library or entry: %+v", entry)
	}
	if len(funcs) != 2 || funcs[0].Index != 4 || funcs[1] != entry {
		t.Fatalf("reachable functions: %+v", funcs)
	}
	if len(p.Funcs) != 2 || p.FragmentFunc.Block != block {
		t.Fatal("original program changed")
	}
}
