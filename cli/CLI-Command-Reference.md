# CLI Command Reference (CLI-Command-Reference.md)

This document provides the command reference lookup table for all commands in the PromptEngine CLI binary.

---

## 1. Global Command Mapping Table

| Command | Action Description | Core Flags | Related Workflows | CLI Version |
| :--- | :--- | :--- | :--- | :--- |
| **`init`** | Greenfield project bootstrapping | `--stack`, `--non-interactive` | NewProjectBootstrap | v0.1 |
| **`migrate`** | Adopts PromptEngine in legacy codebases | `--legacy-rules`, `--force` | ExistingProjectBootstrap | v1.0 |
| **`doctor`** | Validates link and manifest integrity | `--fix`, `--verbose` | Audit Standard | v0.2 |
| **`analyze`** | Traces changed files and suggests spec changes | `--commit`, `--branch` | Refactoring/Git standard | v0.3 |
| **`scan`** | Disk audits for schemas/dependencies | `--output`, `--path` | ExistingProjectBootstrap | v0.1 |
| **`review`** | Audits security, performance, & a11y standards | `--severity`, `--target` | Code Review Standard | v0.3 |
| **`sync`** | Reconciles Progress and checklists logs | `--commit-msg` | Daily development workflow | v0.3 |
| **`update`** | Migrates templates and playbooks versions | `--dry-run`, `--all` | Upgrade standard | v1.0 |
| **`context`** | Resolves minimal context array mappings | `--task`, `--output-format` | AI Agent Engineering | v0.2 |
| **`prompt`** | Renders variables-injected prompt templates | `--prompt-id`, `--vars` | Reusable Prompt Library | v0.2 |
| **`docs`** | Generates blank templates under `docs/` | `--doc-name`, `--force` | Project Documentation | v0.2 |
| **`generate`** | Stub code/schema mock builder | `--type`, `--output` | API / Database design | v0.3 |
| **`config`** | Read/write global/local parameters values | `--get`, `--set`, `--global` | Config Standard | v0.2 |
| **`install`** | Pre-configures specific stacks playbooks | `--plugin-id` | Stacks setup | v1.0 |
| **`health`** | Computes diagnostic scores for workspace | `--json`, `--output` | Health Engine | v0.3 |
| **`hooks`** | Installs git client hooks triggers | `--uninstall`, `--ci` | Git integration | v1.0 |
| **`plugins`** | List/install external plugins adapters | `--list`, `--install` | Plugins system | v1.0 |
| **`version`** | Print PromptEngine CLI build data | `--short` | System utility | v0.1 |
| **`help`** | Displays console help screens layout | `[command]` | System utility | v0.1 |

---

## 2. Global Syntax Flags

Every CLI command supports these standard parameters:
- `-h`, `--help`: Show specific command parameters descriptions.
- `-v`, `--verbose`: Enable debug level output log dumps.
- `-q`, `--quiet`: Run silently, suppressing spinners and normal output logs (ideal for scripts integrations).
- `-n`, `--non-interactive`: Run in scripting mode, bypass TTY prompts, and use defaults.
- `--json`: Format command outputs as machine-readable JSON.

---

## 3. Shell Exit Codes

The binary exits with standard POSIX error codes:
* `0`: Success (All checks, migrations, or scans succeeded).
* `1`: General Error (Invalid arguments, missing files).
* `2`: Lint/Format Warning (In non-interactive mode, indicates link warning or documentation drift warning).
* `3`: Health Score Block (Health score falls below block limits, in strict validation runs).
* `127`: command not found.
