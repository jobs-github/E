package ast

import (
	"fmt"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
	"github.com/jobs-github/Q/token"
)

type Nodes []Node

type Node interface {
	Encode() interface{}
	Decode(b []byte) error
	String() string
	Eval(env object.Env, insideLoop bool) (object.Object, error)
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

func (this *ExpressionSlice) eval(env object.Env, insideLoop bool) (object.Objects, error) {
	result := object.Objects{}
	for _, expr := range *this {
		evaluated, err := expr.Eval(env, insideLoop)
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

func (this *StatementSlice) eval(isBlockStmts bool, env object.Env, insideLoop bool) (object.Object, error) {
	var result object.Object
	for _, stmt := range *this {
		if v, err := stmt.Eval(env, insideLoop); nil != err {
			return object.Nil, function.NewError(err)
		} else {
			if needReturn, returnValue := v.Return(); needReturn {
				if isBlockStmts {
					// it stops execution in a possible deeper block statement and bubbles up to Program.Eval
					// where it finally get's unwrapped
					return v, nil
				} else {
					return returnValue, nil
				}
			}
			if insideLoop {
				isBreak, _ := v.Break()
				if isBreak {
					return v, nil
				}
			} else { // outside loop
				isBreak, breakCount := v.Break()
				if isBreak && 1 == breakCount { // orginal break
					err := fmt.Errorf("'break' outside loop, stmt(`%v`)", stmt.String())
					return object.Nil, function.NewError(err)
				}
			}
			result = v
		}
	}
	if err := this.doDefer(env); nil != err {
		return object.Nil, function.NewError(err)
	}
	return result, nil
}

func (this *StatementSlice) doDefer(env object.Env) error {
	sz := len(*this)
	for i := sz - 1; i >= 0; i-- {
		stmt := (*this)[i]
		if err := stmt.doDefer(env); nil != err {
			return function.NewError(err)
		}
	}
	return nil
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
	insideLoop bool,
	name *Identifier,
	value Expression,
) (object.Object, error) {
	val, err := value.Eval(env, insideLoop)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	env.Set(name.Value, val)
	return val, nil
}
