# CLI Command Specification: `sync`

---

## 1. Purpose
Programmatically synchronizes the development checklist in `docs/Progress.md` based on active branch commits or staged files.

## 2. Use Cases
- Auto-checking completed task checkboxes in `docs/Progress.md` when a developer finishes implementing a code path.
- Generating progress statistics for release logs.

## 3. Inputs
- **Flags**:
  - `-m`, `--message`: Associate sync update with a custom message.
  - `--dry-run`: View check changes without writing to `docs/Progress.md`.

## 4. Outputs
- **Files Modified**:
  - `docs/Progress.md` (updates task list statuses).
- **Console Dumps**:
  - Diff representation of checks marked completed/added.

---

## 5. Interactive Behaviour
1. Scans code files and checklist items.
2. Identifies matching file paths and checklist items (e.g. check `- [ ] Implement UserController` maps to `app/Http/Controllers/UserController.php` existence).
3. Prompts the user: `Staged UserController found on disk. Mark check 'Implement UserController' as completed in Progress.md? (y/n)`.

## 6. Non-Interactive Behaviour
1. Compares checklist paths against git index data.
2. Automatically marks items as completed if matching files exist.
3. Overwrites files silently.

---

## 7. Expected Workflow
- Staged code changes.
- Run `promptengine sync`.
- Docs update automatically.

## 8. Error Handling
- **Progress File Missing**: Returns error log and exits code `1` if `docs/Progress.md` cannot be located.

## 9. Future Extensions
- Automated synchronization with external ticketing software (such as Jira or GitHub Issues API).

## 10. Related PromptEngine Workflows
- **[Documentation Sync Prompt](../../prompts/08-documentation-sync.md)**.
