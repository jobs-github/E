package scanner

import "github.com/jobs-github/Q/token"

const (
	_ int = iota
	PRECED_LOWEST
	PRECED_QUESTION // ?
	PRECED_OR       // ||
	PRECED_AND      // &&
	PRECED_EQ       // ==
	PRECED_NEQ      // !=
	PRECED_LT       // < > >= <=
	PRECED_ADD      // +
	PRECED_MUL      // *
	PRECED_PREFIX   // -x !x
	PRECED_CALL     // myFn(x)
	PRECED_HIGHEST
	PRECED_INDEX  = PRECED_HIGHEST // array[index]
	PRECED_PERIOD = PRECED_HIGHEST // s.len()
)

var (
	precedences = map[token.TokenType]int{
		token.LT:       PRECED_LT,
		token.GT:       PRECED_LT,
		token.ADD:      PRECED_ADD,
		token.SUB:      PRECED_ADD,
		token.MUL:      PRECED_MUL,
		token.DIV:      PRECED_MUL,
		token.MOD:      PRECED_MUL,
		token.EQ:       PRECED_EQ,
		token.NEQ:      PRECED_NEQ,
		token.LEQ:      PRECED_LT,
		token.GEQ:      PRECED_LT,
		token.AND:      PRECED_AND,
		token.OR:       PRECED_OR,
		token.LPAREN:   PRECED_CALL,
		token.LBRACK:   PRECED_INDEX,
		token.PERIOD:   PRECED_PERIOD,
		token.QUESTION: PRECED_QUESTION,
	}
)

func getPrecedence(tok *token.Token) int {
	if v, ok := precedences[tok.Type]; ok {
		return v
	}
	return PRECED_LOWEST
}
