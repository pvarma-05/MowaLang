package parser

import (
	"fmt"

	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

type type_nud_handler func(p *parser) ast.Type
type type_led_handler func(p *parser, left ast.Type, bp binding_power) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]binding_power

var type_bp_lu = type_bp_lookup{}
var type_nud_lu = type_nud_lookup{}
var type_led_lu = type_led_lookup{}

// func type_led(kind lexer.TokenKind, bp binding_power, led_fn type_led_handler) {
// 	type_bp_lu[kind] = bp
// 	type_led_lu[kind] = led_fn
// }

func type_nud(kind lexer.TokenKind, nud_fn type_nud_handler) {
	type_nud_lu[kind] = nud_fn
}

func createTypeTokenLookups() {
	type_nud(lexer.IDENTIFIER, parse_symbol_type)
	type_nud(lexer.OPEN_BRACKET, parse_array_type)
}

func parse_symbol_type(p *parser) ast.Type {
	return ast.SymbolType{
		Name: p.expect(lexer.IDENTIFIER).Value,
	}
}

func parse_array_type(p *parser) ast.Type {
	p.advance() // Consume OPEN_BRACKET
	underlying := parse_type(p, default_bp)
	if underlying == nil {
		p.errors.Report("Mowa, array type lo underlying type undaali mowa!", p.line)
		return nil
	}
	var size ast.Expr
	if p.currentTokenKind() == lexer.OPEN_PAREN {
		p.advance() // Consume OPEN_PAREN
		size = parse_expr(p, default_bp)
		if size == nil {
			p.errors.Report("Mowa, array size expression undaali mowa!", p.line)
			return nil
		}
		p.expect(lexer.CLOSE_PAREN)
	}
	p.expect(lexer.CLOSE_BRACKET)
	if sym, ok := underlying.(ast.SymbolType); ok {
		if sym.Name != "number" && sym.Name != "string" {
			p.errors.Report(fmt.Sprintf("Arrays ki ayithe 'number' undaali lekapothey 'string' Mowa! '%s' undakodadhu", sym.Name), p.line)
			return nil
		}
	}
	return ast.ArrayType{Underlying: underlying, Size: size}
}

func parse_type(p *parser, bp binding_power) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]
	if !exists {
		p.errors.Report(fmt.Sprintf("Type handler missing mowa for '%s'", lexer.TokenKindString(tokenKind)), p.line)
		return nil
	}
	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := type_led_lu[tokenKind]
		if !exists {
			panic(fmt.Sprintf("Token %s ki LED Handler expect chesthunna Mowa\n", lexer.TokenKindString(tokenKind)))
		}
		left = led_fn(p, left, type_bp_lu[p.currentTokenKind()])
	}
	return left
}
