package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// lbrack : implement tokenDecoder
type lbrack struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *lbrack) decode() (ast.Expression, error) {
	expr := ast.NewArray()
	v, err := this.p.ParseExpressions(token.RBRACK)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Items = v
	return expr, nil
}
