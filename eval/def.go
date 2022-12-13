package eval

import (
	"io"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/parser"
)

type Eval interface {
	Repl(in io.Reader, out io.Writer)
	EvalJson(path string)
	EvalScript(path string)
	EvalCode(code string)

	DumpAst(path string) (string, error)
	LoadJson(path string) (ast.Node, error)
	LoadAst(code string) (ast.Node, error)
}

func NewInterpreter() Eval { return interpreter{} }
func NewState() Eval       { return virtualMachine{} }

func loadAst(code string) (ast.Node, error) {
	l := lexer.New(code)
	p, err := parser.New(l)
	if nil != err {
		return nil, function.NewError(err)
	}
	return p.ParseProgram()
}
