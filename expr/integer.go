package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/scanner"
)

// integer : implement tokenDecoder
type integer struct {
	scanner scanner.Scanner
}

func (this *integer) decode() (ast.Expression, error) {
	return this.scanner.NewInteger()
}
