package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// assignIndexStmt : implement stmtDecoder
type assignIndexStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *assignIndexStmt) match() bool {
	if !this.scanner.ExpectCur2(token.IDENT, token.LBRACK) {
		return false
	}
	this.scanner.NextToken()
	this.scanner.NextToken()
	if _, err := this.p.ParseExpression(scanner.PRECED_LOWEST); nil != err {
		return false
	}
	return this.scanner.ExpectPeek2(token.RBRACK, token.ASSIGN)
}

func (this *assignIndexStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := this.scanner.NewAssignIndex()
	if err := this.scanner.ExpectPeek(token.LBRACK); nil != err {
		return nil, function.NewError(err)
	}
	this.scanner.NextToken()

	idx, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}

	if err := this.scanner.ExpectPeek(token.RBRACK); nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.ASSIGN); nil != err {
		return nil, function.NewError(err)
	}

	this.scanner.NextToken()

	val, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Idx = idx
	stmt.Value = val

	for !this.scanner.StmtEnd(endTok) {
		this.scanner.NextToken()
	}
	return stmt, nil
}
