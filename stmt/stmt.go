package stmt

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/builtin"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

type StmtParser interface {
	Decode(t token.TokenType, endTok token.TokenType) (ast.Statement, error)
}

func NewStmtParser(
	s scanner.Scanner,
	p interfaces.Parser,
	newParser func(s scanner.Scanner) interfaces.Parser,
) StmtParser {
	return &stmtParser{
		scanner:         s,
		p:               p,
		newParser:       newParser,
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

// stmtParser : implement StmtParser
type stmtParser struct {
	scanner         scanner.Scanner
	p               interfaces.Parser
	newParser       func(s scanner.Scanner) interfaces.Parser
	functionDecoder stmtDecoder
	exprDecoder     stmtDecoder
	m               map[token.TokenType]stmtDecoder
}

func (this *stmtParser) Decode(t token.TokenType, endTok token.TokenType) (ast.Statement, error) {
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

func (this *stmtParser) isFunctionStmt() bool {
	return this.scanner.ExpectCur2(token.FUNC, token.IDENT)
}

// constStmt : implement stmtDecoder
type constStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *constStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := ast.NewConst()
	if err := this.scanner.ExpectPeek(token.IDENT); nil != err {
		return nil, function.NewError(err)
	}
	stmt.Name = this.scanner.GetIdentifier()

	if builtin.IsBuiltin(stmt.Name.Value) {
		err := fmt.Errorf("`%v` is built-in function", stmt.Name.Value)
		return nil, function.NewError(err)
	}

	if err := this.scanner.ExpectPeek(token.ASSIGN); nil != err {
		return nil, function.NewError(err)
	}

	this.scanner.NextToken()

	expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if fn, err := expr.AsFunction(); nil == err {
		fn.Lambda = stmt.Name.Value
	}
	stmt.Value = expr

	for !this.scanner.StmtEnd(endTok) {
		this.scanner.NextToken()
	}
	return stmt, nil
}

// exprStmt : implement stmtDecoder
type exprStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *exprStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	stmt := ast.NewExpr()
	expr, err := this.p.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Expr = expr

	if err := this.scanner.PeekIs(token.SEMICOLON); nil == err {
		this.scanner.NextToken()
	}
	return stmt, nil
}

// functionStmt : implement stmtDecoder
type functionStmt struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *functionStmt) decode(endTok token.TokenType) (ast.Statement, error) {
	this.scanner.NextToken()
	stmt := this.scanner.NewFunction()
	fn, err := this.scanner.ParseFunction(false, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	stmt.Value = fn
	for !this.scanner.StmtEnd(endTok) {
		this.scanner.NextToken()
	}
	return stmt, nil
}
