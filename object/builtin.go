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
	Fn   BuiltinFunction
	Name string
	fns  objectBuiltins
}

func (this *Builtin) String() string {
	return fmt.Sprintf("<built-in function %v>", this.Name)
}

func (this *Builtin) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Builtin) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
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

func (this *Builtin) True() bool {
	return false
}

func (this *Builtin) AsState() (*State, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Builtin) AsByteFunc() (*ByteFunc, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Builtin) getType() ObjectType {
	return objectTypeBuiltin
}

func (this *Builtin) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *Builtin) equal(other Object) error {
	return other.equalBuiltin(this)
}

func (this *Builtin) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalBuiltin(other *Builtin) error {
	if this.Name != other.Name {
		return fmt.Errorf("name mismatch, this: %v, other: %v", this.Name, other.Name)
	}
	return nil
}

func (this *Builtin) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalByteFunc(other *ByteFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Builtin) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *Builtin) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcByteFunc(op *token.Token, left *ByteFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Builtin) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *Builtin) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
