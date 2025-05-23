package parser

import (
	"fmt"

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
		return ast.ExprStmt{Expression: expression}
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

	if p.currentTokenKind() == lexer.ASSIGNMENT {
		p.advance()
		assignedValue = parse_expr(p, assignment)
	} else if p.currentTokenKind() != lexer.SEMI_COLON {
		lastToken := p.currentToken()
		if p.pos > 0 {
			lastToken = p.tokens[p.pos-1]
		}
		p.errors.Report(fmt.Sprintf("Mowa anukundhi okati : 'semi_colon' vachindhi okati : '%s': Variable declaration ki semicolon undaali Mowa '%s' tarvatha",
			lexer.TokenKindString(p.currentTokenKind()), lastToken.Value), lastToken.Line)
		return nil
	}

	if explicitType == nil && assignedValue == nil {
		p.errors.Report("Mowa, variable ki ayithe value undali lekapothey type undaali!", p.line)
		return nil
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
			p.errors.Report("Comma lekapothey semicolon expect chesthunna Mowa!", p.line)
			break
		}
	}
	return ast.PrintStmt{Expressions: expressions}
}

func parse_input_stmt(p *parser) ast.Stmt {
	p.expect(lexer.INPUT)
	if p.currentTokenKind() == lexer.IDENTIFIER {
		varName := p.advance().Value
		if p.currentTokenKind() == lexer.OPEN_BRACKET {
			p.advance() // Consume OPEN_BRACKET
			index := parse_expr(p, default_bp)
			if index == nil {
				p.errors.Report("Mowa, array index expression undaali mowa!", p.line)
				return nil
			}
			p.expect(lexer.CLOSE_BRACKET)
			return ast.InputIndexStmt{
				Array: ast.SymbolExpr{Value: varName},
				Index: index,
			}
		}
		return ast.InputStmt{VarName: varName}
	}
	p.expectError(lexer.IDENTIFIER, "Input ki variable name undaali Mowa")
	return nil
}

func parse_if_stmt(p *parser) ast.Stmt {
	p.expect(lexer.IF)
	p.expect(lexer.OPEN_PAREN)
	condition := parse_expr(p, default_bp)
	if condition == nil {
		p.errors.Report("Mowa, if ki condition undaali mowa!", p.line)
		return nil
	}
	p.expect(lexer.CLOSE_PAREN)
	p.expect(lexer.THEN)
	p.expect(lexer.OPEN_CURLY)
	thenBody := parse_block(p)
	p.expect(lexer.CLOSE_CURLY)

	var elseIfs []ast.ElseIfBranch
	for p.hasTokens() && p.currentTokenKind() == lexer.ELSE_IF {
		p.advance()
		p.expect(lexer.OPEN_PAREN)
		elseIfCondition := parse_expr(p, default_bp)
		if elseIfCondition == nil {
			p.errors.Report("Mowa, else if ki condition undaali mowa!", p.line)
			return nil
		}
		p.expect(lexer.CLOSE_PAREN)
		p.expect(lexer.OPEN_CURLY)
		elseIfBody := parse_block(p)
		p.expect(lexer.CLOSE_CURLY)
		elseIfs = append(elseIfs, ast.ElseIfBranch{
			Condition: elseIfCondition,
			Body:      elseIfBody,
		})
	}

	var elseBranch *ast.BlockStmt
	if p.hasTokens() && p.currentTokenKind() == lexer.ELSE {
		p.advance()
		p.expect(lexer.OPEN_CURLY)
		elseBody := parse_block(p)
		p.expect(lexer.CLOSE_CURLY)
		elseBranch = &elseBody
	}

	return ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBody,
		ElseIfs:    elseIfs,
		ElseBranch: elseBranch,
	}
}

func parse_switch_stmt(p *parser) ast.Stmt {
	p.expect(lexer.SWITCH)
	p.expect(lexer.OPEN_PAREN)
	expr := parse_expr(p, default_bp)
	if expr == nil {
		p.errors.Report("Mowa, switch ki expression undaali mowa!", p.line)
		return nil
	}
	p.expect(lexer.CLOSE_PAREN)
	p.expect(lexer.OPEN_CURLY)

	cases := []ast.CaseBranch{}
	var defaultBranch *ast.BlockStmt

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		if p.currentTokenKind() == lexer.CASE {
			p.advance()
			caseValue := parse_expr(p, default_bp)
			if caseValue == nil {
				p.errors.Report("Mowa, case ki value undaali mowa!", p.line)
				return nil
			}
			p.expect(lexer.COLON)
			body := []ast.Stmt{}
			for p.hasTokens() && p.currentTokenKind() != lexer.CASE && p.currentTokenKind() != lexer.DEFAULT && p.currentTokenKind() != lexer.CLOSE_CURLY {
				stmt := parse_stmt(p)
				if stmt != nil {
					body = append(body, stmt)
				}
				if p.errors.HasErrors() {
					break
				}
				if p.hasTokens() && p.currentTokenKind() == lexer.SEMI_COLON {
					p.advance()
				}
			}
			cases = append(cases, ast.CaseBranch{
				Value: caseValue,
				Body:  ast.BlockStmt{Body: body},
			})
		} else if p.currentTokenKind() == lexer.DEFAULT {
			p.advance()
			p.expect(lexer.COLON)
			body := []ast.Stmt{}
			for p.hasTokens() && p.currentTokenKind() != lexer.CASE && p.currentTokenKind() != lexer.DEFAULT && p.currentTokenKind() != lexer.CLOSE_CURLY {
				stmt := parse_stmt(p)
				if stmt != nil {
					body = append(body, stmt)
				}
				if p.errors.HasErrors() {
					break
				}
				if p.hasTokens() && p.currentTokenKind() == lexer.SEMI_COLON {
					p.advance()
				}
			}
			if defaultBranch != nil {
				p.errors.Report("Mowa, switch lo okka default matrame undali mowa!", p.line)
				return nil
			}
			defaultBranch = &ast.BlockStmt{Body: body}
		} else {
			p.errors.Report(fmt.Sprintf("Mowa, 'case' or 'default' expect chesthunna, '%s' vachindhi mowa!", lexer.TokenKindString(p.currentTokenKind())), p.line)
			return nil
		}
	}
	p.expect(lexer.CLOSE_CURLY)

	return ast.SwitchStmt{
		Expression: expr,
		Cases:      cases,
		Default:    defaultBranch,
	}
}

