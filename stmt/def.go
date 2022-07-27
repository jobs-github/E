package stmt

import (
	"github.com/jobs-github/escript/ast"
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
			token.VAR: &varStmt{s, p},
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
