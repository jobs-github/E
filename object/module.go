package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewModule(name string, as string, env Env) Object {
	return &Module{name: name, as: as, env: env}
}

// Module : implement Object
type Module struct {
	name string
	as   string
	env  Env
}

func (this *Module) String() string {
	return fmt.Sprintf("<module '%v' from '%v'>", this.as, this.name)
}

func (this *Module) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Module) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Module) Calc(op *token.Token, right Object) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *Module) CallMember(name string, args Objects) (Object, error) {
	obj, ok := this.env.Get(name)
	if !ok {
		err := fmt.Errorf("symbol `%v` missing in module (`%v`)", name, this.String())
		return Nil, function.NewError(err)
	}
	return obj.Call(args)
}

func (this *Module) GetMember(name string) (Object, error) {
	obj, ok := this.env.Get(name)
	if !ok {
		err := fmt.Errorf("symbol `%v` missing in module (`%v`)", name, this.String())
		return Nil, function.NewError(err)
	}
	return obj, nil
}

func (this *Module) True() bool {
	return false
}

func (this *Module) Return() (bool, Object) {
	return false, nil
}

func (this *Module) Break() (bool, int) {
	return false, 0
}

func (this *Module) getType() ObjectType {
	return objectTypeModule
}

func (this *Module) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *Module) equal(other Object) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalInteger(other *Integer) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalString(other *String) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalBoolean(other *Boolean) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalNull(other *Null) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalArray(other *Array) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalHash(other *Hash) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalBuiltin(other *Builtin) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalFunction(other *Function) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalObjectFunc(other *ObjectFunc) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalArrayIter(other *ArrayIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) equalHashIter(other *HashIterator) error {
	return unsupported(function.GetFunc(), this)
}

func (this *Module) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcString(op *token.Token, left *String) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcNull(op *token.Token, left *Null) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcArray(op *token.Token, left *Array) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcHash(op *token.Token, left *Hash) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}
func (this *Module) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcFunction(op *token.Token, left *Function) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}

func (this *Module) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return Nil, unsupportedOp(function.GetFunc(), op, this)
}
