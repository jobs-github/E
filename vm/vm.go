package vm

import (
	"errors"
	"fmt"

	"github.com/jobs-github/escript/builtin"
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
	errNotCallable   = errors.New("not callable")
)

func NewGlobals() object.Objects {
	return make(object.Objects, GlobalsSize)
}

func Make(b compiler.Bytecode, c object.Objects, globals object.Objects) VM {
	return &virtualMachine{
		b:         b,
		constants: c,
		stack:     make(object.Objects, StackSize),
		globals:   globals,
		sp:        0,
		frames:    NewCallFrame(b, MaxFrames),
		ip:        -1,
		ins:       nil,
	}
}

func New(b compiler.Bytecode, c object.Objects) VM {
	return Make(b, c, NewGlobals())
}

type VM interface {
	Run() error
	StackTop() object.Object
	LastPopped() object.Object
}

// virtualMachine : implement VM
type virtualMachine struct {
	b         compiler.Bytecode
	constants object.Objects
	stack     object.Objects
	globals   object.Objects
	sp        int // top stack [sp - 1]
	frames    CallFrame
	ip        int
	ins       code.Instructions
}

func (this *virtualMachine) decodeUint16() uint16 {
	return code.DecodeUint16(this.ins[this.ip+1:])
}

func (this *virtualMachine) decodeUint8() uint8 {
	return code.DecodeUint8(this.ins[this.ip+1:])
}

func (this *virtualMachine) fetchUint16() uint16 {
	v := this.decodeUint16()
	this.frames.incrby(2)
	return v
}

func (this *virtualMachine) fetchUint8() uint8 {
	v := this.decodeUint8()
	this.frames.incr()
	return v
}

func (this *virtualMachine) fetchClosure() (uint16, int) {
	idx := this.fetchUint16()
	frees := code.DecodeUint8(this.ins[this.ip+3:])
	this.frames.incr()
	return idx, int(frees)
}

func (this *virtualMachine) Run() error {
	for !this.frames.eof() {
		this.frames.incr()
		this.ip = this.frames.ip()
		this.ins = this.frames.instructions()
		op := code.Opcode(this.ins[this.ip])
		switch op {
		case code.OpConst:
			{
				idx := this.fetchUint16()
				err := this.push(this.constants[idx])
				if nil != err {
					return err
				}
			}
		case code.OpSetGlobal:
			{
				idx := this.fetchUint16()
				this.globals[idx] = this.pop() // bind
			}
		case code.OpGetGlobal:
			{
				idx := this.fetchUint16()
				// resolve
				if err := this.push(this.globals[idx]); nil != err {
					return err
				}
			}
		case code.OpSetLocal: // pop the stack and fill the hole
			{
				localIndex := this.fetchUint8()
				idx := this.frames.basePointer() + int(localIndex)
				this.stack[idx] = this.pop()
			}
		case code.OpGetLocal:
			{
				localIndex := this.fetchUint8()
				idx := this.frames.basePointer() + int(localIndex)
				if err := this.push(this.stack[idx]); nil != err {
					return err
				}
			}
		case code.OpIncLocal:
			{
				localIndex := this.fetchUint8()
				idx := this.frames.basePointer() + int(localIndex)
				this.stack[idx].Incr()
			}
		case code.OpJump:
			{
				pos := this.decodeUint16()
				// in a loop that increments ip with each iteration
				// we need to set ip to the offset right before the one we want
				this.frames.jmp(int(pos - 1))
			}
		case code.OpJumpWhenFalse:
			{
				pos := this.fetchUint16()
				cond := this.pop()
				if !cond.True() {
					this.frames.jmp(int(pos - 1))
				}
			}
		case code.OpArrayLen:
			{
				if err := this.doArrayLen(); nil != err {
					return err
				}
			}
		case code.OpArrayNew:
			{
				if err := this.doArrayNew(); nil != err {
					return err
				}
			}
		case code.OpArrayAppend:
			{
				if err := this.doArrayAppend(); nil != err {
					return err
				}
			}
		case code.OpArraySet:
			{
				if err := this.doArraySet(); nil != err {
					return err
				}
			}
		case code.OpGetBuiltin: // pair with OpCall
			{
				if err := this.doGetBuiltin(); nil != err {
					return err
				}
			}
		case code.OpGetObjectFn: // pair with OpCall
			{
				if err := this.doGetObjectFn(); nil != err {
					return err
				}
			}
		case code.OpGetFree:
			{
				if err := this.doGetFree(); nil != err {
					return err
				}
			}
		case code.OpGetLambda:
			{
				if err := this.push(this.frames.current().fn); nil != err {
					return err
				}
			}
		case code.OpClosure: // pair with OpCall
			{
				if err := this.doClosure(); nil != err {
					return err
				}
			}
		case code.OpCall:
			{
				if err := this.doCall(); nil != err {
					return err
				}
			}
		case code.OpReturn:
			{
				if err := this.doReturn(); nil != err {
					return err
				}
			}
		case code.OpArray:
			{
				if err := this.doArray(); nil != err {
					return err
				}
			}
		case code.OpHash:
			{
				if err := this.doHash(); nil != err {
					return err
				}
			}
		case code.OpPop:
			{
				this.pop()
			}
		case code.OpTrue:
			{
				if err := this.push(object.True); nil != err {
					return err
				}
			}
		case code.OpFalse:
			{
				if err := this.push(object.False); nil != err {
					return err
				}
			}
		case code.OpNull:
			{
				if err := this.push(object.Nil); nil != err {
					return err
				}
			}
		case code.OpNot:
			{
				if err := this.doPrefix(object.FnNot); nil != err {
					return err
				}
			}
		case code.OpNeg:
			{
				if err := this.doPrefix(object.FnNeg); nil != err {
					return err
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
					return err
				}
			}
		case code.OpIndex:
			{
				if err := this.doIndex(); nil != err {
					return err
				}
			}
		}
	}
	return nil
}

