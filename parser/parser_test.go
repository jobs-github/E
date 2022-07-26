package parser

import (
	"reflect"
	"testing"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/lexer"
)

func parseProgram(t *testing.T, p interfaces.Parser) *ast.Program {
	program, err := p.ParseProgram()
	if nil != err {
		t.Fatal(err)
	}
	rc, ok := program.(*ast.Program)
	if !ok {
		t.Fatal(err)
	}
	return rc
}

func testVarStatements(t *testing.T, s ast.Statement, name string) bool {
	varStmt, ok := s.(*ast.VarStmt)
	if !ok {
		t.Errorf("s is not *ast.VarStatement, got=%v", reflect.TypeOf(s).String())
		return false
	}
	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value != %v, got=%v", name, varStmt.Name.Value)
		return false
	}
	return true
}

func TestVarStatements(t *testing.T) {
	cases := []struct {
		input     string
		wantIdent string
		wantValue interface{}
	}{
		{"var x = 5;", "x", 5},
		{"var y = true;", "y", true},
		{"var foobar = y;", "foobar", "y"},
	}
	for _, tt := range cases {
		l := lexer.New(tt.input)
		p, err := New(l)
		if nil != err {
			t.Fatal(err)
		}

		program := parseProgram(t, p)
		if len(program.Stmts) != 1 {
			t.Fatalf("number of program Statements: %v", len(program.Stmts))
		}
		stmt := program.Stmts[0]
		if !testVarStatements(t, stmt, tt.wantIdent) {
			return
		}
		val := stmt.(*ast.VarStmt).Value
		if !testLiteralExpression(t, val, tt.wantValue) {
			return
		}
	}
}

func TestIdentExpr(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}

	program := parseProgram(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("number of program Statements: %v", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	ident, ok := stmt.Expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expr is not *ast.Identifier, got %v", reflect.TypeOf(stmt.Expr).String())
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value != `foobar`, got %v", ident.Value)
	}
}

