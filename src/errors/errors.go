package errors

import (
	"fmt"
	"math/rand"
	"time"
)

type MowaError struct {
	Message string
}

func (e MowaError) Error() string {
	return e.Message
}

type ErrorReporter struct {
	Errors []MowaError
}

func NewErrorReporter() *ErrorReporter {
	rand.Seed(time.Now().UnixNano())
	return &ErrorReporter{Errors: []MowaError{}}
}

func (r *ErrorReporter) Report(message string) {
	r.Errors = append(r.Errors, MowaError{Message: message})
}

func (r *ErrorReporter) HasErrors() bool {
	return len(r.Errors) > 0
}

var successDialogues = []string{
	"Baahubali: 'Naa saamrajyam lo error undadu ra!'",
	"Pushpa: 'Pushpa success ayipoyindi ra, flower anukunnava?'",
	"RRR: 'Fire and water combine ayi perfect run ayindi ra!'",
}

var failureDialogues = []string{
	"Baahubali: 'Naa kingdom lo ee error enduku ra!'",
	"Pushpa: 'Error vachindi ra, thaggede le!'",
	"KGF: 'Rocky bhai disappointment ayya ra, ee code fail!'",
}

func (r *ErrorReporter) PrintErrors() {
	if !r.HasErrors() {
		dialogue := successDialogues[rand.Intn(len(successDialogues))]
		fmt.Printf("%s\n\nOutput Mowa:\n", dialogue)
		return
	}

	dialogue := failureDialogues[rand.Intn(len(failureDialogues))]
	fmt.Printf("%s\nArey, errors unnai ra:\n", dialogue)
	for i, err := range r.Errors {
		fmt.Printf("%d. %s\n", i+1, err.Message)
	}
}

func (r *ErrorReporter) Clear() {
	r.Errors = []MowaError{}
}
