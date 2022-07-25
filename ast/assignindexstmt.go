package ast

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// AssignIndexStmt : implement Statement
type AssignIndexStmt struct {
	Name  *Identifier
	Idx   Expression
	Value Expression
}

func (this *AssignIndexStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeStmtAssignindex,
		keyValue: map[string]interface{}{
			"name":  this.Name.Encode(),
			"idx":   this.Idx.Encode(),
			"value": this.Value.Encode(),
		},
	}
}

func (this *AssignIndexStmt) Decode(b []byte) error {
	var v struct {
		Name  JsonNode `json:"name"`
		Idx   JsonNode `json:"idx"`
		Value JsonNode `json:"value"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Name, err = v.Name.decodeIdent()
	if nil != err {
		return function.NewError(err)
	}
	this.Idx, err = v.Idx.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Value, err = v.Value.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *AssignIndexStmt) statementNode() {}

func (this *AssignIndexStmt) String() string {
	var out bytes.Buffer
	out.WriteString(this.Name.String())
	out.WriteString("[")
	out.WriteString(this.Idx.String())
	out.WriteString("]")
	out.WriteString(" = ")
	out.WriteString(this.Value.String())
	out.WriteString(";")
	return out.String()
}
func (this *AssignIndexStmt) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	obj, ok := env.Get(this.Name.Value)
	if !ok {
		err := fmt.Errorf("identifier `%v` missing in env, stmt(`%v`)", this.Name.Value, this.String())
		return object.Nil, function.NewError(err)
	}
	idx, err := this.Idx.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	val, err := this.Value.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return obj.CallMember(object.FnSet, object.Objects{idx, val})
}
func (this *AssignIndexStmt) walk(cb func(module string))  {}
func (this *AssignIndexStmt) doDefer(env object.Env) error { return nil }
