package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func newBoolean(v bool) *Boolean {
	obj := &Boolean{
		Value: v,
	}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
		FnNeg: obj.builtinNeg,
		FnInt: obj.builtinInt,
	}
	return obj
}

func NewBoolean(v bool) Object {
	return newBoolean(v)
}

// Boolean : implement Object
type Boolean struct {
	defaultObject
	Value bool
	fns   objectBuiltins
}

func (this *Boolean) String() string {
	return fmt.Sprintf("%v", this.Value)
}

func (this *Boolean) Hash() (*HashKey, error) {
	return &HashKey{Type: this.getType(), Value: uint64(toInt64(this.Value))}, nil
}

func (this *Boolean) Dump() (interface{}, error) {
	return this.Value, nil
}

func (this *Boolean) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcBoolean(op, this)
}

func (this *Boolean) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Boolean) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Boolean) True() bool {
	return this.Value
}

func (this *Boolean) getType() ObjectType {
	return objectTypeBoolean
}

func (this *Boolean) asInteger() (int64, error) {
	return toInt64(this.Value), nil
}

func (this *Boolean) equal(other Object) error {
	return other.equalBoolean(this)
}

func (this *Boolean) equalBoolean(other *Boolean) error {
	if this.Value != other.Value {
		return fmt.Errorf("value mismatch, this: %v, other: %v", this.Value, other.Value)
	}
	return nil
}

func (this *Boolean) calcInteger(op *token.Token, left *Integer) (Object, error) {
	right := toInteger(this.Value)
	return right.calcInteger(op, left)
}

func (this *Boolean) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	switch op.Type {
	case token.ADD:
		return NewInteger(toInt64(left.Value) + toInt64(this.Value)), nil
	case token.SUB:
		return NewInteger(toInt64(left.Value) - toInt64(this.Value)), nil
	case token.MUL:
		return NewInteger(toInt64(left.Value) * toInt64(this.Value)), nil
	case token.DIV:
		return NewInteger(toInt64(left.Value) / toInt64(this.Value)), nil
	case token.MOD:
		return NewInteger(toInt64(left.Value) % toInt64(this.Value)), nil
	case token.LT:
		return ToBoolean(toInt64(left.Value) < toInt64(this.Value)), nil
	case token.LEQ:
		return ToBoolean(toInt64(left.Value) <= toInt64(this.Value)), nil
	case token.GT:
		return ToBoolean(toInt64(left.Value) > toInt64(this.Value)), nil
	case token.GEQ:
		return ToBoolean(toInt64(left.Value) >= toInt64(this.Value)), nil
	case token.EQ:
		return ToBoolean(left.Value == this.Value), nil
	case token.NEQ:
		return ToBoolean(left.Value != this.Value), nil
	case token.AND:
		return ToBoolean(left.Value && this.Value), nil
	case token.OR:
		return ToBoolean(left.Value || this.Value), nil
	default:
		return Nil, unsupportedOp(function.GetFunc(), op, this)
	}
}

func (this *Boolean) calcNull(op *token.Token, left *Null) (Object, error) {
	return infixNull(op, this, function.GetFunc())
}

// builtin
func (this *Boolean) builtinNot(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return newBoolean(false), fmt.Errorf("not() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.Value {
		return newBoolean(false), nil
	} else {
		return newBoolean(true), nil
	}
}

func (this *Boolean) builtinNeg(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return NewInteger(0), fmt.Errorf("neg() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.Value {
		return NewInteger(-1), nil
	} else {
		return NewInteger(0), nil
	}
}

func (this *Boolean) builtinInt(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return NewInteger(0), fmt.Errorf("int() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.Value {
		return NewInteger(1), nil
	} else {
		return NewInteger(0), nil
	}
}
