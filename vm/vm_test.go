package vm

import (
	"fmt"
	"testing"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
)

type testHashType = map[object.HashKey]int64

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

func testBooleanObject(want bool, obj object.Object) error {
	result, ok := obj.(*object.Boolean)
	if !ok {
		return function.NewError(fmt.Errorf("object is not boolean, got=%v", obj))
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

func testIntSliceObject(want []int, obj object.Object) error {
	result, ok := obj.(*object.Array)
	if !ok {
		return function.NewError(fmt.Errorf("object is not Array, got=%v", obj))
	}
	if len(result.Items) != len(want) {
		return function.NewError(fmt.Errorf("wrong size, got=%v, want: %v", len(result.Items), len(want)))
	}
	for i, wi := range want {
		if err := testIntegerObject(int64(wi), result.Items[i]); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func testIntMapObject(want testHashType, obj object.Object) error {
	result, ok := obj.(*object.Hash)
	if !ok {
		return function.NewError(fmt.Errorf("object is not Hash, got=%v", obj))
	}
	if len(result.Pairs) != len(want) {
		return function.NewError(fmt.Errorf("wrong size, got=%v, want: %v", len(result.Pairs), len(want)))
	}
	for k, v := range want {
		pair, ok := result.Pairs[k]
		if !ok {
			return function.NewError(fmt.Errorf("key missing, got=%v", k))
		}
		if err := testIntegerObject(v, pair.Value); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

type vmTestCase struct {
	name  string
	input string
	want  interface{}
}

func testExpectedObject(t *testing.T, want interface{}, v object.Object) {
	switch et := want.(type) {
	case int:
		if err := testIntegerObject(int64(et), v); nil != err {
			t.Errorf("testIntegerObject failed, err: %v", err)
		}
	case bool:
		if err := testBooleanObject(et, v); nil != err {
			t.Errorf("testBooleanObject failed, err: %v", err)
		}
	case string:
		if err := testStringObject(et, v); nil != err {
			t.Errorf("testStringObject failed, err: %v", err)
		}
	case []int:
		if err := testIntSliceObject(et, v); nil != err {
			t.Errorf("testIntSliceObject failed, err: %v", err)
		}
	case testHashType:
		if err := testIntMapObject(et, v); nil != err {
			t.Errorf("testIntMapObject failed, err: %v", err)
		}
	}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := parse(t, tt.input)
			c := compiler.New()
			err := c.Compile(program)
			if nil != err {
				t.Fatal(err)
			}
			vm := New(c.Bytecode(), c.Constants())
			if err := vm.Run(); nil != err {
				t.Fatal(err)
			}
			e := vm.LastPopped()
			testExpectedObject(t, tt.want, e)
		})
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", "1", 1},
		{"case_2", "2", 2},
		{"case_3", "1 + 2", 3},
		{"case_4", "1 - 2", -1},
		{"case_5", "1 * 2", 2},
		{"case_6", "4 / 2", 2},
		{"case_7", "50 / 2 * 2 + 10 - 5", 55},
		{"case_8", "5 + 5 + 5 + 5 - 10", 10},
		{"case_9", "2 * 2 * 2 * 2 * 2", 32},
		{"case_10", "5 * 2 + 10", 20},
		{"case_11", "5 + 2 * 10", 25},
		{"case_12", "5 * (2 + 10)", 60},
		{"case_13", "-5", -5},
		{"case_14", "-10", -10},
		{"case_15", "-50 + 100 + -50", 0},
		{"case_16", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVmTests(t, tests)
}

func TestBooleanExpr(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", "true", true},
		{"case_2", "false", false},
		{"case_3", "1 < 2", true},
		{"case_4", "1 > 2", false},
		{"case_5", "1 == 2", false},
		{"case_6", "1 >= 2", false},
		{"case_7", "1 <= 2", true},
		{"case_8", "1 != 2", true},
		{"case_9", "true == true", true},
		{"case_10", "true == false", false},
		{"case_11", "true != false", true},
		{"case_12", "true != true", false},
		{"case_13", "(1 < 2) == true", true},
		{"case_14", "(1 < 2) == false", false},
		{"case_15", "(1 > 2) == true", false},
		{"case_16", "(1 > 2) == false", true},
		{"case_17", "!true", false},
		{"case_18", "!false", true},
		{"case_19", "!5", false},
		{"case_20", "!!true", true},
		{"case_21", "!!false", false},
		{"case_22", "!!5", true},
	}
	runVmTests(t, tests)
}

func TestConditional(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", "true ? 10 : 20", 10},
		{"case_2", "false ? 10 : 20", 20},
		{"case_3", "1 ? 10 : 20", 10},
		{"case_4", "(1 < 2) ? 10 : 20", 10},
		{"case_5", "(1 > 2) ? 10 : 20", 20},
	}
	runVmTests(t, tests)
}

func TestGlobalConstStmt(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", "const one = 1; one;", 1},
		{"case_2", "const one = 1; const two = 2; one + two;", 3},
		{"case_3", "const one = 1; const two = one + one; one + two;", 3},
	}
	runVmTests(t, tests)
}

func TestString(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", `"hello"`, "hello"},
		{"case_2", `"hello" + " " + "world"`, "hello world"},
	}
	runVmTests(t, tests)
}

func TestArray(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", "[]", []int{}},
		{"case_2", "[1,2,3]", []int{1, 2, 3}},
		{"case_3", "[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}
	runVmTests(t, tests)
}

func TestHash(t *testing.T) {
	v1 := object.NewInteger(1)
	h1, _ := v1.Hash()
	v2 := object.NewInteger(2)
	h2, _ := v2.Hash()
	v6 := object.NewInteger(6)
	h6, _ := v6.Hash()
	tests := []vmTestCase{
		{"case_1", "{}", testHashType{}},
		{"case_2", "{1:2,2:3}", testHashType{*h1: 2, *h2: 3}},
		{"case_3", "{1 + 1 : 2 * 2, 3 + 3 : 4 * 4}", testHashType{*h2: 4, *h6: 16}},
	}
	runVmTests(t, tests)
}

func TestIndexExpr(t *testing.T) {
	tests := []vmTestCase{
		{"case_1", "[1,2,3][1]", 2},
		{"case_2", "[1,2,3][0 + 2]", 3},
		{"case_3", "[[1,1,1]][0][0]", 1},
		{"case_4", "{1:1, 2:2}[1]", 1},
		{"case_5", "{1:1, 2:2}[2]", 2},
	}
	runVmTests(t, tests)
}

func TestCall(t *testing.T) {
	tests := []vmTestCase{
		{
			"case_1",
			`
			const fn = func() { 5 + 10 };
			fn();
			`,
			15,
		},
		{
			"case_2",
			`
			const f1 = func() { 1 };
			const f2 = func() { 2 };
			f1() + f2();
			`,
			3,
		},
		{
			"case_3",
			`
			const a = func() { 1 };
			const b = func() { a() + 1 };
			const c = func() { b() + 1 };
			c();
			`,
			3,
		},
		{
			"case_4",
			`
			const f1 = func() { 1 };
			const f2 = func() { f1 };
			f2()();
			`,
			1,
		},
		{
			"case_5",
			`
			func add() { 5 + 10 };
			add();
			`,
			15,
		},
		{
			"case_6",
			`
			func add(x, y) { x + y };
			add(10, 5);
			`,
			15,
		},
	}
	runVmTests(t, tests)
}
