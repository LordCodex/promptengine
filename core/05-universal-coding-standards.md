---
document_id: core-universal-coding-standards
title: Universal Coding Standards
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Universal Coding Standards

This document defines the universal coding standards for all codebases. Framework-specific guides (PHP, Laravel, Dart, Vue, etc.) extend this document and inherit its principles.

---

## 1. Naming Conventions and Taxonomy

### Purpose
To establish a predictable, self-documenting naming convention across all files, classes, variables, and database objects.

### Why
Names represent the semantic intent of code. Vague or inconsistent names force developers and AI agents to analyze the internal execution paths of classes to understand their purpose, creating cognitive friction.

### Benefits
- Code becomes readable as prose.
- AI agents locate target logic paths instantly.
- Minimizes naming discrepancies between separate team branches.

### Trade-offs
Descriptive, explicit names can be longer. This requires slightly more typing but is offset by modern IDE autocomplete features.

### Anti-patterns & Common Mistakes
- **Type-Encoding (Hungarian Notation)**: Naming variables with their primitive type (e.g. `array_users`, `stringName`).
- **Meaningless Abbreviations**: Using single letters or truncated terms (e.g. `c` for customer, `msg` for message) that have multiple interpretations.

### Good Examples
- `UserRepository` (cohesive class name expressing domain and pattern)
- `isTransactionPending` (boolean expressing state clearly)
- `calculateOrderTax()` (verb-noun combination for functions)

### Bad Examples
- `DataClass` (vague noun)
- `flag` (unclear boolean status)
- `doWork()` (generic function action verb)

### Review Checklist
- [ ] Are class names UpperCamelCase and represented by clear nouns?
- [ ] Are functions and variable names lowerCamelCase and descriptive?
- [ ] Do Boolean variables begin with state-checking verbs (e.g., `is`, `has`, `should`)?

### References
- Database Modeling: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)

---

## 2. Code Readability and Cognitive Load

### Purpose
To reduce the mental bandwidth required for developers and AI agents to parse logical flows.

### Why
Complexity accumulates incrementally. Code that utilizes clever language tricks, deep nests, or dense Boolean conditions is prone to maintenance errors and logical omissions.

### Benefits
- Lowers Mean Time to Resolution (MTTR) on production bugs.
- Simplifies writing automated tests.

### Trade-offs
Flattening nested loops and conditions sometimes requires creating extra variables or splitting code into smaller helper methods, increasing overall file line count.

### Anti-patterns & Common Mistakes
- **The Arrow Anti-Pattern**: Nesting multiple `if` conditions, loops, and callbacks, creating deep indentation corridors.
- **Clever One-Liners**: Writing complex mathematical evaluations or recursive statements on a single line to minimize code space.

### Good Examples
```typescript
function chargeCustomer(user, amount) {
  if (user === null) return;
  if (!user.isBillingActive) throw new InactiveBillingException();
  if (amount <= 0) throw new InvalidAmountException();

  gateway.charge(user.paymentMethodToken, amount);
}
```

### Bad Examples
```typescript
function chargeCustomer(user, amount) {
  if (user !== null) {
    if (user.isBillingActive) {
      if (amount > 0) {
        gateway.charge(user.paymentMethodToken, amount);
      }
    }
  }
}
```

### Review Checklist
- [ ] Is code nesting limited to a maximum of **three levels**?
- [ ] Are complex conditions extracted to named Boolean helper variables?
- [ ] Did you use Guard Clauses to handle error/edge paths first?

### References
- Planning Phase: [core/01-thinking-and-planning.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/01-thinking-and-planning.md)

---

## 3. Comments and Documentation Boundaries

### Purpose
To define the boundary between code that documents itself and necessary inline technical comments.

### Why
Obvious comments (restating what the code does) add visual clutter. Outdated comments that do not match modified source code mislead developers and AI agents.

### Benefits
- Keeps source files clean and readable.
- Ensures documentation remains accurate since only complex, stable rules are documented.

### Trade-offs
Relying entirely on self-documenting code can sometimes obscure complex business requirements that cannot be inferred from function signatures.

### Anti-patterns & Common Mistakes
- **The Inline Translator**: Writing comments for every single line of standard code.
- **Tombstones**: Commenting out unused blocks of code instead of deleting them.

### Good Examples
```typescript
// Stripe API requires values in cents (integers) to prevent
// floating-point precision rounding errors on currency calculations.
const paymentAmountInCents = Math.round(amount * 100);
```

### Bad Examples
```typescript
// Increment index by 1
index++;
```

### Review Checklist
- [ ] Does the comment explain the *why* (business logic/decisions) rather than the *what* (syntax)?
- [ ] Are all commented-out code blocks deleted?

