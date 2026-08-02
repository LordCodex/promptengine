---
document_id: core-project-bootstrap-and-developer-onboarding-standard
title: Project Bootstrap and Developer Onboarding Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-testing-engineering-standard
  - core-git-and-collaboration-standard
  - core-cicd-and-deployment-standard
  - core-infrastructure-and-devops-standard
  - core-legacy-modernization-and-refactoring-standard
  - core-refactoring-standards-and-safe-migration-workflow
  - core-code-review-engineering-standard
  - core-documentation-engineering-standard
  - core-ai-agent-engineering-workflow-standard
  - core-reusable-ai-prompt-template-library
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Project Bootstrap and Developer Onboarding Standard

## Purpose & Inheritance
This document defines the core standards for initializing new software repositories, setting up local development environments, and onboarding human engineers and AI agents. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), the [Infrastructure SRE Standard](14-infrastructure-and-devops-standard.md), and the [Documentation Engineering Standard](19-documentation-engineering-standard.md). It establishes project discovery protocols, stack selection grids, initial repo layouts, and bootstrap checklists.

---

## 1. Project Initialization Philosophy

The initial decisions made during repository setup determine the long-term maintainability, security, and scaling potential of the application.
- **Enforce Simple Foundations**: Start with the simplest architecture that satisfies the requirements. Avoid overengineering, complex patterns, or microservice structures before a product has users or proven scale requirements.
- **Explicit Version Pinning**: Pin all runtime versions (Node, PHP, Flutter SDK, Docker base images) and dependency packages on Day 1. Avoid loose range selectors (`^` or `*`) to prevent compilation errors caused by third-party updates.
- **Zero-Trust Defaults**: Configure security structures (HTTPS validation, route authentication policies, database parameterization layers, and secrets managers) before writing business logic features.

---

## 2. Project Discovery Phase

Before initializing a codebase repository, the technical lead must document these two categories in `/docs/discovery/`:

### Business Discovery
- **System Purpose**: High-level target goals and customer flows.
- **Core Actors**: Types of users and access permission boundaries.
- **Business Invariants**: Key calculations (pricing models, commission algorithms, data lifetimes).
- **Constraints**: Compliance boundaries (e.g. GDPR, HIPAA, local payment controls).

### Technical Discovery
- **Target Scale**: Estimated user actions, request limits, data retention sizes.
- **External Integrations**: Third-party APIs (Stripe, Twilio, OAuth).
- **Performance Budgets**: Max API latencies ($<200\text{ ms}$), mobile startup budgets ($<2\text{ s}$).

---

## 3. Stack Selection Process

Do not select runtime technologies based on popularity or personal preference. Use this structural grid to align technologies with project goals:

### Runtime Stack Selection Matrix
| Tier | Stack Option | Primary Fit Criteria | Key Trade-offs |
| :--- | :--- | :--- | :--- |
| **Backend** | **PHP / Laravel** | Fast crud loops, server-rendered views, structured queue setups. | Memory overhead; single-threaded execution. |
| | **Go** | High concurrency processing, microservices, minimal memory footprints. | High boilerplate code; smaller ecosystem. |
| | **Node (TS)** | Full-stack JS architectures, real-time WebSockets sync pipelines. | Single-threaded thread blocking risks. |
| **Frontend** | **Vue 3** | Single Page Applications, dashboard panels, widgets. | Client-side hydration compilation required. |
| | **Nuxt 3** | Server-Side Rendering (SSR), SEO-focused web apps. | Nitro execution layer adds server complexity. |
| | **Inertia.js** | Monolithic bridge routing; no dedicated REST API required. | Coupled views; requires matching backend adapters. |
| **Mobile** | **Flutter / Dart**| Cross-platform apps (Android + iOS) with high custom UI. | Large binary payloads footprint. |
| **Infra** | **VPS / Docker** | Low-cost projects, simple monolithic setups. | Manual system configurations management. |
| | **Managed App** | Zero ops requirements; automatic scaling. | Higher cost boundaries; vendor lock-in. |

---

## 4. Architecture & Module Planning

Before writing code, map the application layers using these parameters:
- **Folder Boundaries**: Isolate business workflows into domain modules (e.g., `Billing`, `Identity`, `Catalog`).
- **Data Flow Guidelines**: Enforce unidirectional data flow: UI views trigger Action controllers, Action controllers query Repositories, Repositories write to DB models.
- **Authentication & Security Strategy**: Establish JWT rotation lifecycles, OAuth endpoints, and cookie parameters.

