package json

import (
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// stringDecoder : implement decoder
type stringDecoder struct{}

func (this *stringDecoder) decode(s string, depth int, maxDepth int) (object.Object, string, error) {
	v, tail, err := parseRawString(s[1:])
	if err != nil {
		return nil, tail, function.NewError(err)
	}
	return v, tail, nil
}
