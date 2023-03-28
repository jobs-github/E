package vm

import (
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/compiler"
	"github.com/jobs-github/escript/object"
)

func NewCallFrame(b compiler.Bytecode, frameSize int) CallFrame {
	fn := object.NewByteFn(b.Instructions(), 0)
	mainFrame := NewFrame(object.NewClosure(fn, nil), 0)
	frames := make([]*Frame, frameSize)
	frames[0] = mainFrame
	return &callFrame{
		frames:     frames,
		frameIndex: 1,
	}
}

func NewFrame(fn *object.Closure, basePointer int) *Frame {
	return &Frame{
		fn: fn,
		ip: -1,
		bp: basePointer,
	}
}

type Frame struct {
	fn *object.Closure
	ip int // => bytecode
	bp int // => stack
}

func (this *Frame) reset() {
	this.ip = -1
}

type CallFrame interface {
	reset()
	instructions() code.Instructions
	ip() int
	basePointer() int
	eof() bool
	jmp(ip int)
	incr()
	incrby(sz int)
	current() *Frame
	push(f *Frame)
	pop() *Frame
}

// callFrame : implement CallFrame
type callFrame struct {
	frames     []*Frame
	frameIndex int
}

func (this *callFrame) reset() {
	this.frameIndex = 1
	this.current().reset()
}

func (this *callFrame) instructions() code.Instructions {
	return this.current().fn.Fn.Ins
}

func (this *callFrame) ip() int {
	return this.current().ip
}

func (this *callFrame) basePointer() int {
	return this.current().bp
}

func (this *callFrame) eof() bool {
	return this.ip() >= len(this.instructions())-1
}

func (this *callFrame) jmp(ip int) {
	this.current().ip = ip
}

func (this *callFrame) incr() {
	this.current().ip++
}

func (this *callFrame) incrby(sz int) {
	this.current().ip += sz
}

func (this *callFrame) current() *Frame {
	return this.frames[this.frameIndex-1]
}

func (this *callFrame) push(f *Frame) {
	this.frames[this.frameIndex] = f
	this.frameIndex++
}

func (this *callFrame) pop() *Frame {
	this.frameIndex--
	return this.frames[this.frameIndex]
}
