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
	Obj  Object
	Name string
	Fn   objectFn
	fns  objectBuiltins
}

func (this *ObjectFunc) String() string {
	return fmt.Sprintf("<built-in method %v of %v object>", this.Name, Typeof(this.Obj))
}

func (this *ObjectFunc) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ObjectFunc) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
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

func (this *ObjectFunc) True() bool {
	return false
}

func (this *ObjectFunc) Return() (bool, Object) {
	return false, nil
}

func (this *ObjectFunc) Break() (bool, int) {
	return false, 0
}

func (this *ObjectFunc) getType() ObjectType {
	return objectTypeObjectFunc
}

func (this *ObjectFunc) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *ObjectFunc) equal(other Object) error {
	return other.equalObjectFunc(this)
}

func (this *ObjectFunc) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
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

func (this *ObjectFunc) equalArrayIter(other *ArrayIterator) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) equalHashIter(other *HashIterator) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ObjectFunc) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *ObjectFunc) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ObjectFunc) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *ObjectFunc) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
