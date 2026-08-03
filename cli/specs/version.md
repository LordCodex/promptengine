# CLI Command Specification: `version`

---

## 1. Purpose
Prints build, template package, and compiler version details for the PromptEngine CLI binary.

## 2. Use Cases
- Verification of version compatibility between CLI template engines and core playbook specifications.
- Log outputs during CI pipeline runs.

## 3. Inputs
- **Flags**:
  - `-s`, `--short`: Print only the version number string (e.g. `1.0.0`) without extra metadata lines.

## 4. Outputs
- **Console Dumps**: Version number string, build timestamp, commit SHA, compiler version.

---

## 5. Interactive Behaviour
1. Directly prints version data. No prompts.

## 6. Non-Interactive Behaviour
1. Same as interactive. Writes to standard output.

---

## 7. Expected Workflow
- Run `promptengine version` in terminal.
- Review output data.

## 8. Error Handling
- None.

## 9. Future Extensions
- Automatic updates availability checks.

## 10. Related PromptEngine Workflows
- **[Versioning Strategy Specification](../CLI-Roadmap.md)**.
