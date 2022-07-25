package stmt

import (
	"fmt"

	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/builtin"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// varStmt : implement stmtDecoder
type varStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *varStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := this.scanner.NewVar()
	if err := this.scanner.ExpectPeek(token.IDENT); nil != err {
		return nil, function.NewError(err)
	}
	stmt.Name = this.scanner.GetIdentifier()

	if builtin.IsBuiltin(stmt.Name.Value) {
		err := fmt.Errorf("`%v` is built-in function", stmt.Name.Value)
		return nil, function.NewError(err)
	}

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
