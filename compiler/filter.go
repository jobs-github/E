package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
)

// filterImpl : implement loop
type filterImpl struct {
	define       func(name string) *Symbol
	doLimitFn    func(arr *ast.Identifier) error
	doBodyArgsFn func(i *ast.Identifier, arr *ast.Identifier) error
	doBodyRetFn  func(i *ast.Identifier, arr *ast.Identifier, ri *Symbol) error
	doReturnFn   func(res *ast.Identifier) error
	i            *ast.Identifier
	arr          *ast.Identifier
	res          *ast.Identifier
	ri           *Symbol
}

func (this *filterImpl) prepare()          { this.ri = this.define(this.res.Value) }
func (this *filterImpl) doLimit() error    { return this.doLimitFn(this.arr) }
func (this *filterImpl) doBodyArgs() error { return this.doBodyArgsFn(this.i, this.arr) }
func (this *filterImpl) doBodyRet() error  { return this.doBodyRetFn(this.i, this.arr, this.ri) }
func (this *filterImpl) doReturn() error   { return this.doReturnFn(this.res) }

func (this *visitor) DoFilter(v *ast.FilterExpr) error {
	i := newIdent(loopIter)
	arr := newIdent(loopArray)
	res := newIdent(loopResult)
	l := &filterImpl{
		define:       this.define,
		doLimitFn:    this.doLimitArrLen,
		doBodyArgsFn: this.doMapArgs,
		doBodyRetFn:  this.doFilterRet,
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
	if _, err := this.c.encode(code.OpArrayNew, 0); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 3); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doFilterRet(i *ast.Identifier, arr *ast.Identifier, ri *Symbol) error {
	if err := this.doIndex(arr, i); nil != err { // push arr[i]
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpArrayAppend, ri.Index); nil != err {
		return function.NewError(err)
	}
	return nil
}
