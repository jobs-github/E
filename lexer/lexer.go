package lexer

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

type Lexer interface {
	Parse() ([]*token.Token, error)
	nextToken() (*token.Token, error)
}

// lexer : implement Lexer
type lexer struct {
	input        string
	position     int
	nextPosition int
	ch           byte
}

func New(input string) Lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, ch byte) *token.Token {
	return &token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '_' == ch
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isWhitespace(c byte) bool {
	return c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\v' ||
		c == '\f' ||
		c == '\r'
}

func (this *lexer) startofComment() bool {
	return this.ch == '/' && this.peekChar() == '*'
}

func (this *lexer) endofComment() bool {
	return this.ch == '*' && this.peekChar() == '/'
}

func (this *lexer) skipWhitespace() {
	this.readChar()

	for isWhitespace(this.ch) {
		this.readChar()

		if this.eof() {
			break
		}
	}
}

func (this *lexer) skipComment() {
	this.readChar()
	this.readChar()

	if this.eof() {
		return
	}

	for !this.endofComment() {
		this.readChar()

		if this.eof() {
			return
		}
	}
	this.readChar()
	this.readChar()
}

func (this *lexer) skip() {
	for {
		if isWhitespace(this.ch) {
			this.skipWhitespace()
		} else if this.startofComment() {
			this.skipComment()
		} else {
			break
		}
	}
}

func (this *lexer) twoCharToken(tokenType token.TokenType, expectedNextChar byte, tokenType2 token.TokenType, literal string) *token.Token {
	if expectedNextChar == this.peekChar() {
		this.readChar()
		return &token.Token{Type: tokenType2, Literal: literal}
	} else {
		return newToken(tokenType, this.ch)
	}
}

func (this *lexer) Parse() ([]*token.Token, error) {
	toks := []*token.Token{}
	for {
		tok, err := this.nextToken()
		if nil != err {
			return nil, function.NewError(err)
		}
		toks = append(toks, tok)
		if tok.Eof() {
			break
		}
	}
	return toks, nil
}

func (this *lexer) eof() bool {
	return 0 == this.ch
}

func (this *lexer) nextToken() (*token.Token, error) {
	var tok *token.Token
	this.skip()

	if this.eof() {
		return &token.Token{Type: token.EOF, Literal: ""}, nil
	}

	switch this.ch {
	case '"':
		if s, err := this.readString(); nil != err {
			return newToken(token.ILLEGAL, this.ch), function.NewError(err)
		} else {
			tok = &token.Token{Type: token.STRING, Literal: s}
		}
	case '&':
		tok = this.twoCharToken(token.ILLEGAL, '&', token.AND, "&&")
	case '|':
		tok = this.twoCharToken(token.ILLEGAL, '|', token.OR, "||")
	case '=':
		tok = this.twoCharToken(token.ASSIGN, '=', token.EQ, "==")
	case '!':
		tok = this.twoCharToken(token.NOT, '=', token.NEQ, "!=")
	case '<':
		tok = this.twoCharToken(token.LT, '=', token.LEQ, "<=")
	case '>':
		tok = this.twoCharToken(token.GT, '=', token.GEQ, ">=")
	default:
		tt, ok := token.GetTokenType(this.ch)
		if ok {
			tok = newToken(tt, this.ch)
		} else {
			if isLetter(this.ch) {
				literal := this.readIdentifier()
				return &token.Token{Type: token.LookupIdent(literal), Literal: literal}, nil
			} else if isDigit(this.ch) {
				return &token.Token{Type: token.INT, Literal: this.readNumber()}, nil
			} else {
				tok = newToken(token.ILLEGAL, this.ch)
			}
		}
	}

	this.readChar()
	return tok, nil
}

func (this *lexer) readNumber() string {
	pos := this.position
	for isDigit(this.ch) {
		this.readChar()
	}
	return this.input[pos:this.position]
}

func (this *lexer) readIdentifier() string {
	pos := this.position
	for isLetter(this.ch) || isDigit(this.ch) {
		this.readChar()
	}
	return this.input[pos:this.position]
}

func (this *lexer) checkEscape(ch byte) bool {
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
		return true
	}
	return false
}

func (this *lexer) readString() (string, error) {
	start := this.position + 1
	for {
		this.readChar()
		if this.ch == '"' {
			break
		}
		if this.ch == '\\' {
			this.readChar()
			if !this.checkEscape(this.ch) {
				err := fmt.Errorf("unexpected escape, pos: %v, %v", this.position, this.input[this.position:])
				return "", function.NewError(err)
			}
		}
	}
	return this.input[start:this.position], nil
}

func (this *lexer) peekChar() byte {
	if this.nextPosition >= len(this.input) {
		return 0
	} else {
		return this.input[this.nextPosition]
	}
}

func (this *lexer) readChar() {
	if this.nextPosition >= len(this.input) {
		this.ch = 0
	} else {
		this.ch = this.input[this.nextPosition]
	}
	this.position = this.nextPosition
	this.nextPosition = this.nextPosition + 1
}
