package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// Hash : implement Expression
type Hash struct {
	defaultNode
	Pairs ExpressionMap
}

func (this *Hash) Do(v Visitor) error {
	return v.DoHash(this)
}

func (this *Hash) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprHash,
		keyValue: this.Pairs.encode(),
	}
}
func (this *Hash) Decode(b []byte) error {
	var err error
	this.Pairs, err = decodeExprMap(b)
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *Hash) expressionNode() {}

func (this *Hash) String() string {
	var out bytes.Buffer
	items := []string{}
	for k, v := range this.Pairs {
		items = append(items, fmt.Sprintf("%v:%v", k.String(), v.String()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("}")
	return out.String()
}

func (this *Hash) Eval(e object.Env) (object.Object, error) {
	m := object.HashMap{}
	for k, v := range this.Pairs {
		key, err := k.Eval(e)
		if nil != err {
			return object.Nil, err
		}
		h, err := key.Hash()
		if nil != err {
			return object.Nil, err
		}
		val, err := v.Eval(e)
		if nil != err {
			return object.Nil, err
		}
		m[*h] = &object.HashPair{Key: key, Value: val}
	}
	return object.NewHash(m), nil
}
