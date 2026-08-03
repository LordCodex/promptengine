# CLI Command Specification: `update`

---

## 1. Purpose
Checks for updates to the PromptEngine repository or standard stack playbooks, downloading or merging template patches safely.

## 2. Use Cases
- Upgrading standard stacks playbooks templates (e.g. updating Nuxt 3 templates rules to Nuxt 4 standard).
- Migrating local project specs structure to newer version schemas.

## 3. Inputs
- **Flags**:
  - `--dry-run`: View templates modified list without modifying disk files.
  - `--all`: Upgrades all local stack playbooks.

## 4. Outputs
- **Files Modified**:
  - Updates standard templates under local `promptengine/` clone.
- **Console Dumps**: Check report showing files updated, version changelog, and compat warnings.

---

## 5. Interactive Behaviour
1. Contacts registry/repository backend (or local path source).
2. Lists templates with available updates.
3. Prompts user: `Update 'stacks/php-laravel/laravel-logic.md' to v1.2.0? (y/n)`.

## 6. Non-Interactive Behaviour
1. Quietly runs check.
2. Automatically pulls and overrides files in the local copy directory.
3. Prints list of modified playbooks.

---

## 7. Expected Workflow
- Run `promptengine update` inside repository root.
- Review and select playbooks to patch.
- Apply.

## 8. Error Handling
- **Merge Conflict**: If a developer modified a standard playbook file locally, update prints error mapping and copies the update as `[filename].new.md` to prevent local changes loss.

## 9. Future Extensions
- Automated conflict resolution helper directly in the terminal interface.

## 10. Related PromptEngine Workflows
- **[Upgrade System Specifications](../CLI-Architecture.md)**.
