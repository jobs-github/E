package lexer

import "github.com/jobs-github/Q/token"

type Lexer interface {
	Parse() ([]*token.Token, error)
	nextToken() (*token.Token, error)
}
