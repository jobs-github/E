package object

import (
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

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
	objectTypeState
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
	FnLen    = "len"
	FnIndex  = "index"
	FnNot    = "not"
	FnNeg    = "neg"
	FnInt    = "int"
	FnMap    = "map"
	FnReduce = "reduce"
	FnFilter = "filter"
	FnFirst  = "first"
	FnLast   = "last"
	FnTail   = "tail"
	FnPush   = "push"
	FnKeys   = "keys"
	FnValue  = "value"
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
		objectTypeState:      "state",
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

func Callable(v Object) bool {
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

type BuiltinFunction func(args Objects) (Object, error)

type HashKey struct {
	Type  ObjectType
	Value uint64
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

func ToBoolean(v bool) *Boolean {
	if v {
		return True
	} else {
		return False
	}
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
