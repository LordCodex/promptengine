---
document_id: core-refactoring-standards-and-safe-migration-workflow
title: Refactoring Standards and Safe Migration Workflow
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
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Refactoring Standards and Safe Migration Workflow

## Purpose & Inheritance
This document defines the core standards for refactoring code bases and migrating software architectures. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), the [Security Engineering Standard](08-security-engineering-standard.md), the [Testing Engineering Standard](11-testing-engineering-standard.md), and the [Legacy Modernization & Safe Refactoring Standard](16-legacy-modernization-and-refactoring-standard.md). It establishes workflows for code-smell correction, structural updates, Strangler/Branch-by-Abstraction pattern migrations, and overengineering mitigation.

---

## 1. Refactoring Philosophy

Refactoring is the process of **improving the internal structure of code without altering its external behavior**.
- **Behavior Preservation**: Refactoring must not introduce new features or fix unrelated bugs in the same commit. External inputs, validation logic, error envelopes, and outputs must remain identical.
- **Complexity Reduction**: Clean code is code that is easy to read and cheap to change. If a refactoring change increases code complexity or requires developers to learn new custom abstractions, the refactoring has failed.
- **Maintainability Over Aesthetic**: We reject refactoring changes that are driven by personal style preferences. Refactoring is justified only when it makes the codebase more testable, readable, and resilient to regression.

---

## 2. When to Refactor

Enforce this decision framework before starting a refactoring task:

```text
Change is Risky / Bugs Repeat / Code is Coupled ──> Refactor
Syntactic Style Choice / Unused Pattern / Personal Preference ──> Do Not Refactor
```

### Justified Refactoring Scenarios
- **Mixed Responsibilities**: A class or component is handling UI rendering, database access, and validation rules simultaneously.
- **Frequent Bug Regressions**: A specific class is a hotspot for recurring bug reports.
- **High Friction to Change**: Adding a simple feature requires modifying five unrelated files.
- **Severe Code Duplication**: Identical business rules are copied across three distinct locations, leading to drift in logic.

---

## 3. Step-by-Step Refactoring Workflow

Developers and AI agents must follow this five-step workflow for all refactoring tasks:

```text
[Step 1: Understand] ──> [Step 2: Define Goal] ──> [Step 3: Protect Behavior]
                                                          │
                              [Step 5: Validate] <── [Step 4: Edit Small] ┘
```

1. **Step 1 — Understand**: Analyze the current file execution path, identify business invariants, scan code dependencies, and evaluate existing test coverage.
2. **Step 2 — Define Goal**: State the exact scope of the refactor (e.g. "Extract third-party API adapter to separate service class; reduce controller size").
3. **Step 3 — Protect Behaviour**: Ensure you have a test safety net. If test coverage is low, write characterization tests to capture the existing outputs. Wrap high-risk changes in Feature Flags.
4. **Step 4 — Make Small Changes**: Make changes incrementally. Commit small, logical edits. Run tests after every change.
5. **Step 5 — Validate**: Run the full validation suite, run performance latency benchmarks, audit security authorization policies, and verify logs to ensure zero regression.

---

## 4. Code Smell Identification & Resolution

### Large Classes
- **Symptoms**: Classes exceeding 300 lines of code; multiple public methods handling unrelated responsibilities; difficult testing setups.
- **Resolution**: Extract focused single-responsibility classes. Move business rules to Actions and third-party integrations to Services.

### Large Functions
- **Symptoms**: Methods containing multiple nested blocks, long sequential code steps, and hard-to-read calculations.
- **Resolution**: Apply the "Extract Method" pattern. Extract steps to private helper functions with descriptive naming. Keep methods under 30 lines.

### Duplicate Code
- **Refactor when**: Duplicate code contains core business calculations (like tax rates, interest formulas, or permission authorizations).
- **Acceptable duplication**: Duplicating simple visual formatting rules or minor layout properties is acceptable if merging them introduces complex, over-parameterized helper functions.

### Deep Nesting
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

## 5. Architectural Refactoring