### References
- Project conventions: [CONTRIBUTING.md](file:///Users/kodexkode/Documents/workspace/promptengine/CONTRIBUTING.md)

---

## 4. Function Design

### Purpose
To ensure functions are modular, testable, and have isolated execution boundaries.

### Why
Large, multi-purpose functions are prone to logical side-effects and are difficult to cover with clean unit tests.

### Benefits
- Enhances code reusability.
- Simplifies writing unit tests by minimizing dependencies.

### Trade-offs
Splitting tasks into smaller functions can create extra function call definitions and trace overhead.

### Anti-patterns & Common Mistakes
- **Flag Parameters**: Passing boolean flags that change the internal execution path of the function.
- **Hidden Side-Effects**: Modifying global variables or parent scope parameters without returning a value.

### Good Examples
```typescript
function calculateDiscount(price: number, isPremium: boolean): number {
  return isPremium ? price * 0.15 : 0;
}
```

### Bad Examples
```typescript
function calculateDiscountAndSave(price: number, isPremium: boolean, user: User): number {
  const discount = isPremium ? price * 0.15 : 0;
  user.discountCount++;
  db.save(user); // Hidden database write side-effect
  return discount;
}
```

### Review Checklist
- [ ] Does the function do exactly **one thing**?
- [ ] Is the function length under **30 lines** of execution code?
- [ ] Are there **3 or fewer** parameters?

### References
- Testing functions: [core/04-testing-philosophy.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/04-testing-philosophy.md)

---

## 5. Class Design

### Purpose
To build robust, object-oriented logical entities that prevent internal state corruption.

### Why
Improper encapsulation and rigid inheritance lead to fragile classes that break when external parameters change.

### Benefits
- Enforces strict encapsulation of state.
- Simplifies class substitution (Liskov Substitution Principle).

### Trade-offs
Requires declaring private properties and public getters/setters, slightly increasing class boilerplate code.

### Anti-patterns & Common Mistakes
- **God Classes**: Writing monolithic classes that hold thousands of lines of logic and manage multiple unrelated concerns.
- **Deep Inheritance Chains**: Extending multiple levels of base classes, causing subclass behavior to become unpredictable.

### Good Examples
```typescript
class BankAccount {
  private _balance: number = 0;

  public get balance(): number { return this._balance; }

  public deposit(amount: number): void {
    if (amount <= 0) throw new InvalidAmountException();
    this._balance += amount;
  }
}
```

### Bad Examples
```typescript
class BankAccount {
  public balance: number = 0; // Accessible and mutable globally
}
```

### Review Checklist
- [ ] Are class properties private by default, exposing change actions through public methods?
- [ ] Are dependencies injected via constructor arguments rather than initialized inline?

### References
- Simplicity Standards: [core/02-architecture-and-simplicity.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/02-architecture-and-simplicity.md)

---

## 6. Module and Architectural Boundaries

### Purpose
To isolate application layers, enforce high cohesion, and keep business logic decoupled from frameworks.

### Why
If business rules are tightly coupled to framework boundaries (like database engines, routing modules, or UI kits), upgrading or swapping frameworks breaks the application.

### Benefits
- Code remains portable and reusable.
- Upgrading framework versions does not corrupt business logic.

### Trade-offs
Requires mapping parameters between layers (e.g. converting database entity properties to public API resources), creating extra transformer files.

### Anti-patterns & Common Mistakes
- **UI Data Leaks**: Passing database connection objects directly down to UI elements.
- **Smart Controllers**: Writing raw SQL queries inside HTTP route handlers or mobile views.

### Good Examples
- A controller intercepts a request, parses parameters into a DTO, calls a pure Domain Action class, and passes the result to an API Resource class for serialization.

### Bad Examples
- A controller queries the database, formats the result into raw HTML strings, and echoes the output directly inside the method execution.

### Review Checklist
- [ ] Does the business layer remain free of framework dependencies?
- [ ] Do communications between layers cross borders using defined contracts or primitive types?

### References
- API contract patterns: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)

---

## 7. Code Duplication and Abstraction

### Purpose
To define when to dry up code and when to tolerate duplication to preserve simplicity.

### Why
Premature abstraction creates highly coupled classes that are difficult to modify. Tolerating mild duplication is often cheaper than maintaining a wrong abstraction.

### Benefits
- Prevents complex dependency trees.
- Keeps code modules simple and self-contained.

### Trade-offs
Changes to a duplicated code segment must be applied to all copies, which increases developer maintenance efforts if the duplication is widespread.

### Anti-patterns & Common Mistakes
- **Premature DRYing**: Merging two functions that share identical structures but have entirely different domain responsibilities.
- **Utility Dump**: Creating a massive `utils.ts` or `Helpers.php` file containing hundreds of unrelated functions.

