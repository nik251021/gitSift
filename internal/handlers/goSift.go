package handlers

import (
	"fmt"
)

func SiftHandler(args []string) error {
	if len(args) < 1 || args[0] == "" {
		return fmt.Errorf("please provide a repository URL. Example: goSift sift <url>")
	}
	repoURL := args[0]

	repoCode, err := GetRepoFiles(args)
	if err != nil {
		return fmt.Errorf("failed to collect repository files: %w", err)
	}

	fakeArgs := []string{repoURL, repoCode}

	err = CreateContext(fakeArgs)
	if err != nil {
		return fmt.Errorf("failed to generate and save AI context: %w", err)
	}

	fmt.Println("\n🚀 [ALL DONE] project is ready!")
	return nil
}
