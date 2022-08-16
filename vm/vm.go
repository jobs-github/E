package vm

import (
	"errors"

	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

const (
	StackSize = 2048
)

var (
	errStackOverflow = errors.New("stack overflow")
)

func New(b compiler.Bytecode) VM {
	return &virtualMachine{
		b:     b,
		stack: make(object.Objects, StackSize),
		sp:    0,
	}
}

type VM interface {
	Run() error
	StackTop() object.Object
	LastPopped() object.Object
}

// virtualMachine : implement VM
type virtualMachine struct {
	b     compiler.Bytecode
	stack object.Objects
	sp    int // top stack [sp - 1]
}

func (this *virtualMachine) Run() error {
	ins := this.b.Instructions()
	consts := this.b.Constants()
	sz := len(ins)
	for ip := 0; ip < sz; ip++ {
		op := code.Opcode(ins[ip])
		switch op {
		case code.OpConst:
			constIndex := code.DecodeUint16(ins[ip+1:])
			ip += 2
			err := this.push(consts[constIndex])
			if nil != err {
				return function.NewError(err)
			}
		case code.OpPop:
			this.pop()
		case code.OpAdd:
			right := this.pop()
			left := this.pop()
			r, err := left.Calc(token.Add, right)
			if nil != err {
				return function.NewError(err)
			}
			this.push(r)
		}
	}
	return nil
}

func (this *virtualMachine) StackTop() object.Object {
	if this.sp == 0 {
		return nil
	}
	return this.stack[this.sp-1]
}

func (this *virtualMachine) LastPopped() object.Object {
	return this.stack[this.sp]
}

func (this *virtualMachine) push(o object.Object) error {
	if this.sp >= StackSize {
		return function.NewError(errStackOverflow)
	}

	this.stack[this.sp] = o
	this.sp++

	return nil
}

func (this *virtualMachine) pop() object.Object {
	o := this.stack[this.sp-1]
	this.sp--
	return o
}
