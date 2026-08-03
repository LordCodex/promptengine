# PromptEngine Command Line Interface (CLI) Foundation

This directory houses the foundational design, architecture, configuration schemas, and command specifications for the future production-quality PromptEngine CLI.

---

## Foundation Design Index

* **[CLI Specification (CLI-Spec.md)](CLI-Spec.md)**: Core design goals, interactive vs. non-interactive modes, and console UX patterns.
* **[CLI Architecture (CLI-Architecture.md)](CLI-Architecture.md)**: Conceptual layers of the CLI, the Detection Engine, Health Engine, Hook System, and AI provider abstractions.
* **[CLI Implementation Architecture (CLI-Implementation-Architecture.md)](CLI-Implementation-Architecture.md)**: Go repository folder structure, package boundary designs, dynamic commands registries, and detailed engines architectures.
* **[CLI Roadmap (CLI-Roadmap.md)](CLI-Roadmap.md)**: Phased implementation release schedule (v0.1 through v1.0).
* **[CLI Configuration (CLI-Configuration.md)](CLI-Configuration.md)**: Formats and preferences schema (local and global levels).
* **[CLI Command Reference (CLI-Command-Reference.md)](CLI-Command-Reference.md)**: Comprehensive command mapping tables, flag patterns, and shell exit codes.

---

## Command Specifications (`specs/`)

Every command is documented in detail inside the **[specs/](specs/)** sub-folder:

| Command | Purpose |
| :--- | :--- |
| **[init](specs/init.md)** | Bootstrap a greenfield project constitution and core docs. |
| **[migrate](specs/migrate.md)** | Consolidate and adopt PromptEngine in legacy systems. |
| **[doctor](specs/doctor.md)** | Validate local repository linkages, syntax integrity, and links. |
| **[analyze](specs/analyze.md)** | Trace code changes and diffs to identify out-of-sync specs. |
| **[scan](specs/scan.md)** | Reverse-engineer dependencies, routing patterns, and schema tables. |
| **[review](specs/review.md)** | Execute security, performance, and accessibility reviews. |
| **[sync](specs/sync.md)** | Automatically sync progress and track tasks. |
| **[update](specs/update.md)** | Migrate templates or upgrade standard stack files. |
| **[context](specs/context.md)** | Resolve minimal playbook context arrays for specific tasks (token management). |
| **[prompt](specs/prompt.md)** | Output copy-and-paste prompts with runtime parameters injected. |
| **[docs](specs/docs.md)** | Scaffolds or validates specific documentation specifications. |
| **[generate](specs/generate.md)** | Generates database schemes or mock data specs. |
| **[config](specs/config.md)** | Query and configure user preferences. |
| **[install](specs/install.md)** | Download and configure stacks patterns plugins. |
| **[health](specs/health.md)** | Compute diagnostic scoring across documentation and testing coverages. |
| **[hooks](specs/hooks.md)** | Inject git pre-commit triggers and CI pipeline triggers. |
| **[plugins](specs/plugins.md)** | Manage stack-specific or proprietary plugins. |
| **[version](specs/version.md)** | Display engine version logs. |
| **[help](specs/help.md)** | Print command-line parameter help screens. |
