package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	apiKey := os.Getenv("AGENT_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("AGENT_API_KEY is not set in .env file")
	}

	fmt.Println("sending code to gemini for analyse")

	systemPrompt := "Ты — анализатор кода. Проанализируй этот слепок репозитория. " +
		"Составь подробную, но сжатую техническую справку для кэша. Без лишних слов. Добавь в конце примеры кода для самых важных функций по твоему усмотрению.\n\n"

	fullPrompt := systemPrompt + repoCode

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
		return fmt.Errorf("error in serilisation of json: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error in http-request to gemini: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error in reading answer of API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Gemini API returned status code %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return fmt.Errorf("error in parsing of JSON answer: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("Gemini returned nothing")
	}

	aiResponse := geminiResp.Candidates[0].Content.Parts[0].Text

	contextsDir := ".gosift_contexts"
	if err := os.MkdirAll(contextsDir, 0755); err != nil {
		return fmt.Errorf("we cant create folder for context...: %w", err)
	}

	safeFileName := strings.ReplaceAll(repoURL, "://", "___")
	safeFileName = strings.ReplaceAll(safeFileName, "/", "_")
	safeFileName = safeFileName + ".txt"

	finalPath := contextsDir + "/" + safeFileName

	err = os.WriteFile(finalPath, []byte(aiResponse), 0644)
	if err != nil {
		return fmt.Errorf("cant write to file of context: %w", err)
	}

	fmt.Printf("context was succesfully loaded in: %s\n", finalPath)
	return nil
}
