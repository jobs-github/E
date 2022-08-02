package builtin

import (
	"fmt"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

func builtinFor(args object.Objects) (object.Object, error) {
	argc := len(args)
	if argc != 5 {
		return object.Nil, fmt.Errorf("for() takes exactly 5 arguments (%v given)", argc)
	}
	i, err := object.ToInteger(args[0])
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	cond := args[1]
	if !object.Callable(cond) {
		return object.Nil, fmt.Errorf("the second argument (cond) should be callable (%v given)", object.Typeof(cond))
	}
	next := args[2]
	if !object.Callable(next) {
		return object.Nil, fmt.Errorf("the third argument (next) should be callable (%v given)", object.Typeof(next))
	}
	fn := args[3]
	if !object.Callable(fn) {
		return object.Nil, fmt.Errorf("the fourth argument (fn) should be callable (%v given)", object.Typeof(next))
	}
	state, err := args[4].AsState()
	if nil != err {
		return object.Nil, function.NewError(err)
	}
	iter := object.NewInteger(i)
	for {
		r, err := cond.Call(object.Objects{iter})
		if nil != err {
			return object.Nil, function.NewError(err)
		}
		if !r.True() {
			break
		}
		v, err := fn.Call(object.Objects{iter, state})
		if nil != err {
			return object.Nil, function.NewError(err)
		}
		if s, err := v.AsState(); nil != err {
			return object.Nil, function.NewError(err)
		} else {
			state = s
		}
		if state.Quit {
			break
		}
		nextVal, err := next.Call(object.Objects{iter})
		if nil != err {
			return object.Nil, function.NewError(err)
		}
		iter = nextVal
	}
	return state, nil
}
