package compiler

import (
	"github.com/jobs-github/escript/code"
)

func newScopeBytecode(mainScope Bytecode) Bytecode {
	return &scopeBytecode{
		scopes:     []Bytecode{mainScope},
		scopeIndex: 0,
	}
}

func newBytecode(i code.Instructions) Bytecode {
	return &bytecode{
		instructions: i,
		lastIns:      EncodedInstruction{},
		prevLastIns:  EncodedInstruction{},
	}
}

type EncodedInstruction struct {
	Opcode code.Opcode
	Pos    int
}

type Bytecode interface {
	Instructions() code.Instructions
	ScopeCode() Bytecode // current
	Scope() int

	enterScope()
	leaveScope() Bytecode
	opCode(pos int) code.Opcode
	addInstruction(ins []byte) int
	replaceInstruction(pos int, newInstruction []byte)
	setLastInstruction(op code.Opcode, pos int)
	lastCode() code.Opcode
	prevLastCode() code.Opcode
}

// bytecode : implement Bytecode
type bytecode struct {
	instructions code.Instructions
	lastIns      EncodedInstruction
	prevLastIns  EncodedInstruction
}

func (this *bytecode) Instructions() code.Instructions {
	return this.instructions
}

func (this *bytecode) ScopeCode() Bytecode  { return nil }
func (this *bytecode) Scope() int           { return -1 }
func (this *bytecode) enterScope()          {}
func (this *bytecode) leaveScope() Bytecode { return nil }

func (this *bytecode) opCode(pos int) code.Opcode {
	return code.Opcode(this.instructions[pos])
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

//func (this *bytecode) lastInstructionIsPop() bool {
//	return this.lastIns.Opcode == code.OpPop
//}
//func (this *bytecode) removeLastInstruction() {
//	this.removeTail(this.lastIns.Pos)
//	this.lastIns = this.prevLastIns
//}

func (this *bytecode) setLastInstruction(op code.Opcode, pos int) {
	this.prevLastIns = this.lastIns
	this.lastIns = EncodedInstruction{Opcode: op, Pos: pos}
}

func (this *bytecode) lastCode() code.Opcode {
	return this.lastIns.Opcode
}

func (this *bytecode) prevLastCode() code.Opcode {
	return this.prevLastIns.Opcode
}

// scopeBytecode : implement Bytecode
type scopeBytecode struct {
	scopes     []Bytecode
	scopeIndex int
}

func (this *scopeBytecode) Instructions() code.Instructions {
	return this.scopes[this.scopeIndex].Instructions()
}

func (this *scopeBytecode) ScopeCode() Bytecode {
	return this.scopes[this.scopeIndex]
}

func (this *scopeBytecode) Scope() int {
	return this.scopeIndex
}

func (this *scopeBytecode) enterScope() {
	b := newBytecode(code.Instructions{})
	this.scopes = append(this.scopes, b)
	this.scopeIndex++
}

func (this *scopeBytecode) leaveScope() Bytecode {
	top := this.ScopeCode()
	this.scopes = this.scopes[:len(this.scopes)-1]
	this.scopeIndex--
	return top
}

func (this *scopeBytecode) opCode(pos int) code.Opcode {
	return this.ScopeCode().opCode(pos)
}

func (this *scopeBytecode) addInstruction(ins []byte) int {
	return this.ScopeCode().addInstruction(ins)
}

func (this *scopeBytecode) replaceInstruction(pos int, newInstruction []byte) {
	this.scopes[this.scopeIndex].replaceInstruction(pos, newInstruction)
}

func (this *scopeBytecode) setLastInstruction(op code.Opcode, pos int) {
	this.scopes[this.scopeIndex].setLastInstruction(op, pos)
}

func (this *scopeBytecode) lastCode() code.Opcode {
	return this.ScopeCode().lastCode()
}

func (this *scopeBytecode) prevLastCode() code.Opcode {
	return this.ScopeCode().prevLastCode()
}
