package lexer

import "github.com/jobs-github/escript/token"

type Lexer interface {
	Parse() ([]*token.Token, error)
	nextToken() (*token.Token, error)
}
