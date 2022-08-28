package eval

import (
	"fmt"
	"io"
	"strings"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/parser"
)

func EvalAst(node ast.Node, env object.Env) (object.Object, error) {
	v := newVisitor(env)
	if err := node.Do(v); nil != err {
		return object.Nil, function.NewError(err)
	}
	return v.r, nil
}

type Eval interface {
	Repl(baseDir string, in io.Reader, out io.Writer)
	EvalJson(path string, args []string)
	EvalScript(path string, args []string)
	EvalCode(code string, args []string)

	NewEnv(args []string) object.Env

	DumpAst(path string) (string, error)
	LoadJson(path string) (ast.Node, error)
	LoadAst(code string) (ast.Node, error)
}

func New(nonrecursive bool) Eval { return evalImpl{nonrecursive: nonrecursive} }

func loadCode(path string) ([]byte, error) {
	if !strings.HasSuffix(path, ast.SuffixQs) {
		err := fmt.Errorf(`file "%v" not endwith ".qs"`, path)
		return nil, function.NewError(err)
	}
	b, err := function.LoadFile(path)
	if nil != err {
		return nil, function.NewError(err)
	}
	return b, nil
}

func loadAst(code string) (ast.Node, error) {
	l := lexer.New(code)
	p, err := parser.New(l)
	if nil != err {
		return nil, function.NewError(err)
	}
	return p.ParseProgram()
}
