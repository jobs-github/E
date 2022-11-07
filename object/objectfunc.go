package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewObjectFunc(obj Object, name string, fn objectFn) Object {
	f := &ObjectFunc{
		Obj:  obj,
		Name: name,
		Fn:   fn,
	}
	f.fns = objectBuiltins{
		FnNot: f.builtinNot,
	}
	return f
}

// ObjectFunc : implement Object
type ObjectFunc struct {
	defaultObject
	Obj  Object
	Name string
	Fn   objectFn
	fns  objectBuiltins
}

func (this *ObjectFunc) String() string {
	return fmt.Sprintf("<built-in method %v of %v object>", this.Name, Typeof(this.Obj))
}

func (this *ObjectFunc) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcObjectFunc(op, this)
}

func (this *ObjectFunc) Call(args Objects) (Object, error) {
	return this.Fn(args)
}

func (this *ObjectFunc) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *ObjectFunc) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *ObjectFunc) getType() ObjectType {
	return objectTypeObjectFunc
}

func (this *ObjectFunc) equal(other Object) error {
	return other.equalObjectFunc(this)
}

func (this *ObjectFunc) equalObjectFunc(other *ObjectFunc) error {
	if this.Name != other.Name {
		return fmt.Errorf("name mismatch, this: %v, other: %v", this.Name, other.Name)
	}
	if err := other.Obj.equal(this.Obj); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *ObjectFunc) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *ObjectFunc) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
