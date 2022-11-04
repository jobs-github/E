package compiler

import "fmt"

type SymbolScope uint

const (
	ScopeGlobal SymbolScope = iota
	ScopeLocal
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
	return &symbolTable{
		parent:      parent,
		m:           map[string]*Symbol{},
		sz:          0,
		freeSymbols: Symbols{},
	}
}

type SymbolTable interface {
	newEnclosed() SymbolTable
	size() int
	outer() SymbolTable
	define(key string) *Symbol
	resolve(key string) (*Symbol, error)
	free() Symbols
	defineFree(s *Symbol) *Symbol
}

// symbolTable : implement SymbolTable
type symbolTable struct {
	parent      SymbolTable
	m           map[string]*Symbol
	sz          int
	freeSymbols Symbols
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

func (this *symbolTable) resolve(key string) (*Symbol, error) {
	v, ok := this.m[key]
	if !ok {
		if this.parent != nil {
			return this.parent.resolve(key)
		} else {
			return nil, fmt.Errorf("symbol `%v` missing", key)
		}
	}
	return v, nil
}

func (this *symbolTable) free() Symbols {
	return this.freeSymbols
}

func (this *symbolTable) defineFree(s *Symbol) *Symbol {
	idx := len(this.freeSymbols)
	this.freeSymbols = append(this.freeSymbols, s)
	symbol := newSymbol(s.Name, ScopeFree, idx)
	this.m[s.Name] = symbol
	return symbol
}
