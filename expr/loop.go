package expr

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
)

// loopExpr : implement tokenDecoder
type loopExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *loopExpr) decode() (ast.Expression, error) {
	expr := ast.NewLoop()
	data, body, err := decodeLoopFn(this.scanner, this.p)
	if nil != err {
		return nil, function.NewError(err)
	}
	expr.Cnt = data
	expr.Body = body
	return expr, nil
}
