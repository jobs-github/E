package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// MapExpr : implement Expression
type MapExpr struct {
	defaultNode
	Arr  Expression // Array or Identifier
	Body Expression // Function or Identifier
}

func (this *MapExpr) expressionNode() {}

func (this *MapExpr) Do(v Visitor) error {
	return v.DoMap(this)
}

func (this *MapExpr) value() map[string]interface{} {
	m := map[string]interface{}{
		"arr":  this.Arr.Encode(),
		"body": this.Body.Encode(),
	}
	return m
}

func (this *MapExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprMap,
		keyValue: this.value(),
	}
}
func (this *MapExpr) Decode(b []byte) error {
	var v struct {
		Arr  JsonNode `json:"arr"`
		Body JsonNode `json:"body"`
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
	return nil
}
func (this *MapExpr) String() string {
	var out bytes.Buffer
	out.WriteString("map(")
	out.WriteString(this.Arr.String())
	out.WriteString(",")
	out.WriteString(this.Body.String())
	out.WriteString(")")
	return out.String()
}
func (this *MapExpr) Eval(e object.Env) (object.Object, error) {
	v, err := this.Arr.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	arr, err := v.AsArray()
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
		return object.NewArray(object.Objects{}), nil
	}
	sz := len(arr.Items)
	r := make(object.Objects, sz)
	for i, item := range arr.Items {
		v, err := cb.Call(object.Objects{object.NewInteger(int64(i)), item})
		if nil != err {
			return object.Nil, err
		}
		r[i] = v
	}
	return object.NewArray(r), nil
}
