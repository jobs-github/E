package object

import (
	"fmt"

	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewByteFn(ins code.Instructions) *ByteFunc {
	obj := &ByteFunc{Instructions: ins}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

func NewByteFunc(ins code.Instructions) Object {
	return NewByteFn(ins)
}

// ByteFunc : implement Object
type ByteFunc struct {
	Instructions code.Instructions
	fns          objectBuiltins
}

func (this *ByteFunc) String() string {
	return fmt.Sprintf("byte_func[%p]", this)
}

func (this *ByteFunc) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ByteFunc) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ByteFunc) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcByteFunc(op, this)
}

func (this *ByteFunc) Call(args Objects) (Object, error) {
	// TODO
	return nil, nil
}

func (this *ByteFunc) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *ByteFunc) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *ByteFunc) True() bool {
	return false
}

func (this *ByteFunc) AsState() (*State, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ByteFunc) getType() ObjectType {
	return objectTypeByteFunc
}

func (this *ByteFunc) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *ByteFunc) equal(other Object) error {
	return other.equalByteFunc(this)
}

func (this *ByteFunc) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) equalByteFunc(other *ByteFunc) error {
	// TODO
	return nil
}

func (this *ByteFunc) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ByteFunc) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcByteFunc(op *token.Token, left *ByteFunc) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *ByteFunc) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ByteFunc) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *ByteFunc) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