---

## 5. Initial Repository Structure

Every new project must bootstrap using this standardized repository directory tree:

```text
├── .github/
│   └── workflows/      # Pinned CI pipeline actions configs
├── config/             # Framework-level environment parameters
├── core/               # Shared logic utility wrappers
├── docs/               # Architecture records (ADRs) and setup guides
├── database/
│   ├── migrations/     # Database version schemas
│   └── seeders/        # Structured local mock database values
├── tests/              # Core unit, integration, and E2E suites
├── .env.example        # Environment variables templates (no live values)
├── .gitignore          # Pinned temp files ignore list
├── README.md           # Entry developer manual
└── docker-compose.yml  # Local dev container configurations
```

---

## 6. Development Environment Setup

- **Version Managers**: Enforce language version managers in configurations (e.g., `.nvmrc` for Node, `.sdkmanrc` for Java, `.tool-versions` for `asdf`).
- **Mock Seeding**: Database seeders must populate the database with enough realistic data to allow offline development (e.g., admin users, sample catalog invoices).
- **Setup Diagnostics Script**: Include a setup script (e.g., `bin/setup.sh`) that installs dependencies, copies `.env.example`, runs migrations, and verifies that the compiler passes.

---

## 7. Foundations: Security & Database

### Security Foundation
- **KMS / Secrets Managers**: Store developer credentials exclusively in local `.env` files (ignored by git).
- **Least Privilege Access**: Configure database container configurations with minimal access user rules (no root access permissions for local application runners).
- **Validation Engine**: Standardize on schema validation middleware handlers (Zod, Form Requests) from Day 1.

### Database Foundation
- **Model Key Schema**: Map table primary keys to UUIDv4 strings rather than auto-incrementing integers.
- **Referential Integrity Constraints**: Every relationship migration must explicitly declare foreign key behaviors (e.g., `onDelete('restrict')`).

---

## 8. First Development Tasks (Execution Sequence)

Follow this exact order of operations when launching a newly initialized project:

```text
1. Repo Setup ──> 2. Env Config ──> 3. Database Schema ──> 4. Auth Policies 
                                                                  │
      7. CI/CD Pipeline <── 6. Testing Harness <── 5. Core Actions ┘
```

1. **Project Setup**: Initialize the repository folders structure and framework packages.
2. **Environment Configuration**: Create `.env.example` configurations and container dockerfiles.
3. **Database Foundation**: Run base migrations schemas and mock seeder classes.
4. **Authentication Foundation**: Configure JWT keys and routing access gates middleware.
5. **Core Architecture**: Build base repository adapters and domain Action files.
6. **Testing Foundation**: Configure the test runner framework, environment config overrides, and mock HTTP hooks.
7. **CI/CD Foundation**: Deploy build workflows checking code styling, lints, and test execution.

---

## 9. Feature Development Workflow

All code modifications for new features must execute through this pipeline:

```text
[1. Understand] ──> [2. Plan] ──> [3. Design] ──> [4. Implement]
                                                        │
         [7. Document] <── [6. Review] <── [5. Test] <──┘
```

- **Understand**: Clarify edge cases and constraints with product leads.
- **Plan**: Create the implementation plan and verify impact surfaces.
- **Design**: Map database schemas, endpoints schemas, and class diagrams.
- **Implement**: Write target code blocks matching style conventions.
- **Test**: Execute tests verifying success pathways, validations, and failures.
- **Review**: Propose PR changes and verify against review checklist guidelines.
- **Document**: Update architecture docs, ADR logs, and API specifications.

---

## 10. Developer Onboarding

Onboard new human engineers and AI agents using this checklist:
- **Day 1 Goal**: Run the local setup script (`bin/setup.sh`) and verify that all automated tests pass successfully.
- **Day 1 Reading List**:
  - Read [README.md](../README.md) for local run command parameters.
  - Read [docs/architecture/](19-documentation-engineering-standard.md) for data flows.
  - Read [docs/decisions/](19-documentation-engineering-standard.md) to understand current technology choices.

---

## 11. Legacy Repository Bootstraps

