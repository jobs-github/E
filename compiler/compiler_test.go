package compiler

import (
	"fmt"
	"testing"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
)

type compilerTestCase struct {
	name             string
	input            string
	wantConstants    []interface{}
	wantInstrcutions []code.Instructions
}

func newCode(op code.Opcode, operands ...int) code.Instructions {
	r, err := code.Make(op, operands...)
	if nil != err {
		return code.Instructions{}
	}
	return r
}

func parse(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p, err := parser.New(l)
	if nil != err {
		t.Fatal(err)
	}
	r, err := p.ParseProgram()
	if nil != err {
		t.Fatal(err)
	}
	program, ok := r.(*ast.Program)
	if !ok {
		t.Fatal("parse error, type not Program")
	}
	return program
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parse(t, tt.input)
			cpl := New()
			err := cpl.Compile(program)
			if nil != err {
				t.Fatal(err)
			}
			b := cpl.Bytecode()
			err = testInstructions(tt.wantInstrcutions, b.Instructions())
			if nil != err {
				t.Fatal(err)
			}
			err = testConstants(tt.wantConstants, b.Constants())
			if nil != err {
				t.Fatal(err)
			}
		})
	}
}

func testInstructions(want []code.Instructions, got code.Instructions) error {
	r := joinInstructions(want)
	if len(got) != len(r) {
		return function.NewError(fmt.Errorf("wrong len, want: %v, got: %v", len(r), len(got)))
	}
	for i, ins := range r {
		if got[i] != ins {
			return function.NewError(fmt.Errorf("wrong byte at pos %v\nwant: %q\ngot: %q", i, r, got))
		}
	}
	return nil
}

func testConstants(want []interface{}, got object.Objects) error {
	if len(want) != len(got) {
		return function.NewError(fmt.Errorf("wrong len, want: %v, got: %v", len(want), len(got)))
	}
	for i, v := range want {
		switch wantVal := v.(type) {
		case int:
			err := testIntegerObject(int64(wantVal), got[i])
			if nil != err {
				return function.NewError(err)
			}
		case string:
			err := testStringObject(wantVal, got[i])
			if nil != err {
				return function.NewError(err)
			}
		}
	}
	return nil
}

func joinInstructions(s []code.Instructions) code.Instructions {
	r := code.Instructions{}
	for _, ins := range s {
		r = append(r, ins...)
	}
	return r
}

func testIntegerObject(want int64, obj object.Object) error {
	result, ok := obj.(*object.Integer)
	if !ok {
		return function.NewError(fmt.Errorf("object is not integer, got=%v", obj))
	}
	if result.Value != want {
		return function.NewError(fmt.Errorf("object has wrong value, got=%v, want: %v", result.Value, want))
	}
	return nil
}

func testStringObject(want string, obj object.Object) error {
	result, ok := obj.(*object.String)
	if !ok {
		return function.NewError(fmt.Errorf("object is not string, got=%v", obj))
	}
	if result.Value != want {
		return function.NewError(fmt.Errorf("object has wrong value, got=%v, want: %v", result.Value, want))
	}
	return nil
}

