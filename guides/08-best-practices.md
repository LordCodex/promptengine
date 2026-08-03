# 08. Developer Best Practices Guide

This guide outlines advanced tips for developers to optimize their interactions with AI coding tools, minimize token consumption, prevent context window saturation, and coordinate documentation.

---

## 1. Context Optimization & Token Management

AI models charge tokens for both input (context) and output. If you feed the AI too many files, it degrades execution quality and increases latency.

- **Load Minimal Context**: Do not feed the entire repository into a chat thread. Identify the affected modules and select only those files for the prompt.
- **Automated Context Minimization**: Use the future CLI command `promptengine context --task <task>` (specified in [context command spec](../cli/specs/context.md)) to instantly resolve the minimal file path dependencies mapping to the current issue, keeping the context window clean.
- **Reference Playbooks, Don't Copy**: AI agents in IDEs (like Cursor or Windsurf) can index PromptEngine. Use `@` symbols or path pointers to link to playbooks (e.g. `@Universal Coding Standards`) rather than copying rules text into prompts.
- **Wipe Chat Thread History regularly**: When a feature implementation or bug fix is complete and checked in, start a brand-new chat thread. This clears conversational clutter and resets the context window, avoiding hallucinations.

---

## 2. Managing Project Specifications (When to Update)

- **Database Changes**: Update `docs/Database.md` immediately in the same commit as your migration script. Never defer this, or future AI iterations will generate code using stale schemas.
- **API Payloads**: Update `docs/API.md` when adding fields to requests/responses or altering endpoint routes.
- **Business Logic Rules**: Update `docs/BusinessRules.md` immediately when validation rules, equations, status codes, or calculations change.
- **Decisions & ADRs**: Record a decision in `docs/Decisions.md` (ADR log) whenever you choose between major technical paths (e.g. Postgres vs. MySQL, Sanctum vs. Custom JWT, VPS vs. AWS).
- **Regenerating AGENTS.md**: Regenerate or manually update Section 2 of `AGENTS.md` ONLY when the primary technology stack, caching layers, authentication approach, or approved exceptions change.

---

## 3. Tool-Specific Integration Patterns

Here is how to set up PromptEngine with leading AI coding assistants:

### Cursor
1. Add PromptEngine to your `.cursorrules` file or workspace config.
2. Configure Cursor’s system instructions to point to `AGENTS.md` and `promptengine/ai/bootstrap.md`.
3. Use the `@` symbol to load specifications (e.g. `@API.md` or `@Database.md`) when planning changes.

### Claude Code
1. Start `claude` in your project root.
2. In your initial prompt, instruct Claude: *"Inspect AGENTS.md in the project root and align context with its stack definitions."*
3. Rely on Claude's CLI to scan folders based on manifest instructions.

### Windsurf
1. Add the path to `AGENTS.md` and PromptEngine bootstrap rules into Windsurf's global or directory-level system instructions.
2. Use Windsurf's interactive terminal to run automated check scripts after generation.

### Codex CLI
1. Use Codex commands to query the index manifest (`playbook-manifest.json`).
2. Generate migration steps by piping output logs to the CLI.

### ChatGPT (Web Interface)
1. Since the web interface lacks filesystem access, copy the generated `AGENTS.md` and paste it into the first message.
2. Drag and drop only the target code files and docs (e.g. `docs/PRD.md`) into the chat window when requesting edits.
