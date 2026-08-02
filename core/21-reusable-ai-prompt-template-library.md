---
document_id: core-reusable-ai-prompt-template-library
title: Reusable AI Prompt Template Library
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
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Reusable AI Prompt Template Library

## Purpose & Inheritance
This library defines the standard reusable templates for prompting AI coding agents and assistants (such as Codex, Claude Code, Cursor, and Gemini). It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md) and the [AI Agent Engineering Workflow Standard](20-ai-agent-engineering-workflow-standard.md). These templates enforce consistent quality across project analysis, feature additions, bug fixes, refactoring, security reviews, and deployment prep tasks.

---

## 1. Prompt Design Principles

To maximize AI code quality, all custom prompts must follow this five-part structure:
- **Context**: Define the current state, tech stack versions, framework patterns, and configurations.
- **Goal**: Define the exact problem to be solved.
- **Constraints**: Set bounds (e.g. no new external packages, lock specific files, preserve backward compatibility, performance targets).
- **Expected Output**: Define the required files, code structures, diff formats, and explanations.
- **Validation**: Define success checks (e.g. compile requirements, specific tests execution, latency profile).

---

## 2. Reusable Prompt Templates

### 1. Project Analysis Prompt
Use this prompt when onboarding the AI to an unfamiliar codebase.

```markdown
# Project Analysis Prompt

## Context
I am introducing you to a new, unfamiliar repository. 

## Goal
Perform a high-level scan and analysis of this project to build a correct architectural mental map.

## Constraints
- Read only configuration manifests (like composer.json, package.json, pubspec.yaml, docker-compose.yml) and root folder structures.
- Do not read individual business logic files yet.

## Expected Output
Provide a report containing:
1. **Core Stack & Versions**: Major framework, language, and runner versions.
2. **Architecture Model**: Monolith vs. modular monolith vs. microservices; client-server bridges (e.g. Inertia, REST, GraphQL).
3. **Directory Map**: High-level folder purposes.
4. **Existing Abstraction Patterns**: Active patterns (e.g. Actions, Services, Repositories, Pinia stores, Riverpod controllers).
5. **Onboarding Recommendations**: Key files to inspect next to understand core logic.
```

---

### 2. Feature Implementation Prompt
Use this prompt when requesting the addition of a new business feature.

```markdown
# Feature Implementation Prompt

## Context
- Tech Stack: {STACK}
- Existing Feature Abstractions: {ABSTRACTIONS}
- Affected Subsystem: {SUBSYSTEM}

## Goal
Implement the following feature: {FEATURE_DESCRIPTION}.

## Constraints
- Enforce the Single Responsibility Principle.
- Follow existing directory conventions and naming styles (casing, suffixes).
- Do not introduce new external package dependencies unless explicitly approved.
- Ensure all queries are parameterized. Enforce server-side validation.

## Expected Output
1. **Implementation Plan**: Proposed classes, files to modify, database schema migrations.
2. **Target Code Edits**: Atomic, formatted code implementations.
3. **Automated Tests**: Unit and integration test specifications.

## Validation
- The code must compile with zero errors.
- Verify security (authorization rules checked) and performance (no N+1 queries introduced).
```

---

### 3. Bug Investigation Prompt
Use this prompt to isolate and remediate application bugs.

```markdown
# Bug Investigation Prompt

## Context
- Subsystem: {SUBSYSTEM}
- Symptoms: {SYMPTOMS}
- Reproduction Steps: {REPRODUCTION_STEPS}

## Goal
Investigate and fix the root cause of this bug without introducing regressions.

## Constraints
- Do not guess or suggest random configuration tweaks.
- Propose the minimum possible code change to resolve the issue correctly.
- Do not change unrelated code formatting or business rules.

## Expected Output
1. **Root Cause Analysis**: Explanation of why the bug occurs.
2. **Minimal Bugfix**: The target code diff resolving the issue.
3. **Reproduction Test**: A test case that fails before the fix and passes after.
```

---

### 4. Legacy Code Analysis Prompt
Use this prompt when analyzing undocumented or older codebase components.

```markdown
# Legacy Code Analysis Prompt

## Context
- Target Legacy Module: {MODULE_PATH}
- Technology / PHP version: {LANGUAGE_VERSION}

## Goal
Analyze the legacy module to understand its execution paths and dependencies.

## Constraints
- Do not propose refactoring or rewrites yet.
- Focus strictly on documenting actual current behavior (including known bugs).

## Expected Output
1. **Execution Workflow**: Mappings of inputs to outputs.
2. **Data Dependencies**: Relational tables, integrations, and session flags read/written.
3. **Risk Analysis**: Dependencies on global parameters, performance bottlenecks, and security gaps.
```

---

### 5. Safe Refactoring Prompt
Use this prompt when restructuring code to improve clarity or testability.

