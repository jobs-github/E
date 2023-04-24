package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

type loop interface {
	prepare()
	doLimit() error
	doBodyArgs() error
	doBodyRet() error
	doReturn() error
}

func (this *visitor) doBody(l loop, body ast.Expression) error {
	if err := body.Do(this); nil != err {
		return function.NewError(err)
	}
	if err := l.doBodyArgs(); nil != err {
		return function.NewError(err)
	}
	if err := l.doBodyRet(); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doLoop(l loop, body ast.Expression, i *ast.Identifier, up *ast.Identifier) error {
	this.c.enterScope()

	// args
	si := this.c.define(i.Value)
	this.c.define(up.Value)
	l.prepare()

	startPos, err := this.doCond(l, i)
	if nil != err {
		return function.NewError(err)
	}
	endPos, err := this.c.encode(code.OpJumpWhenFalse, -1)
	if nil != err {
		return function.NewError(err)
	}
	// push closure
	if err := this.doBody(l, body); nil != err {
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
	if err := l.doReturn(); nil != err {
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

func (this *visitor) doCond(l loop, i *ast.Identifier) (int, error) {
	// if loopStartPos == 0
	//     vm will quit after jmp (ip = unsigned(-1))
	if _, err := this.c.encode(code.OpPlaceholder); nil != err {
		return -1, function.NewError(err)
	}
	loopStartPos, err := this.doIdent(i) // push i
	if nil != err {
		return -1, function.NewError(err)
	}
	if err := l.doLimit(); nil != err {
		return -1, function.NewError(err)
	}
	if _, err := this.c.encode(code.OpLt); nil != err {
		return -1, function.NewError(err)
	}
	return loopStartPos, nil
}

func (this *visitor) doLimitCnt(cnt *ast.Identifier) error {
	if _, err := this.doIdent(cnt); nil != err { // push cnt
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doLimitArrLen(arr *ast.Identifier) error {
	if _, err := this.doIdent(arr); nil != err { // push arr
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpArrayLen); nil != err { // pop arr & push len
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doReturnNil() error {
	if _, err := this.doConst(object.Nil); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doReturnRes(res *ast.Identifier) error {
	if _, err := this.doIdent(res); nil != err { // push res
		return function.NewError(err)
	}
	return nil
}

// loopImpl : implement loop
type loopImpl struct {
	doLimitFn    func(cnt *ast.Identifier) error
	doBodyArgsFn func(i *ast.Identifier) error
	doBodyRetFn  func() error
	doReturnFn   func() error
	doIdent      func(v *ast.Identifier) (int, error)
	v            *ast.LoopExpr
	i            *ast.Identifier
	cnt          *ast.Identifier
}

func (this *loopImpl) prepare()          {}
func (this *loopImpl) doLimit() error    { return this.doLimitFn(this.cnt) }
func (this *loopImpl) doBodyArgs() error { return this.doBodyArgsFn(this.i) }
func (this *loopImpl) doBodyRet() error  { return this.doBodyRetFn() }
func (this *loopImpl) doReturn() error   { return this.doReturnFn() }

// LoopExpr bytecode format
//
//	     init i cnt
//	|--->cond
//	|    OpJumpWhenFalse--|
//	|    loop             |
//	|    next             |
//	|----OpJump           |
//	     ...<-------------|
func (this *visitor) DoLoop(v *ast.LoopExpr) error {
	i := newIdent(loopIter)
	cnt := newIdent(loopCnt)

	l := &loopImpl{
		doLimitFn:    this.doLimitCnt,
		doBodyArgsFn: this.doLoopArgs,
		doBodyRetFn:  this.doLoopRet,
		doReturnFn:   this.doReturnNil,
		doIdent:      this.doIdent,
		v:            v,
		i:            i,
		cnt:          cnt,
	}
	if err := this.doLoop(l, v.Body, i, cnt); nil != err {
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

func (this *visitor) doLoopArgs(i *ast.Identifier) error {
	// push args
	if _, err := this.doIdent(i); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 1); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doLoopRet() error {
	// discard return value
	if _, err := this.c.encode(code.OpPop); nil != err {
		return function.NewError(err)
	}
	return nil
}
