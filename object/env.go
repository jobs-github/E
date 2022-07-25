package object

import (
	"fmt"

	"github.com/jobs-github/escript/function"
)

type Env interface {
	Import(module string, asKey string) error
	Symbol(asKey string) (*Symbol, bool)
	Get(name string) (Object, bool)
	Set(name string, val Object) Object
	Assign(name string, val Object) error
	NewEnclosedEnv() Env
}

type Importer interface {
	Load(module string) (Env, error)
}

type Symbol struct {
	ModuleName string
	AsKey      string
	E          Env
}

// environment : implement Env
type environment struct {
	importer Importer
	parent   Env
	m        *modules
	e        map[string]Object
}

func NewEnv(importer Importer) Env {
	return &environment{
		importer: importer,
		parent:   nil,
		m:        newModules(),
		e:        map[string]Object{},
	}
}

func (this *environment) Import(module string, asKey string) error {
	e, err := this.importer.Load(module)
	if nil != err {
		return function.NewError(err)
	}
	return this.m.add(module, asKey, e)
}

func (this *environment) Symbol(asKey string) (*Symbol, bool) {
	e, ok := this.m.s.get(asKey)
	if !ok && nil != this.parent {
		e, ok = this.parent.Symbol(asKey)
		return e, ok
	}
	return e, ok
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
		importer: this.importer,
		parent:   this,
		m:        newModules(),
		e:        map[string]Object{},
	}
}

func newFunctionEnv(outer Env, args []string, values Objects) Env {
	env := outer.NewEnclosedEnv()
	for i, name := range args {
		env.Set(name, values[i])
	}
	return env
}

func newSymbol(name string, asKey string, e Env) *Symbol {
	return &Symbol{ModuleName: name, AsKey: asKey, E: e}
}

type symbols map[string]*Symbol

func (this *symbols) get(k string) (*Symbol, bool) {
	v, ok := (*this)[k]
	return v, ok
}

func (this *symbols) set(k string, v *Symbol) {
	(*this)[k] = v
}

type modules struct {
	s symbols
	m function.StringSet
}

func newModules() *modules {
	return &modules{
		s: symbols{},
		m: function.StringSet{},
	}
}

func (this *modules) add(module string, asKey string, e Env) error {
	if _, ok := this.s.get(asKey); ok {
		return fmt.Errorf("duplicated symbol %v", asKey)
	}
	if this.m.Find(module) {
		return fmt.Errorf("duplicated module %v", module)
	}
	this.s.set(asKey, newSymbol(module, asKey, e))
	this.m.Add(module)
	return nil
}
