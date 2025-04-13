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

type IfStmt struct {
	Condition  Expr           // Condition to evaluate (e.g., a > b)
	ThenBranch BlockStmt      // Statements if condition is true
	ElseIfs    []ElseIfBranch // Optional else-if branches
	ElseBranch *BlockStmt     // Optional else branch (nil if absent)
}

type ElseIfBranch struct {
	Condition Expr
	Body      BlockStmt
}

func (n IfStmt) stmt() {}

type SwitchStmt struct {
	Expression Expr         // The value to switch on (e.g., a)
	Cases      []CaseBranch // List of case branches
	Default    *BlockStmt   // Optional default branch (nil if absent)
}

type CaseBranch struct {
	Value Expr      // The value to match (e.g., 1, "hello")
	Body  BlockStmt // Statements to execute if matched
}

func (n SwitchStmt) stmt() {}

type BreakStmt struct{}

func (n BreakStmt) stmt() {}

type ForStmt struct {
	Init      Stmt      // e.g., VarDeclStmt or AssignmentExpr
	Condition Expr      // e.g., BinaryExpr (i < 5)
	Increment Expr      // e.g., PostfixExpr (i++)
	Body      BlockStmt // Loop body
}

func (n ForStmt) stmt() {}

type ContinueStmt struct{}

func (n ContinueStmt) stmt() {}
