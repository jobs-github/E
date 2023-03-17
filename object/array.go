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
		FnLen:    obj.builtinLen,
		FnIndex:  obj.builtinIndex,
		FnNot:    obj.builtinNot,
		FnFilter: obj.builtinFilter,
		FnFirst:  obj.builtinFirst,
		FnLast:   obj.builtinLast,
		FnTail:   obj.builtinTail,
		FnPush:   obj.builtinPush,
	}
	return obj
}

// Array : implement Object
type Array struct {
	defaultObject
	Items Objects
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

func (this *Array) Dump() (interface{}, error) {
	arr := []interface{}{}
	if nil == this.Items || len(this.Items) < 1 {
		return arr, nil
	}
	for _, item := range this.Items {
		v, err := item.Dump()
		if nil != err {
			return nil, err
		}
		arr = append(arr, v)
	}
	return arr, nil
}

func (this *Array) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcArray(op, this)
}

func (this *Array) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Array) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Array) AsArray() (*Array, error) {
	return this, nil
}

func (this *Array) Set(idx Object, val Object) error {
	i, err := idx.asInteger()
	if nil != err {
		return err
	}
	sz := int64(len(this.Items))
	if i < 0 || i > sz-1 {
		return errInvalidIndex
	}
	this.Items[i] = val
	return nil
}

func (this *Array) New() Object {
	if this.Items == nil || len(this.Items) == 0 {
		return NewArray(Objects{})
	}
	arr := make(Objects, len(this.Items))
	return NewArray(arr)
}

func (this *Array) getType() ObjectType {
	return objectTypeArray
}

func (this *Array) equal(other Object) error {
	return other.equalArray(this)
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
			return err
		}
	}
	return nil
}

func (this *Array) calcArray(op *token.Token, left *Array) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

// builtin
func (this *Array) builtinLen(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("len() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return NewInteger(int64(len(this.Items))), nil
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
		return Nil, err
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

func (this *Array) builtinFilter(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("filter() takes exactly one argument (%v given), (`%v`)", argc, this.String())
	}
	cb := args[0]
	if !Callable(cb) {
		return Nil, errors.New("argument 1 is not callable")
	}
	if nil == this.Items || len(this.Items) < 1 {
		return NewArray(Objects{}), nil
	}
	r := Objects{}
	for _, item := range this.Items {
		v, err := cb.Call(Objects{item})
		if nil != err {
			return Nil, err
		}
		if v.True() {
			r = append(r, item)
		}
	}
	return NewArray(r), nil
}
