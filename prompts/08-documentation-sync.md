# 08. Documentation Sync Prompt

---

## Purpose
Instructs the AI to scan recent code modifications (e.g. git diffs or edited files) and update the corresponding markdown specifications under `docs/` and `docs/Progress.md`.

## When to use
Use before checking in a feature or bug fix, to ensure documentation is kept synchronized with implementation.

## Example
Updating the API response shape in `docs/API.md` after adding a field to a Laravel controller.

---

## Copy-and-Paste Prompt

```markdown
I have modified some source files in the codebase. I want you to synchronize our project documentation.

Modified files list: {MODIFIED_FILES}
Current Git Diff (if available): {GIT_DIFF}

Please follow these steps:
1. Scan the modified files or diff content.
2. Identify which specifications under `docs/` are affected (e.g. database schema columns, API endpoint shapes, business calculations formulas).
3. Update ONLY the affected sections of the corresponding documents (e.g. `docs/Database.md`, `docs/API.md`, `docs/BusinessRules.md`). Do not rewrite unrelated text.
4. Record the completed task list and files changed in `docs/Progress.md`.
5. Present the list of modified docs with clickable file links for my review.
```

---

## Expected AI Behaviour
1. The AI reads the diff or source files.
2. It parses the matching files under `docs/`.
3. It makes surgical modifications to only the affected blocks or tables.
4. It lists the changes made and presents relative file links.

## Common Mistakes
- **Rewriting the entire file**: The AI erasing historical logs or adjacent documentation. Command it: *Only update the affected lines/tables, preserving context.*
- **Missing Progress update**: Forgetting to check off tasks in `docs/Progress.md`.
