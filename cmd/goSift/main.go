package main

import (
	"fmt"
	"os"

	"goSift/internal/handlers"

	"github.com/joho/godotenv"
)

type CommandHandler func(args []string) error

func main() {
	godotenv.Load()

	commands := map[string]CommandHandler{
		"help":                 handlers.HelpHandler,
		"getAPIkey":            handlers.GetApiKeyHandler,
		"setAPIkey":            handlers.SetApiKeyHandler,
		"setCurrentLinkToRepo": handlers.SetCurrentRepoHandler,
		"getCurrentLinkToRepo": handlers.GetCurrentRepoHandler,
	}
	if len(os.Args) < 2 {
		fmt.Println("Please, enter any arguments, or use goSift help")
		return
	}
	curcmd := os.Args[1]
	args := os.Args[2:]

	if handler, exists := commands[curcmd]; exists {
		err := handler(args)
		if err != nil {
			fmt.Println("Error in command", curcmd, err)
			return
		}
	} else {
		fmt.Println("Command not found", curcmd)
		return
	}
}
