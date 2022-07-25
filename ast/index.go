package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// IndexExpr : implement Expression
type IndexExpr struct {
	Left  Expression // array
	Index Expression
}

func (this *IndexExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprIndex,
		keyValue: map[string]interface{}{
			"left":  this.Left.Encode(),
			"index": this.Index.Encode(),
		},
	}
}
func (this *IndexExpr) Decode(b []byte) error {
	var v struct {
		Left  JsonNode `json:"left"`
		Index JsonNode `json:"index"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Left, err = v.Left.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Index, err = v.Index.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *IndexExpr) expressionNode() {}

func (this *IndexExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(this.Left.String())
	out.WriteString("[")
	out.WriteString(this.Index.String())
	out.WriteString("])")
	return out.String()
}

// IndexExpr -> ast -> Eval() -> object -> Index
func (this *IndexExpr) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	left, err := this.Left.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	idx, err := this.Index.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return left.CallMember(object.FnIndex, object.Objects{idx})
}
func (this *IndexExpr) walk(cb func(module string))  {}
func (this *IndexExpr) doDefer(env object.Env) error { return nil }
