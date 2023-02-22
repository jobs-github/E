package expr

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

type ExprParser interface {
	Decode(tok token.TokenType) (ast.Expression, error)
}

func NewExprParser(s scanner.Scanner, p interfaces.Parser) ExprParser {
	bd := &boolean{s}
	pd := &prefixExpr{s, p}

	return &exprParser{
		m: map[token.TokenType]tokenDecoder{
			token.IDENT:  &identifier{s},
			token.INT:    &integer{s},
			token.STRING: &string{s},
			token.TRUE:   bd,
			token.FALSE:  bd,
			token.NULL:   &null{s},
			token.NOT:    pd,
			token.SUB:    pd,
			token.LPAREN: &lparen{s, p},
			token.LBRACK: &lbrack{s, p},
			token.LBRACE: &lbrace{s, p},
			token.FUNC:   &lambdaFunction{s, p},
			token.FOR:    &forExpr{s, p},
		},
	}
}

type tokenDecoder interface {
	decode() (ast.Expression, error)
}

// exprParser : implement ExprParser
type exprParser struct {
	m map[token.TokenType]tokenDecoder
}

func (this *exprParser) Decode(tok token.TokenType) (ast.Expression, error) {
	fn, ok := this.m[tok]
	if !ok {
		err := fmt.Errorf("%v has no parser", token.ToString(tok))
		return nil, function.NewError(err)
	}
	return fn.decode()
}
