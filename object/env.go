package object

type Env interface {
	Symbol(name string) (Callable, bool)
	Get(name string) (Object, bool)
	Set(name string, val Object) Object
	NewEnclosedEnv() Env
}

// environment : implement Env
type environment struct {
	s      Symbols
	parent Env
	e      SymbolTable
}

func NewEnv(s Symbols) Env {
	return &environment{
		s:      s,
		parent: nil,
		e:      SymbolTable{},
	}
}

func (this *environment) Symbol(name string) (Callable, bool) {
	v, ok := this.s[name]
	return v, ok
}

func (this *environment) Get(name string) (Object, bool) {
	v, ok := this.e[name]
	if !ok && nil != this.parent {
		v, ok = this.parent.Get(name)
		return v, ok
	}
	return v, ok
}

func (this *environment) Set(name string, val Object) Object {
	this.e[name] = val
	return val
}

func (this *environment) NewEnclosedEnv() Env {
	return &environment{
		s:      this.s,
		parent: this,
		e:      SymbolTable{},
	}
}

func newFunctionEnv(outer Env, args []string, values Objects) Env {
	env := outer.NewEnclosedEnv()
	for i, name := range args {
		env.Set(name, values[i])
	}
	return env
}
