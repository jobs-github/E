package ast

import (
	"bytes"
	"strings"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// Array : implement Expression
type Array struct {
	Items ExpressionSlice
}

func (this *Array) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprArray,
		keyValue: this.Items.encode(),
	}
}

func (this *Array) Decode(b []byte) error {
	var err error
	this.Items, err = decodeExprs(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *Array) expressionNode() {}

func (this *Array) String() string {
	var out bytes.Buffer
	items := []string{}
	for _, v := range this.Items {
		items = append(items, v.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")
	return out.String()
}
func (this *Array) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	items, err := this.Items.eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return object.NewArray(items), nil
}
func (this *Array) walk(cb func(module string))  {}
func (this *Array) doDefer(env object.Env) error { return nil }
