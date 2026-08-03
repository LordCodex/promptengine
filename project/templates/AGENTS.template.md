# Project Constitution & AI Rules (AGENTS.md)

## Playbook Metadata
- **Purpose**: Defines the project-specific AI Constitution, combining core PromptEngine rules with the project-specific technical stack, architecture style, and business constraints.
- **Scope**: Applied automatically by AI agents for all project-related tasks.
- **When to Read**: First-step entry point for any AI agent interacting with the project.
- **Related Playbooks**: [Project Onboarding Standard](../01-project-bootstrap-standard.md), [Project Documentation Standard](../02-project-documentation-standard.md).
- **Version**: 1.0.0
- **Last Reviewed**: 2026-08-03

---

# AI Constitution: [Project Name]

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

This section defines the specific boundaries, technology stack, constraints, and exceptions of the project.

*(Note for the AI Agent: For **New Project Bootstrap**, populate these fields from the discovery interview. For **Existing Project Bootstrap**, reverse-engineer them from scanning the repository files and configurations on disk.)*

### 2.1 Project Overview
- **Project Name**: [Enter Project Name]
- **Description**: [Provide a brief, high-level summary of the application's purpose and goals.]

### 2.2 Technology Stack & Architecture
- **Architecture Style**: [e.g., Monolith, Microservices, Clean Architecture, Vertical Slice]
- **Backend Framework**: [e.g., Laravel, Node/Express, Go/Gin, None]
- **Frontend Framework**: [e.g., Vue 3 (Composition API), Nuxt 3, React/Next.js, Tailwind, Vanilla HTML/CSS/JS]
- **Mobile Framework**: [e.g., Flutter, React Native, None]
- **Database**: [e.g., PostgreSQL, MySQL, MongoDB, SQLite]
- **Cache**: [e.g., Redis, Memcached, In-Memory]
- **Queue / Background Jobs**: [e.g., Laravel Horizon/Redis, BullMQ, Go Channels, SQS]
- **API Style**: [e.g., RESTful (RFC 7807 for errors), GraphQL, gRPC]
- **Authentication Approach**: [e.g., JWT with cookie storage, Sanctum session, OAuth2, Firebase Auth]
- **Primary Key Strategy**: [e.g., UUIDv7, UUIDv4, Auto-increment Integer, Snowflake ID]

### 2.3 Operating Constraints & Requirements
- **Business Constraints**: [e.g., Payment integrations locked to Stripe, compliance rules like GDPR/HIPAA]
- **Legacy Constraints**: [e.g., Must remain backwards compatible with legacy API versions on disk, database tables that cannot be renamed]
- **Performance Requirements**: [e.g., API P95 latency < 200ms, page load time < 1.5s, mobile app bundle size < 50MB]
- **Security Constraints**: [e.g., Input parameter validation using Zod, role-based access control policies on all endpoints]
- **Deployment Strategy**: [e.g., Docker containers deployed via GitHub Actions to AWS ECS, fly.io VPS deployment]

### 2.4 Coding Conventions & Rules
- **Coding Style**: [e.g., PSR-12 for PHP, ESLint/Prettier for JS/TS, Dart Lints for Flutter]
- **Testing Approach**: [e.g., Unit tests with Pest, Component tests with Vitest, Integration E2E tests]
- **Project-Specific Rules**: [e.g., "All transaction values must use integer cents (never floats)", "Vue component styles must use scoped CSS"]

### 2.5 Approved Exceptions & ADR Index
- **Exceptions to Standards**: [e.g., "Approved exception: Using auto-increment integers for logs table for performance, documented in Decision ADR-002", "None"]
- **ADR Index**:
  - **ADR-001**: [Initial stack selection]
  - [List subsequent architectural decision records here]

---

## 3. Onboarding & Execution Steps

When starting a development thread or session, the AI agent must:
1. Load `AGENTS.md` to align context with the Project Constitution.
2. Read the local documentation in the `docs/` folder (such as `PRD.md`, `Architecture.md`, `Database.md`, `API.md`) to verify details.
3. Locate the `playbook-manifest.json` from the configured PromptEngine directory.
4. Select and load only the required playbooks based on task requirements.
5. Formulate an implementation plan in accordance with `core/01-thinking-and-planning.md`.
