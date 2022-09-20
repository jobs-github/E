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
	Constants() object.Objects

	enterScope()
	leaveScope() Bytecode
	addConst(obj object.Object) int
	// return pos before encode
	encode(op code.Opcode, operands ...int) (int, error)
	pos() int

	changeOperand(opPos int, operand int) error

	define(key string) *Symbol
	resolve(key string) (*Symbol, error)
	symbols() int
}

func Make(s SymbolTable, consts object.Objects) Compiler {
	return &compilerImpl{
		st:        s,
		b:         newScopeBytecode(newBytecode(code.Instructions{})),
		constants: consts,
	}
}

func New() Compiler {
	return Make(NewSymbolTable(nil), object.Objects{})
}

// compilerImpl : implement Compiler
type compilerImpl struct {
	st        SymbolTable
	b         Bytecode
	constants object.Objects
}

func (this *compilerImpl) Compile(node ast.Node) error {
	return node.Do(newVisitor(this, nil))
}

func (this *compilerImpl) Bytecode() Bytecode {
	return this.b
}

func (this *compilerImpl) Constants() object.Objects {
	return this.constants
}

func (this *compilerImpl) enterScope() {
	this.b.enterScope()
	this.st = this.st.newEnclosed()
}

func (this *compilerImpl) leaveScope() Bytecode {
	this.st = this.st.outer()
	return this.b.leaveScope()
}

func (this *compilerImpl) addConst(obj object.Object) int {
	this.constants = append(this.constants, obj)
	return len(this.constants) - 1
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

func (this *compilerImpl) symbols() int {
	return this.st.size()
}
