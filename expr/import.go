package expr

import (
	"path/filepath"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/interfaces"
	"github.com/jobs-github/escript/scanner"
	"github.com/jobs-github/escript/token"
)

// importExpr : implement tokenDecoder
type importExpr struct {
	scanner scanner.Scanner
	p       interfaces.Parser
}

func (this *importExpr) decodeAbbr() ast.Expression {
	this.scanner.NextToken()
	expr := this.scanner.NewImport()
	expr.ModuleName = this.scanner.GetIdentifier().Value
	expr.AsKey = filepath.Base(expr.ModuleName)
	return expr
}

func (this *importExpr) decode() (ast.Expression, error) {
	expr := this.scanner.NewImport()
	if nil == this.scanner.PeekIs(token.STRING) {
		return this.decodeAbbr(), nil
	}
	if err := this.scanner.ExpectPeek(token.IDENT); nil != err {
		return nil, function.NewError(err)
	}
	expr.AsKey = this.scanner.GetIdentifier().Value
	if err := this.scanner.ExpectPeek(token.STRING); nil != err {
		return nil, function.NewError(err)
	}
	expr.ModuleName = this.scanner.GetIdentifier().Value
	return expr, nil
}
