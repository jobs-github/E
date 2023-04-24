package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
)

// mapImpl : implement loop
type mapImpl struct {
	define       func(name string) *Symbol
	doLimitFn    func(arr *ast.Identifier) error
	doBodyArgsFn func(i *ast.Identifier, arr *ast.Identifier) error
	doBodyRetFn  func(i *ast.Identifier, ri *Symbol) error
	doReturnFn   func(res *ast.Identifier) error
	i            *ast.Identifier
	arr          *ast.Identifier
	res          *ast.Identifier
	ri           *Symbol
}

func (this *mapImpl) prepare()          { this.ri = this.define(this.res.Value) }
func (this *mapImpl) doLimit() error    { return this.doLimitFn(this.arr) }
func (this *mapImpl) doBodyArgs() error { return this.doBodyArgsFn(this.i, this.arr) }
func (this *mapImpl) doBodyRet() error  { return this.doBodyRetFn(this.i, this.ri) }
func (this *mapImpl) doReturn() error   { return this.doReturnFn(this.res) }

// MapExpr
func (this *visitor) DoMap(v *ast.MapExpr) error {
	i := newIdent(loopIter)
	arr := newIdent(loopArray)
	res := newIdent(loopResult)

	l := &mapImpl{
		define:       this.define,
		doLimitFn:    this.doLimitArrLen,
		doBodyArgsFn: this.doMapArgs,
		doBodyRetFn:  this.doMapRet,
		doReturnFn:   this.doReturnRes,
		i:            i,
		arr:          arr,
		res:          res,
	}
	if err := this.doLoop(l, v.Body, i, arr); nil != err {
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
	if _, err := this.c.encode(code.OpArrayNew, 1); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 3); nil != err {
		return function.NewError(err)
	}
	return nil
}

// filter
func (this *visitor) doMapArgs(i *ast.Identifier, arr *ast.Identifier) error {
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
	return nil
}

// range
func (this *visitor) doMapRet(i *ast.Identifier, ri *Symbol) error {
	if _, err := this.doIdent(i); nil != err { // push i
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpArraySet, ri.Index); nil != err {
		return function.NewError(err)
	}
	return nil
}
