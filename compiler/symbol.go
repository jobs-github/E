package compiler

import (
	"fmt"

	"github.com/jobs-github/escript/builtin"
	"github.com/jobs-github/escript/object"
)

type SymbolScope uint

const (
	ScopeGlobal SymbolScope = iota
	ScopeLocal
	ScopeBuiltin
	ScopeObjectFn
	ScopeFree
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

func (this *Symbol) equal(other *Symbol) error {
	if this.Name != other.Name {
		return fmt.Errorf("name mismatch, this: %v, other: %v", this.Name, other.Name)
	}
	if this.Scope != other.Scope {
		return fmt.Errorf("scope mismatch, this: %v, other: %v", this.Scope, other.Scope)
	}
	if this.Index != other.Index {
		return fmt.Errorf("index mismatch, this: %v, other: %v", this.Index, other.Index)
	}
	return nil
}

func newSymbol(name string, scope SymbolScope, index int) *Symbol {
	return &Symbol{
		Name:  name,
		Scope: scope,
		Index: index,
	}
}

type Symbols []*Symbol

func NewSymbolTable(parent SymbolTable) SymbolTable {
	s := &symbolTable{
		parent: parent,
		m:      map[string]*Symbol{},
		sz:     0,
		frees:  Symbols{},
	}
	builtin.Traverse(func(i int, name string) {
		s.defineBuiltin(i, name)
	})
	object.Traverse(func(i int, name string) {
		s.defineObjectFn(i, name)
	})
	return s
}

type SymbolTable interface {
	newEnclosed() SymbolTable
	size() int
	outer() SymbolTable
	define(key string) *Symbol
	defineBuiltin(index int, name string) *Symbol
	defineObjectFn(index int, name string) *Symbol
	resolve(key string) (*Symbol, error)
	freeSymbols() Symbols
	defineFree(orginal *Symbol) *Symbol
}

// symbolTable : implement SymbolTable
type symbolTable struct {
	parent SymbolTable
	m      map[string]*Symbol
	sz     int
	frees  Symbols
}

func (this *symbolTable) newEnclosed() SymbolTable {
	return NewSymbolTable(this)
}

func (this *symbolTable) size() int {
	return this.sz
}

func (this *symbolTable) outer() SymbolTable {
	return this.parent
}

func (this *symbolTable) define(key string) *Symbol {
	s := newSymbol(key, ScopeGlobal, this.sz)
	if nil == this.parent {
		s.Scope = ScopeGlobal
	} else {
		s.Scope = ScopeLocal
	}
	this.m[key] = s
	this.sz++
	return s
}

func (this *symbolTable) defineBuiltin(index int, name string) *Symbol {
	s := newSymbol(name, ScopeBuiltin, index)
	this.m[name] = s
	return s
}

func (this *symbolTable) defineObjectFn(index int, name string) *Symbol {
	s := newSymbol(name, ScopeObjectFn, index)
	this.m[name] = s
	return s
}

func (this *symbolTable) resolve(key string) (*Symbol, error) {
	if v, ok := this.m[key]; ok {
		return v, nil
	}
	// not resolved, search parent
	if this.parent == nil {
		return nil, fmt.Errorf("symbol `%v` missing", key)
	}
	pv, err := this.parent.resolve(key)
	if nil != err {
		return nil, err
	}
	// pv is LOCAL BINDING in parent, but free binding here
	if pv.Scope == ScopeGlobal || pv.Scope == ScopeBuiltin {
		return pv, nil
	}

	// not resolved in self scope, but resolved in parent scope
	// not a global binding, or a built-in function
	// it’s a free variable
	return this.defineFree(pv), nil
}

func (this *symbolTable) freeSymbols() Symbols {
	return this.frees
}

func (this *symbolTable) defineFree(orginal *Symbol) *Symbol {
	this.frees = append(this.frees, orginal)
	symbol := newSymbol(orginal.Name, ScopeFree, len(this.frees)-1)
	this.m[orginal.Name] = symbol
	return symbol
}
