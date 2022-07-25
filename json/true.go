package json

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

const (
	trueStr = "true"
	trueLen = len(trueStr)
)

// trueDecoder : implement decoder
type trueDecoder struct{}

func (this *trueDecoder) decode(s string, depth int, maxDepth int) (object.Object, string, error) {
	if len(s) < trueLen || s[:trueLen] != trueStr {
		err := fmt.Errorf("unexpected value found: %q", s)
		return nil, s, function.NewError(err)
	}
	return object.True, s[trueLen:], nil
}
