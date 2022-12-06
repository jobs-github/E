package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// ConditionalExpr : implement Expression
type ConditionalExpr struct {
	defaultNode
	Cond Expression
	Yes  Expression
	No   Expression
}

func (this *ConditionalExpr) Do(v Visitor) error {
	return v.DoConditional(this)
}

func (this *ConditionalExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprConditional,
		keyValue: map[string]interface{}{
			"cond": this.Cond.Encode(),
			"yes":  this.Yes.Encode(),
			"no":   this.No.Encode(),
		},
	}
}
func (this *ConditionalExpr) Decode(b []byte) error {
	var v struct {
		Cond JsonNode `json:"cond"`
		Yes  JsonNode `json:"yes"`
		No   JsonNode `json:"no"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Cond, err = v.Cond.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Yes, err = v.Yes.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.No, err = v.No.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ConditionalExpr) expressionNode() {}

func (this *ConditionalExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(this.Cond.String())
	out.WriteString(") ? (")
	out.WriteString(this.Yes.String())
	out.WriteString(") : (")
	out.WriteString(this.No.String())
	out.WriteString(")")
	return out.String()
}

func (this *ConditionalExpr) Eval(e object.Env) (object.Object, error) {
	r, err := this.Cond.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	node := this.getCondNode(r.True())
	return node.Eval(e.NewEnclosedEnv())
}

func (this *ConditionalExpr) getCondNode(yes bool) Expression {
	if yes {
		return this.Yes
	} else {
		return this.No
	}
}
