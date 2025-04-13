package eval

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pvarma-05/MowaLang/src/ast"
	"github.com/pvarma-05/MowaLang/src/errors"
	"github.com/pvarma-05/MowaLang/src/lexer"
)

// Value represents a runtime value in MowaLang
type Value interface{}

// Environment holds variables and their values
type Environment struct {
	variables map[string]Value
	declared  map[string]bool
	types     map[string]string
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]Value),
		declared:  make(map[string]bool),
		types:     make(map[string]string),
	}
}

// Declare marks a variable as declared with a type
func (env *Environment) Declare(name string, typ ast.Type) {
	env.declared[name] = true
	if symType, ok := typ.(ast.SymbolType); ok {
		env.types[name] = symType.Name
	}
}

// Define sets a variable in the environment with type checking
func (env *Environment) Define(name string, value Value, lineNumber int, errors *errors.ErrorReporter) {
	if typ, exists := env.types[name]; exists {
		switch typ {
		case "number":
			if _, ok := value.(float64); !ok {
				errors.Report(fmt.Sprintf("Mowa, '%s' number type ki '%v' (string) ivvakudadhu ra!", name, value), lineNumber)
				return
			}
		case "string":
			if _, ok := value.(string); !ok {
				errors.Report(fmt.Sprintf("Mowa, '%s' string type ki '%v' (number) ivvakudadhu ra!", name, value), lineNumber)
				return
			}
		}
	}
	env.variables[name] = value
	env.declared[name] = true
}

// Get retrieves a variable’s value, returning an error if undefined
func (env *Environment) Get(name string) (Value, error) {
	if value, exists := env.variables[name]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("mowa, '%s' ane variable define cheyaledhu ra", name)
}

// IsDeclared checks if a variable was declared
func (env *Environment) IsDeclared(name string) bool {
	return env.declared[name]
}

// Evaluator executes the AST
type Evaluator struct {
	env    *Environment
	errors *errors.ErrorReporter
	line   int
}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{
		env:    NewEnvironment(),
		errors: errors.NewErrorReporter(),
		line:   1,
	}
}

// Evaluate runs the program
func (e *Evaluator) Evaluate(program ast.BlockStmt) *errors.ErrorReporter {
	for _, stmt := range program.Body {
		if stmt != nil {
			e.evalStmt(stmt)
			e.line++
		}
		if e.errors.HasErrors() {
			break
		}
	}
	return e.errors
}

func (e *Evaluator) evalStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case ast.VarDeclStmt:
		e.evalVarDeclStmt(s)
	case ast.InputStmt:
		e.evalInputStmt(s)
	case ast.PrintStmt:
		e.evalPrintStmt(s)
	case ast.ExprStmt:
		e.evalExpr(s.Expression)
	case ast.IfStmt:
		e.evalIfStmt(s)
	case ast.SwitchStmt:
		e.evalSwitchStmt(s)
	case ast.BreakStmt:
		// Handled in context
	case ast.ContinueStmt:
		// Handled in loop context
	case ast.ForStmt:
		e.evalForStmt(s)
	default:
		e.errors.Report(fmt.Sprintf("Mowa, ee statement type (%T) handle cheyalenu mowa!", s), e.line)
	}
}

func (e *Evaluator) evalVarDeclStmt(stmt ast.VarDeclStmt) {
	var value Value
	if stmt.AssignedValue != nil {
		value = e.evalExpr(stmt.AssignedValue)
	}
	if e.errors.HasErrors() {
		return
	}
	e.env.Declare(stmt.VarName, stmt.ExplicitType)
	if value != nil {
		e.env.Define(stmt.VarName, value, e.line, e.errors)
	}
}

func (e *Evaluator) evalInputStmt(stmt ast.InputStmt) {
	if !e.env.IsDeclared(stmt.VarName) {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' declare chesaka theesko vaadu mowa!", stmt.VarName), e.line)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if num, err := strconv.ParseFloat(input, 64); err == nil {
		e.env.Define(stmt.VarName, num, e.line, e.errors)
	} else {
		e.env.Define(stmt.VarName, input, e.line, e.errors)
	}
}

