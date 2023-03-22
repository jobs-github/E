package builtin

import (
	"github.com/jobs-github/escript/object"
)

type symbol struct {
	name string
	fn   object.Object
}

func newSymbol(name string, fn object.BuiltinFunction) *symbol {
	return &symbol{
		name: name,
		fn:   object.NewBuiltin(fn, name),
	}
}

type symbolTable []*symbol

func (this *symbolTable) newMap() map[string]object.Object {
	m := map[string]object.Object{}
	for _, v := range *this {
		m[v.name] = v.fn
	}
	return m
}

func (this *symbolTable) traverse(cb func(i int, name string)) {
	for i, v := range *this {
		cb(i, v.name)
	}
}

var (
	builtinSymbolTable = symbolTable{
		newSymbol("type", builtinType),
		newSymbol("str", builtinStr),
		newSymbol("print", builtinPrint),
		newSymbol("println", builtinPrintln),
		newSymbol("printf", builtinPrintf),
		newSymbol("sprintf", builtinSprintf),
		newSymbol("loads", builtinLoads),
		newSymbol("dumps", builtinDumps),
	}
	builtins = builtinSymbolTable.newMap()
)

func IsBuiltin(key string) bool {
	_, ok := builtins[key]
	return ok
}

func Get(key string) object.Object {
	if fn, ok := builtins[key]; ok {
		return fn
	} else {
		return nil
	}
}

func Resolve(idx int) object.Object {
	return builtinSymbolTable[idx].fn
}

func Traverse(cb func(i int, name string)) {
	builtinSymbolTable.traverse(cb)
}
