# CLI Architecture (CLI-Architecture.md)

This document specifies the internal module architecture of the PromptEngine CLI, detailing its analytical engines, verification loops, and provider adapters.

---

## 1. System Layers

The CLI is structured into three execution layers:

```text
┌────────────────────────────────────────────────────────┐
│ 1. CLI Presentation Layer (Command Router, Console UX) │
└───────────────────────────┬────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│ 2. Core Analysis Engines (Detection, Health, Changes)  │
└───────────────────────────┬────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│ 3. Provider Adaptor Layer (GitHub, Claude, Cursor)     │
└────────────────────────────────────────────────────────┘
```

1. **CLI Presentation Layer**: Handles argument parsing, interactive questions routing, colored terminal logs, and output streams.
2. **Core Analysis Engines**: Programmatic analyzers that read files on disk, compute metrics, map changes, and resolve playbooks.
3. **Provider Adaptor Layer**: Abstractions to export prompts or context to Cursor, Windsurf, Claude Code, or plain Markdown files.

---

## 2. Analytical Engines

### Project Detection Engine
- **Strategy**: Scans the current directory on startup to reverse-engineer project dependencies.
- **Detections**:
  - *Languages*: scans extensions (`.py`, `.php`, `.ts`, `.dart`).
  - *Frameworks*: parses package manifests (e.g. `composer.json` for Laravel, `package.json` for Next.js/Nuxt).
  - *Databases*: checks config files, docker-compose models, and drivers dependencies.
  - *Queue & Caching*: identifies drivers configurations (e.g. `redis`, `memcached`, `sqs`).
  - *Testing*: maps test folders and packages (e.g., Pest, Jest, Vitest, JUnit).

### Documentation Health Engine
- **Methodology**: Evaluates the state of specifications in `docs/` and `AGENTS.md`.
- **Health Checks**:
  - Missing file check (ensuring all 10 core specs exist).
  - Links target check (validates relative link integrity).
  - Stale spec check (flags documentation that has not been edited in 30 days).
  - ADR decision review (validates that decisions are recorded in `docs/Decisions.md`).
- **Scoring Formulas**:
  - Let $D$ be the percentage of existing core files ($0-100$).
  - Let $L$ be the percentage of unbroken relative links ($0-100$).
  - Let $S$ be the percentage of synchronized files (no change flags).
  - $\text{Health Score} = 0.4D + 0.3L + 0.3S$.
  - Output rating: $90-100$ (A - Excellent), $70-89$ (B - Good), $50-69$ (C - Warning), $<50$ (F - Critical).

### Change Analysis Engine
- **Methodology**: Run git diff between active HEAD and target commit to detect changes requiring documentation updates.
- **Mapping Matrix**:
  - *Database migrations* -> Updates needed in `docs/Database.md`.
  - *API route changes* -> Updates needed in `docs/API.md`.
  - *Env variables changes* -> Updates needed in `docs/Deployment.md` and `.env.example`.
  - *Auth/Permission additions* -> Updates needed in `docs/BusinessRules.md`.
  - *New vendor package addition* -> Updates needed in `docs/Architecture.md` and `AGENTS.md`.

### Context Builder
- **Methodology**: Reduces token consumption by selecting only the minimal files required for a task.
- **Context Selection Mappings**:
  - `New Feature`: `AGENTS.md` + `docs/PRD.md` + `docs/Database.md` (if database affected) + relevant stack playbook.
  - `Bug Fix`: `AGENTS.md` + `docs/Troubleshooting.md` + target source file + stack playbook.
  - `Refactor`: `AGENTS.md` + `docs/Architecture.md` + target file.
  - `Architecture Change`: `AGENTS.md` + `docs/Decisions.md` + `docs/Database.md` + `docs/Architecture.md`.

---

## 3. Review Engine

- **Methodology**: Evaluates code quality and architectural compliance by executing local static analysis rules.
- **Scoring Categories**:
  - *Security*: Check for parameterized SQL strings, authentication policies, and secret exposures.
  - *Performance*: Check for loops containing database loads.
  - *Maintainability*: Check for separation of concerns compliance (e.g. controllers must remain thin).
  - *PromptEngine Compliance*: Check if files follow naming conventions.

---

## 4. Hook System

Exposes git automation to prevent documentation decay without blocking developers.
- **Git pre-commit Hook**: Audits the staging area. If code has schema changes but no docs updates are staged, it prints a warning: `[PromptEngine Warning] Staged database migration found, but docs/Database.md is not modified.` (Warning only by default, does not block commits unless `--strict` flag is set).
- **CI/CD Hook**: Executed during pipeline runs to generate build reports and verify health scores before deployment.

---

## 5. AI Provider Abstraction Layer

The CLI must remain provider-agnostic. It abstracts the prompt output format into standard interfaces:
- **Clipboard Output**: Standard text block for copy-and-pasting into ChatGPT Web, Claude Web, or Gemini.
- **System Rules Export**: Outputs configurations targeting IDE systems (`.cursorrules` for Cursor, system instructions for Windsurf, or config JSONs for Claude CLI).
