package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
)

// rangeImpl : implement loop
type rangeImpl struct {
	define       func(name string) *Symbol
	doLimitFn    func(cnt *ast.Identifier) error
	doBodyArgsFn func(i *ast.Identifier) error
	doBodyRetFn  func(i *ast.Identifier, ri *Symbol) error
	doReturnFn   func(res *ast.Identifier) error
	i            *ast.Identifier
	cnt          *ast.Identifier
	res          *ast.Identifier
	ri           *Symbol
}

func (this *rangeImpl) prepare()          { this.ri = this.define(this.res.Value) }
func (this *rangeImpl) doLimit() error    { return this.doLimitFn(this.cnt) }
func (this *rangeImpl) doBodyArgs() error { return this.doBodyArgsFn(this.i) }
func (this *rangeImpl) doBodyRet() error  { return this.doBodyRetFn(this.i, this.ri) }
func (this *rangeImpl) doReturn() error   { return this.doReturnFn(this.res) }

// RangeExpr
func (this *visitor) DoRange(v *ast.RangeExpr) error {
	i := newIdent(loopIter)
	cnt := newIdent(loopCnt)
	res := newIdent(loopResult)
	l := &rangeImpl{
		define:       this.define,
		doLimitFn:    this.doLimitCnt,
		doBodyArgsFn: this.doRangeArgs,
		doBodyRetFn:  this.doMapRet,
		doReturnFn:   this.doReturnRes,
		i:            i,
		cnt:          cnt,
		res:          res,
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
	// push res
	if _, err := this.c.encode(code.OpArrayReserve); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 3); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doRangeArgs(i *ast.Identifier) error {
	// push args
	if _, err := this.doIdent(i); nil != err { // push i
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 1); nil != err {
		return function.NewError(err)
	}
	return nil
}
