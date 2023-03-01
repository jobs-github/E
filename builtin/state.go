package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/object"
)

func builtinState(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 && argc != 2 {
		return object.Nil, fmt.Errorf("state() takes exactly 1 or 2 arguments (%v given)", argc)
	}
	if argc == 1 {
		return object.NewState(args[0], false), nil
	} else {
		return object.NewState(args[0], args[1].True()), nil
	}
}
