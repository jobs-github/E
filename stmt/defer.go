package stmt

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
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
