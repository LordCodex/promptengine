# CLI Command Specification: `init`

---

## 1. Purpose
Initializes a new PromptEngine workspace in an empty directory or brand-new project. It generates `AGENTS.md` (AI Constitution) and placeholders/templates for the 10 core documents under `docs/`.

## 2. Use Cases
- A developer starts a brand-new greenfield repository and wants to structure it according to PromptEngine documentation rules.
- Onboarding PromptEngine in an empty workspace folder before coding begins.

## 3. Inputs
- **Arguments**:
  - `[project_name]`: Optional name of the target application.
- **Flags**:
  - `-s`, `--stack`: Pre-configured stacks selection (e.g. `php-laravel`, `react-next`, `dart-flutter`).
  - `-d`, `--docs-dir`: Customize documentation directory path (defaults to `docs/`).
  - `-n`, `--non-interactive`: Run without prompting, using standard template defaults.

## 4. Outputs
- **Files Created**:
  - `[project-root]/AGENTS.md` (AI Constitution populated with Section 1 and template Section 2).
  - `[project-root]/[docs-dir]/` folder containing the 10 core markdown specification files (`PRD.md`, `Architecture.md`, etc.).
  - `[project-root]/.promptengine.json` (Local project config file).
- **Console Dumps**: Status log showing files generated and relative linkages markdown check.

---

## 5. Interactive Behaviour
1. Prompts the user: `Enter project name:`.
2. Prompts with a selection list: `Select core technology stack:`.
3. Prompts: `Where should specifications reside? (default: docs/)`.
4. Executes discovery questionnaire (brief grouping of high-impact questions like authentication style and database type).
5. Displays a confirmation message showing paths to be created, asking: `Proceed with generation? (y/n)`.

## 6. Non-Interactive Behaviour
1. Automatically assigns standard defaults if flags are omitted.
2. Creates `.promptengine.json`, `AGENTS.md`, and `docs/` directly without questioning.
3. Suppresses confirmation dialogs.

---

## 7. Expected Workflow
- Run `promptengine init my-project --stack php-laravel` in your shell.
- Complete the 3-question prompt loop.
- The CLI constructs directories and fills out files.
- Human review is conducted before checking specs into git.

## 8. Error Handling
- **Directory Not Empty**: Warns the developer if files already exist. If `--force` is set, overwrites files. Otherwise, exits with code `1`.
- **Invalid Stack selection**: If `--stack` maps to a stack not present in standard stack playbooks, prints warning and lists approved stack playbooks.

## 9. Future Extensions
- Support for generating scaffold code repositories based on popular framework generators directly (e.g. `npx create-next-app`).

## 10. Related PromptEngine Workflows
- **[New Project Bootstrap Workflow](../../workflows/NewProjectBootstrap.md)**.
