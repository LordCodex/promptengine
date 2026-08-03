# CLI Command Specification: `scan`

---

## 1. Purpose
Performs a deep static analysis of files on disk to map languages, dependencies, frameworks, database engines, routes, and tests setups.

## 2. Use Cases
- Auto-detecting the technology stack parameters during project onboarding or reverse-engineering cycles.

## 3. Inputs
- **Arguments**:
  - `[target_path]`: Path to directory to scan (defaults to current folder).
- **Flags**:
  - `-o`, `--output`: Format and save results to a specific file target.

## 4. Outputs
- **Console Dumps**:
  - Tables listing languages detected, versions mapped, routes files identified, database tables detected.
- **Files Created (optional)**:
  - JSON dump of codebase metadata.

---

## 5. Interactive Behaviour
1. Spins up the directory crawler showing files processed count.
2. If multiple frameworks are detected (e.g. both React and Laravel files), prompts user: `Detected both PHP (Laravel) and JS (React). Is this a Laravel backend with Inertia.js React frontend? (y/n)`.

## 6. Non-Interactive Behaviour
1. Quietly runs crawler.
2. Logs errors or unresolved detections to `stderr`.
3. Returns JSON array output.

---

## 7. Expected Workflow
- Run `promptengine scan` inside a legacy repository.
- Review mapped framework parameters.
- Run `promptengine init` or `promptengine migrate` using the output settings.

## 8. Error Handling
- **Perms Denied**: If files are locked, CLI prints warnings and skips path.

## 9. Future Extensions
- Automated detection of third-party SaaS integrations (like Stripe, SendGrid, Firebase).

## 10. Related PromptEngine Workflows
- **[Active Existing Projects Guide](../../guides/03-existing-project.md)**.
