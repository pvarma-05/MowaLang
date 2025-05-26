package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pvarma-05/MowaLang/src/errors"
	"github.com/pvarma-05/MowaLang/src/eval"
	"github.com/pvarma-05/MowaLang/src/lexer"
	"github.com/pvarma-05/MowaLang/src/parser"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

func runMowa(filePath string) *errors.ErrorReporter {
	if !strings.HasSuffix(strings.ToLower(filePath), ".mowa") {
		return &errors.ErrorReporter{
			Errors: []errors.MowaError{{Message: "Mowa, file '.mowa' extension tho undaali mowa!", LineNumber: 0}},
		}
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return &errors.ErrorReporter{
			Errors: []errors.MowaError{{Message: fmt.Sprintf("File ledhu Mowa: %v", err), LineNumber: 0}},
		}
	}

	tokens, lexErrors := lexer.Tokenize(string(bytes))
	if lexErrors.HasErrors() {
		return lexErrors
	}

	ast, parseErrors := parser.Parse(tokens)
	if parseErrors.HasErrors() {
		return parseErrors
	}

	evaluator := eval.NewEvaluator()
	return evaluator.Evaluate(ast)
}

func main() {
	var runFile string
	var prefActor string
	var showVersion bool

	var rootCmd = &cobra.Command{
		Use:                   "mowa <command>",
		Short:                 "MowaLang: A Custom Built Interpreter for Telugu Users",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			if showVersion {
				fmt.Printf("MowaLang v%s\n", version)
				return
			}

			if prefActor != "" {
				data, err := errors.DialoguesFS.ReadFile("dialogues.json")
				if err != nil {
					fmt.Printf("Mowa, dialogues.json load cheyalenu mowa: %v\n", err)
					os.Exit(1)
				}
				var dialogues map[string]any
				if err := json.Unmarshal(data, &dialogues); err != nil {
					fmt.Printf("Mowa, dialogues.json parse cheyalenu mowa: %v\n", err)
					os.Exit(1)
				}
				actor := strings.ToLower(prefActor)
				if _, exists := dialogues[actor]; !exists {
					fmt.Printf("Mowa, actor '%s' dialogues.json lo undaali mowa! Available: %v\n", actor, getActorKeys(dialogues))
					os.Exit(1)
				}

				configDir := filepath.Join(os.Getenv("HOME"), ".mowalang")
				if err := os.MkdirAll(configDir, 0755); err != nil {
					fmt.Printf("Mowa, config directory create cheyalenu mowa: %v\n", err)
					os.Exit(1)
				}
				configPath := filepath.Join(configDir, "config.json")
				config := map[string]string{"actor": actor}
				data, err = json.MarshalIndent(config, "", "  ")
				if err != nil {
					fmt.Printf("Mowa, config encode cheyalenu mowa: %v\n", err)
					os.Exit(1)
				}
				if err := os.WriteFile(configPath, data, 0644); err != nil {
					fmt.Printf("Mowa, config save cheyalenu mowa: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Mowa, actor preference set to '%s'!\n", actor)
				return
			}

			if runFile != "" {
				errors := runMowa(runFile)
				errors.PrintErrors()
				if errors.HasErrors() {
					os.Exit(1)
				}
				return
			}

			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Mowa, help display cheyalenu mowa: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.Flags().StringVarP(&runFile, "run", "r", "", "Run a MowaLang file")
	rootCmd.Flags().StringVarP(&prefActor, "preference", "p", "", "Set the dialogue actor preference")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print the version of MowaLang")

	rootCmd.SetHelpTemplate(`

 /$$      /$$                                         /$$
| $$$    /$$$                                        | $$
| $$$$  /$$$$  /$$$$$$  /$$  /$$  /$$  /$$$$$$       | $$        /$$$$$$  /$$$$$$$   /$$$$$$
| $$ $$/$$ $$ /$$__  $$| $$ | $$ | $$ |____  $$      | $$       |____  $$| $$__  $$ /$$__  $$
| $$  $$$| $$| $$  \ $$| $$ | $$ | $$  /$$$$$$$ /$$$$| $$        /$$$$$$$| $$  \ $$| $$  \ $$
| $$\  $ | $$| $$  | $$| $$ | $$ | $$ /$$__  $$|____/| $$       /$$__  $$| $$  | $$| $$  | $$
| $$ \/  | $$|  $$$$$$/|  $$$$$/$$$$/|  $$$$$$$      | $$$$$$$$|  $$$$$$$| $$  | $$|  $$$$$$$
|__/     |__/ \______/  \_____/\___/  \_______/      |________/ \_______/|__/  |__/ \____  $$
                                                                                    /$$  \ $$
                                                                                   |  $$$$$$/
                                                                                    \______/

----------------------------------------------------------------------------------------------

Thanks For Downloading MowaLang - A Custom Built Interpreter for Telugu Users

Description:
MowaLang is a dynamically typed, Telugu-based programming language built with Go.
It lets you write expressive programs using Telugu syntax and shows iconic Telugu movie dialogues when your code succeeds or fails.

Developed by Pradeep Varma

----------------------------------------------------------------------------------------------

USAGE:
  {{.UseLine}}

FLAGS:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}
`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Mowa, error mowa: %v\n", err)
		os.Exit(1)
	}
}

func getActorKeys(dialogues map[string]any) []string {
	keys := make([]string, 0, len(dialogues))
	for k := range dialogues {
		keys = append(keys, k)
	}
	return keys
}
