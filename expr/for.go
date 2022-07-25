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

func (this *forExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewFor()

	this.scanner.NextToken()

	if nil != this.scanner.CurrentIs(token.LBRACE) {
		if err := this.decodeEnv(expr); nil != err {
			return nil, function.NewError(err)
		}
	}

	stmt, err := this.p.ParseBlockStmt()
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Loop = stmt
	return expr, nil
}

func (this *forExpr) decodeEnv(expr *ast.ForExpr) error {
	if err := this.decodeInit(expr); nil != err {
		return function.NewError(err)
	}
	if err := this.decodeCond(expr); nil != err {
		return function.NewError(err)
	}
	if err := this.decodePost(expr); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *forExpr) skip() bool {
	return nil == this.scanner.CurrentIs(token.SEMICOLON)
}

func (this *forExpr) decodeInit(f *ast.ForExpr) error {
	if this.skip() {
		this.scanner.NextToken()
		return nil
	}
	stmt, err := this.p.ParseStmt(token.SEMICOLON)
	if nil != err {
		return function.NewError(err)
	}
	f.Init = stmt
	if err := this.scanner.CurrentIs(token.SEMICOLON); nil != err {
		return function.NewError(err)
	}
	this.scanner.NextToken()
	return nil
}

func (this *forExpr) decodeCond(f *ast.ForExpr) error {
	if this.skip() {
		return nil
	}
	expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return function.NewError(err)
	}
	f.Cond = expr
	if err := this.scanner.PeekIs(token.SEMICOLON); nil != err {
		return function.NewError(err)
	}
	this.scanner.NextToken()
	return nil
}

func (this *forExpr) decodePost(f *ast.ForExpr) error {
	if err := this.scanner.PeekIs(token.LBRACE); nil == err {
		this.scanner.NextToken()
		return nil
	}
	this.scanner.NextToken() // ship ";"
	stmt, err := this.p.ParseStmt(token.LBRACE)
	if nil != err {
		return function.NewError(err)
	}
	f.Post = stmt
	return nil
}
