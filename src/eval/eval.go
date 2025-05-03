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

// Function represents a function definition at runtime.
type Function struct {
	Decl    ast.FunctionDeclStmt // Function declaration
	Closure *Environment         // Captured environment
}

// Environment holds variables, functions, and nested scopes.
type Environment struct {
	variables map[string]Value
	functions map[string]Function
	declared  map[string]bool
	types     map[string]string
	parent    *Environment // Parent environment for scope chaining
}

// NewEnvironment creates a new top-level environment.
func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]Value),
		functions: make(map[string]Function),
		declared:  make(map[string]bool),
		types:     make(map[string]string),
		parent:    nil,
	}
}

// NewChildEnvironment creates a nested environment.
func (env *Environment) NewChildEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]Value),
		functions: make(map[string]Function),
		declared:  make(map[string]bool),
		types:     make(map[string]string),
		parent:    env,
	}
}

// Declare marks a variable as declared with a type
func (env *Environment) Declare(name string, typ ast.Type) {
	env.declared[name] = true
	if symType, ok := typ.(ast.SymbolType); ok {
		env.types[name] = symType.Name
	} else if arrType, ok := typ.(ast.ArrayType); ok {
		if symType, ok := arrType.Underlying.(ast.SymbolType); ok {
			env.types[name] = fmt.Sprintf("[%s]", symType.Name)
		}
	}
}

// DeclareFunction stores a function in the environment.
func (env *Environment) DeclareFunction(name string, fn Function) {
	env.functions[name] = fn
	env.declared[name] = true
}

// Define sets a variable in the environment with type checking
func (env *Environment) Define(name string, value Value, lineNumber int, errors *errors.ErrorReporter) {
	if typ, exists := env.types[name]; exists {
		switch typ {
		case "number":
			if _, ok := value.(float64); !ok {
				errors.Report(fmt.Sprintf("Mowa, '%s' number type ki '%v' ivvakudadhu ra!", name, value), lineNumber)
				return
			}
		case "string":
			if _, ok := value.(string); !ok {
				errors.Report(fmt.Sprintf("Mowa, '%s' string type ki '%v' ivvakudadhu ra!", name, value), lineNumber)
				return
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				errors.Report(fmt.Sprintf("Mowa, '%s' boolean type ki '%v' ivvakudadhu ra!", name, value), lineNumber)
				return
			}
		case "[number]":
			if arr, ok := value.([]Value); ok {
				for _, elem := range arr {
					if _, ok := elem.(float64); !ok {
						errors.Report(fmt.Sprintf("Mowa, '%s' array lo number undaali, '%v' undhi ra!", name, elem), lineNumber)
						return
					}
				}
			} else {
				errors.Report(fmt.Sprintf("Mowa, '%s' array type ki '%v' ivvakudadhu ra!", name, value), lineNumber)
				return
			}
		case "[string]":
			if arr, ok := value.([]Value); ok {
				for _, elem := range arr {
					if _, ok := elem.(string); !ok {
						errors.Report(fmt.Sprintf("Mowa, '%s' array lo string undaali, '%v' undhi ra!", name, elem), lineNumber)
						return
					}
				}
			} else {
				errors.Report(fmt.Sprintf("Mowa, '%s' array type ki '%v' ivvakudadhu ra!", name, value), lineNumber)
				return
			}
		}
	}
	env.variables[name] = value
	env.declared[name] = true
}

// Get retrieves a variable’s value, searching parent scopes if needed.
func (env *Environment) Get(name string) (Value, error) {
	if value, exists := env.variables[name]; exists {
		return value, nil
	}
	if env.parent != nil {
		return env.parent.Get(name)
	}
	return nil, fmt.Errorf("mowa, '%s' ane variable define cheyaledhu ra", name)
}

// GetFunction retrieves a function, searching parent scopes if needed.
func (env *Environment) GetFunction(name string) (Function, error) {
	if fn, exists := env.functions[name]; exists {
		return fn, nil
	}
	if env.parent != nil {
		return env.parent.GetFunction(name)
	}
	return Function{}, fmt.Errorf("mowa, '%s' ane function define cheyaledhu ra", name)
}

