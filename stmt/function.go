package stmt

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
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
