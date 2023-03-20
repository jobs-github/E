package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
)

// filterExpr : implement tokenDecoder
type filterExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *filterExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewFilter()
	data, body, err := decodeLoopFn(this.scanner, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Arr = data
	expr.Body = body
	return expr, nil
}
