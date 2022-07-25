package object

import (
	"fmt"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/token"
)

func NewArrayIterator(arr *Array) Object {
	obj := &ArrayIterator{arr: arr, offset: 0, sz: len(arr.Items)}
	obj.fns = objectBuiltins{
		FnNot:   obj.builtinNot,
		"next":  obj.builtinNext,
		"key":   obj.builtinKey,
		"value": obj.builtinValue,
	}
	return obj
}

// ArrayIterator : implement Object
type ArrayIterator struct {
	arr    *Array
	offset int
	sz     int
	fns    objectBuiltins
}

func (this *ArrayIterator) String() string {
	return toString(objectTypeArrayIter)
}

func (this *ArrayIterator) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ArrayIterator) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ArrayIterator) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcArrayIter(op, this)
}

func (this *ArrayIterator) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *ArrayIterator) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *ArrayIterator) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *ArrayIterator) True() bool {
	return false
}

func (this *ArrayIterator) Return() (bool, Object) {
	return false, nil
}

func (this *ArrayIterator) Break() (bool, int) {
	return false, 0
}

func (this *ArrayIterator) getType() ObjectType {
	return objectTypeArrayIter
}

func (this *ArrayIterator) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *ArrayIterator) equal(other Object) error {
	return other.equalArrayIter(this)
}

func (this *ArrayIterator) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) equalArrayIter(other *ArrayIterator) error {
	if this.sz != other.sz {
		return fmt.Errorf("size mismatch, this: %v, other: %v", this.sz, other.sz)
	}
	if this.offset != other.offset {
		return fmt.Errorf("offset mismatch, this: %v, other: %v", this.offset, other.offset)
	}
	if err := this.arr.equal(other.arr); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *ArrayIterator) equalHashIter(other *HashIterator) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *ArrayIterator) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *ArrayIterator) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *ArrayIterator) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *ArrayIterator) builtinNext(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("next() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.offset < this.sz-1 {
		this.offset = this.offset + 1
		return this, nil
	}
	return Nil, nil
}

func (this *ArrayIterator) builtinKey(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("key() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.offset < 0 || this.offset > this.sz-1 {
		return Nil, fmt.Errorf("index out of range , size: %v, offset: %v, (`%v`)", this.sz, this.offset, this.String())
	}
	return NewInteger(int64(this.offset)), nil
}

func (this *ArrayIterator) builtinValue(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("index() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.offset > this.sz-1 {
		return Nil, fmt.Errorf("index out of range , size: %v, offset: %v, (`%v`)", this.sz, this.offset, this.String())
	}
	return this.arr.Items[this.offset], nil
}

func (this *ArrayIterator) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
