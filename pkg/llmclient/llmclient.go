package llmclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const endpoint = "https://api.openai.com/v1/chat/completions"

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestBody struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type choice struct {
	Message message `json:"message"`
}

type responseBody struct {
	Choices []choice `json:"choices"`
}

func ProposeFix(issueText, codeSnippet string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	prompt := fmt.Sprintf(
		"Here is a GitHub issue:\n\n%s\n\nAnd here is the relevant code:\n\n%s\n\nPlease propose a fix or patch that resolves the issue.",
		issueText, codeSnippet,
	)

	body := requestBody{
		Model: "gpt-4",
		Messages: []message{
			{Role: "system", Content: "You are an expert Go developer contributing to an open-source Kubernetes project."},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.2,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respText, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error: %s", string(respText))
	}

	var result responseBody
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return result.Choices[0].Message.Content, nil
}
