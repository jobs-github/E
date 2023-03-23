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

// symbol : implement tokenDecoder
type symbol struct {
	scanner scanner.Scanner
}

func (this *symbol) decode() (ast.Expression, error) {
	// TODO
	return nil, nil
}

// identifier : implement tokenDecoder
type identifier struct {
	scanner scanner.Scanner
}

func (this *identifier) decode() (ast.Expression, error) {
	return this.scanner.GetIdentifier(), nil
}

// integer : implement tokenDecoder
type integer struct {
	scanner scanner.Scanner
}

func (this *integer) decode() (ast.Expression, error) {
	return this.scanner.NewInteger()
}

// stringExpr : implement tokenDecoder
type stringExpr struct {
	scanner scanner.Scanner
}

func (this *stringExpr) decode() (ast.Expression, error) {
	return this.scanner.NewString(), nil
}

// boolean : implement tokenDecoder
type boolean struct {
	scanner scanner.Scanner
}

func (this *boolean) decode() (ast.Expression, error) {
	return this.scanner.NewBoolean(), nil
}

// null : implement tokenDecoder
type null struct {
	scanner scanner.Scanner
}

func (this *null) decode() (ast.Expression, error) {
	return ast.NewNull(), nil
}

// prefixExpr : implement tokenDecoder
type prefixExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *prefixExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewPrefix()
	this.scanner.NextToken()
	right, err := this.p.ParseExpression(scanner.PRECED_PREFIX)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Right = right
	return expr, nil
}

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

// lbrack : implement tokenDecoder
type lbrack struct {
	scanner scanner.Scanner
	p       interfaces.Parser
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
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *lbrace) eof() bool {
	return nil == this.scanner.PeekIs(token.RBRACE)
}

func (this *lbrace) decode() (ast.Expression, error) {
	h := ast.NewHash()
	h.Pairs = ast.ExpressionMap{}
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

// lambdaFunction : implement tokenDecoder
type lambdaFunction struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *lambdaFunction) decode() (ast.Expression, error) {
	return this.scanner.ParseFunction(true, this.p)
}

// loopExpr : implement tokenDecoder
type loopExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *loopExpr) decode() (ast.Expression, error) {
	expr := ast.NewLoop()
	data, body, err := decodeLoopFn(this.scanner, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Cnt = data
	expr.Body = body
	return expr, nil
}

// mapExpr : implement tokenDecoder
type mapExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *mapExpr) decode() (ast.Expression, error) {
	expr := ast.NewMap()
	data, body, err := decodeLoopFn(this.scanner, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Arr = data
	expr.Body = body
	return expr, nil
}

// reduceExpr : implement tokenDecoder
type reduceExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *reduceExpr) decode() (ast.Expression, error) {
	expr := ast.NewReduce()
	if err := this.scanner.ExpectPeek(token.LPAREN); nil != err {
		return nil, function.NewError(err)
	}
	var err error
	expr.Arr, err = decodeExpr(this.scanner, this.p, false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Body, err = decodeExpr(this.scanner, this.p, false)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Init, err = decodeExpr(this.scanner, this.p, true)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return expr, nil
}

// filterExpr : implement tokenDecoder
type filterExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *filterExpr) decode() (ast.Expression, error) {
	expr := ast.NewFilter()
	data, body, err := decodeLoopFn(this.scanner, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Arr = data
	expr.Body = body
	return expr, nil
}
