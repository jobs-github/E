package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/scanner"
)

// integer : implement tokenDecoder
type integer struct {
	scanner scanner.Scanner
}

func (this *integer) decode() (ast.Expression, error) {
	return this.scanner.NewInteger()
}
