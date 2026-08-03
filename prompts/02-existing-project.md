# 02. Existing Project Bootstrap Prompt

---

## Purpose
Instructs the AI to perform a codebase audit, reverse-engineer active dependencies and patterns, and auto-generate `AGENTS.md` and the 10 core documents under `docs/`.

## When to use
Use when onboarding an AI coding assistant to an active repository that lacks PromptEngine structure or documentation.

## Example
Adopting PromptEngine in a legacy Node/Express API with a PostgreSQL database.

---

## Copy-and-Paste Prompt

```markdown
You are going to adopt PromptEngine in this existing codebase.

Follow the Existing Project Bootstrap Workflow:
1. Scan the repository directories and configuration manifests (e.g., composer.json, package.json, pubspec.yaml, docker-compose.yml) to identify active language and framework versions.
2. Scan routing files, controllers, and database schema/migrations files on disk to map our models, API routes, and architectural patterns.
3. Automatically generate the project constitution (`AGENTS.md`) in the project root using `project/templates/AGENTS.template.md`, reverse-engineering the Stack details, constraints, and exceptions from the audit.
4. If a `docs/` folder does not exist, create it and generate the 10 core specs mapping the current state of our codebase on disk. If `docs/` already exists, audit and update the outdated sections.
5. Present a report summarizing the stack versions, detected pattern, any documentation gaps, and wait for my review and validation.

Begin by scanning the directories.
```

---

## Expected AI Behaviour
1. The AI will read package and lock files.
2. It will generate `AGENTS.md` in the project root.
3. It will populate or update files in `docs/` to reflect the database schemas and routes.
4. It will provide a delta analysis of any discrepancies between documentation and code.

## Common Mistakes
- **Failing to check lock files**: The AI might make assumptions about package versions if lock files are omitted. Verify that lock files are committed and readable.
- **Speculating missing requirements**: The AI writing documentation for non-existent routes or database fields. Remind it: *Do not document what is not on disk.*
