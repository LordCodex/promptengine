---
document_id: core-legacy-modernization-and-refactoring-standard
title: Legacy Code Modernisation and Safe Refactoring Standard
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
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Legacy Code Modernisation and Safe Refactoring Standard

## Playbook Metadata
- **Purpose**: Establishes standards for code smell identification, refactoring workflows, and architectural migration patterns (Strangler, Branch-by-Abstraction) for legacy codebase upgrades.
- **Scope**: Code refactoring, code review, framework upgrades, database schema expansions, and legacy system modernization.
- **When to Read**: Before performing any refactoring, migrating architectural blocks, correcting code smells, or upgrading language/framework versions.
- **Related Playbooks**: [Universal Coding Standards](05-universal-coding-standards.md), [Architecture Standards](02-architecture-and-simplicity.md), [Database Engineering Standard](06-database-engineering-standard.md), [Testing Standard](11-testing-engineering-standard.md).
- **Canonical Source**: This is the canonical document for Legacy Code Modernization and Refactoring standards.
- **Version**: 1.1.0
- **Last Reviewed**: 2026-08-03

---

## 1. Legacy Engineering & Refactoring Philosophy

Legacy code is not "bad code"; it is **working code that generates business value but carries operational risk**.
- **Refactoring Definition**: Refactoring is the process of improving the internal structure of code without altering its external behavior.
- **Value Over Cleanliness**: Legacy systems have users, execute business logic, and store valuable data. Do not refactor code just because its syntax is outdated.
- **Reject the Green-Field Rewrite Trap**: Avoid complete, top-down rewrites. Rewrites are high-risk, expensive, and often fail because they ignore the edge cases and business rules embedded in the legacy codebase.
- **Incremental Modernization**: Refactor code systematically using the "Boy Scout Rule" (leave the codebase cleaner than you found it).
- **Behavior Preservation**: Refactoring must not introduce new features or fix unrelated bugs in the same commit. External inputs, validation logic, error envelopes, and outputs must remain identical.
- **Complexity Reduction**: Clean code is code that is easy to read and cheap to change. If a refactoring change increases code complexity or requires developers to learn new custom abstractions, the refactoring has failed.

---

## 2. Pre-Refactoring Discovery Protocol

You must not modify a legacy module until you have completed these discovery steps:

```text
Map Dependencies ──> Identify Business Rules ──> Locate Data Schemas ──> Establish Test Baseline
```

1. **Understand Business Intent**: Document who uses the module, what input triggers it, and what success outputs are expected.
2. **Scan Architecture Dependencies**: Identify which database tables, files, or external third-party APIs the legacy module communicates with.
3. **Audit the Deployment Track**: Verify that you can deploy, verify, and roll back changes safely using the active CI/CD pipeline.

---

## 3. Legacy Code Auditing

Before refactoring, run a localized audit of the target module using these criteria:

### Audit Parameters
- **Code Quality**: Identify deep nesting levels ($>3$ loops), large functions ($>100\text{ lines}$), coupled classes, and copy-paste code duplicates.
- **Security**: Verify input validation points, session boundaries, and database query bindings. Check for hardcoded API keys or credentials.
- **Performance**: Scan for $N+1$ query patterns, memory leaks in data processing loops, and un-indexed SQL search columns.
- **Maintainability**: Check if the module has automated test coverage and documentation.

---

## 4. Code Smell Identification & Resolution

### Large Classes / Large Functions
- **Large Classes**: Classes exceeding 300 lines of code, or multiple public methods handling unrelated responsibilities. Resolve by extracting focused, single-responsibility classes. Move business rules to Actions and third-party integrations to Services.
- **Large Functions**: Methods containing multiple nested blocks, long sequential steps, and complex calculations. Resolve by applying the "Extract Method" pattern. Keep methods under 30 lines.

### Duplicate Code
- **Refactor when**: Duplicate code contains core business calculations (like tax rates, interest formulas, or permission authorizations).
- **Acceptable duplication**: Duplicating simple visual formatting rules or minor layout properties is acceptable if merging them introduces complex, over-parameterized helper functions.

### Deep Nesting (Guard Clauses)
- **Symptoms**: Deeply nested logic statements ($>3$ indentation levels).
- **Resolution**: Enforce Guard Clauses (early returns) to exit functions early if requirements fail. This flattens control structures:

