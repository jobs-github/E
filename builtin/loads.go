package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/json"
	"github.com/jobs-github/escript/object"
)

func builtinLoads(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("loads() takes exactly one argument (%v given)", argc)
	}
	if !object.IsString(args[0]) {
		return object.Nil, fmt.Errorf("loads the first argument should be string (%v given)", object.Typeof(args[0]))
	}
	s := unquote(args[0].String())
	v, err := json.Decode(s)
	if nil != err {
		return object.Nil, err
	}
	return v, nil
}
