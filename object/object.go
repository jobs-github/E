package object

import (
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

type Callable func() (Object, error)

type ObjectType uint8

type ObjectTypes []ObjectType

const (
	objectTypeBuiltin ObjectType = iota
	objectTypeInteger
	objectTypeString
	objectTypeBoolean
	objectTypeNull
	objectTypeFunction
	objectTypeByteFunc
	objectTypeClosure
	objectTypeArray
	objectTypeHash
	objectTypeObjectFunc
)

const (
	TypeHash    = "hash"
	TypeArray   = "array"
	TypeBool    = "boolean"
	TypeInt     = "integer"
	TypeStr     = "string"
	TypeBuiltin = "builtin"
)

const (
	FnLen   = "len"
	FnIndex = "index"
	FnNot   = "not"
	FnNeg   = "neg"
	FnInt   = "int"
	FnFirst = "first"
	FnLast  = "last"
	FnTail  = "tail"
	FnPush  = "push"
	FnKeys  = "keys"
)

var (
	Nil   = newNull()
	True  = newBoolean(true)
	False = newBoolean(false)
)

var (
	errStringEmpty = errors.New("string is empty")
	errListEmpty   = errors.New("list is empty")
	errArrayEmpty  = errors.New("array is empty")
	errHashEmpty   = errors.New("hash is empty")
)

var (
	objectTypeStrings = map[ObjectType]string{
		objectTypeBuiltin:    TypeBuiltin,
		objectTypeInteger:    TypeInt,
		objectTypeString:     TypeStr,
		objectTypeBoolean:    TypeBool,
		objectTypeNull:       token.Null,
		objectTypeFunction:   "function",
		objectTypeByteFunc:   "byte_func",
		objectTypeClosure:    "closure",
		objectTypeArray:      TypeArray,
		objectTypeHash:       TypeHash,
		objectTypeObjectFunc: "object_func",
	}
)

func IsString(v Object) bool {
	return v.getType() == objectTypeString
}

func IsNull(v Object) bool {
	return v.getType() == objectTypeNull
}

func IsInteger(v Object) bool {
	return v.getType() == objectTypeInteger
}

func IsBuiltin(v Object) bool {
	return v.getType() == objectTypeBuiltin
}

func IsObjectFunc(v Object) bool {
	return v.getType() == objectTypeObjectFunc
}

func IsClosure(v Object) bool {
	return v.getType() == objectTypeClosure
}

func IsCallable(v Object) bool {
	t := v.getType()
	return t == objectTypeFunction ||
		t == objectTypeObjectFunc ||
		t == objectTypeByteFunc ||
		t == objectTypeClosure ||
		t == objectTypeBuiltin
}

func Typeof(v Object) string {
	return toString(v.getType())
}

func ToInteger(v Object) (int64, error) {
	return v.asInteger()
}

func ToBoolean(v bool) *Boolean {
	if v {
		return True
	} else {
		return False
	}
}

type ByteEncoder interface {
	encode(op code.Opcode, operands ...int) (int, error)
}

type BuiltinFunction func(args Objects) (Object, error)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

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
	Incr()
	AsByteFunc() (*ByteFunc, error)
	AsClosure() (*Closure, error)
	AsArray() (*Array, error)

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
	equalByteFunc(other *ByteFunc) error
	equalClosure(other *Closure) error
	equalObjectFunc(other *ObjectFunc) error
	// calc
	calcInteger(op *token.Token, left *Integer) (Object, error)
	calcString(op *token.Token, left *String) (Object, error)
	calcBoolean(op *token.Token, left *Boolean) (Object, error)
	calcNull(op *token.Token, left *Null) (Object, error)
	calcArray(op *token.Token, left *Array) (Object, error)
	calcHash(op *token.Token, left *Hash) (Object, error)
	calcBuiltin(op *token.Token, left *Builtin) (Object, error)
	calcFunction(op *token.Token, left *Function) (Object, error)
	calcByteFunc(op *token.Token, left *ByteFunc) (Object, error)
	calcClosure(op *token.Token, left *Closure) (Object, error)
	calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error)
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

