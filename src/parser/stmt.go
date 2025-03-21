package parser

import (
	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

func parse_stmt(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]
	if exists {
		return stmt_fn(p)
	}
	expression := parse_expr(p, default_bp)
	if expression != nil {
		return ast.ExprStmt{Expression: expression} // Semicolon handled in Parse
	}
	return nil
}

func parse_decl_stmt(p *parser) ast.Stmt {
	var explicitType ast.Type
	var assignedValue ast.Expr
	p.expect(lexer.IDHI)
	VarName := p.expectError(lexer.IDENTIFIER, "variable name expect chesthunna Mowa!").Value

	if p.currentTokenKind() == lexer.COLON {
		p.advance()
		explicitType = parse_type(p, default_bp)
	}

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assignedValue = parse_expr(p, assignment)
	} else if explicitType == nil {
		p.errors.Report("Mowa, variable ki ayithe value undali lekapothey type undaali!")
	}

	return ast.VarDeclStmt{
		ExplicitType:  explicitType,
		VarName:       VarName,
		AssignedValue: assignedValue,
	}
}

func parse_print_stmt(p *parser) ast.Stmt {
	p.expect(lexer.PRINT)
	expressions := []ast.Expr{}
	for p.currentTokenKind() != lexer.SEMI_COLON && p.hasTokens() {
		expr := parse_expr(p, default_bp)
		if expr != nil {
			expressions = append(expressions, expr)
		}
		if p.currentTokenKind() == lexer.COMMA {
			p.advance()
		} else if p.currentTokenKind() != lexer.SEMI_COLON {
			p.errors.Report("Comma lekapothey semicolon expect chesthunna Mowa!")
			break // Let Parse handle recovery
		}
	}
	return ast.PrintStmt{Expressions: expressions}
}

func parse_input_stmt(p *parser) ast.Stmt {
	p.expect(lexer.INPUT)
	if p.currentTokenKind() != lexer.IDENTIFIER {
		p.expectError(lexer.IDENTIFIER, "Input ki variable name undaali Mowa")
		return nil
	}
	varName := p.advance().Value
	return ast.InputStmt{
		VarName: varName,
	}
}
