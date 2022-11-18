package object

import (
	"errors"

	"github.com/jobs-github/escript/token"
)

var (
	errNotSupportHash = errors.New("not support hash func")
	errNotSupportDump = errors.New("not support dump func")

	errNotSupportCall            = errors.New("not support call func")
	errNotSupportEqual           = errors.New("not support equal func")
	errNotSupportEqualInteger    = errors.New("not support equalInteger func")
	errNotSupportEqualString     = errors.New("not support equalString func")
	errNotSupportEqualBoolean    = errors.New("not support equalBoolean func")
	errNotSupportEqualNull       = errors.New("not support equalNull func")
	errNotSupportEqualArray      = errors.New("not support equalArray func")
	errNotSupportEqualHash       = errors.New("not support equalHash func")
	errNotSupportEqualBuiltin    = errors.New("not support equalBuiltin func")
	errNotSupportEqualFunction   = errors.New("not support equalFunction func")
	errNotSupportEqualByteFunc   = errors.New("not support equalByteFunc func")
	errNotSupportEqualClosure    = errors.New("not support equalClosure func")
	errNotSupportEqualObjectFunc = errors.New("not support equalObjectFunc func")

	errInvalidOperation = errors.New("invalid operation")
	errNotSupportCalc   = errors.New("not support calc func")

	errTypeIsNotState    = errors.New("type is not State")
	errTypeIsNotByteFunc = errors.New("type is not ByteFunc")
	errTypeIsNotClosure  = errors.New("type is not Closure")
	errTypeIsNotInt      = errors.New("type is not Integer")
)

type defaultObject struct {
	fns objectBuiltins
}

func (this *defaultObject) Hash() (*HashKey, error) {
	return nil, errNotSupportHash
}

func (this *defaultObject) Dump() (interface{}, error) {
	return nil, errNotSupportDump
}

func (this *defaultObject) Calc(op *token.Token, right Object) (Object, error) {
	return nil, errNotSupportCalc
}

func (this *defaultObject) Call(args Objects) (Object, error) {
	return Nil, errNotSupportCall
}

func (this *defaultObject) True() bool {
	return false
}

func (this *defaultObject) AsState() (*State, error) {
	return nil, errTypeIsNotState
}

func (this *defaultObject) AsByteFunc() (*ByteFunc, error) {
	return nil, errTypeIsNotByteFunc
}

func (this *defaultObject) AsClosure() (*Closure, error) {
	return nil, errTypeIsNotClosure
}

func (this *defaultObject) asInteger() (int64, error) {
	return 0, errTypeIsNotInt
}

func (this *defaultObject) equal(other Object) error {
	return errNotSupportEqual
}

func (this *defaultObject) equalInteger(other *Integer) error {
	return errNotSupportEqualInteger
}

func (this *defaultObject) equalString(other *String) error {
	return errNotSupportEqualString
}

func (this *defaultObject) equalBoolean(other *Boolean) error {
	return errNotSupportEqualBoolean
}

func (this *defaultObject) equalNull(other *Null) error {
	return errNotSupportEqualNull
}

func (this *defaultObject) equalArray(other *Array) error {
	return errNotSupportEqualArray
}

func (this *defaultObject) equalHash(other *Hash) error {
	return errNotSupportEqualHash
}

func (this *defaultObject) equalBuiltin(other *Builtin) error {
	return errNotSupportEqualBuiltin
}

func (this *defaultObject) equalFunction(other *Function) error {
	return errNotSupportEqualFunction
}

func (this *defaultObject) equalByteFunc(other *ByteFunc) error {
	return errNotSupportEqualByteFunc
}

func (this *defaultObject) equalClosure(other *Closure) error {
	return errNotSupportEqualClosure
}

func (this *defaultObject) equalObjectFunc(other *ObjectFunc) error {
	return errNotSupportEqualObjectFunc
}

func (this *defaultObject) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcByteFunc(op *token.Token, left *ByteFunc) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcClosure(op *token.Token, left *Closure) (Object, error) {
	return notEqual(op)
}

func (this *defaultObject) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(op)
}
