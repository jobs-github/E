package ast

import (
	"bytes"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// AssignStmt : implement Statement
type AssignStmt struct {
	Name  *Identifier
	Value Expression
}

func (this *AssignStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeStmtAssign,
		keyValue: map[string]interface{}{
			"name":  this.Name.Encode(),
			"value": this.Value.Encode(),
		},
	}
}

func (this *AssignStmt) Decode(b []byte) error {
	var err error
	this.Name, this.Value, err = decodeKv(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *AssignStmt) statementNode() {}

func (this *AssignStmt) String() string {
	var out bytes.Buffer
	out.WriteString(this.Name.String())
	out.WriteString(" = ")
	if nil != this.Value {
		out.WriteString(this.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (this *AssignStmt) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	val, err := this.Value.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	if err := env.Assign(this.Name.Value, val); nil != err {
		return object.Nil, function.NewError(err)
	}
	return val, nil
}
func (this *AssignStmt) walk(cb func(module string))  {}
func (this *AssignStmt) doDefer(env object.Env) error { return nil }
