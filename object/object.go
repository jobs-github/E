package object

import (
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

type Object interface {
	String() string
	Hash() (*HashKey, error)
	Dump() (interface{}, error)

	// public
	Calc(op *token.Token, right Object) (Object, error)
	Call(args Objects) (Object, error)
	CallMember(name string, args Objects) (Object, error)
	GetMember(name string) (Object, error)
	True() bool

	getType() ObjectType
	asInteger() (int64, error)
	equal(other Object) error
	equalInteger(other *Integer) error
	equalString(other *String) error
	equalBoolean(other *Boolean) error
	equalNull(other *Null) error
	equalArray(other *Array) error
	equalHash(other *Hash) error
	equalBuiltin(other *Builtin) error
	equalFunction(other *Function) error
	equalObjectFunc(other *ObjectFunc) error
	equalArrayIter(other *ArrayIterator) error
	equalHashIter(other *HashIterator) error
	// calc
	calcInteger(op *token.Token, left *Integer) (Object, error)
	calcString(op *token.Token, left *String) (Object, error)
	calcBoolean(op *token.Token, left *Boolean) (Object, error)
	calcNull(op *token.Token, left *Null) (Object, error)
	calcArray(op *token.Token, left *Array) (Object, error)
	calcHash(op *token.Token, left *Hash) (Object, error)
	calcBuiltin(op *token.Token, left *Builtin) (Object, error)
	calcFunction(op *token.Token, left *Function) (Object, error)
	calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error)
	calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error)
	calcHashIter(op *token.Token, left *HashIterator) (Object, error)
}

type objectFn func(args Objects) (Object, error)
type objectBuiltins map[string]objectFn

func (this *objectBuiltins) get(name string) (objectFn, bool) {
	v, ok := (*this)[name]
	return v, ok
}

type Objects []Object

func (this *Objects) Append(item Object) {
	*this = append(*this, item)
}

func (this *Objects) first() (Object, error) {
	if len(*this) < 1 {
		return Nil, function.NewError(errListEmpty)
	}
	return (*this)[0], nil
}

func (this *Objects) last() (Object, error) {
	sz := len(*this)
	if sz < 1 {
		return Nil, function.NewError(errListEmpty)
	}
	return (*this)[sz-1], nil
}

func (this *Objects) tail() (Object, error) {
	sz := len(*this)
	if sz < 1 {
		return Nil, function.NewError(errListEmpty)
	}
	items := make(Objects, sz-1, sz-1)
	copy(items, (*this)[1:sz])
	return NewArray(items), nil
}

func (this *Objects) push(item Object) (Object, error) {
	this.Append(item)
	return NewArray(*this), nil
}
