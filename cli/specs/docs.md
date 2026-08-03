# CLI Command Specification: `docs`

---

## 1. Purpose
Generates boilerplate structures or validates specific documents inside the `docs/` folder.

## 2. Use Cases
- Generating a blank Architecture Decision Record (ADR) template inside `docs/Decisions.md` containing all standard headers.
- Re-scaffolding a missing `docs/Troubleshooting.md` file.

## 3. Inputs
- **Flags**:
  - `-d`, `--doc-name`: The target document key (e.g. `PRD`, `Database`, `Decisions`).
  - `-f`, `--force`: Overwrite existing files without asking.

## 4. Outputs
- **Files Created**:
  - The specified document template is created or appended under the mapped directory.
- **Console Dumps**: Confirmation link indicating file path.

---

## 5. Interactive Behaviour
1. Prompts user: `Select document to scaffold:`.
2. Asks for validation if file exists: `File docs/PRD.md already exists. Overwrite? (y/n)`.
3. If creating an ADR in `docs/Decisions.md`, prompts user: `Enter ADR title:`, and generates the templated section containing Context, Decision, and Consequences blocks.

## 6. Non-Interactive Behaviour
1. Directly scaffolds files if flags are passed.
2. If files exist and `--force` is missing, halts and prints error.

---

## 7. Expected Workflow
- Run `promptengine docs -d Decisions` in terminal.
- Fill out prompt titles.
- The new ADR block is appended.

## 8. Error Handling
- **Invalid Doc Key**: Exits code `1` and prints list of valid document options.

## 9. Future Extensions
- Automated formatting of tables and links during generation.

## 10. Related PromptEngine Workflows
- **[Managing Project Documentation Guide](../../guides/06-project-documentation.md)**.
