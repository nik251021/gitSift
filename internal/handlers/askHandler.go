package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getCacheFilePath(repoURL string) string {
	url := strings.TrimSpace(repoURL)
	url = strings.TrimSuffix(url, ".git")

	replacer := strings.NewReplacer(
		"://", "___",
		"/", "_",
	)
	safeName := replacer.Replace(url)
	return filepath.Join(".gosift_contexts", safeName+".txt")
}

func AskHandler(args []string) error {
	if len(args) < 1 {
		return errors.New("no question provided for AI")
	}

	apiKey := os.Getenv("AGENT_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("AGENT_API_KEY is not set in .env file")
	}

	question := strings.Join(args, " ")
	repoURL := os.Getenv("CURRENT_REPO_LINK")
	if repoURL == "" {
		return fmt.Errorf("CURRENT_REPO_LINK is not set in .env file or environment")
	}

	repoCodePath := getCacheFilePath(repoURL)

	fmt.Printf("Question: %s\n", question)

	if _, err := os.Stat(repoCodePath); os.IsNotExist(err) {
		fmt.Printf("Context file not found. Auto-running Sift for: %s\n", repoURL)
		siftArgs := []string{repoURL}
		if err := SiftHandler(siftArgs); err != nil {
			return fmt.Errorf("auto-sift failed: %w", err)
		}
	}

	fmt.Printf("Loading context from file: %s\n", repoCodePath)
	repoContent, err := os.ReadFile(repoCodePath)
	if err != nil {
		return fmt.Errorf("failed to read repository cache file: %w", err)
	}

	systemPrompt := "You are a code analyzer. Scan the provided repository code and answer the user's question precisely, short and without water.\n\n"
	fullPrompt := fmt.Sprintf("%s=== REPOSITORY CODE ===\n%s\n=====================\n\nUser Question: %s",
		systemPrompt, string(repoContent), question)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s", apiKey)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: fullPrompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to serialize JSON request: %w", err)
	}

	var resp *http.Response
	var body []byte

	for i := 1; i <= 3; i++ {
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("HTTP request to Gemini failed: %w", err)
		}

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read API response body: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		if resp.StatusCode == http.StatusServiceUnavailable && i < 3 {
			fmt.Printf("Gemini is overloaded (503). Retrying in 3 seconds... (Attempt %d/3)\n", i)
			time.Sleep(3 * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests && i < 3 {
			fmt.Printf("Quota exceeded (429). Waiting 20 seconds for rate-limit reset... (Attempt %d/3)\n", i)
			time.Sleep(20 * time.Second)
			continue
		}

		return fmt.Errorf("Gemini API returned status code %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return errors.New("Gemini returned an empty response")
	}

	aiResponse := geminiResp.Candidates[0].Content.Parts[0].Text

	fmt.Println("\n--- Gemini Response ---")
	fmt.Println(aiResponse)
	fmt.Println("-----------------------")

	return nil
}
