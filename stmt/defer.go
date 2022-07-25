package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// deferStmt : implement stmtDecoder
type deferStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *deferStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := this.scanner.NewDefer()
	this.scanner.NextToken()
	do, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Do = do
	if err := this.scanner.PeekIs(token.SEMICOLON); nil == err {
		this.scanner.NextToken()
	}
	return stmt, nil
}
