package escript

import (
	"github.com/jobs-github/escript/ast"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/parser"
)

func LoadAst(code string) (ast.Node, error) {
	l := lexer.New(code)
	p, err := parser.New(l)
	if nil != err {
		return nil, function.NewError(err)
	}
	return p.ParseProgram()
}
