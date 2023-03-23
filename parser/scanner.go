package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

var (
	errNoTok = errors.New("no valid token")
)

type scanner interface {
	ParseFunction(lambda bool, p Parser) (*ast.Function, error)
	NewPrefix() *ast.PrefixExpr
	NewInfix(left ast.Expression) *ast.InfixExpr
	NewIndex(left ast.Expression) *ast.IndexExpr
	NewCallMember(left ast.Expression) *ast.CallMember
	NewObjectMember(left ast.Expression) *ast.ObjectMember
	NewCall(left ast.Expression) *ast.Call
	NewConditional(left ast.Expression) *ast.ConditionalExpr
	NewFunction() *ast.FunctionStmt
	NewBoolean() *ast.Boolean
	NewInteger() (*ast.Integer, error)
	NewString() *ast.String

	Clone() scanner
	String() string
	StmtEnd(endTok token.TokenType) bool
	GetIdentifier() *ast.Identifier
	Eof() bool
	PeekPrecedence() int
	CurPrecedence() int
	CurTokenType() token.TokenType
	PeekTokenType() token.TokenType
	NextToken()
	ExpectPeek(t token.TokenType) error
	PeekIs(t token.TokenType) error
	Peek2Is(t token.TokenType) error
	CurrentIs(t token.TokenType) error
	ExpectPeek2(t1 token.TokenType, t2 token.TokenType) bool
	ExpectCur2(cur token.TokenType, peek token.TokenType) bool
}

// scannerImpl : implement Scanner
type scannerImpl struct {
	toks     []*token.Token
	pos      int
	curTok   *token.Token
	peekTok  *token.Token
	peekTok2 *token.Token
}

func newScanner(l Lexer) (scanner, error) {
	toks, err := l.Parse()
	if nil != err {
		return nil, function.NewError(err)
	}
	if nil == toks || len(toks) < 1 {
		return nil, function.NewError(errNoTok)
	}
	s := &scannerImpl{toks: toks, pos: 0}
	s.curTok = toks[0]
	sz := len(toks)
	if sz == 1 {
		s.peekTok = toks[0]
		s.peekTok2 = toks[0]
	} else if sz == 2 {
		s.peekTok = toks[1]
		s.peekTok2 = toks[1]
	} else {
		s.peekTok = toks[1]
		s.peekTok2 = toks[2]
	}
	return s, nil
}

func (this *scannerImpl) fnName(lambda bool) string {
	if lambda {
		return ""
	} else {
		return this.curTok.Literal
	}
}

func (this *scannerImpl) ParseFunction(lambda bool, p Parser) (*ast.Function, error) {
	fn := &ast.Function{Name: this.fnName(lambda)}
	if err := this.ExpectPeek(token.LPAREN); nil != err {
		return nil, function.NewError(err)
	}
	args, err := this.parseArgs()
	if nil != err {
		return nil, function.NewError(err)
	}
	fn.Args = args
	if err := this.ExpectPeek(token.LBRACE); nil != err {
		return nil, function.NewError(err)
	}
	stmt, err := p.ParseBlockStmt()
	if nil != err {
		return nil, function.NewError(err)
	}
	fn.Body = stmt
	return fn, nil
}

func (this *scannerImpl) parseArgs() (ast.IdentifierSlice, error) {
	args := ast.IdentifierSlice{}
	if err := this.PeekIs(token.RPAREN); nil == err {
		this.NextToken()
		return args, nil
	}
	this.NextToken()
	ident := this.GetIdentifier()
	args = append(args, ident)

	for nil == this.PeekIs(token.COMMA) {
		this.NextToken()
		this.NextToken()
		ident := this.GetIdentifier()
		args = append(args, ident)
	}

	if err := this.ExpectPeek(token.RPAREN); nil != err {
		return nil, function.NewError(err)
	}
	return args, nil
}

func (this *scannerImpl) NewPrefix() *ast.PrefixExpr {
	return &ast.PrefixExpr{Op: this.curTok}
}

func (this *scannerImpl) NewInfix(left ast.Expression) *ast.InfixExpr {
	return &ast.InfixExpr{
		Op:   this.curTok,
		Left: left,
	}
}

func (this *scannerImpl) NewIndex(left ast.Expression) *ast.IndexExpr {
	return &ast.IndexExpr{Left: left}
}

func (this *scannerImpl) NewCallMember(left ast.Expression) *ast.CallMember {
	return &ast.CallMember{Left: left}
}

