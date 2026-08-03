# CLI Command Specification: `generate`

---

## 1. Purpose
Generates mock data models, API schema responses templates, or database index maps based on parameters defined in local specifications.

## 2. Use Cases
- Generating test JSON response mock structures mapping to a route definition in `docs/API.md`.
- Generating SQL DDL schemas mapping to a table declaration in `docs/Database.md`.

## 3. Inputs
- **Flags**:
  - `-t`, `--type`: Output generation type (`sql`, `json-mock`, `ts-interface`).
  - `-o`, `--output`: Save path destination file.

## 4. Outputs
- **Files Created**:
  - Generated code stub or schema file.
- **Console Dumps**: Renders generated block code to console.

---

## 5. Interactive Behaviour
1. Prompts: `Select output type:`.
2. Scans docs files (like `docs/API.md` and `docs/Database.md`).
3. If multiple targets are found, asks: `Generate SQL schema for database tables? (y/n)`.
4. Outputs the syntax to the terminal.

## 6. Non-Interactive Behaviour
1. Parses flags.
2. Directly reads specs files and prints raw code blocks to standard output.
3. Exits code `1` if dependencies cannot compile.

---

## 7. Expected Workflow
- Run `promptengine generate -t sql` in terminal.
- Review SQL code output.
- Pipe directly to db shell: `promptengine generate -t sql | mysql my_db`.

## 8. Error Handling
- **Missing Spec context**: If `docs/Database.md` lacks table column variables definition, generate logs warning and exits with code `1`.

## 9. Future Extensions
- Integrations with Prisma or schema migrators.

## 10. Related PromptEngine Workflows
- **[Database Design Prompt](../../prompts/07-architecture-change.md)**.
