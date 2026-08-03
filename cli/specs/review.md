# CLI Command Specification: `review`

---

## 1. Purpose
Performs static analysis of codebase components against PromptEngine standards (Security, Performance, Accessibility, and Maintainability) and outputs score parameters.

## 2. Use Cases
- Developer checks their local code modifications before committing to verify standard compliance.
- Pull request validation pipeline runs.

## 3. Inputs
- **Arguments**:
  - `[target_path]`: Specific file path or directory folder to review (defaults to current staging changes).
- **Flags**:
  - `-s`, `--severity`: Filter issues by severity (`block`, `important`, `suggestion`).
  - `-o`, `--output`: Save markdown review report output path.

## 4. Outputs
- **Console Dumps**:
  - Status log showing code checks executed.
  - Bullet tables showing issues grouped by severity with references to universal/stack playbooks.
  - Compliance Score (e.g. `Security Compliance: 85%`).

---

## 5. Interactive Behaviour
1. Executes static audits showing diagnostic progress.
2. Lists findings in a scrollable terminal interface.
3. For each finding, asks user: `Would you like to open the corresponding PromptEngine playbook link? (y/n)`.

## 6. Non-Interactive Behaviour
1. Executes audits.
2. Returns findings mapping logs directly.
3. Exits with code `3` if any severity `block` finding is detected.

---

## 7. Expected Workflow
- Run `promptengine review app/Actions/RegisterUser.php` in console.
- Fix any `[BLOCK]` findings.
- Re-run check to verify.

## 8. Error Handling
- **AST Parse failure**: If codebase syntax is invalid (fails compiling), prints exception details and exits with code `1`.

## 9. Future Extensions
- Integrations with local LLM APIs to perform semantic code analysis instead of purely regex-based static analysis.

## 10. Related PromptEngine Workflows
- **[Code Review Prompt](../../prompts/09-project-review.md)**.
