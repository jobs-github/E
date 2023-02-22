package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
)

// forExpr : implement tokenDecoder
type forExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *forExpr) decode() (ast.Expression, error) {
	// TODO
	return nil, nil
}
