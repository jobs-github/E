package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/object"
)

func builtinInt(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("int() takes exactly one argument (%v given)", argc)
	}

	return args[0].CallMember(object.FnInt, object.Objects{})
}
