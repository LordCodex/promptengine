# 06. Managing Project Documentation

This guide describes the Project Knowledge System and details how developers and AI agents must maintain the 10 core documentation specifications under the `docs/` (or `.agents/`) folder.

---

## The 10 Core Specs

### 1. `docs/PRD.md` (Product Requirements Document)
- **Purpose**: Defines **what** the product does. Lists user personas, target features, user stories, and acceptance criteria.
- **When to Edit**: When adding features, changing user roles, or editing acceptance criteria.
- **Developer Rule**: Do not write technical details (like SQL schemas or endpoint names) here. Focus strictly on user-facing requirements.

### 2. `docs/Architecture.md`
- **Purpose**: Defines system boundaries, folder structures, component responsibilities, and data flow guidelines.
- **When to Edit**: When adding service layers, changing modules boundaries, or introducing new vendor APIs.

### 3. `docs/BusinessRules.md`
- **Purpose**: Documents domain calculations, interest rates, tax formulas, state-machine transitions, and permissions rules.
- **When to Edit**: When a business formula or validation constraint changes.
- **Developer Rule**: Ensure math equations are documented in clear prose or markdown lists so the AI can translate them directly to code.

### 4. `docs/Database.md`
- **Purpose**: Relational schema tables, column types, primary key strategy, database index allocations, and foreign relationship mappings.
- **When to Edit**: In the same commit as database schema migration scripts.

### 5. `docs/API.md`
- **Purpose**: Documents endpoint URLs, request parameters, JSON headers, query configurations, validation scopes, and JSON response models (including RFC 7807 error envelopes).
- **When to Edit**: In the same pull request as controller routing edits.

### 6. `docs/Progress.md`
- **Purpose**: A living sprint check log (`task.md`) detailing active tasks, completed features, file modifications, and test results.
- **When to Edit**: Daily, as features are implemented and tested.

### 7. `docs/Roadmap.md`
- **Purpose**: High-level future releases, tech debt targets, and deprecated modules.
- **When to Edit**: Quarterly, or when planning upcoming releases.

### 8. `docs/Decisions.md`
- **Purpose**: Logs Architecture Decision Records (ADRs) detailing design options, trade-offs, and choices made.
- **When to Edit**: On major database, authentication, infrastructure, or dependency changes.

### 9. `docs/Deployment.md`
- **Purpose**: Documents CI/CD pipelines, Docker container settings, environment variable overrides (no live keys), and rollback instructions.
- **When to Edit**: When adjusting Dockerfiles, workflows, or target cloud hosts.

### 10. `docs/Troubleshooting.md`
- **Purpose**: Diagnostics checklists, error code mappings, logs patterns, and recovery sequences.
- **When to Edit**: After resolving an incident, to prevent future recurrence.

---

## Best Practices for Markdown Maintenance

To maintain clean and readable documentation:
- **No Source Code Duplication**: Do not copy-paste code files or classes into markdown files. Code updates will instantly make them stale. Explain the **intent** and **structure** instead.
- **Use Relative Paths**: Always reference other specs using relative markdown links (e.g. `[Database Schema](Database.md)`). Never use absolute local file paths or URLs.
- **Section Isolation**: When updating a document, modify only the relevant subsection. Do not let the AI rewrite the entire file as it wastes tokens and introduces hallucinations.
