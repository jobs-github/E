package ast

import (
	"fmt"
	"strconv"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// Integer : implement Expression
type Integer struct {
	Value int64
}

func (this *Integer) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprInteger,
		keyValue: this.Value,
	}
}
func (this *Integer) Decode(b []byte) error {
	v := function.BytesToString(b)
	i, err := strconv.ParseInt(v, 10, 64)
	if nil != err {
		return function.NewError(err)
	}
	this.Value = i
	return nil
}
func (this *Integer) expressionNode() {}

func (this *Integer) String() string {
	return fmt.Sprintf("%v", this.Value)
}
func (this *Integer) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return object.NewInteger(this.Value), nil
}
func (this *Integer) walk(cb func(module string))  {}
func (this *Integer) doDefer(env object.Env) error { return nil }
