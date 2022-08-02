package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/object"
)

func builtinState(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 2 {
		return object.Nil, fmt.Errorf("state() takes exactly 2 arguments (%v given)", argc)
	}
	return object.NewState(args[0].True(), args[1]), nil
}
