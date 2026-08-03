# 07. AI Prompt Library Integration Guide

This guide explains how to locate, configure, and copy-and-paste prompts from the PromptEngine Prompts Library to control AI coding assistants.

---

## 1. Locating the Prompts

The reusable prompts are located in the top-level `prompts/` directory. Each file is dedicated to a specific development workflow:

| Workflow File | Purpose |
| :--- | :--- |
| **[prompts/01-new-project.md](../prompts/01-new-project.md)** | Conducting discovery interviews and bootstrapping new projects. |
| **[prompts/02-existing-project.md](../prompts/02-existing-project.md)** | Reverse-engineering and documenting existing codebases. |
| **[prompts/03-migrate-existing-project.md](../prompts/03-migrate-existing-project.md)** | Migrating legacy codebases and importing guidelines. |
| **[prompts/04-add-feature.md](../prompts/04-add-feature.md)** | Scoping and implementing new features. |
| **[prompts/05-bug-fix.md](../prompts/05-bug-fix.md)** | Isolating, reproducing, and fixing codebase bugs. |
| **[prompts/06-refactor.md](../prompts/06-refactor.md)** | Safe refactoring and behavioral verification. |
| **[prompts/07-architecture-change.md](../prompts/07-architecture-change.md)** | Managing database schema updates and module layers changes. |
| **[prompts/08-documentation-sync.md](../prompts/08-documentation-sync.md)** | Synchronizing code changes with API, DB, or PRD specs. |
| **[prompts/09-project-review.md](../prompts/09-project-review.md)** | Requesting security, performance, or accessibility reviews. |
| **[prompts/10-release.md](../prompts/10-release.md)** | Auditing project readiness for production launch. |

---

## 2. How to Use the Prompts

All PromptEngine prompts are designed to be **AI-agnostic** and work across tools like ChatGPT, Claude, Cursor, Windsurf, or Codex.

To use a prompt:
1. **Open the prompt file** matching your current task.
2. **Review the guidelines** (Purpose, expected behavior, common mistakes) to align your expectations.
3. **Locate the block** labeled `### Copy-and-Paste Prompt`.
4. **Replace the bracketed variables** (such as `{STACK}`, `{FEATURE_DESCRIPTION}`, or `{SYMPTOMS}`) with details matching your current project.
5. **Paste the prompt** into your AI chat box or agent configuration.

---

## 3. Feeding Rules to External AI Agents

If you are using an external AI tool (like ChatGPT Web or Claude Web) that does not have access to your local files:
- Copy the contents of your local `AGENTS.md` and paste it as the very first message in the chat session. This establishes the Project Constitution.
- When asking the AI to perform a task, paste the relevant prompt from the library and attach only the specific source files or documentation markdown files affected (e.g. `docs/PRD.md` for a new feature). This keeps the context small and prevents memory decay.
