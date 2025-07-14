package githubclient

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	client *github.Client
	ctx    context.Context
}

func NewContext() context.Context {
	return context.Background()
}

func NewClient(ctx context.Context) *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN not set")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	return &GitHubClient{client: client, ctx: ctx}
}

func (g *GitHubClient) GetIssueDetails(issueNum string) (string, string) {
	owner := "kubescape"
	repo := "kubescape"
	issueNumber := parseIssueNumber(issueNum)

	issue, _, err := g.client.Issues.Get(g.ctx, owner, repo, issueNumber)
	if err != nil {
		log.Fatalf("Failed to fetch issue: %v", err)
	}
	return issue.GetTitle(), issue.GetBody()
}

func parseIssueNumber(s string) int {
	var num int
	_, err := fmt.Sscanf(s, "%d", &num)
	if err != nil {
		log.Fatalf("Invalid issue number: %v", err)
	}
	return num
}
