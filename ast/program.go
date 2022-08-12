package ast

import (
	"bytes"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// Program : implement Node
type Program struct {
	Stmts StatementSlice
}

func (this *Program) Do(v Visitor) error {
	return v.DoProgram(this)
}

func (this *Program) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeNodeProgram,
		keyValue: this.Stmts.encode(),
	}
}
func (this *Program) Decode(b []byte) error {
	var err error
	this.Stmts, err = decodeStmts(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *Program) String() string {
	var out bytes.Buffer
	for _, s := range this.Stmts {
		out.WriteString(s.String())
	}
	return out.String()
}

func (this *Program) Eval(env object.Env) (object.Object, error) {
	return this.Stmts.Eval(env)
}
