package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// RangeExpr
func (this *visitor) DoRange(v *ast.RangeExpr) error {
	i := newIdent(loopIter)
	cnt := newIdent(loopCnt)
	res := newIdent(loopResult)
	if err := this.doRangeFn(v, i, cnt, res); nil != err {
		return function.NewError(err)
	}
	// push 0
	if err := ast.NewInteger().Do(this); nil != err {
		return function.NewError(err)
	}
	// push cnt
	if err := v.Cnt.Do(this); nil != err {
		return function.NewError(err)
	}
	// push res
	if _, err := this.c.encode(code.OpArrayReserve); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 3); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doRangeFn(v *ast.RangeExpr, i *ast.Identifier, cnt *ast.Identifier, res *ast.Identifier) error {
	this.c.enterScope()

	// args
	si := this.c.define(i.Value)
	this.c.define(cnt.Value)
	ri := this.c.define(res.Value)

	startPos, err := this.doLoopCond(i, cnt)
	if nil != err {
		return function.NewError(err)
	}
	endPos, err := this.c.encode(code.OpJumpWhenFalse, -1)
	if nil != err {
		return function.NewError(err)
	}
	if err := this.doRangeBody(i, v, ri); nil != err {
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
	if _, err := this.c.encode(code.OpReturn); nil != err {
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

func (this *visitor) doRangeBody(i *ast.Identifier, v *ast.RangeExpr, res *Symbol) error {
	// push closure
	if err := v.Body.Do(this); nil != err {
		return function.NewError(err)
	}
	// push args
	if _, err := this.doIdent(i); nil != err { // push i
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 1); nil != err {
		return function.NewError(err)
	}
	if _, err := this.doIdent(i); nil != err { // push i
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpArraySet, res.Index); nil != err {
		return function.NewError(err)
	}
	return nil
}
