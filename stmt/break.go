package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// breakStmt : implement stmtDecoder
type breakStmt struct {
	scanner scanner.Scanner
}

func (this *breakStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	return this.scanner.NewBreak(), nil
}
