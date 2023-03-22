package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// reduceExpr : implement tokenDecoder
type reduceExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *reduceExpr) decode() (ast.Expression, error) {
	expr := ast.NewReduce()
	if err := this.scanner.ExpectPeek(token.LPAREN); nil != err {
		return nil, function.NewError(err)
	}
	var err error
	expr.Arr, err = decodeExpr(this.scanner, this.p, false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Body, err = decodeExpr(this.scanner, this.p, false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Init, err = decodeExpr(this.scanner, this.p, true)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}
