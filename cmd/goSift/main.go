package main

import (
	"fmt"
	"os"

	"goSift/internal/handlers"

	"github.com/joho/godotenv"
)

type CommandHandler func(args []string) error

func main() {
	_ = godotenv.Load()

	commands := map[string]CommandHandler{
		"help":                 handlers.HelpHandler,
		"getAPIkey":            handlers.GetApiKeyHandler,
		"setAPIkey":            handlers.SetApiKeyHandler,
		"setCurrentLinkToRepo": handlers.SetCurrentRepoHandler,
		"getCurrentLinkToRepo": handlers.GetCurrentRepoHandler,
		"createContext":        handlers.CreateContext,
		"sift":                 handlers.SiftHandler,
		"ask":                  handlers.AskHandler,
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage error: please provide a command or run 'goSift help'")
		return
	}

	curcmd := os.Args[1]
	args := os.Args[2:]

	handler, exists := commands[curcmd]
	if !exists {
		fmt.Printf("Error: command '%s' not found\n", curcmd)
		return
	}

	if err := handler(args); err != nil {
		fmt.Printf("Error execution failed for command '%s': %v\n", curcmd, err)
		return
	}
}
