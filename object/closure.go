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
	Fn   *ByteFunc
	Free Objects
	fns  objectBuiltins
}

func (this *Closure) String() string {
	return fmt.Sprintf("closure[%p]", this)
}

func (this *Closure) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Closure) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Closure) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcClosure(op, this)
}

func (this *Closure) Call(args Objects) (Object, error) {
	// TODO
	return nil, nil
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

func (this *Closure) AsState() (*State, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Closure) AsByteFunc() (*ByteFunc, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Closure) AsClosure() (*Closure, error) {
	return this, nil
}

func (this *Closure) getType() ObjectType {
	return objectTypeClosure
}

func (this *Closure) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *Closure) equal(other Object) error {
	return other.equalClosure(this)
}

func (this *Closure) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalByteFunc(other *ByteFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) equalClosure(other *Closure) error {
	// TODO
	return nil
}

func (this *Closure) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Closure) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcByteFunc(op *token.Token, left *ByteFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcClosure(op *token.Token, left *Closure) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *Closure) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Closure) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *Closure) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
