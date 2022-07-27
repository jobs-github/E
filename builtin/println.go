package builtin

import (
	"fmt"
	"strings"

	"github.com/jobs-github/escript/object"
)

func builtinPrintln(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc == 0 {
		return object.NewString(""), nil
	}
	var s strings.Builder
	for _, arg := range args {
		s.WriteString(arg.String())
	}
	fmt.Println(s.String())
	return object.NewString(""), nil
}