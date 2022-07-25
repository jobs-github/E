package builtin

import (
	"fmt"
	"os"

	"github.com/jobs-github/Q/object"
)

func builtinGetenv(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 1 {
		return object.Nil, fmt.Errorf("getenv() takes exactly one argument (%v given)", argc)
	}
	if !object.IsString(args[0]) {
		return object.Nil, fmt.Errorf("getenv the first argument should be string (%v given)", object.Typeof(args[0]))
	}
	return object.NewString(os.Getenv(args[0].String())), nil
}
