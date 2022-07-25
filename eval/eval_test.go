package eval

import (
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
)

func scriptsDir() string {
	_, filename, _, _ := runtime.Caller(1)
	cur := path.Dir(filename)
	return filepath.Join(filepath.Dir(cur), "scripts")
}

func TestEvalExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`var s = "\"hello\""; s`, `\"hello\"`},
		{`var a = [1,2,3]; (a[1] == 2) ? true : false`, true},
		{`var a = [1,2,3]; var r = (a[1] == 2) ? (1 + 1) : (10 % 3); r;`, 2},
		{`func f1() { return true; }; func f2() { return true; }; f1 == f2;`, false},
		{`func f() { return true; }; var f1 = f; var f2 = f; f1 == f2;`, true},
		{`var f1 = func() { return true; }; var f2 = func() { return true; }; f1 == f2;`, true},
		{`var f1 = func() { return true; }; var f2 = func() { return false; }; f1 == f2;`, false},
		{`var d1 = {"k1": "v1", "k2": "v2"}; var d2 = {"k1": "v1", "k2": "v2"}; d1 == d2`, true},
		{`var d1 = {"k1": "v1", "k2": "v2"}; var d2 = {"k1": "v1", "k2": "v"}; d1 == d2`, false},
		{`var a1 = [1,2,3]; var a2 = [1,2,3]; a1 == a2;`, true},
		{`var a1 = [1,2,3]; var a2 = [1,2]; a1 == a2;`, false},
		{`var arr = [1,2,3]; var f1 = arr.len; var f2 = arr.len; f1 == f2;`, true},
		{`var a = [1,2,3]; var i1 = a.iter(); var i2 = a.iter(); i1 == i2;`, true},
		{`var a = [1,2,3]; var i1 = a.iter(); var i2 = a.iter(); i1.next(); i1 == i2;`, false},
		{`var d = {"k1": "v1", "k2": "v2"}; var i1 = d.iter(); var i2 = d.iter(); i1 == i2;`, true},
		{`var d = {"k1": "v1", "k2": "v2"}; var i1 = d.iter(); var i2 = d.iter(); i1.next(); i1 == i2;`, false},
		{`var f1 = str; var f2 = str; f1 == f2;`, true},
		{`var f1 = str; var f2 = len; f1 == f2;`, false},
		{`var s1 = "hello"; var s2 = "hello"; s1 == s2;`, true},
		{`var s1 = "hello1"; var s2 = "hello2"; s1 == s2;`, false},

		{`var arr = [1,2,3]; var f = arr.len; f();`, 3},
		{`var arr = [1,2,3]; arr.len();`, 3},
		{`var arr = [1,2,3,4,5]; arr[1] = 20; arr[1]`, 20},
		{`var arr = [1,2,3,4,5]; arr[5-4] = 20; arr[2-1]`, 20},
		{`var m = {"k1": "v1", "k2": "v2"}; m["k1"] = "hello"; m["k1"]`, "hello"},
		{`var m = {"k1": "v1", "k2": "v2"}; var k = "k1"; m[k] = "hello"; m[k]`, "hello"},

		{`{"foo": 5}["foo"]`, 5},
		{`var k = "foo"; {"foo": 5}[k]`, 5},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},

		{`var arr = [1,2,4]; arr.first()`, 1},
		{`var arr = [1,2,4]; arr.last()`, 4},
		{`var arr = [1,2,4]; arr.tail()`, []int64{2, 4}},
		{`var arr = [1,2,4]; arr.tail().tail()`, []int64{4}},
		{`var arr = [1,2,4]; arr.push(8)`, []int64{1, 2, 4, 8}},

		{`"123"[1]`, "2"},
		{`var s = "123"; s[2]`, "3"},

		{`[1,2,4][0]`, 1},
		{`[1,2,4][1]`, 2},
		{`[1,2,4][2]`, 4},
		{`var i = 0; [1][i];`, 1},
		{`[1,2,4][1+1]`, 4},
		{`var arr = [1,2,4]; arr[2];`, 4},
		{`var arr = [1,2,4]; arr[0] + arr[1] + arr[2];`, 7},
		{`var arr = [1,2,4]; var i = arr[0];`, 1},

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
		testEvalObject(t, evaluated, tt.expected)
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
	return program.Eval(env, false)
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

