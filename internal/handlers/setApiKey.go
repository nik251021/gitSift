package handlers

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

func SetApiKeyHandler(args []string) error {
	if len(args) < 1 {
		fmt.Println("Please, enter an api key")
		return errors.New("No api key found in args")
	}
	env, _ := godotenv.Read()

	env["AGENT_API_KEY"] = args[0]

	return godotenv.Write(env, ".env")
}
