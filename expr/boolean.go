package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/scanner"
)

// boolean : implement tokenDecoder
type boolean struct {
	scanner scanner.Scanner
}

func (this *boolean) decode() (ast.Expression, error) {
	return this.scanner.NewBoolean(), nil
}
