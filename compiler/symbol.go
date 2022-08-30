package compiler

import "fmt"

type SymbolScope uint

const (
	GlobalScope SymbolScope = iota
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

func NewSymbolTable() SymbolTable {
	return &symbolTable{
		m:   map[string]*Symbol{},
		cnt: 0,
	}
}

type SymbolTable interface {
	define(key string) *Symbol
	resolve(key string) (*Symbol, error)
}

// symbolTable : implement SymbolTable
type symbolTable struct {
	m   map[string]*Symbol
	cnt int
}

func (this *symbolTable) define(key string) *Symbol {
	s := newSymbol(key, GlobalScope, this.cnt)
	this.m[key] = s
	this.cnt++
	return s
}

func (this *symbolTable) resolve(key string) (*Symbol, error) {
	v, ok := this.m[key]
	if !ok {
		return nil, fmt.Errorf("symbol `%v` missing", key)
	}
	return v, nil
}
