package parser

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

type exprParser interface {
	Decode(tok token.TokenType) (ast.Expression, error)
}

func newExprParser(s scanner, p Parser) exprParser {
	bd := &boolean{s}
	pd := &prefixExpr{s, p}

	return &exprParserImpl{
		m: map[token.TokenType]tokenDecoder{
			token.SYMBOL: &symbol{s},
			token.IDENT:  &identifier{s},
			token.INT:    &integer{s},
			token.STRING: &stringExpr{s},
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
			token.REDUCE: &reduceExpr{s, p},
			token.FILTER: &filterExpr{s, p},
		},
	}
}

type tokenDecoder interface {
	decode() (ast.Expression, error)
}

// exprParserImpl : implement exprParser
type exprParserImpl struct {
	m map[token.TokenType]tokenDecoder
}

func (this *exprParserImpl) Decode(tok token.TokenType) (ast.Expression, error) {
	fn, ok := this.m[tok]
	if !ok {
		err := fmt.Errorf("%v has no parser", token.ToString(tok))
		return nil, function.NewError(err)
	}
	return fn.decode()
}

func decodeExpr(s scanner, p Parser, final bool) (ast.Expression, error) {
	s.NextToken()
	v, err := p.ParseExpression(PRECED_LOWEST)
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

func decodeLoopFn(s scanner, p Parser) (ast.Expression, ast.Expression, error) {
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

// symbol : implement tokenDecoder
type symbol struct {
	s scanner
}

func (this *symbol) decode() (ast.Expression, error) {
	// TODO
	return nil, nil
}

// identifier : implement tokenDecoder
type identifier struct {
	s scanner
}

func (this *identifier) decode() (ast.Expression, error) {
	return this.s.GetIdentifier(), nil
}

// integer : implement tokenDecoder
type integer struct {
	s scanner
}

func (this *integer) decode() (ast.Expression, error) {
	return this.s.NewInteger()
}

// stringExpr : implement tokenDecoder
type stringExpr struct {
	s scanner
}

func (this *stringExpr) decode() (ast.Expression, error) {
	return this.s.NewString(), nil
}

// boolean : implement tokenDecoder
type boolean struct {
	s scanner
}

func (this *boolean) decode() (ast.Expression, error) {
	return this.s.NewBoolean(), nil
}

// null : implement tokenDecoder
type null struct {
	s scanner
}

func (this *null) decode() (ast.Expression, error) {
	return ast.NewNull(), nil
}

// prefixExpr : implement tokenDecoder
type prefixExpr struct {
	s scanner
	p Parser
}

func (this *prefixExpr) decode() (ast.Expression, error) {
	expr := this.s.NewPrefix()
	this.s.NextToken()
	right, err := this.p.ParseExpression(PRECED_PREFIX)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Right = right
	return expr, nil
}

// lparen : implement tokenDecoder
type lparen struct {
	s scanner
	p Parser
}

func (this *lparen) decode() (ast.Expression, error) {
	this.s.NextToken()
	expr, err := this.p.ParseExpression(PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.s.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}

// lbrack : implement tokenDecoder
type lbrack struct {
	s scanner
	p Parser
}

func (this *lbrack) decode() (ast.Expression, error) {
	expr := ast.NewArray()
	v, err := this.p.ParseExpressions(token.RBRACK)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Items = v
	return expr, nil
}

// lbrace : implement tokenDecoder
type lbrace struct {
	s scanner
	p Parser
}

func (this *lbrace) eof() bool {
	return nil == this.s.PeekIs(token.RBRACE)
}

func (this *lbrace) decode() (ast.Expression, error) {
	h := ast.NewHash()
	h.Pairs = ast.ExpressionMap{}
	// empty hash
	if this.eof() {
		this.s.NextToken()
		return h, nil
	}
	for !this.eof() {
		this.s.NextToken()
		key, err := this.p.ParseExpression(PRECED_LOWEST)
		if nil != err {
			return nil, function.NewError(err)
		}
		if err := this.s.ExpectPeek(token.COLON); nil != err {
			return nil, function.NewError(err)
		}
		this.s.NextToken()
		val, err := this.p.ParseExpression(PRECED_LOWEST)
		if nil != err {
			return nil, function.NewError(err)
		}
		h.Pairs[key] = val
		if this.eof() {
			break
		}
		if err := this.s.ExpectPeek(token.COMMA); nil != err {
			return nil, function.NewError(err)
		}
	}
	if err := this.s.ExpectPeek(token.RBRACE); nil != err {
		return nil, function.NewError(err)
	}
	return h, nil
}

// lambdaFunction : implement tokenDecoder
type lambdaFunction struct {
	s scanner
	p Parser
}

func (this *lambdaFunction) decode() (ast.Expression, error) {
	return this.s.ParseFunction(true, this.p)
}

// loopExpr : implement tokenDecoder
type loopExpr struct {
	s scanner
	p Parser
}

func (this *loopExpr) decode() (ast.Expression, error) {
	expr := ast.NewLoop()
	data, body, err := decodeLoopFn(this.s, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Cnt = data
	expr.Body = body
	return expr, nil
}

// mapExpr : implement tokenDecoder
type mapExpr struct {
	s scanner
	p Parser
}

func (this *mapExpr) decode() (ast.Expression, error) {
	expr := ast.NewMap()
	data, body, err := decodeLoopFn(this.s, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Arr = data
	expr.Body = body
	return expr, nil
}

// reduceExpr : implement tokenDecoder
type reduceExpr struct {
	s scanner
	p Parser
}

func (this *reduceExpr) decode() (ast.Expression, error) {
	expr := ast.NewReduce()
	if err := this.s.ExpectPeek(token.LPAREN); nil != err {
		return nil, function.NewError(err)
	}
	var err error
	expr.Arr, err = decodeExpr(this.s, this.p, false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Body, err = decodeExpr(this.s, this.p, false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Init, err = decodeExpr(this.s, this.p, true)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.s.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}

// filterExpr : implement tokenDecoder
type filterExpr struct {
	s scanner
	p Parser
}

func (this *filterExpr) decode() (ast.Expression, error) {
	expr := ast.NewFilter()
	data, body, err := decodeLoopFn(this.s, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Arr = data
	expr.Body = body
	return expr, nil
}