func TestIntExpr(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}

	program := parseProgram(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("number of program Statements: %v", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	literal, ok := stmt.Expr.(*ast.Integer)
	if !ok {
		t.Fatalf("Expr is not *ast.Integer, got %v", reflect.TypeOf(stmt.Expr).String())
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value != 5, got %v", literal.Value)
	}
}

func TestStringExpr(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}

	program := parseProgram(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("number of program Statements: %v", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	literal, ok := stmt.Expr.(*ast.String)
	if !ok {
		t.Fatalf("Expr is not *ast.String, got %v", reflect.TypeOf(stmt.Expr).String())
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value != hello world, got %v", literal.Value)
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr not *ast.Identifier, got %v", reflect.TypeOf(expr).String())
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value != %v, got %v", value, ident.Value)
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, val int64) bool {
	iv, ok := expr.(*ast.Integer)
	if !ok {
		t.Errorf("expr not *ast.IntegerLiteral, got %v", reflect.TypeOf(expr).String())
		return false
	}
	if iv.Value != val {
		t.Errorf("val not %v, got %v", val, iv.Value)
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, val bool) bool {
	iv, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("expr not *ast.Boolean, got %v", reflect.TypeOf(expr).String())
		return false
	}
	if iv.Value != val {
		t.Errorf("val not %v, got %v", val, iv.Value)
		return false
	}
	return true
}

func TestBoolExpr(t *testing.T) {
	input := `true;`

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}

	program := parseProgram(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("number of program Statements: %v", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	literal, ok := stmt.Expr.(*ast.Boolean)
	if !ok {
		t.Fatalf("Expr is not *ast.Boolean, got %v", reflect.TypeOf(stmt.Expr).String())
	}
	if literal.Value != true {
		t.Errorf("literal.Value != true, got %v", literal.Value)
	}
}

func testLiteralExpression(t *testing.T, expr ast.Expression, want interface{}) bool {
	switch v := want.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}
	t.Errorf("type of expr unsupported: %v", reflect.TypeOf(expr).String())
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, op string, right interface{}) bool {
	infixExpr, ok := expr.(*ast.InfixExpr)
	if !ok {
		t.Errorf("expr is not ast.InfixExpression, got %v", reflect.TypeOf(expr).String())
		return false
	}
	if !testLiteralExpression(t, infixExpr.Left, left) {
		return false
	}
	if infixExpr.Op.Literal != op {
		t.Errorf("infixExpr.Op != %v, got %v", op, infixExpr.Op)
		return false
	}
	if !testLiteralExpression(t, infixExpr.Right, right) {
		return false
	}
	return true
}

func TestParsingPrefixExpressions(t *testing.T) {
	cases := []struct {
		input string
		op    string
		val   interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}
	for _, tt := range cases {
		l := lexer.New(tt.input)
		p, err := New(l)
		if nil != err {
			t.Fatal(err)
		}
		program := parseProgram(t, p)
		if len(program.Stmts) != 1 {
			t.Fatalf("number of program Statements: %v", len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
		}
		expr, ok := stmt.Expr.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("Expr is not *ast.PrefixExpression, got %v", reflect.TypeOf(stmt.Expr).String())
		}
		if expr.Op.Literal != tt.op {
			t.Errorf("expr.Op != %v, got %v", tt.op, expr.Op)
		}
		if testLiteralExpression(t, expr.Right, tt.val) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	cases := []struct {
		input string
		left  interface{}
		op    string
		right interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range cases {
		l := lexer.New(tt.input)
		p, err := New(l)
		if nil != err {
			t.Fatal(err)
		}
		program := parseProgram(t, p)
		if len(program.Stmts) != 1 {
			t.Fatalf("number of program Statements: %v", len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
		}
		expr, ok := stmt.Expr.(*ast.InfixExpr)
		if !ok {
			t.Fatalf("Expr is not *ast.InfixExpression, got %v", reflect.TypeOf(stmt.Expr).String())
		}
		if !testLiteralExpression(t, expr.Left, tt.left) {
			return
		}
		if expr.Op.Literal != tt.op {
			t.Errorf("expr.Op != %v, got %v", tt.op, expr.Op)
		}
		if !testLiteralExpression(t, expr.Right, tt.right) {
			return
		}
	}
}

func TestOpPrecedParsing(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b % c", "(a + (b % c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 >= 4 == 3 <= 4", "((5 >= 4) == (3 <= 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"a && b || c", "((a && b) || c)"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}
	for _, tt := range cases {
		l := lexer.New(tt.input)
		p, err := New(l)
		if nil != err {
			t.Fatal(err)
		}
		program := parseProgram(t, p)
		str := program.String()
		if tt.want != str {
			t.Errorf("expected %v, want %v", tt.want, str)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("number of program Statements: %v", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	expr, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("Expr is not *ast.IfExpression, got %v", reflect.TypeOf(stmt.Expr).String())
	}
	if 1 != len(expr.Clauses) {
		t.Fatalf("number of expr.Clauses: %v", len(expr.Clauses))
	}

	if !testInfixExpression(t, expr.Clauses[0].If, "x", "<", "y") {
		return
	}
	if 1 != len(expr.Clauses[0].Then.Stmts) {
		t.Fatalf("number of expr.Clauses[0].Then.Stmts: %v", len(expr.Clauses[0].Then.Stmts))
	}
	thenstmt := expr.Clauses[0].Then.Stmts[0]
	then, ok := thenstmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt, got %v", reflect.TypeOf(thenstmt).String())
	}
	if !testIdentifier(t, then.Expr, "x") {
		return
	}
	if expr.Else != nil {
		t.Fatalf("expr.Else != nil")
	}
}

func TestIfClausesExpression(t *testing.T) {
	input := `
	if (x < y) { 
		x
	} else if (x > y) {
		y
	} else {
		z
	}
	`
	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("number of program Statements: %v", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	expr, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("Expr is not *ast.IfExpression, got %v", reflect.TypeOf(stmt.Expr).String())
	}
	if 2 != len(expr.Clauses) {
		t.Fatalf("number of expr.Clauses: %v", len(expr.Clauses))
	}

	if !testInfixExpression(t, expr.Clauses[0].If, "x", "<", "y") {
		return
	}
	if 1 != len(expr.Clauses[0].Then.Stmts) {
		t.Fatalf("number of expr.Clauses[0].Then.Stmts: %v", len(expr.Clauses[0].Then.Stmts))
	}
	thenstmt := expr.Clauses[0].Then.Stmts[0]
	then, ok := thenstmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt, got %v", reflect.TypeOf(thenstmt).String())
	}
	if !testIdentifier(t, then.Expr, "x") {
		return
	}

	if !testInfixExpression(t, expr.Clauses[1].If, "x", ">", "y") {
		return
	}
	if 1 != len(expr.Clauses[1].Then.Stmts) {
		t.Fatalf("number of expr.Clauses[0].Then.Stmts: %v", len(expr.Clauses[1].Then.Stmts))
	}
	thenstmt2 := expr.Clauses[1].Then.Stmts[0]
	then2, ok := thenstmt2.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStmt, got %v", reflect.TypeOf(thenstmt2).String())
	}
	if !testIdentifier(t, then2.Expr, "y") {
		return
	}

	if expr.Else == nil {
		t.Fatalf("expr.Else == nil")
	}

	if 1 != len(expr.Else.Stmts) {
		t.Fatalf("number of expr.Else.Stmts: %v", len(expr.Else.Stmts))
	}
	elsestmt := expr.Else.Stmts[0]
	thenelse, ok := elsestmt.(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("elsestmt is not *ast.ExpressionStmt, got %v", reflect.TypeOf(elsestmt).String())
	}
	if !testIdentifier(t, thenelse.Expr, "z") {
		return
	}
}

func TestFunctionParsing(t *testing.T) {
	input := `func(x, y) { x + y }`
	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("number of program Statements: %v", len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	expr, ok := stmt.Expr.(*ast.Function)
	if !ok {
		t.Fatalf("Expr is not *ast.Function, got %v", reflect.TypeOf(stmt.Expr).String())
	}
	if 2 != len(expr.Args) {
		t.Fatalf("number of expr.Args: %v", len(expr.Args))
	}
	testLiteralExpression(t, expr.Args[0], "x")
	testLiteralExpression(t, expr.Args[1], "y")

	if len(expr.Body.Stmts) != 1 {
		t.Fatalf("number of expr.Body.Stmts: %v", len(expr.Body.Stmts))
	}
	bodyStmt, ok := expr.Body.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("expr.Body.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	testInfixExpression(t, bodyStmt.Expr, "x", "+", "y")
}

func TestFuncArgsParsing(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{input: "func() {};", want: []string{}},
		{input: "func(x) {};", want: []string{"x"}},
		{input: "func(x, y, z) {};", want: []string{"x", "y", "z"}},
	}

	for _, tt := range cases {
		l := lexer.New(tt.input)
		p, err := New(l)
		if nil != err {
			t.Fatal(err)
		}
		program := parseProgram(t, p)
		stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
		}
		function, ok := stmt.Expr.(*ast.Function)
		if !ok {
			t.Fatalf("stmt.Expr is not *ast.Function, got %v", reflect.TypeOf(program.Stmts[0]).String())
		}
		if len(function.Args) != len(tt.want) {
			t.Fatalf("len(function.Args) != %v, got %v", len(tt.want), len(function.Args))
		}
		for i, ident := range tt.want {
			testLiteralExpression(t, function.Args[i], ident)
		}
	}
}

func TestCallParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	expr, ok := stmt.Expr.(*ast.Call)
	if !ok {
		t.Fatalf("stmt.Expr is not *ast.Call, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	if !testIdentifier(t, expr.Func, "add") {
		return
	}
	if len(expr.Args) != 3 {
		t.Fatalf("wrong len of args, got %v", len(expr.Args))
	}
	testLiteralExpression(t, expr.Args[0], 1)
	testInfixExpression(t, expr.Args[1], 2, "*", 3)
	testInfixExpression(t, expr.Args[2], 4, "+", 5)
}

func TestArrayParsing(t *testing.T) {
	input := "[1, 2 * 3, 4 + 5]"

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	arr, ok := stmt.Expr.(*ast.Array)
	if !ok {
		t.Fatalf("stmt.Expr is not *ast.Array, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	if len(arr.Items) != 3 {
		t.Fatalf("wrong len of args, got %v", len(arr.Items))
	}
	testIntegerLiteral(t, arr.Items[0], 1)
	testInfixExpression(t, arr.Items[1], 2, "*", 3)
	testInfixExpression(t, arr.Items[2], 4, "+", 5)
}

func TestIndexParsing(t *testing.T) {
	input := "testArr[1 + 10]"

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	idx, ok := stmt.Expr.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not *ast.IndexExpr, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	testIdentifier(t, idx.Left, "testArr")
	testInfixExpression(t, idx.Index, 1, "+", 10)
}

func TestHashParsing(t *testing.T) {
	input := `{"k1": 1, "k2": 2, "k3": 3}`

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	h, ok := stmt.Expr.(*ast.Hash)
	if !ok {
		t.Fatalf("stmt.Expr is not *ast.Hash, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	if len(h.Pairs) != 3 {
		t.Fatalf("Pairs has wrong len, got %v", len(h.Pairs))
	}
	expected := map[string]int64{"k1": 1, "k2": 2, "k3": 3}
	for k, v := range h.Pairs {
		s, ok := k.(*ast.String)
		if !ok {
			t.Errorf("k is not *ast.String, got %v", reflect.TypeOf(k))
		}
		ev := expected[s.String()]
		testIntegerLiteral(t, v, ev)
	}
}

func TestHashExprParsing(t *testing.T) {
	input := `{"k1": 1 + 1, "k2": 100 - 90, "k3": 30 / 10}`

	l := lexer.New(input)
	p, err := New(l)
	if nil != err {
		t.Fatal(err)
	}
	program := parseProgram(t, p)
	stmt, ok := program.Stmts[0].(*ast.ExpressionStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not *ast.ExpressionStmt, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	h, ok := stmt.Expr.(*ast.Hash)
	if !ok {
		t.Fatalf("stmt.Expr is not *ast.Hash, got %v", reflect.TypeOf(program.Stmts[0]).String())
	}
	if len(h.Pairs) != 3 {
		t.Fatalf("Pairs has wrong len, got %v", len(h.Pairs))
	}
	expected := map[string]func(ast.Expression){
		"k1": func(e ast.Expression) {
			testInfixExpression(t, e, 1, "+", 1)
		},
		"k2": func(e ast.Expression) {
			testInfixExpression(t, e, 100, "-", 90)
		},
		"k3": func(e ast.Expression) {
			testInfixExpression(t, e, 30, "/", 10)
		},
	}
	for k, v := range h.Pairs {
		s, ok := k.(*ast.String)
		if !ok {
			t.Errorf("k is not *ast.String, got %v", reflect.TypeOf(k))
		}
		fn := expected[s.String()]
		fn(v)
	}
}
