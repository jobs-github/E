package stmt

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
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
