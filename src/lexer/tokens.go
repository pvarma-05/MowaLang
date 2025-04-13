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

	PLUS_PLUS   // ++
	MINUS_MINUS // --
	STAR_STAR
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

	IDHI  // Variable declaration
	ANUKO // in loop variable declaration

	IF      // okavela (cond) ayithe
	THEN    // ayithe
	ELSE_IF // ledhante (cond)
	ELSE    // ledha

	SWITCH  // enchuko
	CASE    // case
	DEFAULT // default

	FN     // pani
	RETURN // ichey

	PRINT // mowa
	INPUT // theesko

	TRUE  // nijam
	FALSE // abadham
	NULL  // bokka

	BREAK    // aagipo
	CONTINUE // kaani
	TYPEOF   // rakam

	FOR // varaku (anuko variable ; cond ; inc/dec)
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"idhi":     IDHI,
	"anuko":    ANUKO,
	"okavela":  IF,
	"ayithe":   THEN,
	"ledhante": ELSE_IF,
	"ledha":    ELSE,
	"enchuko":  SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"pani":     FN,
	"ichey":    RETURN,
	"mowa":     PRINT,
	"theesko":  INPUT,
	"nijam":    TRUE,
	"abadham":  FALSE,
	"bokka":    NULL,
	"aagipo":   BREAK,
	"kaani":    CONTINUE,
	"rakam":    TYPEOF,
	"varaku":   FOR,
}

type Token struct {
	Kind  TokenKind
	Value string
	Line  int // Added for line number tracking
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
		fmt.Printf("%s (%s) [line %d]\n", TokenKindString(token.Kind), token.Value, token.Line) // Updated to show line
	} else {
		fmt.Printf("%s () [line %d]\n", TokenKindString(token.Kind), token.Line)
	}
}

func newToken(kind TokenKind, value string, line int) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Line:  line,
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
	case STAR_STAR:
		return "star_star"
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
	case PRINT:
		return "mowa"
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
		return "kaani"
	case TYPEOF:
		return "rakam"
	case FOR:
		return "varaku"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
}
