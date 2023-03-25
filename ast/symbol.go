package ast

import (
	"encoding/json"
	"fmt"

	"github.com/jobs-github/escript/object"
)

// SymbolExpr : implement Expression
type SymbolExpr struct {
	defaultNode
	Value string
}

func (this *SymbolExpr) Do(v Visitor) error {
	return v.DoSymbol(this)
}

func (this *SymbolExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprSymbol,
		keyValue: this.Value,
	}
}
func (this *SymbolExpr) Decode(b []byte) error {
	return json.Unmarshal(b, &this.Value)
}
func (this *SymbolExpr) expressionNode() {}

func (this *SymbolExpr) String() string {
	return fmt.Sprintf("$%v", this.Value)
}

func (this *SymbolExpr) Eval(e object.Env) (object.Object, error) {
	cb, ok := e.Symbol(this.Value)
	if !ok {
		return object.Nil, fmt.Errorf("symbol `%v` missing", this.Value)
	}
	return cb()
}
