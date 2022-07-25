package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/scanner"
)

// identifier : implement tokenDecoder
type identifier struct {
	scanner scanner.Scanner
}

func (this *identifier) decode() (ast.Expression, error) {
	return this.scanner.GetIdentifier(), nil
}
