# CLI Command Specification: `migrate`

---

## 1. Purpose
Migrates a legacy codebase to PromptEngine, merging existing developer guidelines and text specifications into the newly generated constitution and specs folder.

## 2. Use Cases
- Adopting PromptEngine in a mature, active codebase containing existing rules (like `instructions.txt` or `README.md` guideline blocks).

## 3. Inputs
- **Arguments**:
  - `[legacy_rules_path]`: Path to existing manual rules files or directories to read and consolidate.
- **Flags**:
  - `-f`, `--force`: Overwrite existing files without asking.
  - `-n`, `--non-interactive`: Run migration using auto-inferred parameters.

## 4. Outputs
- **Files Created/Modified**:
  - `AGENTS.md` (filled with reverse-engineered details and legacy overrides).
  - Updates to the 10 specification files in `docs/`.
- **Console Dumps**: Summary of legacy instructions blocks imported, showing what document target they were consolidated into.

---

## 5. Interactive Behaviour
1. Prompts the user: `Confirm path to legacy guidelines file:`.
2. Reads files and displays extracted key parameters (e.g. detected database choices and custom limits).
3. Asks user: `Should these rules overwrite standard PromptEngine configurations? (y/n)`.
4. Asks: `Do you want to delete the legacy rules file after migration is successful? (y/n)`.

## 6. Non-Interactive Behaviour
1. Automatically scans for file keywords (e.g. `instructions.txt`, `rules.md`).
2. Generates `AGENTS.md` and spec templates, appending custom legacy rules to `docs/BusinessRules.md`.
3. Suppresses file deletion prompts.

---

## 7. Expected Workflow
- Run `promptengine migrate instructions.txt` in project directory.
- Verify consolidation report logs.
- Commit generated documents to git.

## 8. Error Handling
- **Parsing Failure**: If legacy rules use unreadable formats, prints error log and falls back to generating standard specifications templates while copying legacy content to `docs/BusinessRules.md`.

## 9. Future Extensions
- Automated scanning of wiki pages (such as GitHub Wiki paths) to pull remote requirements.

## 10. Related PromptEngine Workflows
- **[Migrating Existing Projects Guide](../../guides/04-migrate-existing-project.md)**.
