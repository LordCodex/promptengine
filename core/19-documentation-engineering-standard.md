---
document_id: core-documentation-engineering-standard
title: Documentation Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-git-and-collaboration-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Documentation Engineering Standard

## Purpose & Inheritance
This document defines the core standards for writing, organizing, and maintaining technical software documentation. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), the [Security Engineering Standard](08-security-engineering-standard.md), and the [Git & Collaboration Standard](12-git-and-collaboration-standard.md). It establishes documentation models, Architecture Decision Record (ADR) formats, and documentation hygiene guidelines.

---

## 1. Documentation Philosophy

Good documentation is a **critical component of developer experience (DX) and system durability**, not a chore to complete after coding.
- **Explain the "Why", Not Just the "What"**: Code shows *what* the system does; documentation must explain *why* the system was built that way, what design trade-offs were made, and what business rules dictate the implementation.
- **Minimize Maintenance Overhead**: Do not duplicate active code details in documentation. Avoid writing low-level code descriptions that become outdated as soon as a file is modified.
- **Match Actual Code Behavior**: Documentation that is inaccurate or outdated is more harmful than having no documentation at all, as it leads to incorrect assumptions and onboarding errors.

---

## 2. Documentation Types & Standards

We classify documentation into seven primary categories:

### 1. Project README (`README.md`)
Every repository must contain a root-level `README.md` serving as the main entry checkpoint. It must document:
- **Project Purpose**: High-level business goals, core workflows, and target features.
- **Requirements**: Pinned runtime versions (e.g. PHP 8.3, Node 20 LTS, Flutter 3.19).
- **Installation & Local Dev Setup**: Exact commands required to build, seed databases, install dependencies, and run tests.
- **Environment Setup**: Pointers to configuration files.
- **Build and Deployment Overview**: Folder structure mappings and build pipeline configurations.
Avoid turning the README into a tutorial.

### 2. Architecture Documentation
Located inside `docs/architecture/` or recorded via ADRs, this describes:
- **System Component Diagrams**: Component mappings and data ingestion loops.
- **Integration Points**: Critical external service adapters.
- **Business Logic Invariants**: Core workflows (e.g., payment ledgers, checkout pipelines).

### 3. API Documentation
- **API Spec Mappings**: Keep documentation synchronized with endpoints. Document authentication keys, endpoint paths, request formats, response formats, validation errors (`400`/`422`), rate limits, and versioning. Ensure all code examples remain accurate.

### 4. Database Documentation
- **Entity Schemas**: Document critical database tables, entity relationships (ERD), unique indexing patterns, and data ownership. Do not document standard system column tables without purpose.

### 5. Developer Setup & Onboarding
- **Bootstrapping Guides**: Document required tooling, dependencies installation paths, local Docker configs, and coding standards. A new developer must be able to boot, test, and deploy the application without requiring tribal knowledge.

### 6. Environment Variables (`.env.example`)
Document every environment variable in `.env.example`:
- **Description**: Explain the purpose of each variable.
- **Status**: Clarify if it is required or optional.
- **Example Value**: Provide a realistic, safe placeholder value.
- **Security Guidance**: Document key rotation schedules or access rules. **Never document real secrets or production credentials.**

### 7. Project Changelog (`CHANGELOG.md`)
Maintain a structured changelog at the repository root:
- **Scope**: Document breaking changes, new features, security fixes, and performance improvements.
- **Naming**: Use clear description entries; avoid vague log texts (e.g. "fixes", "updates").

---

## 3. Architecture Decision Records (ADR)

An ADR is a short, version-controlled document that records a significant technical decision, its context, and its consequences.

### ADR Workflow
Create an ADR for high-risk technical decisions, such as:
- Choosing a specific database engine (e.g., PostgreSQL vs. MongoDB).
- Swapping state management libraries (e.g., Bloc vs. Riverpod).
- Adopting a specific architecture model (e.g., Monolith vs. Microservices).

### ADR Template Blueprint (`docs/decisions/YYYY-MM-DD-title.md`)
```markdown
# ADR [Index]: [Title Description]

## Status
[Proposed | Accepted | Superseded by ADR-XX]

## Context
Describe the current system state, requirements, limitations, and the problem being solved.

## Decision
State the exact technical path selected. Explain the rationale.

## Alternatives Considered
- **Alternative A**: Description and trade-offs.
- **Alternative B**: Description and trade-offs.

## Consequences
- **Positive**: What benefits do we gain?
- **Negative / Risk**: What technical debt or complexity is introduced?

## Future Considerations
Describe future checkpoints or triggers that would require revisiting this decision.
```

---

## 4. Code Commenting Rules

Code comments must explain **intent and non-obvious business rules**, not self-explanatory syntax.

```text
Bad Comment:  // Loop through user models
Good Comment: // Users are processed in batches of 100 to prevent PHP memory limits exhaustion
```

- **Explain Workarounds**: If code contains a workaround for a third-party API limitation or OS bug, document it in a comment alongside a link to the corresponding issue ticket.
- **Naming Reduces Comments**: Choose intention-revealing class, method, and variable names (e.g., `calculateMonthlyRevenue()`) to eliminate the need for comments.

---

## 5. Documenting Business Logic

