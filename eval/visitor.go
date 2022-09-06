package eval

import (
	"fmt"

	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/builtin"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
	"github.com/jobs-github/escript/token"
)

func evalBody(v *ast.Function) func(env object.Env) (object.Object, error) {
	return func(env object.Env) (object.Object, error) {
		return EvalAst(v.Body, env)
	}
}

func getCondNode(v *ast.ConditionalExpr, yes bool) ast.Expression {
	if yes {
		return v.Yes
	} else {
		return v.No
	}
}

func newVisitor(env object.Env) *visitor {
	return &visitor{e: env, r: object.Nil}
}

// visitor : implement ast.Visitor
type visitor struct {
	e object.Env
	r object.Object
}

func (this *visitor) enclosedVisitor() *visitor {
	return newVisitor(this.e.NewEnclosedEnv())
}

func (this *visitor) doVar(name *ast.Identifier, value ast.Expression) error {
	err := value.Do(this)
	if nil != err {
		return function.NewError(err)
	}
	this.e.Set(name.Value, this.r)
	return nil
}

func (this *visitor) doPrefixExpr(op *token.Token, right object.Object) (object.Object, error) {
	switch op.Type {
	case token.NOT:
		return right.CallMember(object.FnNot, object.Objects{})
	case token.SUB:
		return right.CallMember(object.FnNeg, object.Objects{})
	default:
		err := fmt.Errorf("unsupport op %v(%v)", op.Literal, token.ToString(op.Type))
		return object.Nil, function.NewError(err)
	}
}

func (this *visitor) doExprs(v ast.ExpressionSlice) (object.Objects, error) {
	result := object.Objects{}
	for _, expr := range v {
		if err := expr.Do(this); nil != err {
			return nil, function.NewError(err)
		}
		result.Append(this.r)
	}
	return result, nil
}

func (this *visitor) doStmts(v ast.StatementSlice) error {
	for _, stmt := range v {
		if err := stmt.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func (this *visitor) DoProgram(v *ast.Program) error {
	return this.doStmts(v.Stmts)
}

func (this *visitor) DoConst(v *ast.ConstStmt) error {
	return this.doVar(v.Name, v.Value)
}

func (this *visitor) DoBlock(v *ast.BlockStmt) error {
	return v.Stmt.Do(this)
}

func (this *visitor) DoExpr(v *ast.ExpressionStmt) error {
	return v.Expr.Do(this)
}

func (this *visitor) DoFunction(v *ast.FunctionStmt) error {
	return this.doVar(v.Name, v.Value)
}

func (this *visitor) DoPrefix(v *ast.PrefixExpr) error {
	if err := v.Right.Do(this); nil != err {
		return function.NewError(err)
	}
	if r, err := this.doPrefixExpr(v.Op, this.r); nil != err {
		return function.NewError(err)
	} else {
		this.r = r
		return nil
	}
}

func (this *visitor) DoInfix(v *ast.InfixExpr) error {
	if err := v.Left.Do(this); nil != err {
		return function.NewError(err)
	}
	left := this.r
	if err := v.Right.Do(this); nil != err {
		return function.NewError(err)
	}
	right := this.r
	if r, err := left.Calc(v.Op, right); nil != err {
		return function.NewError(err)
	} else {
		this.r = r
		return nil
	}
}

func (this *visitor) DoIdent(v *ast.Identifier) error {
	if val, ok := this.e.Get(v.Value); ok {
		this.r = val
		return nil
	}
	if fn := builtin.Get(v.Value); nil != fn {
		this.r = fn
		return nil
	}
	err := fmt.Errorf("`%v` not found", v.Value)
	return function.NewError(err)
}

func (this *visitor) DoConditional(v *ast.ConditionalExpr) error {
	if err := v.Cond.Do(this); nil != err {
		return function.NewError(err)
	}
	innerVisitor := this.enclosedVisitor()
	node := getCondNode(v, this.r.True())
	if err := node.Do(innerVisitor); nil != err {
		return function.NewError(err)
	} else {
		// important
		this.r = innerVisitor.r
		return nil
	}
}

func (this *visitor) DoFn(v *ast.Function) error {
	this.r = object.NewFunction(
		v.Name,
		v.Args.Values(),
		evalBody(v),
		this.e,
	)
	return nil
}

func (this *visitor) DoCall(v *ast.Call) error {
	if err := v.Func.Do(this); nil != err {
		return function.NewError(err)
	}
	fn := this.r
	args, err := this.doExprs(v.Args)
	if nil != err {
		return function.NewError(err)
	}
	if r, err := fn.Call(args); nil != err {
		return function.NewError(err)
	} else {
		this.r = r
		return nil
	}
}

func (this *visitor) DoCallMember(v *ast.CallMember) error {
	if err := v.Left.Do(this); nil != err {
		return function.NewError(err)
	}
	obj := this.r
	args, err := this.doExprs(v.Args)
	if nil != err {
		return function.NewError(err)
	}
	if r, err := obj.CallMember(v.Func.Value, args); nil != err {
		return function.NewError(err)
	} else {
		this.r = r
		return nil
	}
}

func (this *visitor) DoObjectMember(v *ast.ObjectMember) error {
	if err := v.Left.Do(this); nil != err {
		return function.NewError(err)
	}
	obj := this.r
	if r, err := obj.GetMember(v.Member.Value); nil != err {
		return function.NewError(err)
	} else {
		this.r = r
		return nil
	}
}

func (this *visitor) DoIndex(v *ast.IndexExpr) error {
	if err := v.Left.Do(this); nil != err {
		return function.NewError(err)
	}
	left := this.r
	if err := v.Index.Do(this); nil != err {
		return function.NewError(err)
	}
	idx := this.r
	if r, err := left.CallMember(object.FnIndex, object.Objects{idx}); nil != err {
		return function.NewError(err)
	} else {
		this.r = r
		return nil
	}
}

func (this *visitor) DoNull(v *ast.Null) error {
	this.r = object.Nil
	return nil
}

func (this *visitor) DoInteger(v *ast.Integer) error {
	this.r = object.NewInteger(v.Value)
	return nil
}

func (this *visitor) DoBoolean(v *ast.Boolean) error {
	this.r = object.ToBoolean(v.Value)
	return nil
}

func (this *visitor) DoString(v *ast.String) error {
	this.r = object.NewString(v.Value)
	return nil
}

func (this *visitor) DoArray(v *ast.Array) error {
	items, err := this.doExprs(v.Items)
	if nil != err {
		return function.NewError(err)
	}
	this.r = object.NewArray(items)
	return nil
}

func (this *visitor) DoHash(v *ast.Hash) error {
	m := object.HashMap{}
	for k, v := range v.Pairs {
		if err := k.Do(this); nil != err {
			return function.NewError(err)
		}
		key := this.r
		h, err := key.Hash()
		if nil != err {
			return function.NewError(err)
		}
		if err := v.Do(this); nil != err {
			return function.NewError(err)
		}
		val := this.r
		m[*h] = &object.HashPair{Key: key, Value: val}
	}
	this.r = object.NewHash(m)
	return nil
}
