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

	tokens, lexErrors := lexer.Tokenize(string(bytes))
	if lexErrors.HasErrors() {
		lexErrors.PrintErrors()
		os.Exit(1)
	}

	ast, parseErrors := parser.Parse(tokens)
	if parseErrors.HasErrors() {
		parseErrors.PrintErrors()
		os.Exit(1)
	}

	evaluator := eval.NewEvaluator()
	evalErrors := evaluator.Evaluate(ast)
	evalErrors.PrintErrors()
	if evalErrors.HasErrors() {
		os.Exit(1)
	}
}
