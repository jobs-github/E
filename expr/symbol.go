package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/scanner"
)

// symbol : implement tokenDecoder
type symbol struct {
	scanner scanner.Scanner
}

func (this *symbol) decode() (ast.Expression, error) {
	// TODO
	return nil, nil
}
