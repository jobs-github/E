package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

func (this *visitor) DoProgram(v *ast.Program) error {
	for _, s := range v.Stmts {
		if err := s.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func (this *visitor) DoConst(v *ast.ConstStmt) error {
	return this.doBind(v.Name, v.Value)
}

func (this *visitor) DoBlock(v *ast.BlockStmt) error {
	if err := v.Stmt.Do(this); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoExpr(v *ast.ExpressionStmt) error {
	if err := v.Expr.Do(this); nil != err {
		return function.NewError(err)
	}
	option := this.optionExpr()
	switch option {
	case optionEncodeNothing:
		return nil
	case optionEncodePop:
		if _, err := this.c.encode(code.OpPop); nil != err {
			return function.NewError(err)
		}
	case optionEncodeReturn:
		if _, err := this.c.encode(code.OpReturn); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func (this *visitor) DoFunction(v *ast.FunctionStmt) error {
	return this.doBind(v.Name, v.Value)
}

func (this *visitor) DoPrefix(v *ast.PrefixExpr) error {
	if err := v.Right.Do(this); nil != err {
		return function.NewError(err)
	}
	opCode, err := code.PrefixCode(v.Op.Type)
	if nil != err {
		return unsupportedOp(function.GetFunc(), v.Op, v.Right)
	}
	if _, err := this.c.encode(opCode); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoInfix(v *ast.InfixExpr) error {
	if err := v.Left.Do(this); nil != err {
		return function.NewError(err)
	}
	if err := v.Right.Do(this); nil != err {
		return function.NewError(err)
	}
	opCode, err := code.InfixCode(v.Op.Type)
	if nil != err {
		return unsupportedOp(function.GetFunc(), v.Op, v.Right)
	}
	if _, err := this.c.encode(opCode); nil != err {
		return function.NewError(err)
	}
	return nil
}

// ConditionalExpr bytecode format
//
//	     cond
//	     OpJumpWhenFalse--|
//	     Yes              |
//	|----OpJump           |
//	|    No<--------------|
//	|--->...
func (this *visitor) DoConditional(v *ast.ConditionalExpr) error {
	if err := v.Cond.Do(this); nil != err {
		return function.NewError(err)
	}
	posJumpWhenFalse, err := this.c.encode(code.OpJumpWhenFalse, -1)
	if nil != err {
		return function.NewError(err)
	}
	if err := v.Yes.Do(this.enclosed(optionEncodeNothing)); nil != err {
		return function.NewError(err)
	}
	posJump, err := this.c.encode(code.OpJump, -1)
	if nil != err {
		return function.NewError(err)
	}
	// back-patching
	if err := this.c.changeOperand(posJumpWhenFalse, this.c.pos()); nil != err {
		return function.NewError(err)
	}
	if err := v.No.Do(this.enclosed(optionEncodeNothing)); nil != err {
		return function.NewError(err)
	}
	// back-patching
	if err := this.c.changeOperand(posJump, this.c.pos()); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoFn(v *ast.Function) error {
	this.c.enterScope()

	if v.Lambda != "" {
		// add the lambda function's name to the symbol table
		this.c.defineLambda(v.Lambda)
	}

	for _, a := range v.Args {
		this.c.define(a.Value)
	}

	if err := v.Body.Do(this.enclosed(optionEncodeReturn)); nil != err {
		return function.NewError(err)
	}

	// after compiled a function’s body, capture the FreeSymbols before leave scope
	freeSymbols := this.c.freeSymbols()
	symbols := this.c.symbols()
	r := this.c.leaveScope()

	// vm will put the free variables on to the stack
	// waiting to be merged with an ByteFunc into an Closure.
	for _, s := range freeSymbols {
		if _, err := this.doLoadSymbol(s); nil != err {
			return function.NewError(err)
		}
	}

	fn := object.NewByteFunc(r.Instructions(), symbols)
	idx := this.c.addConst(fn)
	// not OpConst here
	if _, err := this.c.encode(code.OpClosure, idx, len(freeSymbols)); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doCall(fn ast.Expression, args ast.ExpressionSlice) error {
	if err := fn.Do(this); nil != err {
		return function.NewError(err)
	}
	for _, a := range args {
		if err := a.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	if _, err := this.c.encode(code.OpCall, len(args)); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoCall(v *ast.Call) error {
	return this.doCall(v.Func, v.Args)
}

func (this *visitor) DoCallMember(v *ast.CallMember) error {
	if err := v.Left.Do(this.enclosed(optionEncodeNothing)); nil != err {
		return function.NewError(err)
	}
	return this.doCall(v.Func, v.Args)
}

func (this *visitor) DoObjectMember(v *ast.ObjectMember) error {
	if err := v.Left.Do(this.enclosed(optionEncodeNothing)); nil != err {
		return function.NewError(err)
	}
	if err := v.Member.Do(this); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoIndex(v *ast.IndexExpr) error {
	if err := v.Left.Do(this); nil != err {
		return function.NewError(err)
	}
	if err := v.Index.Do(this); nil != err {
		return function.NewError(err)
	}
	if _, err := this.c.encode(code.OpIndex); nil != err {
		return function.NewError(err)
	}
	return nil
}
