package stmt

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/interfaces"
	"github.com/jobs-github/Q/scanner"
	"github.com/jobs-github/Q/token"
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
		scanner:            s,
		p:                  p,
		newParser:          newParser,
		functionDecoder:    &functionStmt{s, p},
		assignDecoder:      &assignStmt{s, p},
		assignIndexDecoder: &assignIndexStmt{s, p},
		exprDecoder:        &exprStmt{s, p},
		deferDecoder:       &deferStmt{s, p},
		m: map[token.TokenType]stmtDecoder{
			token.VAR:    &varStmt{s, p},
			token.RETURN: &returnStmt{s, p},
			token.BREAK:  &breakStmt{s},
		},
	}
}

type stmtDecoder interface {
	decode(endTok token.TokenType) (ast.Statement, error)
}

// stmtParser : implement StmtParser
type stmtParser struct {
	scanner            scanner.Scanner
	p                  interfaces.Parser
	newParser          func(s scanner.Scanner) interfaces.Parser
	functionDecoder    stmtDecoder
	assignDecoder      stmtDecoder
	assignIndexDecoder stmtDecoder
	exprDecoder        stmtDecoder
	deferDecoder       stmtDecoder
	m                  map[token.TokenType]stmtDecoder
}

func (this *stmtParser) Decode(t token.TokenType, endTok token.TokenType) (ast.Statement, error) {
	parser, ok := this.m[t]
	if ok {
		return parser.decode(endTok)
	} else {
		if this.isFunctionStmt() {
			return this.functionDecoder.decode(endTok)
		} else if this.isAssignStmt() {
			return this.assignDecoder.decode(endTok)
		} else if this.isAssignIndexStmt() {
			return this.assignIndexDecoder.decode(endTok)
		} else if this.isDeferStmt() {
			return this.deferDecoder.decode(endTok)
		} else {
			return this.exprDecoder.decode(endTok)
		}
	}
}

func (this *stmtParser) isFunctionStmt() bool {
	return this.scanner.ExpectCur2(token.FUNC, token.IDENT)
}

func (this *stmtParser) isAssignStmt() bool {
	return this.scanner.ExpectCur2(token.IDENT, token.ASSIGN)
}

func (this *stmtParser) isAssignIndexStmt() bool {
	s := this.scanner.Clone()
	parser := &assignIndexStmt{s, this.newParser(s)}
	return parser.match()
}

func (this *stmtParser) isDeferStmt() bool {
	return nil == this.scanner.CurrentIs(token.DEFER)
}
