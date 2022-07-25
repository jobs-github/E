package json

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// reference : https://github.com/valyala/fastjson

const (
	MaxDepth = 300
)

var (
	errEmptyString        = errors.New("cannot parse empty string")
	errMissingClosing     = errors.New(`missing closing '"'`)
	errMaxDepth           = errors.New("too big depth for the nested JSON; it exceeds max depth")
	errMissingObjectEnd   = errors.New("missing '}'")
	errObjectStart        = errors.New(`cannot find opening '"" for object key`)
	errObjectKeySeparator = errors.New("missing ':' after object key")
	errObjectValSeparator = errors.New("missing ',' after object value")
	errObjectEnd          = errors.New("unexpected end of object")
	errMissingArrayEnd    = errors.New("missing ']'")
	errArrayEnd           = errors.New("unexpected end of array")
	errArrayValSeparator  = errors.New("missing ',' after array value")
)

func skipWSSlow(s string) string {
	if len(s) == 0 || s[0] != 0x20 && s[0] != 0x0A && s[0] != 0x09 && s[0] != 0x0D {
		return s
	}
	for i := 1; i < len(s); i++ {
		if s[i] != 0x20 && s[i] != 0x0A && s[i] != 0x09 && s[i] != 0x0D {
			return s[i:]
		}
	}
	return ""
}

func skipWS(s string) string {
	if len(s) == 0 || s[0] > 0x20 {
		return s
	}
	return skipWSSlow(s)
}

func parseRawString(s string) (object.Object, string, error) {
	n := strings.IndexByte(s, '"')
	if n < 0 {
		return nil, "", function.NewError(errMissingClosing)
	}
	if n == 0 || s[n-1] != '\\' {
		return object.NewString(s[:n]), s[n+1:], nil
	}

	ss := s
	for {
		i := n - 1
		for i > 0 && s[i-1] == '\\' {
			i--
		}
		if uint(n-i)%2 == 0 {
			return object.NewString(ss[:len(ss)-len(s)+n]), s[n+1:], nil
		}
		s = s[n+1:]

		n = strings.IndexByte(s, '"')
		if n < 0 {
			return object.NewString(ss), "", function.NewError(errMissingClosing)
		}
		if n == 0 || s[n-1] != '\\' {
			return object.NewString(ss[:len(ss)-len(s)+n]), s[n+1:], nil
		}
	}
}

func parseRawKey(s string) (object.Object, string, error) {
	sz := len(s)
	for i := 0; i < sz; i++ {
		if s[i] == '"' {
			return object.NewString(s[:i]), s[i+1:], nil
		}
		if s[i] == '\\' {
			return parseRawString(s)
		}
	}
	return nil, "", function.NewError(errMissingClosing)
}

func isDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || ch == '-'
}

func parseNumber(s string) (object.Object, string, error) {
	sz := len(s)
	var i int
	for i = 0; i < sz; i++ {
		if !isDigit(s[i]) {
			break
		}
	}
	ns := s[:i]
	s = s[i:]
	v, err := strconv.ParseInt(ns, 10, 64)
	if nil != err {
		return nil, s, function.NewError(err)
	}
	return object.NewInteger(v), s, nil
}
