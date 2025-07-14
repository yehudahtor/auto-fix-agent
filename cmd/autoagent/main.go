package main

import (
	"fmt"
	"log"
	"os"
	"github.com/armo/auto-fix-agent/pkg/githubclient"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: autoagent <issue_number>")
	}

	issueNum := os.Args[1]
	ctx := githubclient.NewContext()
	client := githubclient.NewClient(ctx)

	title, body := client.GetIssueDetails(issueNum)
	fmt.Println("Issue Title:", title)
	fmt.Println("Issue Body:", body)
}
