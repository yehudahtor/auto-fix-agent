// pkg/analysis/code_search.go
package analysis

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var supportedExtensions = []string{".go", ".yaml", ".yml", ".md"}

func EnsureRepoExists(repoPath string) error {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Printf("📥 Cloning Kubescape repo into %s", repoPath)
		cmd := exec.Command("git", "clone", "https://github.com/kubescape/kubescape.git", repoPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}
	}
	return nil
}

func FindRelevantFilesByKeywords(keywords []string, repoPath string) ([]string, error) {
	var relevantFiles []string

	err := filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(d.Name())
		if !contains(supportedExtensions, ext) {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil // skip unreadable files
		}

		content := strings.ToLower(string(contentBytes))
		for _, keyword := range keywords {
			if strings.Contains(content, strings.ToLower(keyword)) {
				relevantFiles = append(relevantFiles, path)
				log.Printf("Keyword match: '%s' in file %s", keyword, path)
				break // no need to keep scanning once matched
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(relevantFiles) == 0 {
		log.Println("No relevant files matched the provided keywords.")
	} else {
		log.Printf("%d relevant file(s) found.", len(relevantFiles))
	}

	return relevantFiles, nil
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
