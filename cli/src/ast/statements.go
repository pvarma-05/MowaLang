package ast

type BlockStmt struct {
	Body []Stmt
}

func (BlockStmt) stmt() {}

type ExprStmt struct {
	Expression Expr
}

func (ExprStmt) stmt() {}

type VarDeclStmt struct {
	VarName       string
	AssignedValue Expr
	ExplicitType  Type
}

func (VarDeclStmt) stmt() {}

type PrintStmt struct {
	Expressions []Expr
}

func (PrintStmt) stmt() {}

type InputStmt struct {
	VarName string
}

func (InputStmt) stmt() {}

type InputIndexStmt struct {
	Array Expr
	Index Expr
}

func (InputIndexStmt) stmt() {}

type IfStmt struct {
	Condition  Expr
	ThenBranch BlockStmt
	ElseIfs    []ElseIfBranch
	ElseBranch *BlockStmt
}

type ElseIfBranch struct {
	Condition Expr
	Body      BlockStmt
}

func (IfStmt) stmt() {}

type SwitchStmt struct {
	Expression Expr
	Cases      []CaseBranch
	Default    *BlockStmt
}

type CaseBranch struct {
	Value Expr
	Body  BlockStmt
}

func (SwitchStmt) stmt() {}

type BreakStmt struct{}

func (BreakStmt) stmt() {}

type ContinueStmt struct{}

func (ContinueStmt) stmt() {}

type ForStmt struct {
	Init      Stmt
	Condition Expr
	Increment Expr
	Body      BlockStmt
}

func (ForStmt) stmt() {}

type FunctionDeclStmt struct {
	Name       string
	Parameters []Parameter
	ReturnType Type
	Body       BlockStmt
}

type Parameter struct {
	Name string
	Type Type
}

func (FunctionDeclStmt) stmt() {}

type ReturnStmt struct {
	Value Expr
}

func (ReturnStmt) stmt() {}