func testEvalObject(t *testing.T, evaluated object.Object, expected interface{}) {
	switch et := expected.(type) {
	case bool:
		testBooleanObject(t, evaluated, et)
	case object.Null:
		testNullObject(t, evaluated)
	case int:
		testIntegerObject(t, evaluated, int64(et))
	case string:
		testStringObject(t, evaluated, expected.(string))
	case []int64:
		testIntegerSliceObject(t, evaluated, expected.([]int64))
	case []string:
		testStringSliceObject(t, evaluated, expected.([]string))
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (2 < 1) { 10 } else if (2 < 2) { 20 } else { 30 }", 30},
		{"if (2 < 1) { 10 } else if (2 > 1) { 20 } else { 30 }", 20},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestReturnStmts(t *testing.T) {
	stmt := `
	if (10 > 1) {
		if (10 > 1) {
			return 10;
		}
		return 1;
	}
	`
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9", 10},
		{"9; return 2 * 5; 9", 10},
		{stmt, 10},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestVarStmts(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
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
	argument := fn.Fn.ArgumentOf(0)
	if argument != "x" {
		t.Fatalf("argument of 0 not x, got `%v`", argument)
	}
	body := fn.Fn.Body()
	expected := "(x + 2)"
	if body != expected {
		t.Fatalf("body not (x + 2), got `%v`", body)
	}
}

func TestFunctionCases(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"var identity = func(x) { x; }; identity(5)", 5},
		{"var identity = func(x) { return x; }; identity(5)", 5},
		{"var double = func(x) { x * 2; }; double(5)", 10},
		{"var add = func(x, y) { x + y; }; add(5, 5)", 10},
		{"var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5))", 20},
		{"func(x) { x; }(5)", 5},
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

func TestAssignStmts(t *testing.T) {
	stmts := `
	var a = 5; 
	var b = a; 
	var c = a + b + 5; 
	a = c; 
	a;
	`
	tests := []struct {
		input    string
		expected interface{}
	}{
		{stmts, 15},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestForCases(t *testing.T) {
	stmt1 := `
	var i = 0;
	var result = 0;
	for ;; {
		if (i > 4) {
			break;
		}
		result = result + i;
		i = i + 1;
	}
	result;
	`

	stmt2 := `
	var result = 0;
	for var i = 0; i < 5; i = i + 1 {
		result = result + i;
	}
	result;
	`

	stmt3 := `
	var result = 0;
	for var i = 0; i < 5;  {
		result = result + i;
		i = i + 1;
	}
	result;
	`

	stmt4 := `
	var i = 0;
	var result = 0;
	for ; i < 5;  {
		result = result + i;
		i = i + 1;
	}
	result;
	`

	stmt5 := `
	var i = 0;
	var result = 0;
	for i = 0; i < 5; i = i + 1 {
		result = result + i;
	}
	result;
	`

	stmt6 := `
	var i = 0;
	var result = 0;
	for {
		if (i > 4) {
			break;
		}
		result = result + i;
		i = i + 1;
	}
	result;
	`

	tests := []struct {
		input    string
		expected interface{}
	}{
		{stmt1, 10},
		{stmt2, 10},
		{stmt3, 10},
		{stmt4, 10},
		{stmt5, 10},
		{stmt6, 10},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestIterCases(t *testing.T) {
	stmt1 := `
	var arr = [10,20,30];
	var result = [];
    for var iter = arr.iter(); iter != null; iter = iter.next() {
        var idx = iter.key();
        var val = iter.value();
		result.push(idx);
		result.push(val);
    };
	result;
	`

	stmt2 := `
	var h = {"k1": "v1", "k2": "v2", "k3": "v3"};
	var result = [];
    for var iter = h.iter(); iter != null; iter = iter.next() {
        var key = iter.key();
        var val = iter.value();
		result.push(key);
		result.push(val);
    };
	result;
	`

	tests := []struct {
		input    string
		expected interface{}
	}{
		{stmt1, []int64{0, 10, 1, 20, 2, 30}},
		{stmt2, []string{"k1", "v1", "k2", "v2", "k3", "v3"}},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestScopeCases(t *testing.T) {
	stmt1 := `
	var i = 10;
	var j = 0;
	if (true) {
		var i = 1;
		var j = 0;
		if (i == 1) {
			j = 1;
		}
	}
	j;
	`

	stmt2 := `
	var i = 10;
	var j = 0;
	if (true) {
		var i = 1;
		if (i == 1) {
			j = 1;
		}
	}
	j;
	`

	stmt3 := `
	var i = 1;
	for var i = 0; i < 5; i = i + 1 {
		i;
	}
	i;
	`

	stmt4 := `
	var i = 1;
	var result = 0;
	for i = 0; i < 5; i = i + 1 {
		result = result + i;
	}
	i;
	`

	tests := []struct {
		input    string
		expected interface{}
	}{
		{stmt1, 0},
		{stmt2, 1},
		{stmt3, 1},
		{stmt4, 5},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestComplexFunc(t *testing.T) {
	stmt1 := `
	var map = func(arr, fn) {
		var iter = func(arr, accumulated) {
			if (len(arr) == 0) {
				return accumulated;
			} else {
				return iter(arr.tail(), accumulated.push(fn(arr.first())));
			}
		};
		return iter(arr, []);
	};
	var double = func(x) {
		return x * 2;
	};
	var a = [1,2,3,4,5];
	var result = map(a, double);
	result;
	`

	stmt2 := `
	var reduce = func(arr, initial, fn) {
		var iter = func(arr, result) {
			if (len(arr) == 0) {
				return result;
			} else {
				return iter(arr.tail(), fn(result, arr.first()));
			}
		};
		return iter(arr, initial);
	};
	var add = func(x, y) {
		return x + y;
	};
	var sum = func(arr) {
		return reduce(arr, 0, add);
	};
	var rc = sum([1,2,3,4,5]);
	rc;
	`

	tests := []struct {
		input    string
		expected interface{}
	}{
		{stmt1, []int64{2, 4, 6, 8, 10}},
		{stmt2, 15},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestComplexFuncStmt(t *testing.T) {
	stmt1 := `
	func map(arr, fn) {
		func iter(arr, accumulated) {
			if (len(arr) == 0) {
				return accumulated;
			} else {
				return iter(arr.tail(), accumulated.push(fn(arr.first())));
			}
		};
		return iter(arr, []);
	};
	func double(x) {
		return x * 2;
	};
	var a = [1,2,3,4,5];
	var result = map(a, double);
	result;
	`

	stmt2 := `
	func reduce(arr, initial, fn) {
		func iter(arr, result) {
			if (len(arr) == 0) {
				return result;
			} else {
				return iter(arr.tail(), fn(result, arr.first()));
			}
		};
		return iter(arr, initial);
	};
	func add(x, y) {
		return x + y;
	};
	func sum(arr) {
		return reduce(arr, 0, add);
	};
	var rc = sum([1,2,3,4,5]);
	rc;
	`

	tests := []struct {
		input    string
		expected interface{}
	}{
		{stmt1, []int64{2, 4, 6, 8, 10}},
		{stmt2, 15},
	}
	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		if nil != err {
			t.Fatal(err)
		}
		testEvalObject(t, evaluated, tt.expected)
	}
}

func TestHash(t *testing.T) {
	input := `
	var k2 = "k2";
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
