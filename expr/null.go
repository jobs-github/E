package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/scanner"
)

// null : implement tokenDecoder
type null struct {
	scanner scanner.Scanner
}

func (this *null) decode() (ast.Expression, error) {
	return this.scanner.NewNull(), nil
}
