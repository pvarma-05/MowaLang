package errors

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

//go:embed dialogues.json
var DialoguesFS embed.FS

type MowaError struct {
	Message    string
	LineNumber int
}

func (e MowaError) Error() string {
	return e.Message
}

type ErrorReporter struct {
	Errors []MowaError
}

type ActorDialogues struct {
	Success []string `json:"success"`
	Failure []string `json:"failure"`
}

var dialogues map[string]ActorDialogues

func init() {
	// Load dialogues from embedded JSON
	data, err := DialoguesFS.ReadFile("dialogues.json")
	if err != nil {
		panic(fmt.Sprintf("Mowa, dialogues.json load cheyalenu mowa: %v", err))
	}
	if err := json.Unmarshal(data, &dialogues); err != nil {
		panic(fmt.Sprintf("Mowa, dialogues.json parse cheyalenu mowa: %v", err))
	}
	rand.Seed(time.Now().UnixNano())
}

func NewErrorReporter() *ErrorReporter {
	return &ErrorReporter{Errors: []MowaError{}}
}

func (r *ErrorReporter) Report(message string, lineNumber int) {
	r.Errors = append(r.Errors, MowaError{Message: message, LineNumber: lineNumber})
}

func (r *ErrorReporter) ReportSimple(message string) {
	r.Report(message, 0)
}

func (r *ErrorReporter) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ErrorReporter) PrintErrors() {
	// Load actor preference
	actor := "all"
	configPath := filepath.Join(os.Getenv("HOME"), ".mowalang", "config.json")
	if configData, err := os.ReadFile(configPath); err == nil {
		var config map[string]string
		if json.Unmarshal(configData, &config) == nil {
			if pref, exists := config["actor"]; exists {
				actor = pref
			}
		}
	}

	// Select dialogues based on preference
	actorDialogues, exists := dialogues[actor]
	if !exists {
		actorDialogues = dialogues["all"] // Fallback to "all"
	}

	if !r.HasErrors() {
		dialogue := actorDialogues.Success[rand.Intn(len(actorDialogues.Success))]
		fmt.Printf("\n\n\033[32m%s\033[0m\n", dialogue)
		return
	}
	dialogue := actorDialogues.Failure[rand.Intn(len(actorDialogues.Failure))]
	fmt.Printf("\n\033[31m%s\033[0m\nMowa, errors unnai mowa:\n", dialogue)
	for _, err := range r.Errors {
		if err.LineNumber > 0 {
			fmt.Printf("ln %d: %s\n", err.LineNumber, err.Message)
		} else {
			fmt.Printf("-: %s\n", err.Message)
		}
	}
}
