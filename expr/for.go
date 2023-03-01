package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// forExpr : implement tokenDecoder
type forExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *forExpr) decodeExpr(final bool) (ast.Expression, error) {
	this.scanner.NextToken()
	v, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if !final {
		if err := this.scanner.ExpectPeek(token.COMMA); nil != err {
			return nil, function.NewError(err)
		}
	}
	return v, nil
}

func (this *forExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewFor()
	if err := this.scanner.ExpectPeek(token.LPAREN); nil != err {
		return nil, function.NewError(err)
	}
	var err error
	expr.Init, err = this.decodeExpr(false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Cond, err = this.decodeExpr(false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Next, err = this.decodeExpr(false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.LoopFn, err = this.decodeExpr(true)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}
