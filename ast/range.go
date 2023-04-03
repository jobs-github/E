package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// RangeExpr : implement Expression
type RangeExpr struct {
	defaultNode
	Cnt  Expression // loop count
	Body Expression // Function or Identifier
}

func (this *RangeExpr) expressionNode() {}

func (this *RangeExpr) Do(v Visitor) error {
	return v.DoRange(this)
}

func (this *RangeExpr) value() map[string]interface{} {
	m := map[string]interface{}{
		"cnt":  this.Cnt.Encode(),
		"body": this.Body.Encode(),
	}
	return m
}

func (this *RangeExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprRange,
		keyValue: this.value(),
	}
}
func (this *RangeExpr) Decode(b []byte) error {
	var v struct {
		Cnt  JsonNode `json:"start"`
		Body JsonNode `json:"body"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Cnt, err = v.Cnt.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Body, err = v.Body.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *RangeExpr) String() string {
	var out bytes.Buffer
	out.WriteString("range(")
	out.WriteString(this.Cnt.String())
	out.WriteString(",")
	out.WriteString(this.Body.String())
	out.WriteString(")")
	return out.String()
}
func (this *RangeExpr) Eval(e object.Env) (object.Object, error) {
	v, err := this.Cnt.Eval(e)
	cnt, err := object.ToInteger(v)
	if nil != err {
		return object.Nil, err
	}
	fn, err := this.Body.Eval(e)
	if !object.IsCallable(fn) {
		return object.Nil, err
	}
	r := make(object.Objects, cnt)
	for i := int64(0); i < cnt; i++ {
		v, err := fn.Call(object.Objects{object.NewInteger(i)})
		if nil != err {
			return object.Nil, err
		}
		r[i] = v
	}
	return object.NewArray(r), nil
}
