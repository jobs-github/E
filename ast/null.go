package ast

import (
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// Null : implement Expression
type Null struct{}

func (this *Null) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprNull,
	}
}
func (this *Null) Decode(b []byte) error {
	return nil
}
func (this *Null) expressionNode() {}

func (this *Null) String() string {
	return token.Null
}
func (this *Null) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return object.Nil, nil
}
func (this *Null) walk(cb func(module string))  {}
func (this *Null) doDefer(env object.Env) error { return nil }
