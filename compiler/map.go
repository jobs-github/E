package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// MapExpr
func (this *visitor) DoMap(v *ast.MapExpr) error {
	i := newIdent(loopIter)
	arr := newIdent(loopArray)
	res := newIdent(loopResult)
	if err := this.doMapFn(v, i, arr, res); nil != err {
		return function.NewError(err)
	}
	// push 0
	if err := ast.NewInteger().Do(this); nil != err {
		return function.NewError(err)
	}
	// push arr
	if err := v.Arr.Do(this); nil != err {
		return function.NewError(err)
	}
	// push res
	if _, err := this.c.encode(code.OpNewArray); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 3); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doMapFn(v *ast.MapExpr, i *ast.Identifier, arr *ast.Identifier, res *ast.Identifier) error {
	this.c.enterScope()

	// args
	si := this.c.define(i.Value)
	this.c.define(arr.Value)
	ri := this.c.define(res.Value)

	startPos, err := this.doArrayCond(i, arr)
	if nil != err {
		return function.NewError(err)
	}
	endPos, err := this.c.encode(code.OpJumpWhenFalse, -1)
	if nil != err {
		return function.NewError(err)
	}
	if err := this.doMapBody(i, arr, v, ri); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpIncLocal, si.Index); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpJump, startPos); nil != err {
		return function.NewError(err)
	}
	// back-patching
	if err := this.c.changeOperand(endPos, this.c.pos()); nil != err {
		return function.NewError(err)
	}
	if _, err := this.doIdent(res); nil != err { // push res
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpReturn); nil != err { // return res
		return function.NewError(err)
	}
	symbols := this.c.symbols()
	r := this.c.leaveScope()

	fn := object.NewByteFunc(r.Instructions(), symbols)
	idx := this.c.addConst(fn)
	if _, err := this.c.encode(code.OpClosure, idx, 0); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doMapBody(i *ast.Identifier, arr *ast.Identifier, v *ast.MapExpr, res *Symbol) error {
	// push closure
	if err := v.Body.Do(this); nil != err {
		return function.NewError(err)
	}
	// push args
	if _, err := this.doIdent(i); nil != err { // push i
		return function.NewError(err)
	}
	if err := this.doIndex(arr, i); nil != err { // push arr[i]
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 2); nil != err {
		return function.NewError(err)
	}
	if _, err := this.doIdent(i); nil != err { // push i
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpSetLocalIdx, res.Index); nil != err {
		return function.NewError(err)
	}
	return nil
}
