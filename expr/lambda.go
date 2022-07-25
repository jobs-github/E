package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
)

// lambdaFunction : implement tokenDecoder
type lambdaFunction struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *lambdaFunction) decode() (ast.Expression, error) {
	return this.scanner.ParseFunction(true, this.p)
}
