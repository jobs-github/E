package ast

import (
	"errors"
	"fmt"
	"sort"

	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

var (
	errNotFunction = errors.New("ast type is not function")
	errNotCallable = errors.New("object is not callable")
)

type Visitor interface {
	DoProgram(v *Program) error
	DoConst(v *ConstStmt) error
	DoBlock(v *BlockStmt) error
	DoExpr(v *ExpressionStmt) error
	DoLoop(v *LoopExpr) error
	DoMap(v *MapExpr) error
	DoReduce(v *ReduceExpr) error
	DoFilter(v *FilterExpr) error
	DoRange(v *RangeExpr) error
	DoFunction(v *FunctionStmt) error
	DoPrefix(v *PrefixExpr) error
	DoInfix(v *InfixExpr) error
	DoIdent(v *Identifier) error
	DoSymbol(v *SymbolExpr) error
	DoConditional(v *ConditionalExpr) error
	DoFn(v *Function) error
	DoCall(v *Call) error
	DoCallMember(v *CallMember) error
	DoObjectMember(v *ObjectMember) error
	DoIndex(v *IndexExpr) error
	DoNull(v *Null) error
	DoInteger(v *Integer) error
	DoBoolean(v *Boolean) error
	DoString(v *String) error
	DoArray(v *Array) error
	DoHash(v *Hash) error
}

type defaultNode struct{}

func (this *defaultNode) AsFunction() (*Function, error) { return nil, errNotFunction }

type Node interface {
	Do(v Visitor) error
	Encode() interface{}
	Decode(b []byte) error
	String() string
	Eval(e object.Env) (object.Object, error)
	AsFunction() (*Function, error)
}

type Nodes []Node

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionSlice []Expression

func (this *ExpressionSlice) eval(e object.Env) (object.Objects, error) {
	result := object.Objects{}
	for _, expr := range *this {
		r, err := expr.Eval(e)
		if nil != err {
			return nil, err
		}
		result.Append(r)
	}
	return result, nil
}

func (this *ExpressionSlice) encode() interface{} {
	r := []interface{}{}
	for _, v := range *this {
		r = append(r, v.Encode())
	}
	return r
}

type ExpressionMap map[Expression]Expression

func (this *ExpressionMap) SortedKeys() ExpressionSlice {
	keys := ExpressionSlice{}
	for k, _ := range *this {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	return keys
}

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

func (this *StatementSlice) Eval(e object.Env) (object.Object, error) {
	var r object.Object
	for _, stmt := range *this {
		if v, err := stmt.Eval(e); nil != err {
			return object.Nil, err
		} else {
			r = v
		}
	}
	return r, nil
}

func evalVar(name *Identifier, value Expression, e object.Env) (object.Object, error) {
	r, err := value.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	e.Set(name.Value, r)
	return r, nil
}

func evalPrefix(op *token.Token, right object.Object) (object.Object, error) {
	switch op.Type {
	case token.NOT:
		return right.CallMember(object.FnNot, object.Objects{})
	case token.SUB:
		return right.CallMember(object.FnNeg, object.Objects{})
	default:
		err := fmt.Errorf("unsupport op %v(%v)", op.Literal, token.ToString(op.Type))
		return object.Nil, err
	}
}
