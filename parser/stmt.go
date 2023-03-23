package parser

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/builtin"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

type stmtParser interface {
	Decode(t token.TokenType, endTok token.TokenType) (ast.Statement, error)
}

func newStmtParser(
	s scanner,
	p Parser,
) stmtParser {
	return &stmtParserImpl{
		scanner:         s,
		p:               p,
		functionDecoder: &functionStmt{s, p},
		exprDecoder:     &exprStmt{s, p},
		m: map[token.TokenType]stmtDecoder{
			token.CONST: &constStmt{s, p},
		},
	}
}

type stmtDecoder interface {
	decode(endTok token.TokenType) (ast.Statement, error)
}

// stmtParserImpl : implement StmtParser
type stmtParserImpl struct {
	scanner         scanner
	p               Parser
	functionDecoder stmtDecoder
	exprDecoder     stmtDecoder
	m               map[token.TokenType]stmtDecoder
}

func (this *stmtParserImpl) Decode(t token.TokenType, endTok token.TokenType) (ast.Statement, error) {
	parser, ok := this.m[t]
	if ok {
		return parser.decode(endTok)
	} else {
		if this.isFunctionStmt() {
			return this.functionDecoder.decode(endTok)
		} else {
			return this.exprDecoder.decode(endTok)
		}
	}
}

func (this *stmtParserImpl) isFunctionStmt() bool {
	return this.scanner.ExpectCur2(token.FUNC, token.IDENT)
}

// constStmt : implement stmtDecoder
type constStmt struct {
	s scanner
	p Parser
}

func (this *constStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := ast.NewConst()
	if err := this.s.ExpectPeek(token.IDENT); nil != err {
		return nil, function.NewError(err)
	}
	stmt.Name = this.s.GetIdentifier()

	if builtin.IsBuiltin(stmt.Name.Value) {
		err := fmt.Errorf("`%v` is built-in function", stmt.Name.Value)
		return nil, function.NewError(err)
	}

	if err := this.s.ExpectPeek(token.ASSIGN); nil != err {
		return nil, function.NewError(err)
	}

	this.s.NextToken()

	expr, err := this.p.ParseExpression(PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if fn, err := expr.AsFunction(); nil == err {
		fn.Lambda = stmt.Name.Value
	}
	stmt.Value = expr

	for !this.s.StmtEnd(endTok) {
		this.s.NextToken()
	}
	return stmt, nil
}

// exprStmt : implement stmtDecoder
type exprStmt struct {
	s scanner
	p Parser
}

func (this *exprStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := ast.NewExpr()
	expr, err := this.p.ParseExpression(PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Expr = expr

	if err := this.s.PeekIs(token.SEMICOLON); nil == err {
		this.s.NextToken()
	}
	return stmt, nil
}

// functionStmt : implement stmtDecoder
type functionStmt struct {
	s scanner
	p Parser
}

func (this *functionStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	this.s.NextToken()
	stmt := this.s.NewFunction()
	fn, err := this.s.ParseFunction(false, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Value = fn
	for !this.s.StmtEnd(endTok) {
		this.s.NextToken()
	}
	return stmt, nil
}
