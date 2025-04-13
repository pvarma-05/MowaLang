package errors

import (
	"fmt"
	"math/rand"
	"time"
)

// stores a message and line number
type MowaError struct {
	Message    string
	LineNumber int
}

func (e MowaError) Error() string {
	return e.Message
}

// to collect and manage errors
type ErrorReporter struct {
	Errors []MowaError
}

func NewErrorReporter() *ErrorReporter {
	rand.Seed(time.Now().UnixNano())
	return &ErrorReporter{Errors: []MowaError{}}
}

func (r *ErrorReporter) Report(message string, lineNumber int) {
	r.Errors = append(r.Errors, MowaError{Message: message, LineNumber: lineNumber})
}

// Report with a line number
func (r *ErrorReporter) ReportSimple(message string) {
	r.Report(message, 0) // Fallback to 0 for compatibility
}

func (r *ErrorReporter) HasErrors() bool {
	return len(r.Errors) > 0
}

// Movie dialogues
var successDialogues = []string{
	"Prabhas: 'Jai Maahishmathi!'",
	"AA: 'Mowa... Assala Thaggedeley'",
	"NTR: 'Devara code raasinadu ante program run ayindhani ardham.'",
}

var failureDialogues = []string{
	"Prabhas: 'Thappu chesaav Devasena ...'",
	"AA: 'Konni Saarlu gelavadam kante, odipodame goppa'",
	"Bhramhanandham: 'Arey Tuppas Edhava, Thappu Chesavu ra'",
}

// PrintErrors with movie dialogues and colors
func (r *ErrorReporter) PrintErrors() {
	if !r.HasErrors() {
		// Success case: green dialogue
		dialogue := successDialogues[rand.Intn(len(successDialogues))]
		fmt.Printf("\n\n\033[32m%s\033[0m\n", dialogue)
		return
	}
	// Failure case: red dialogue
	dialogue := failureDialogues[rand.Intn(len(failureDialogues))]
	fmt.Printf("\n\033[31m%s\033[0m\nArey, errors unnai ra:\n", dialogue)
	for _, err := range r.Errors {
		if err.LineNumber > 0 {
			fmt.Printf("ln %d: %s\n", err.LineNumber, err.Message)
		} else {
			fmt.Printf("-: %s\n", err.Message) // Fallback
		}
	}
}

func (r *ErrorReporter) Clear() {
	r.Errors = []MowaError{}
}
