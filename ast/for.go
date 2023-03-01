package ast

import (
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// ForExpr : implement Expression
type ForExpr struct {
	defaultNode
	Init   Expression // initial state
	Cond   Expression // Function or Identifier
	Next   Expression // Function or Identifier
	LoopFn Expression // Function or Identifier
}

func (this *ForExpr) expressionNode() {}

func (this *ForExpr) Do(v Visitor) error {
	return v.DoFor(this)
}

func (this *ForExpr) value() map[string]interface{} {
	m := map[string]interface{}{
		"init": this.Init.Encode(),
		"cond": this.Cond.Encode(),
		"next": this.Next.Encode(),
		"loop": this.LoopFn.Encode(),
	}
	return m
}

func (this *ForExpr) Encode() interface{} {
	return map[string]interface{}{
		keyType:  typeExprFor,
		keyValue: this.value(),
	}
}
func (this *ForExpr) Decode(b []byte) error {
	var v struct {
		Init JsonNode `json:"init"`
		Cond JsonNode `json:"cond"`
		Next JsonNode `json:"next"`
		Loop JsonNode `json:"loop"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	this.Init, err = v.Init.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Cond, err = v.Cond.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.Next, err = v.Next.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	this.LoopFn, err = v.Loop.decodeExpr()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ForExpr) String() string {
	return ""
}
func (this *ForExpr) Eval(e object.Env) (object.Object, error) {
	v, err := this.Init.Eval(e)
	if nil != err {
		return object.Nil, err
	}
	state, err := v.AsState()
	if nil != err {
		return object.Nil, err
	}
	cond, err := this.Cond.Eval(e)
	if !object.Callable(cond) {
		return object.Nil, err
	}
	next, err := this.Next.Eval(e)
	if !object.Callable(next) {
		return object.Nil, err
	}
	fn, err := this.LoopFn.Eval(e)
	if !object.Callable(fn) {
		return object.Nil, err
	}
	return this.do(state, cond, next, fn)
}

func (this *ForExpr) do(
	state *object.State,
	cond object.Object,
	next object.Object,
	fn object.Object,
) (object.Object, error) {
	i, err := object.ToInteger(state.Value)
	if nil != err {
		return object.Nil, err
	}
	iter := object.NewInteger(i)
	for {
		r, err := cond.Call(object.Objects{iter})
		if nil != err {
			return object.Nil, err
		}
		if !r.True() {
			break
		}
		v, err := fn.Call(object.Objects{iter, state})
		if nil != err {
			return object.Nil, err
		}
		if s, err := v.AsState(); nil != err {
			return object.Nil, err
		} else {
			state = s
		}
		if state.Quit {
			break
		}
		nextVal, err := next.Call(object.Objects{iter})
		if nil != err {
			return object.Nil, err
		}
		iter = nextVal
	}
	return state, nil
}
