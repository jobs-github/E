package object

import "fmt"

type objectFnSlice []string

func (this *objectFnSlice) reverseIndex() objectFnMap {
	m := objectFnMap{}
	for i, v := range *this {
		m[v] = uint8(i)
	}
	return m
}

type objectFnMap map[string]uint8

type objectsFnSlice []objectFnSlice

func (this *objectsFnSlice) get(t int, idx int) string {
	return (*this)[t][idx]
}

func (this *objectsFnSlice) reverseIndex() objectFnMapSlice {
	sz := len(*this)
	r := make(objectFnMapSlice, sz)
	for i, v := range *this {
		r[i] = v.reverseIndex()
	}
	return r
}

type objectFnMapSlice []objectFnMap

func (this *objectFnMapSlice) get(t ObjectType, fn string) (uint8, error) {
	i, ok := (*this)[t][fn]
	if !ok {
		return 0, fmt.Errorf("unsupported func: %v, type: %v", fn, t)
	}
	return i, nil
}

var (
	// vm
	objectsFns = objectsFnSlice{
		objectFnSlice{ // objectTypeBuiltin
			FnNot,
		},
		objectFnSlice{ // objectTypeInteger
			FnNot,
			FnNeg,
			FnInt,
		},
		objectFnSlice{ // objectTypeString
			FnLen,
			FnIndex,
			FnNot,
			FnInt,
		},
		objectFnSlice{ // objectTypeBoolean
			FnNot,
			FnNeg,
			FnInt,
		},
		objectFnSlice{ // objectTypeNull
			FnNot,
		},
		objectFnSlice{ // objectTypeFunction
			FnNot,
		},
		objectFnSlice{ // objectTypeByteFunc
			FnNot,
		},
		objectFnSlice{ // objectTypeClosure
			FnNot,
		},
		objectFnSlice{ // objectTypeArray
			FnLen,
			FnIndex,
			FnNot,
			FnMap,
			FnReduce,
			FnFilter,
			FnFirst,
			FnLast,
			FnTail,
			FnPush,
		},
		objectFnSlice{ // objectTypeHash
			FnLen,
			FnIndex,
			FnNot,
			FnMap,
			FnReduce,
			FnFilter,
			FnKeys,
		},
		objectFnSlice{ // objectTypeObjectFunc
			FnNot,
		},
		objectFnSlice{ // objectTypeState
			FnValue,
		},
	}

	// compiler
	objectsFnMap = objectsFns.reverseIndex()
)

func GetObjectFuncName(t ObjectType, idx uint8) string {
	return objectsFns.get(int(t), int(idx))
}

func GetObjectFuncIndex(t ObjectType, fn string) (uint8, error) {
	return objectsFnMap.get(t, fn)
}
