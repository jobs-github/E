package parser

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/expr"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/lexer"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/stmt"
	"github.com/jobs-github/escript/token"
)

type infixDecoderMap map[token.TokenType](func(ast.Expression) (ast.Expression, error))

// parserImpl : implement interfaces.Parser
type parserImpl struct {
	scanner       scanner.Scanner
	stmtParser    stmt.StmtParser
	exprParser    expr.ExprParser
	infixDecoders infixDecoderMap
}

func newInfixDecoders(p interfaces.Parser) infixDecoderMap {
	return infixDecoderMap{
		token.LT:       p.ParseInfixExpression,
		token.GT:       p.ParseInfixExpression,
		token.ADD:      p.ParseInfixExpression,
		token.SUB:      p.ParseInfixExpression,
		token.MUL:      p.ParseInfixExpression,
		token.DIV:      p.ParseInfixExpression,
		token.MOD:      p.ParseInfixExpression,
		token.EQ:       p.ParseInfixExpression,
		token.NEQ:      p.ParseInfixExpression,
		token.LEQ:      p.ParseInfixExpression,
		token.GEQ:      p.ParseInfixExpression,
		token.AND:      p.ParseInfixExpression,
		token.OR:       p.ParseInfixExpression,
		token.LPAREN:   p.ParseCallExpression,
		token.LBRACK:   p.ParseIndexExpression,
		token.PERIOD:   p.ParseMemberExpression,
		token.QUESTION: p.ParseConditionalExpression,
	}
}

func newParser(s scanner.Scanner) interfaces.Parser {
	p := &parserImpl{scanner: s}
	p.stmtParser = stmt.NewStmtParser(s, p, newParser)
	p.exprParser = expr.NewExprParser(s, p)
	p.infixDecoders = newInfixDecoders(p)
	return p
}

func New(l lexer.Lexer) (interfaces.Parser, error) {
	s, err := scanner.New(l)
	if nil == s {
		return nil, err
	}
	return newParser(s), nil
}

func (this *parserImpl) ParseStmt(endTok token.TokenType) (ast.Statement, error) {
	return this.stmtParser.Decode(this.scanner.CurTokenType(), endTok)
}

func (this *parserImpl) ParseProgram() (ast.Node, error) {
	program := &ast.Program{Stmts: ast.StatementSlice{}}
	for !this.scanner.Eof() {
		// need to skip ;
		if nil == this.scanner.CurrentIs(token.SEMICOLON) {
			this.scanner.NextToken()
			continue
		}

		stmt, err := this.ParseStmt(token.SEMICOLON)
		if nil != err {
			return nil, function.NewError(err)
		}
		program.Stmts = append(program.Stmts, stmt)
		this.scanner.NextToken()
	}
	return program, nil
}

func (this *parserImpl) ParseBlockStmt() (*ast.BlockStmt, error) {
	block := ast.NewBlock()
	this.scanner.NextToken()
	stmt, err := this.ParseStmt(token.SEMICOLON)
	if nil != err {
		return nil, function.NewError(err)
	}
	block.Stmt = stmt
	if err := this.scanner.ExpectPeek(token.RBRACE); nil != err {
		return nil, function.NewError(err)
	}
	return block, nil
}

func (this *parserImpl) ParseExpression(precedence int) (ast.Expression, error) {
	leftExpr, err := this.exprParser.Decode(this.scanner.CurTokenType())
	if nil != err {
		return nil, function.NewError(err)
	}

	for nil != this.scanner.PeekIs(token.SEMICOLON) && precedence < this.scanner.PeekPrecedence() {
		fn, ok := this.infixDecoders[this.scanner.PeekTokenType()]
		if !ok {
			return leftExpr, nil
		}
		this.scanner.NextToken()
		expr, err := fn(leftExpr)
		if nil != err {
			return nil, function.NewError(err)
		}
		leftExpr = expr
	}
	return leftExpr, nil
}

func (this *parserImpl) ParseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.scanner.NewInfix(left)
	preced := this.scanner.CurPrecedence()
	this.scanner.NextToken()
	right, err := this.ParseExpression(preced)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Right = right
	return expr, nil
}

func (this *parserImpl) ParseCallExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.scanner.NewCall(left)
	args, err := this.ParseExpressions(token.RPAREN)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Args = args
	return expr, nil
}

func (this *parserImpl) ParseIndexExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.scanner.NewIndex(left)
	this.scanner.NextToken()
	idx, err := this.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.RBRACK); nil != err {
		return nil, function.NewError(err)
	}
	expr.Index = idx
	return expr, nil
}

func (this *parserImpl) isCallMemeber() bool {
	return this.scanner.ExpectPeek2(token.IDENT, token.LPAREN)
}

func (this *parserImpl) isObjectMember() bool {
	return nil == this.scanner.PeekIs(token.IDENT) && nil != this.scanner.Peek2Is(token.LPAREN)
}

func (this *parserImpl) parseCallMemberExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.scanner.NewCallMember(left)
	this.scanner.NextToken()

	expr.Func = this.scanner.GetIdentifier()
	this.scanner.NextToken()
	args, err := this.ParseExpressions(token.RPAREN)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Args = args
	return expr, nil
}

func (this *parserImpl) parseObjectMemberExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.scanner.NewObjectMember(left)
	this.scanner.NextToken()
	expr.Member = this.scanner.GetIdentifier()
	return expr, nil
}

func (this *parserImpl) ParseMemberExpression(left ast.Expression) (ast.Expression, error) {
	if this.isCallMemeber() {
		return this.parseCallMemberExpression(left)
	} else if this.isObjectMember() {
		return this.parseObjectMemberExpression(left)
	} else {
		err := fmt.Errorf("unknown pattern, %v", this.scanner.String())
		return nil, function.NewError(err)
	}
}

func (this *parserImpl) ParseConditionalExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.scanner.NewConditional(left)
	this.scanner.NextToken()
	yes, err := this.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.scanner.ExpectPeek(token.COLON); nil != err {
		return nil, function.NewError(err)
	}
	this.scanner.NextToken()
	no, err := this.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Yes = yes
	expr.No = no
	return expr, nil
}

func (this *parserImpl) ParseExpressions(endTok token.TokenType) (ast.ExpressionSlice, error) {
	args := ast.ExpressionSlice{}
	if nil == this.scanner.PeekIs(endTok) {
		this.scanner.NextToken()
		return args, nil
	}
	this.scanner.NextToken()
	expr, err := this.ParseExpression(scanner.PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	args = append(args, expr)

	for nil == this.scanner.PeekIs(token.COMMA) {
		this.scanner.NextToken()
		this.scanner.NextToken()

		expr, err := this.ParseExpression(scanner.PRECED_LOWEST)
		if nil != err {
			return nil, function.NewError(err)
		}
		args = append(args, expr)
	}
	if err := this.scanner.ExpectPeek(endTok); nil != err {
		return nil, function.NewError(err)
	}
	return args, nil
}
