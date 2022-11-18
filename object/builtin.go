package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewBuiltin(fn BuiltinFunction, name string) Object {
	obj := &Builtin{
		Fn:   fn,
		Name: name,
	}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

// Builtin : implement Object
type Builtin struct {
	defaultObject
	Fn   BuiltinFunction
	Name string
}

func (this *Builtin) String() string {
	return fmt.Sprintf("<built-in function %v>", this.Name)
}

func (this *Builtin) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcBuiltin(op, this)
}

func (this *Builtin) Call(args Objects) (Object, error) {
	return this.Fn(args)
}

func (this *Builtin) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Builtin) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Builtin) getType() ObjectType {
	return objectTypeBuiltin
}

func (this *Builtin) equal(other Object) error {
	return other.equalBuiltin(this)
}

func (this *Builtin) equalBuiltin(other *Builtin) error {
	if this.Name != other.Name {
		return fmt.Errorf("name mismatch, this: %v, other: %v", this.Name, other.Name)
	}
	return nil
}

func (this *Builtin) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *Builtin) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