```typescript
// Bad: Deeply nested conditions
function processPayment(invoice) {
    if (invoice !== null) {
        if (invoice.status === 'pending') {
            if (invoice.amount_cents > 0) {
                executePayment(invoice);
            }
        }
    }
}

// Good: Flattened control structure using Guard Clauses
function processPayment(invoice) {
    if (!invoice) return;
    if (invoice.status !== 'pending') return;
    if (invoice.amount_cents <= 0) return;

    executePayment(invoice);
}
```

### Poor Naming
- **Naming Principles**: Variable, method, and class names must be intention-revealing (e.g. `isSubscriptionExpired` instead of `checkExp`). Use names that match the business domain language. Avoid abbreviations (`inv` should be `invoice`).

---

## 5. Characterization Testing Safety Nets

A **Characterization Test** asserts the *actual current behavior* of a legacy module (including its bugs), establishing a safety baseline before you modify the code.

```text
Define Inputs ──> Execute Legacy Logic ──> Capture Actual Output ──> Codify as Strict Assertions
```

### Characterization Testing Rules
- **Capture Current Outputs**: Do not write tests based on what the documentation *says* the code should do. Write tests that assert what the code *actually does* right now.
- **Cover Edge Cases**: Pass extreme inputs (empty strings, null values, out-of-bound numbers) to record how the legacy module responds.
- **Use as a Refactoring Guard**: Run these characterization tests after every small code edit. If a test fails, you have broken existing behavior.

---

## 6. Safe Refactoring Workflow

We enforce a strict seven-step workflow for modifying legacy code:

```text
[1. Understand] ──> [2. Identify Risk] ──> [3. Add Protection Tests]
                                                     │
[7. Monitor] <── [6. Canary Deploy] <── [5. Review Diff] <── [4. Edit Code] 
```

1. **Understand**: Analyze the current file execution path and dependency models.
2. **Identify Risk**: Assess the impact of changes on database locks, external APIs, and client compatibility.
3. **Add Protection Tests**: Write characterization integration tests to verify the module's behavior.
4. **Make Small Improvements**: Apply minor refactors (e.g. renaming variables, extracting single-responsibility methods). Run tests after every change.
5. **Review Diff**: Verify that only the intended code blocks were changed (no accidental formatting modifications).
6. **Canary Deploy**: Deploy changes to a subset of staging or production servers.
7. **Monitor**: Verify error telemetry and application logs post-deployment.

---

## 7. Code Modernization Areas

### Architecture & Boundaries
- **Introduce Dependency Injection**: Replace hardcoded class initializations with Dependency Injection (DI) containers.
- **Isolate Code Boundaries**: Wrap legacy modules in Service classes or Adapters to prevent legacy patterns from leaking into newer codebase features.

### Database Schema Evolution
- **Never Modify Live Columns Directly**: When refactoring database columns, use the **Expand and Contract** pattern (see [Database Engineering Standard](06-database-engineering-standard.md)).
- **Backfill Data in Batches**: When adding non-nullable columns, run data backfill tasks in small, throttled database chunks to prevent table lockups.

### Dependency Upgrades (Framework Migration)
- **Incremental Upgrades**: When upgrading major framework versions (e.g., Laravel 10 to 11, Vue 2 to 3), upgrade one minor version at a time. Run the test suite between each upgrade step.
- **Rollback Preparation**: Ensure you can revert to the previous framework package version instantly if deployment checks fail.

---

## 8. Modernizing Legacy PHP & Laravel Systems

### Legacy Pure PHP (No Framework)
- **Gradually Introduce Autoloading**: Replace manual `require` or `include` calls with Composer PSR-4 autoloading.
- **Extract DB Queries**: Extract raw SQL queries embedded in HTML templates to dedicated Data Access Objects (DAOs) or Repository classes.
- **Introduce Validation Middleware**: Wrap page request parameters in validation checks before executing business logic.

### Legacy Laravel Applications
- **Slim Down Fat Controllers**: Move business calculations, mail dispatches, and third-party integrations from controller files into single-responsibility Actions or Services:
  ```php
  // Good: Business logic extracted from controller to Action
  class RegisterUserAction
  {
      public function execute(array $data): User
      {
          return DB::transaction(fn() => User::create($data));
      }
  }
  ```
