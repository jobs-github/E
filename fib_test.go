package escript

import (
	"testing"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/vm"
)

// go test -v fib_test.go

const (
	// 30
	BENCH_CODE = `
	func fib(x) {
		(x == 0) ? 0 : (
			(x == 1) ? 1 : (
				fib(x - 1) + fib(x - 2)
			)
		)
	};
	fib(3);
	`
)

var (
	BENCH_AST = newAst(BENCH_CODE)
	BENCH_VM  = newVM(BENCH_AST)
)

func newAst(code string) ast.Node {
	node, err := LoadAst(code)
	if nil != err {
		panic(err)
	}
	return node
}

func newVM(node ast.Node) vm.VM {
	c := compiler.New()
	if err := c.Compile(node); nil != err {
		panic(err)
	}
	return vm.New(c.Bytecode(), c.Constants())
}

func TestFibExpr(t *testing.T) {
	r, err := BENCH_AST.Eval(object.NewEnv())
	if err != nil {
		panic(err)
	}
	t.Logf("result: %v", r)
}

func TestFibVM(t *testing.T) {
	err := BENCH_VM.Run()
	if err != nil {
		panic(err)
	}
	r := BENCH_VM.LastPopped()
	t.Logf("result: %v", r)
}
