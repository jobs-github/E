package ast

import (
	"bytes"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// VarStmt : implement Statement
type VarStmt struct {
	Name  *Identifier
	Value Expression
}

func (this *VarStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeStmtVar,
		keyValue: map[string]interface{}{
			"name":  this.Name.Encode(),
			"value": this.Value.Encode(),
		},
	}
}
func (this *VarStmt) Decode(b []byte) error {
	var err error
	this.Name, this.Value, err = decodeKv(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *VarStmt) statementNode() {}

func (this *VarStmt) String() string {
	var out bytes.Buffer
	out.WriteString(token.Var)
	out.WriteString(" ")
	out.WriteString(this.Name.String())
	out.WriteString(" = ")
	if nil != this.Value {
		out.WriteString(this.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (this *VarStmt) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return evalVar(env, insideLoop, this.Name, this.Value)
}
func (this *VarStmt) walk(cb func(module string)) {
	this.Value.walk(cb)
}
func (this *VarStmt) doDefer(env object.Env) error { return nil }
