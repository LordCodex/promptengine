# AI Constitution: Goshenwell Admin Portal

This project uses PromptEngine to enforce strict, high-quality software engineering standards.

> [!IMPORTANT]
> **Instructions for the AI Agent:**
> - Read this file first before analyzing the codebase or starting implementation.
> - Under Section 1 (PromptEngine Core Rules), adhere strictly to the listed guidelines.
> - Under Section 2 (Project Constitution), follow the project-specific boundaries and tech stacks defined below.
> - Do not modify Section 1 rules. Section 2 should evolve with the project.

---

## 1. PromptEngine Core Rules

This section contains the fundamental, immutable engineering rules inherited from PromptEngine. These rules apply to all tasks and files in this repository.

1. **AI Entry Point**: Always read `AGENTS.md` first before starting any task or generating code.
2. **Understand Before Implementation**: Always read the relevant project documentation (e.g., PRD, Architecture, Database, API, Business Rules) and load appropriate PromptEngine standards before writing any code.
3. **Follow the Workflow**: Classify the task and follow the matching workflow (e.g., Feature Implementation, Refactoring, Security Hotfix).
4. **Synchronized Documentation**: Keep project documentation fully synchronized with the codebase. If code changes affect requirements, APIs, databases, or deployment configurations, update the corresponding documentation files (under `docs/`) in the same pull request.
5. **No Speculative Requirements**: Never invent requirements, business logic formulas, or API contracts. If a requirement is ambiguous or undocumented, seek developer clarification.
6. **Respect the Source of Truth**: Treat the codebase and approved documentation as the source of truth. Respect local project documentation overrides when they conflict with generic PromptEngine standards.
7. **Simplicity & Maintainability**: Avoid premature optimization or unnecessary abstraction layers. Keep code boring, clean, and maintainable.
8. **Do Not Blindly Trust Defaults**: Critical choices must be validated against the project's architecture, security posture, and performance budgets. Do not accept framework scaffolding or defaults without review.
9. **Document Architectural Decisions**: Major decisions (e.g., databases, API structures, auth mechanisms) must be logged as Architecture Decision Records (ADRs) in `docs/Decisions.md` (or `docs/decisions/`).
10. **Dual-Layer Validation Gate**: Verify code correctness using automated tests and perform security/performance reviews before completing tasks.

---

## 2. Project Constitution

This section defines the specific boundaries, technology stack, constraints, and exceptions of the Goshenwell Admin Portal project.

### 2.1 Project Overview
- **Project Name**: Goshenwell Admin Portal
- **Description**: An administrative management system for DOVAL parishes to orchestrate automated email broadcasts, schedule parish notifications, and audit user logs.

### 2.2 Technology Stack & Architecture
- **Architecture Style**: Monolith with decoupled frontend layer
- **Backend Framework**: PHP 8.2 / Laravel 10 (without Octane)
- **Frontend Framework**: Vue 3 (Composition API) integrated via Inertia.js with Tailwind CSS
- **Mobile Framework**: None
- **Database**: MySQL 8.0
- **Cache**: Redis (used for session storage and cache tags)
- **Queue / Background Jobs**: Redis queues managed by Laravel Horizon
- **API Style**: RESTful API endpoints returning RFC 7807 problem details on error
- **Authentication Approach**: Session-based cookie auth via Laravel Sanctum
- **Primary Key Strategy**: UUIDv7 for all primary keys to guarantee ordering and decentralized ID generation

### 2.3 Operating Constraints & Requirements
- **Business Constraints**: Integrates exclusively with Paystack and Flutterwave payment gateways. All financial calculations must be run in integer cents.
- **Legacy Constraints**: Must support compatibility with the legacy `.php` scripts in the `cron/` folder (such as `parish_notifications.php` and `email_broadcast.php`). Do not delete or refactor these without explicit confirmation.
- **Performance Requirements**: P95 endpoint latency must be under 150ms. Dashboard widget updates must execute asynchronously using queue workers.
- **Security Constraints**: Strict role-based access control (RBAC) checked on every request. Input parameter validation using Form Requests.

### 2.4 Coding Conventions & Rules
- **Coding Style**: PSR-12 for PHP, ESLint + Prettier for JS/TS.
- **Testing Approach**: Unit and feature testing with Pest. Minimum 80% coverage required on all new modules.
- **Project-Specific Rules**: All database schemas must use raw foreign keys targeting UUID fields; do not use framework shortcuts that implicitly create auto-increment columns.

### 2.5 Approved Exceptions & ADR Index
- **Exceptions to Standards**:
  - **Exception 1**: Using traditional auto-increment integers on the `audit_logs` table (instead of UUIDv7) for indexing and insertion throughput performance (Documented in Decisions/ADR-003).
- **ADR Index**:
  - **ADR-001**: Stack Selection (Laravel, Vue 3, Inertia)
  - **ADR-002**: Database ID strategy (UUIDv7 default)
  - **ADR-003**: Audit log table key exception

---

## 3. Onboarding & Execution Steps

When starting a development thread or session, the AI agent must:
1. Load `AGENTS.md` to align context with the Project Constitution.
2. Read the local documentation in the `docs/` folder (such as `PRD.md`, `Architecture.md`, `Database.md`, `API.md`) to verify details.
3. Locate the `playbook-manifest.json` from the configured PromptEngine directory.
4. Select and load only the required playbooks based on task requirements.
5. Formulate an implementation plan in accordance with `core/01-thinking-and-planning.md`.
