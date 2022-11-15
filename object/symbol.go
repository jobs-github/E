package object

import "fmt"

type symbolSlice []string

func (this *symbolSlice) reverseIndex() symbolMap {
	m := symbolMap{}
	for i, v := range *this {
		m[v] = uint8(i)
	}
	return m
}

type symbolMap map[string]uint8

type symbolTable []symbolSlice

func (this *symbolTable) get(t int, idx int) string {
	return (*this)[t][idx]
}

func (this *symbolTable) reverseIndex() symbolIndex {
	sz := len(*this)
	r := make(symbolIndex, sz)
	for i, v := range *this {
		r[i] = v.reverseIndex()
	}
	return r
}

type symbolIndex []symbolMap

func (this *symbolIndex) get(t ObjectType, fn string) (uint8, error) {
	i, ok := (*this)[t][fn]
	if !ok {
		return 0, fmt.Errorf("unsupported func: %v, type: %v", fn, t)
	}
	return i, nil
}

var (
	// vm
	objectsSymbolTable = symbolTable{
		symbolSlice{ // objectTypeBuiltin
			FnNot,
		},
		symbolSlice{ // objectTypeInteger
			FnNot,
			FnNeg,
			FnInt,
		},
		symbolSlice{ // objectTypeString
			FnLen,
			FnIndex,
			FnNot,
			FnInt,
		},
		symbolSlice{ // objectTypeBoolean
			FnNot,
			FnNeg,
			FnInt,
		},
		symbolSlice{ // objectTypeNull
			FnNot,
		},
		symbolSlice{ // objectTypeFunction
			FnNot,
		},
		symbolSlice{ // objectTypeByteFunc
			FnNot,
		},
		symbolSlice{ // objectTypeClosure
			FnNot,
		},
		symbolSlice{ // objectTypeArray
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
		symbolSlice{ // objectTypeHash
			FnLen,
			FnIndex,
			FnNot,
			FnMap,
			FnReduce,
			FnFilter,
			FnKeys,
		},
		symbolSlice{ // objectTypeObjectFunc
			FnNot,
		},
		symbolSlice{ // objectTypeState
			FnValue,
		},
	}

	// compiler
	objectSymbolIndex = objectsSymbolTable.reverseIndex()
)

func SymbolName(t ObjectType, idx uint8) string {
	return objectsSymbolTable.get(int(t), int(idx))
}

func SymbolIndex(t ObjectType, fn string) (uint8, error) {
	return objectSymbolIndex.get(t, fn)
}
