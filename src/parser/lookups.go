package parser

import (
	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

type binding_power int

const (
	default_bp binding_power = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type stmt_handler func(p *parser) ast.Stmt
type nud_handler func(p *parser) ast.Expr
type led_handler func(p *parser, left ast.Expr, bp binding_power) ast.Expr

type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]binding_power

var bp_lu = bp_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var stmt_lu = stmt_lookup{}

func led(kind lexer.TokenKind, bp binding_power, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nud_handler) {
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	stmt_lu[kind] = stmt_fn
	bp_lu[kind] = default_bp
}

func createTokenLookups() {
	// Assignment operators
	led(lexer.ASSIGNMENT, assignment, parse_assignment_expr)
	led(lexer.PLUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.MINUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.STAR_EQUALS, assignment, parse_assignment_expr)
	led(lexer.SLASH_EQUALS, assignment, parse_assignment_expr)
	led(lexer.PERCENT_EQUALS, assignment, parse_assignment_expr)

	// Postfix operators
	led(lexer.PLUS_PLUS, unary, parse_postfix_expr)   // a++
	led(lexer.MINUS_MINUS, unary, parse_postfix_expr) // a--

	// Prefix operators
	nud(lexer.PLUS_PLUS, parse_prefix_expr)   // ++a
	nud(lexer.MINUS_MINUS, parse_prefix_expr) // --a

	// Numbers, Symbols, Booleans, Grouping
	nud(lexer.NUMBER, parse_primary_expr)
	nud(lexer.STRING, parse_primary_expr)
	nud(lexer.IDENTIFIER, parse_primary_expr)
	nud(lexer.DASH, parse_prefix_expr)
	nud(lexer.OPEN_PAREN, parse_grouping_expr)
	nud(lexer.TRUE, parse_primary_expr)
	nud(lexer.FALSE, parse_primary_expr)

	// Binary Expr
	led(lexer.AND, logical, parse_binary_expr)
	led(lexer.OR, logical, parse_binary_expr)
	led(lexer.DOT_DOT, logical, parse_binary_expr)
	led(lexer.LESS, relational, parse_binary_expr)
	led(lexer.LESS_EQUALS, relational, parse_binary_expr)
	led(lexer.GREATER, relational, parse_binary_expr)
	led(lexer.GREATER_EQUALS, relational, parse_binary_expr)
	led(lexer.EQUALS, relational, parse_binary_expr)
	led(lexer.NOT_EQUALS, relational, parse_binary_expr)
	led(lexer.PLUS, additive, parse_binary_expr)
	led(lexer.DASH, additive, parse_binary_expr)
	led(lexer.STAR, multiplicative, parse_binary_expr)
	led(lexer.SLASH, multiplicative, parse_binary_expr)
	led(lexer.PERCENT, multiplicative, parse_binary_expr)
	led(lexer.STAR_STAR, multiplicative, parse_binary_expr)

	// Statements
	stmt(lexer.IDHI, parse_decl_stmt)
	stmt(lexer.PRINT, parse_print_stmt)
	stmt(lexer.INPUT, parse_input_stmt)
	stmt(lexer.IF, parse_if_stmt)
	stmt(lexer.SWITCH, parse_switch_stmt)
	stmt(lexer.BREAK, parse_break_stmt) // Added BREAK statement handler
	stmt(lexer.FOR, parse_for_stmt)
	stmt(lexer.CONTINUE, parse_continue_stmt)
}

func parse_postfix_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	return ast.PostfixExpr{
		LeftExpr: left, // Operand is on the left for postfix
		Operator: operatorToken,
	}
}

func parse_break_stmt(p *parser) ast.Stmt {
	p.expect(lexer.BREAK) // Consume 'aagipo'
	return ast.BreakStmt{}
}

func parse_continue_stmt(p *parser) ast.Stmt {
	p.expect(lexer.CONTINUE) // Consume 'kaani'
	return ast.ContinueStmt{}
}
