package ast

import (
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

const (
	typeNodeProgram      = "program"
	typeStmtConst        = token.Const
	typeStmtFn           = token.Func
	typeStmtExpr         = "expr"
	typeStmtBlock        = "block"
	typeExprIdent        = "ident"
	typeExprLoop         = token.Loop
	typeExprMap          = token.Map
	typeExprReduce       = token.Reduce
	typeExprFilter       = token.Filter
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
	typeExprIndex        = "index"
	typeExprInfix        = "infix"
	typeExprPrefix       = "prefix"
)

func NewConst() *ConstStmt             { return &ConstStmt{} }
func NewFunction() *FunctionStmt       { return &FunctionStmt{} }
func NewExpr() *ExpressionStmt         { return &ExpressionStmt{} }
func NewBlock() *BlockStmt             { return &BlockStmt{} }
func NewIdent() *Identifier            { return &Identifier{} }
func NewLoop() *LoopExpr               { return &LoopExpr{} }
func NewMap() *MapExpr                 { return &MapExpr{} }
func NewReduce() *ReduceExpr           { return &ReduceExpr{} }
func NewFilter() *FilterExpr           { return &FilterExpr{} }
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
func NewIndex() *IndexExpr             { return &IndexExpr{} }
func NewInfix() *InfixExpr             { return &InfixExpr{} }
func NewPrefix() *PrefixExpr           { return &PrefixExpr{} }

var (
	stmtFactory = map[string]func() Statement{
		typeStmtConst: func() Statement { return NewConst() },
		typeStmtFn:    func() Statement { return NewFunction() },
		typeStmtExpr:  func() Statement { return NewExpr() },
		typeStmtBlock: func() Statement { return NewBlock() },
	}
	exprFactory = map[string]func() Expression{
		typeExprIdent:        func() Expression { return NewIdent() },
		typeExprLoop:         func() Expression { return NewLoop() },
		typeExprMap:          func() Expression { return NewMap() },
		typeExprReduce:       func() Expression { return NewReduce() },
		typeExprFilter:       func() Expression { return NewFilter() },
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
		typeExprIndex:        func() Expression { return NewIndex() },
		typeExprInfix:        func() Expression { return NewInfix() },
		typeExprPrefix:       func() Expression { return NewPrefix() },
	}
)
