package json

import (
	"fmt"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

const (
	nullStr = "null"
	nullLen = len(nullStr)
)

// nullDecoder : implement decoder
type nullDecoder struct{}

func (this *nullDecoder) decode(s string, depth int, maxDepth int) (object.Object, string, error) {
	if len(s) < nullLen || s[:nullLen] != nullStr {
		err := fmt.Errorf("unexpected value found: %q", s)
		return nil, s, function.NewError(err)
	}
	return object.Nil, s[nullLen:], nil
}
