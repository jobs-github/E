package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/object"
)

func builtinPrintf(args object.Objects) (object.Object, error) {
	r, err := newFormatArgs("printf()", args)
	if nil != err {
		return object.Nil, err
	}
	fmt.Printf(r.format, r.args...)
	return object.NewString(""), nil
}
