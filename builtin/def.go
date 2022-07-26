package builtin

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jobs-github/escript/object"
)

var (
	errMode = errors.New(`mode string must begin with one of "r", "w", "a", "r+", "w+", "a+"`)
)

const (
	// read: Open file for input operations. The file must exist.
	modeRead = "r"
	// write: Create an empty file for output operations.
	// If a file with the same name already exists, its contents are discarded and
	// the file is treated as a new empty file.
	modeWrite = "w"
	// append: Open file for output at the end of a file. Output operations always
	// write data at the end of the file, expanding it.
	// Repositioning operations (fseek, fsetpos, rewind) are ignored.
	// The file is created if it does not exist.
	modeAppend = "a"
	// read/update: Open a file for update (both for input and output).
	// The file must exist.
	modeReadUpdate = "r+"
	// write/update: Create an empty file and open it for update (both for input and output).
	// If a file with the same name already exists its contents are discarded and
	// the file is treated as a new empty file.
	modeWriteUpdate = "w+"
	// append/update: Open a file for update (both for input and output) with
	// all output operations writing data at the end of the file.
	// Repositioning operations (fseek, fsetpos, rewind) affects the next
	// input operations, but output operations move the position back to the end
	// of file. The file is created if it does not exist.
	modeAppendUpdate = "a+"
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
