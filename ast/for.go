package ast

import (
	"bytes"
	"encoding/json"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/object"
)

// ForExpr : implement Expression
type ForExpr struct {
	Init Statement  // initialization statement; or nil
	Cond Expression // condition; or nil
	Post Statement  // post iteration statement; or nil
	Loop *BlockStmt
}

func (this *ForExpr) value() map[string]interface{} {
	m := map[string]interface{}{}
	if this.Init != nil {
		m["init"] = this.Init.Encode()
	}
	if this.Cond != nil {
		m["cond"] = this.Cond.Encode()
	}
	if this.Post != nil {
		m["post"] = this.Post.Encode()
	}
	if this.Loop != nil {
		m["loop"] = this.Loop.Encode()
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
		Init JsonNode `json:"init,omitempty"`
		Cond JsonNode `json:"cond,omitempty"`
		Post JsonNode `json:"post,omitempty"`
		Loop JsonNode `json:"loop"`
	}
	var err error
	if err = json.Unmarshal(b, &v); nil != err {
		return function.NewError(err)
	}
	if !v.Init.empty() {
		this.Init, err = v.Init.decodeStmt()
		if nil != err {
			return function.NewError(err)
		}
	}
	if !v.Cond.empty() {
		this.Cond, err = v.Cond.decodeExpr()
		if nil != err {
			return function.NewError(err)
		}
	}
	if !v.Post.empty() {
		this.Post, err = v.Post.decodeStmt()
		if nil != err {
			return function.NewError(err)
		}
	}
	this.Loop, err = v.Loop.decodeBlockStmt()
	if nil != err {
		return function.NewError(err)
	}
	return nil
}
func (this *ForExpr) expressionNode() {}

func (this *ForExpr) String() string {
	var out bytes.Buffer
	if nil == this.Init && nil == this.Cond {
		out.WriteString("for {")
	} else {
		out.WriteString("for (")
		if nil != this.Init {
			out.WriteString(this.Init.String())
		}
		out.WriteString(";")
		if nil != this.Cond {
			out.WriteString(this.Cond.String())
		}
		out.WriteString(";")
		if nil != this.Post {
			out.WriteString(this.Post.String())
		}
		out.WriteString(") {")
	}
	if nil != this.Loop {
		out.WriteString(this.Loop.String())
	}
	out.WriteString("}")
	return out.String()
}
func (this *ForExpr) Eval(env object.Env, insideLoop bool) (object.Object, error) {
	innerEnv := env.NewEnclosedEnv()
	if nil != this.Init {
		if _, err := this.Init.Eval(innerEnv, insideLoop); nil != err {
			return object.Nil, function.NewError(err)
		}
	}
	var rc object.Object
	var last object.Object
	for {
		if nil != this.Cond {
			cond, err := this.Cond.Eval(innerEnv, true)
			if nil != err {
				return object.Nil, function.NewError(err)
			}
			if !cond.True() {
				rc = last
				break
			}
		}

		v, err := this.Loop.Eval(innerEnv, true)
		if nil != err {
			return object.Nil, function.NewError(err)
		}
		last = v
		if isBreak, _ := v.Break(); isBreak {
			rc = v
			break
		}
		if needReturn, _ := v.Return(); needReturn {
			return v, nil
		}

		if nil != this.Post {
			if _, err := this.Post.Eval(innerEnv, true); nil != err {
				return object.Nil, function.NewError(err)
			}
		}
	}
	return rc, nil
}
func (this *ForExpr) walk(cb func(module string)) {
	if this.Loop != nil {
		this.Loop.walk(cb)
	}
}
func (this *ForExpr) doDefer(env object.Env) error { return nil }
