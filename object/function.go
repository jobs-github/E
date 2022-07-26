package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewFunction(
	name string,
	fn function.Function,
	args []string,
	evalBody func(env Env, insideLoop bool) (Object, error),
	env Env,
) Object {
	obj := &Function{
		Name:     name,
		Fn:       fn,
		Args:     args,
		EvalBody: evalBody,
		Env:      env,
	}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

// Function : implement Object
type Function struct {
	Name     string
	Fn       function.Function
	Args     []string
	EvalBody func(env Env, insideLoop bool) (Object, error)
	Env      Env
	fns      objectBuiltins
}

func (this *Function) String() string {
	return this.Fn.String()
}

func (this *Function) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Function) Dump() (interface{}, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Function) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcFunction(op, this)
}

func (this *Function) Call(args Objects) (Object, error) {
	if len(args) != len(this.Args) {
		err := fmt.Errorf("%v args provided, but %v args required, (`%v`)", len(args), len(this.Args), this.String())
		return Nil, function.NewError(err)
	}
	innerEnv := newFunctionEnv(this.Env, this.Args, args)
	evaluated, err := this.EvalBody(innerEnv, false)
	if nil != err {
		return Nil, function.NewError(err)
	}
	return evaluated, nil
}

func (this *Function) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Function) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Function) True() bool {
	return false
}

func (this *Function) getType() ObjectType {
	return objectTypeFunction
}

func (this *Function) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *Function) equal(other Object) error {
	return other.equalFunction(this)
}

func (this *Function) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalFunction(other *Function) error {
	if this.Name != "" {
		if this.Name != other.Name {
			return fmt.Errorf("function mismatch, this: %v, other: %v", this.Name, other.Name)
		}
	}
	src := this.String()
	dst := other.String()
	if src != dst {
		return fmt.Errorf("function mismatch, this: %v, other: %v", src, dst)
	}
	return nil
}

func (this *Function) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalArrayIter(other *ArrayIterator) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) equalHashIter(other *HashIterator) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Function) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcFunction(op *token.Token, left *Function) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *Function) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcArrayIter(op *token.Token, left *ArrayIterator) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcHashIter(op *token.Token, left *HashIterator) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Function) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *Function) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
