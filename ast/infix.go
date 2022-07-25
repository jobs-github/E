package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
	"github.com/jobs-github/Q/token"
)

// InfixExpr : implement Expression
type InfixExpr struct {
	Left  Expression
	Op    *token.Token
	Right Expression
}

func (this *InfixExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprInfix,
		keyValue: map[string]interface{}{
			"left":  this.Left.Encode(),
			"op":    this.Op.Literal,
			"right": this.Right.Encode(),
		},
	}
}
func (this *InfixExpr) Decode(b []byte) error {
	var v struct {
		Left  JsonNode `json:"left"`
		Op    string   `json:"op"`
		Right JsonNode `json:"right"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Left, err = v.Left.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Op, err = token.GetInfixToken(v.Op)
	if nil != err {
		return function.NewError(err)
	}
	this.Right, err = v.Right.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *InfixExpr) expressionNode() {}

func (this *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(this.Left.String())
	out.WriteString(" ")
	out.WriteString(this.Op.Literal)
	out.WriteString(" ")
	out.WriteString(this.Right.String())
	out.WriteString(")")
	return out.String()
}
func (this *InfixExpr) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	left, err := this.Left.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	right, err := this.Right.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return left.Calc(this.Op, right)
}
func (this *InfixExpr) walk(cb func(module string)) {
	this.Left.walk(cb)
	this.Right.walk(cb)
}
func (this *InfixExpr) doDefer(env object.Env) error { return nil }
