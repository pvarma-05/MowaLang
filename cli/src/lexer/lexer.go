package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pvarma-05/MowaLang/src/errors"
)

// regexHandler defines a function to handle a matched regex pattern.
type regexHandler func(lex *lexer, regex *regexp.Regexp)

// regexPattern pairs a regex with its handler function.
type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

// lexer holds the state of the lexing process.
type lexer struct {
	pos      int                   // Current position in the source code
	source   string                // Input source code
	Tokens   []Token               // List of generated tokens
	patterns []regexPattern        // Regex patterns and their handlers
	errors   *errors.ErrorReporter // Error reporter for lexing errors
	line     int                   // Current line number for error reporting
}

// advanceN moves the lexer position forward by n characters and updates line count.
func (lex *lexer) advanceN(n int) {
	// Count newlines in the advanced section to track line numbers.
	slice := lex.source[lex.pos : lex.pos+n]
	newlines := strings.Count(slice, "\n")
	lex.line += newlines
	lex.pos += n
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}

// Returns the tokens and any lexing errors.
func Tokenize(source string) ([]Token, *errors.ErrorReporter) {
	lex := createLexer(source)

	for !lex.at_eof() {
		matched := false
		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}
		if !matched {
			lex.errors.Report(fmt.Sprintf("lexer problem Mowa, unrecognized token near '%s'", lex.remainder()), lex.line)
			break
		}
	}
	lex.push(newToken(EOF, "EOF", lex.line))
	return lex.Tokens, lex.errors
}

// defaultHandler creates a handler for simple tokens (e.g., operators, punctuation).
func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(newToken(kind, value, lex.line))
	}
}

// createLexer initializes a lexer with the source code and regex patterns.
func createLexer(source string) *lexer {
	return &lexer{
		pos:    0,
		source: source,
		Tokens: make([]Token, 0),
		errors: errors.NewErrorReporter(),
		line:   1, // Start at line 1
		patterns: []regexPattern{
			// Identifiers and keywords (e.g., "a", "idhi")
			{regexp.MustCompile(`[a-zA-Z_][A-Za-z0-9_]*`), symbolHandler},
			// Numbers (e.g., "123", "3.14")
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			// Strings (e.g., "\"hello\"")
			{regexp.MustCompile(`"[^"]*"|'[^']*'`), stringHandler},
			// Comments (e.g., "// comment")
			{regexp.MustCompile(`\/\/.*`), skipHandler},
			// Whitespace (skipped)
			{regexp.MustCompile(`\s+`), skipHandler},
			// Punctuation and operators
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\*\*`), defaultHandler(STAR_STAR, "**")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\*=`), defaultHandler(STAR_EQUALS, "*=")},
			{regexp.MustCompile(`\/=`), defaultHandler(SLASH_EQUALS, "/=")},
			{regexp.MustCompile(`%=`), defaultHandler(PERCENT_EQUALS, "%=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

// skipHandler skips whitespace and comments without generating tokens.
func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

// stringHandler processes string literals (e.g., "\"hello\"").
func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]
// Normalize to double quotes
	normalized := stringLiteral
	if stringLiteral[0] == '\'' {
		normalized = `"` + stringLiteral[1:len(stringLiteral)-1] + `"`
	}
	lex.push(newToken(STRING, normalized, lex.line))
	lex.advanceN(len(stringLiteral))
}

// numberHandler processes numeric literals (e.g., "123", "3.14").
func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(newToken(NUMBER, match, lex.line))
	lex.advanceN(len(match))
}

// symbolHandler processes identifiers and keywords (e.g., "a", "idhi").
func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	remainder := lex.remainder()

	// Check for the longest matching keyword (e.g., "idhi" over "i").
	longestMatch := ""
	matchingToken := IDENTIFIER

	for key, kind := range reserved_lu {
		if len(remainder) >= len(key) && remainder[:len(key)] == key {
			if len(key) > len(longestMatch) {
				longestMatch = key
				matchingToken = kind
			}
		}
	}

	if longestMatch != "" {
		// Keyword found (e.g., "idhi" → IDHI).
		lex.advanceN(len(longestMatch))
		lex.push(newToken(matchingToken, longestMatch, lex.line))
	} else {
		// Identifier (e.g., "a").
		match := regex.FindString(lex.remainder())
		lex.advanceN(len(match))
		lex.push(newToken(IDENTIFIER, match, lex.line))
	}
}
