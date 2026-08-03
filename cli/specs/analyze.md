# CLI Command Specification: `analyze`

---

## 1. Purpose
Parses the git diff or specified file changes and outputs recommendations for which specifications (under `docs/`) require updates.

## 2. Use Cases
- Before staging a commit, a developer checks which documentation files need to be modified alongside their code edits.
- Hook system verification.

## 3. Inputs
- **Flags**:
  - `-c`, `--commit`: Compare against a specific git commit hash (defaults to active staging area).
  - `-b`, `--branch`: Compare against a target branch (e.g. `main`).

## 4. Outputs
- **Console Dumps**:
  - List of modified code files.
  - Matching documentation files that are affected.
  - Warning logs for code changes that lack corresponding docs updates in the staging index.

---

## 5. Interactive Behaviour
1. Evaluates git logs.
2. Displays a table listing the files modified and the docs files that should be updated.
3. Asks user: `Do you want to run sync for these documents? (y/n)`.

## 6. Non-Interactive Behaviour
1. Executes diff analysis.
2. Returns a structured JSON list of out-of-sync docs mapping types (ideal for hook scripts).

---

## 7. Expected Workflow
- Run `promptengine analyze` before staging commits.
- Read warnings list.
- Run `promptengine sync` or modify docs to clean alerts.

## 8. Error Handling
- **No Git Repository**: Returns error log and exits with code `1` if repository is not a git workspace.

## 9. Future Extensions
- Automated draft generation of specs changes inside the CLI using AST parsing.

## 10. Related PromptEngine Workflows
- **[Documentation Sync Workflow](../../workflows/01-feature-implementation.md)**.
