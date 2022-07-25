package object

import (
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/token"
)

func NewReturn(v Object) Object {
	return &ReturnValue{
		Value: v,
	}
}

// ReturnValue : implement Object
type ReturnValue struct {
	Value Object
}

func (this *ReturnValue) String() string {
	return this.Value.String()
}

func (this *ReturnValue) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) Calc(op *token.Token, right Object) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) CallMember(name string, args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) GetMember(name string) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) True() bool {
	return false
}

func (this *ReturnValue) Return() (bool, Object) {
	return true, this.Value
}

func (this *ReturnValue) Break() (bool, int) {
	return false, 0
}

func (this *ReturnValue) getType() ObjectType {
	return objectTypeReturnValue
}

func (this *ReturnValue) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equal(other Object) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalInteger(other *Integer) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalString(other *String) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalBoolean(other *Boolean) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalNull(other *Null) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalArray(other *Array) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalHash(other *Hash) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalBuiltin(other *Builtin) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalFunction(other *Function) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalObjectFunc(other *ObjectFunc) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalArrayIter(other *ArrayIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) equalHashIter(other *HashIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *ReturnValue) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcString(op *token.Token, left *String) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcNull(op *token.Token, left *Null) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcArray(op *token.Token, left *Array) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcHash(op *token.Token, left *Hash) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcFunction(op *token.Token, left *Function) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *ReturnValue) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}
