# CLI Command Specification: `health`

---

## 1. Purpose
Computes a comprehensive compliance and coverage score for the project across Documentation, Architecture, Security, Testing, and Maintainability rules.

## 2. Use Cases
- Pre-release gate audit to ensure codebase health does not fall below standard limits.
- CI/CD build script verification block.

## 3. Inputs
- **Flags**:
  - `--json`: Output raw data schema format (machine readable).
  - `--output [path]`: Path to save generated health status report markdown.

## 4. Outputs
- **Console Dumps**:
  - Breakdown lists of compliance scores (e.g. `Documentation Health: A (95%)`, `Security Coverage: B (88%)`).
  - Total health score.

---

## 5. Interactive Behaviour
1. Executes complete scans showing sub-checks execution logs (real-time diagnostics progress).
2. Renders ANSI-colored score charts directly in console.
3. If scores are below targets, highlights specific actions to fix (e.g. *"docs/Database.md does not match migration app/database/migration.php"*).

## 6. Non-Interactive Behaviour
1. Quietly runs audits.
2. Returns total health integer value (e.g. `87`).
3. If score falls below target threshold config (e.g. `<70`), exits with code `3`.

---

## 7. Expected Workflow
- Run `promptengine health` in terminal.
- Review total score breakdown.
- Commit clean results report.

## 8. Error Handling
- **Missing Core playbooks**: If standard playbook templates cannot be loaded, logs warning and exits code `1`.

## 9. Future Extensions
- Historical charts generation logging score data over git commits timeline.

## 10. Related PromptEngine Workflows
- **[Documentation Health Engine Specs](../CLI-Architecture.md)**.
