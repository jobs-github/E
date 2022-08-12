package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

func NewBytecode(i code.Instructions, c object.Objects) Bytecode {
	return &bytecode{i, c}
}

type Bytecode interface {
	Instructions() code.Instructions
	Constants() object.Objects

	addConst(obj object.Object) int
	addInstruction(ins []byte) int
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

func (this *bytecode) addConst(obj object.Object) int {
	this.constants = append(this.constants, obj)
	return len(this.constants) - 1
}

func (this *bytecode) addInstruction(ins []byte) int {
	lastPos := len(this.instructions)
	this.instructions = append(this.instructions, ins...)
	return lastPos
}

type Compiler interface {
	Compile(node ast.Node) error
	Bytecode() Bytecode

	addConst(obj object.Object) int
	encode(op code.Opcode, operands ...int) (int, error)
}

func New() Compiler {
	return &compilerImpl{
		b: NewBytecode(code.Instructions{}, object.Objects{}),
	}
}

// compilerImpl : implement Compiler
type compilerImpl struct {
	b Bytecode
}

func (this *compilerImpl) Compile(node ast.Node) error {
	return node.Do(&visitor{this})
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
	return lastPos, nil
}

func (this *compilerImpl) addInstruction(ins []byte) int {
	return this.b.addInstruction(ins)
}
