# CLI Command Specification: `config`

---

## 1. Purpose
Reads or modifies user preferences in global configuration files or project settings in local configurations files.

## 2. Use Cases
- Registering a custom local path for PromptEngine clone.
- Updating default AI client model preference settings.

## 3. Inputs
- **Flags**:
  - `--get [parameter]`: Query parameter value.
  - `--set [parameter]=[value]`: Set parameter value.
  - `--global`: Apply change to user global configuration level (defaults to local workspace `.promptengine.json`).

## 4. Outputs
- **Console Dumps**: Renders the queried config value or writes status confirmation line.

---

## 5. Interactive Behaviour
1. If no flags are passed, prints the active configuration file parameters.
2. Prompts user: `Modify configuration parameter? (y/n)`.
3. If yes, prompts with list: `Select parameter category to change:`.
4. Executes the edit and validates inputs type.

## 6. Non-Interactive Behaviour
1. Directly reads or updates the settings values.
2. If schema validation fails, prints errors list to `stderr` and exits code `1`.

---

## 7. Expected Workflow
- Run `promptengine config --set ai.default_provider=gemini` in project directory.
- Verify setting update logs.

## 8. Error Handling
- **Invalid Key**: If key is not in JSON schema structure, exits with error code `1`.

## 9. Future Extensions
- Automated synchronization of configuration files between developer teams.

## 10. Related PromptEngine Workflows
- **[Configuration System Specifications](../CLI-Configuration.md)**.
