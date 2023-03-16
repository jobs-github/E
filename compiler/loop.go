package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

func (this *visitor) doLoopFn(v *ast.LoopExpr, i *ast.Identifier, cnt *ast.Identifier) error {
	this.c.enterScope()

	// args
	si := this.c.define(i.Value)
	this.c.define(cnt.Value)

	startPos, err := this.doLoopCond(i, cnt)
	if nil != err {
		return function.NewError(err)
	}
	endPos, err := this.c.encode(code.OpJumpWhenFalse, -1)
	if nil != err {
		return function.NewError(err)
	}
	if err := this.doLoopBody(i, v); nil != err {
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
	if _, err := this.doConst(object.Nil); nil != err {
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

func (this *visitor) doLoopCond(i *ast.Identifier, cnt *ast.Identifier) (int, error) {
	// if loopStartPos == 0
	//     vm will quit after jmp (ip = unsigned(-1))
	if _, err := this.c.encode(code.OpPlaceholder); nil != err {
		return -1, function.NewError(err)
	}
	loopStartPos, err := this.doIdent(i) // push i
	if nil != err {
		return -1, function.NewError(err)
	}
	if _, err := this.doIdent(cnt); nil != err { // push cnt
		return -1, function.NewError(err)
	}
	if _, err := this.c.encode(code.OpLt); nil != err {
		return -1, function.NewError(err)
	}
	return loopStartPos, nil
}

// like DoCall
func (this *visitor) doLoopBody(i *ast.Identifier, v *ast.LoopExpr) error {
	// push closure
	if err := v.Body.Do(this); nil != err {
		return function.NewError(err)
	}
	// push args
	if _, err := this.doIdent(i); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 1); nil != err {
		return function.NewError(err)
	}
	// discard return value
	if _, err := this.c.encode(code.OpPop); nil != err {
		return function.NewError(err)
	}
	return nil
}

// MapExpr
func (this *visitor) doMapFn(v *ast.MapExpr, i *ast.Identifier, arr *ast.Identifier, res *ast.Identifier) error {
	this.c.enterScope()

	// args
	si := this.c.define(i.Value)
	this.c.define(arr.Value)
	ri := this.c.define(res.Value)

	startPos, err := this.doMapCond(i, arr)
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

func (this *visitor) doMapCond(i *ast.Identifier, arr *ast.Identifier) (int, error) {
	// if loopStartPos == 0
	//     vm will quit after jmp (ip = unsigned(-1))
	if _, err := this.c.encode(code.OpPlaceholder); nil != err {
		return -1, function.NewError(err)
	}
	loopStartPos, err := this.doIdent(i) // push i
	if nil != err {
		return -1, function.NewError(err)
	}
	if _, err := this.doIdent(arr); nil != err { // push arr
		return -1, function.NewError(err)
	}
	if _, err := this.c.encode(code.OpLen); nil != err { // pop arr & push len
		return -1, function.NewError(err)
	}
	if _, err := this.c.encode(code.OpLt); nil != err {
		return -1, function.NewError(err)
	}
	return loopStartPos, nil
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
