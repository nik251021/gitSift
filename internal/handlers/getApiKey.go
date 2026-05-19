package handlers

import (
	"fmt"
	"os"
)

func GetApiKeyHandler(args []string) error {
	fmt.Println("Current key =", os.Getenv("AGENT_API_KEY"))
	return nil
}