```markdown
# Safe Refactoring Prompt

## Context
- Class / Module under refactor: {MODULE_PATH}
- Current Test Coverage: {TESTS_STATE}

## Goal
Improve the internal structure of the target module without modifying its external behavior.

## Constraints
- Do not modify API routes, return payload keys, database fields, or validation limits.
- Step-by-step edits only. Run tests after every change.
- Flatten control structures (flatten nested blocks using early exits and guard clauses).

## Expected Output
1. **Refactoring Plan**: Proposed code extractions (e.g. methods, service helpers).
2. **Target Code Diffs**: Clean, reviewable git diff formats.
3. **Post-Refactor Validation**: Verification that all unit/integration tests pass.
```

---

### 6. Security Audit Prompt
Use this prompt to audit code segments for security vulnerabilities.

```markdown
# Security Audit Prompt

## Context
- Target Code Module: {MODULE_PATH}
- Sensitive Data Handled: {DATA_TYPE}

## Goal
Audit the code segment for security vulnerabilities and compliance gaps.

## Constraints
- Inspect authentication validations, authorization gates (RBAC/Policies), server-side input schema checks, output encoding patterns, secrets hygiene, and dependency vulnerabilities.

## Expected Output
Provide a report listing:
1. **Vulnerability Findings**: Specific issues identified (e.g. missing permission checks, unescaped raw outputs).
2. **Severity Classification**: BLOCK (immediate fix) vs. IMPORTANT (remediate soon) classification.
3. **Remediation Recommendations**: Secure code blueprints resolving each finding.
```

---

### 7. Performance Audit Prompt
Use this prompt to profile and optimize slow executing logic.

```markdown
# Performance Audit Prompt

## Context
- Code Pathway: {MODULE_PATH}
- Performance Bottleneck: {BOTTLENECK}

## Goal
Audit execution latency, query counts, and memory allocations to identify optimizations.

## Constraints
- Focus on database query optimizations (N+1, missing index, large select blocks), memory usage (streaming vs. loading arrays), network calls, and caching patterns.
- Do not add complex abstractions that degrade code readability.

## Expected Output
1. **Bottleneck Source**: Clear diagnostics (e.g., "$invoice->customer relation loaded inside loop").
2. **Optimized Code Proposal**: Performance-optimized code changes.
3. **Expected Metrics Improvement**: Estimated database query reductions or memory savings.
```

---

### 8. Database Design Prompt
Use this prompt to design relational database schemas.

```markdown
# Database Design Prompt

## Context
- Feature / Concept to model: {BUSINESS_CONCEPT}
- Target DB Engine: {DB_ENGINE}

## Goal
Design a database schema to store the target concept efficiently.

## Constraints
- Enforce referential integrity (foreign keys, cascade parameters).
- Choose appropriate data types (e.g. unsigned integers, precise decimals for currencies).
- Add indexes for search columns. Avoid redundant indexes.
- Do not create over-normalized or over-complex schemas.

## Expected Output
1. **Schema Design**: SQL DDL commands or framework migration schemas.
2. **Index Strategy**: Explanation of columns indexed and why.
3. **Entity Relationship Mapping**: Descriptions of foreign relationships (1:M, M:M).
```

---

### 9. API Design Prompt
Use this prompt to design RESTful HTTP APIs.

```markdown
# API Design Prompt

## Context
- Target Resource: {RESOURCE}
- Consumer Audience: {CONSUMERS}

## Goal
Design a RESTful API structure for the target resource.

## Constraints
- Follow standard HTTP status codes (200, 201, 400, 401, 403, 404, 422, 500).
- Wrap responses in standard envelopes (e.g. data keys).
- Enforce input validations, pagination models, and authentication policies.

## Expected Output
1. **Endpoint Routing List**: URL patterns and matching HTTP verbs.
2. **JSON Payloads Specifications**: JSON request schemas and response formats.
3. **Error Envelopes Blueprint**: Format for validations and authorization errors.
```

---

### 10. Frontend UI Implementation Prompt
Use this prompt to generate frontend interfaces.

```markdown
# Frontend UI Implementation Prompt

## Context
- Framework / UI Stack: {UI_STACK}
- Stylings: {STYLING_STANDARD}
- Design Concept: {UI_DESCRIPTION}

## Goal
Build a responsive, accessible, and performant user interface component.

## Constraints
- Separate presentation components from business logic (extract state to Pinia/Composables).
- Ensure accessibility compatibility (screen reader ARIA labels, semantic markup, keyboard focus).
- Avoid inline layout hacks. Ensure responsiveness on desktop, tablet, and mobile views.

## Expected Output
1. **Component Template**: Accessible markup.
2. **State Logic Script**: Setup composition scripts or class controllers.
3. **Styling Block**: Scoped styling classes.
```

---

### 11. Backend Implementation Prompt
Use this prompt to generate server-side processing workflows.

