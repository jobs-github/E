package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewClosure(fn *ByteFunc) *Closure {
	obj := &Closure{Fn: fn}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

// Closure : implement Object
type Closure struct {
	defaultObject
	Fn   *ByteFunc
	Free Objects
}

func (this *Closure) String() string {
	return fmt.Sprintf("closure[%p]", this)
}

func (this *Closure) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcClosure(op, this)
}

func (this *Closure) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Closure) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Closure) True() bool {
	return false
}

func (this *Closure) AsClosure() (*Closure, error) {
	return this, nil
}

func (this *Closure) getType() ObjectType {
	return objectTypeClosure
}

func (this *Closure) equal(other Object) error {
	return other.equalClosure(this)
}

func (this *Closure) equalClosure(other *Closure) error {
	// TODO
	return nil
}

func (this *Closure) calcClosure(op *token.Token, left *Closure) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *Closure) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