// IsDeclared checks if a variable or function was declared, searching parent scopes.
func (env *Environment) IsDeclared(name string) bool {
	if env.declared[name] {
		return true
	}
	if env.parent != nil {
		return env.parent.IsDeclared(name)
	}
	return false
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
	case ast.InputIndexStmt:
		e.evalInputIndexStmt(s)
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
	case ast.FunctionDeclStmt:
		e.evalFunctionDeclStmt(s)
	case ast.ReturnStmt:
		e.evalReturnStmt(s)
	default:
		e.errors.Report(fmt.Sprintf("Mowa, ee statement type (%T) handle cheyalenu mowa!", s), e.line)
	}
}

func (e *Evaluator) evalFunctionDeclStmt(stmt ast.FunctionDeclStmt) {
	fn := Function{
		Decl:    stmt,
		Closure: e.env,
	}
	e.env.DeclareFunction(stmt.Name, fn)
}

func (e *Evaluator) evalReturnStmt(stmt ast.ReturnStmt) {
	var value Value
	if stmt.Value != nil {
		value = e.evalExpr(stmt.Value)
		if e.errors.HasErrors() {
			return
		}
	}
	// Store return value in a special variable to propagate it
	e.env.Define("_return", value, e.line, e.errors)
}

func (e *Evaluator) evalVarDeclStmt(stmt ast.VarDeclStmt) {
	var value Value
	if stmt.AssignedValue != nil {
		value = e.evalExpr(stmt.AssignedValue)
		if e.errors.HasErrors() {
			return
		}
	}

	// Handle array type with size expression
	if arrType, ok := stmt.ExplicitType.(ast.ArrayType); ok && arrType.Size != nil {
		sizeVal := e.evalExpr(arrType.Size)
		if e.errors.HasErrors() {
			return
		}
		size, ok := sizeVal.(float64)
		if !ok || size != float64(int(size)) || size < 0 {
			e.errors.Report(fmt.Sprintf("Mowa, array size '%v' positive integer undaali ra!", sizeVal), e.line)
			return
		}
		// Initialize empty array of specified size
		var underlyingType string
		if symType, ok := arrType.Underlying.(ast.SymbolType); ok {
			underlyingType = symType.Name
		}
		arr := make([]Value, int(size))
		if underlyingType == "number" {
			for i := range arr {
				arr[i] = float64(0)
			}
		} else if underlyingType == "string" {
			for i := range arr {
				arr[i] = ""
			}
		}
		value = arr
	}

	// Infer type for implicit array declarations
	if stmt.ExplicitType == nil && value != nil {
		if arr, ok := value.([]Value); ok && len(arr) > 0 {
			firstElem := arr[0]
			var inferredType ast.Type
			switch firstElem.(type) {
			case float64:
				inferredType = ast.ArrayType{Underlying: ast.SymbolType{Name: "number"}}
			case string:
				inferredType = ast.ArrayType{Underlying: ast.SymbolType{Name: "string"}}
			default:
				e.errors.Report("Mowa, array type infer cheyalenu ra!", e.line)
				return
			}
			e.env.Declare(stmt.VarName, inferredType)
		}
	} else {
		e.env.Declare(stmt.VarName, stmt.ExplicitType)
	}

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
	typ, exists := e.env.types[stmt.VarName]
	if !exists {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' type define cheyaledhu ra!", stmt.VarName), e.line)
		return
	}
	if typ == "number" {
		if num, err := strconv.ParseFloat(input, 64); err == nil {
			e.env.Define(stmt.VarName, num, e.line, e.errors)
		} else {
			e.errors.Report(fmt.Sprintf("Mowa, '%s' number undaali ra, '%s' ichav!", stmt.VarName, input), e.line)
		}
	} else if typ == "string" {
		e.env.Define(stmt.VarName, input, e.line, e.errors)
	} else if typ == "boolean" {
		input = strings.ToLower(input)
		if input == "nijam" || input == "true" {
			e.env.Define(stmt.VarName, true, e.line, e.errors)
		} else if input == "adhey" || input == "false" {
			e.env.Define(stmt.VarName, false, e.line, e.errors)
		} else {
			e.errors.Report(fmt.Sprintf("Mowa, '%s' boolean undaali ra, '%s' ichav!", stmt.VarName, input), e.line)
		}
	} else {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' scalar type undaali ra for direct input!", stmt.VarName), e.line)
	}
}

