package compiler

import (
	"errors"
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

var (
	errUnsupportedVisitor = errors.New("unsupported visitor")
)

func unsupportedOp(entry string, op *token.Token, node ast.Node) error {
	return fmt.Errorf("%v -> unsupported op %v(%v), (`%v`)", entry, op.Literal, token.ToString(op.Type), node.String())
}

func newVisitor(c Compiler, o *options) ast.Visitor {
	return &visitor{c, o}
}

func newOptions(skipEncodePop bool) *options {
	return &options{
		skipEncodePop: skipEncodePop,
	}
}

type options struct {
	skipEncodePop bool
}

// visitor : implement ast.Visitor
type visitor struct {
	c Compiler
	o *options
}

func (this *visitor) skipEncodePop() bool {
	return nil != this.o && this.o.skipEncodePop
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

func (this *visitor) DoProgram(v *ast.Program) error {
	for _, s := range v.Stmts {
		if err := s.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func (this *visitor) DoConst(v *ast.ConstStmt) error {
	if err := v.Value.Do(this); nil != err {
		return function.NewError(err)
	}
	symbol := this.c.define(v.Name.Value)
	if _, err := this.c.encode(code.OpSetGlobal, symbol.Index); nil != err {
		return function.NewError(err)
	}
	return nil
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
	if !this.skipEncodePop() {
		if _, err := this.c.encode(code.OpPop); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func (this *visitor) DoFunction(v *ast.FunctionStmt) error {
	return function.NewError(errUnsupportedVisitor)
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
	if _, err := this.c.encode(code.OpGetGlobal, s.Index); nil != err {
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
	if err := v.Yes.Do(newVisitor(this.c, newOptions(true))); nil != err {
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
	if err := v.No.Do(newVisitor(this.c, newOptions(true))); nil != err {
		return function.NewError(err)
	}
	// back-patching
	if err := this.c.changeOperand(posJump, this.c.pos()); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoFn(v *ast.Function) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoCall(v *ast.Call) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoCallMember(v *ast.CallMember) error {
	return function.NewError(errUnsupportedVisitor)
}

func (this *visitor) DoObjectMember(v *ast.ObjectMember) error {
	return function.NewError(errUnsupportedVisitor)
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
	return function.NewError(errUnsupportedVisitor)
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
