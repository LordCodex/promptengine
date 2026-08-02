---
document_id: checklists-feature-implementation
title: Feature Implementation Checklist
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-security-engineering-standard
  - core-performance-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Feature Implementation Checklist

Use this checklist when implementing any new feature, regardless of language or framework. It is designed to be run twice: **before writing code** and **after writing code**.

---

## Before Writing Code

### 1. Understanding
- [ ] Do you fully understand the requirement? If anything is unclear, ask before writing.
- [ ] Have you confirmed the existing and expected business behaviour?
- [ ] Have you identified which existing patterns, utilities, or helpers already handle any part of this?
- [ ] Have you checked `composables/`, `utils/`, `App\Services`, `App\Helpers` before writing new logic?

### 2. Scope Definition
- [ ] Can you list every file that will be created or modified?
- [ ] Does this change touch multiple unrelated modules? If yes, stop and get explicit confirmation before proceeding.
- [ ] Does this change alter business behaviour, persistence, authorization, or deployment? State it explicitly.

### 3. Architecture Alignment
- [ ] Does the approach follow the existing architectural patterns in this codebase?
- [ ] Will the implementation use the smallest correct abstraction? (Do not wrap a simple function in a service, repository, and interface unless the problem requires it.)
- [ ] Are you solving the problem in front of you — not a speculative future version of it?

### 4. Security Pre-Check
Answer the Three Questions before writing the first line:
- [ ] **Who sends this?** → Is authentication required? Is middleware applied?
- [ ] **Are they allowed?** → Is authorization (role, ownership, policy) enforced?
- [ ] **Is the data safe?** → Is all input validated server-side using a strict whitelist?

---

## After Writing Code

### 5. Code Quality
- [ ] Does every function do exactly one thing?
- [ ] Is nesting limited to a maximum of three levels?
- [ ] Are all names descriptive and consistent with existing conventions?
- [ ] Are there zero magic numbers or strings (use named constants)?
- [ ] Are there zero `TODO` comments, placeholder methods, stub classes, or dead code?
- [ ] Are there zero debug output statements (`dd`, `var_dump`, `console.log`, `debugger`, `ray`)?
- [ ] Are there zero emojis in code, comments, strings, or logs?
- [ ] Does the code read as if written by an experienced, senior engineer — not auto-generated?

### 6. Security Review
- [ ] Are all inputs validated with a strict whitelist (not blacklist)?
- [ ] Is all dynamic output escaped before rendering in browser templates?
- [ ] Are file uploads validated by MIME type (not just extension) and stored outside `public/`?
- [ ] Is no sensitive data logged, returned to the client unnecessarily, or hardcoded?
- [ ] Are tokens compared with constant-time functions (not `==` or `===`)?
- [ ] Are financial operations wrapped in a transaction with a pessimistic lock?
- [ ] Are security headers applied to responses?

### 7. Simplicity Check
- [ ] Is there a simpler way to solve this?
- [ ] Am I adding layers the problem does not require?
- [ ] Can someone read and understand this code in a few minutes?
- [ ] Am I building for a requirement that does not exist yet?

### 8. Performance Check
- [ ] Are there any N+1 queries?
- [ ] Are all large datasets paginated or chunked?
- [ ] Is heavy work (emails, reports, image processing, webhooks) offloaded to a queue?
- [ ] Am I loading only the columns the caller needs (no `SELECT *`)?
- [ ] Am I importing more frontend code than I need?

### 9. Test Coverage
- [ ] Are there automated tests for the success path?
- [ ] Are there automated tests for the failure/edge case paths?
- [ ] Do tests cover security scenarios (unauthorized access, invalid input)?
- [ ] Does the full automated test suite pass?

### 10. Final Gate
- [ ] Would another senior developer approve this code in review without hesitation?
- [ ] Have no unrelated files been modified?
- [ ] Have no unnecessary dependencies been introduced?
- [ ] Is the implementation production-ready as submitted?

---

## References
- Security Engineering: [core/08-security-engineering-standard.md](../core/08-security-engineering-standard.md)
- Performance Engineering: [core/10-performance-engineering-standard.md](../core/10-performance-engineering-standard.md)
- Testing Standard: [core/11-testing-engineering-standard.md](../core/11-testing-engineering-standard.md)
- Universal Coding Standards: [core/05-universal-coding-standards.md](../core/05-universal-coding-standards.md)
