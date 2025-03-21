package parser

import (
	"fmt"

	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/errors"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
	errors *errors.ErrorReporter
}

func createParser(tokens []lexer.Token) *parser {
	createTokenLookups()
	createTypeTokenLookups()
	return &parser{
		tokens: tokens,
		pos:    0,
		errors: errors.NewErrorReporter(),
	}
}

func Parse(tokens []lexer.Token) (ast.BlockStmt, *errors.ErrorReporter) {
	Body := make([]ast.Stmt, 0)
	p := createParser(tokens)

	for p.hasTokens() {
		stmt := parse_stmt(p)
		if stmt != nil {
			Body = append(Body, stmt)
		}
		if p.errors.HasErrors() {
			// Skip to next semicolon or EOF to recover, then stop
			for p.hasTokens() && p.currentTokenKind() != lexer.SEMI_COLON && p.currentTokenKind() != lexer.EOF {
				p.advance()
			}
			if p.currentTokenKind() == lexer.SEMI_COLON {
				p.advance()
			}
			break
		}
		// Expect a semicolon after each statement; if missing, report and recover
		if p.hasTokens() && p.currentTokenKind() != lexer.SEMI_COLON {
			p.expectError(lexer.SEMI_COLON, "Statement end ki semicolon undaali Mowa")
			for p.hasTokens() && p.currentTokenKind() != lexer.SEMI_COLON && p.currentTokenKind() != lexer.EOF {
				p.advance()
			}
			if p.currentTokenKind() == lexer.SEMI_COLON {
				p.advance()
			}
		} else if p.hasTokens() {
			p.advance() // Consume the semicolon
		}
	}

	return ast.BlockStmt{Body: Body}, p.errors
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++
	return tk
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind
	if kind != expectedKind {
		msg := fmt.Sprintf("Mowa anukundhi okati : '%s' vachindhi okati : '%s'", lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind))
		if err != nil {
			msg = fmt.Sprintf("%s: %v", msg, err)
		}
		p.errors.Report(msg)
	}
	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}
