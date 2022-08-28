package ast

import (
	"github.com/jobs-github/escript/function"
)

// ExpressionStmt : implement Statement
type ExpressionStmt struct {
	Expr Expression
}

func (this *ExpressionStmt) Do(v Visitor) error {
	return v.DoExpr(this)
}

func (this *ExpressionStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeStmtExpr,
		keyValue: this.Expr.Encode(),
	}
}
func (this *ExpressionStmt) Decode(b []byte) error {
	var err error
	this.Expr, err = decodeExpr(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ExpressionStmt) statementNode() {}

func (this *ExpressionStmt) String() string {
	if this.Expr != nil {
		return this.Expr.String()
	}
	return ""
}
