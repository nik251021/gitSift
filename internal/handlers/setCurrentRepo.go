package handlers

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

func SetCurrentRepoHandler(args []string) error {
	if len(args) < 1 {
		return errors.New("No repository given to function 'SetCurrentRepoHandler'")
	}
	env, _ := godotenv.Read()

	env["CURRENT_REPO_LINK"] = args[0]
	fmt.Print("che-to", env)

	return godotenv.Write(env, ".env")
}
