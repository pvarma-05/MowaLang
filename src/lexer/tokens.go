package lexer

import "fmt"

type TokenKind int

const (
	EOF TokenKind = iota
	NUMBER
	STRING
	IDENTIFIER

	OPEN_BRACKET  // [
	CLOSE_BRACKET // ]
	OPEN_CURLY    // {
	CLOSE_CURLY   // }
	OPEN_PAREN    // (
	CLOSE_PAREN   // )

	ASSIGNMENT // =
	EQUALS     // ==
	NOT        // !
	NOT_EQUALS // !=

	LESS           // <
	LESS_EQUALS    // <=
	GREATER        // >
	GREATER_EQUALS // >=

	OR  // ||
	AND // &&

	DOT        // .
	DOT_DOT    // ..
	SEMI_COLON // ;
	COLON      // :
	QUESTION   // ?
	COMMA      // ,

	PLUS_PLUS      // ++
	MINUS_MINUS    // --
	PLUS_EQUALS    // +=
	MINUS_EQUALS   // -=
	STAR_EQUALS    // *=
	SLASH_EQUALS   // /=
	PERCENT_EQUALS // %=

	PLUS    // +
	DASH    // -
	STAR    // *
	SLASH   // /
	PERCENT // %

	MOWA_START // mowa start
	MOWA_END   // mowa end

	IDHI  // Variable declaration
	ANUKO // in loop variable declaration

	IF      // mowa okavela (cond) ayithe
	THEN    // ayithe
	ELSE_IF // ledhante (cond)
	ELSE    // ledha

	SWITCH  // enchuko
	CASE    // case
	DEFAULT // default

	FN     // pani
	RETURN // ichey

	PRINT // mowa cheppu
	INPUT // theesko

	TRUE  // nijam
	FALSE // abadham
	NULL  // bokka

	BREAK    // aagipo
	CONTINUE // konasaagu
	TYPEOF   // rakam

	WHILE    // mowa ayye varaku
	DO_WHILE // mowa modhalu { } ayye varaku
	FOR      // mowa sarlu (anuko variable ; cond ; inc/dec)
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"mowa start":       MOWA_START,
	"mowa end":         MOWA_END,
	"mowa idhi":        IDHI,
	"anuko":            ANUKO,
	"mowa okavela":     IF,
	"ayithe":           THEN,
	"ledhante":         ELSE_IF,
	"ledha":            ELSE,
	"enchuko":          SWITCH,
	"case":             CASE,
	"default":          DEFAULT,
	"pani":             FN,
	"ichey":            RETURN,
	"mowa cheppu":      PRINT,
	"mowa theesko":     INPUT,
	"nijam":            TRUE,
	"abadham":          FALSE,
	"bokka":            NULL,
	"aagipo":           BREAK,
	"konasaagu":        CONTINUE,
	"rakam":            TYPEOF,
	"mowa ayye varaku": WHILE,
	"mowa modhalu":     DO_WHILE,
	"mowa sarlu":       FOR,
}

type Token struct {
	Kind  TokenKind
	Value string
}

func (token Token) isOneOfMany(expectedTokens ...TokenKind) bool {
	for _, expected := range expectedTokens {
		if expected == token.Kind {
			return true
		}
	}
	return false
}

func (token Token) Debug() {
	if token.isOneOfMany(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", TokenKindString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s ()\n", TokenKindString(token.Kind))
	}
}

func newToken(kind TokenKind, value string) Token {
	return Token{
		kind, value,
	}
}

func TokenKindString(kind TokenKind) string {
	switch kind {
	case EOF:
		return "eof"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case IDENTIFIER:
		return "identifier"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case EQUALS:
		return "equals"
	case NOT_EQUALS:
		return "not_equals"
	case NOT:
		return "not"
	case LESS:
		return "less"
	case LESS_EQUALS:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUALS:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case DOT_DOT:
		return "dot_dot"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS_MINUS:
		return "minus_minus"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
	case PLUS:
		return "plus"
	case DASH:
		return "dash"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case PERCENT:
		return "percent"
	case MOWA_START:
		return "mowa_start"
	case MOWA_END:
		return "mowa_end"
	case PRINT:
		return "cheppu"
	case IDHI:
		return "idhi"
	case ANUKO:
		return "anuko"
	case FN:
		return "pani"
	case RETURN:
		return "ichey"
	case SWITCH:
		return "enchuko"
	case CASE:
		return "case"
	case DEFAULT:
		return "default"
	case IF:
		return "okavela"
	case THEN:
		return "ayithe"
	case ELSE_IF:
		return "ledhantey"
	case ELSE:
		return "ledha"
	case TRUE:
		return "nijam"
	case FALSE:
		return "abadham"
	case NULL:
		return "bokka"
	case BREAK:
		return "aagipo"
	case CONTINUE:
		return "konasaagu"
	case TYPEOF:
		return "rakam"
	case WHILE:
		return "mowa ayye varaku"
	case DO_WHILE:
		return "mowa modhalu"
	case FOR:
		return "mowa sarlu"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
