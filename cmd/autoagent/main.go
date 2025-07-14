package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yehudahtor/auto-fix-agent/pkg/analysis"
	"github.com/yehudahtor/auto-fix-agent/pkg/githubclient"
	"github.com/yehudahtor/auto-fix-agent/pkg/llmclient"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: autoagent <issue_number> [path_to_kubescape_repo]")
	}

	issueNum := os.Args[1]

	defaultRepoPath := os.Getenv("KUBESCAPE_PATH")
	if defaultRepoPath == "" {
		defaultRepoPath = "/Users/yehudahtor/Projects/kubescape" // adjust this as needed
	}

	repoPath := defaultRepoPath
	if len(os.Args) >= 3 {
		repoPath = os.Args[2]
	}

	ctx := githubclient.NewContext()
	client := githubclient.NewClient(ctx)

	title, body := client.GetIssueDetails(issueNum)
	fmt.Println("Issue Title:", title)
	fmt.Println("Issue Body:", body)

	matches, err := analysis.FindRelevantFiles(title+" "+body, repoPath)
	if err != nil {
		log.Fatalf("Error searching files: %v", err)
	}

	fmt.Println("Relevant files found:")
	for _, file := range matches {
		fmt.Println(" -", file)
	}

	if len(matches) == 0 {
		fmt.Println("No relevant files found. Exiting.")
		return
	}

	// For now, just read the first matching file (in real use you'd aggregate context)
	content, err := os.ReadFile(matches[0])
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", matches[0], err)
	}

	suggestion, err := llmclient.ProposeFix(title+"\n"+body, string(content))
	if err != nil {
		log.Fatalf("LLM call failed: %v", err)
	}

	fmt.Println("\n💡 Suggested Fix:\n")
	fmt.Println(strings.TrimSpace(suggestion))
}
