package builtin

import (
	"fmt"

	"github.com/jobs-github/Q/object"
)

func builtinLen(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("len() takes exactly one argument (%v given)", argc)
	}
	return args[0].CallMember(object.FnLen, object.Objects{})
}