func (this *scannerImpl) NewObjectMember(left ast.Expression) *ast.ObjectMember {
	return &ast.ObjectMember{Left: left}
}

func (this *scannerImpl) NewCall(left ast.Expression) *ast.Call {
	return &ast.Call{Func: left}
}

func (this *scannerImpl) NewConditional(left ast.Expression) *ast.ConditionalExpr {
	return &ast.ConditionalExpr{Cond: left}
}

func (this *scannerImpl) NewFunction() *ast.FunctionStmt {
	return &ast.FunctionStmt{Name: &ast.Identifier{Value: this.curTok.Literal}}
}

func (this *scannerImpl) NewBoolean() *ast.Boolean {
	return &ast.Boolean{Value: this.curTok.TypeIs(token.TRUE)}
}

func (this *scannerImpl) NewInteger() (*ast.Integer, error) {
	expr := &ast.Integer{}
	val, err := strconv.ParseInt(this.curTok.Literal, 0, 64)
	if nil != err {
		err := fmt.Errorf("could not parse %v as integer", this.curTok.Literal)
		return nil, function.NewError(err)
	}
	expr.Value = val
	return expr, nil
}

func (this *scannerImpl) NewString() *ast.String {
	return &ast.String{Value: this.curTok.Literal}
}

func (this *scannerImpl) Clone() scanner {
	return &scannerImpl{
		toks:     this.toks,
		pos:      this.pos,
		curTok:   this.curTok,
		peekTok:  this.peekTok,
		peekTok2: this.peekTok2,
	}
}

func (this *scannerImpl) String() string {
	return fmt.Sprintf("`%v %v %v`", this.curTok.Literal, this.peekTok.Literal, this.peekTok2.Literal)
}

func (this *scannerImpl) StmtEnd(endTok token.TokenType) bool {
	return this.curTok.TypeIs(endTok) || this.curTok.Eof()
}

func (this *scannerImpl) GetIdentifier() *ast.Identifier {
	return &ast.Identifier{Value: this.curTok.Literal}
}

func (this *scannerImpl) Eof() bool {
	return this.curTok.Eof()
}

func (this *scannerImpl) PeekPrecedence() int {
	return getPrecedence(this.peekTok)
}

func (this *scannerImpl) CurPrecedence() int {
	return getPrecedence(this.curTok)
}

func (this *scannerImpl) CurTokenType() token.TokenType {
	return this.curTok.Type
}

func (this *scannerImpl) PeekTokenType() token.TokenType {
	return this.peekTok.Type
}

func (this *scannerImpl) NextToken() {
	if this.Eof() {
		return
	}
	this.pos++
	this.curTok = this.toks[this.pos]
	if !this.curTok.Eof() {
		this.peekTok = this.toks[this.pos+1]
		if !this.peekTok.Eof() {
			this.peekTok2 = this.toks[this.pos+2]
		}
	}
}

func (this *scannerImpl) ExpectPeek(t token.TokenType) error {
	if this.peekTok.TypeIs(t) {
		this.NextToken()
		return nil
	}
	err := fmt.Errorf("expected next token to be %v, got %v instead", token.ToString(t), token.ToString(this.peekTok.Type))
	return function.NewError(err)
}

func (this *scannerImpl) PeekIs(t token.TokenType) error {
	if !this.peekTok.TypeIs(t) {
		err := fmt.Errorf("expected peek token to be %v, got %v instead", token.ToString(t), token.ToString(this.peekTok.Type))
		return function.NewError(err)
	}
	return nil
}

func (this *scannerImpl) Peek2Is(t token.TokenType) error {
	if !this.peekTok2.TypeIs(t) {
		err := fmt.Errorf("expected peek2 token to be %v, got %v instead", token.ToString(t), token.ToString(this.peekTok2.Type))
		return function.NewError(err)
	}
	return nil
}

func (this *scannerImpl) CurrentIs(t token.TokenType) error {
	if !this.curTok.TypeIs(t) {
		err := fmt.Errorf("expected current token to be %v, got %v instead", token.ToString(t), token.ToString(this.curTok.Type))
		return function.NewError(err)
	}
	return nil
}

func (this *scannerImpl) ExpectPeek2(t1 token.TokenType, t2 token.TokenType) bool {
	return this.peekTok.TypeIs(t1) && this.peekTok2.TypeIs(t2)
}

func (this *scannerImpl) ExpectCur2(cur token.TokenType, peek token.TokenType) bool {
	return this.curTok.TypeIs(cur) && this.peekTok.TypeIs(peek)
}
