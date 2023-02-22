package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// ForExpr : implement Expression
type ForExpr struct {
	defaultNode
	Init Expression
	Cond *Function
	Next *Function
	Loop *Function
	St   Expression // initial state
}

func (this *ForExpr) expressionNode() {}

func (this *ForExpr) Do(v Visitor) error {
	return v.DoFor(this)
}

func (this *ForExpr) value() map[string]interface{} {
	m := map[string]interface{}{
		"init": this.Init.Encode(),
		"cond": this.Cond.Encode(),
		"next": this.Next.Encode(),
		"loop": this.Loop.Encode(),
		"st":   this.St.Encode(),
	}
	return m
}

func (this *ForExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprFor,
		keyValue: this.value(),
	}
}
func (this *ForExpr) Decode(b []byte) error {
	var v struct {
		Init JsonNode `json:"init"`
		Cond JsonNode `json:"cond"`
		Next JsonNode `json:"next"`
		Loop JsonNode `json:"loop"`
		St   JsonNode `json:"st"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Init, err = v.Init.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Cond, err = v.Cond.decodeFn()
	if nil != err {
		return function.NewError(err)
	}
	this.Next, err = v.Next.decodeFn()
	if nil != err {
		return function.NewError(err)
	}
	this.Loop, err = v.Loop.decodeFn()
	if nil != err {
		return function.NewError(err)
	}
	this.St, err = v.St.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ForExpr) String() string {
	// TODO
	var out bytes.Buffer
	out.WriteString("for")
	return out.String()
}
func (this *ForExpr) Eval(e object.Env) (object.Object, error) {
	// TODO
	return nil, nil
}
