package builtin

import (
	"fmt"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

func builtinPrintf(args object.Objects) (object.Object, error) {
	r, err := newFormatArgs("printf()", args)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	fmt.Printf(r.format, r.args...)
	return object.NewString(""), nil
}