- **Introduce Form Requests**: Replace inline `$request->validate()` blocks in controllers with dedicated Form Request classes.

---

## 9. Modernizing Legacy Frontend Client Systems (Vue/Nuxt)

When upgrading legacy Vue/Nuxt files:
- **Break Up God Components**: Extract nested sub-templates and repeated loops into smaller, stateless visual components. UI components must only manage presentation layout.
- **Extract Stateful Composables**: Move logic, state variables, and HTTP client calls from script setup blocks into reusable Composables (`composables/useDomain.ts`).
- **Clean Up Shared Global State**: Replace dynamic, untracked global variables with structured Pinia store models.
- **Avoid Wrapper Component Bloat**: Do not create custom components that merely wrap standard HTML tags without adding new behaviors or styles.

---

## 10. Modernizing Legacy Mobile Client Systems (Flutter)

When upgrading legacy Flutter files:
- **Decompose Widget Trees**: Extract large widget layouts into small, focused stateless widgets. Ensure child components compile with `const` constructors to prevent unnecessary redraw cycles.
- **Extract Controller Logic**: Move state management, validation logic, and repository calls out of widget files and into dedicated Riverpod providers or Bloc controllers.
- **Enforce Repository Boundary**: Widgets and view state controllers must never interact directly with database files or API client configurations. Route all data requests through Repository interfaces.

---

## 11. Architectural Migration Strategies

### 1. The Strangler Pattern
- **Concept**: Gradually replace legacy system components with modern implementations, route by route.
- **When to Use**: Upgrading legacy monolith systems to modular monoliths or distributed architectures.
- **Execution**: Route incoming traffic to the new implementation at the load balancer, proxy, or gateway routing layer, leaving the legacy code intact for un-migrated paths.

### 2. Branch By Abstraction
- **Concept**: Replace an internal code dependency dynamically using abstraction wrappers.
- **When to Use**: Swapping core backend systems (e.g. replacing a legacy PDF builder or changing SMS gateways).
- **Execution**:
  1. Define an Interface abstraction wrapper around the target component.
  2. Modify the codebase to consume this Interface.
  3. Build the new component implementation matching this Interface.
  4. Swap the active class binding in the dependency injector container.

### 3. Feature Flags
- **Concept**: Toggle runtime execution paths dynamically using configuration flags.
- **When to Use**: Releasing high-risk refactored code safely to production.
- **Execution**: Run the refactored logic for a small percentage of users (canary release), and monitor error rates before rolling out the change to all users.

---

## 12. Backward Compatibility Safeguards

Modernizing code must not break downstream API consumers or mobile clients.
- **Preserve API Response Contracts**: Ensure that refactored API routes return identical JSON keys and data types. Use API Resource wrappers to maintain response structures.
- **Version Routes**: If an API contract change is unavoidable, version the route (e.g., `/api/v2/invoices`). Keep the old `/api/v1/invoices` route active until legacy consumers migrate.

---

## 13. Security & Performance Verification

- **Maintain Authorization Policies**: When refactoring routing layers or controllers, verify that permission gates and middleware policies remain active. Do not bypass policy validations.
- **Benchmark Performance**: Run latency and CPU benchmarks before and after refactoring. Ensure that new abstractions have not introduced performance regressions (like N+1 queries).

---

## 14. Concurrency Safety in Legacy Systems

Legacy applications often have hidden concurrency bugs because they were originally written for single-server deployments. Modernization must address these systematically.

### Principles
- Assume **thousands of concurrent requests** and **multiple application servers** running simultaneously at all times.
- Never rely on local filesystem state for correctness of business logic — filesystem locks prevent duplicate execution on a **single server only**. Correctness must come from database constraints, transactions, atomic updates, or idempotency keys.
- Prefer **idempotent operations** — running the same operation twice must produce the same result.

### Transaction Rules
- Use database transactions for all multi-query operations. If any query in a sequence fails, the entire operation must roll back.
- Keep transactions **as short as possible**. Do not perform slow operations (HTTP calls, file writes, queue dispatches) inside a database transaction.
- Never hold locks longer than necessary. Acquire the lock, read and write the data, then commit and release immediately.
- Do not perform external HTTP requests, email sends, or queue dispatches while holding a database lock.

