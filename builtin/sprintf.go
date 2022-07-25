package builtin

import (
	"fmt"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

func builtinSprintf(args object.Objects) (object.Object, error) {
	r, err := newFormatArgs("printf()", args)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return object.NewString(fmt.Sprintf(r.format, r.args...)), nil
}
