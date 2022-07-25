package ast

import (
	"bytes"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
	"github.com/jobs-github/Q/token"
)

// DeferStmt : implement Statement
type DeferStmt struct {
	Do Expression
}

func (this *DeferStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeStmtDefer,
		keyValue: this.Do.Encode(),
	}
}
func (this *DeferStmt) Decode(b []byte) error {
	var err error
	this.Do, err = decodeExpr(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *DeferStmt) statementNode() {}

func (this *DeferStmt) String() string {
	var out bytes.Buffer
	out.WriteString(token.Defer)
	out.WriteString(" ")
	out.WriteString(this.Do.String())
	out.WriteString(";")
	return out.String()
}
func (this *DeferStmt) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return object.Nil, nil
}
func (this *DeferStmt) walk(cb func(module string)) {}
func (this *DeferStmt) doDefer(env object.Env) error {
	if _, err := this.Do.Eval(env, false); nil != err {
		return function.NewError(err)
	}
	return nil
}