### Race Condition Prevention
- Never check then act on a value without first acquiring a lock:
  ```sql
  -- Correct: Lock before reading
  SELECT balance FROM wallets WHERE id = ? FOR UPDATE;
  ```
- Two simultaneous requests must never both pass the same balance check, inventory availability check, or booking availability check.
- Use atomic database operations (`UPDATE wallets SET balance = balance - ? WHERE balance >= ?`) where supported.

### Scheduled Job Safety
- Background and scheduled jobs must:
  - Use non-blocking overlap protection (database flag or lock) to prevent two instances running simultaneously.
  - Perform **bounded work per run** (process a limited number of rows, then exit).
  - Log summaries without sensitive payloads.
  - Return a non-zero exit code on failure.
  - Be registered in the existing scheduler — do not require a separate, undocumented scheduling mechanism.

---

## 15. Avoiding Overengineering

We value simplicity. Do not introduce:
- **Patterns Without Problems**: Do not apply design patterns (like Command Query Responsibility Segregation or Repository patterns) unless the codebase has structural problems that require them.
- **Extra Layers Without Need**: Do not insert intermediate service boundaries between simple CRUD controllers and Eloquent models.
- **Microservices Unnecessarily**: Keep applications structured as Monoliths or Modular Monoliths until deployment scale demands distributed services.

---

## 16. Legacy Refactoring Review Checklist

Use this checklist during code review to evaluate legacy code modifications.

### Discovery & Safety
- [ ] Has the module's business intent and integration dependencies been documented?
- [ ] Have characterization tests been written to establish a behavior baseline?
- [ ] Does the codebase pass all tests after the refactor?

### Architecture & Quality
- [ ] Have large files ($>100\text{ lines}$) or deeply nested code blocks (indentation level $>3$) been simplified?
- [ ] Are business logic calculations separated from controllers and templates?
- [ ] Has dependency injection been used instead of hardcoded class initializations?

### Database & Compatibility
- [ ] Are schema upgrades backward-compatible (no direct column drops or modifications)?
- [ ] Have API response contracts been preserved (no changes to JSON key keys or data types)?

### Security & Performance
- [ ] Are input validation and authorization policies active in the refactored code?
- [ ] Has performance benchmarking verified that no latency regressions have been introduced?
- [ ] Have all hardcoded credentials or API keys been removed from the code?

---

## 17. Decision Matrices

Use these matrices to identify the correct refactoring decision based on project context.

### Matrix 1: Rewrite vs. Refactor
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| High database dependency, active customer traffic, complex business logic | **Refactor** | Lowers deployment risks; preserves built-in edge case logic. |
| Obsolete language runtime, zero test coverage, simple single-purpose app | **Rewrite** | Low complexity; allows starting fresh with modern security controls. |

### Matrix 2: Extract Service vs. Keep Logic Inline
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Calculations used across multiple features (e.g. tax calculators) | **Extract Service** | Centralizes updates; enforces DRY (Don't Repeat Yourself). |
| Page-specific visual formatting, one-off state conversions | **Keep Inline** | Avoids overengineering; keeps simple logic easy to read. |

---

## 18. AI Modernization Rules

AI agents modernizing or refactoring code in this repository must follow these rules:

1. **Never Rewrite Code Automatically**: Do not suggest rewriting working legacy systems in bulk. Focus changes strictly on the requested files.
2. **Do Not Delete Unused Code Blindly**: Before deleting code that appears unused, run a codebase-wide search to verify it is not queried by reflection, dynamic routes, or background scripts.
3. **Write Protection Tests First**: Before modifying a legacy file, write characterization tests to assert its current output behaviors.
4. **No Direct Schema Alterations**: Do not suggest dropping database tables or columns directly in migrations. Enforce the Expand-and-Contract migration pattern.
5. **Explain Refactoring Risks**: When proposing structural changes, clearly document potential side-effects on APIs, integrations, and performance.

---

## References
- Universal Naming Rules: [core/05-universal-coding-standards.md](05-universal-coding-standards.md)
- Safe Database Migrations: [core/06-database-engineering-standard.md](06-database-engineering-standard.md)
- QA Test Verification: [core/11-testing-engineering-standard.md](11-testing-engineering-standard.md)
- Rollback Deployment: [core/13-cicd-and-deployment-standard.md](13-cicd-and-deployment-standard.md)
