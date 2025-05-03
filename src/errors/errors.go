package errors

import (
	"fmt"
	"math/rand"
	"time"
)

type MowaError struct {
	Message    string
	LineNumber int
}

// Error implements the error interface for MowaError.
func (e MowaError) Error() string {
	return e.Message
}

// ErrorReporter collects and manages errors during lexing, parsing, or evaluation.
type ErrorReporter struct {
	Errors []MowaError
}

func NewErrorReporter() *ErrorReporter {
	rand.Seed(time.Now().UnixNano())
	return &ErrorReporter{Errors: []MowaError{}}
}

// Report adds an error with a message and line number to the Errors list.
func (r *ErrorReporter) Report(message string, lineNumber int) {
	r.Errors = append(r.Errors, MowaError{Message: message, LineNumber: lineNumber})
}

// ReportSimple adds an error without a line number (fallback).
func (r *ErrorReporter) ReportSimple(message string) {
	r.Report(message, 0)
}

func (r *ErrorReporter) HasErrors() bool {
	return len(r.Errors) > 0
}

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

// Uses ANSI colors: red for failure, green for success.
func (r *ErrorReporter) PrintErrors() {
	if !r.HasErrors() {
		dialogue := successDialogues[rand.Intn(len(successDialogues))]
		fmt.Printf("\n\n\033[32m%s\033[0m\n", dialogue)
		return
	}
	dialogue := failureDialogues[rand.Intn(len(failureDialogues))]
	fmt.Printf("\n\033[31m%s\033[0m\nArey, errors unnai ra:\n", dialogue)
	for _, err := range r.Errors {
		if err.LineNumber > 0 {
			fmt.Printf("ln %d: %s\n", err.LineNumber, err.Message)
		} else {
			fmt.Printf("-: %s\n", err.Message) // Fallback for errors without line numbers
		}
	}
}

// // Clear resets the error list (not currently used but available for future extensions).
// func (r *ErrorReporter) Clear() {
// 	r.Errors = []MowaError{}
// }
