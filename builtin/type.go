package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/object"
)

func builtinType(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("type() takes exactly one argument (%v given)", argc)
	}
	return object.NewString(object.Typeof(args[0])), nil
}