### Good Examples
- Applying the **Rule of Three**: Keep duplicate blocks of logic separate until they appear in at least three distinct contexts before refactoring into a shared helper or service.

### Bad Examples
- Creating a base abstract class to share single-line logic between two entirely unrelated models (e.g. `User` and `Invoice`).

### Review Checklist
- [ ] Is this helper or service abstraction justified by at least three distinct occurrences?
- [ ] Can the duplicated code block be edited without breaking other systems?

### References
- Decoupling rules: [core/02-architecture-and-simplicity.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/02-architecture-and-simplicity.md)

---

## 8. Error Handling and Resilience

### Purpose
To manage failures gracefully, prevent crashes, and log errors securely.

### Why
Uncaught exceptions cause system crashes. Swallowed exceptions hide logical bugs and result in data corruption.

### Benefits
- Protects runtime stability.
- Structured logging details assist in immediate diagnostics.

### Trade-offs
Requires developers to write custom exception classes and map try-catch logic handlers on all boundaries, increasing file count.

### Anti-patterns & Common Mistakes
- **The Empty Catch**: Catching exceptions and doing nothing with them.
- **Generic Exception Raising**: Raising generic `Exception` messages with no typing or parameters.

### Good Examples
```typescript
try {
  await gateway.processPayment(amount);
} catch (error) {
  logger.error('Payment failure occurred', { userId, amount, error });
  throw new PaymentProcessingException('Gateway failure', { parent: error });
}
```

### Bad Examples
```typescript
try {
  gateway.processPayment(amount);
} catch (e) {
  // Silent fallback
}
```

### Review Checklist
- [ ] Are custom exceptions thrown instead of generic class errors?
- [ ] Are secrets (keys, passwords) filtered out of logging payloads?

### References
- Deployment safety: [legacy/03-deployment-risk-reduction.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/03-deployment-risk-reduction.md)

---

## 9. Defensive Programming and Input Boundaries

### Purpose
To safeguard application execution against malformed parameters and corrupted states.

### Why
Data entering the application (via user inputs, database states, third-party APIs) can be invalid. System stability relies on enforcing explicit boundary validation.

### Benefits
- Prevents memory crashes and SQL injections.
- Validates system assertions early.

### Trade-offs
Slightly increases memory usage and CPU cycles due to parameter checking at boundaries.

### Anti-patterns & Common Mistakes
- **Implicit Faith**: Assuming that data stored in databases or received from internal services is always correct and well-formed.
- **Unchecked Nulls**: Performing method calls on objects without asserting null states.

### Good Examples
- Enforcing input parameters type-checking, checking array keys existence, and asserting null states before running inner services.

### Bad Examples
- Reading user parameters from a raw, unvalidated JSON payload directly inside core calculations.

### Review Checklist
- [ ] Are array keys checked for existence before they are queried?
- [ ] Are inputs verified at the absolute boundary of the application?

### References
- API contract payloads: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)

---

## 10. Maintainability and Technical Debt Management

### Purpose
To write code that is clean, modular, and optimized for future developers and AI agents.

### Why
Technical debt slows development, increases bug density, and degrades developer experience. Code must be maintained actively on every commit.

### Benefits
- Keeps the code base highly adaptable to changing business targets.
- Lowers MTTR and developer onboarding times.

### Trade-offs
Refactoring technical debt takes time away from writing new product features.

### Anti-patterns & Common Mistakes
- **Hack-and-Leave**: Writing quick, messy workarounds and promising to refactor them "later".
- **Rewriting Everything**: Rewriting healthy, working modules because of minor style differences.

### Good Examples
- Following the **Boy Scout Rule**: Always leave the code file slightly cleaner than you found it (e.g. adding strict types, cleaning up a condition) when executing updates.

### Bad Examples
- Layering hacks on top of existing spaghetti code to avoid refactoring the root logic.

### Review Checklist
- [ ] Does this pull request decrease the overall technical debt of the modified files?
- [ ] Have all workarounds been documented with links to tracking tickets?

### References
- Safe modernization methods: [legacy/01-safe-refactoring.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/01-safe-refactoring.md)

---

## 11. Code Review Standards

### Purpose
To verify logic correctness, style conformance, and security parameters before merging code changes.

### Why
Code reviews represent the final filter to prevent bugs, logic gaps, and security vulnerabilities from reaching production.

### Benefits
- Promotes uniform coding standards.
- Serves as an educational tool for developers.

### Trade-offs
Requires developer time and blocks branch merges until review iterations complete.

### Anti-patterns & Common Mistakes
- **Formatting Arguments**: Arguing about style parameters (spaces, quotes) during code reviews instead of letting automated linters handle them.
- **Rubber-Stamping**: Approving code changes instantly without reading or testing the changes.

### Good Examples
- Standardizing review checks against a clear checklist, keeping comments focused on logic intent, and testing branches locally.

