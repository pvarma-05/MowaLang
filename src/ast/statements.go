package ast

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

type ExprStmt struct {
	Expression Expr
}

func (n ExprStmt) stmt() {}

type VarDeclStmt struct {
	VarName       string
	AssignedValue Expr
	ExplicitType  Type
}

func (n VarDeclStmt) stmt() {}

type PrintStmt struct {
	Expressions []Expr
}

func (n PrintStmt) stmt() {}

type InputStmt struct {
	VarName string
}

func (n InputStmt) stmt() {}
