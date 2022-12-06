package ast

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

// Call : implement Expression
type Call struct {
	defaultNode
	Func Expression
	Args ExpressionSlice
}

func (this *Call) Do(v Visitor) error {
	return v.DoCall(this)
}

func (this *Call) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprCall,
		keyValue: map[string]interface{}{
			token.Func: this.Func.Encode(),
			"args":     this.Args.encode(),
		},
	}
}
func (this *Call) Decode(b []byte) error {
	var v struct {
		Func JsonNode        `json:"func"`
		Args json.RawMessage `json:"args"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Func, err = v.Func.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Args, err = decodeExprs(v.Args)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *Call) expressionNode() {}

func (this *Call) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range this.Args {
		args = append(args, a.String())
	}

	out.WriteString(this.Func.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (this *Call) Eval(e object.Env) (object.Object, error) {
	fn, err := this.Func.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	args, err := this.Args.eval(e)
	if nil != err {
		return object.Nil, err
	}
	return fn.Call(args)
}
