# CollectTODO

A GitHub Action to automatically collect TODO comments in your codebase and post a Markdown TODO summary as a comment on Pull Requests.

It is basically not recommended to merge any developing branch with unresolved TODOs into main branch, but a issued TODO. For a new work, always create an Issue and a new branch.

---

## How It Works

- When you open or update a Pull Request, this Action scans your code for `TODO[tag]: Description` comments.
- It generates a categorized TODO summary in Markdown.
- The summary is posted as a comment on the PR (not pushed to any branch or file).

---

## Usage Example

Add this Action to your workflow (e.g. `.github/workflows/todo-summary.yml`):

```yaml
name: TODO-Summary
permissions:
  contents: read
  pull-requests: write # Required to update PR with TODO summary

on:
  pull_request:
    branches: ["main", "develop/*"]

jobs:
  test:
    name: TODO Summary
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.24"

      - name: Run CollectTODO
        uses: kao-fu/CollectTODO@main
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          # Optional inputs:
          root_dir: ./src
          blacklist: node_modules,.git
          whitelist: main.go,utils.go
```

### Action Inputs

| Name      | Type   | Required | Description                                                                           |
| --------- | ------ | -------- | ------------------------------------------------------------------------------------- |
| root_dir  | string | No       | Root directory to scan for TODOs. Default is the repository root.                     |
| blacklist | string | No       | Comma-separated list of base names/extensions/paths to ignore.                        |
| whitelist | string | No       | Comma-separated list of base names/extensions/paths to include (overrides blacklist). |

---

## TODO Comment Format

Use the following format in your code:

- Go as example

```go
// TODO[tag]: Description

// TODO[urgent]: Refactor this function for better readability
```

- Python as example

```python
# TODO[tag]: Description

# TODO[urgent]: Refactor this function for better readability
```

---

## Tag Examples (By ChatGPT)

### Priority-Based Tags

| Tag       | Meaning                                              |
| --------- | ---------------------------------------------------- |
| `now`     | Must be fixed immediately. Blocking.                 |
| `urgent`  | High priority, should be done soon.                  |
| `later`   | Can be delayed, not urgent.                          |
| `low`     | importance; almost optional.                         |
| `next`    | Should be picked up in the next sprint or dev cycle. |
| `blocker` | Blocking another feature or fix.                     |

### Type of Work

| Tag        | Meaning                                   |
| ---------- | ----------------------------------------- |
| `refactor` | Needs structural improvement.             |
| `fix`      | Known bug or broken behavior.             |
| `test`     | Missing or inadequate tests.              |
| `perf`     | Performance optimization needed.          |
| `doc`      | Requires documentation or comment update. |
| `review`   | Needs code review or feedback.            |
| `cleanup`  | Temporary or messy code needing polish.   |
| `hack`     | Non-ideal workaround; revisit later.      |
| `optimize` | Improve efficiency or code quality.       |
| `remove`   | Dead code or obsolete section.            |

### Process / Planning

| Tag        | Meaning                                         |
| ---------- | ----------------------------------------------- |
| `todo`     | General pending task.                           |
| `idea`     | Idea for possible improvement or feature.       |
| `research` | Needs investigation or prototyping.             |
| `discuss`  | Needs team discussion.                          |
| `block`    | Currently blocked by external/internal factors. |
| `depend`   | Depends on another component or library.        |
| `upgrade`  | Requires version update or migration.           |

### Contextual Tags

| Tag        | Meaning                               |
| ---------- | ------------------------------------- |
| `ui`       | User interface related.               |
| `api`      | API-level concern.                    |
| `db`       | Database issue or improvement.        |
| `security` | Security-related fix or audit needed. |
| `infra`    | Infrastructure or deployment issue.   |
| `sre`      | Related to site reliability / ops.    |
| `config`   | Configuration or environment-related. |

---

## Example Output (as PR Comment)

```
# TODO Summary

## refactor
- **2025-06-19** (main.go:25, foo_project/main.go): Refactor this function for better readability

## research
- **2025-06-19** (utils.go:11, foo_project/utils.go): Investigate edge case handling
```

---

## Directory Structure

- `main.go` — Main logic for scanning and generating TODO summary.
- `.github/workflows/todo-summary.yml` — Example workflow file.
- `action.yml` — Action definition.

---

## Contributing

Contributions are welcome! Please open issues or pull requests for improvements.