func (e *Evaluator) evalInputIndexStmt(stmt ast.InputIndexStmt) {
	array := e.evalExpr(stmt.Array)
	if e.errors.HasErrors() {
		return
	}
	arr, ok := array.([]Value)
	if !ok {
		e.errors.Report(fmt.Sprintf("Mowa, '%v' array kaadhu ra!", array), e.line)
		return
	}
	index := e.evalExpr(stmt.Index)
	if e.errors.HasErrors() {
		return
	}
	idx, ok := index.(float64)
	if !ok || idx != float64(int(idx)) || idx < 0 || int(idx) >= len(arr) {
		e.errors.Report(fmt.Sprintf("Mowa, invalid array index '%v' ra!", index), e.line)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	varName := stmt.Array.(ast.SymbolExpr).Value
	typ, exists := e.env.types[varName]
	if !exists {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' type define cheyaledhu ra!", varName), e.line)
		return
	}
	var value Value
	if typ == "[number]" {
		if num, err := strconv.ParseFloat(input, 64); err == nil {
			value = num
		} else {
			e.errors.Report(fmt.Sprintf("Mowa, array '%s' lo number undaali ra, '%s' ichav!", varName, input), e.line)
			return
		}
	} else if typ == "[string]" {
		value = input
	} else {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' array type undaali ra!", varName), e.line)
		return
	}
	arr[int(idx)] = value
	e.env.Define(varName, arr, e.line, e.errors)
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
			// Handle escaped sequences
			if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
				v = v[1 : len(v)-1]
			}
			v = strings.ReplaceAll(v, `\n`, "\n")
			v = strings.ReplaceAll(v, `\t`, "\t")
			v = strings.ReplaceAll(v, `\"`, `"`)
			v = strings.ReplaceAll(v, `\\`, `\`)
			output.WriteString(v)
		case []Value:
			output.WriteString("[")
			for i, elem := range v {
				if i > 0 {
					output.WriteString(", ")
				}
				switch elem := elem.(type) {
				case float64:
					output.WriteString(fmt.Sprintf("%g", elem))
				case string:
					output.WriteString(fmt.Sprintf("%q", elem))
				default:
					output.WriteString(fmt.Sprint(elem))
				}
			}
			output.WriteString("]")
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
		return e.evalBlockWithReturn(stmt.ThenBranch)
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
			return e.evalBlockWithReturn(elseIf.Body)
		}
	}

	if stmt.ElseBranch != nil {
		return e.evalBlockWithReturn(*stmt.ElseBranch)
	}
	return "none"
}

func (e *Evaluator) evalSwitchStmt(stmt ast.SwitchStmt) string {
	switchValue := e.evalExpr(stmt.Expression)
	if e.errors.HasErrors() {
		return "none"
	}

	matched := false
	for i, caseBranch := range stmt.Cases {
		caseValue := e.evalExpr(caseBranch.Value)
		if e.errors.HasErrors() {
			return "none"
		}
		if matched || equals(switchValue, caseValue) {
			matched = true
			if result := e.evalBlockWithReturn(caseBranch.Body); result != "none" {
				return result
			}
			for j := i + 1; j < len(stmt.Cases); j++ {
				if result := e.evalBlockWithReturn(stmt.Cases[j].Body); result != "none" {
					return result
				}
			}
			if stmt.Default != nil {
				if result := e.evalBlockWithReturn(*stmt.Default); result != "none" {
					return result
				}
			}
			return "none"
		}
	}

	if !matched && stmt.Default != nil {
		return e.evalBlockWithReturn(*stmt.Default)
	}
	return "none"
}

func (e *Evaluator) evalForStmt(stmt ast.ForStmt) string {
	if stmt.Init != nil {
		switch init := stmt.Init.(type) {
		case ast.VarDeclStmt:
			e.evalVarDeclStmt(init)
		case ast.ExprStmt:
			if assignExpr, ok := init.Expression.(ast.AssignmentExpr); ok {
				symbol, ok := assignExpr.Assignee.(ast.SymbolExpr)
				if !ok {
					e.errors.Report("Mowa, for loop init lo assignment ki variable undaali ra!", e.line)
					return "none"
				}
				if !e.env.IsDeclared(symbol.Value) {
					e.errors.Report(fmt.Sprintf("Mowa, '%s' declare chesaka init cheyyali ra!", symbol.Value), e.line)
					return "none"
				}
				rhs := e.evalExpr(assignExpr.Value)
				if e.errors.HasErrors() {
					return "none"
				}
				e.env.Define(symbol.Value, rhs, e.line, e.errors)
			} else {
				e.evalExpr(init.Expression)
			}
		}
	}
	if e.errors.HasErrors() {
		return "none"
	}

	for {
		condVal := e.evalExpr(stmt.Condition)
		if e.errors.HasErrors() {
			return "none"
		}
		condBool := toBool(condVal, e.line, e.errors)
		if e.errors.HasErrors() {
			return "none"
		}
		if !condBool {
			break
		}

		switch e.evalBlockWithReturn(stmt.Body) {
		case "break":
			return "none"
		case "continue":
			// Skip to increment
		case "return":
			return "return"
		case "none":
			// Normal execution
		}

		e.evalExpr(stmt.Increment)
		if e.errors.HasErrors() {
			return "none"
		}
	}
	return "none"
}

func (e *Evaluator) evalBlockWithReturn(block ast.BlockStmt) string {
	for _, stmt := range block.Body {
		if stmt != nil {
			switch s := stmt.(type) {
			case ast.IfStmt:
				result := e.evalIfStmt(s)
				if result == "break" || result == "continue" || result == "return" {
					return result
				}
			case ast.ReturnStmt:
				e.evalReturnStmt(s)
				if e.errors.HasErrors() {
					return "none"
				}
				return "return"
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
	case ast.BoolExpr:
		return exprType.Value
	case ast.SymbolExpr:
		value, err := e.env.Get(exprType.Value)
		if err != nil {
			// Check if it's a function
			if _, fnErr := e.env.GetFunction(exprType.Value); fnErr == nil {
				return exprType.Value // Return symbol for function names
			}
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
	case ast.ArrayLiteralExpr:
		return e.evalArrayLiteralExpr(exprType)
	case ast.ArrayIndexExpr:
		return e.evalArrayIndexExpr(exprType)
	case ast.MemberAccessExpr:
		return e.evalMemberAccessExpr(exprType)
	case ast.CallExpr:
		return e.evalCallExpr(exprType)
	case ast.TypeofExpr:
		return e.evalTypeofExpr(exprType)
	default:
		e.errors.Report(fmt.Sprintf("Mowa, ee expression type (%T) evaluate cheyalenu mowa!", expr), e.line)
		return nil
	}
}

func (e *Evaluator) evalTypeofExpr(expr ast.TypeofExpr) Value {
	var value Value
	var varName string
	switch arg := expr.Right.(type) {
	case ast.SymbolExpr:
		varName = arg.Value
		// Check if it's a function
		if _, err := e.env.GetFunction(arg.Value); err == nil {
			return "function"
		}
		var err error
		value, err = e.env.Get(arg.Value)
		if err != nil {
			e.errors.Report(err.Error(), e.line)
			return nil
		}
	case ast.ArrayIndexExpr:
		varName = arg.Array.(ast.SymbolExpr).Value
		value = e.evalExpr(arg)
		if e.errors.HasErrors() {
			return nil
		}
	default:
		e.errors.Report("Mowa, rakam function variable or array index thiskuntadhi ra!", e.line)
		return nil
	}

	switch v := value.(type) {
	case float64:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	case []Value:
		if len(v) > 0 {
			switch v[0].(type) {
			case float64:
				return "[number]"
			case string:
				return "[string]"
			}
		}
		if typ, exists := e.env.types[varName]; exists {
			return typ
		}
		return "[unknown]"
	default:
		e.errors.Report(fmt.Sprintf("Mowa, '%v' type determine cheyalenu ra!", v), e.line)
		return nil
	}
}

func (e *Evaluator) evalCallExpr(expr ast.CallExpr) Value {
	symbol, isSymbol := expr.Function.(ast.SymbolExpr)
	if !isSymbol {
		e.errors.Report("Mowa, function call lo variable name undaali ra!", e.line)
		return nil
	}

	fn, err := e.env.GetFunction(symbol.Value)
	if err != nil {
		e.errors.Report(err.Error(), e.line)
		return nil
	}

	args := []Value{}
	for _, arg := range expr.Arguments {
		value := e.evalExpr(arg)
		if e.errors.HasErrors() {
			return nil
		}
		args = append(args, value)
	}

	if len(args) != len(fn.Decl.Parameters) {
		e.errors.Report(fmt.Sprintf("Mowa, '%s' function %d parameters expect chesthundi, %d ichav ra!", fn.Decl.Name, len(fn.Decl.Parameters), len(args)), e.line)
		return nil
	}

	callEnv := fn.Closure.NewChildEnvironment()
	prevEnv := e.env
	e.env = callEnv
	defer func() { e.env = prevEnv }()

	for i, param := range fn.Decl.Parameters {
		callEnv.Declare(param.Name, param.Type)
		callEnv.Define(param.Name, args[i], e.line, e.errors)
		if e.errors.HasErrors() {
			return nil
		}
	}

	_ = e.evalBlockWithReturn(fn.Decl.Body)
	if e.errors.HasErrors() {
		return nil
	}

	retVal, err := callEnv.Get("_return")
	if err != nil {
		if fn.Decl.ReturnType != nil {
			e.errors.Report(fmt.Sprintf("Mowa, '%s' function return cheyyali ra!", fn.Decl.Name), e.line)
			return nil
		}
		return nil
	}

	if fn.Decl.ReturnType != nil {
		switch rt := fn.Decl.ReturnType.(type) {
		case ast.SymbolType:
			switch rt.Name {
			case "number":
				if _, ok := retVal.(float64); !ok {
					e.errors.Report(fmt.Sprintf("Mowa, '%s' function number return cheyyali, '%v' ichav ra!", fn.Decl.Name, retVal), e.line)
					return nil
				}
			case "string":
				if _, ok := retVal.(string); !ok {
					e.errors.Report(fmt.Sprintf("Mowa, '%s' function string return cheyyali, '%v' ichav ra!", fn.Decl.Name, retVal), e.line)
					return nil
				}
			case "boolean":
				if _, ok := retVal.(bool); !ok {
					e.errors.Report(fmt.Sprintf("Mowa, '%s' function boolean return cheyyali, '%v' ichav ra!", fn.Decl.Name, retVal), e.line)
					return nil
				}
			}
		case ast.ArrayType:
			if arr, ok := retVal.([]Value); ok {
				underlying, ok := rt.Underlying.(ast.SymbolType)
				if !ok {
					e.errors.Report("Mowa, array underlying type symbol undaali ra!", e.line)
					return nil
				}
				for _, elem := range arr {
					if underlying.Name == "number" && !isNumber(elem) {
						e.errors.Report(fmt.Sprintf("Mowa, '%s' array lo number undaali, '%v' undhi ra!", fn.Decl.Name, elem), e.line)
						return nil
					}
					if underlying.Name == "string" && !isString(elem) {
						e.errors.Report(fmt.Sprintf("Mowa, '%s' array lo string undaali, '%v' undhi ra!", fn.Decl.Name, elem), e.line)
						return nil
					}
				}
			} else {
				e.errors.Report(fmt.Sprintf("Mowa, '%s' array return cheyyali ra!", fn.Decl.Name), e.line)
				return nil
			}
		}
	}

	return retVal
}

// Helper functions for type checking
func isNumber(v Value) bool {
	_, ok := v.(float64)
	return ok
}

func isString(v Value) bool {
	_, ok := v.(string)
	return ok
}

func (e *Evaluator) evalArrayLiteralExpr(expr ast.ArrayLiteralExpr) Value {
	elements := []Value{}
	for _, elem := range expr.Elements {
		value := e.evalExpr(elem)
		if e.errors.HasErrors() {
			return nil
		}
		elements = append(elements, value)
	}
	return elements
}

func (e *Evaluator) evalArrayIndexExpr(expr ast.ArrayIndexExpr) Value {
	array := e.evalExpr(expr.Array)
	index := e.evalExpr(expr.Index)
	if e.errors.HasErrors() {
		return nil
	}
	arr, ok := array.([]Value)
	if !ok {
		e.errors.Report(fmt.Sprintf("Mowa, '%v' array kaadhu ra!", array), e.line)
		return nil
	}
	idx, ok := index.(float64)
	if !ok || idx != float64(int(idx)) || idx < 0 || int(idx) >= len(arr) {
		e.errors.Report(fmt.Sprintf("Mowa, invalid array index '%v' ra!", index), e.line)
		return nil
	}
	return arr[int(idx)]
}

func (e *Evaluator) evalMemberAccessExpr(expr ast.MemberAccessExpr) Value {
	object := e.evalExpr(expr.Object)
	if e.errors.HasErrors() {
		return nil
	}
	arr, ok := object.([]Value)
	if !ok {
		e.errors.Report(fmt.Sprintf("Mowa, '%v' array kaadhu ra!", object), e.line)
		return nil
	}
	if expr.Property == "length" {
		return float64(len(arr))
	}
	e.errors.Report(fmt.Sprintf("Mowa, unknown property '%s' ra!", expr.Property), e.line)
	return nil
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
	rhs := e.evalExpr(expr.Value)
	if e.errors.HasErrors() {
		return nil
	}

	switch assignee := expr.Assignee.(type) {
	case ast.SymbolExpr:
		if !e.env.IsDeclared(assignee.Value) {
			e.errors.Report(fmt.Sprintf("Mowa, '%s' declare cheyaledhu ra before assignment!", assignee.Value), e.line)
			return nil
		}
		currentVal, err := e.env.Get(assignee.Value)
		if err != nil {
			e.env.Define(assignee.Value, rhs, e.line, e.errors)
			return rhs
		}
		lhsNum, lhsOk := currentVal.(float64)
		rhsNum, rhsOk := rhs.(float64)
		if !lhsOk || !rhsOk {
			e.errors.Report(fmt.Sprintf("Mowa, '%s' and '%v' numbers undaali ra for %s!", assignee.Value, rhs, expr.Operator.Value), e.line)
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
		e.env.Define(assignee.Value, result, e.line, e.errors)
		return result
	case ast.ArrayIndexExpr:
		array, err := e.env.Get(assignee.Array.(ast.SymbolExpr).Value)
		if err != nil {
			e.errors.Report(err.Error(), e.line)
			return nil
		}
		arr, ok := array.([]Value)
		if !ok {
			e.errors.Report(fmt.Sprintf("Mowa, '%v' array kaadhu ra!", array), e.line)
			return nil
		}
		index := e.evalExpr(assignee.Index)
		if e.errors.HasErrors() {
			return nil
		}
		idx, ok := index.(float64)
		if !ok || idx != float64(int(idx)) || idx < 0 || int(idx) >= len(arr) {
			e.errors.Report(fmt.Sprintf("Mowa, invalid array index '%v' ra!", index), e.line)
			return nil
		}
		typ, exists := e.env.types[assignee.Array.(ast.SymbolExpr).Value]
		if exists {
			if typ == "[number]" {
				if _, ok := rhs.(float64); !ok {
					e.errors.Report(fmt.Sprintf("Mowa, array '%s' lo number undaali, '%v' undhi ra!", assignee.Array.(ast.SymbolExpr).Value, rhs), e.line)
					return nil
				}
			} else if typ == "[string]" {
				if _, ok := rhs.(string); !ok {
					e.errors.Report(fmt.Sprintf("Mowa, array '%s' lo string undaali, '%v' undhi ra!", assignee.Array.(ast.SymbolExpr).Value, rhs), e.line)
					return nil
				}
			}
		}
		arr[int(idx)] = rhs
		e.env.Define(assignee.Array.(ast.SymbolExpr).Value, arr, e.line, e.errors)
		return rhs
	default:
		e.errors.Report("Mowa, assignment ki variable or array index undaali ra!", e.line)
		return nil
	}
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