func hash64(b []byte) uint64 {
	h := fnv.New64a()
	h.Write(b)
	return h.Sum64()
}

func unsupported(entry string, obj Object) error {
	return fmt.Errorf("%v -> unsupported, (`%v`)", entry, obj.String())
}

func unsupportedOp(entry string, op *token.Token, obj Object) error {
	return fmt.Errorf("%v -> unsupported op %v(%v), (`%v`)", entry, op.Literal, token.ToString(op.Type), obj.String())
}

func toString(t ObjectType) string {
	s, ok := objectTypeStrings[t]
	if ok {
		return s
	}
	return "undefined type"
}

func toInt64(v bool) int64 {
	if v {
		return 1
	} else {
		return 0
	}
}

func toInteger(v bool) *Integer {
	return newInteger(toInt64(v))
}

func infixNull(op *token.Token, right Object, method string) (Object, error) {
	switch op.Type {
	case token.LT:
		return ToBoolean(true), nil
	case token.LEQ:
		return ToBoolean(true), nil
	case token.GT:
		return ToBoolean(false), nil
	case token.GEQ:
		return ToBoolean(false), nil
	case token.EQ:
		return ToBoolean(false), nil
	case token.NEQ:
		return ToBoolean(true), nil
	case token.AND:
		return Nil, nil
	case token.OR:
		return right, nil
	default:
		err := fmt.Errorf("(%v) unsupported op %v(%v)", method, op.Literal, token.ToString(op.Type))
		return Nil, err
	}
}

func checkIdx(idx int64, sz int64) error {
	if idx < 0 {
		err := fmt.Errorf("list index out of range, idx: %v", idx)
		return err
	}
	if idx > sz-1 {
		err := fmt.Errorf("list index out of range, idx: %v, len: %v", idx, sz)
		return err
	}
	return nil
}

func indexofArray(items Objects, idx int64) (Object, error) {
	sz := int64(len(items))
	if err := checkIdx(idx, sz); nil != err {
		return Nil, err
	}
	return items[idx], nil
}

func setValue(items Objects, idx int64, v Object) (Object, error) {
	sz := int64(len(items))
	if err := checkIdx(idx, sz); nil != err {
		return Nil, err
	}
	items[idx] = v
	return NewString(""), nil
}

func indexofString(s string, idx int64) (Object, error) {
	sz := int64(len(s))
	if err := checkIdx(idx, sz); nil != err {
		return Nil, err
	}
	return NewString(s[idx : idx+1]), nil
}

func keyofHash(m HashMap, key Object) (Object, error) {
	k, err := key.Hash()
	if nil != err {
		return Nil, err
	}
	v, ok := m.get(k)
	if !ok {
		err := fmt.Errorf("key `%v` missing", key.String())
		return Nil, err
	}
	return v.Value, nil
}

func notEqual(op *token.Token) (Object, error) {
	switch op.Type {
	case token.EQ:
		return ToBoolean(false), nil
	case token.NEQ:
		return ToBoolean(true), nil
	default:
		return Nil, errInvalidOperation
	}
}

func compare(entry string, this Object, left Object, op *token.Token) (Object, error) {
	switch op.Type {
	case token.EQ:
		return ToBoolean(nil == this.equal(left)), nil
	case token.NEQ:
		return ToBoolean(nil != this.equal(left)), nil
	default:
		return Nil, unsupportedOp(function.GetFunc(), op, this)
	}
}

func callMember(this Object, fns objectBuiltins, name string, args Objects) (Object, error) {
	fn, ok := fns.get(name)
	if !ok {
		err := fmt.Errorf("no attribute '%v' in %v, (`%v`)", name, Typeof(this), this.String())
		return Nil, err
	}
	return fn(args)
}

func getMember(this Object, fns objectBuiltins, name string) (Object, error) {
	fn, ok := fns.get(name)
	if !ok {
		err := fmt.Errorf("no attribute '%v' in %v, (`%v`)", name, Typeof(this), this.String())
		return Nil, err
	}
	return NewObjectFunc(this, name, fn), nil
}

func defaultNot(this Object, args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("not() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return False, nil
}
