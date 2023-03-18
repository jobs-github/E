package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// ReduceExpr : implement Expression
type ReduceExpr struct {
	defaultNode
	Arr  Expression // Array or Identifier
	Body Expression // Function or Identifier
	Init Expression
}

func (this *ReduceExpr) expressionNode() {}

func (this *ReduceExpr) Do(v Visitor) error {
	return v.DoReduce(this)
}

func (this *ReduceExpr) value() map[string]interface{} {
	m := map[string]interface{}{
		"arr":  this.Arr.Encode(),
		"body": this.Body.Encode(),
		"init": this.Init.Encode(),
	}
	return m
}

func (this *ReduceExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprReduce,
		keyValue: this.value(),
	}
}
func (this *ReduceExpr) Decode(b []byte) error {
	var v struct {
		Arr  JsonNode `json:"arr"`
		Body JsonNode `json:"body"`
		Init JsonNode `json:"init"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Arr, err = v.Arr.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Body, err = v.Body.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Init, err = v.Init.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ReduceExpr) String() string {
	var out bytes.Buffer
	out.WriteString("reduce(")
	out.WriteString(this.Arr.String())
	out.WriteString(",")
	out.WriteString(this.Body.String())
	out.WriteString(",")
	out.WriteString(this.Init.String())
	out.WriteString(")")
	return out.String()
}
func (this *ReduceExpr) Eval(e object.Env) (object.Object, error) {
	v, err := this.Arr.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	arr, err := v.AsArray()
	if nil != err {
		return object.Nil, err
	}
	acc, err := this.Init.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	cb, err := this.Body.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	if !object.Callable(cb) {
		return object.Nil, errNotCallable
	}
	if arr.Items == nil || len(arr.Items) < 1 {
		return acc, nil
	}
	for _, item := range arr.Items {
		v, err := cb.Call(object.Objects{acc, item})
		if nil != err {
			return object.Nil, err
		}
		acc = v
	}
	return acc, nil
}
