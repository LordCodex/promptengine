---
document_id: core-code-review-engineering-standard
title: Code Review Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-testing-engineering-standard
  - core-git-and-collaboration-standard
  - core-legacy-modernization-and-refactoring-standard
  - core-refactoring-standards-and-safe-migration-workflow
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Code Review Engineering Standard

## Purpose & Inheritance
This document defines the core standards for peer code reviews and automated pull request (PR) evaluations. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), the [Git & Collaboration Standard](12-git-and-collaboration-standard.md), and the [Testing Engineering Standard](11-testing-engineering-standard.md). It establishes quality verification metrics, severity classifications, and constructive communication protocols for developers and AI agents.

---

## 1. Code Review Philosophy

Code review is a **knowledge-sharing and risk-reduction process**, not an ideological validation gate or a personal styling critique. 

### Core Review Objectives
Every code review must evaluate the pull request against these four questions:
1. **Goal Alignment**: Does this change solve the target problem correctly according to requirements?
2. **Operational Safety**: Is the implementation secure, performant, and resilient under failure states?
3. **Maintainability**: Will a developer unfamiliar with this change understand it easily six months from now?
4. **Regression Prevention**: Could this change introduce unexpected side-effects in other parts of the application?

---

## 2. Review Roles & Responsibilities

### The Author
- **Atomic Submissions**: Submit small, logically isolated Pull Requests that address a single requirement.
- **Provide Context**: Document the problem solved, testing logs, and design trade-offs inside the PR description.
- **Perform Self-Review**: Review your own git diff before requesting peer review to clean up debug comments and formatting noise.

### The Reviewer
- **Find Technical Risks**: Focus on finding security vulnerabilities, logical flaws, N+1 query patterns, and testing gaps.
- **Avoid Subjective Rewrites**: Do not request that the author rewrite working, readable code merely to match your personal implementation preferences.
- **Provide Constructive Feedback**: Ask questions and explain the technical reasoning behind your change requests.

---

## 3. Before Reviewing (Preparation Protocol)

To prevent shallow reviews ("LGTM" on massive diffs), reviewers must execute this protocol before looking at code:
1. **Read the Ticket/Issue**: Understand the business requirements and edge cases.
2. **Inspect the Architecture Mappings**: Review which models, API routes, or database tables are modified.
3. **Verify the Deployment Plan**: Check for database migrations or potential breaking changes.

---

## 4. Five-Stage Review Process

Reviewers must execute reviews in these sequential stages:

```text
[Stage 1: Understand] ──> [Stage 2: Design] ──> [Stage 3: Implementation]
                                                       │
         [Stage 5: Testing] <── [Stage 4: Risks] <─────┘
```

- **Stage 1 — Understand The Change**: Verify *why* the change is occurring. Check if the PR description matches the file modifications.
- **Stage 2 — Review Design**: Verify that the changes fit the playbook architecture. Check that concerns are separated (e.g. no SQL inside visual templates).
- **Stage 3 — Review Implementation**: Evaluate naming clarity, nesting flattens, method sizes, and duplicate code logic.
- **Stage 4 — Review Risks**: Evaluate performance latency, security policies, database locks, and API backward-compatibility contracts.
- **Stage 5 — Review Testing**: Ensure business behaviors and validation boundaries have corresponding integration tests.

---

## 5. Code Quality & Architecture Review

### Code Quality Checklist
- **Intention-Revealing Naming**: Verify that variables, classes, and methods clearly describe their purpose (e.g., `isPaymentAuthorized` instead of `chkPay`). Avoid abbreviations.
- **Flat Control Flow**: Check that code blocks do not exceed 3 levels of nested loops. Enforce Guard Clauses (early returns) to exit functions immediately on failures.
- **Readability**: Ensure that the code flow is logical and does not require complex comments to explain basic behavior.

### Architectural Compliance
- **Separation of Concerns**: Ensure UI presentation files, business logic modules, and database query abstractions do not leak into each other.
- **Module Boundaries**: Verify that dependencies flow in the correct direction (e.g., domain modules must not import presentation widgets).

---

## 6. Security Review Checklist

Reviewers must evaluate security parameters on every pull request:

- [ ] **Authentication**: Are route access gates active? Are tokens parsed securely?
- [ ] **Authorization**: Are backend policies (Laravel Policies, RBAC rules) checked (never trust client UI visibilities)?
- [ ] **Input Validation**: Is all incoming data validated on the server using typed schemas (Form Requests, Zod)?
- [ ] **SQL Injections**: Are database queries parameterized? Are raw SQL injections prevented?
- [ ] **XSS Prevention**: Are outputs HTML-encoded? Is the use of raw rendering blocks (`v-html`, `Blade raw`) restricted and sanitized?
- [ ] **Secrets Hygiene**: Are API tokens, database passwords, or certificates excluded from commits?

---

## 7. Performance & Database Review

