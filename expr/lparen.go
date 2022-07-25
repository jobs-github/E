package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// lparen : implement tokenDecoder
type lparen struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *lparen) decode() (ast.Expression, error) {
	this.scanner.NextToken()
	expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}
