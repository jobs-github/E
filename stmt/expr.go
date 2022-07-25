package stmt

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
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
