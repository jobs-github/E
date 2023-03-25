package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// FilterExpr : implement Expression
type FilterExpr struct {
	defaultNode
	Arr  Expression // Array or Identifier
	Body Expression // Function or Identifier
}

func (this *FilterExpr) expressionNode() {}

func (this *FilterExpr) Do(v Visitor) error {
	return v.DoFilter(this)
}

func (this *FilterExpr) value() map[string]interface{} {
	m := map[string]interface{}{
		"arr":  this.Arr.Encode(),
		"body": this.Body.Encode(),
	}
	return m
}

func (this *FilterExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprFilter,
		keyValue: this.value(),
	}
}
func (this *FilterExpr) Decode(b []byte) error {
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
func (this *FilterExpr) String() string {
	var out bytes.Buffer
	out.WriteString("filter(")
	out.WriteString(this.Arr.String())
	out.WriteString(",")
	out.WriteString(this.Body.String())
	out.WriteString(")")
	return out.String()
}
func (this *FilterExpr) Eval(e object.Env) (object.Object, error) {
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
	if !object.IsCallable(cb) {
		return object.Nil, errNotCallable
	}
	if arr.Items == nil || len(arr.Items) < 1 {
		return object.NewArray(object.Objects{}), nil
	}
	r := object.Objects{}
	for i, item := range arr.Items {
		v, err := cb.Call(object.Objects{object.NewInteger(int64(i)), item})
		if nil != err {
			return object.Nil, err
		}
		if v.True() {
			r = append(r, item)
		}
	}
	return object.NewArray(r), nil
}
