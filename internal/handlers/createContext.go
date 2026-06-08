package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func CreateContext(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments: need repoURL (args[0]) and repoCode (args[1])")
	}

	repoURL := args[0]
	repoCode := args[1]

	contextsDir := ".gosift_contexts"
	if err := os.MkdirAll(contextsDir, 0755); err != nil {
		return fmt.Errorf("failed to create contexts directory: %w", err)
	}

	url := strings.TrimSpace(repoURL)
	url = strings.TrimSuffix(url, ".git")

	replacer := strings.NewReplacer(
		"://", "___",
		"/", "_",
	)
	safeName := replacer.Replace(url)
	finalPath := filepath.Join(contextsDir, safeName+".txt")

	if err := os.WriteFile(finalPath, []byte(repoCode), 0644); err != nil {
		return fmt.Errorf("failed to write context cache file: %w", err)
	}

	fmt.Printf("Context successfully loaded into: %s\n", finalPath)
	return nil
}
