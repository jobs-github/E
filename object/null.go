package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func newNull() *Null {
	obj := &Null{}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

// Null : implement Object
type Null struct {
	defaultObject
}

func (this *Null) String() string {
	return toString(objectTypeNull)
}

func (this *Null) Dump() (interface{}, error) {
	return nil, nil
}

func (this *Null) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcNull(op, this)
}

func (this *Null) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Null) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Null) getType() ObjectType {
	return objectTypeNull
}

func (this *Null) equal(other Object) error {
	return other.equalNull(this)
}

func (this *Null) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Null) equalNull(other *Null) error {
	return nil
}

func (this *Null) andInteger(left *Integer) Object {
	if 0 == left.Value {
		return left
	}
	return Nil
}

func (this *Null) andBoolean(left *Boolean) Object {
	if false == left.Value {
		return left
	}
	return Nil
}

func (this *Null) orInteger(left *Integer) Object {
	if 0 != left.Value {
		return left
	}
	return Nil
}

func (this *Null) orBoolean(left *Boolean) Object {
	if false != left.Value {
		return left
	}
	return Nil
}

func (this *Null) calcInteger(op *token.Token, left *Integer) (Object, error) {
	switch op.Type {
	case token.LT:
		return ToBoolean(false), nil
	case token.LEQ:
		return ToBoolean(false), nil
	case token.GT:
		return ToBoolean(true), nil
	case token.GEQ:
		return ToBoolean(true), nil
	case token.EQ:
		return ToBoolean(false), nil
	case token.NEQ:
		return ToBoolean(true), nil
	case token.AND:
		return this.andInteger(left), nil
	case token.OR:
		return this.orInteger(left), nil
	default:
		return Nil, unsupportedOp(function.GetFunc(), op, this)
	}
}

func (this *Null) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	switch op.Type {
	case token.AND:
		return this.andBoolean(left), nil
	case token.OR:
		return this.orBoolean(left), nil
	default:
		return this.calcInteger(op, toInteger(left.Value))
	}
}

func (this *Null) calcNull(op *token.Token, left *Null) (Object, error) {
	switch op.Type {
	case token.LT:
		return ToBoolean(false), nil
	case token.LEQ:
		return ToBoolean(true), nil
	case token.GT:
		return ToBoolean(false), nil
	case token.GEQ:
		return ToBoolean(true), nil
	case token.EQ:
		return ToBoolean(true), nil
	case token.NEQ:
		return ToBoolean(false), nil
	case token.AND:
		return this, nil
	case token.OR:
		return this, nil
	default:
		return Nil, unsupportedOp(function.GetFunc(), op, this)
	}
}

// builtin
func (this *Null) builtinNot(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return False, fmt.Errorf("not() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return True, nil
}
