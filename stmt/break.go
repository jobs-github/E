package stmt

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// breakStmt : implement stmtDecoder
type breakStmt struct {
	scanner scanner.Scanner
}

func (this *breakStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	return this.scanner.NewBreak(), nil
}
