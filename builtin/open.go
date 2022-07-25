package builtin

import (
	"github.com/jobs-github/escript/object"
)

func builtinOpen(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc < 1 {
		return object.Nil, nil
	}
	oa, err := newOpenArgs(args)
	if nil != err {
		return object.Nil, nil
	}
	return oa.open(), nil
}
