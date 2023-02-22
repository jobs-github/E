package compiler

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

type doOption uint

const (
	optionEncodePop     doOption = 1
	optionEncodeReturn  doOption = 2
	optionEncodeNothing doOption = 3
)

func unsupportedOp(entry string, op *token.Token, node ast.Node) error {
	return fmt.Errorf("%v -> unsupported op %v(%v), (`%v`)", entry, op.Literal, token.ToString(op.Type), node.String())
}

func newVisitor(c Compiler, o *options) ast.Visitor {
	return &visitor{c, o}
}

func newOptions(optionExpr doOption) *options {
	return &options{
		optionExpr: optionExpr,
	}
}

type options struct {
	optionExpr doOption
}

// visitor : implement ast.Visitor
type visitor struct {
	c Compiler
	o *options
}

func (this *visitor) enclosed(op doOption) ast.Visitor {
	return newVisitor(this.c, newOptions(op))
}

func (this *visitor) optionExpr() doOption {
	if nil == this.o {
		return optionEncodePop
	} else {
		return this.o.optionExpr
	}
}

func (this *visitor) opCodeSymbolSet(s *Symbol) code.Opcode {
	if s.Scope == ScopeGlobal {
		return code.OpSetGlobal
	} else {
		return code.OpSetLocal
	}
}

// loadSymbol
func (this *visitor) opCodeSymbolGet(s *Symbol) code.Opcode {
	if s.Scope == ScopeGlobal {
		return code.OpGetGlobal
	} else if s.Scope == ScopeBuiltin {
		return code.OpGetBuiltin
	} else if s.Scope == ScopeObjectFn {
		return code.OpGetObjectFn
	} else if s.Scope == ScopeFree {
		return code.OpGetFree
	} else if s.Scope == ScopeLambda {
		return code.OpGetLambda
	} else {
		return code.OpGetLocal
	}
}

func (this *visitor) opCodeBoolean(v *ast.Boolean) code.Opcode {
	if v.Value {
		return code.OpTrue
	} else {
		return code.OpFalse
	}
}

func (this *visitor) doConst(v object.Object) error {
	idx := this.c.addConst(v)
	if _, err := this.c.encode(code.OpConst, idx); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doBind(name *ast.Identifier, value ast.Expression) error {
	symbol := this.c.define(name.Value)

	if err := value.Do(this); nil != err {
		return function.NewError(err)
	}

	if _, err := this.c.encode(this.opCodeSymbolSet(symbol), symbol.Index); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doLoadSymbol(s *Symbol) error {
	if _, err := this.c.encode(this.opCodeSymbolGet(s), s.Index); nil != err {
		return function.NewError(err)
	} else {
		return nil
	}
}

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

func (this *visitor) DoFor(v *ast.ForExpr) error {
	// TODO
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

func (this *visitor) DoIdent(v *ast.Identifier) error {
	s, err := this.c.resolve(v.Value)
	if nil != err {
		return function.NewError(err)
	}
	if err := this.doLoadSymbol(s); nil != err {
		return function.NewError(err)
	}
	return nil
}

// ConditionalExpr bytecode format
// cond
// OpJumpWhenFalse
// Yes
// OpJump
// No
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
		if err := this.doLoadSymbol(s); nil != err {
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

func (this *visitor) DoNull(v *ast.Null) error {
	return this.doConst(object.Nil)
}

func (this *visitor) DoInteger(v *ast.Integer) error {
	return this.doConst(object.NewInteger(v.Value))
}

func (this *visitor) DoBoolean(v *ast.Boolean) error {
	if _, err := this.c.encode(this.opCodeBoolean(v)); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoString(v *ast.String) error {
	return this.doConst(object.NewString(v.Value))
}

func (this *visitor) DoArray(v *ast.Array) error {
	// pattern: compile data first, op last
	for _, e := range v.Items {
		if err := e.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	if _, err := this.c.encode(code.OpArray, len(v.Items)); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoHash(v *ast.Hash) error {
	keys := v.Pairs.SortedKeys()
	for _, k := range keys {
		if err := k.Do(this); nil != err {
			return function.NewError(err)
		}
		v := v.Pairs[k]
		if err := v.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	if _, err := this.c.encode(code.OpHash, len(v.Pairs)); nil != err {
		return function.NewError(err)
	}
	return nil
}
