package object

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jobs-github/escript/function"
	"github.com/jobs-github/escript/token"
)

type HashPair struct {
	Key   Object
	Value Object
}

type HashMap map[HashKey]*HashPair

func (this *HashMap) get(k *HashKey) (*HashPair, bool) {
	v, ok := (*this)[*k]
	return v, ok
}

func (this *HashMap) set(k *HashKey, v *HashPair) {
	(*this)[*k] = v
}

func (this *HashMap) keys() Objects {
	arr := Objects{}
	for _, v := range *this {
		arr = append(arr, v.Key)
	}
	sort.SliceStable(arr, func(i, j int) bool {
		if r, err := arr[i].Calc(&token.Token{Type: token.LT}, arr[j]); nil != err {
			return false
		} else {
			return r.True()
		}
	})
	return arr
}

func NewHash(pairs HashMap) Object {
	obj := &Hash{
		Pairs: pairs,
	}
	obj.fns = objectBuiltins{
		FnLen:    obj.builtinLen,
		FnIndex:  obj.builtinIndex,
		FnNot:    obj.builtinNot,
		FnMap:    obj.builtinMap,
		FnReduce: obj.builtinReduce,
		FnFilter: obj.builtinFilter,
		"keys":   obj.builtinKeys,
	}
	return obj
}

// Hash : implement Object
type Hash struct {
	Pairs HashMap
	fns   objectBuiltins
}

func (this *Hash) String() string {
	var out bytes.Buffer
	items := []string{}
	for _, v := range this.Pairs {
		items = append(items, fmt.Sprintf("%v: %v", v.Key.String(), v.Value.String()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("}")
	return out.String()
}

func (this *Hash) Hash() (*HashKey, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Hash) Dump() (interface{}, error) {
	m := map[string]interface{}{}
	if nil == this.Pairs || len(this.Pairs) < 1 {
		return m, nil
	}
	for _, item := range this.Pairs {
		if !IsString(item.Key) {
			err := fmt.Errorf("`%v` (%v) is not string", item.Key.String(), Typeof(item.Key))
			return nil, function.NewError(err)
		}
		v, err := item.Value.Dump()
		if nil != err {
			return nil, function.NewError(err)
		}
		m[item.Key.String()] = v
	}
	return m, nil
}

func (this *Hash) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcHash(op, this)
}

func (this *Hash) Call(args Objects) (Object, error) {
	return Nil, unsupported(function.GetFunc(), this)
}

func (this *Hash) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Hash) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Hash) True() bool {
	return false
}

func (this *Hash) AsState() (*State, error) {
	return nil, unsupported(function.GetFunc(), this)
}

func (this *Hash) getType() ObjectType {
	return objectTypeHash
}

func (this *Hash) asInteger() (int64, error) {
	return 0, unsupported(function.GetFunc(), this)
}

func (this *Hash) equal(other Object) error {
	return other.equalHash(this)
}

func (this *Hash) equalInteger(other *Integer) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalString(other *String) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalBoolean(other *Boolean) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalNull(other *Null) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalArray(other *Array) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalHash(other *Hash) error {
	szSrc := len(this.Pairs)
	szDst := len(other.Pairs)
	if szSrc != szDst {
		return fmt.Errorf("hash size mismatch, this: %v, other: %v", szSrc, szDst)
	}
	for k, valSrc := range this.Pairs {
		if valDst, ok := other.Pairs[k]; !ok {
			return fmt.Errorf("other hash missing key `%v`", valSrc.Key.String())
		} else {
			if err := valDst.Value.equal(valSrc.Value); nil != err {
				return function.NewError(err)
			}
		}
	}
	return nil
}

func (this *Hash) equalBuiltin(other *Builtin) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalFunction(other *Function) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalByteFunc(other *ByteFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) equalObjectFunc(other *ObjectFunc) error {
	return fmt.Errorf("type mismatch, this: %v, other: %v", Typeof(this), Typeof(other))
}

func (this *Hash) calcInteger(op *token.Token, left *Integer) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcString(op *token.Token, left *String) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcBoolean(op *token.Token, left *Boolean) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcNull(op *token.Token, left *Null) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcArray(op *token.Token, left *Array) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcHash(op *token.Token, left *Hash) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
}

func (this *Hash) calcBuiltin(op *token.Token, left *Builtin) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcFunction(op *token.Token, left *Function) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcByteFunc(op *token.Token, left *ByteFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

func (this *Hash) calcObjectFunc(op *token.Token, left *ObjectFunc) (Object, error) {
	return notEqual(function.GetFunc(), this, op)
}

// builtin
func (this *Hash) builtinLen(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("len() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	return NewInteger(int64(len(this.Pairs))), nil
}

func (this *Hash) builtinIndex(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("index() takes exactly one argument (%v given)", argc)
	}
	idx := args[0]
	if nil == this.Pairs || len(this.Pairs) < 1 {
		return Nil, function.NewError(errHashEmpty)
	}
	return keyofHash(this.Pairs, idx)
}

func (this *Hash) builtinNot(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("not() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if nil == this.Pairs || len(this.Pairs) < 1 {
		return True, nil
	} else {
		return False, nil
	}
}

func (this *Hash) builtinMap(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("map() takes exactly one argument (%v given)", argc)
	}
	cb := args[0]
	if !Callable(cb) {
		return Nil, errors.New("argument is not callable")
	}
	if nil == this.Pairs || len(this.Pairs) < 1 {
		return NewHash(HashMap{}), nil
	}
	m := HashMap{}
	for k, v := range this.Pairs {
		val, err := cb.Call(Objects{v.Key, v.Value})
		if nil != err {
			return Nil, function.NewError(err)
		}
		m[k] = &HashPair{v.Key, val}
	}
	return NewHash(m), nil
}

func (this *Hash) builtinReduce(args Objects) (Object, error) {
	argc := len(args)
	if argc != 2 {
		return Nil, fmt.Errorf("reduce() takes 2 arguments (%v given), (`%v`)", argc, this.String())
	}
	cb := args[0]
	if !Callable(cb) {
		return Nil, errors.New("argument 1 is not callable")
	}
	acc := args[1]
	if nil == this.Pairs || len(this.Pairs) < 1 {
		return acc, nil
	}
	for _, v := range this.Pairs {
		val, err := cb.Call(Objects{acc, v.Value})
		if nil != err {
			return Nil, function.NewError(err)
		}
		acc = val
	}
	return acc, nil
}

func (this *Hash) builtinFilter(args Objects) (Object, error) {
	argc := len(args)
	if argc != 1 {
		return Nil, fmt.Errorf("filter() takes exactly one argument (%v given), (`%v`)", argc, this.String())
	}
	cb := args[0]
	if !Callable(cb) {
		return Nil, errors.New("argument 1 is not callable")
	}
	if nil == this.Pairs || len(this.Pairs) < 1 {
		return NewHash(HashMap{}), nil
	}
	m := HashMap{}
	for k, v := range this.Pairs {
		val, err := cb.Call(Objects{v.Key, v.Value})
		if nil != err {
			return Nil, function.NewError(err)
		}
		if val.True() {
			m[k] = v
		}
	}
	return NewHash(m), nil
}

func (this *Hash) builtinKeys(args Objects) (Object, error) {
	argc := len(args)
	if argc != 0 {
		return Nil, fmt.Errorf("keys() takes no argument (%v given), (`%v`)", argc, this.String())
	}
	if nil == this.Pairs || len(this.Pairs) < 1 {
		return NewArray(Objects{}), nil
	}
	return NewArray(this.Pairs.keys()), nil
}
