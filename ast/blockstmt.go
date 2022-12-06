package ast

import (
	"bytes"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// BlockStmt : implement Statement
type BlockStmt struct {
	defaultNode
	Stmt Statement
}

func (this *BlockStmt) Do(v Visitor) error {
	return v.DoBlock(this)
}

func (this *BlockStmt) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeStmtBlock,
		keyValue: this.Stmt.Encode(),
	}
}
func (this *BlockStmt) Decode(b []byte) error {
	var err error
	this.Stmt, err = decodeStmt(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *BlockStmt) statementNode() {}

func (this *BlockStmt) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	out.WriteString(this.Stmt.String())
	out.WriteString("}")
	return out.String()
}

func (this *BlockStmt) Eval(e object.Env) (object.Object, error) {
	return this.Stmt.Eval(e)
}
