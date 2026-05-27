package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetRepoFiles(args []string) (string, error) {
	if len(args) < 1 || args[0] == "" {
		return "", fmt.Errorf("no repository URL provided")
	}
	repoURL := args[0]

	fmt.Println("Getting files. . .")

	tmpDir, err := os.MkdirTemp("", "gosift-")
	if err != nil {
		return "", fmt.Errorf("Error in creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cloneCmd := exec.Command("git", "clone", "--depth", "1", repoURL, tmpDir)
	if err := cloneCmd.Run(); err != nil {
		return "", fmt.Errorf("error in running git clone for %s: %w", repoURL, err)
	}
	fmt.Println("File was downloaded, collecting files")

	listCmd := exec.Command("git", "ls-files")
	listCmd.Dir = tmpDir

	out, err := listCmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error in getting files: %w", err)
	}

	files := strings.Split(string(out), "\n")
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("=== string of repo: %s ===\n\n", repoURL))

	for _, filePath := range files {
		filePath = strings.TrimSpace(filePath)
		if filePath == "" {
			continue
		}

		if filePath == ".env" || strings.HasSuffix(filePath, ".sum") {
			continue
		}

		fullPath := tmpDir + "/" + filePath

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		builder.WriteString(fmt.Sprintf("--- File: %s ---\n", filePath))
		builder.WriteString("```\n")
		builder.Write(content)
		builder.WriteString("\n```\n\n")
	}

	return builder.String(), nil
}
