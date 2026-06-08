package handlers

import (
	"fmt"
	"os"
)

func GetCurrentRepoHandler(args []string) error {
	fmt.Println("Current repo: ", os.Getenv("CURRENT_REPO_LINK"))
	return nil
}