```markdown
# Backend Implementation Prompt

## Context
- Framework / Runtime: {BACKEND_STACK}
- Target Workflow: {WORKFLOW_DESCRIPTION}

## Goal
Implement a performant, secure backend business workflow.

## Constraints
- Place business logic in single-responsibility Actions or Services (controllers must remain thin).
- Enforce strict database transactions (`DB::transaction`).
- Write feature integration tests validating inputs, permissions, and DB states.

## Expected Output
1. **Backend Classes**: Proposed Action/Service classes.
2. **Controller Endpoint Mapping**: Wire routes to new abstractions.
3. **Automated Feature Tests**: Verification test suites.
```

---

### 12. Code Review Prompt
Use this prompt to review git diff changes.

```markdown
# Code Review Prompt

## Context
- Git Diff Content: {GIT_DIFF}
- Repository Architecture Rules: {ARCHITECTURE_RULES}

## Goal
Evaluate the provided git diff for correctness, security, performance, and architecture.

## Constraints
- Categorize comments using severity prefixes: `[BLOCK]` (critical fix), `[IMPORTANT]` (remediate soon), `[SUGGESTION]` (optional style).
- Provide constructive, technical arguments. Avoid subjective style styling comments.

## Expected Output
A code review report containing:
1. **Design & Architectural Check**: Does code follow separation of concerns?
2. **Security Issues**: Missing permission gates or validation holes.
3. **Performance Regressions**: N+1 queries, memory leaks, or unoptimized DB reads.
4. **List of Actionable Review Comments**: Constructive remediation advice.
```

---

### 13. Test Creation Prompt
Use this prompt to generate automated test assertions.

```markdown
# Test Creation Prompt

## Context
- Target Code under test: {FILE_CONTENT}
- Testing Framework: {TEST_FRAMEWORK}

## Goal
Write a comprehensive test suite covering the target code module.

## Constraints
- Test core business outcomes, not language syntax or framework internals.
- Cover success flows, validation failures, boundary value limits, and exceptions.
- Mock external third-party integrations (HTTP calls, email dispatches).

## Expected Output
1. **Mocking Setup**: Mock configurations for external libraries.
2. **Test File Code**: Clean, assertion-driven tests code structure.
3. **Explanation of Edge Cases Covered**: Brief list of cases verified.
```

---

### 14. Documentation Generation Prompt
Use this prompt to generate markdown manuals.

```markdown
# Documentation Generation Prompt

## Context
- Target Module / Subsystem: {MODULE_PATH}
- Target Audience: {AUDIENCE}

## Goal
Write accurate, developer-friendly documentation for the target module.

## Constraints
- Focus on the "why" and usage guides. Avoid writing low-level code descriptions that duplicate variable declarations.
- Format documentation in clean GitHub Markdown using relative file links.

## Expected Output
1. **System Overview**: Business goals and architectural purpose.
2. **Usage Instructions**: Code samples or command structures.
3. **Dependencies Mapping**: Databases and third-party APIs used.
```

---

### 15. Dependency/Framework Migration Prompt
Use this prompt when executing framework or major library upgrades.

```markdown
# Dependency/Framework Migration Prompt

## Context
- Current Version: {CURRENT_VERSION}
- Target Upgrade Version: {UPGRADE_VERSION}
- Package/Runtime Manager: {PACKAGE_MANAGER}

## Goal
Plan and execute a safe, step-by-step upgrade of the target runtime or package dependency.

## Constraints
- Avoid bulk changes. Target upgrading one version increment at a time.
- Identify deprecated methods and replace them with modern equivalents.
- Provide a rollback plan if compilation or test execution fails.

## Expected Output
1. **Migration Roadmap**: Sequential list of package file updates.
2. **Deprecated Code Replacements**: Code diff drafts replacing EOL calls.
3. **Verification & Rollback Instructions**: Check commands and recovery plans.
```

---

### 16. Production Readiness Review Prompt
Use this prompt before launching code changes to production.

```markdown
# Production Readiness Review Prompt

## Context
- Feature Subsystem: {FEATURE_PATH}
- Target Environment: {ENVIRONMENT}

## Goal
Verify that the target code is production-ready across security, scaling, monitoring, and recovery categories.

## Constraints
- Verify: Secrets configurations (no keys in code), database scaling (indexes, migration locks), logging diagnostics setup, API rate-limiting configs, and rollback setups.

## Expected Output
Provide a checklist report:
1. **Ready Parameters**: Confirmed configuration areas.
2. **Production Vulnerabilities**: Identified risks (e.g. missing timeout limits, inadequate logger traces).
3. **Pre-Launch Recommendations**: Remediation steps.
```

---

### 17. AI Self-Critique Prompt
Enforce this prompt internally to force the AI to review its own proposals before output.

