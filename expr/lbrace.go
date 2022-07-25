package expr

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
)

// lbrace : implement tokenDecoder
type lbrace struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *lbrace) eof() bool {
	return nil == this.scanner.PeekIs(token.RBRACE)
}

func (this *lbrace) decode() (ast.Expression, error) {
	h := this.scanner.NewHash()
	// empty hash
	if this.eof() {
		this.scanner.NextToken()
		return h, nil
	}
	for !this.eof() {
		this.scanner.NextToken()
		key, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
		if nil != err {
			return nil, function.NewError(err)
		}
		if err := this.scanner.ExpectPeek(token.COLON); nil != err {
			return nil, function.NewError(err)
		}
		this.scanner.NextToken()
		val, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
		if nil != err {
			return nil, function.NewError(err)
		}
		h.Pairs[key] = val
		if this.eof() {
			break
		}
		if err := this.scanner.ExpectPeek(token.COMMA); nil != err {
			return nil, function.NewError(err)
		}
	}
	if err := this.scanner.ExpectPeek(token.RBRACE); nil != err {
		return nil, function.NewError(err)
	}
	return h, nil
}
