# CLI Roadmap (CLI-Roadmap.md)

This document maps out the phased implementation plan for the PromptEngine CLI. Development is split into core diagnostic, context resolution, validation, and full pipeline automation milestones.

---

## Roadmap Milestones

```text
  [ v0.1: Bootstrap & Diagnostics ] ──> [ v0.2: Context & Prompts ]
                                                  │
  [ v1.0: Hooks & Plugins Exts ]    <── [ v0.3: Health & Changes Analyzer ]
```

---

### v0.1: Workspace Bootstrap & Local Audit (Alpha)
*Focus: Getting repositories adopted and audited locally.*

* **Command Scope**:
  - `init`: Bootstraps greenfield constitution.
  - `scan`: Audits package managers and files on disk.
  - `version` & `help`: Engine outputs and parameters descriptions.
* **Underlying Engines**:
  - Initial version of the **Project Detection Engine** (scans package manifests).
* **Success Criteria**: A user can run `promptengine init` in a directory and get a structurally correct `AGENTS.md` and `docs/` templates.

---

### v0.2: Context, Prompts & Doctor (Beta)
*Focus: Token management and local linking integrity.*

* **Command Scope**:
  - `context`: Assembles required specs maps for target tasks.
  - `prompt`: Injects variables and outputs copy-and-paste blocks.
  - `doctor`: Validates relative link targets.
  - `config`: Handles global and local preference settings.
* **Underlying Engines**:
  - Initial **Context Builder** (reads manifest mappings).
* **Success Criteria**: Integrations with Cursor and Claude CLI are functional using `promptengine context --task add-feature` output feeds.

---

### v0.3: Change Analysis, Health, & Reviews (Staging)
*Focus: Automated auditing, change tracking, and scoring.*

* **Command Scope**:
  - `analyze`: Scans git diff files and flags out-of-sync specs.
  - `health`: Computes metric logs for documentation and tests.
  - `review`: Scans security, performance, and UI patterns.
  - `sync`: Updates `Progress.md` and check states.
* **Underlying Engines**:
  - **Change Analysis Engine** (git diff parser).
  - **Documentation Health Engine** (scoring logic).
  - **Review Engine** (static logic analysis).
* **Success Criteria**: Health reports compile automatically in under 5 seconds.

---

### v1.0: Git Hooks, Plugins, & Upgrades (Production Release)
*Focus: Pipeline automation, framework custom plugins, and lifecycle maintenance.*

* **Command Scope**:
  - `hooks`: Installs git pre-commit triggers and CI/CD pipelines.
  - `plugins`: Lists and installs stack-specific adapters.
  - `install`: Pre-configures libraries.
  - `update`: Safe migrations for templates.
  - `migrate`: Adopts PromptEngine in legacy systems with full rule translation.
* **Success Criteria**: Zero-config pre-commit hook runs on client machines, warning if a schema diff lacks matching `docs/Database.md` updates.