### Bad Examples
- Approving a 1,000-line pull request in 3 minutes without compiling the codebase.

### Review Checklist
- [ ] **Correctness**: Have the changes been verified and tested locally?
- [ ] **Test Coverage**: Do new files have corresponding tests?
- [ ] **Security**: Are all inputs validated and secrets isolated?

### References
- Pipelines checks: [environment/03-ci-cd-pipelines.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/03-ci-cd-pipelines.md)

---

## 12. Professional Output Quality

### Purpose
To ensure that all code — whether written by a developer or an AI agent — reads as if it was crafted by an experienced, senior engineer and not auto-generated.

### Why
Code that looks "AI-generated" (overly verbose, template-heavy, pattern-obsessed, annotated with obvious comments) degrades team trust and introduces unnecessary complexity. Code that is left with debug outputs or placeholder stubs creates production risk.

### Benefits
- Code passes peer review without hesitation.
- Production incidents are not triggered by leftover debugging artifacts.
- New engineers understand the codebase without needing to decipher AI scaffolding.

### Trade-offs
Requires AI agents and developers to take additional time to review and clean output before submitting.

### Rules

#### Write Boring, Predictable Code
- Write code that reads as if it was authored by an experienced engineer over many years — not generated in one pass.
- Prefer **obvious over impressive**. If two implementations are equally correct, always choose the simpler one.
- Avoid "clever" code, over-abstracted patterns, or architectural experiments not required by the problem.
- Another developer must be able to understand any file **within 30 seconds**.
- Every class must have a **single, clear purpose**. Every method must be easy to understand.
- Code must be **boring and predictable** — surprising behaviour is always a bug.

#### No Emojis in Code
- Do **not** use emojis anywhere in the codebase — not in source code, comments, string literals, log output, or console messages.
- Zero exceptions unless explicitly authorized by the project owner for a specific user-facing string.

#### No Placeholder or Dead Code
- Never generate `TODO` comments, placeholder methods, stub classes, or dead code in production output.
- Never leave commented-out code blocks in the source without an explicit, dated explanation of why they are preserved.
- If a feature is not yet implemented, do not generate its scaffold — wait until the actual requirement is confirmed.

#### No Debug Output in Source
- Remove all debug output before submitting any code:
  - PHP: `var_dump()`, `print_r()`, `dd()`, `dump()`, `ray()`
  - JavaScript/TypeScript: `console.log()`, `console.debug()`, `debugger`
  - Dart: `debugPrint()` in production paths, `print()` statements
- Debug output in production code is a **security risk** (information leakage) and a **quality failure**.

### Anti-patterns & Common Mistakes
- **Scaffold Sprawl**: Generating full module structures with placeholder methods and TODOs for a simple bug fix.
- **Emoji Decoration**: Adding 🚀, ✅, or ⚠️ to log strings or comments for visual appeal.
- **Debug Remnants**: Leaving `console.log('here')` or `dd($variable)` in production-bound code.

### Review Checklist
- [ ] Does the code read as natural, senior-engineer-quality output?
- [ ] Are there zero emojis anywhere in the submitted code?
- [ ] Are there zero TODO/placeholder/stub constructs?
- [ ] Are there zero debug output statements (`dd`, `console.log`, `var_dump`, etc.)?

---

## 13. AI Agent Coding Directives

### Purpose
To guide automated AI coding systems to write clean, context-appropriate, and consistent code in this repository.

### Why
Without strict constraints, AI agents can generate redundant patterns, introduce unnecessary files, or execute destructive rewrites.

### Benefits
- Reduces token consumption and execution costs.
- Enforces codebase pattern consistency.

### Trade-offs
AI agents operate under stricter boundaries, which may require human oversight to approve deviation exceptions.

### Anti-patterns & Common Mistakes
- **Pattern Synthesis**: Introducing a new design pattern or folder structure that does not match existing conventions.
- **Unnecessary Abstractions**: Creating extra interfaces or service layers for a simple bug fix.

### Good Examples
- Querying the `playbook-manifest.json` first, modifying only targeted lines, and following established architectural models.

### Bad Examples
- Generating a complete file rewrite when only modifying a single line of business logic.

### Review Checklist
- [ ] Did the AI agent match variable naming, class structures, and file formatting conventions?
- [ ] Have the changes been verified against the existing test suite?
- [ ] Did the AI agent produce zero TODO/placeholder/stub constructs?
- [ ] Did the AI agent remove all debug output before submitting?

### References
- AI instruction manifests: [playbook-manifest.json](file:///Users/kodexkode/Documents/workspace/promptengine/playbook-manifest.json)
- AI workflow standard: [core/20-ai-agent-engineering-workflow-standard.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/20-ai-agent-engineering-workflow-standard.md)
