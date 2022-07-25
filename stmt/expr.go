package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// exprStmt : implement stmtDecoder
type exprStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *exprStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := this.scanner.NewExpr()
	expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Expr = expr

	if err := this.scanner.PeekIs(token.SEMICOLON); nil == err {
		this.scanner.NextToken()
	}
	return stmt, nil
}
