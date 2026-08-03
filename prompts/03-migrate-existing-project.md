# 03. Migrate Existing Project Prompt

---

## Purpose
Instructs the AI to import old or fragmented developer instructions files, merge them into the PromptEngine structure, and generate `AGENTS.md` and documentation.

## When to use
Use when onboarding PromptEngine in a project that already has old instructions text files, wiki pages, or separate guidelines.

## Example
Consolidating old `instructions.txt` and `RULES.md` into `AGENTS.md` and `docs/`.

---

## Copy-and-Paste Prompt

```markdown
I want to migrate this project under PromptEngine governance.

We have legacy guidelines or old rules located in: {LEGACY_RULES_PATHS}.

Follow these steps:
1. Scan the codebase files on disk to map the current technology stack, frameworks, and architecture.
2. Read the legacy guidelines files specified above.
3. Consolidate these guidelines:
   - Move project-level stack parameters and constraints into a new project-specific root `AGENTS.md` generated from `project/templates/AGENTS.template.md`.
   - Merge business logic calculations, constants, and custom validation limits into `docs/BusinessRules.md`.
   - Merge product context files into the corresponding templates under `docs/` (e.g. `docs/PRD.md` or `docs/Architecture.md`).
4. Propose a plan to delete the old redundant guidelines files, and wait for my approval before modifying the repository.
```

---

## Expected AI Behaviour
1. The AI reads both the source files and the specified legacy text rules.
2. It generates `AGENTS.md` and maps the custom constraints (such as *"All tables must use integer cents"*) into Section 2.
3. It places old wiki-like requirements inside `docs/` specs.
4. It provides a diff list and proposal to clean up (delete) the old rules.

## Common Mistakes
- **Redundant rule retention**: Leaving the old text rules in place, which confuses the AI. Ensure you delete them once they are merged.
- **Overwriting local specifications**: overwriting files in `docs/` without merging the legacy instructions. Instruct the AI explicitly to *merge*, not replace.
