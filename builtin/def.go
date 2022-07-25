package builtin

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jobs-github/escript/function"
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

type openArgs struct {
	url  string
	mode string
	flag int
}

func (this *openArgs) size() int64 {
	n, _ := function.FileSize(this.url)
	return n
}

func (this *openArgs) open() object.Object {
	f, err := os.OpenFile(this.url, this.flag, 0666)
	if nil != err {
		return object.Nil
	}
	return object.NewFile(this.url, this.mode, f, this.size())
}

func getFlag(mode string) (int, error) {
	if modeRead == mode {
		return os.O_RDONLY, nil
	} else if modeWrite == mode {
		return os.O_WRONLY, nil
	} else if modeAppend == mode {
		return os.O_WRONLY | os.O_APPEND | os.O_CREATE, nil
	} else if modeReadUpdate == mode {
		return os.O_RDWR, nil
	} else if modeWriteUpdate == mode {
		return os.O_RDWR | os.O_CREATE | os.O_TRUNC, nil
	} else if modeAppendUpdate == mode {
		return os.O_RDWR | os.O_APPEND | os.O_CREATE, nil
	} else {
		return 0, function.NewError(errMode)
	}
}

func newOpenArgs(args object.Objects) (*openArgs, error) {
	url := args[0]
	if !object.IsString(url) {
		return nil, fmt.Errorf("open() the first argument should be string (%v given)", object.Typeof(url))
	}
	if len(args) == 1 {
		return &openArgs{url.String(), modeRead, os.O_RDONLY}, nil
	}
	m := args[1]
	if !object.IsString(m) {
		return nil, fmt.Errorf("open() argument 2 must be string (%v given)", object.Typeof(m))
	}
	mode := m.String()
	flag, err := getFlag(mode)
	if nil != err {
		return nil, function.NewError(err)
	}
	return &openArgs{url.String(), mode, flag}, nil
}
