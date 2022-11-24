package eval

import (
	"reflect"
	"testing"

	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
)

func TestEvalExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`const s = "\"hello\""; s`, `\"hello\"`},
		{`const a = [1,2,3]; (a[1] == 2) ? true : false`, true},
		{`const a = [1,2,3]; const r = (a[1] == 2) ? (1 + 1) : (10 % 3); r;`, 2},
		{`func f1() { true }; func f2() { true }; f1 == f2;`, false},
		{`func f() { true }; const f1 = f; const f2 = f; f1 == f2;`, true},
		{`const f1 = func() { true }; const f2 = func() { true }; f1 == f2;`, false},
		{`const f1 = func() { true }; const f2 = func() { false }; f1 == f2;`, false},
		{`const d1 = {"k1": "v1", "k2": "v2"}; const d2 = {"k1": "v1", "k2": "v2"}; d1 == d2`, true},
		{`const d1 = {"k1": "v1", "k2": "v2"}; const d2 = {"k1": "v1", "k2": "v"}; d1 == d2`, false},
		{`const a1 = [1,2,3]; const a2 = [1,2,3]; a1 == a2;`, true},
		{`const a1 = [1,2,3]; const a2 = [1,2]; a1 == a2;`, false},
		{`const arr = [1,2,3]; const f1 = arr.len; const f2 = arr.len; f1 == f2;`, true},
		{`const f1 = str; const f2 = str; f1 == f2;`, true},
		{`const f1 = str; const f2 = len; f1 == f2;`, false},
		{`const s1 = "hello"; const s2 = "hello"; s1 == s2;`, true},
		{`const s1 = "hello1"; const s2 = "hello2"; s1 == s2;`, false},

		{`const arr = [1,2,3]; const f = arr.len; f();`, 3},
		{`const arr = [1,2,3]; arr.len();`, 3},

		{`{"foo": 5}["foo"]`, 5},
		{`const k = "foo"; {"foo": 5}[k]`, 5},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},

		{`const arr = [1,2,4]; arr.first()`, 1},
		{`const arr = [1,2,4]; arr.last()`, 4},
		{`const arr = [1,2,4]; arr.tail()`, []int64{2, 4}},
		{`const arr = [1,2,4]; arr.tail().tail()`, []int64{4}},
		{`const arr = [1,2,4]; arr.push(8)`, []int64{1, 2, 4, 8}},

		{`"123"[1]`, "2"},
		{`const s = "123"; s[2]`, "3"},

		{`[1,2,4][0]`, 1},
		{`[1,2,4][1]`, 2},
		{`[1,2,4][2]`, 4},
		{`const i = 0; [1][i];`, 1},
		{`[1,2,4][1+1]`, 4},
		{`const arr = [1,2,4]; arr[2];`, 4},
		{`const arr = [1,2,4]; arr[0] + arr[1] + arr[2];`, 7},
		{`const arr = [1,2,4]; const i = arr[0];`, 1},

		{`len("")`, 0},
		{`len("four")`, 4},

		{`"hello world"`, "hello world"},
		{`"hello" + " " + "world"`, "hello world"},

		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!null", true},
		{"!!null", false},

		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"15 % 10", 5},

		{"1 && 2", 2},
		{"2 && 1", 1},
		{"0 && 2", 0},
		{"1 || 2", 1},
		{"2 || 1", 2},
		{"0 || 2", 2},

		{"true + 1", 2},
		{"false + 1", 1},
		{"true - 1", 0},
		{"false - 1", -1},
		{"true * 2", 2},
		{"false * 2", 0},
		{"true / 2", 0},
		{"false / 2", 0},
		{"true % 2", 1},
		{"false % 2", 0},

		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 >= 1", true},
		{"1 <= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},

		{"true > 1", false},
		{"true >= 1", true},
		{"true < 1", false},
		{"true <= 1", true},
		{"true == 1", true},
		{"false > 0", false},
		{"false >= 0", true},
		{"false < 0", false},
		{"false <= 0", true},
		{"false == 0", true},

		{"true && 2", 2},
		{"2 && true", 1},
		{"false && 2", 0},
		{"true || 2", 1},
		{"2 || true", 2},
		{"false || 2", 2},

		{"null > 0", false},
		{"null >= 0", false},
		{"null < 0", true},
		{"null <= 0", true},
		{"null != 0", true},
		{"null == 0", false},
		{"null && 0", object.Nil},
		{"null || 0", 0},

		{"null > false", false},
		{"null >= false", false},
		{"null < false", true},
		{"null <= false", true},
		{"null != false", true},
		{"null == false", false},
		{"null && false", object.Nil},
		{"null || false", false},

		{"0 > null", true},
		{"0 >= null", true},
		{"0 < null", false},
		{"0 <= null", false},
		{"0 != null", true},
		{"0 == null", false},
		{"0 && null", 0},
		{"1 && null", object.Nil},
		{"0 || null", object.Nil},
		{"1 || null", 1},

		{"false > null", true},
		{"false >= null", true},
		{"false < null", false},
		{"false <= null", false},
		{"false != null", true},
		{"false == null", false},
		{"false && null", false},
		{"true && null", object.Nil},
		{"false || null", object.Nil},
		{"true || null", true},

		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"true && true", true},
		{"true && false", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || false", false},
		{"true + true", 2},
		{"true - true", 0},
		{"true * true", 1},
		{"true / true", 1},
		{"true % true", 0},
		{"true > true", false},
		{"true >= true", true},
		{"true < true", false},
		{"true <= true", true},

		{"null == null", true},
		{"null != null", false},
		{"null > null", false},
		{"null >= null", true},
		{"null < null", false},
		{"null <= null", true},
		{"null && null", object.Nil},
		{"null || null", object.Nil},
	}
	for i, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatalf("i: %v, err: %v", i, err)
		}
		if !testEvalObject(t, evaluated, tt.expected) {
			t.Fatalf("i: %v", i)
		}
	}
}

