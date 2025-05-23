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
	line   int
}

//  parser
func createParser(tokens []lexer.Token) *parser {
	// Initialize token and type lookup tables
	createTokenLookups()
	createTypeTokenLookups()
	p := &parser{
		tokens: tokens,
		pos:    0,
		errors: errors.NewErrorReporter(),
		line:   1,
	}
	if len(tokens) > 0 {
		p.line = tokens[0].Line
	}
	return p
}

// Parse converts tokens into an AST
// Returns the AST and any parsing errors.
func Parse(tokens []lexer.Token) (ast.BlockStmt, *errors.ErrorReporter) {
	Body := make([]ast.Stmt, 0)
	p := createParser(tokens)

	for p.hasTokens() {
		stmt := parse_stmt(p)
		if stmt != nil {
			Body = append(Body, stmt)
		}
		if p.errors.HasErrors() {
			for p.hasTokens() && p.currentTokenKind() != lexer.SEMI_COLON && p.currentTokenKind() != lexer.EOF {
				p.advance()
			}
			if p.currentTokenKind() == lexer.SEMI_COLON {
				p.advance()
			}
			continue
		}
		if p.hasTokens() && p.currentTokenKind() != lexer.SEMI_COLON {
			// Check if the statement is a block statement (no semicolon required)
			isBlockStmt := false
			if stmt != nil {
				switch stmt.(type) {
				case ast.IfStmt, ast.ForStmt, ast.SwitchStmt, ast.FunctionDeclStmt:
					isBlockStmt = true
				}
			}
			if !isBlockStmt {
				// Report missing semicolon for non-block statements
				lastToken := p.currentToken()
				errorLine := p.line
				if p.pos > 0 {
					lastToken = p.tokens[p.pos-1]
					errorLine = lastToken.Line
				}
				p.errors.Report(fmt.Sprintf("Mowa ravalsindhi 'semi_colon' kaani vachindhi '%s'.: Statement end ki semicolon undaali Mowa '%s' tarvatha",
					lexer.TokenKindString(p.currentTokenKind()), lastToken.Value), errorLine)
				for p.hasTokens() && p.currentTokenKind() != lexer.SEMI_COLON && p.currentTokenKind() != lexer.EOF {
					p.advance()
				}
				if p.currentTokenKind() == lexer.SEMI_COLON {
					p.advance()
				}
			}
		} else if p.hasTokens() {
			// Consume semicolon if present
			p.advance()
			if p.pos < len(p.tokens) {
				p.line = p.tokens[p.pos].Line
			}
		}
	}

	return ast.BlockStmt{Body: Body}, p.errors
}

func (p *parser) currentToken() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Kind: lexer.EOF, Value: "EOF", Line: p.line}
	}
	return p.tokens[p.pos]
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++
	if p.pos < len(p.tokens) {
		p.line = p.tokens[p.pos].Line
	}
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
		p.errors.Report(msg, token.Line)
	}
	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}
