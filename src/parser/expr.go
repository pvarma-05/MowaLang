package parser

import (
	"fmt"
	"strconv"

	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

func parse_expr(p *parser, bp binding_power) ast.Expr {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]
	if !exists {
		p.errors.Report(fmt.Sprintf("token '%s' ki NUD handler undaali Mowa", lexer.TokenKindString(tokenKind)), p.line)
		return nil
	}
	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]
		if !exists {
			p.errors.Report(fmt.Sprintf("token '%s' ki LED handler undaali Mowa", lexer.TokenKindString(tokenKind)), p.line)
			return left
		}
		left = led_fn(p, left, bp_lu[p.currentTokenKind()])
	}
	return left
}

func parse_primary_expr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER, lexer.STRING, lexer.IDENTIFIER, lexer.TRUE, lexer.FALSE:
		return p.parse_primary_expr_helper()
	default:
		p.errors.Report(fmt.Sprintf("'%s' vaadi expression cheyyalem Mowa", lexer.TokenKindString(p.currentTokenKind())), p.line)
		return nil
	}
}

func parse_binary_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operatorToken := p.advance()
	right := parse_expr(p, bp)
	return ast.BinaryExpr{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}

func parse_prefix_expr(p *parser) ast.Expr {
	op := p.advance()
	rhs := parse_expr(p, unary)
	return ast.PrefixExpr{
		Operator:  op,
		RightExpr: rhs,
	}
}

func parse_assignment_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	op := p.advance()
	rhs := parse_expr(p, bp)
	return ast.AssignmentExpr{
		Operator: op,
		Value:    rhs,
		Assignee: left,
	}
}

func parse_grouping_expr(p *parser) ast.Expr {
	p.advance()
	expr := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parse_array_literal_expr(p *parser) ast.Expr {
	p.advance()
	elements := []ast.Expr{}
	if p.currentTokenKind() != lexer.CLOSE_BRACKET {
		elements = append(elements, parse_expr(p, default_bp))
		for p.currentTokenKind() == lexer.COMMA {
			p.advance()
			if p.currentTokenKind() == lexer.CLOSE_BRACKET {
				break
			}
			elements = append(elements, parse_expr(p, default_bp))
		}
	}
	p.expect(lexer.CLOSE_BRACKET)
	return ast.ArrayLiteralExpr{Elements: elements}
}

func parse_array_index_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance()
	index := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_BRACKET)
	return ast.ArrayIndexExpr{Array: left, Index: index}
}

func parse_member_access_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance()
	if p.currentTokenKind() != lexer.IDENTIFIER {
		p.errors.Report(fmt.Sprintf("Expected identifier after '.', got '%s'", lexer.TokenKindString(p.currentTokenKind())), p.line)
		return nil
	}
	property := p.advance().Value
	return ast.MemberAccessExpr{Object: left, Property: property}
}

func parse_call_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.advance() // Consume OPEN_PAREN
	arguments := []ast.Expr{}
	if p.currentTokenKind() != lexer.CLOSE_PAREN {
		for {
			arg := parse_expr(p, default_bp)
			if arg == nil {
				p.errors.Report("Mowa, argument expression parse cheyalenu mowa!", p.line)
				return nil
			}
			arguments = append(arguments, arg)
			if p.currentTokenKind() != lexer.COMMA {
				break
			}
			p.advance() // Consume comma
		}
	}
	p.expect(lexer.CLOSE_PAREN)
	return ast.CallExpr{Function: left, Arguments: arguments}
}

func (p *parser) parse_primary_expr_helper() ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumberExpr{Value: number}
	case lexer.STRING:
		return ast.StringExpr{Value: p.advance().Value}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{Value: p.advance().Value}
	case lexer.TRUE:
		p.advance()
		return ast.BoolExpr{Value: true}
	case lexer.FALSE:
		p.advance()
		return ast.BoolExpr{Value: false}
	}
	return nil
}
