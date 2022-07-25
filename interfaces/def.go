package interfaces

import (
	"github.com/jobs-github/Q/ast"
	"github.com/jobs-github/Q/token"
)

type Parser interface {
	ParseProgram() (ast.Node, error)

	ParseInfixExpression(left ast.Expression) (ast.Expression, error)
	ParseCallExpression(left ast.Expression) (ast.Expression, error)
	ParseIndexExpression(left ast.Expression) (ast.Expression, error)
	ParseMemberExpression(left ast.Expression) (ast.Expression, error)
	ParseConditionalExpression(left ast.Expression) (ast.Expression, error)

	ParseStmt(endTok token.TokenType) (ast.Statement, error)
	ParseBlockStmt() (*ast.BlockStmt, error)
	ParseExpression(precedence int) (ast.Expression, error)
	ParseExpressions(endTok token.TokenType) (ast.ExpressionSlice, error)
}
