# 04. Add Feature Prompt

---

## Purpose
Instructs the AI to design and plan the implementation of a new software feature, write the code and tests, and update specifications.

## When to use
Use when you want to add a new business workflow, user flow, or logic module.

## Example
Adding a user checkout and billing flow.

---

## Copy-and-Paste Prompt

```markdown
I want you to implement a new feature in our codebase.

Feature Description: {FEATURE_DESCRIPTION}
Target stack: {STACK}
Target paths: {PATHS}

Before writing any code:
1. Read the project constitution (`AGENTS.md`) and the relevant files under `docs/` (such as `docs/PRD.md`, `docs/Database.md`, or `docs/API.md`).
2. Write a detailed implementation plan in your thinking block outlining:
   - Affected database schemas or route configurations.
   - Proposed classes, actions, or service files.
   - Pinned dependencies to use (no dynamic additions without approval).
   - Automated Pest/Vitest test assertions to cover edge and failure paths.
3. Wait for my explicit approval before modifying the codebase files.
4. Once approved, implement the files, verify that automated tests pass, and update the corresponding documentation files (under `docs/`) and `docs/Progress.md`.
```

---

## Expected AI Behaviour
1. The AI reads `AGENTS.md` and related specs.
2. It outputs a structured plan showing code modifications, database scripts, and Pest/Vitest tests.
3. It pauses and waits for your review.
4. After approval, it writes the code and updates the docs.

## Common Mistakes
- **Skipping the planning step**: Letting the AI output 500 lines of code without plan validation. Tell it: *Stop and present the plan first.*
- **Stale documentation**: Failing to update `docs/API.md` when payload shapes are altered.