- **Separation of Concerns**: Separate the codebase into distinct layers: Presentation UI layer (Blade, Vue widgets), Business Logic layer (Actions, domain rules), Data Access layer (Eloquent models, SQL queries), and External Integration adapters (Service classes).
- **Dependency Management**: Inject class dependencies using constructors (Dependency Injection) instead of hardcoding static class instantiation calls.
- **Interfaces usage rules**:
  - Do not create interfaces for every class.
  - Create interfaces only when the system has multiple concrete implementations (e.g., matching a `SmsProviderInterface` with `TwilioSms` and `NexmoSms` implementations).

---

## 6. Backend Stack Refactoring (PHP / Laravel)

When refactoring Laravel structures, distribute responsibilities according to this matrix:

| Abstraction | Scope | Justification |
| :--- | :--- | :--- |
| **Action** | Single business workflow (e.g., `PayInvoiceAction`). | Encapsulates a transaction; reusable across HTTP controllers and console commands. |
| **Service** | Third-party integrations (e.g., Stripe Payment SDK). | Isolates external dependencies; prevents SDK models from leaking into domain layers. |
| **Policy** | User permissions validations (e.g., `InvoicePolicy`). | Centralizes authorization checks; prevents custom logic inside controllers. |
| **Job** | Asynchronous task (e.g., `GeneratePdfReport`). | Offloads slow operations to background workers. |

---

## 7. Frontend Stack Refactoring (Vue / Nuxt)

- **UI Component Refactors**: Extract nested template blocks to smaller stateless child widgets. UI components must only manage presentation layout.
- **Composable Extractions**: Extract UI helper variables, local loaders, and Axios client calls from templates into reusable composables.
- **State stores**: Extract shared page attributes to feature-scoped Pinia stores. Remove untracked global variable systems.
- **Avoid Wrapper Component Bloat**: Do not create custom components that merely wrap standard HTML tags without adding new behaviors or styles.

---

## 8. Mobile Stack Refactoring (Flutter)

- **Decompose Widget Trees**: Extract large widget layouts into small, focused stateless widgets. Ensure child components compile with `const` constructors to prevent unnecessary redraw cycles.
- **Extract Controller Logic**: Move state management, validation logic, and repository calls out of widget files and into dedicated Riverpod providers or Bloc controllers.
- **Enforce Repository Boundary**: Widgets and view state controllers must never interact directly with database files or API client configurations. Route all data requests through Repository interfaces.

---

## 9. Database Refactoring Safety

- **Expand and Contract Pattern**: Never drop database columns or tables directly. Run migrations in distinct steps:
  1. *Expand*: Add the new column or table; modify the application code to write to both the old and new columns.
  2. *Backfill*: Run a background script to copy historical data from the old column to the new column in batches.
  3. *Contract*: Update the application code to read from the new column only. Drop the old column in a later release.
- **Preserve Data Meaning**: Do not reuse existing database fields to store new, unrelated data models. This corrupts data history.

---

## 10. Migration Strategies

### 1. Strangler Pattern
- **Concept**: Gradually replace legacy system components with modern implementations, route by route.
- **When to Use**: Upgrading legacy monolith systems to modular monoliths or distributed architectures.
- **Execution**: Route incoming traffic to the new implementation at the load balancer or proxy layer, leaving the legacy code intact for un-migrated paths.

### 2. Branch By Abstraction
- **Concept**: Replace an internal code dependency dynamically using abstraction wrappers.
- **When to Use**: Swapping core system components (such as database engines or payment providers).
- **Execution**:
  1. Define an Interface abstraction wrapper around the target component.
  2. Modify the codebase to consume this Interface.
  3. Build the new component implementation matching this Interface.
  4. Swap the active class binding in the dependency injector container.

### 3. Feature Flags
- **Concept**: Toggle runtime execution paths dynamically using configuration flags.
- **When to Use**: Releasing high-risk refactored code safely.
- **Execution**: Run the refactored logic for a small percentage of users, and monitor error rates before rolling out the change to all users.

---

## 11. Avoiding Overengineering

We value simplicity. Do not introduce:
- **Patterns Without Problems**: Do not apply design patterns (like Command Query Responsibility Segregation or Repository patterns) unless the codebase has structural problems that require them.
- **Extra Layers Without Need**: Do not insert intermediate service boundaries between simple CRUD controllers and Eloquent models.
- **Microservices Unnecessarily**: Keep applications structured as Monoliths or Modular Monoliths until deployment scale demands distributed services.

