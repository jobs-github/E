package builtin

import (
	"os"

	"github.com/jobs-github/Q/object"
)

func builtinRemove(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc < 1 {
		return object.Nil, nil
	}
	url := args[0]
	if !object.IsString(url) {
		return object.Nil, nil
	}
	os.Remove(url.String())
	return object.Nil, nil
}