```markdown
# AI Self-Critique Prompt

## Context
- Original Prompt Goal: {GOAL}
- Proposed Code Solution: {PROPOSED_CODE}

## Goal
Critique your own proposed solution to identify flaws, code smells, security risks, or overengineering.

## Constraints
- Evaluate against: Overengineering (patterns without problems), security holes, naming violations, and performance issues.

## Expected Output
1. **Critique Findings**: Honest assessment of potential problems.
2. **Trade-offs Explanation**: Explicit comparison of alternative implementations.
3. **Refined Code Draft**: Optimizations and cleanups to the original code proposal.
```

---

## 3. Stack-Specific Prompt Templates

### PHP & Laravel Feature Addition
```markdown
# PHP & Laravel Feature Addition Prompt

## Context
- Target Laravel Version: 11.x
- Target Path: /app/Actions/ and /app/Http/Controllers/

## Goal
Implement the following backend controller and Action logic: {WORKFLOW_GOAL}.

## Constraints
- Do not place business logic in controllers. Extract execution steps to an Action class.
- Use Laravel Form Request classes for validation.
- Enforce route authorization using Laravel Policies (`$this->authorize`).
- Avoid N+1 queries. Eager load relations.
- Write tests using Pest/PHPUnit.

## Expected Output
1. **Migration & Models Code**
2. **Form Request & Action Classes**
3. **Pest Feature Test File**
```

---

### Vue 3, Nuxt 3 & Inertia Component Addition
```markdown
# Vue 3, Nuxt 3 & Inertia Component Addition

## Context
- Frontend Framework: Vue 3 Composition API (<script setup lang="ts">)
- Rendering Bridge: Inertia.js (passing server attributes)

## Goal
Build the frontend visual component: {COMPONENT_GOAL}.

## Constraints
- Use strict TypeScript interfaces.
- Separate UI layout from business logic (move state calculations to composables).
- Ensure WCAG contrast ratios and add appropriate screen reader labels.

## Expected Output
1. **Vue Component Code File**
2. **TS Typings & Interfaces**
3. **Composable Logic File**
```

---

### Flutter & Dart State Refactoring
```markdown
# Flutter & Dart State Refactoring

## Context
- State Management: {STATE_MANAGER_LIBRARY}
- Target View: {VIEW_PATH}

## Goal
Refactor state management logic in the view to improve separation of concerns.

## Constraints
- View widgets must remain stateless and visual-only.
- Move state, validation, and database updates to a view model controller.
- Use Riverpod/Bloc configurations.
- Widgets must use const constructors.

## Expected Output
1. **Refactored View Widget Code**
2. **State Controller Class File**
3. **Widget State Integration Test**
```

---

### Docker & CI/CD Pipeline Audit
```markdown
# Docker & CI/CD Pipeline Audit

## Context
- Target Dockerfile: {DOCKERFILE_PATH}
- Target CI Pipeline: {CI_CONFIG_PATH}

## Goal
Audit Docker containers and CI workflows for security vulnerabilities and pipeline speed bottlenecks.

## Constraints
- Verify: Multi-stage Docker optimization (pinned base tags, non-root users), runner step execution speeds (caching node/pub modules, lock files), and secrets handling.

## Expected Output
1. **Dockerfile Vulnerabilities & Patches**
2. **CI Pipeline Step Speed Optimizations**
3. **Safe Credentials Storage Recommendations**
```

---

## 4. Prompt Variable Placeholders

When executing prompt templates, inject parameters using this standardized syntax:
- `Technology`: Language and framework version (e.g. PHP 8.3, Laravel 11.x).
- `Path`: Absolute file system path (e.g. `[RegisterUserAction.php](file:///app/Actions/RegisterUserAction.php)`).
- `Task`: Concise goal description.
- `Constraints`: Strict architectural rules (e.g., "no raw database queries").
- `Validation`: Verification steps (e.g., "run `composer test`").

---

## 5. Prompt Usage & Context Efficiency Guide

To optimize context window usage and avoid AI response decay:

- **Use Task-Focused Micro-Prompts**: Do not compile all 17 prompts into a single master request. Choose the smallest relevant prompt for the task.
- **Provide Targeted Context Files**: Pass only the specific files relevant to the prompt target.
- **Do Not Repost Baseline Playbooks**: AI agents should fetch general rules from local playbook references rather than having them copied directly into every prompt message.

---

## References
- Simplicity Standards: [02-architecture-and-simplicity.md](02-architecture-and-simplicity.md)
- Universal Naming Standards: [05-universal-coding-standards.md](05-universal-coding-standards.md)
- Testing & QA Harnesses: [11-testing-engineering-standard.md](11-testing-engineering-standard.md)
- AI agent workflow parameters: [20-ai-agent-engineering-workflow-standard.md](20-ai-agent-engineering-workflow-standard.md)
