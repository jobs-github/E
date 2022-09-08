package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

type Compiler interface {
	Compile(node ast.Node) error
	Bytecode() Bytecode

	addConst(obj object.Object) int
	// return pos before encode
	encode(op code.Opcode, operands ...int) (int, error)
	pos() int

	changeOperand(opPos int, operand int) error

	define(key string) *Symbol
	resolve(key string) (*Symbol, error)
}

func Make(s SymbolTable, consts object.Objects) Compiler {
	return &compilerImpl{
		b:  newBytecode(code.Instructions{}, consts),
		st: s,
	}
}

func New() Compiler {
	return Make(NewSymbolTable(), object.Objects{})
}

// compilerImpl : implement Compiler
type compilerImpl struct {
	b  Bytecode
	st SymbolTable
}

func (this *compilerImpl) Compile(node ast.Node) error {
	return node.Do(newVisitor(this, nil))
}

func (this *compilerImpl) Bytecode() Bytecode {
	return this.b
}

func (this *compilerImpl) addConst(obj object.Object) int {
	return this.b.addConst(obj)
}

func (this *compilerImpl) encode(op code.Opcode, operands ...int) (int, error) {
	ins, err := code.Make(op, operands...)
	if nil != err {
		return -1, function.NewError(err)
	}
	lastPos := this.addInstruction(ins)
	this.b.setLastInstruction(op, lastPos)
	return lastPos, nil
}

func (this *compilerImpl) pos() int {
	return len(this.b.Instructions())
}

func (this *compilerImpl) changeOperand(opPos int, operand int) error {
	op := this.b.opCode(opPos)
	newIns, err := code.Make(op, operand)
	if nil != err {
		return function.NewError(err)
	}
	this.b.replaceInstruction(opPos, newIns)
	return nil
}

func (this *compilerImpl) addInstruction(ins []byte) int {
	return this.b.addInstruction(ins)
}

func (this *compilerImpl) define(key string) *Symbol {
	return this.st.define(key)
}

func (this *compilerImpl) resolve(key string) (*Symbol, error) {
	return this.st.resolve(key)
}
