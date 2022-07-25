package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// functionStmt : implement stmtDecoder
type functionStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *functionStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	this.scanner.NextToken()
	stmt := this.scanner.NewFunction()
	fn, err := this.scanner.ParseFunction(false, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Value = fn
	for !this.scanner.StmtEnd(endTok) {
		this.scanner.NextToken()
	}
	return stmt, nil
}
