package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
)

// reduceImpl : implement loop
type reduceImpl struct {
	define       func(name string) *Symbol
	doLimitFn    func(arr *ast.Identifier) error
	doBodyArgsFn func(i *ast.Identifier, arr *ast.Identifier, acc *ast.Identifier) error
	doBodyRetFn  func(ri *Symbol) error
	doReturnFn   func(res *ast.Identifier) error
	doIdent      func(v *ast.Identifier) (int, error)
	doIndex      func(arr *ast.Identifier, i *ast.Identifier) error
	i            *ast.Identifier
	arr          *ast.Identifier
	res          *ast.Identifier
	ri           *Symbol
}

func (this *reduceImpl) prepare()          { this.ri = this.define(this.res.Value) }
func (this *reduceImpl) doLimit() error    { return this.doLimitFn(this.arr) }
func (this *reduceImpl) doBodyArgs() error { return this.doBodyArgsFn(this.i, this.arr, this.res) }
func (this *reduceImpl) doBodyRet() error  { return this.doBodyRetFn(this.ri) }
func (this *reduceImpl) doReturn() error   { return this.doReturnFn(this.res) }

// ReduceExpr
func (this *visitor) DoReduce(v *ast.ReduceExpr) error {
	i := newIdent(loopIter)
	arr := newIdent(loopArray)
	res := newIdent(loopResult)
	l := &reduceImpl{
		define:       this.define,
		doLimitFn:    this.doLimitArrLen,
		doBodyArgsFn: this.doReduceArgs,
		doBodyRetFn:  this.doReduceRet,
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
	// push init
	if err := v.Init.Do(this); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpCall, 3); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doReduceRet(ri *Symbol) error {
	if _, err := this.c.encode(code.OpSetLocal, ri.Index); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doReduceArgs(i *ast.Identifier, arr *ast.Identifier, acc *ast.Identifier) error {
	// push args
	if _, err := this.doIdent(acc); nil != err { // push acc
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
