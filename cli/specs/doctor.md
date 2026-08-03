# CLI Command Specification: `doctor`

---

## 1. Purpose
Checks the repository state to locate broken relative file links, missing core documents, or manifest syntax errors, and optionally resolves them.

## 2. Use Cases
- A developer runs a health check to verify that all markdown file linkages are correct before committing.
- Automated pipeline validation runs.

## 3. Inputs
- **Flags**:
  - `--fix`: Automatically resolve minor issues (like correcting link targets or creating missing blank templates).
  - `-v`, `--verbose`: Prints detailed logs listing every evaluated link.

## 4. Outputs
- **Console Dumps**:
  - Detailed diagnostic report showing passed checks, warnings, and errors.
  - Final exit status block.

---

## 5. Interactive Behaviour
1. Scans workspace showing a real-time progress spinner.
2. If link targets are broken (e.g. `[Schema](Data.md)` but file is named `Database.md`), presents a suggestion prompt: `Link target 'Data.md' in PRD.md does not exist. Change to 'Database.md'? (y/n)`.

## 6. Non-Interactive Behaviour
1. Executes complete scan.
2. Prints error descriptions to `stderr`.
3. If errors are found, exits with code `2`.

---

## 7. Expected Workflow
- Run `promptengine doctor` in project terminal.
- Review link output tables.
- Address errors to pass pipeline constraints.

## 8. Error Handling
- **Missing AGENTS.md**: CLI returns error and halts. Exits with code `1`.

## 9. Future Extensions
- Automated lint checks validating that all headers follow hierarchy rules.

## 10. Related PromptEngine Workflows
- **[Troubleshooting Guide](../../guides/09-troubleshooting.md)**.
