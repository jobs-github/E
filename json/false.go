package json

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

const (
	falseStr = "false"
	falseLen = len(falseStr)
)

// falseDecoder : implement decoder
type falseDecoder struct{}

func (this *falseDecoder) decode(s string, depth int, maxDepth int) (object.Object, string, error) {
	if len(s) < falseLen || s[:falseLen] != falseStr {
		err := fmt.Errorf("unexpected value found: %q", s)
		return nil, s, function.NewError(err)
	}
	return object.False, s[falseLen:], nil
}
