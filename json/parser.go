package json

import (
	"fmt"

	"github.com/jobs-github/Q/function"
	"github.com/jobs-github/Q/object"
)

func Decode(s string) (object.Object, error) {
	p := newParser()
	v, err := p.Parse(s, MaxDepth)
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	return v, nil
}

type Parser interface {
	Parse(s string, maxDepth int) (object.Object, error)
	parseValue(s string, depth int, maxDepth int) (object.Object, string, error)
}

func newParser() Parser {
	p := &parser{}
	p.decoders = map[byte]decoder{
		'{': &objectDecoder{p: p},
		'[': &arrayDecoder{p: p},
		'"': &stringDecoder{},
		't': &trueDecoder{},
		'f': &falseDecoder{},
		'n': &nullDecoder{},
	}
	return p
}

type decoder interface {
	decode(s string, depth int, maxDepth int) (object.Object, string, error)
}

// parser : implement Parser
type parser struct {
	b        []byte
	decoders map[byte]decoder
}

func (this *parser) Parse(s string, maxDepth int) (object.Object, error) {
	s = skipWS(s)
	this.b = append(this.b[:0], s...)
	v, tail, err := this.parseValue(function.BytesToString(this.b), 0, maxDepth)
	if err != nil {
		return object.Nil, function.NewError(err)
	}
	tail = skipWS(tail)
	if len(tail) > 0 {
		return object.Nil, fmt.Errorf("unexpected tail: %q", tail)
	}
	return v, nil
}

func (this *parser) parseValue(s string, depth int, maxDepth int) (object.Object, string, error) {
	if len(s) == 0 {
		return nil, s, function.NewError(errEmptyString)
	}
	depth++
	if depth > maxDepth {
		return nil, s, function.NewError(errMaxDepth)
	}
	if d, ok := this.decoders[s[0]]; ok {
		v, tail, err := d.decode(s, depth, maxDepth)
		if err != nil {
			return nil, tail, function.NewError(err)
		}
		return v, tail, nil
	}
	v, tail, err := parseNumber(s)
	if err != nil {
		return nil, tail, function.NewError(err)
	}
	return v, tail, nil
}
