package vm

import (
	"errors"

	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
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
		case code.OpJump:
			pos := code.DecodeUint16(ins[ip+1:])
			// in a loop that increments ip with each iteration
			// we need to set ip to the offset right before the one we want
			ip = pos - 1
		case code.OpJumpWhenFalse:
			pos := code.DecodeUint16(ins[ip+1:])
			ip = ip + 2
			cond := this.pop()
			if !cond.True() {
				ip = pos - 1 // jump
			}
		case code.OpConst:
			constIndex := code.DecodeUint16(ins[ip+1:])
			ip += 2
			err := this.push(consts[constIndex])
			if nil != err {
				return function.NewError(err)
			}
		case code.OpPop:
			this.pop()
		case code.OpTrue:
			if err := this.push(object.True); nil != err {
				return function.NewError(err)
			}
		case code.OpFalse:
			if err := this.push(object.False); nil != err {
				return function.NewError(err)
			}
		case code.OpNull:
			if err := this.push(object.Nil); nil != err {
				return function.NewError(err)
			}
		case code.OpNot:
			if err := this.execPrefix(object.FnNot); nil != err {
				return function.NewError(err)
			}
		case code.OpNeg:
			if err := this.execPrefix(object.FnNeg); nil != err {
				return function.NewError(err)
			}
		case code.OpAdd,
			code.OpSub,
			code.OpMul,
			code.OpDiv,
			code.OpMod,
			code.OpLt,
			code.OpGt,
			code.OpEq,
			code.OpNeq,
			code.OpLeq,
			code.OpGeq,
			code.OpAnd,
			code.OpOr:
			if err := this.execInfix(op); nil != err {
				return function.NewError(err)
			}
		}
	}
	return nil
}

func (this *virtualMachine) execPrefix(fn string) error {
	right := this.pop()
	if r, err := right.CallMember(fn, object.Objects{}); nil != err {
		return function.NewError(err)
	} else {
		return this.push(r)
	}
}

func (this *virtualMachine) execInfix(op code.Opcode) error {
	t, err := code.InfixToken(op)
	if nil != err {
		return function.NewError(err)
	}
	right := this.pop()
	left := this.pop()
	r, err := left.Calc(t, right)
	if nil != err {
		return function.NewError(err)
	}
	this.push(r)
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
