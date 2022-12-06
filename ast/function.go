package ast

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// Function : implement Expression
type Function struct {
	defaultNode
	Lambda string
	Name   string
	Args   IdentifierSlice
	Body   *BlockStmt
}

func (this *Function) Do(v Visitor) error {
	return v.DoFn(this)
}

func (this *Function) value() map[string]interface{} {
	m := map[string]interface{}{
		"lambda": this.Lambda,
		"name":   this.Name,
		"args":   this.Args.encode(),
		"body":   this.Body.Encode(),
	}
	return m
}

func (this *Function) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprFn,
		keyValue: this.value(),
	}
}
func (this *Function) Decode(b []byte) error {
	var v struct {
		Lambda string          `json:"lambda"`
		Name   string          `json:"name"`
		Args   json.RawMessage `json:"args"`
		Body   JsonNode        `json:"body"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Lambda = v.Lambda
	this.Name = v.Name
	this.Args, err = decodeIdents(v.Args)
	if nil != err {
		return function.NewError(err)
	}
	this.Body, err = v.Body.decodeBlockStmt()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *Function) expressionNode() {}

func (this *Function) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, p := range this.Args {
		args = append(args, p.String())
	}
	if "" == this.Name {
		out.WriteString("func ")
	} else {
		out.WriteString(this.Name)
	}

	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	out.WriteString(this.Body.String())

	return out.String()
}

func (this *Function) Eval(e object.Env) (object.Object, error) {
	return object.NewFunction(
		this.Name,
		this.Args.Values(),
		this.evalBody(),
		e,
	), nil
}

func (this *Function) AsFunction() (*Function, error) {
	return this, nil
}

func (this *Function) evalBody() func(e object.Env) (object.Object, error) {
	return func(e object.Env) (object.Object, error) {
		return this.Body.Eval(e)
	}
}
