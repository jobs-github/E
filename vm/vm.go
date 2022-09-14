package vm

import (
	"errors"

	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
	MaxFrames   = 1024
)

var (
	errStackOverflow = errors.New("stack overflow")
)

func NewGlobals() object.Objects {
	return make(object.Objects, GlobalsSize)
}

func Make(b compiler.Bytecode, globals object.Objects) VM {
	return &virtualMachine{
		b:       b,
		stack:   make(object.Objects, StackSize),
		globals: globals,
		sp:      0,
		frames:  NewCallFrame(b, MaxFrames),
	}
}

func New(b compiler.Bytecode) VM {
	return Make(b, NewGlobals())
}

type VM interface {
	Run() error
	StackTop() object.Object
	LastPopped() object.Object
}

// virtualMachine : implement VM
type virtualMachine struct {
	b       compiler.Bytecode
	stack   object.Objects
	globals object.Objects
	sp      int // top stack [sp - 1]
	frames  CallFrame
}

func (this *virtualMachine) decodeUint16(ip int, ins code.Instructions) uint16 {
	return code.DecodeUint16(ins[ip+1:])
}

func (this *virtualMachine) fetchUint16(ip int, ins code.Instructions) uint16 {
	v := this.decodeUint16(ip, ins)
	this.frames.incrby(2)
	return v
}

func (this *virtualMachine) Run() error {
	var ip int
	var ins code.Instructions

	for !this.frames.eof() {
		this.frames.incr()
		ip = this.frames.ip()
		ins = this.frames.Instructions()
		consts := this.frames.Constants()
		op := code.Opcode(ins[ip])
		switch op {
		case code.OpSetGlobal:
			{
				idx := this.fetchUint16(ip, ins)
				this.globals[idx] = this.pop() // bind
			}
		case code.OpGetGlobal:
			{
				idx := this.fetchUint16(ip, ins)
				// resolve
				if err := this.push(this.globals[idx]); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpJump:
			{
				pos := this.decodeUint16(ip, ins)
				// in a loop that increments ip with each iteration
				// we need to set ip to the offset right before the one we want
				this.frames.jmp(int(pos - 1))
			}
		case code.OpJumpWhenFalse:
			{
				pos := this.fetchUint16(ip, ins)
				cond := this.pop()
				if !cond.True() {
					this.frames.jmp(int(pos - 1))
				}
			}
		case code.OpConst:
			{
				idx := this.fetchUint16(ip, ins)
				err := this.push(consts[idx])
				if nil != err {
					return function.NewError(err)
				}
			}
		case code.OpArray:
			{
				sz := int(this.fetchUint16(ip, ins))
				arr := this.doArray(sz)
				if err := this.push(arr); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpHash:
			{
				sz := int(this.fetchUint16(ip, ins))
				h, err := this.doHash(sz)
				if nil != err {
					return function.NewError(err)
				}
				if err := this.push(h); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpPop:
			{
				this.pop()
			}
		case code.OpTrue:
			{
				if err := this.push(object.True); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpFalse:
			{
				if err := this.push(object.False); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpNull:
			{
				if err := this.push(object.Nil); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpNot:
			{
				if err := this.doPrefix(object.FnNot); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpNeg:
			{
				if err := this.doPrefix(object.FnNeg); nil != err {
					return function.NewError(err)
				}
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
			{
				if err := this.doInfix(op); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpIndex:
			{
				idx := this.pop()
				left := this.pop()
				if err := this.doIndex(left, idx); nil != err {
					return function.NewError(err)
				}
			}
		case code.OpCall:
			{
				fn, err := this.top().AsByteFunc()
				if nil != err {
					return function.NewError(err)
				}
				frame := NewFrame(fn)
				this.frames.push(frame)
			}
		case code.OpReturn:
			{
				returnValue := this.pop()
				this.frames.pop()
				this.pop()
				if err := this.push(returnValue); nil != err {
					return function.NewError(err)
				}
			}
		}
	}
	return nil
}

func (this *virtualMachine) doHash(sz int) (object.Object, error) {
	h := object.HashMap{}
	for i := 0; i < sz; i++ {
		v := this.pop()
		k := this.pop()

		pair := object.HashPair{Key: k, Value: v}

		key, err := k.Hash()
		if nil != err {
			return nil, function.NewError(err)
		}
		h[*key] = &pair
	}
	return object.NewHash(h), nil
}

func (this *virtualMachine) doArray(sz int) object.Object {
	arr := make(object.Objects, sz)
	for i := 0; i < sz; i++ {
		arr[sz-i-1] = this.pop()
	}
	return object.NewArray(arr)
}

func (this *virtualMachine) doPrefix(fn string) error {
	right := this.pop()
	if r, err := right.CallMember(fn, object.Objects{}); nil != err {
		return function.NewError(err)
	} else {
		return this.push(r)
	}
}

func (this *virtualMachine) doIndex(left object.Object, idx object.Object) error {
	if r, err := left.CallMember(object.FnIndex, object.Objects{idx}); nil != err {
		return function.NewError(err)
	} else {
		this.push(r)
		return nil
	}
}

func (this *virtualMachine) doInfix(op code.Opcode) error {
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
	top := this.top()
	this.sp--
	return top
}

func (this *virtualMachine) top() object.Object {
	return this.stack[this.sp-1]
}
