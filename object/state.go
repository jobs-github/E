package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewState(quit bool, v Object) *State {
	obj := &State{
		Quit:  quit,
		Value: v,
	}
	obj.fns = objectBuiltins{
		"value": obj.builtinValue,
	}
	return obj
}

// State : implement Object
type State struct {
	Quit  bool
	Value Object
	fns   objectBuiltins
}

func (this *State) String() string {
	return this.Value.String()
}

func (this *State) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *State) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *State) Calc(op *token.Token, right Object) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *State) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *State) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *State) True() bool {
	return false
}

func (this *State) AsState() (*State, error) {
	return this, nil
}

func (this *State) getType() ObjectType {
	return objectTypeState
}

func (this *State) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *State) equal(other Object) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalInteger(other *Integer) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalString(other *String) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalBoolean(other *Boolean) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalNull(other *Null) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalArray(other *Array) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalHash(other *Hash) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalBuiltin(other *Builtin) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalFunction(other *Function) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) equalObjectFunc(other *ObjectFunc) error {
	return unsupported(function.GetFunc(), this)
}

func (this *State) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcString(op *token.Token, left *String) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcNull(op *token.Token, left *Null) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcArray(op *token.Token, left *Array) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcHash(op *token.Token, left *Hash) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcFunction(op *token.Token, left *Function) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *State) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

// builtin
func (this *State) builtinValue(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("value() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return this.Value, nil
}
