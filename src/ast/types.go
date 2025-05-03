package ast

type SymbolType struct {
	Name string // [number | string]
}

func (SymbolType) _type() {}

// "[number]", "[number(a)]"
type ArrayType struct {
	Underlying Type
	Size       Expr
}

func (ArrayType) _type() {}
