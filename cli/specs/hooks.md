# CLI Command Specification: `hooks`

---

## 1. Purpose
Installs or configures client-side Git hooks or CI/CD pipelines validation steps to automate PromptEngine compliance checking.

## 2. Use Cases
- A developer wants to run `promptengine analyze` automatically before every git commit to ensure docs remain synchronized.

## 3. Inputs
- **Flags**:
  - `--uninstall`: Remove all PromptEngine git hook bindings.
  - `--ci [provider]`: Scaffold CI runner scripts mapping to specific providers (e.g. `github-actions`, `gitlab-ci`).

## 4. Outputs
- **Files Created/Modified**:
  - Writes hooks script inside `.git/hooks/pre-commit` folder.
  - Generates pipeline config stubs (e.g. `.github/workflows/promptengine.yml`).
- **Console Dumps**: Confirmation status of hook registration actions.

---

## 5. Interactive Behaviour
1. Prompts user: `Select hooks to register:`.
2. Lists options (e.g. `Pre-commit validation`, `GitHub Actions CI workflow`).
3. Prompts for target behaviors: `Should hooks block commit on warnings? (y/n)`.
4. Writes hook configuration script and verifies file permissions.

## 6. Non-Interactive Behaviour
1. Automatically registers standard warning-only pre-commit hook mapping config settings.
2. If `--ci` flag is passed, directly copies workflow templates files.
3. Exits with code `1` if `.git` workspace folder cannot be found.

---

## 7. Expected Workflow
- Run `promptengine hooks` in project console.
- Confirm TTY setup options.
- The pre-commit check script is active.

## 8. Error Handling
- **Missing Git Index**: Returns warning log if repository has never initialized git.

## 9. Future Extensions
- Verification triggers for other version control tools (e.g. Mercurial).

## 10. Related PromptEngine Workflows
- **[Hook System Specifications](../CLI-Architecture.md)**.
