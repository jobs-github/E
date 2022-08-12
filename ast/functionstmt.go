package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// FunctionStmt : implement Statement
type FunctionStmt struct {
	Name  *Identifier
	Value *Function
}

func (this *FunctionStmt) Do(v Visitor) error {
	return v.DoFunction(this)
}

func (this *FunctionStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeStmtFn,
		keyValue: map[string]interface{}{
			"name":  this.Name.Encode(),
			"value": this.Value.Encode(),
		},
	}
}
func (this *FunctionStmt) Decode(b []byte) error {
	var v struct {
		Name  JsonNode `json:"name"`
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
	this.Value, err = v.Value.decodeFn()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *FunctionStmt) statementNode() {}

func (this *FunctionStmt) String() string {
	var out bytes.Buffer
	out.WriteString("func ")
	out.WriteString(this.Name.String())
	out.WriteString(this.Value.String())
	out.WriteString(";")
	return out.String()
}
func (this *FunctionStmt) Eval(env object.Env) (object.Object, error) {
	return evalVar(env, this.Name, this.Value)
}
