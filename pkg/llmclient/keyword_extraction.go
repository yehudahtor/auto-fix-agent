// pkg/llmclient/keyword_extraction.go
package llmclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func ExtractKeywords(text string) ([]string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	prompt := fmt.Sprintf(`Extract 5-10 high-signal keywords (no explanations) from the following GitHub issue text:

%s`, text)

	body := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "system", "content": "You extract keywords from GitHub issues for a code scanning agent."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// DEBUG: Log full response for analysis
	log.Printf("🛰️ OpenAI API Status: %s", resp.Status)

	rawResp, _ := io.ReadAll(resp.Body)
	log.Printf("📦 Raw OpenAI Response: %s", rawResp)

	// Decode again after reading manually
	resp.Body = io.NopCloser(bytes.NewBuffer(rawResp))

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		log.Println("⚠️ No response choices returned by OpenAI")
		return fallbackKeywords(text), nil
	}

	raw := response.Choices[0].Message.Content
	log.Printf("🧠 Raw GPT Keyword Output: %s", raw)

	// Basic tokenization and cleanup
	lines := strings.Split(raw, "\n")
	var keywords []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "• ")
		line = strings.TrimPrefix(line, "* ")
		if line != "" {
			keywords = append(keywords, line)
		}
	}

	if len(keywords) == 0 {
		log.Println("⚠️ GPT response had no usable keywords, falling back")
		return fallbackKeywords(text), nil
	}

	return keywords, nil
}

func fallbackKeywords(text string) []string {
	log.Println("✅ Using fallback keyword extraction")
	text = strings.ToLower(text)
	possible := []string{"role", "rolebinding", "helm", "operator", "scan", "control", "workload", "cluster"}
	var results []string
	for _, word := range possible {
		if strings.Contains(text, word) {
			results = append(results, word)
		}
	}
	return results
}
