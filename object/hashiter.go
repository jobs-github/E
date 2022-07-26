package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewHashIterator(h *Hash) Object {
	keys := h.Pairs.keys()
	obj := &HashIterator{h: h, keys: keys, offset: 0, sz: len(keys)}
	obj.fns = objectBuiltins{
		FnNot:   obj.builtinNot,
		"next":  obj.builtinNext,
		"key":   obj.builtinKey,
		"value": obj.builtinValue,
	}
	return obj
}

// HashIterator : implement Object
type HashIterator struct {
	h      *Hash
	keys   Objects
	offset int
	sz     int
	fns    objectBuiltins
}

func (this *HashIterator) String() string {
	return toString(objectTypeHashIter)
}

func (this *HashIterator) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *HashIterator) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *HashIterator) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcHashIter(op, this)
}

func (this *HashIterator) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *HashIterator) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *HashIterator) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *HashIterator) True() bool {
	return false
}

func (this *HashIterator) getType() ObjectType {
	return objectTypeHashIter
}

func (this *HashIterator) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *HashIterator) equal(other Object) error {
	return other.equalHashIter(this)
}

func (this *HashIterator) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalArrayIter(other *ArrayIterator) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *HashIterator) equalHashIter(other *HashIterator) error {
	if this.sz != other.sz {
		return fmt.Errorf("size mismatch, this: %v, other: %v", this.sz, other.sz)
	}
	if this.offset != other.offset {
		return fmt.Errorf("offset mismatch, this: %v, other: %v", this.offset, other.offset)
	}
	if err := this.h.equal(other.h); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *HashIterator) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *HashIterator) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *HashIterator) builtinNext(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("HashIterator next() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.offset < this.sz-1 {
		this.offset = this.offset + 1
		return this, nil
	}
	return Nil, nil
}

func (this *HashIterator) pair() (*HashPair, error) {
	key := this.keys[this.offset]
	k, err := key.Hash()
	if nil != err {
		return nil, function.NewError(err)
	}
	v, ok := this.h.Pairs.get(k)
	if !ok {
		return nil, fmt.Errorf("key `%v` missing in hash", key.String())
	}
	return v, nil
}

func (this *HashIterator) builtinKey(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("key() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.offset < 0 || this.offset > this.sz-1 {
		return Nil, fmt.Errorf("index out of range , size: %v, offset: %v, (`%v`)", this.sz, this.offset, this.String())
	}
	pair, err := this.pair()
	if nil != err {
		return Nil, function.NewError(err)
	}
	return pair.Key, nil
}

func (this *HashIterator) builtinValue(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("index() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if this.offset > this.sz-1 {
		return Nil, fmt.Errorf("index out of range , size: %v, offset: %v, (`%v`)", this.sz, this.offset, this.String())
	}
	pair, err := this.pair()
	if nil != err {
		return Nil, function.NewError(err)
	}
	return pair.Value, nil
}

func (this *HashIterator) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
