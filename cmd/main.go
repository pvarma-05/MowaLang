package main

import (
	"fmt"
	"os"

	"github.com/pvarma-05/MowaLang/src/lexer"
	"github.com/pvarma-05/MowaLang/src/parser"
	"github.com/sanity-io/litter"
)

func main() {
	bytes, err := os.ReadFile("examples/hello.mowa")
	if err != nil {
		fmt.Printf("File read cheyaleka ra: %v\n", err)
		os.Exit(1)
	}

	tokens, lexErrors := lexer.Tokenize(string(bytes))
	if lexErrors.HasErrors() {
		lexErrors.PrintErrors()
		os.Exit(1)
	}

	ast, parseErrors := parser.Parse(tokens)
	parseErrors.PrintErrors()
	if parseErrors.HasErrors() {
		os.Exit(1)
	}
	litter.Dump(ast)
}
