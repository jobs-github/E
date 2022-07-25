package json

import (
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// arrayDecoder : implement decoder
type arrayDecoder struct {
	p Parser
}

func (this *arrayDecoder) decode(s string, depth int, maxDepth int) (object.Object, string, error) {
	v, tail, err := this.parseArray(s[1:], depth, maxDepth)
	if err != nil {
		return nil, tail, function.NewError(err)
	}
	return v, tail, nil
}

func (this *arrayDecoder) parseArray(s string, depth int, maxDepth int) (object.Object, string, error) {
	s = skipWS(s)
	if len(s) == 0 {
		return nil, s, function.NewError(errMissingArrayEnd)
	}
	items := object.Objects{}
	if s[0] == ']' {
		return object.NewArray(items), s[1:], nil
	}
	var v object.Object
	var err error
	for {
		s = skipWS(s)
		v, s, err = this.p.parseValue(s, depth, maxDepth)
		if err != nil {
			return nil, s, function.NewError(err)
		}
		items = append(items, v)
		s = skipWS(s)
		if len(s) == 0 {
			return nil, s, function.NewError(errArrayEnd)
		}
		if s[0] == ',' {
			s = s[1:]
			continue
		}
		if s[0] == ']' {
			break
		}
		return nil, s, function.NewError(errArrayValSeparator)
	}
	return object.NewArray(items), s[1:], nil
}
