package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

func newBytecode(i code.Instructions, c object.Objects) Bytecode {
	return &bytecode{i, c}
}

type Bytecode interface {
	Instructions() code.Instructions
	Constants() object.Objects

	opCode(pos int) code.Opcode
	addConst(obj object.Object) int
	addInstruction(ins []byte) int
	removeTail(pos int)
	replaceInstruction(pos int, newInstruction []byte)
}

// bytecode : implement Bytecode
type bytecode struct {
	instructions code.Instructions
	constants    object.Objects
}

func (this *bytecode) Instructions() code.Instructions {
	return this.instructions
}
func (this *bytecode) Constants() object.Objects {
	return this.constants
}

func (this *bytecode) opCode(pos int) code.Opcode {
	return code.Opcode(this.instructions[pos])
}

func (this *bytecode) addConst(obj object.Object) int {
	this.constants = append(this.constants, obj)
	return len(this.constants) - 1
}

func (this *bytecode) addInstruction(ins []byte) int {
	lastPos := len(this.instructions)
	this.instructions = append(this.instructions, ins...)
	return lastPos
}

func (this *bytecode) removeTail(pos int) {
	this.instructions = this.instructions[:pos]
}

func (this *bytecode) replaceInstruction(pos int, newInstruction []byte) {
	sz := len(newInstruction)
	for i := 0; i < sz; i++ {
		this.instructions[pos+i] = newInstruction[i]
	}
}

type EncodedInstruction struct {
	Opcode code.Opcode
	Pos    int
}

type Compiler interface {
	Compile(node ast.Node) error
	Bytecode() Bytecode

	addConst(obj object.Object) int
	// return pos before encode
	encode(op code.Opcode, operands ...int) (int, error)
	pos() int
	lastInstructionIsPop() bool // TODO
	removeLastInstruction()     // TODO
	changeOperand(opPos int, operand int) error

	define(key string) *Symbol
	resolve(key string) (*Symbol, error)
}

func Make(s SymbolTable, consts object.Objects) Compiler {
	return &compilerImpl{
		b:           newBytecode(code.Instructions{}, consts),
		st:          s,
		lastIns:     EncodedInstruction{},
		prevLastIns: EncodedInstruction{},
	}
}

func New() Compiler {
	return Make(NewSymbolTable(), object.Objects{})
}

// compilerImpl : implement Compiler
type compilerImpl struct {
	b           Bytecode
	st          SymbolTable
	lastIns     EncodedInstruction
	prevLastIns EncodedInstruction
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
	this.setLastInstruction(op, lastPos)
	return lastPos, nil
}

func (this *compilerImpl) pos() int {
	return len(this.b.Instructions())
}

func (this *compilerImpl) lastInstructionIsPop() bool {
	return this.lastIns.Opcode == code.OpPop
}

func (this *compilerImpl) removeLastInstruction() {
	this.b.removeTail(this.lastIns.Pos)
	this.lastIns = this.prevLastIns
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

func (this *compilerImpl) setLastInstruction(op code.Opcode, pos int) {
	this.prevLastIns = this.lastIns
	this.lastIns = EncodedInstruction{Opcode: op, Pos: pos}
}

func (this *compilerImpl) define(key string) *Symbol {
	return this.st.define(key)
}

func (this *compilerImpl) resolve(key string) (*Symbol, error) {
	return this.st.resolve(key)
}
