package ast

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

// CallMember : implement Expression
type CallMember struct {
	defaultNode
	Left Expression
	Func *Identifier
	Args ExpressionSlice
}

func (this *CallMember) Do(v Visitor) error {
	return v.DoCallMember(this)
}

func (this *CallMember) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprCallmember,
		keyValue: map[string]interface{}{
			"left":     this.Left.Encode(),
			token.Func: this.Func.Encode(),
			"args":     this.Args.encode(),
		},
	}
}
func (this *CallMember) Decode(b []byte) error {
	var v struct {
		Left JsonNode        `json:"left"`
		Func JsonNode        `json:"func"`
		Args json.RawMessage `json:"args"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Left, err = v.Left.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Func, err = v.Func.decodeIdent()
	if nil != err {
		return function.NewError(err)
	}
	this.Args, err = decodeExprs(v.Args)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *CallMember) expressionNode() {}

func (this *CallMember) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range this.Args {
		args = append(args, a.String())
	}

	out.WriteString(this.Left.String())
	out.WriteString(".")
	out.WriteString(this.Func.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
