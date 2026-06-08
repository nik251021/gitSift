package handlers

import (
	"fmt"
	"os"
)

func GetApiKeyHandler(args []string) error {
	apiKey := os.Getenv("AGENT_API_KEY")
	if apiKey == "" {
		fmt.Println("Current AGENT_API_KEY is not set")
		return nil
	}

	fmt.Printf("Current AGENT_API_KEY = %s\n", apiKey)
	return nil
}
