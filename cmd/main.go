package main

import (
	"os"

	"github.com/pvarma-05/MowaLang/src/lexer"
	"github.com/pvarma-05/MowaLang/src/parser"
	"github.com/sanity-io/litter"
)

func main() {
	bytes, err := os.ReadFile("examples/hello.mowa")

	if err != nil {
		panic(err)
	}

	tokens := lexer.Tokenize(string(bytes))

	// for _, token := range tokens {
	// 	token.Debug()
	// }

	ast := parser.Parse(tokens) // PARSER TESTING
	litter.Dump(ast)
}
