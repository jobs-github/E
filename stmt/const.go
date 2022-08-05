package stmt

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/builtin"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// constStmt : implement stmtDecoder
type constStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *constStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := this.scanner.NewConst()
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
