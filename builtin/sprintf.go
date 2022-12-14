package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/object"
)

func builtinSprintf(args object.Objects) (object.Object, error) {
	r, err := newFormatArgs("printf()", args)
	if nil != err {
		return object.Nil, err
	}
	return object.NewString(fmt.Sprintf(r.format, r.args...)), nil
}
