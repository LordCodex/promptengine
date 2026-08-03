# CLI Command Specification: `context`

---

## 1. Purpose
Determines and outputs the minimum set of PromptEngine playbooks and project specifications required for a target task to minimize token usage and prevent AI context loss.

## 2. Use Cases
- A developer or agent shell wrapper queries which files are required for implementing a new feature in order to feed only those to the LLM.

## 3. Inputs
- **Flags**:
  - `-t`, `--task`: Target task type (`new-feature`, `bug-fix`, `review`, `refactor`, `migrate`, `deployment`, `architecture`).
  - `-o`, `--output-format`: Output configuration format (`json`, `text`, `cursor`).

## 4. Outputs
- **Console Dumps**:
  - List of absolute file paths to load.
  - Estimated token size block.

---

## 5. Interactive Behaviour
1. Prompts user: `Select active task:`.
2. Prompts: `Which codebase components are affected? (e.g. app/Models/User.php)`.
3. Computes dependencies and prints the list of files.
4. Asks: `Export paths list to clipboard? (y/n)`.

## 6. Non-Interactive Behaviour
1. Resolves paths based on task mapping configuration rules.
2. Returns space-separated absolute paths list (ideal for piping to CLI tools).

---

## 7. Expected Workflow
- Run `promptengine context --task new-feature` in terminal.
- CLI outputs minimal paths list.
- Pipe list directly to your IDE assistant (e.g. `cursor --read-context $(promptengine context --task new-feature)`).

## 8. Error Handling
- **Invalid Task Flag**: Exits with code `1` and prints list of valid tasks.

## 9. Future Extensions
- Automated scanning of modified file imports to dynamically add dependent files.

## 10. Related PromptEngine Workflows
- **[Context Builder Specification](../CLI-Architecture.md)**.