- **N+1 Query Detection**: Audit Eloquent/SQL query calls inside loops. Verify that relationships are eager-loaded (`with()`) before iterating.
- **Database Locks mitigation**: Check that database migrations do not perform blocking locks on high-traffic tables.
- **Memory Optimization**: Verify that large datasets are streamed or paginated (no loading millions of rows into a single array).
- **CPU & Main Thread Checks**: On mobile clients (Flutter), ensure heavy parsing operations are offloaded to background Isolates to prevent UI lag.

---

## 8. API, Frontend & Mobile Review

### API Design Review
- **Backward Compatibility**: Ensure changing API response payloads does not break existing clients (verify resources wrappers).
- **HTTP Status Codes**: Check that routes return semantic status codes (e.g. `201 Created` on creation, `422 Unprocessable` on validation failures).

### Frontend Review (Vue / Nuxt)
- **Component Size**: Verify that components remain small (under 300 lines) with state logic extracted to composables or Pinia stores.
- **Accessibility (A11y)**: Check for semantic HTML tags, keyboard navigation capabilities, and screen reader ARIA labels.

### Mobile Review (Flutter)
- **Widget Const Constructors**: Verify that immutable widget constructors use the `const` keyword.
- **Secure Token Caches**: Check that JWTs and secrets are stored in secure OS Keychains (never in plain SharedPreferences).

---

## 9. Review Comments Standard

Review comments must explain the technical concern, describe the potential impact, and suggest a resolution path. Avoid vague, opinionated, or hostile remarks.

```text
Bad Comment:  This is terrible code.
Bad Comment:  Fix this query.
```

```text
Good Comment: [BLOCK] This loop triggers an N+1 query vulnerability because it calls '$invoice->customer' for every row. This will cause database performance issues. Consider eager loading the relationship using 'with('customer')' in the controller query before passing the dataset.
```

---

## 10. Review Severity Levels

Reviewers must prefix their comments with one of these severity tags to establish clear action items:

### 1. `[BLOCK]` (Blocking)
- **Action**: Must be resolved before the pull request can be merged.
- **Examples**: Security vulnerabilities, potential data loss, broken compiler states, missing tests for critical paths.

### 2. `[IMPORTANT]` (Important)
- **Action**: Should be addressed, but can be deferred to a follow-up ticket if release timelines are tight.
- **Examples**: Performance regressions, minor maintainability issues, missing docstrings on public APIs.

### 3. `[SUGGESTION]` (Suggestion)
- **Action**: Optional improvement. The author can choose to implement it or ignore it.
- **Examples**: Alternative naming structures, minor style cleanups, code style preferences.

---

## 11. Reviewing Legacy Repositories

When evaluating changes in legacy repositories:
- **Prioritize Risk Reduction**: Focus on verifying that modifications do not break backward compatibility or degrade current performance limits.
- **Accept Pragmatic Modernization**: Do not block a pull request because it touches a legacy file without rewriting the entire file. Approve incremental upgrades that improve the code structure.

---

## 12. AI Code Review Workflow

AI agents reviewing code in this repository must execute these checks:

1. **Verify Problem Solution**: Check that the code changes solve the user's specific request.
2. **Enforce Playbook Architecture**: Validate that files are placed in correct folders and follow stack rules (e.g., Composition API setup, Laravel Form Requests).
3. **Audit Security Checkpoints**: Check for missing authorization gates or unparameterized queries.
4. **Identify Code Smells**: Identify deep nesting levels or giant classes, and suggest early exits or decomposition.
5. **Evaluate Test Coverage**: Verify that tests assert actual behavioral outcomes, not framework internals.

---

## 13. Pull Request Review Checklist

Use this checklist to evaluate a pull request before granting approval.

### Functionality & Design
- [ ] Do the changes satisfy the business requirements?
- [ ] Does the implementation fit the existing architecture (no mixed concerns)?
- [ ] Is complexity minimized (no premature abstractions)?

### Quality & Naming
- [ ] Are variables, methods, and classes named clearly according to intent?
- [ ] Is the code structure easy to read (no deep nesting or giant methods)?

### Security Hardening
- [ ] Are input validations and authorization gates enforced on the server?
- [ ] Are SQL queries parameterized (no raw injection vulnerability)?
- [ ] Are secrets and keys excluded from the commits?

### Performance & DB
- [ ] Are database queries optimized (no N+1 query patterns)?
- [ ] Are migrations backward-compatible?

### Testing & Verification
- [ ] Are automated tests included (covering success, error, and edge cases)?
- [ ] Do all tests pass in the local and CI runners?

---

## References
- Safe Refactoring Workflow: [core/16-legacy-modernization-and-refactoring-standard.md](16-legacy-modernization-and-refactoring-standard.md)
- Testing & CI Integrations: [core/11-testing-engineering-standard.md](11-testing-engineering-standard.md)
- Conventional Commit Rules: [core/12-git-and-collaboration-standard.md](12-git-and-collaboration-standard.md)
