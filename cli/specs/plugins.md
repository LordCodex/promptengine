# CLI Command Specification: `plugins`

---

## 1. Purpose
Lists, disables, enables, or uninstalls custom extension plugins that load proprietary workflows, custom stacks playbooks, or extra compliance rules.

## 2. Use Cases
- A team wants to configure a proprietary in-house security checklist rules plugin across their developers local environments.

## 3. Inputs
- **Flags**:
  - `--list`: Print all registered plugin adapters.
  - `--install [id]`: Install plugin by key.
  - `--uninstall [id]`: Remove plugin configuration.

## 4. Outputs
- **Console Dumps**:
  - List of active plugins, authors, source locations, and enabled statuses.

---

## 5. Interactive Behaviour
1. If no flags are passed, prints list showing active plugins.
2. Prompts user: `Manage active plugins? (y/n)`.
3. If yes, displays checkbox array to enable/disable specific custom adapters.

## 6. Non-Interactive Behaviour
1. Executes list or install targets directly.
2. If registry target is invalid, prints diagnostic errors to `stderr` and exits code `1`.

---

## 7. Expected Workflow
- Run `promptengine plugins --list` to check integrations.
- Toggle plugins using the CLI.

## 8. Error Handling
- **Plugin Incompatible**: Warns user if plugin requirements map to incompatible PromptEngine version limits.

## 9. Future Extensions
- GPG signature checks for verified plugin authors.

## 10. Related PromptEngine Workflows
- **[Plugin Architecture Guide](../CLI-Architecture.md)**.
