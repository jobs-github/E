package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
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
