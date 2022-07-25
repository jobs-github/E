package json

import (
	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

// objectDecoder : implement decoder
type objectDecoder struct {
	p Parser
}

func (this *objectDecoder) decode(s string, depth int, maxDepth int) (object.Object, string, error) {
	v, tail, err := this.parseObject(s[1:], depth, maxDepth)
	if err != nil {
		return nil, tail, function.NewError(err)
	}
	return v, tail, nil
}

func (this *objectDecoder) parseObject(s string, depth int, maxDepth int) (object.Object, string, error) {
	s = skipWS(s)
	if len(s) == 0 {
		return nil, s, function.NewError(errMissingObjectEnd)
	}
	m := object.HashMap{}
	if s[0] == '}' {
		return object.NewHash(m), s[1:], nil
	}
	var key object.Object
	var val object.Object
	var err error
	for {
		s = skipWS(s)
		if len(s) == 0 || s[0] != '"' {
			return nil, s, function.NewError(errObjectStart)
		}
		key, s, err = parseRawKey(s[1:])
		if nil != err {
			return nil, s, function.NewError(err)
		}
		h, err := key.Hash()
		if nil != err {
			return nil, s, function.NewError(err)
		}
		s = skipWS(s)
		if len(s) == 0 || s[0] != ':' {
			return nil, s, function.NewError(errObjectKeySeparator)
		}
		s = s[1:]
		s = skipWS(s)
		val, s, err = this.p.parseValue(s, depth, maxDepth)
		if err != nil {
			return nil, s, function.NewError(err)
		}

		m[*h] = &object.HashPair{Key: key, Value: val}

		s = skipWS(s)
		if len(s) == 0 {
			return nil, s, function.NewError(errObjectEnd)
		}
		if s[0] == ',' {
			s = s[1:]
			continue
		}
		if s[0] == '}' {
			break
		}
		return nil, s, function.NewError(errObjectValSeparator)
	}
	return object.NewHash(m), s[1:], nil
}
