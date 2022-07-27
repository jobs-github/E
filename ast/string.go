package ast

import (
	"encoding/json"

	"github.com/jobs-github/escript/object"
)

// String : implement Expression
type String struct {
	Value string
}

func (this *String) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprString,
		keyValue: this.Value,
	}
}
func (this *String) Decode(b []byte) error {
	return json.Unmarshal(b, &this.Value)
}
func (this *String) expressionNode() {}

func (this *String) String() string {
	return this.Value
}
func (this *String) Eval(env object.Env) (object.Object, error) {
	return object.NewString(this.Value), nil
}
func (this *String) walk(cb func(module string))  {}
func (this *String) doDefer(env object.Env) error { return nil }
