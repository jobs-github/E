package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/scanner"
)

// string : implement tokenDecoder
type string struct {
	scanner scanner.Scanner
}

func (this *string) decode() (ast.Expression, error) {
	return this.scanner.NewString(), nil
}
