package compiler

import (
	"github.com/jobs-github/escript/ast"
	"github.com/jobs-github/escript/code"
	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

func (this *visitor) doIdent(v *ast.Identifier) (int, error) {
	s, err := this.c.resolve(v.Value)
	if nil != err {
		return -1, function.NewError(err)
	}
	return this.doLoadSymbol(s)
}

func (this *visitor) DoIdent(v *ast.Identifier) error {
	if _, err := this.doIdent(v); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoNull(v *ast.Null) error {
	_, err := this.doConst(object.Nil)
	return err
}

func (this *visitor) DoInteger(v *ast.Integer) error {
	_, err := this.doConst(object.NewInteger(v.Value))
	return err
}

func (this *visitor) DoBoolean(v *ast.Boolean) error {
	if _, err := this.c.encode(this.opCodeBoolean(v)); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoString(v *ast.String) error {
	_, err := this.doConst(object.NewString(v.Value))
	return err
}

func (this *visitor) DoArray(v *ast.Array) error {
	// pattern: compile data first, op last
	for _, e := range v.Items {
		if err := e.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	if _, err := this.c.encode(code.OpArray, len(v.Items)); nil != err {
		return function.NewError(err)
	}
	return nil
}

func (this *visitor) DoHash(v *ast.Hash) error {
	keys := v.Pairs.SortedKeys()
	for _, k := range keys {
		if err := k.Do(this); nil != err {
			return function.NewError(err)
		}
		v := v.Pairs[k]
		if err := v.Do(this); nil != err {
			return function.NewError(err)
		}
	}
	if _, err := this.c.encode(code.OpHash, len(v.Pairs)); nil != err {
		return function.NewError(err)
	}
	return nil
}
