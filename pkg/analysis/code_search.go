package analysis

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FindRelevantFiles searches files with relevant extensions for keyword matches in the issue text.
func FindRelevantFiles(issueText string, baseDir string) ([]string, error) {
	var matches []string
	keywords := extractKeywords(issueText)

	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !isRelevantFile(path) {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		lowerContent := strings.ToLower(string(content))
		for _, keyword := range keywords {
			if strings.Contains(lowerContent, strings.ToLower(keyword)) {
				matches = append(matches, path)
				break
			}
		}
		return nil
	})

	return matches, err
}

func extractKeywords(text string) []string {
	words := strings.Fields(text)
	var keywords []string
	for _, word := range words {
		word = strings.Trim(word, ".,():\"") // Clean punctuation
		if len(word) > 4 {
			keywords = append(keywords, word)
		}
	}
	return keywords
}

func isRelevantFile(path string) bool {
	extensions := []string{".go", ".yaml", ".yml", ".md"}
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
