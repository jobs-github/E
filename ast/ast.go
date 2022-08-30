package ast

type Visitor interface {
	DoProgram(v *Program) error
	DoConst(v *ConstStmt) error
	DoBlock(v *BlockStmt) error
	DoExpr(v *ExpressionStmt) error
	DoFunction(v *FunctionStmt) error
	DoPrefix(v *PrefixExpr) error
	DoInfix(v *InfixExpr) error
	DoIdent(v *Identifier) error
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

type Node interface {
	Do(v Visitor) error
	Encode() interface{}
	Decode(b []byte) error
	String() string
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

func (this *ExpressionSlice) encode() interface{} {
	r := []interface{}{}
	for _, v := range *this {
		r = append(r, v.Encode())
	}
	return r
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
