package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/scanner"
)

// null : implement tokenDecoder
type null struct {
	scanner scanner.Scanner
}

func (this *null) decode() (ast.Expression, error) {
	return this.scanner.NewNull(), nil
}
