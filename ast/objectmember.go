package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// ObjectMember : implement Expression
type ObjectMember struct {
	Left   Expression
	Member *Identifier
}

func (this *ObjectMember) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprObjectmember,
		keyValue: map[string]interface{}{
			"left":   this.Left.Encode(),
			"member": this.Member.Encode(),
		},
	}
}
func (this *ObjectMember) Decode(b []byte) error {
	var v struct {
		Left   JsonNode `json:"left"`
		Member JsonNode `json:"member"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Left, err = v.Left.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Member, err = v.Member.decodeIdent()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ObjectMember) expressionNode() {}

func (this *ObjectMember) String() string {
	var out bytes.Buffer

	out.WriteString(this.Left.String())
	out.WriteString(".")
	out.WriteString(this.Member.String())

	return out.String()
}
func (this *ObjectMember) Eval(env object.Env) (object.Object, error) {
	obj, err := this.Left.Eval(env)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return obj.GetMember(this.Member.Value)
}
