package builtin

import (
	"encoding/json"
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

func builtinDumps(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("dumps() takes exactly one argument (%v given)", argc)
	}
	v, err := args[0].Dump()
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	b, err := json.Marshal(v)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return object.NewString(function.BytesToString(b)), nil
}
