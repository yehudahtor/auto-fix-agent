# 🛠️ Auto Fix Agent for Kubescape

This agent automates triaging and proposing fixes for GitHub issues in the [Kubescape](https://github.com/kubescape/kubescape) project.

---

## 🔍 What It Does

Given a GitHub issue number, this CLI tool:

1. Retrieves the issue title and body using the GitHub API.
2. Uses OpenAI (GPT-4o-mini) to extract relevant keywords.
3. Searches through the Kubescape codebase for relevant `.go`, `.yaml`, `.yml`, and `.md` files.
4. Proposes a fix using LLM-based code analysis and summarization (early stage).
5. (Planned) Asks for human approval before modifying the code and posting a comment.

---

## 🧱 Project Structure

```
auto-fix-agent/
├── cmd/autoagent/              # CLI entrypoint
│   └── main.go
├── pkg/
│   ├── analysis/               # File scanning and repo management
│   │   └── code_search.go
│   ├── githubclient/           # GitHub API client (TBD)
│   │   └── github.go
│   └── llmclient/              # LLM keyword extraction + fix generation
│       ├── keyword_extraction.go
│       └── llmclient.go (TBD)
├── go.mod
├── go.sum
└── README.md
```

---

## 🚀 How to Run

### ✅ Prerequisites

- Go 1.20+
- [OpenAI API Key](https://platform.openai.com/account/api-keys)
- [GitHub Token](https://github.com/settings/tokens) with `repo` and `read:org` scopes

### ⚖️ Setup

```bash
export OPENAI_API_KEY=your_openai_key_here
export GITHUB_TOKEN=your_github_token_here
```

### 🧲 Run the Agent

```bash
go run cmd/autoagent/main.go <issue_number>
```

If no local Kubescape repo exists, it will automatically clone it into:

```bash
/Users/yehudahtor/Projects/kubescape
```

Override this path by setting:

```bash
export KUBESCAPE_PATH=/your/path/here
```

---

## ✨ Example Output

```bash
Issue Title: Kubescape operator does not scan roles
Issue Body: ...
🔍 LLM-Extracted Keywords:
[role, rolebinding, helm, operator, scan, control, workload]
Relevant files found:
 - pkg/...
💡 Suggested Fix:
<Proposed patch by LLM>
```

---

## 📌 Roadmap

- [x] Keyword extraction via OpenAI
- [x] File search based on content matching
- [ ] Diff patching and code modifications
- [ ] Unit test validation
- [ ] GitHub issue auto-comment
- [ ] Continuous mode (cron/CI)

---

## 🧠 Credits

Built by the Kubescape Engineering Team @ [Armo](https://www.armosec.io/).