Before starting work on an existing, undocumented legacy codebase, complete this analysis phase:
1. **Inventory Scan**: List all active services, APIs, database engines, and third-party vendor dependencies.
2. **Dependency Audit**: Scan version dependencies against security alerts directories (CVE databases).
3. **Draft Architecture Map**: Document how data flows through the application's layers.

---

## 12. Monorepos & Multi-Project Framework

When structuring monorepo systems containing multiple backend/frontend/mobile applications:
- **Shared Domain Libraries**: Place common domain interfaces, API contract structures, and visual designs inside a shared package (`/packages/shared/`).
- **Independent Routing boundaries**: Applications must not import modules directly from other app domains (enforce clean import boundaries).
- **Deployment Coordination**: Use version boundaries mapping when executing cross-project modifications.

---

## 13. Decision Matrices

### Matrix 1: Monolith vs. Microservices
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Startup team, simple business workflows, low scaling complexity | **Monolith** | Low operations overhead; simple deployment pipelines. |
| Distinct engineering teams, independent scale points, decoupled domain contexts | **Microservices** | Allows separate delivery teams; high operations complexity. |

### Matrix 2: Relational SQL vs. NoSQL
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Financial databases, strict schema validation requirements, structured relationships | **SQL** | Strong ACID compliance guarantees. |
| High throughput writes of unstructured data streams, logs, document caches | **NoSQL** | Dynamic schemas; cheap scaling storage. |

### Matrix 3: API-First vs. Server-Rendered (Blade/Livewire)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Single Page Applications, mobile companion apps, third-party consumers | **API-First** | Reusable endpoints contract layers. |
| Simple internal admin dashboards, rapid prototypes, CRUD loops | **Server-Rendered**| Eliminates the need for separate client APIs and deployment paths. |

### Matrix 4: Simple CRUD structure vs. Clean Architecture (DDD)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Standard content systems, low complexity workflows, simple models | **Simple CRUD** | Minimal boilerplate; quick development speed. |
| Complex enterprise business rules, multiple actors, long-lived workflows | **Clean Architecture**| Isolates business rules from external framework dependencies. |

### Matrix 5: Managed App Service vs. Self-Hosted VPS
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Lean developer teams, zero ops engineers, budget for managed infrastructure | **Managed** | Automated backups and scaling. |
| Custom container setups, database sizes that exceed cloud budgets | **Self-Hosted** | Control over system configurations; manual maintenance. |

---

## 14. AI Agent Bootstrap Rules

AI agents initializing or writing code in a new repository must follow these rules:

1. **Never Generate Code Immediately**: Do not generate massive, template-driven code bases on the first prompt. Ask clarifying questions first.
2. **Never Create Redundant Folders**: Follow the initial repository folder organization. Do not introduce custom folders (like `/helpers/` or `/utils/`) without approval.
3. **No Dynamic Dependencies Additions**: Do not suggest adding package libraries unless the current pinned dependencies cannot support the feature.
4. **Enforce the Workflow Sequence**: Follow the exact First Development Tasks sequence (Repo setup $\rightarrow$ Env configuration $\rightarrow$ Migrations $\rightarrow$ Auth policies).

---

## 15. Project Bootstrap Quality Checklist

Verify these criteria before declaring project bootstrap complete:

### Architectural Setup
- [ ] Is the repository folder layout aligned with the standard tree structure?
- [ ] Are runtimes and package versions pinned explicitly (no loose tags)?

### Local Environment
- [ ] Does the setup script copy `.env.example` and boot containers?
- [ ] Does the database seeder populate the database with realistic mock data?
- [ ] Do all automated tests pass successfully?

### Security Verification
- [ ] Are all credentials and API keys stored exclusively in local `.env` files?
- [ ] Is user input schema validation configured?
- [ ] Are database container users configured with minimal privileges?

---

## References
- Simplicity Rules: [core/02-architecture-and-simplicity.md](02-architecture-and-simplicity.md)
- Universal Coding Guidelines: [core/05-universal-coding-standards.md](05-universal-coding-standards.md)
- Testing Suites Setup: [core/11-testing-engineering-standard.md](11-testing-engineering-standard.md)
- CI Pipeline Automation: [core/13-cicd-and-deployment-standard.md](13-cicd-and-deployment-standard.md)
- SRE Configurations: [core/14-infrastructure-and-devops-standard.md](14-infrastructure-and-devops-standard.md)
