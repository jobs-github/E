package ast

import (
	"encoding/json"

	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// Boolean : implement Expression
type Boolean struct {
	Value bool
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
func (this *Boolean) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return object.ToBoolean(this.Value), nil
}
func (this *Boolean) walk(cb func(module string))  {}
func (this *Boolean) doDefer(env object.Env) error { return nil }
