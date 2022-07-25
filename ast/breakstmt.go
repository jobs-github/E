package ast

import (
	"bytes"

	"github.com/jobs-github/Q/object"
	"github.com/jobs-github/Q/token"
)

// BreakStmt : implement Statement
type BreakStmt struct {
}

func (this *BreakStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeStmtBreak,
	}
}
func (this *BreakStmt) Decode(b []byte) error {
	return nil
}
func (this *BreakStmt) statementNode() {}

func (this *BreakStmt) String() string {
	var out bytes.Buffer
	out.WriteString(token.Break)
	out.WriteString(";")
	return out.String()
}
func (this *BreakStmt) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return object.NewBreak(), nil
}
func (this *BreakStmt) walk(cb func(module string))  {}
func (this *BreakStmt) doDefer(env object.Env) error { return nil }