func Test_IntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			"case_1",
			"1 + 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpAdd),
				newCode(code.OpPop),
			},
		},
		{
			"case_2",
			"1;2;",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpPop),
				newCode(code.OpConst, 1),
				newCode(code.OpPop),
			},
		},
		{
			"case_3",
			"1 - 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpSub),
				newCode(code.OpPop),
			},
		},
		{
			"case_4",
			"1 * 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpMul),
				newCode(code.OpPop),
			},
		},
		{
			"case_5",
			"2 / 1",
			[]interface{}{2, 1},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpDiv),
				newCode(code.OpPop),
			},
		},
		{
			"case_6",
			"10 % 3",
			[]interface{}{10, 3},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpMod),
				newCode(code.OpPop),
			},
		},
		{
			"case_7",
			"-1",
			[]interface{}{1},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpNeg),
				newCode(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func Test_BooleanExpr(t *testing.T) {
	tests := []compilerTestCase{
		{
			"case_1",
			"true",
			[]interface{}{},
			[]code.Instructions{
				newCode(code.OpTrue),
				newCode(code.OpPop),
			},
		},
		{
			"case_2",
			"false",
			[]interface{}{},
			[]code.Instructions{
				newCode(code.OpFalse),
				newCode(code.OpPop),
			},
		},
		{
			"case_3",
			"1 > 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpGt),
				newCode(code.OpPop),
			},
		},
		{
			"case_4",
			"1 < 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpLt),
				newCode(code.OpPop),
			},
		},
		{
			"case_5",
			"1 == 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpEq),
				newCode(code.OpPop),
			},
		},
		{
			"case_6",
			"1 != 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpNeq),
				newCode(code.OpPop),
			},
		},
		{
			"case_6",
			"1 >= 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpGeq),
				newCode(code.OpPop),
			},
		},
		{
			"case_7",
			"1 <= 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpLeq),
				newCode(code.OpPop),
			},
		},
		{
			"case_8",
			"true == false",
			[]interface{}{},
			[]code.Instructions{
				newCode(code.OpTrue),
				newCode(code.OpFalse),
				newCode(code.OpEq),
				newCode(code.OpPop),
			},
		},
		{
			"case_9",
			"true != false",
			[]interface{}{},
			[]code.Instructions{
				newCode(code.OpTrue),
				newCode(code.OpFalse),
				newCode(code.OpNeq),
				newCode(code.OpPop),
			},
		},
		{
			"case_10",
			"1 && 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpAnd),
				newCode(code.OpPop),
			},
		},
		{
			"case_11",
			"1 || 2",
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpOr),
				newCode(code.OpPop),
			},
		},
		{
			"case_12",
			"!true",
			[]interface{}{},
			[]code.Instructions{
				newCode(code.OpTrue),
				newCode(code.OpNot),
				newCode(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func Test_Conditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			"case_1",
			"true ? 10 : 20;3333;",
			[]interface{}{10, 20, 3333},
			[]code.Instructions{
				// 0000
				newCode(code.OpTrue),
				// 0001
				newCode(code.OpJumpWhenFalse, 10),
				// 0004
				newCode(code.OpConst, 0),
				// 0007
				newCode(code.OpJump, 13),
				// 0010
				newCode(code.OpConst, 1),
				// 0013
				newCode(code.OpPop),
				// 0014
				newCode(code.OpConst, 2),
				// 0017
				newCode(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func Test_GlobalConstStmts(t *testing.T) {
	tests := []compilerTestCase{
		{
			"case_1",
			`
			const one = 1;
			const two = 2;
			`,
			[]interface{}{1, 2},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpSetGlobal, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpSetGlobal, 1),
			},
		},
		{
			"case_2",
			`
			const one = 1;
			one;
			`,
			[]interface{}{1},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpSetGlobal, 0),
				newCode(code.OpGetGlobal, 0),
				newCode(code.OpPop),
			},
		},
		{
			"case_3",
			`
			const one = 1;
			const two = one;
			two;
			`,
			[]interface{}{1},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpSetGlobal, 0),
				newCode(code.OpGetGlobal, 0),
				newCode(code.OpSetGlobal, 1),
				newCode(code.OpGetGlobal, 1),
				newCode(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func Test_Define(t *testing.T) {
	expected := map[string]*Symbol{
		"a": &Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": &Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	st := NewSymbolTable()
	a := st.define("a")
	if err := a.equal(expected["a"]); nil != err {
		t.Errorf("expected a=%+v, got=%+v, err: %v", expected["a"], a, err)
	}

	b := st.define("b")
	if err := b.equal(expected["b"]); nil != err {
		t.Errorf("expected b=%+v, got=%+v, err: %v", expected["b"], b, err)
	}
}

func Test_ResolveGlobal(t *testing.T) {
	expected := map[string]*Symbol{
		"a": &Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": &Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	st := NewSymbolTable()
	st.define("a")
	st.define("b")

	for _, sy := range expected {
		r, err := st.resolve(sy.Name)
		if nil != err {
			t.Errorf("name %v not resolvable", sy.Name)
			continue
		}
		if err := r.equal(sy); nil != err {
			t.Errorf("expected %v to resolve to %+v, got=%+v, err: %v", sy.Name, sy, r, err)
		}
	}
}

func Test_String(t *testing.T) {
	tests := []compilerTestCase{
		{
			"case_1",
			`"hello"`,
			[]interface{}{"hello"},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpPop),
			},
		},
		{
			"case_2",
			`"hello" + "world"`,
			[]interface{}{"hello", "world"},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpAdd),
				newCode(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func Test_Array(t *testing.T) {
	tests := []compilerTestCase{
		{
			"case_1",
			"[]",
			[]interface{}{},
			[]code.Instructions{
				newCode(code.OpArray, 0),
				newCode(code.OpPop),
			},
		},
		{
			"case_2",
			"[1,2,3]",
			[]interface{}{1, 2, 3},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpConst, 2),
				newCode(code.OpArray, 3),
				newCode(code.OpPop),
			},
		},
		{
			"case_3",
			"[1 + 2, 3 - 4, 5 * 6]",
			[]interface{}{1, 2, 3, 4, 5, 6},
			[]code.Instructions{
				newCode(code.OpConst, 0),
				newCode(code.OpConst, 1),
				newCode(code.OpAdd),
				newCode(code.OpConst, 2),
				newCode(code.OpConst, 3),
				newCode(code.OpSub),
				newCode(code.OpConst, 4),
				newCode(code.OpConst, 5),
				newCode(code.OpMul),
				newCode(code.OpArray, 3),
				newCode(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}
