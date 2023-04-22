package parser

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

// public
type Parser interface {
	ParseProgram() (ast.Node, error)

	ParseInfixExpression(left ast.Expression) (ast.Expression, error)
	ParseCallExpression(left ast.Expression) (ast.Expression, error)
	ParseIndexExpression(left ast.Expression) (ast.Expression, error)
	ParseMemberExpression(left ast.Expression) (ast.Expression, error)
	ParseConditionalExpression(left ast.Expression) (ast.Expression, error)

	ParseStmt(endTok token.TokenType) (ast.Statement, error)
	ParseBlockStmt() (*ast.BlockStmt, error)
	ParseExpression(precedence int) (ast.Expression, error)
	ParseExpressions(endTok token.TokenType) (ast.ExpressionSlice, error)
}

func New(code string) (Parser, error) {
	l := newLexer(code)
	s, err := newScanner(l)
	if nil == s {
		return nil, err
	}
	p := &parserImpl{s: s}
	p.sp = newStmtParser(s, p)
	p.ep = newExprParser(s, p)
	p.im = newInfixDecoders(p)
	return p, nil
}

type infixDecoderMap map[token.TokenType](func(ast.Expression) (ast.Expression, error))

// parserImpl : implement Parser
type parserImpl struct {
	s  scanner
	sp stmtParser
	ep exprParser
	im infixDecoderMap
}

func newInfixDecoders(p Parser) infixDecoderMap {
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

func (this *parserImpl) ParseStmt(endTok token.TokenType) (ast.Statement, error) {
	return this.sp.Decode(this.s.CurTokenType(), endTok)
}

func (this *parserImpl) ParseProgram() (ast.Node, error) {
	program := &ast.Program{Stmts: ast.StatementSlice{}}
	for !this.s.Eof() {
		// need to skip ;
		if nil == this.s.CurrentIs(token.SEMICOLON) {
			this.s.NextToken()
			continue
		}

		stmt, err := this.ParseStmt(token.SEMICOLON)
		if nil != err {
			return nil, function.NewError(err)
		}
		program.Stmts = append(program.Stmts, stmt)
		this.s.NextToken()
	}
	return program, nil
}

func (this *parserImpl) ParseBlockStmt() (*ast.BlockStmt, error) {
	block := ast.NewBlock()
	this.s.NextToken()
	stmt, err := this.ParseStmt(token.SEMICOLON)
	if nil != err {
		return nil, function.NewError(err)
	}
	block.Stmt = stmt
	if err := this.s.ExpectPeek(token.RBRACE); nil != err {
		return nil, function.NewError(err)
	}
	return block, nil
}

func (this *parserImpl) ParseExpression(precedence int) (ast.Expression, error) {
	leftExpr, err := this.ep.Decode(this.s.CurTokenType())
	if nil != err {
		return nil, function.NewError(err)
	}

	for nil != this.s.PeekIs(token.SEMICOLON) && precedence < this.s.PeekPrecedence() {
		fn, ok := this.im[this.s.PeekTokenType()]
		if !ok {
			return leftExpr, nil
		}
		this.s.NextToken()
		expr, err := fn(leftExpr)
		if nil != err {
			return nil, function.NewError(err)
		}
		leftExpr = expr
	}
	return leftExpr, nil
}

func (this *parserImpl) ParseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.s.NewInfix(left)
	preced := this.s.CurPrecedence()
	this.s.NextToken()
	right, err := this.ParseExpression(preced)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Right = right
	return expr, nil
}

func (this *parserImpl) ParseCallExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.s.NewCall(left)
	args, err := this.ParseExpressions(token.RPAREN)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Args = args
	return expr, nil
}

func (this *parserImpl) ParseIndexExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.s.NewIndex(left)
	this.s.NextToken()
	idx, err := this.ParseExpression(PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.s.ExpectPeek(token.RBRACK); nil != err {
		return nil, function.NewError(err)
	}
	expr.Index = idx
	return expr, nil
}

func (this *parserImpl) isCallMemeber() bool {
	return this.s.ExpectPeek2(token.IDENT, token.LPAREN)
}

func (this *parserImpl) isObjectMember() bool {
	return nil == this.s.PeekIs(token.IDENT) && nil != this.s.Peek2Is(token.LPAREN)
}

func (this *parserImpl) parseCallMemberExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.s.NewCallMember(left)
	this.s.NextToken()

	expr.Func = this.s.GetIdentifier()
	this.s.NextToken()
	args, err := this.ParseExpressions(token.RPAREN)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Args = args
	return expr, nil
}

func (this *parserImpl) parseObjectMemberExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.s.NewObjectMember(left)
	this.s.NextToken()
	expr.Member = this.s.GetIdentifier()
	return expr, nil
}

func (this *parserImpl) ParseMemberExpression(left ast.Expression) (ast.Expression, error) {
	if this.isCallMemeber() {
		return this.parseCallMemberExpression(left)
	} else if this.isObjectMember() {
		return this.parseObjectMemberExpression(left)
	} else {
		err := fmt.Errorf("unknown pattern, %v", this.s.String())
		return nil, function.NewError(err)
	}
}

func (this *parserImpl) ParseConditionalExpression(left ast.Expression) (ast.Expression, error) {
	expr := this.s.NewConditional(left)
	this.s.NextToken()
	yes, err := this.ParseExpression(PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	if err := this.s.ExpectPeek(token.COLON); nil != err {
		return nil, function.NewError(err)
	}
	this.s.NextToken()
	no, err := this.ParseExpression(PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Yes = yes
	expr.No = no
	return expr, nil
}

func (this *parserImpl) ParseExpressions(endTok token.TokenType) (ast.ExpressionSlice, error) {
	args := ast.ExpressionSlice{}
	if nil == this.s.PeekIs(endTok) {
		this.s.NextToken()
		return args, nil
	}
	this.s.NextToken()
	expr, err := this.ParseExpression(PRECED_LOWEST)
	if nil != err {
		return nil, function.NewError(err)
	}
	args = append(args, expr)

	for nil == this.s.PeekIs(token.COMMA) {
		this.s.NextToken()
		this.s.NextToken()

		expr, err := this.ParseExpression(PRECED_LOWEST)
		if nil != err {
			return nil, function.NewError(err)
		}
		args = append(args, expr)
	}
	if err := this.s.ExpectPeek(endTok); nil != err {
		return nil, function.NewError(err)
	}
	return args, nil
}
