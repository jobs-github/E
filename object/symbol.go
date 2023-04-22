package object

type Symbols map[string]Callable

type SymbolTable map[string]Object

type objectSymbols []string

func (this *objectSymbols) traverse(cb func(i int, name string)) {
	for i, fn := range *this {
		cb(i, fn)
	}
}

var (
	objectSymbolTable = objectSymbols{
		FnLen,
		FnIndex,
		FnNot,
		FnNeg,
		FnInt,
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
