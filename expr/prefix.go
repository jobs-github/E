package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
)

// prefixExpr : implement tokenDecoder
type prefixExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *prefixExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewPrefix()
	this.scanner.NextToken()
	right, err := this.p.ParseExpression(scanner.PRECED_PREFIX)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Right = right
	return expr, nil
}
