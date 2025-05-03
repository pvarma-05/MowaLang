package parser

import (
	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

type binding_power int

const (
	default_bp     binding_power = iota // Lowest precedence
	comma                               // ,
	assignment                          // =, +=, -=, etc.
	logical                             // &&, ||
	relational                          // ==, !=, <, <=, >, >=
	additive                            // +, -
	multiplicative                      // *, /, %, **
	unary                               // ++, --
	call                                // Function calls
	member                              // ., []
	primary                             // Literals, identifiers
)

// Handler types for parsing statements and expressions.
type stmt_handler func(p *parser) ast.Stmt
type nud_handler func(p *parser) ast.Expr
type led_handler func(p *parser, left ast.Expr, bp binding_power) ast.Expr

// Lookup tables for parsing rules.
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

// createTokenLookups initializes the parsing rules for statements and expressions.
func createTokenLookups() {
	// Assignment operators (e.g., a = 5, a += 2)
	led(lexer.ASSIGNMENT, assignment, parse_assignment_expr)
	led(lexer.PLUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.MINUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.STAR_EQUALS, assignment, parse_assignment_expr)
	led(lexer.SLASH_EQUALS, assignment, parse_assignment_expr)
	led(lexer.PERCENT_EQUALS, assignment, parse_assignment_expr)

	// Postfix operators (e.g., a++, a--)
	led(lexer.PLUS_PLUS, unary, parse_postfix_expr)
	led(lexer.MINUS_MINUS, unary, parse_postfix_expr)

	// Prefix operators (e.g., ++a, --a, rakam)
	nud(lexer.PLUS_PLUS, parse_prefix_expr)
	nud(lexer.MINUS_MINUS, parse_prefix_expr)
	nud(lexer.TYPEOF, parse_typeof_expr)

	// Array operations (e.g., [1, 2], arr[0], arr.length)
	nud(lexer.OPEN_BRACKET, parse_array_literal_expr)
	led(lexer.OPEN_BRACKET, member, parse_array_index_expr)
	led(lexer.DOT, member, parse_member_access_expr)

	// Function calls (e.g., add(1, 2))
	led(lexer.OPEN_PAREN, call, parse_call_expr)

	// Literals and grouping (e.g., 123, "hello", a, (a + b))
	nud(lexer.NUMBER, parse_primary_expr)
	nud(lexer.STRING, parse_primary_expr)
	nud(lexer.IDENTIFIER, parse_primary_expr)
	nud(lexer.DASH, parse_prefix_expr)
	nud(lexer.OPEN_PAREN, parse_grouping_expr)
	nud(lexer.TRUE, parse_primary_expr)
	nud(lexer.FALSE, parse_primary_expr)

	// Binary operators
	led(lexer.AND, logical, parse_binary_expr)
	led(lexer.OR, logical, parse_binary_expr)
	led(lexer.DOT_DOT, logical, parse_binary_expr) // Not used
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
	stmt(lexer.BREAK, parse_break_stmt)
	stmt(lexer.FOR, parse_for_stmt)
	stmt(lexer.CONTINUE, parse_continue_stmt)
	stmt(lexer.FN, parse_function_decl_stmt)
	stmt(lexer.RETURN, parse_return_stmt)
}

func parse_postfix_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	return ast.PostfixExpr{
		LeftExpr: left,
		Operator: operatorToken,
	}
}

func parse_break_stmt(p *parser) ast.Stmt {
	p.expect(lexer.BREAK)
	return ast.BreakStmt{}
}

func parse_continue_stmt(p *parser) ast.Stmt {
	p.expect(lexer.CONTINUE)
	return ast.ContinueStmt{}
}

func parse_typeof_expr(p *parser) ast.Expr {
	op := p.advance() // Consume TYPEOF (rakam)
	p.expect(lexer.OPEN_PAREN)
	expr := parse_expr(p, default_bp)
	if expr == nil {
		p.errors.Report("Mowa, rakam argument parse cheyalenu ra!", p.line)
		return nil
	}
	// Restrict to SymbolExpr or ArrayIndexExpr
	if _, ok := expr.(ast.SymbolExpr); !ok {
		if _, ok := expr.(ast.ArrayIndexExpr); !ok {
			p.errors.Report("Mowa, rakam function variable or array index thiskuntadhi ra!", p.line)
			return nil
		}
	}
	p.expect(lexer.CLOSE_PAREN)
	return ast.TypeofExpr{
		Operator: op,
		Right:    expr,
	}
}
