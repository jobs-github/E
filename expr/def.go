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
			token.LOOP:   &loopExpr{s, p},
			token.MAP:    &mapExpr{s, p},
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

func decodeExpr(s scanner.Scanner, p interfaces.Parser, final bool) (ast.Expression, error) {
	s.NextToken()
	v, err := p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if !final {
		if err := s.ExpectPeek(token.COMMA); nil != err {
			return nil, function.NewError(err)
		}
	}
	return v, nil
}

func decodeLoopFn(s scanner.Scanner, p interfaces.Parser) (ast.Expression, ast.Expression, error) {
	if err := s.ExpectPeek(token.LPAREN); nil != err {
		return nil, nil, function.NewError(err)
	}
	data, err := decodeExpr(s, p, false)
	if nil != err {
		return nil, nil, function.NewError(err)
	}
	body, err := decodeExpr(s, p, true)
	if nil != err {
		return nil, nil, function.NewError(err)
	}
	if err := s.ExpectPeek(token.RPAREN); nil != err {
		return nil, nil, function.NewError(err)
	}
	return data, body, nil
}
