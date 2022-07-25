package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
)

type Env interface {
	Get(name string) (Object, bool)
	Set(name string, val Object) Object
	Assign(name string, val Object) error
	NewEnclosedEnv() Env
}

// environment : implement Env
type environment struct {
	parent Env
	e      map[string]Object
}

func NewEnv() Env {
	return &environment{
		parent: nil,
		e:      map[string]Object{},
	}
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

func (this *environment) Assign(name string, val Object) error {
	if _, ok := this.e[name]; ok {
		this.e[name] = val
	} else {
		if nil == this.parent {
			err := fmt.Errorf("`%v` undefined", name)
			return function.NewError(err)
		} else {
			return this.parent.Assign(name, val)
		}
	}
	return nil
}

func (this *environment) NewEnclosedEnv() Env {
	return &environment{
		parent: this,
		e:      map[string]Object{},
	}
}

func newFunctionEnv(outer Env, args []string, values Objects) Env {
	env := outer.NewEnclosedEnv()
	for i, name := range args {
		env.Set(name, values[i])
	}
	return env
}