---

## 12. Performance & Security Checks

- **Run Performance Baselines**: Measure database query counts, memory allocation, and response times before and after refactoring. Clean code must not degrade performance.
- **Preserve Security Controls**: Ensure validation checks, authorization policies, rate limiters, and CSRF protection layers remain active throughout the refactoring process. Never disable safety controls to simplify code.

---

## 13. Decision Matrices

Use these matrices to identify the correct refactoring decision based on project context.

### Matrix 1: Refactor vs. Rewrite
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Core features require upgrades; code has low-to-medium complexity | **Refactor** | Maintains working business logic; low deployment risk. |
| Obsolete technology stack, severe architectural decay, zero tests | **Rewrite** | Defer only when refactoring cost exceeds rewrite time. |

### Matrix 2: Extract Class vs. Keep Together
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Class manages multiple unrelated tasks (e.g. database + email) | **Extract Class** | Enforces single-responsibility; simplifies unit testing. |
| Cohesive logic fields that are always modified together | **Keep Together** | Avoids unnecessary class separation and directory noise. |

### Matrix 3: Service Layer vs. Direct Logic (Inline)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Third-party API integration (Stripe, Twilio) or shared logic | **Service Layer** | Decouples external library dependencies from models. |
| Database updates that are unique to a single controller | **Direct Logic** | Keeps simple CRUD operations readable. |

### Matrix 4: Repository Pattern vs. Active Record ORM Direct Queries
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Multi-storage engines execution (SQLite + remote APIs) | **Repository** | Abstrakt data access details behind clean interfaces. |
| Standard relational database CRUD query operations | **Active Record** | Leverages Laravel's native Eloquent optimization tools. |

### Matrix 5: Interface vs. Concrete Class
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Multiple implementations exist or are planned (e.g., Mailers) | **Interface** | Decouples caller code from concrete implementations. |
| Single concrete implementation that is unlikely to change | **Concrete Class** | Avoids creating empty interface wrappers that add search noise. |

### Matrix 6: Component Extraction vs. Keep Component Unified
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Large template layout list containing nested forms and dialogs | **Extract Component**| Simplifies testing and reuse of layout modules. |
| Simple rendering screen with low HTML markup lines | **Keep Unified** | Low complexity; avoids directory clutter. |

---

## 14. AI Refactoring Rules

AI agents refactoring code in this repository must follow these rules:

1. **Submit Minimal Changes**: Propose targeted, step-by-step refactors. Avoid large changes that modify multiple unrelated files.
2. **Never Apply Patterns Blindly**: Do not suggest introducing Repository patterns, Factories, or Interfaces unless the existing code violates architectural boundaries.
3. **No Behavioural Alterations**: Verify that all inputs, validation limits, return types, and exceptions remain identical during refactoring.
4. **Preserve Security Checkpoints**: Ensure route authorization checks and validation requests are not removed during refactoring.
5. **Flatten Control Structures**: Flat code structures using guard clauses and early exits.

---

## 15. Refactoring Review Checklist

Use this checklist during code review to evaluate refactoring and migration changes.

### Before Refactoring
- [ ] Has the refactoring goal been defined?
- [ ] Has the target code been audited for security, performance, and complexity?
- [ ] Is there an active test safety net covering the code under refactor?

### During Refactoring
- [ ] Are code modifications applied incrementally in small commits?
- [ ] Do variables, methods, and classes use intention-revealing names?
- [ ] Are control structures flattened using guard clauses (early exits)?
- [ ] Has overengineering been avoided (no unnecessary interfaces or abstractions)?

### After Refactoring
- [ ] Do all unit and integration tests pass successfully?
- [ ] Have performance latency metrics been verified (no N+1 query regressions)?
- [ ] Are all authorization policy checks intact?
- [ ] Have documentation and playbook manifests been updated?

---

## References
- Universal Coding Standards: [05-universal-coding-standards.md](05-universal-coding-standards.md)
- Safe Database Migrations Expand-Contract: [06-database-engineering-standard.md](06-database-engineering-standard.md)
- Test Harness Verification: [11-testing-engineering-standard.md](11-testing-engineering-standard.md)
- Git atomic commits conventions: [12-git-and-collaboration-standard.md](12-git-and-collaboration-standard.md)
