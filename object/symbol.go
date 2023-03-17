package object

type symbolTable []string

func (this *symbolTable) traverse(cb func(i int, name string)) {
	for i, fn := range *this {
		cb(i, fn)
	}
}

var (
	// vm: array
	objectSymbolTable = symbolTable{
		FnLen,
		FnIndex,
		FnNot,
		FnNeg,
		FnInt,
		FnFilter,
		FnFirst,
		FnLast,
		FnTail,
		FnPush,
		FnKeys,
	}
)

func Resolve(idx int) string {
	return objectSymbolTable[idx]
}

func Traverse(cb func(i int, name string)) {
	objectSymbolTable.traverse(cb)
}
