package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func newInteger(v int64) *Integer {
	obj := &Integer{
		Value: v,
	}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
		FnNeg: obj.builtinNeg,
		FnInt: obj.builtinInt,
	}
	return obj
}

func NewInteger(v int64) Object {
	return newInteger(v)
}

// Integer : implement Object
type Integer struct {
	Value int64
	fns   objectBuiltins
}

func (this *Integer) String() string {
	return fmt.Sprintf("%v", this.Value)
}

func (this *Integer) Hash() (*HashKey, error) {
	return &HashKey{Type: this.getType(), Value: uint64(this.Value)}, nil
}

func (this *Integer) Dump() (interface{}, error) {
	return this.Value, nil
}

func (this *Integer) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcInteger(op, this)
}

func (this *Integer) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *Integer) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Integer) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Integer) True() bool {
	if 0 == this.Value {
		return false
	}
	return true
}

func (this *Integer) AsState() (*State, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Integer) getType() ObjectType {
	return objectTypeInteger
}

func (this *Integer) asInteger() (int64, error) {
	return this.Value, nil
}

func (this *Integer) equal(other Object) error {
	return other.equalInteger(this)
}

func (this *Integer) equalInteger(other *Integer) error {
	if this.Value != other.Value {
		return fmt.Errorf("value mismatch, this: %v, other: %v", this.Value, other.Value)
	}
	return nil
}

func (this *Integer) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalByteFunc(other *ByteFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Integer) calcInteger(op *token.Token, left *Integer) (Object, error) {
	switch op.Type {
	case token.ADD:
		return NewInteger(left.Value + this.Value), nil
	case token.SUB:
		return NewInteger(left.Value - this.Value), nil
	case token.MUL:
		return NewInteger(left.Value * this.Value), nil
	case token.DIV:
		return NewInteger(left.Value / this.Value), nil
	case token.MOD:
		return NewInteger(left.Value % this.Value), nil
	case token.LT:
		return ToBoolean(left.Value < this.Value), nil
	case token.LEQ:
		return ToBoolean(left.Value <= this.Value), nil
	case token.GT:
		return ToBoolean(left.Value > this.Value), nil
	case token.GEQ:
		return ToBoolean(left.Value >= this.Value), nil
	case token.EQ:
		return ToBoolean(left.Value == this.Value), nil
	case token.NEQ:
		return ToBoolean(left.Value != this.Value), nil
	case token.AND:
		return this.and(left)
	case token.OR:
		return this.or(left)
	default:
		return Nil, unsupportedOp(function.GetFunc(), op, this)
	}
}

func (this *Integer) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Integer) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return this.calcInteger(op, toInteger(left.Value))
}

func (this *Integer) calcNull(op *token.Token, left *Null) (Object, error) {
	return infixNull(op, this, function.GetFunc())
}

func (this *Integer) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Integer) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Integer) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Integer) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Integer) calcByteFunc(op *token.Token, left *ByteFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Integer) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Integer) and(left *Integer) (Object, error) {
	if 0 == left.Value {
		return left, nil
	}
	return this, nil
}

func (this *Integer) or(left *Integer) (Object, error) {
	if 0 != left.Value {
		return left, nil
	}
	return this, nil
}

// builtin
func (this *Integer) builtinNot(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return newBoolean(false), fmt.Errorf("not() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if 0 == this.Value {
		return newBoolean(true), nil
	} else {
		return newBoolean(false), nil
	}
}

func (this *Integer) builtinNeg(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return NewInteger(0), fmt.Errorf("neg() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return NewInteger(-this.Value), nil
}

func (this *Integer) builtinInt(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return NewInteger(0), fmt.Errorf("int() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return this, nil
}
