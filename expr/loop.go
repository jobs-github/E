package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// loopExpr : implement tokenDecoder
type loopExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *loopExpr) decodeExpr(final bool) (ast.Expression, error) {
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

func (this *loopExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewLoop()
	if err := this.scanner.ExpectPeek(token.LPAREN); nil != err {
		return nil, function.NewError(err)
	}
	var err error
	expr.Cnt, err = this.decodeExpr(false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Body, err = this.decodeExpr(true)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}
