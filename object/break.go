package object

import (
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewBreak() Object {
	return &BreakObject{count: 0}
}

// BreakObject : implement Object
type BreakObject struct {
	count int
}

func (this *BreakObject) String() string {
	return toString(objectTypeBreakObject)
}

func (this *BreakObject) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *BreakObject) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *BreakObject) Calc(op *token.Token, right Object) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *BreakObject) CallMember(name string, args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *BreakObject) GetMember(name string) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *BreakObject) True() bool {
	return false
}

func (this *BreakObject) Return() (bool, Object) {
	return false, nil
}

func (this *BreakObject) Break() (bool, int) {
	this.count++
	return true, this.count
}

func (this *BreakObject) getType() ObjectType {
	return objectTypeBreakObject
}

func (this *BreakObject) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equal(other Object) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalInteger(other *Integer) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalString(other *String) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalBoolean(other *Boolean) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalNull(other *Null) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalArray(other *Array) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalHash(other *Hash) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalBuiltin(other *Builtin) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalFunction(other *Function) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalObjectFunc(other *ObjectFunc) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalArrayIter(other *ArrayIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) equalHashIter(other *HashIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *BreakObject) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcString(op *token.Token, left *String) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcNull(op *token.Token, left *Null) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcArray(op *token.Token, left *Array) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcHash(op *token.Token, left *Hash) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}
func (this *BreakObject) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcFunction(op *token.Token, left *Function) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *BreakObject) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}