Crucial business logic must be documented in dedicated markdown files inside `docs/business-rules/`:
- **Pricing & Tax Calculations**: Exact rounding calculations rules and margin allocations.
- **Commission Calculations**: Split parameters, payouts thresholds, and transactional schedules.
- **Authentication Flows**: JWT rotations lifecycle, OAuth handshake boundaries.
- **Data Retention Policies**: Archiving schedules and user data deletion processes.

---

## 5.5. Operational Error Documentation

Complex operational failures or transient exception loops in production must be documented in dedicated runbooks:
- **Symptoms**: Detail how the failure manifests in client browsers, mobile devices, and server logs.
- **Cause**: Explain the underlying technical trigger or resource bottleneck (e.g., memory exhaustion, locked rows).
- **Resolution**: Step-by-step runbook instructions for engineers to resolve the immediate error state.
- **Prevention**: Long-term prevention strategies, code fixes, and configuration adjustments to prevent recurrence.

---

## 6. Standardized Documentation Folder Hierarchy

Maintain this directory layout at the repository root to keep documentation discoverable:

```text
docs/
├── architecture/      # High-level system design and component mappings
├── decisions/         # Architecture Decision Records (ADRs)
├── api/               # API endpoint specifications and JSON contracts
├── business-rules/    # Pricing schemas and commission equations docs
├── deployment/        # CI/CD pipelines and infrastructure environments setup
└── development/       # Local setup guides and debugging commands
```

---

## 7. Documentation Maintenance & Hygiene

- **Keep Docs in Sync with Code**: When a pull request modifies system architecture, configuration variables, API contracts, or setup steps, updating the corresponding documentation files is **mandatory before merge**.
- **No Plain Secrets**: Never document live passwords, private keys, database credentials, or API tokens. Always use env placeholders.

---

## 8. Legacy Project Documentation

When taking over undocumented legacy applications:
1. **Gradually Document the System**: Do not halt development to write comprehensive documentation.
2. **Start with the Setup**: Document the local environment setup first so teammates can boot the app.
3. **Record Discoveries**: As you analyze legacy code modules, document your findings (known issues, data structures, and edge cases) in `/docs` rather than leaving them undocumented.

---

## 9. Caching, Scaling & Third-Party Docs

- **Caching & Caching Infrastructure**: Document caching key namespaces, eviction policies, and invalidation triggers.
- **Integrations Documentation**: Document third-party integration rate limits, failover behaviors, and authentication protocols.

---

## 10. Documentation Anti-Patterns

- **Documenting Obvious Code**: Writing comments for self-explanatory variable declarations or standard framework methods.
- **Outdated Setup Guides**: Failing to update local setup commands when node or library versions change, blocking new developers.
- **Copy-Pasting Diffs**: Copy-pasting massive raw file diffs or generated JSON outputs without explaining their context.
- **Documenting Assumptions as Facts**: Documenting unverified behaviors or speculative architectures.

---

## 11. AI Documentation Rules

AI agents generating, modifying, or maintaining documentation in this repository must follow these rules:

1. **Match Active Implementation**: Do not write documentation for features that do not exist in the codebase.
2. **Use Relative Markdown Links**: Connect related documents using relative markdown links (e.g., `[naming standard](../core/05-universal-coding-standards.md)`).
3. **No Key Exposures**: Ensure no mock API keys or passwords look like production credentials.
4. **Enforce ADR Scopes**: Suggest creating an ADR file when proposing high-impact architectural changes.
5. **Keep It Concise**: Focus on clarity; avoid generating long, generic documentation pages.
6. **Trigger Updates Dynamically**: You must update the relevant documentation if a pull request triggers change in:
   - **Public behaviors** (routing, endpoints, layouts, client operations).
   - **System configuration** (environment variables, config scopes, local ports).
   - **Deployment procedures** (CI/CD integrations, Docker files).
   - **Architecture boundaries** (component relationships, dependencies direction).
   - **APIs and payload structures** (requests, responses, interfaces).

---

## 12. Documentation Review Checklist

Before completing software work, verify the documentation against this checklist:
- [ ] **README accurate**: The root-level README matches the current runtime requirements and local commands.
- [ ] **Environment documented**: Every environment variable is mapped in `.env.example` with description and safe placeholders.
- [ ] **APIs documented**: Public API specifications correctly represent paths, queries, request bodies, and error codes.
- [ ] **Architecture decisions recorded**: High-impact design selections and pattern deviations are logged via ADRs.
- [ ] **Complex logic explained**: Inherent complexity and legacy workarounds are explained via code comments ("Why", not "What").
- [ ] **No redundant comments**: Self-explanatory code lines are uncommented to prevent noise.
- [ ] **No outdated documentation**: Stale configuration guidelines, instructions, or component guides are deleted or updated.
- [ ] **Examples still valid**: JSON blocks, code blocks, and commands execute without syntax or environment failures.
- [ ] **Setup instructions verified**: Local installation and testing setups work cleanly without requiring local tribal knowledge.

---

## References
- Universal Naming Rules: [core/05-universal-coding-standards.md](05-universal-coding-standards.md)
- Secure Database Schemas: [core/06-database-engineering-standard.md](06-database-engineering-standard.md)
- Git atomic commit conventions: [core/12-git-and-collaboration-standard.md](12-git-and-collaboration-standard.md)
- Design system decisions log: [design/07-design-resources.md](../design/07-design-resources.md)
