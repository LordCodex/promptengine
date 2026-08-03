---
document_id: workflows-feature-implementation
title: Feature Implementation Workflow
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-03
---

# Feature Implementation Workflow

This document defines the structured workflow for implementing new features in PromptEngine-compliant applications. Both human developers and AI agents must follow this workflow to guarantee code correctness, simplicity, and architectural coherence.

---

## 1. Requirement Dissection and Alignment
Before writing any code or modifying the schema:
- Deconstruct the feature request into inputs, outputs, and business invariants.
- Cross-reference with the [Universal Coding Standards](../core/05-universal-coding-standards.md) to ensure naming conventions and design parameters are aligned.
- Complete a pre-coding review check as defined in the [Feature Implementation Checklist](../checklists/01-feature-implementation-checklist.md).

---

## 2. Technical Planning & Design
- Create an edge case matrix analyzing zero/null bounds, database constraints, rate limits, and failure states.
- Draft pseudo-code for core business logic, separating data validation, authorization checks, and processing logic.
- Follow the **Three Questions Mnemonic** (Authentication -> Authorization -> Validation) defined in the [Security Engineering Standard](../core/08-security-engineering-standard.md) for every endpoint.

---

## 3. Data & API Modeling
- Model database schema extensions using the expand-and-contract pattern to ensure zero-downtime deployments (refer to [Database Engineering Standard](../core/06-database-engineering-standard.md)).
- Document API payloads using OpenAPI/Swagger formats to maintain client-server boundaries (refer to [Data and API Modeling](../core/03-data-and-api-modeling.md)).

---

## 4. Test-Driven Setup
- Define success, failure, and edge case assertions *before* implementing concrete classes.
- Write unit tests first to lock in method signatures and contract bounds (refer to [Testing Philosophy](../core/04-testing-philosophy.md)).

---

## 5. Implementation Phase
- Write simple, predictable, single-responsibility code.
- Avoid premature optimizations or unnecessary abstractions.
- Sanitize inputs at the boundary and escape all variables before output rendering (refer to [Secure Coding Standards](../security/01-secure-coding.md)).

---

## 6. Verification and Deployment
- Execute the test suite and run standard linters (e.g. PHPStan, ESLint, Flutter Analyze).
- Run the [Pre-Deployment Checklist](../checklists/03-deployment-checklist.md) before pushing changes to release branches.
