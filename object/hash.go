package object

import (
	"bytes"
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
		FnLen:   obj.builtinLen,
		FnIndex: obj.builtinIndex,
		FnNot:   obj.builtinNot,
		FnKeys:  obj.builtinKeys,
	}
	return obj
}

// Hash : implement Object
type Hash struct {
	defaultObject
	Pairs HashMap
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
			return nil, err
		}
		v, err := item.Value.Dump()
		if nil != err {
			return nil, err
		}
		m[item.Key.String()] = v
	}
	return m, nil
}

func (this *Hash) Calc(op *token.Token, right Object) (Object, error) {
	return right.calcHash(op, this)
}

func (this *Hash) CallMember(name string, args Objects) (Object, error) {
	return callMember(this, this.fns, name, args)
}

func (this *Hash) GetMember(name string) (Object, error) {
	return getMember(this, this.fns, name)
}

func (this *Hash) getType() ObjectType {
	return objectTypeHash
}

func (this *Hash) equal(other Object) error {
	return other.equalHash(this)
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
				return err
			}
		}
	}
	return nil
}

func (this *Hash) calcHash(op *token.Token, left *Hash) (Object, error) {
	return compare(function.GetFunc(), this, left, op)
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
