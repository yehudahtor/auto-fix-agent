package llmclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func ProposeFix(issueText string, fileContent string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	prompt := fmt.Sprintf(`You are an expert Go developer and Kubernetes contributor working on Kubescape. 
A user opened the following GitHub issue:\n\n%s\n\n
You are reviewing this related file:\n\n%s\n\n
Suggest a change that could resolve the issue. 
Respond ONLY with a proposed fix or code snippet and a short explanation.`, issueText, fileContent[:min(len(fileContent), 3000)]) // keep token usage reasonable

	body := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "You help maintain open source Go security tooling. Be precise and concise."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
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

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
