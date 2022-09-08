package compiler

import (
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/object"
)

type CompilationScope struct {
}

func newBytecode(i code.Instructions, c object.Objects) Bytecode {
	return &bytecode{
		instructions: i,
		lastIns:      EncodedInstruction{},
		prevLastIns:  EncodedInstruction{},
		constants:    c,
	}
}

type EncodedInstruction struct {
	Opcode code.Opcode
	Pos    int
}

type Bytecode interface {
	Instructions() code.Instructions
	Constants() object.Objects

	opCode(pos int) code.Opcode
	addConst(obj object.Object) int
	addInstruction(ins []byte) int
	replaceInstruction(pos int, newInstruction []byte)
	setLastInstruction(op code.Opcode, pos int)
	// lastInstructionIsPop() bool // TODO
	// removeLastInstruction()     // TODO
}

// bytecode : implement Bytecode
type bytecode struct {
	instructions code.Instructions
	lastIns      EncodedInstruction
	prevLastIns  EncodedInstruction
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

func (this *bytecode) lastInstructionIsPop() bool {
	return this.lastIns.Opcode == code.OpPop
}

func (this *bytecode) removeLastInstruction() {
	this.removeTail(this.lastIns.Pos)
	this.lastIns = this.prevLastIns
}

func (this *bytecode) setLastInstruction(op code.Opcode, pos int) {
	this.prevLastIns = this.lastIns
	this.lastIns = EncodedInstruction{Opcode: op, Pos: pos}
}
