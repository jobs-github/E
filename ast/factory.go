package ast

import (
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

const (
	typeNodeProgram      = "program"
	typeStmtVar          = token.Var
	typeStmtFn           = token.Func
	typeStmtExpr         = "expr"
	typeStmtAssignindex  = "assignindex"
	typeStmtAssign       = "assign"
	typeStmtBlock        = "block"
	typeStmtReturn       = token.Return
	typeExprIdent        = "ident"
	typeExprBuiltin      = object.TypeBuiltin
	typeExprFn           = "fn"
	typeExprArray        = object.TypeArray
	typeExprNull         = token.Null
	typeExprBoolean      = object.TypeBool
	typeExprInteger      = object.TypeInt
	typeExprString       = object.TypeStr
	typeExprCall         = "call"
	typeExprCallmember   = "callmember"
	typeExprObjectmember = "objectmember"
	typeExprConditional  = "conditional"
	typeExprHash         = object.TypeHash
	typeExprIf           = token.If
	typeExprIndex        = "index"
	typeExprInfix        = "infix"
	typeExprPrefix       = "prefix"
)

func NewVar() *VarStmt                 { return &VarStmt{} }
func NewFunction() *FunctionStmt       { return &FunctionStmt{} }
func NewExpr() *ExpressionStmt         { return &ExpressionStmt{} }
func NewAssignIndex() *AssignIndexStmt { return &AssignIndexStmt{} }
func NewAssign() *AssignStmt           { return &AssignStmt{} }
func NewBlock() *BlockStmt             { return &BlockStmt{} }
func NewReturn() *ReturnStmt           { return &ReturnStmt{} }
func NewIdent() *Identifier            { return &Identifier{} }
func NewFn() *Function                 { return &Function{} }
func NewArray() *Array                 { return &Array{} }
func NewNull() *Null                   { return &Null{} }
func NewBoolean() *Boolean             { return &Boolean{} }
func NewInteger() *Integer             { return &Integer{} }
func NewString() *String               { return &String{} }
func NewCall() *Call                   { return &Call{} }
func NewCallMember() *CallMember       { return &CallMember{} }
func NewObjectMember() *ObjectMember   { return &ObjectMember{} }
func NewConditional() *ConditionalExpr { return &ConditionalExpr{} }
func NewHash() *Hash                   { return &Hash{} }
func NewIf() *IfExpr                   { return &IfExpr{} }
func NewIndex() *IndexExpr             { return &IndexExpr{} }
func NewInfix() *InfixExpr             { return &InfixExpr{} }
func NewPrefix() *PrefixExpr           { return &PrefixExpr{} }

var (
	stmtFactory = map[string]func() Statement{
		typeStmtVar:         func() Statement { return NewVar() },
		typeStmtFn:          func() Statement { return NewFunction() },
		typeStmtExpr:        func() Statement { return NewExpr() },
		typeStmtAssignindex: func() Statement { return NewAssignIndex() },
		typeStmtAssign:      func() Statement { return NewAssign() },
		typeStmtBlock:       func() Statement { return NewBlock() },
		typeStmtReturn:      func() Statement { return NewReturn() },
	}
	exprFactory = map[string]func() Expression{
		typeExprIdent:        func() Expression { return NewIdent() },
		typeExprBuiltin:      func() Expression { return NewIdent() },
		typeExprFn:           func() Expression { return NewFn() },
		typeExprArray:        func() Expression { return NewArray() },
		typeExprNull:         func() Expression { return NewNull() },
		typeExprBoolean:      func() Expression { return NewBoolean() },
		typeExprInteger:      func() Expression { return NewInteger() },
		typeExprString:       func() Expression { return NewString() },
		typeExprCall:         func() Expression { return NewCall() },
		typeExprCallmember:   func() Expression { return NewCallMember() },
		typeExprObjectmember: func() Expression { return NewObjectMember() },
		typeExprConditional:  func() Expression { return NewConditional() },
		typeExprHash:         func() Expression { return NewHash() },
		typeExprIf:           func() Expression { return NewIf() },
		typeExprIndex:        func() Expression { return NewIndex() },
		typeExprInfix:        func() Expression { return NewInfix() },
		typeExprPrefix:       func() Expression { return NewPrefix() },
	}
)