func parse_block(p *parser) ast.BlockStmt {
	body := []ast.Stmt{}
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		stmt := parse_stmt(p)
		if stmt != nil {
			body = append(body, stmt)
		}
		if p.errors.HasErrors() {
			break
		}
		// Check if the statement requires a semicolon
		isBlockStmt := false
		if stmt != nil {
			switch stmt.(type) {
			case ast.IfStmt, ast.ForStmt, ast.SwitchStmt, ast.FunctionDeclStmt:
				isBlockStmt = true
			}
		}
		if !isBlockStmt && p.currentTokenKind() != lexer.CLOSE_CURLY {
			// Expect semicolon for non-block statements
			if p.currentTokenKind() != lexer.SEMI_COLON {
				lastToken := p.currentToken()
				errorLine := p.line
				if p.pos > 0 {
					lastToken = p.tokens[p.pos-1]
					errorLine = lastToken.Line
				}
				p.errors.Report(fmt.Sprintf("Mowa ravalsindhi 'semi_colon' kaani vachindhi '%s'.: Statement end ki semicolon undaali Mowa '%s' tarvatha",
					lexer.TokenKindString(p.currentTokenKind()), lastToken.Value), errorLine)
				// Skip to next semicolon or closing curly to recover
				for p.hasTokens() && p.currentTokenKind() != lexer.SEMI_COLON && p.currentTokenKind() != lexer.CLOSE_CURLY {
					p.advance()
				}
				if p.currentTokenKind() == lexer.SEMI_COLON {
					p.advance()
				}
			} else {
				// Consume semicolon
				p.advance()
			}
		}
	}
	return ast.BlockStmt{Body: body}
}

func parse_for_stmt(p *parser) ast.Stmt {
	p.expect(lexer.FOR)
	p.expect(lexer.OPEN_PAREN)

	var init ast.Stmt
	if p.currentTokenKind() == lexer.IDHI {
		init = parse_decl_stmt(p)
	} else {
		initExpr := parse_expr(p, default_bp)
		init = ast.ExprStmt{Expression: initExpr}
	}
	p.expect(lexer.SEMI_COLON)

	condition := parse_expr(p, default_bp)
	p.expect(lexer.SEMI_COLON)

	increment := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)

	p.expect(lexer.OPEN_CURLY)
	body := parse_block(p)
	p.expect(lexer.CLOSE_CURLY)

	return ast.ForStmt{
		Init:      init,
		Condition: condition,
		Increment: increment,
		Body:      body,
	}
}

func parse_function_decl_stmt(p *parser) ast.Stmt {
	p.expect(lexer.FN)
	name := p.expectError(lexer.IDENTIFIER, "Function name expect chesthunna Mowa!").Value

	// Parse parameters (e.g., (a: number, b: number))
	p.expect(lexer.OPEN_PAREN)
	parameters := []ast.Parameter{}
	if p.currentTokenKind() != lexer.CLOSE_PAREN {
		for {
			paramName := p.expectError(lexer.IDENTIFIER, "Parameter name expect chesthunna Mowa!").Value
			p.expect(lexer.COLON)
			paramType := parse_type(p, default_bp)
			if paramType == nil {
				p.errors.Report("Mowa, parameter type undaali mowa!", p.line)
				return nil
			}
			parameters = append(parameters, ast.Parameter{Name: paramName, Type: paramType})
			if p.currentTokenKind() != lexer.COMMA {
				break
			}
			p.advance()
		}
	}
	p.expect(lexer.CLOSE_PAREN)

	// Parse return type (optional)
	var returnType ast.Type
	if p.currentTokenKind() == lexer.COLON {
		p.advance()
		returnType = parse_type(p, default_bp)
		if returnType == nil {
			p.errors.Report("Mowa, return type parse cheyalenu mowa!", p.line)
			return nil
		}
	}

	// Parse function body
	p.expect(lexer.OPEN_CURLY)
	body := parse_block(p)
	p.expect(lexer.CLOSE_CURLY)

	return ast.FunctionDeclStmt{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
	}
}

func parse_return_stmt(p *parser) ast.Stmt {
	p.expect(lexer.RETURN)
	var value ast.Expr
	if p.currentTokenKind() != lexer.SEMI_COLON {
		value = parse_expr(p, default_bp)
		if value == nil {
			p.errors.Report("Mowa, return value parse cheyalenu mowa!", p.line)
			return nil
		}
	}
	return ast.ReturnStmt{Value: value}
}
