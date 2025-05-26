//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"
	"syscall/js"

	"github.com/pvarma-05/MowaLang/src/errors"
	"github.com/pvarma-05/MowaLang/src/eval"
	"github.com/pvarma-05/MowaLang/src/lexer"
	"github.com/pvarma-05/MowaLang/src/parser"
)

// custom io.Writer to capture output in WASM
type outputWriter struct {
	buffer bytes.Buffer
}

func (w *outputWriter) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

func (w *outputWriter) String() string {
	return w.buffer.String()
}

// loads dialogues
func loadDialogues() (errors.ActorDialogues, error) {
	data, err := errors.DialoguesFS.ReadFile("dialogues.json")
	if err != nil {
		return errors.ActorDialogues{}, fmt.Errorf("mowa, dialogues.json load cheyalenu mowa: %v", err)
	}
	var dialogues map[string]errors.ActorDialogues
	if err := json.Unmarshal(data, &dialogues); err != nil {
		return errors.ActorDialogues{}, fmt.Errorf("mowa, dialogues.json parse cheyalenu mowa: %v", err)
	}
	actorDialogues, exists := dialogues["all"]
	if !exists {
		return errors.ActorDialogues{}, fmt.Errorf("mowa, 'all' dialogues not found in dialogues.json mowa")
	}
	return actorDialogues, nil
}

// executes the MowaLang code and returns the combined program output and result
func runMowa(code string, inputCallback js.Value) (string, error) {
	// Load dialogues
	actorDialogues, err := loadDialogues()
	if err != nil {
		return err.Error(), nil
	}

	// Custom writer to capture output
	writer := &outputWriter{}

	// Input handler
	eval.SetInputHandler(func(prompt string) string {
		return inputCallback.Invoke(prompt).String()
	})

	// Tokenize
	tokens, lexErrors := lexer.Tokenize(code)
	if lexErrors.HasErrors() {
		return formatErrorOutput(writer.String(), actorDialogues, lexErrors), nil
	}

	// Parse
	ast, parseErrors := parser.Parse(tokens)
	if parseErrors.HasErrors() {
		return formatErrorOutput(writer.String(), actorDialogues, parseErrors), nil
	}

	// Evaluate, redirecting print output to our writer
	evaluator := eval.NewEvaluator()
	eval.SetPrintHandler(func(s string) {
		fmt.Fprint(writer, s)
	})
	errReporter := evaluator.Evaluate(ast)

	// Combine program output with dialogue/error message
	programOutput := strings.TrimSpace(writer.String())
	result := errReporterToString(errReporter, actorDialogues)
	if programOutput != "" {
		if result != "" {
			return programOutput + "\n" + result, nil
		}
		return programOutput, nil
	}
	return result, nil
}

// lexer/parser errors with a failure dialogue
func formatErrorOutput(programOutput string, actorDialogues errors.ActorDialogues, errReporter *errors.ErrorReporter) string {
	var output strings.Builder
	if len(actorDialogues.Failure) == 0 {
		output.WriteString("[DIALOGUE]Mowa, no failure dialogues found for 'all' mowa![/DIALOGUE]\n")
	} else {
		output.WriteString(fmt.Sprintf("[DIALOGUE]%s[/DIALOGUE]\n", actorDialogues.Failure[rand.IntN(len(actorDialogues.Failure))]))
	}
	output.WriteString("Mowa, errors unnai mowa:\n")
	for _, err := range errReporter.Errors {
		if err.LineNumber > 0 {
			output.WriteString(fmt.Sprintf("ln %d: %s\n", err.LineNumber, err.Message))
		} else {
			output.WriteString(fmt.Sprintf("-: %s\n", err.Message))
		}
	}
	if programOutput != "" {
		return programOutput + "\n" + output.String()
	}
	return output.String()
}

// converts evaluation errors or success dialogues to a string
func errReporterToString(errReporter *errors.ErrorReporter, actorDialogues errors.ActorDialogues) string {
	if !errReporter.HasErrors() {
		if len(actorDialogues.Success) == 0 {
			return "[DIALOGUE]Mowa, no success dialogues found for 'all' mowa![/DIALOGUE]"
		}
		return fmt.Sprintf("[DIALOGUE]%s[/DIALOGUE]", actorDialogues.Success[rand.IntN(len(actorDialogues.Success))])
	}
	return formatErrorOutput("", actorDialogues, errReporter)
}

// jsRunMowa is the JavaScript-exposed function.
func jsRunMowa(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{
			"error": "Mowa, code and input callback arguments undaali mowa!",
		}
	}
	code := args[0].String()
	inputCallback := args[1]
	result, err := runMowa(code, inputCallback)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}
	return map[string]interface{}{
		"result": result,
	}
}

func main() {
	// Keep the WASM module alive
	c := make(chan struct{}, 0)
	// Register the runMowa function to JavaScript
	js.Global().Set("runMowa", js.FuncOf(jsRunMowa))
	<-c
}