func (e *Evaluator) evalPrintStmt(stmt ast.PrintStmt) {
	var output strings.Builder
	for _, expr := range stmt.Expressions {
		value := e.evalExpr(expr)
		if e.errors.HasErrors() {
			return
		}
		switch v := value.(type) {
		case string:
			unescaped := strings.ReplaceAll(v, `\n`, "\n")
			unescaped = strings.ReplaceAll(unescaped, `\t`, "\t")
			unescaped = strings.ReplaceAll(unescaped, `\"`, `"`)
			unescaped = strings.ReplaceAll(unescaped, `\\`, `\`)
			output.WriteString(unescaped)
		default:
			output.WriteString(fmt.Sprint(v))
		}
	}
	fmt.Print(output.String())
}

func (e *Evaluator) evalIfStmt(stmt ast.IfStmt) string {
	conditionVal := e.evalExpr(stmt.Condition)
	if e.errors.HasErrors() {
		return "none"
	}
	condBool := toBool(conditionVal, e.line, e.errors)
	if e.errors.HasErrors() {
		return "none"
	}

	if condBool {
		return e.evalBlockWithControl(stmt.ThenBranch)
	}

	for _, elseIf := range stmt.ElseIfs {
		elseIfCond := e.evalExpr(elseIf.Condition)
		if e.errors.HasErrors() {
			return "none"
		}
		elseIfBool := toBool(elseIfCond, e.line, e.errors)
		if e.errors.HasErrors() {
			return "none"
		}
		if elseIfBool {
			return e.evalBlockWithControl(elseIf.Body)
		}
	}

	if stmt.ElseBranch != nil {
		return e.evalBlockWithControl(*stmt.ElseBranch)
	}
	return "none"
}

func (e *Evaluator) evalSwitchStmt(stmt ast.SwitchStmt) {
	switchValue := e.evalExpr(stmt.Expression)
	if e.errors.HasErrors() {
		return
	}

	matched := false
	for i, caseBranch := range stmt.Cases {
		caseValue := e.evalExpr(caseBranch.Value)
		if e.errors.HasErrors() {
			return
		}
		if matched || equals(switchValue, caseValue) {
			matched = true
			if e.evalBlockWithBreak(caseBranch.Body) {
				return
			}
			for j := i + 1; j < len(stmt.Cases); j++ {
				if e.evalBlockWithBreak(stmt.Cases[j].Body) {
					return
				}
			}
			if stmt.Default != nil {
				if e.evalBlockWithBreak(*stmt.Default) {
					return
				}
			}
			return
		}
	}

	if !matched && stmt.Default != nil {
		e.evalBlockWithBreak(*stmt.Default)
	}
}

func (e *Evaluator) evalForStmt(stmt ast.ForStmt) {
	if stmt.Init != nil {
		switch init := stmt.Init.(type) {
		case ast.VarDeclStmt:
			e.evalVarDeclStmt(init)
		case ast.ExprStmt:
			if assignExpr, ok := init.Expression.(ast.AssignmentExpr); ok {
				symbol, ok := assignExpr.Assignee.(ast.SymbolExpr)
				if !ok {
					e.errors.Report("Mowa, for loop init lo assignment ki variable undaali ra!", e.line)
					return
				}
				if !e.env.IsDeclared(symbol.Value) {
					e.errors.Report(fmt.Sprintf("Mowa, '%s' declare chesaka init cheyyali ra!", symbol.Value), e.line)
					return
				}
				rhs := e.evalExpr(assignExpr.Value)
				if e.errors.HasErrors() {
					return
				}
				e.env.Define(symbol.Value, rhs, e.line, e.errors)
			} else {
				e.evalExpr(init.Expression)
			}
		}
	}
	if e.errors.HasErrors() {
		return
	}

	for {
		condVal := e.evalExpr(stmt.Condition)
		if e.errors.HasErrors() {
			return
		}
		condBool := toBool(condVal, e.line, e.errors)
		if e.errors.HasErrors() {
			return
		}
		if !condBool {
			break
		}

		switch e.evalBlockWithControl(stmt.Body) {
		case "break":
			return
		case "continue":
		case "none":
		}

		e.evalExpr(stmt.Increment)
		if e.errors.HasErrors() {
			return
		}
	}
}

func (e *Evaluator) evalBlockWithControl(block ast.BlockStmt) string {
	for _, stmt := range block.Body {
		if stmt != nil {
			switch s := stmt.(type) {
			case ast.IfStmt:
				result := e.evalIfStmt(s)
				if result == "break" || result == "continue" {
					return result
				}
			case ast.BreakStmt:
				return "break"
			case ast.ContinueStmt:
				return "continue"
			default:
				e.evalStmt(stmt)
				if e.errors.HasErrors() {
					return "none"
				}
			}
		}
	}
	return "none"
}

func (e *Evaluator) evalBlock(block ast.BlockStmt) {
	for _, stmt := range block.Body {
		if stmt != nil {
			e.evalStmt(stmt)
			if e.errors.HasErrors() {
				break
			}
		}
	}
}

func (e *Evaluator) evalBlockWithBreak(block ast.BlockStmt) bool {
	for _, stmt := range block.Body {
		if stmt != nil {
			e.evalStmt(stmt)
			if e.errors.HasErrors() {
				return false
			}
			if _, isBreak := stmt.(ast.BreakStmt); isBreak {
				return true
			}
		}
	}
	return false
}

func toBool(value Value, line int, errors *errors.ErrorReporter) bool {
	switch v := value.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case string:
		return v != ""
	case nil:
		return false
	default:
		errors.Report(fmt.Sprintf("Mowa, '%v' ni boolean ki convert cheyalenu ra!", v), line)
		return false
	}
}

func (e *Evaluator) evalExpr(expr ast.Expr) Value {
	switch exprType := expr.(type) {
	case ast.NumberExpr:
		return exprType.Value
	case ast.StringExpr:
		return exprType.Value[1 : len(exprType.Value)-1]
	case ast.SymbolExpr:
		value, err := e.env.Get(exprType.Value)
		if err != nil {
			e.errors.Report(err.Error(), e.line)
			return nil
		}
		return value
	case ast.BinaryExpr:
		return e.evalBinaryExpr(exprType)
	case ast.PrefixExpr:
		return e.evalPrefixExpr(exprType)
	case ast.PostfixExpr:
		return e.evalPostfixExpr(exprType)
	case ast.AssignmentExpr:
		return e.evalAssignmentExpr(exprType)
	default:
		e.errors.Report(fmt.Sprintf("Mowa, ee expression type (%T) evaluate cheyalenu mowa!", expr), e.line)
		return nil
	}
}

func (e *Evaluator) evalPrefixExpr(expr ast.PrefixExpr) Value {
	switch expr.Operator.Kind {
	case lexer.PLUS_PLUS, lexer.MINUS_MINUS:
		symbolExpr, ok := expr.RightExpr.(ast.SymbolExpr)
		if !ok {
			e.errors.Report("Mowa, prefix ++ or -- ki variable undaali ra!", e.line)
			return nil
		}
		currentVal, err := e.env.Get(symbolExpr.Value)
		if err != nil {
			e.errors.Report(err.Error(), e.line)
			return nil
		}
		num, ok := currentVal.(float64)
		if !ok {
			e.errors.Report(fmt.Sprintf("Mowa, '%s' number undaali ra for ++ or --!", symbolExpr.Value), e.line)
			return nil
		}
		var newVal float64
		if expr.Operator.Kind == lexer.PLUS_PLUS {
			newVal = num + 1
		} else {
			newVal = num - 1
		}
		e.env.Define(symbolExpr.Value, newVal, e.line, e.errors)
		if e.errors.HasErrors() {
			return nil
		}
		return newVal
	default:
		e.errors.Report(fmt.Sprintf("Mowa, operator '%s' handle cheyalenu ra!", expr.Operator.Value), e.line)
		return nil
	}
}

func (e *Evaluator) evalPostfixExpr(expr ast.PostfixExpr) Value {
	switch expr.Operator.Kind {
	case lexer.PLUS_PLUS, lexer.MINUS_MINUS:
		symbolExpr, ok := expr.LeftExpr.(ast.SymbolExpr)
		if !ok {
			e.errors.Report("Mowa, postfix ++ or -- ki variable undaali ra!", e.line)
			return nil
		}
		currentVal, err := e.env.Get(symbolExpr.Value)
		if err != nil {
			e.errors.Report(err.Error(), e.line)
			return nil
		}
		num, ok := currentVal.(float64)
		if !ok {
			e.errors.Report(fmt.Sprintf("Mowa, '%s' number undaali ra for ++ or --!", symbolExpr.Value), e.line)
			return nil
		}
		var newVal float64
		if expr.Operator.Kind == lexer.PLUS_PLUS {
			newVal = num + 1
		} else {
			newVal = num - 1
		}
		e.env.Define(symbolExpr.Value, newVal, e.line, e.errors)
		if e.errors.HasErrors() {
			return nil
		}
		return num
	default:
		e.errors.Report(fmt.Sprintf("Mowa, operator '%s' handle cheyalenu ra!", expr.Operator.Value), e.line)
		return nil
	}
}

func (e *Evaluator) evalAssignmentExpr(expr ast.AssignmentExpr) Value {
	symbol, ok := expr.Assignee.(ast.SymbolExpr)
	if !ok {
		e.errors.Report("Mowa, assignment ki variable undaali ra!", e.line)
		return nil
	}
	if !e.env.IsDeclared(symbol.Value) {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' declare cheyaledhu ra before assignment!", symbol.Value), e.line)
		return nil
	}
	rhs := e.evalExpr(expr.Value)
	if e.errors.HasErrors() {
		return nil
	}
	currentVal, err := e.env.Get(symbol.Value)
	if err != nil {
		e.env.Define(symbol.Value, rhs, e.line, e.errors)
		return rhs
	}
	lhsNum, lhsOk := currentVal.(float64)
	rhsNum, rhsOk := rhs.(float64)
	if !lhsOk || !rhsOk {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' and '%v' numbers undaali ra for %s!", symbol.Value, rhs, expr.Operator.Value), e.line)
		return nil
	}
	var result float64
	switch expr.Operator.Kind {
	case lexer.ASSIGNMENT:
		result = rhsNum
	case lexer.PLUS_EQUALS:
		result = lhsNum + rhsNum
	case lexer.MINUS_EQUALS:
		result = lhsNum - rhsNum
	case lexer.STAR_EQUALS:
		result = lhsNum * rhsNum
	case lexer.SLASH_EQUALS:
		if rhsNum == 0 {
			e.errors.Report("Mowa, zero tho divide cheyakudadhu ra!", e.line)
			return nil
		}
		result = lhsNum / rhsNum
	case lexer.PERCENT_EQUALS:
		if rhsNum == 0 {
			e.errors.Report("Mowa, zero tho modulo cheyakudadhu ra!", e.line)
			return nil
		}
		result = math.Mod(lhsNum, rhsNum)
	default:
		e.errors.Report(fmt.Sprintf("Mowa, assignment operator '%s' handle cheyalenu ra!", expr.Operator.Value), e.line)
		return nil
	}
	e.env.Define(symbol.Value, result, e.line, e.errors)
	return result
}

func (e *Evaluator) evalBinaryExpr(expr ast.BinaryExpr) Value {
	left := e.evalExpr(expr.Left)
	if e.errors.HasErrors() {
		return nil
	}
	right := e.evalExpr(expr.Right)
	if e.errors.HasErrors() {
		return nil
	}
	switch expr.Operator.Kind {
	case lexer.PLUS:
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
			e.errors.Report(fmt.Sprintf("Mowa, number '%v' + non-number '%v' operation cheyakudadhu ra!", l, right), e.line)
		case string:
			if r, ok := right.(string); ok {
				return l + r
			}
			e.errors.Report(fmt.Sprintf("Mowa, string '%v' + non-string '%v' concatenate cheyakudadhu ra!", l, right), e.line)
		default:
			e.errors.Report(fmt.Sprintf("Mowa, '%v' + '%v' types match avvatledu ra!", l, right), e.line)
		}
	case lexer.DASH:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l - r
			}
			e.errors.Report(fmt.Sprintf("Mowa, '%v' - '%v' types match avvatledu ra!", l, right), e.line)
		}
	case lexer.STAR:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l * r
			}
			e.errors.Report(fmt.Sprintf("Mowa, '%v' * '%v' types match avvatledu ra!", l, right), e.line)
		}
	case lexer.SLASH:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					e.errors.Report("Mowa, zero tho divide cheyakudadhu!", e.line)
					return nil
				}
				return l / r
			}
			e.errors.Report(fmt.Sprintf("Mowa, '%v' / '%v' types match avvatledu ra!", l, right), e.line)
		}
	case lexer.PERCENT:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					e.errors.Report("Mowa, zero tho modulo cheyakudadhu ra!", e.line)
					return nil
				}
				return math.Mod(l, r)
			}
			e.errors.Report(fmt.Sprintf("Mowa, '%v' %% '%v' types match avvatledu ra!", l, right), e.line)
		}
	case lexer.STAR_STAR:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return math.Pow(l, r)
			}
			e.errors.Report(fmt.Sprintf("Mowa, '%v' ** '%v' types match avvatledu ra!", l, right), e.line)
		}
	case lexer.EQUALS:
		return equals(left, right)
	case lexer.NOT_EQUALS:
		return !equals(left, right)
	case lexer.LESS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l < r
			}
		}
		e.errors.Report(fmt.Sprintf("Mowa, '%v' < '%v' compare cheyalenu ra!", left, right), e.line)
	case lexer.LESS_EQUALS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l <= r
			}
		}
		e.errors.Report(fmt.Sprintf("Mowa, '%v' <= '%v' compare cheyalenu ra!", left, right), e.line)
	case lexer.GREATER:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l > r
			}
		}
		e.errors.Report(fmt.Sprintf("Mowa, '%v' > '%v' compare cheyalenu ra!", left, right), e.line)
	case lexer.GREATER_EQUALS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l >= r
			}
		}
		e.errors.Report(fmt.Sprintf("Mowa, '%v' >= '%v' compare cheyalenu ra!", left, right), e.line)
	case lexer.AND:
		lBool := toBool(left, e.line, e.errors)
		if e.errors.HasErrors() {
			return nil
		}
		if !lBool {
			return false
		}
		rBool := toBool(right, e.line, e.errors)
		return rBool
	case lexer.OR:
		lBool := toBool(left, e.line, e.errors)
		if e.errors.HasErrors() {
			return nil
		}
		if lBool {
			return true
		}
		rBool := toBool(right, e.line, e.errors)
		return rBool
	default:
		e.errors.Report(fmt.Sprintf("Mowa, '%s' operator handle cheyalenu mowa!", expr.Operator.Value), e.line)
	}
	return nil
}

func equals(left, right Value) bool {
	switch l := left.(type) {
	case float64:
		if r, ok := right.(float64); ok {
			return l == r
		}
	case string:
		if r, ok := right.(string); ok {
			return l == r
		}
	case bool:
		if r, ok := right.(bool); ok {
			return l == r
		}
	}
	return false
}
