# CLI Command Specification: `help`

---

## 1. Purpose
Renders a help catalog detailing usage syntax, commands list, and global flags mapping for the CLI.

## 2. Use Cases
- A user queries command parameters or wants to see a list of valid flags.

## 3. Inputs
- **Arguments**:
  - `[target_command]`: Specific command key to view help for (defaults to global help screen).

## 4. Outputs
- **Console Dumps**:
  - Prints usage manual guidelines to standard output.

---

## 5. Interactive Behaviour
1. Directly prints help details. No TTY prompts.

## 6. Non-Interactive Behaviour
1. Same as interactive.

---

## 7. Expected Workflow
- Run `promptengine help` or `promptengine help init` in terminal.
- Review parameters syntax.

## 8. Error Handling
- **Invalid command argument**: If target command doesn't exist, logs warning and prints global help screen instead.

## 9. Future Extensions
- Automated terminal auto-completion configurations installers for Zsh and Bash.

## 10. Related PromptEngine Workflows
- **[CLI Command Reference Guide](../CLI-Command-Reference.md)**.
