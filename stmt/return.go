package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// returnStmt : implement stmtDecoder
type returnStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *returnStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := this.scanner.NewReturn()
	this.scanner.NextToken()

	if !this.scanner.StmtEnd(endTok) {
		expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
		if nil != err {
			return nil, function.NewError(err)
		}
		stmt.ReturnValue = expr
	} else {
		stmt.ReturnValue = ast.NewNull()
	}

	for !this.scanner.StmtEnd(endTok) {
		this.scanner.NextToken()
	}
	return stmt, nil
}
