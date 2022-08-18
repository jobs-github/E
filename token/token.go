package token

import "fmt"

type TokenType uint

const (
	ILLEGAL TokenType = iota
	EOF

	//literal_beg
	IDENT
	INT
	STRING
	//literal_end

	//operator_beg
	LT        // <
	GT        // >
	ASSIGN    // =
	NOT       // !
	ADD       // +
	SUB       // -
	MUL       // *
	DIV       // /
	MOD       // %
	EQ        // ==
	NEQ       // !=
	LEQ       // <=
	GEQ       // >=
	AND       // &&
	OR        // ||
	COMMA     // ,
	PERIOD    // .
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACK    // [
	RBRACK    // ]
	QUESTION  // ?
	//operator_end

	//keyword_beg
	TRUE
	FALSE
	NULL
	FUNC
	CONST
	//keyword_end
)

const (
	Const = "const"
	Func  = "func"
	Null  = "null"
	True  = "true"
	False = "false"
)

var (
	Add = &Token{ADD, "+"}
	Sub = &Token{SUB, "-"}
	Mul = &Token{MUL, "*"}
	Div = &Token{DIV, "/"}
	Mod = &Token{MOD, "%"}
	Lt  = &Token{LT, "<"}
	Gt  = &Token{GT, ">"}
	Eq  = &Token{EQ, "=="}
	Neq = &Token{NEQ, "!="}
	Leq = &Token{LEQ, "<="}
	Geq = &Token{GEQ, ">="}
	And = &Token{AND, "&&"}
	Or  = &Token{OR, "||"}
)

var (
	infixTokens = map[string]TokenType{
		"<":  LT,
		">":  GT,
		"+":  ADD,
		"-":  SUB,
		"*":  MUL,
		"/":  DIV,
		"%":  MOD,
		"==": EQ,
		"!=": NEQ,
		"<=": LEQ,
		">=": GEQ,
		"&&": AND,
		"||": OR,
		"(":  LPAREN,
		"[":  LBRACK,
		".":  PERIOD,
		"?":  QUESTION,
	}
	tokenTypes = map[byte]TokenType{
		'+': ADD,
		'-': SUB,
		'*': MUL,
		'/': DIV,
		'%': MOD,
		',': COMMA,
		'.': PERIOD,
		';': SEMICOLON,
		':': COLON,
		'(': LPAREN,
		')': RPAREN,
		'{': LBRACE,
		'}': RBRACE,
		'[': LBRACK,
		']': RBRACK,
		'?': QUESTION,
	}
	keywords = map[string]TokenType{
		True:  TRUE,
		False: FALSE,
		Null:  NULL,
		Func:  FUNC,
		Const: CONST,
	}

	tokenTypeStrings = map[TokenType]string{
		EOF:       "EOF",
		IDENT:     "IDENT",
		INT:       "INT",
		STRING:    "STRING",
		LT:        "LT",
		GT:        "GT",
		ASSIGN:    "ASSIGN",
		NOT:       "NOT",
		ADD:       "ADD",
		SUB:       "SUB",
		MUL:       "MUL",
		DIV:       "DIV",
		MOD:       "MOD",
		EQ:        "EQ",
		NEQ:       "NEQ",
		LEQ:       "LEQ",
		GEQ:       "GEQ",
		AND:       "AND",
		OR:        "OR",
		COMMA:     "COMMA",
		PERIOD:    "PERIOD",
		SEMICOLON: "SEMICOLON",
		COLON:     "COLON",
		LPAREN:    "LPAREN",
		RPAREN:    "RPAREN",
		LBRACE:    "LBRACE",
		RBRACE:    "RBRACE",
		LBRACK:    "LBRACK",
		RBRACK:    "RBRACK",
		QUESTION:  "QUESTION",
		TRUE:      "TRUE",
		FALSE:     "FALSE",
		NULL:      "NULL",
		FUNC:      "FUNC",
		CONST:     "CONST",
	}
)

func ToString(t TokenType) string {
	s, ok := tokenTypeStrings[t]
	if ok {
		return s
	}
	return "ILLEGAL"
}

func IsKeywords(ident string) bool {
	_, ok := keywords[ident]
	return ok
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
}

func (this *Token) TypeIs(t TokenType) bool {
	return this.Type == t
}

func (this *Token) Eof() bool {
	return this.TypeIs(EOF)
}

func (this *Token) Illegal() bool {
	return this.TypeIs(ILLEGAL)
}

func (this *Token) String() string {
	return fmt.Sprintf("{\"type\":\"%v\",\"literal\":\"%v\"}", ToString(this.Type), this.Literal)
}

func GetTokenType(ch byte) (TokenType, bool) {
	tt, ok := tokenTypes[ch]
	return tt, ok
}

func GetInfixToken(s string) (*Token, error) {
	tt, ok := infixTokens[s]
	if !ok {
		return nil, fmt.Errorf("not infix token: %v", s)
	}
	return &Token{Type: tt, Literal: s}, nil
}

func Bool(v bool) string {
	if v {
		return True
	} else {
		return False
	}
}
