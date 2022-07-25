package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// ImportExpr : implement Expression
type ImportExpr struct {
	ModuleName string
	AsKey      string
}

func (this *ImportExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprImport,
		keyValue: map[string]interface{}{
			"as":     this.AsKey,
			"module": this.ModuleName,
		},
	}
}
func (this *ImportExpr) Decode(b []byte) error {
	var v struct {
		AsKey      string `json:"as"`
		ModuleName string `json:"module"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.AsKey = v.AsKey
	this.ModuleName = v.ModuleName
	return nil
}
func (this *ImportExpr) expressionNode() {}

func (this *ImportExpr) String() string {
	var out bytes.Buffer
	out.WriteString("import ")
	if this.AsKey != this.ModuleName {
		out.WriteString(this.AsKey)
		out.WriteString(" ")
	}
	out.WriteString("\"")
	out.WriteString(this.ModuleName)
	out.WriteString("\"")
	return out.String()
}
func (this *ImportExpr) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	if err := env.Import(this.ModuleName, this.AsKey); nil != err {
		return object.Nil, function.NewError(err)
	}
	return object.Nil, nil
}
func (this *ImportExpr) walk(cb func(module string)) {
	cb(this.ModuleName)
}
func (this *ImportExpr) doDefer(env object.Env) error { return nil }
