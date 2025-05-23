package lexer

import "fmt"

type TokenKind int

const (
	EOF            TokenKind = iota // End of file
	NUMBER                          // Numeric literal (e.g., "123", "3.14")
	STRING                          // String literal (e.g., "\"hello\"")
	IDENTIFIER                      // Variable or type name (e.g., "a", "number")
	OPEN_BRACKET                    // [ (array start)
	CLOSE_BRACKET                   // ] (array end)
	OPEN_CURLY                      // { (block start)
	CLOSE_CURLY                     // } (block end)
	OPEN_PAREN                      // ( (grouping or condition)
	CLOSE_PAREN                     // ) (grouping or condition end)
	ASSIGNMENT                      // = (assignment)
	EQUALS                          // == (equality)
	NOT                             // ! (logical not)
	NOT_EQUALS                      // != (inequality)
	LESS                            // < (less than)
	LESS_EQUALS                     // <= (less than or equal)
	GREATER                         // > (greater than)
	GREATER_EQUALS                  // >= (greater than or equal)
	OR                              // || (logical or)
	AND                             // && (logical and)
	DOT                             // . (member access, e.g., arr.length)
	DOT_DOT                         // .. (range, not currently used)
	SEMI_COLON                      // ; (statement separator)
	COLON                           // : (type annotation)
	QUESTION                        // ? (not currently used)
	COMMA                           // , (expression separator)
	PLUS_PLUS                       // ++ (increment)
	MINUS_MINUS                     // -- (decrement)
	STAR_STAR                       // ** (exponentiation)
	PLUS_EQUALS                     // += (add and assign)
	MINUS_EQUALS                    // -= (subtract and assign)
	STAR_EQUALS                     // *= (multiply and assign)
	SLASH_EQUALS                    // /= (divide and assign)
	PERCENT_EQUALS                  // %= (modulo and assign)
	PLUS                            // + (addition or concatenation)
	DASH                            // - (subtraction)
	STAR                            // * (multiplication)
	SLASH                           // / (division)
	PERCENT                         // % (modulo)
	IDHI                            // idhi (variable declaration)
	ANUKO                           // anuko (loop variable declaration)
	IF                              // okavela (if)
	THEN                            // ayithe (then)
	ELSE_IF                         // ledhante (else if)
	ELSE                            // ledha (else)
	SWITCH                          // enchuko (switch)
	CASE                            // case (case)
	DEFAULT                         // default (default case)
	FN                              // pani (function, not implemented)
	RETURN                          // ichey (return, not implemented)
	PRINT                           // mowa (print)
	INPUT                           // theesko (input)
	TRUE                            // nijam (true)
	FALSE                           // abadham (false)
	BREAK                           // aagipo (break)
	CONTINUE                        // kaani (continue)
	TYPEOF                          // rakam (typeof, not implemented)
	FOR                             // varaku (for loop)
	// NULL                            bokka (null, not implemented)
)

var reserved_lu = map[string]TokenKind{
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
	"aagipo": BREAK,
	"kaani":  CONTINUE,
	"rakam":  TYPEOF,
	"varaku": FOR,
}

type Token struct {
	Kind  TokenKind
	Value string
	Line  int
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
		fmt.Printf("%s (%s) [line %d]\n", TokenKindString(token.Kind), token.Value, token.Line)
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
