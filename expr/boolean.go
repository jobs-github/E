package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/scanner"
)

// boolean : implement tokenDecoder
type boolean struct {
	scanner scanner.Scanner
}

func (this *boolean) decode() (ast.Expression, error) {
	return this.scanner.NewBoolean(), nil
}
