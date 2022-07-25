package ast

import (
	"bytes"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// ReturnStmt : implement Statement
type ReturnStmt struct {
	ReturnValue Expression
}

func (this *ReturnStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeStmtReturn,
		keyValue: this.ReturnValue.Encode(),
	}
}
func (this *ReturnStmt) Decode(b []byte) error {
	var err error
	this.ReturnValue, err = decodeExpr(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ReturnStmt) statementNode() {}

func (this *ReturnStmt) String() string {
	var out bytes.Buffer
	out.WriteString(token.Return)
	out.WriteString(" ")

	if nil != this.ReturnValue {
		out.WriteString(this.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
func (this *ReturnStmt) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	val, err := this.ReturnValue.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return object.NewReturn(val), nil
}
func (this *ReturnStmt) walk(cb func(module string))  {}
func (this *ReturnStmt) doDefer(env object.Env) error { return nil }
