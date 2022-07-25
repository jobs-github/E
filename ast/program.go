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

func (this *Program) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	return this.Stmts.eval(false, env, false)
}
func (this *Program) walk(cb func(module string)) {
	for _, s := range this.Stmts {
		s.walk(cb)
	}
}
func (this *Program) doDefer(env object.Env) error { return nil }
