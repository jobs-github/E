package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// PrefixExpr : implement Expression
type PrefixExpr struct {
	Op    *token.Token
	Right Expression
}

func (this *PrefixExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprPrefix,
		keyValue: map[string]interface{}{
			"op":    this.Op.Literal,
			"right": this.Right.Encode(),
		},
	}
}
func (this *PrefixExpr) Decode(b []byte) error {
	var v struct {
		Op    string   `json:"op"`
		Right JsonNode `json:"right"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
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
func (this *PrefixExpr) expressionNode() {}

func (this *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(this.Op.Literal)
	out.WriteString(this.Right.String())
	out.WriteString(")")
	return out.String()
}
func (this *PrefixExpr) Eval(env object.Env) (object.Object, error) {
	right, err := this.Right.Eval(env)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return evalPrefixExpression(this.Op, right)
}
