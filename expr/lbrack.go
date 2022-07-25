package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// lbrack : implement tokenDecoder
type lbrack struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *lbrack) decode() (ast.Expression, error) {
	expr := this.scanner.NewArray()
	v, err := this.p.ParseExpressions(token.RBRACK)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Items = v
	return expr, nil
}
