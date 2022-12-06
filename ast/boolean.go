package ast

import (
	"encoding/json"

	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// Boolean : implement Expression
type Boolean struct {
	defaultNode
	Value bool
}

func (this *Boolean) Do(v Visitor) error {
	return v.DoBoolean(this)
}

func (this *Boolean) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprBoolean,
		keyValue: this.Value,
	}
}
func (this *Boolean) Decode(b []byte) error {
	return json.Unmarshal(b, &this.Value)
}
func (this *Boolean) expressionNode() {}

func (this *Boolean) String() string {
	return token.Bool(this.Value)
}

func (this *Boolean) Eval(e object.Env) (object.Object, error) {
	return object.ToBoolean(this.Value), nil
}
