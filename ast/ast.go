package ast

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

type Nodes []Node

type Node interface {
	Encode() interface{}
	Decode(b []byte) error
	String() string
	Eval(env object.Env) (object.Object, error)
	walk(cb func(module string))
	doDefer(env object.Env) error
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionSlice []Expression

func (this *ExpressionSlice) encode() interface{} {
	r := []interface{}{}
	for _, v := range *this {
		r = append(r, v.Encode())
	}
	return r
}

func (this *ExpressionSlice) eval(env object.Env) (object.Objects, error) {
	result := object.Objects{}
	for _, expr := range *this {
		evaluated, err := expr.Eval(env)
		if nil != err {
			return nil, function.NewError(err)
		}
		result.Append(evaluated)
	}
	return result, nil
}

type ExpressionMap map[Expression]Expression

func (this *ExpressionMap) encode() interface{} {
	r := map[string]interface{}{}
	for k, v := range *this {
		r[k.String()] = v.Encode()
	}
	return r
}

type StatementSlice []Statement

func (this *StatementSlice) encode() interface{} {
	r := []interface{}{}
	for _, v := range *this {
		r = append(r, v.Encode())
	}
	return r
}

func (this *StatementSlice) eval(isBlockStmts bool, env object.Env) (object.Object, error) {
	var result object.Object
	for _, stmt := range *this {
		if v, err := stmt.Eval(env); nil != err {
			return object.Nil, function.NewError(err)
		} else {
			result = v
		}
	}
	return result, nil
}

func evalPrefixExpression(op *token.Token, right object.Object) (object.Object, error) {
	switch op.Type {
	case token.NOT:
		return right.CallMember(object.FnNot, object.Objects{})
	case token.SUB:
		return right.CallMember(object.FnOpposite, object.Objects{})
	default:
		err := fmt.Errorf("unsupport op %v(%v)", op.Literal, token.ToString(op.Type))
		return object.Nil, function.NewError(err)
	}
}

func evalVar(
	env object.Env,
	name *Identifier,
	value Expression,
) (object.Object, error) {
	val, err := value.Eval(env)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	env.Set(name.Value, val)
	return val, nil
}
