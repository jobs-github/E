package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// assignStmt : implement stmtDecoder
type assignStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *assignStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := this.scanner.NewAssign()
	if err := this.scanner.ExpectPeek(token.ASSIGN); nil != err {
		return nil, function.NewError(err)
	}

	this.scanner.NextToken()

	expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Value = expr

	for !this.scanner.StmtEnd(endTok) {
		this.scanner.NextToken()
	}
	return stmt, nil
}
