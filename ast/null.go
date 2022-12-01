package ast

import (
	"github.com/jobs-github/escript/token"
)

// Null : implement Expression
type Null struct{ defaultNode }

func (this *Null) Do(v Visitor) error {
	return v.DoNull(this)
}

func (this *Null) Encode() interface{} {
	return map[string]interface{}{
		keyType: typeExprNull,
	}
}
func (this *Null) Decode(b []byte) error {
	return nil
}
func (this *Null) expressionNode() {}

func (this *Null) String() string {
	return token.Null
}
