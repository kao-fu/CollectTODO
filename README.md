# TODO2MD

A GitHub Action to automatically convert TODO comments in your codebase into a Markdown TODO list.

---

## Usage

Add this Action to your workflow:

```yaml
- name: Run TODO2MD
  uses: kao-fu/TODO2MD@v1
  with:
    # Add your inputs here
    # example: path: './src'
```

## Inputs

| Name          | Description                                                                | Required | Default         |
| ------------- | -------------------------------------------------------------------------- | -------- | --------------- |
| path          | Path to search for TODO comments                                           | false    | .               |
| output        | Output markdown file path                                                  | false    | TODO_SUMMARY.md |
| update-readme | Insert TODO summary into README.md (between <!-- TODO_SUMMARY --> markers) | false    | false           |

In your codebase, lines with pattern of `TODO[tag]: Description` will be collected into a `.md` file. The output file can be customized using the `output` input.

Common tags are like:

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

## Outputs

A markdown file (default: `TODO_SUMMARY.md`, customizable via the `output` input)

```markdown
<!-- TODO_SUMMARY -->

# TODO Summary

## todo

- **2025-06-19** (foo:25, api/foo): Add endpoint to get entries/entry by user ID
- **2025-06-19** (foo:59, db/foo): Update the entries in the database - PayAccount
- **2025-06-19** (foo:60, db/foo): Update the entries in the database - ReceiveAccount

## research

- **2025-06-19** (foo_test:11, db/foo_test): Check test assertions and error handling
- **2025-06-19** (foo_test:68, db/foo_test): Detailed checks for the transaction status and other fields
- **2025-06-19** (foo_test:71, db/foo_test): Check the updated balances after the transaction is written

<!-- TODO_SUMMARY -->
```

## Example Workflow

```yaml
name: Generate/Update TODO List
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run TODO2MD
        uses: kao-fu/TODO2MD@v1
        with:
          path: "./src"
          output: "TODO.md" # Customize output file name
          update-readme: true # Optionally update README.md with summary
```

## fake_project Directory

The `fake_project` directory contains example Go code used for demonstration and testing purposes. It helps illustrate how the TODO2MD action collects TODO comments from a sample codebase. You can use or modify this directory to experiment with the action or to test its features.

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open issues or pull requests for improvements.
