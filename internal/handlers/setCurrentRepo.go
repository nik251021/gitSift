package handlers

import (
	"errors"

	"github.com/joho/godotenv"
)

func SetCurrentRepoHandler(args []string) error {
	if len(args) < 1 {
		return errors.New("No repository given to function 'SetCurrentRepoHandler'")
	}
	env, _ := godotenv.Read()

	env["CURRENT_REPO_LINK"] = args[0]

	return godotenv.Write(env, ".env")
}
