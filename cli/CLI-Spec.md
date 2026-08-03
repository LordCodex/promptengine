# CLI Specification (CLI-Spec.md)

---

## 1. Core Principles & Philosophy

The PromptEngine CLI is designed as a **thin orchestration layer**. It contains minimal unique business logic; instead, it exposes, automates, and validates the directory standards, playbooks, prompts, and workflows defined in the PromptEngine repository.

### Rules of Engagement
- **Read-First Principle**: All commands must default to performing a dry-run/read audit before modifying any files.
- **Strict Separation of Concerns**: Coding generation or project refactoring is performed by the AI coding assistant in the IDE, not by the CLI binary. The CLI's job is context assembling, health auditing, change analysis, and configuration routing.
- **Standard Exit Signaling**: Standard Unix exit codes must be followed strictly to allow pipeline automation.

---

## 2. Console User Experience (UX)

To create a premium developer experience, the CLI must adhere to these output standards:

### Typography and Layout
- **Visual Dividers**: Use consistent ASCII lines (`---` or `===`) to group logical operations.
- **Icons & Status Indicators**:
  - `✔` (Green) for successful validations, creations, or health scores.
  - `✖` (Red) for blocks, errors, or critical failures.
  - `⚠` (Yellow) for warnings or sync alerts.
  - `ℹ` (Blue) for context logs or configurations.
- **Spinners & Progress Indicators**: Long-running scans (like dependency searches or directory deep audits) must display active text spinners (e.g. `⠋ Auditing database schemas...`) rather than throwing raw logs.

---

## 3. Execution Modes

The CLI supports two primary operational formats:

### Interactive Mode
Exposed when run in a standard TTY terminal environment. 
- Prompts the user with structured multiple-choice dropdowns, checkbox arrays, and validated inputs (e.g. confirming stack versions).
- Provides instant context formatting (showing what files will change before hitting write).
- Displays colorized highlights (ANSI colors).

### Non-Interactive Mode (CI/CD)
Triggered programmatically when no TTY is detected or when specific flags (e.g. `--non-interactive` or `-n`) are set.
- Automatically selects framework defaults or parses parameters from configs.
- Disables spinners and colored fonts (monochrome output).
- Returns errors to `stderr` and prints raw diagnostic lists.
- Ideal for git pre-commit hooks and pipeline runners.
