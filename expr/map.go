package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
)

// mapExpr : implement tokenDecoder
type mapExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *mapExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewMap()
	data, body, err := decodeLoopFn(this.scanner, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Arr = data
	expr.Body = body
	return expr, nil
}
