package object

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

func NewArray(items Objects) Object {
	obj := &Array{Items: items}
	obj.fns = objectBuiltins{
		FnSet:    obj.builtinSet,
		FnLen:    obj.builtinLen,
		FnIndex:  obj.builtinIndex,
		FnNot:    obj.builtinNot,
		FnMap:    obj.builtinMap,
		FnReduce: obj.builtinReduce,
		FnFilter: obj.builtinFilter,
		"first":  obj.builtinFirst,
		"last":   obj.builtinLast,
		"tail":   obj.builtinTail,
		"push":   obj.builtinPush,
	}
	return obj
}

// Array : implement Object
type Array struct {
	Items Objects
	fns   objectBuiltins
}

func (this *Array) String() string {
	var out bytes.Buffer
	items := []string{}
	for _, v := range this.Items {
		items = append(items, v.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")
	return out.String()
}

func (this *Array) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Array) Dump() (interface{}, error) {
	arr := []interface{}{}
	if nil == this.Items || len(this.Items) < 1 {
		return arr, nil
	}
	for _, item := range this.Items {
		v, err := item.Dump()
		if nil != err {
			return nil, function.NewError(err)
		}
		arr = append(arr, v)
	}
	return arr, nil
}

func (this *Array) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcArray(op, this)
}

func (this *Array) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *Array) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Array) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Array) True() bool {
	return false
}

func (this *Array) getType() ObjectType {
	return objectTypeArray
}

func (this *Array) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *Array) equal(other Object) error {
	return other.equalArray(this)
}

func (this *Array) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) equalArray(other *Array) error {
	szSrc := len(this.Items)
	szDst := len(other.Items)
	if szSrc != szDst {
		return fmt.Errorf("array size mismatch, this: %v, other: %v", szSrc, szDst)
	}
	for i := 0; i < szSrc; i++ {
		src := this.Items[i]
		dst := other.Items[i]
		if err := dst.equal(src); nil != err {
			return function.NewError(err)
		}
	}
	return nil
}

func (this *Array) equalHash(other *Hash) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Array) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Array) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Array) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Array) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Array) calcArray(op *token.Token, left *Array) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *Array) calcHash(op *token.Token, left *Hash) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Array) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Array) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Array) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *Array) builtinLen(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("len() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return NewInteger(int64(len(this.Items))), nil
}

func (this *Array) builtinSet(args Objects) (Object, error) {
	argc := len(args)
	if argc != 2 {
		return Nil, fmt.Errorf("set() takes 2 argument (%v given), (`%v`)", argc, this.String())
	}

	if nil == this.Items || len(this.Items) < 1 {
		return Nil, function.NewError(errArrayEmpty)
	}
	idx, err := args[0].asInteger()
	if nil != err {
		return Nil, function.NewError(err)
	}
	return setValue(this.Items, idx, args[1])
}

func (this *Array) builtinFirst(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("first() takes exactly no argument (%v given)", argc)
	}
	return this.Items.first()
}

func (this *Array) builtinLast(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("last() takes exactly no argument (%v given)", argc)
	}
	return this.Items.last()
}

func (this *Array) builtinTail(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("tail() takes exactly no argument (%v given)", argc)
	}
	return this.Items.tail()
}

func (this *Array) builtinPush(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("push() takes exactly one argument (%v given)", argc)
	}
	return this.Items.push(args[0])
}

func (this *Array) builtinIndex(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("index() takes exactly one argument (%v given)", argc)
	}
	if nil == this.Items || len(this.Items) < 1 {
		return Nil, function.NewError(errArrayEmpty)
	}
	idx, err := args[0].asInteger()
	if nil != err {
		return Nil, function.NewError(err)
	}
	return indexofArray(this.Items, idx)
}

func (this *Array) builtinNot(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("not() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if nil == this.Items || len(this.Items) < 1 {
		return True, nil
	} else {
		return False, nil
	}
}

func (this *Array) builtinMap(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("map() takes exactly one argument (%v given)", argc)
	}
	cb := args[0]
	if !Callable(cb) {
		return Nil, errors.New("argument is not callable")
	}
	if nil == this.Items || len(this.Items) < 1 {
		return NewArray(Objects{}), nil
	}
	r := Objects{}
	for i, item := range this.Items {
		v, err := cb.Call(Objects{NewInteger(int64(i)), item})
		if nil != err {
			return Nil, function.NewError(err)
		}
		r = append(r, v)
	}
	return NewArray(r), nil
}

func (this *Array) builtinReduce(args Objects) (Object, error) {
	argc := len(args)
	if argc != 2 {
		return Nil, fmt.Errorf("reduce() takes 2 argument (%v given), (`%v`)", argc, this.String())
	}
	cb := args[0]
	if !Callable(cb) {
		return Nil, errors.New("argument 1 is not callable")
	}
	acc := args[1]
	if nil == this.Items || len(this.Items) < 1 {
		return acc, nil
	}
	for _, item := range this.Items {
		v, err := cb.Call(Objects{acc, item})
		if nil != err {
			return Nil, function.NewError(err)
		}
		acc = v
	}
	return acc, nil
}

func (this *Array) builtinFilter(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("filter() takes exactly one argument (%v given), (`%v`)", argc, this.String())
	}
	cb := args[0]
	if !Callable(cb) {
		return Nil, errors.New("argument 1 is not callable")
	}
	r := Objects{}
	for _, item := range this.Items {
		v, err := cb.Call(Objects{item})
		if nil != err {
			return Nil, function.NewError(err)
		}
		if v.True() {
			r = append(r, item)
		}
	}
	return NewArray(r), nil
}
