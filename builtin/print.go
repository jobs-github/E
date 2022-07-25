package builtin

import (
	"fmt"
	"strings"

	"github.com/jobs-github/Q/object"
)

func builtinPrint(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc == 0 {
		return object.NewString(""), nil
	}
	var s strings.Builder
	for _, arg := range args {
		s.WriteString(arg.String())
	}
	fmt.Print(s.String())
	return object.NewString(""), nil
}