func (this *virtualMachine) doArrayLen() error {
	arr, err := this.pop().AsArray()
	if nil != err {
		return err
	}
	if err := this.push(object.NewInteger(int64(len(arr.Items)))); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doArrayNew() error {
	flag := this.fetchUint8()
	arr, err := this.pop().AsArray()
	if nil != err {
		return err
	}
	if err := this.push(arr); nil != err {
		return err
	}
	if err := this.push(arr.New(flag)); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doArrayAppend() error {
	localIndex := this.fetchUint8()
	idx := this.frames.basePointer() + int(localIndex)
	arr, err := this.stack[idx].AsArray()
	if nil != err {
		return err
	}
	item := this.pop()
	v := this.pop()
	if v.True() {
		arr.Append(item)
	}
	return nil
}

func (this *virtualMachine) doArraySet() error {
	localIndex := this.fetchUint8()
	idx := this.frames.basePointer() + int(localIndex)
	arr, err := this.stack[idx].AsArray()
	if nil != err {
		return err
	}
	i := this.pop()
	v := this.pop()
	if err := arr.Set(i, v); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doGetBuiltin() error {
	idx := this.fetchUint8()
	builtinFn := builtin.Resolve(int(idx))
	// object.Builtin
	if err := this.push(builtinFn); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doGetObjectFn() error {
	idx := this.fetchUint8()
	obj := this.pop()
	fn := object.Resolve(int(idx))
	r, err := obj.GetMember(fn)
	if nil != err {
		return err
	}
	// object.ObjectFunc
	if err := this.push(r); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doGetFree() error {
	idx := this.fetchUint8()
	fn := this.frames.current().fn
	if err := this.push(fn.Free[idx]); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doClosure() error {
	// exec after a lot OpGetFree, refer to visitor.DoFn
	// idx := this.fetchUint16()
	// frees := int(this.fetchUint8(ip+2, ins))
	idx, frees := this.fetchClosure()
	fn, err := this.constants[idx].AsByteFunc()
	if nil != err {
		return err
	}
	// fetch free symbol from stack top
	freeSymbols := make(object.Objects, frees)
	for i := 0; i < frees; i++ {
		freeSymbols[i] = this.stack[this.sp-frees+i]
	}
	// clean up the stack
	this.sp = this.sp - frees
	if err := this.push(object.NewClosure(fn, freeSymbols)); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doCall() error {
	args := this.fetchUint8()
	obj := this.stack[this.sp-1-int(args)]

	if object.IsBuiltin(obj) || object.IsObjectFunc(obj) {
		arguments := this.stack[this.sp-int(args) : this.sp]
		r, err := obj.Call(arguments)
		if nil != err {
			return err
		}
		this.sp = this.sp - int(args) - 1
		if err := this.push(r); nil != err {
			return err
		}
		return nil
	}

	if object.IsClosure(obj) {
		fn, _ := obj.AsClosure()
		if args != uint8(fn.Fn.Locals) {
			err := fmt.Errorf("wrong number of arguments: want=%v, got=%v", fn.Fn.Locals, args)
			return err
		}
		frame := NewFrame(fn, this.sp-int(args))
		// set env
		this.frames.push(frame)
		this.sp = frame.bp + fn.Fn.Locals // reserverd for local bindings
		return nil
	}

	return errNotCallable
}

func (this *virtualMachine) doReturn() error {
	returnValue := this.pop()
	// recover env
	frame := this.frames.pop()
	this.sp = frame.bp - 1 // frame.bp point to the just-executed function on the stack

	if err := this.push(returnValue); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doHash() error {
	sz := int(this.fetchUint16())
	h := object.HashMap{}
	for i := 0; i < sz; i++ {
		v := this.pop()
		k := this.pop()

		pair := object.HashPair{Key: k, Value: v}

		key, err := k.Hash()
		if nil != err {
			return err
		}
		h[*key] = &pair
	}
	if err := this.push(object.NewHash(h)); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doArray() error {
	sz := int(this.fetchUint16())
	arr := make(object.Objects, sz)
	for i := 0; i < sz; i++ {
		arr[sz-i-1] = this.pop()
	}
	r := object.NewArray(arr)
	if err := this.push(r); nil != err {
		return err
	}
	return nil
}

func (this *virtualMachine) doPrefix(fn string) error {
	right := this.pop()
	if r, err := right.CallMember(fn, object.Objects{}); nil != err {
		return err
	} else {
		return this.push(r)
	}
}

func (this *virtualMachine) doIndex() error {
	idx := this.pop()
	left := this.pop()
	if r, err := left.CallMember(object.FnIndex, object.Objects{idx}); nil != err {
		return err
	} else {
		this.push(r)
		return nil
	}
}

func (this *virtualMachine) doInfix(op code.Opcode) error {
	t, err := code.InfixToken(op)
	if nil != err {
		return err
	}
	right := this.pop()
	left := this.pop()
	r, err := left.Calc(t, right)
	if nil != err {
		return err
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
