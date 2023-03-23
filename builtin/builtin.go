package builtin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jobs-github/escript/function"
	ejson "github.com/jobs-github/escript/json"
	"github.com/jobs-github/escript/object"
)

// public
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

type formatArgs struct {
	format string
	args   []interface{}
}

func unquote(input string) string {
	s, err := strconv.Unquote("\"" + input + "\"")
	if nil == err {
		return s
	}
	return input
}

func newFormatArgs(entry string, args object.Objects) (*formatArgs, error) {
	argc := len(args)
	if argc < 2 {
		return nil, fmt.Errorf("%v takes at least 2 arguments (%v given)", entry, argc)
	}
	s := []interface{}{}
	for i := 1; i < argc; i++ {
		s = append(s, args[i].String())
	}
	format := args[0]
	if !object.IsString(format) {
		return nil, fmt.Errorf("%v the first argument should be string (%v given)", entry, object.Typeof(format))
	}
	return &formatArgs{format: unquote(format.String()), args: s}, nil
}

// implement
func builtinType(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("type() takes exactly one argument (%v given)", argc)
	}
	return object.NewString(object.Typeof(args[0])), nil
}

func builtinStr(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("str() takes exactly one argument (%v given)", argc)
	}
	return object.NewString(args[0].String()), nil
}

func builtinPrint(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc == 0 {
		return object.NewString(""), nil
	}
	var s strings.Builder
	for _, arg := range args {
		s.WriteString(arg.String())
	}
	fmt.Print(s.String())
	return object.NewString(""), nil
}

func builtinPrintln(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc == 0 {
		return object.NewString(""), nil
	}
	var s strings.Builder
	for _, arg := range args {
		s.WriteString(arg.String())
	}
	fmt.Println(s.String())
	return object.NewString(""), nil
}

func builtinPrintf(args object.Objects) (object.Object, error) {
	r, err := newFormatArgs("printf()", args)
	if nil != err {
		return object.Nil, err
	}
	fmt.Printf(r.format, r.args...)
	return object.NewString(""), nil
}

func builtinSprintf(args object.Objects) (object.Object, error) {
	r, err := newFormatArgs("printf()", args)
	if nil != err {
		return object.Nil, err
	}
	return object.NewString(fmt.Sprintf(r.format, r.args...)), nil
}

func builtinLoads(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("loads() takes exactly one argument (%v given)", argc)
	}
	if !object.IsString(args[0]) {
		return object.Nil, fmt.Errorf("loads the first argument should be string (%v given)", object.Typeof(args[0]))
	}
	s := unquote(args[0].String())
	v, err := ejson.Decode(s)
	if nil != err {
		return object.Nil, err
	}
	return v, nil
}

func builtinDumps(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("dumps() takes exactly one argument (%v given)", argc)
	}
	v, err := args[0].Dump()
	if nil != err {
		return object.Nil, err
	}
	b, err := json.Marshal(v)
	if nil != err {
		return object.Nil, err
	}
	return object.NewString(function.BytesToString(b)), nil
}