func testEval(input string) (object.Object, error) {
	l := lexer.New(input)
	p, err := parser.New(l)
	if nil != err {
		return object.Nil, err
	}
	program, err := p.ParseProgram()
	if nil != err {
		return object.Nil, err
	}
	env := object.NewEnv()
	return EvalAst(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not integer, got=%v", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%v, want: %v", result.Value, expected)
		return false
	}
	return true
}

func testIntegerSliceObject(t *testing.T, obj object.Object, expected []int64) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array, got=%v", obj)
		return false
	}
	sz := len(result.Items)
	esz := len(expected)
	if sz != esz {
		t.Errorf("size mismatch, expected=%v, result=%v", esz, sz)
		return false
	}
	for i := 0; i < sz; i++ {
		testIntegerObject(t, result.Items[i], expected[i])
	}
	return true
}

func testStringSliceObject(t *testing.T, obj object.Object, expected []string) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array, got=%v", obj)
		return false
	}
	sz := len(result.Items)
	esz := len(expected)
	if sz != esz {
		t.Errorf("size mismatch, expected=%v, result=%v", esz, sz)
		return false
	}
	for i := 0; i < sz; i++ {
		testStringObject(t, result.Items[i], expected[i])
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not string, got=%v", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%v, want: %v", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("object is not null, got=%v", obj)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean, got=%v", obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%v, want: %v", result.Value, expected)
		return false
	}
	return true
}

func testEvalObject(t *testing.T, evaluated object.Object, expected interface{}) bool {
	switch et := expected.(type) {
	case bool:
		return testBooleanObject(t, evaluated, et)
	case *object.Null:
		return testNullObject(t, evaluated)
	case int:
		return testIntegerObject(t, evaluated, int64(et))
	case int64:
		return testIntegerObject(t, evaluated, et)
	case string:
		return testStringObject(t, evaluated, expected.(string))
	case []int64:
		return testIntegerSliceObject(t, evaluated, expected.([]int64))
	case []string:
		return testStringSliceObject(t, evaluated, expected.([]string))
	default:
		t.Errorf("unsupport type, %v", reflect.TypeOf(expected))
		return false
	}
}

func TestVarStmts(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"const a = 5; a;", 5},
		{"const a = 5 * 5; a;", 25},
		{"const a = 5; const b = a; b;", 5},
		{"const a = 5; const b = a; const c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"
	evaluated, err := testEval(input)
	if nil != err {
		t.Fatal(err)
	}
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not function, got %v", reflect.TypeOf(evaluated).String())
	}
	arguments := len(fn.Args)
	if arguments != 1 {
		t.Fatalf("function has wrong args, got %v", arguments)
	}
	r, err := fn.Call(object.Objects{object.NewInteger(2)})
	if nil != err {
		t.Fatal(err)
	} else {
		ir, ok := r.(*object.Integer)
		if !ok {
			t.Fatalf("result is not integer, got %v", reflect.TypeOf(r).String())
		}
		if ir.Value != 4 {
			t.Fatalf("wrong result, got %v", ir.Value)
		}
	}
}

func TestFunctionCases(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"const identity = func(x) { x; }; identity(5)", 5},
		{"const double = func(x) { x * 2; }; double(5)", 10},
		{"const add = func(x, y) { x + y; }; add(5, 5)", 10},
		{"const add = func(x, y) { x + y; }; add(5 + 5, add(5, 5))", 20},
		{"func(x) { x; }(5)", 5},
		{"func(x) { func(y) { x + y } }(5)(5)", 10},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestArrayCases(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"[1, 2 * 3, 3 * 4]", []int64{1, 6, 12}},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		obj, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("wrong type, got %T (%+v)", evaluated, evaluated)
		}
		sz := len(obj.Items)
		if sz != 3 {
			t.Fatalf("wrong size, got %d", sz)
		}
		for i := 0; i < sz; i++ {
			testEvalObject(t, obj.Items[i], tt.expected[i])
		}
	}
}

func TestHash(t *testing.T) {
	input := `
	const k2 = "k2";
	{
		"k1": 10 - 9,
		k2: 1 + 1,
		"k" + "3": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	};
	`

	evaluated, err := testEval(input)
	if nil != err {
		t.Fatal(err)
	}
	h, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("type is not hash, got=%v", reflect.TypeOf(evaluated).String())
	}
	h1, _ := object.NewString("k1").Hash()
	h2, _ := object.NewString("k2").Hash()
	h3, _ := object.NewString("k3").Hash()
	h4, _ := object.NewInteger(4).Hash()
	h5, _ := object.True.Hash()
	h6, _ := object.False.Hash()
	expected := map[*object.HashKey]int64{
		h1: 1,
		h2: 2,
		h3: 3,
		h4: 4,
		h5: 5,
		h6: 6,
	}
	if len(h.Pairs) != len(expected) {
		t.Fatalf("hash has wrong num of pairs, got=%v", len(h.Pairs))
	}

	for ek, ev := range expected {
		pair, ok := h.Pairs[*ek]
		if !ok {
			t.Fatalf("mismatch key, got=%v", *ek)
		}
		testIntegerObject(t, pair.Value, ev)
	}
}
