package builtin

import (
	"fmt"
	"strconv"

	"github.com/jobs-github/escript/object"
)

type formatArgs struct {
	format string
	args   []interface{}
}

func unquote(input string) string {
	s, err := strconv.Unquote("\"" + input + "\"")
	if nil == err {
		return s
	}
	return input
}

func newFormatArgs(entry string, args object.Objects) (*formatArgs, error) {
	argc := len(args)
	if argc < 2 {
		return nil, fmt.Errorf("%v takes at least 2 arguments (%v given)", entry, argc)
	}
	s := []interface{}{}
	for i := 1; i < argc; i++ {
		s = append(s, args[i].String())
	}
	format := args[0]
	if !object.IsString(format) {
		return nil, fmt.Errorf("%v the first argument should be string (%v given)", entry, object.Typeof(format))
	}
	return &formatArgs{format: unquote(format.String()), args: s}, nil
}
