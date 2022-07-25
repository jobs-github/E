package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
	"github.com/jobs-github/Q/token"
)

type IfClause struct {
	If   Expression
	Then *BlockStmt
}

func (this *IfClause) encode() interface{} {
	return map[string]interface{}{
		token.If: this.If.Encode(),
		"then":   this.Then.Encode(),
	}
}
func (this *IfClause) walk(cb func(module string)) {
	this.Then.walk(cb)
}

type IfClauseSlice []*IfClause

func (this *IfClauseSlice) encode() interface{} {
	arr := []interface{}{}
	for _, v := range *this {
		arr = append(arr, v.encode())
	}
	return arr
}

func (this *IfClauseSlice) walk(cb func(module string)) {
	for _, v := range *this {
		v.walk(cb)
	}
}

// IfExpr : implement Expression
type IfExpr struct {
	Clauses IfClauseSlice
	Else    *BlockStmt
}

func (this *IfExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprIf,
		keyValue: map[string]interface{}{
			"clauses": this.Clauses.encode(),
			"else":    this.Else.Encode(),
		},
	}
}
func (this *IfExpr) Decode(b []byte) error {
	var v struct {
		Clauses json.RawMessage `json:"clauses"`
		Else    JsonNode        `json:"else"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Clauses, err = decodeIfClauses(v.Clauses)
	if nil != err {
		return function.NewError(err)
	}
	this.Else, err = v.Else.decodeBlockStmt()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *IfExpr) expressionNode() {}

func (this *IfExpr) String() string {
	var out bytes.Buffer

	for i, clause := range this.Clauses {
		if 0 == i {
			out.WriteString(token.If)
		} else {
			out.WriteString("else if")
		}
		out.WriteString(clause.If.String())
		out.WriteString("{")
		out.WriteString(clause.Then.String())
		out.WriteString("}")
	}
	if nil != this.Else {
		out.WriteString("else ")
		out.WriteString("{")
		out.WriteString(this.Else.String())
		out.WriteString("}")
	}
	return out.String()
}
func (this *IfExpr) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	for _, clause := range this.Clauses {
		cond, err := clause.If.Eval(env, insideLoop)
		if nil != err {
			return object.Nil, function.NewError(err)
		}
		if cond.True() {
			return clause.Then.Eval(env.NewEnclosedEnv(), insideLoop)
		}
	}
	if nil != this.Else {
		return this.Else.Eval(env.NewEnclosedEnv(), insideLoop)
	}
	return object.Nil, nil
}
func (this *IfExpr) walk(cb func(module string)) {
	this.Clauses.walk(cb)
	if nil != this.Else {
		this.Else.walk(cb)
	}
}
func (this *IfExpr) doDefer(env object.Env) error { return nil }
