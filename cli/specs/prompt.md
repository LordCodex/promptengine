# CLI Command Specification: `prompt`

---

## 1. Purpose
Generates and outputs custom, variables-injected prompt templates from the Prompts Library to copy and paste to AI interfaces.

## 2. Use Cases
- Generating a copy-and-paste prompt for a bug fix or feature implementation directly in the console.

## 3. Inputs
- **Flags**:
  - `-p`, `--prompt-id`: The target prompt filename or ID (e.g. `add-feature`, `bug-fix`).
  - `--vars`: Key-value pairs for template variables injection (e.g. `--vars "STACK=Laravel 11,PATH=app/Actions/Checkout.php"`).

## 4. Outputs
- **Console Dumps**:
  - Renders the variable-injected markdown prompt to standard output.

---

## 5. Interactive Behaviour
1. If `--prompt-id` is omitted, prompts the user with a list: `Select prompt template to generate:`.
2. Reads the variables required for the selected prompt (e.g. `{STACK}`, `{SYMPTOMS}`).
3. Prompts for each variable value: `Enter value for {STACK}:`.
4. Outputs the finished prompt text to console.
5. Prompts: `Copy to clipboard? (y/n)`.

## 6. Non-Interactive Behaviour
1. Requires both `--prompt-id` and `--vars` list to be present.
2. Injects variables and prints raw text.
3. If variables are missing, exits with error code `1`.

---

## 7. Expected Workflow
- Run `promptengine prompt -p bug-fix --vars "STACK=Next.js,SYMPTOMS=hydration error on load"` in terminal.
- Prompt generates instantly.
- Copy text to ChatGPT web interface.

## 8. Error Handling
- **Invalid Prompt ID**: Prints list of valid prompts and exits with code `1`.

## 9. Future Extensions
- Direct piping of generated prompt to local AI binaries like `claude` or `cursor`.

## 10. Related PromptEngine Workflows
- **[AI Prompt Library Integration Guide](../../guides/07-ai-prompt-library.md)**.
