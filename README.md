# Armo Auto-Fix Agent

This is a prototype Go-based agent for analyzing and resolving GitHub issues in the Kubescape repo.

## Usage

```bash
export GITHUB_TOKEN=<your_personal_access_token>
go run cmd/autoagent/main.go <issue_number>
```

## Structure
- `cmd/autoagent/` — main CLI
- `pkg/githubclient/` — wraps GitHub API interactions

## Setup
- Requires Go 1.21+
- Requires a GitHub PAT with `repo` scope
