package ast

import (
	"encoding/json"
	"fmt"

	"github.com/jobs-github/escript/builtin"
	"github.com/jobs-github/escript/object"
)

// Identifier : implement Expression
type Identifier struct {
	defaultNode
	Value string
}

func (this *Identifier) Do(v Visitor) error {
	return v.DoIdent(this)
}

func (this *Identifier) getType() string {
	if fn := builtin.Get(this.Value); nil != fn {
		return typeExprBuiltin
	} else {
		return typeExprIdent
	}
}

func (this *Identifier) Encode() interface{} {
	return map[string]interface{}{
		keyType:  this.getType(),
		keyValue: this.Value,
	}
}
func (this *Identifier) Decode(b []byte) error {
	return json.Unmarshal(b, &this.Value)
}
func (this *Identifier) expressionNode() {}

func (this *Identifier) String() string {
	return this.Value
}

func (this *Identifier) Eval(e object.Env) (object.Object, error) {
	if val, ok := e.Get(this.Value); ok {
		return val, nil
	}
	if fn := builtin.Get(this.Value); nil != fn {
		return fn, nil
	}
	return object.Nil, fmt.Errorf("symbol `%v` missing", this.Value)
}

type IdentifierSlice []*Identifier

func (this *IdentifierSlice) encode() interface{} {
	arr := []interface{}{}
	for _, i := range *this {
		arr = append(arr, i.Encode())
	}
	return arr
}

func (this *IdentifierSlice) Values() []string {
	v := []string{}
	for _, i := range *this {
		v = append(v, i.Value)
	}
	return v
}
