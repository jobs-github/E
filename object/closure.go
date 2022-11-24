package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewClosure(fn *ByteFunc, frees Objects) *Closure {
	obj := &Closure{Fn: fn, Free: frees}
	obj.fns = objectBuiltins{
		FnNot: obj.builtinNot,
	}
	return obj
}

// In summary
// Detect references to free variables while compiling a function,
// get the referenced values on to the stack,
// merge the values and the compiled function into a closure and
// leave it on the stack where it can then be called
//
// The environment in tree-walking interpreter is scattered among the globals store and different regions of the stack,
// all of which can be wiped out with a return from a function.
//
// Give compiled functions the ability to hold bindings that are only created
// at run time and their instructions must already reference said bindings.
//
// The compiler needs to detect references to free variables and emit instructions that will load them on to the stack.
// The VM must not only resolve references to free variables correctly, but also store them on compiled functions.
//
// While compiling the function’s body, inspect each resolved symbol to find out whether it’s a reference to a free variable.
// After the function’s body is compiled and left its scope, SymbolTable should tell us how many free variables
// were referenced and in which scope they were originally defined.
//
// At run time, transfer these free variables to the compiled function.

// Closure : implement Object
type Closure struct {
	defaultObject
	Fn   *ByteFunc
	Free Objects // env
}

func (this *Closure) String() string {
	return fmt.Sprintf("closure[%p]", this)
}

func (this *Closure) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcClosure(op, this)
}

func (this *Closure) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Closure) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Closure) True() bool {
	return false
}

func (this *Closure) AsClosure() (*Closure, error) {
	return this, nil
}

func (this *Closure) getType() ObjectType {
	return objectTypeClosure
}

func (this *Closure) equal(other Object) error {
	return other.equalClosure(this)
}

func (this *Closure) equalClosure(other *Closure) error {
	// TODO
	return nil
}

func (this *Closure) calcClosure(op *token.Token, left *Closure) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *Closure) builtinNot(args Objects) (Object, error) {
	return defaultNot(this, args)
}
