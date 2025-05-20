package ast

import "github.com/pvarma-05/MowaLang/src/lexer"

// -------------------
// LITERAL EXPRESSIONS
// -------------------

type NumberExpr struct {
	Value float64
}

func (NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (SymbolExpr) expr() {}

type BoolExpr struct {
	Value bool
}

func (BoolExpr) expr() {}

// -------------------
// COMPLEX EXPRESSIONS
// -------------------

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (PrefixExpr) expr() {}

type PostfixExpr struct {
	LeftExpr Expr
	Operator lexer.Token
}

func (PostfixExpr) expr() {}

type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Value    Expr
}

func (AssignmentExpr) expr() {}

type ArrayLiteralExpr struct {
	Elements []Expr
}

func (ArrayLiteralExpr) expr() {}

type ArrayIndexExpr struct {
	Array Expr
	Index Expr
}

func (ArrayIndexExpr) expr() {}

type MemberAccessExpr struct {
	Object   Expr
	Property string
}

func (MemberAccessExpr) expr() {}

type CallExpr struct {
	Function  Expr
	Arguments []Expr
}

func (CallExpr) expr() {}

type TypeofExpr struct {
	Operator lexer.Token
	Right    Expr
}

func (TypeofExpr) expr() {}
