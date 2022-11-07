package builtin

import (
	"github.com/jobs-github/escript/object"
)

var builtins = map[string]object.Object{
	object.FnLen: object.NewBuiltin(builtinLen, object.FnLen),
	"type":       object.NewBuiltin(builtinType, "type"),
	"str":        object.NewBuiltin(builtinStr, "str"),
	"int":        object.NewBuiltin(builtinInt, "int"),
	"print":      object.NewBuiltin(builtinPrint, "print"),
	"println":    object.NewBuiltin(builtinPrintln, "println"),
	"printf":     object.NewBuiltin(builtinPrintf, "printf"),
	"sprintf":    object.NewBuiltin(builtinSprintf, "sprintf"),
	"loads":      object.NewBuiltin(builtinLoads, "loads"),
	"dumps":      object.NewBuiltin(builtinDumps, "dumps"),
	"for":        object.NewBuiltin(builtinFor, "for"),
	"state":      object.NewBuiltin(builtinState, "state"),
}

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
