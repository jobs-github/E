package object

// type symbolIndex map[string]int
type symbolTable []string

func (this *symbolTable) traverse(cb func(i int, name string)) {
	for i, fn := range *this {
		cb(i, fn)
	}
}

/*
func (this *symbolTable) newIndex() symbolIndex {
	m := symbolIndex{}
	for i, s := range *this {
		m[s] = i
	}
	return m
}
*/
var (
	// vm: array
	objectSymbolTable = symbolTable{
		FnLen,
		FnIndex,
		FnNot,
		FnNeg,
		FnInt,
		FnMap,
		FnReduce,
		FnFilter,
		FnFirst,
		FnLast,
		FnTail,
		FnPush,
		FnKeys,
		FnValue,
	}

	// compiler: map[string]int
	//objectSymbolIndex = objectSymbolTable.newIndex()
)

/*
func SymbolName(index int) string {
	return objectSymbolTable[index]
}

func SymbolIndex(fn string) int {
	return objectSymbolIndex[fn]
}
*/
func Resolve(idx int) string {
	return objectSymbolTable[idx]
}

func Traverse(cb func(i int, name string)) {
	objectSymbolTable.traverse(cb)
}
