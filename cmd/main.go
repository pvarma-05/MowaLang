package main

import (
	"fmt"
	"os"

	"github.com/pvarma-05/MowaLang/src/eval"
	"github.com/pvarma-05/MowaLang/src/lexer"
	"github.com/pvarma-05/MowaLang/src/parser"
)

func main() {
	bytes, err := os.ReadFile("examples/hello.mowa")
	if err != nil {
		fmt.Printf("File ledhu Mowa : %v\n", err)
		os.Exit(1)
	}

	// Step 1: Lexing
	tokens, lexErrors := lexer.Tokenize(string(bytes))
	if lexErrors.HasErrors() {
		lexErrors.PrintErrors()
		os.Exit(1)
	}

	// Step 2: Parsing
	ast, parseErrors := parser.Parse(tokens)
	if parseErrors.HasErrors() {
		parseErrors.PrintErrors()
		os.Exit(1)
	}

	// Step 3: Evaluation
	evaluator := eval.NewEvaluator()
	evalErrors := evaluator.Evaluate(ast)
	evalErrors.PrintErrors()
	if evalErrors.HasErrors() {
		os.Exit(1)
	}
}
