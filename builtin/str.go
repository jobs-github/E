package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/object"
)

func builtinStr(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("str() takes exactly one argument (%v given)", argc)
	}
	return object.NewString(args[0].String()), nil
}
