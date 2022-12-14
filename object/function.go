package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewFunction(
	name string,
	args []string,
	evalBody func(env Env) (Object, error),
	env Env,
) Object {
	obj := &Function{
		Name:     name,
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
	defaultObject
	Name     string
	Args     []string
	EvalBody func(env Env) (Object, error)
	Env      Env
}

func (this *Function) String() string {
	return fmt.Sprintf("function[%p]", this)
}

func (this *Function) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcFunction(op, this)
}

func (this *Function) Call(args Objects) (Object, error) {
	if len(args) != len(this.Args) {
		err := fmt.Errorf("%v args provided, but %v args required, (`%v`)", len(args), len(this.Args), this.String())
		return Nil, err
	}
	innerEnv := newFunctionEnv(this.Env, this.Args, args)
	evaluated, err := this.EvalBody(innerEnv)
	if nil != err {
		return Nil, err
	}
	return evaluated, nil
}

func (this *Function) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Function) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Function) getType() ObjectType {
	return objectTypeFunction
}

func (this *Function) equal(other Object) error {
	return other.equalFunction(this)
}

func (this *Function) equalFunction(other *Function) error {
	if this != other {
		return fmt.Errorf("function mismatch, this: %p, other: %p", this, other)
	}
	return nil
}

func (this *Function) calcFunction(op *token.Token, left *Function) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *Function) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
