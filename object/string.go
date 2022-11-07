package object

import (
	"fmt"
	"strconv"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewString(v string) Object {
	obj := &String{
		Value: v,
	}
	obj.fns = objectBuiltins{
		FnLen:   obj.builtinLen,
		FnIndex: obj.builtinIndex,
		FnNot:   obj.builtinNot,
		FnInt:   obj.builtinInt,
	}
	return obj
}

// String : implement Object
type String struct {
	defaultObject
	Value string
	fns   objectBuiltins
}

func (this *String) String() string {
	return this.Value
}

func (this *String) Hash() (*HashKey, error) {
	return &HashKey{Type: this.getType(), Value: hash64([]byte(this.Value))}, nil
}

func (this *String) Dump() (interface{}, error) {
	return this.Value, nil
}

func (this *String) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcString(op, this)
}

func (this *String) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *String) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *String) True() bool {
	if "" == this.Value {
		return false
	}
	return true
}

func (this *String) getType() ObjectType {
	return objectTypeString
}

func (this *String) equal(other Object) error {
	return other.equalString(this)
}

func (this *String) equalString(other *String) error {
	if this.Value != other.Value {
		return fmt.Errorf("value mismatch, this: %v, other: %v", this.Value, other.Value)
	}
	return nil
}

func (this *String) calcString(op *token.Token, left *String) (Object, error) {
	switch op.Type {
	case token.ADD:
		return NewString(fmt.Sprintf("%v%v", left.Value, this.Value)), nil
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

func (this *String) and(left *String) (Object, error) {
	if "" == left.Value {
		return left, nil
	}
	return this, nil
}

func (this *String) or(left *String) (Object, error) {
	if "" != left.Value {
		return left, nil
	}
	return this, nil
}

// builtin
func (this *String) builtinLen(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("len() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return NewInteger(int64(len(this.Value))), nil
}

func (this *String) builtinIndex(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("index() takes exactly one argument (%v given)", argc)
	}
	if "" == this.Value {
		return Nil, function.NewError(errStringEmpty)
	}
	idx, err := args[0].asInteger()
	if nil != err {
		return Nil, function.NewError(err)
	}
	return indexofString(this.Value, idx)
}

func (this *String) builtinNot(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("len() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if "" == this.Value {
		return True, nil
	} else {
		return False, nil
	}
}

func (this *String) builtinInt(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("int() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	v, err := strconv.ParseInt(this.Value, 10, 64)
	if nil != err {
		return Nil, function.NewError(err)
	}
	return NewInteger(v), nil
}
