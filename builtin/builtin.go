package builtin

import (
	"github.com/jobs-github/escript/object"
)

type Builtin struct {
	name string
	fn   object.Object
}

func newBuiltin(name string, fn object.BuiltinFunction) *Builtin {
	return &Builtin{
		name: name,
		fn:   object.NewBuiltin(fn, name),
	}
}

type BuiltinSlice []*Builtin

func (this *BuiltinSlice) newMap() map[string]object.Object {
	m := map[string]object.Object{}
	for _, v := range *this {
		m[v.name] = v.fn
	}
	return m
}

func (this *BuiltinSlice) traverse(cb func(i int, name string)) {
	for i, v := range *this {
		cb(i, v.name)
	}
}

var (
	builtinslice = BuiltinSlice{
		newBuiltin(object.FnLen, builtinLen),
		newBuiltin("type", builtinType),
		newBuiltin("str", builtinStr),
		newBuiltin("int", builtinInt),
		newBuiltin("print", builtinPrint),
		newBuiltin("println", builtinPrintln),
		newBuiltin("printf", builtinPrintf),
		newBuiltin("sprintf", builtinSprintf),
		newBuiltin("loads", builtinLoads),
		newBuiltin("dumps", builtinDumps),
		newBuiltin("for", builtinFor),
		newBuiltin("state", builtinState),
	}
	builtins = builtinslice.newMap()
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

func GetFn(idx int) object.Object {
	return builtinslice[idx].fn
}

func Traverse(cb func(i int, name string)) {
	builtinslice.traverse(cb)
}
