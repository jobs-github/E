package object

import (
	"fmt"

	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewByteFn(ins code.Instructions, locals int) *ByteFunc {
	obj := &ByteFunc{Ins: ins, Locals: locals}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

func NewByteFunc(ins code.Instructions, locals int) Object {
	return NewByteFn(ins, locals)
}

// ByteFunc : implement Object
type ByteFunc struct {
	defaultObject
	Ins    code.Instructions
	Locals int
}

func (this *ByteFunc) String() string {
	return fmt.Sprintf("byte_func[%p]", this)
}

func (this *ByteFunc) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcByteFunc(op, this)
}

func (this *ByteFunc) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *ByteFunc) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *ByteFunc) AsByteFunc() (*ByteFunc, error) {
	return this, nil
}

func (this *ByteFunc) getType() ObjectType {
	return objectTypeByteFunc
}

func (this *ByteFunc) equal(other Object) error {
	return other.equalByteFunc(this)
}

func (this *ByteFunc) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalByteFunc(other *ByteFunc) error {
	// TODO
	return nil
}

func (this *ByteFunc) calcByteFunc(op *token.Token, left *ByteFunc) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *ByteFunc) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
