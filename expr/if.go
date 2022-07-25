package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// ifExpr : implement tokenDecoder
type ifExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *ifExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewIf()

	clause, err := this.decodeClause()
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Clauses = append(expr.Clauses, clause)
	for {
		if this.scanner.ExpectPeek2(token.ELSE, token.IF) {
			if err := this.decodeElseIf(expr); nil != err {
				return nil, function.NewError(err)
			}
		} else if nil == this.scanner.PeekIs(token.ELSE) {
			if err := this.decodeElse(expr); nil != err {
				return nil, function.NewError(err)
			}
			break
		} else {
			break
		}
	}
	return expr, nil
}

func (this *ifExpr) decodeClause() (*ast.IfClause, error) {
	clause := &ast.IfClause{}
	if err := this.scanner.ExpectPeek(token.LPAREN); nil != err {
		return nil, function.NewError(err)
	}
	this.scanner.NextToken()
	expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	clause.If = expr
	if err := this.scanner.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.LBRACE); nil != err {
		return nil, function.NewError(err)
	}
	stmt, err := this.p.ParseBlockStmt()
	if nil != err {
		return nil, function.NewError(err)
	}
	clause.Then = stmt
	return clause, nil
}

func (this *ifExpr) decodeElseIf(expr *ast.IfExpr) error {
	this.scanner.NextToken()
	this.scanner.NextToken()

	clause, err := this.decodeClause()
	if nil != err {
		return function.NewError(err)
	}
	expr.Clauses = append(expr.Clauses, clause)
	return nil
}

func (this *ifExpr) decodeElse(expr *ast.IfExpr) error {
	this.scanner.NextToken()
	if err := this.scanner.ExpectPeek(token.LBRACE); nil != err {
		return function.NewError(err)
	}
	stmt, err := this.p.ParseBlockStmt()
	if nil != err {
		return function.NewError(err)
	}
	expr.Else = stmt
	return nil
}
