package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// LoopExpr bytecode format
//
//	         init i cnt
//		|--->cond
//		|    OpJumpWhenFalse--|
//		|    loop             |
//	    |    next             |
//		|----OpJump           |
//		     ...<-------------|
func (this *visitor) DoLoop(v *ast.LoopExpr) error {
	i := newIdent(loopIter)
	cnt := newIdent(loopCnt)

	if err := this.doLoopFn(v, i, cnt); nil != err {
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
	if _, err := this.c.encode(code.OpCall, 2); nil != err {
		return function.NewError(err)
	}
	return nil
}

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

func (this *visitor) doArrayCond(i *ast.Identifier, arr *ast.Identifier) (int, error) {
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
