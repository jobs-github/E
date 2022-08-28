package ast

import (
	"encoding/json"

	"github.com/jobs-github/escript/builtin"
)

// Identifier : implement Expression
type Identifier struct {
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
