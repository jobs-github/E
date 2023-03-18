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
	optionEncodeTODO    doOption = 4
)

const (
	loopIter   = "__i__"
	loopCnt    = "__cnt__"
	loopArray  = "__arr__"
	loopResult = "__res__"
)

func newIdent(name string) *ast.Identifier {
	ident := ast.NewIdent()
	ident.Value = name
	return ident
}

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

// store
func (this *visitor) opCodeSymbolSet(s *Symbol) code.Opcode {
	if s.Scope == ScopeGlobal {
		return code.OpSetGlobal
	} else {
		return code.OpSetLocal
	}
}

// load
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

func (this *visitor) doConst(v object.Object) (int, error) {
	idx := this.c.addConst(v)
	if _, err := this.c.encode(code.OpConst, idx); nil != err {
		return -1, function.NewError(err)
	}
	return idx, nil
}

func (this *visitor) doBind(name *ast.Identifier, value ast.Expression) error {
	s := this.c.define(name.Value)
	if err := value.Do(this); nil != err {
		return function.NewError(err)
	}
	if _, err := this.doStoreSymbol(s); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) doStoreSymbol(s *Symbol) (int, error) {
	return this.c.encode(this.opCodeSymbolSet(s), s.Index)
}

func (this *visitor) doLoadSymbol(s *Symbol) (int, error) {
	return this.c.encode(this.opCodeSymbolGet(s), s.Index)
}

func (this *visitor) doIndex(arr *ast.Identifier, i *ast.Identifier) error {
	// push arr
	if _, err := this.doIdent(arr); nil != err {
		return function.NewError(err)
	}
	// push i
	if _, err := this.doIdent(i); nil != err {
		return function.NewError(err)
	}
	// pop i, arr & push arr[i]
	if _, err := this.c.encode(code.OpIndex); nil != err {
		return function.NewError(err)
	}
	return nil
}
