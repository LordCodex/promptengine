---
document_id: core-testing-engineering-standard
title: Testing Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-security-testing-and-threat-modeling
  - core-performance-engineering-standard
  - stacks-php-conventions
  - stacks-laravel-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Testing Engineering Standard

## Purpose & Inheritance
This document defines the core standards for software testing and QA automation. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), the [Database Engineering Standard](06-database-engineering-standard.md), the [API Engineering Standard](07-api-engineering-standard.md), and the [Security Engineering Standard](08-security-engineering-standard.md). It establishes testing protocols across backend API layers, frontends (Vue/Nuxt), mobile codebases (Flutter/Dart), database persistence layers, and CI/CD pipelines.

---

## 1. Testing Philosophy

Testing is a tool for building **confidence and regression prevention**, not a competition to hit high metrics. 

### Core Focus Areas
- **Behavior Over Implementation**: Test what your code *does* (its public interface and behavior), not *how* it does it (its internal private methods and variables). Testing private variables makes refactoring difficult and creates fragile test suites.
- **Abolish Coverage Obsession**: Chasing 100% code coverage leads to writing meaningless assertions for trivial getter/setter routines. Focus your testing budget on critical paths where failures degrade user experience or cause business loss.
- **Fast Feedback Loop**: A slow test suite is a suite that developers avoid running. Unit tests must execute in milliseconds; integration tests must be isolated and highly optimized.

---

## 2. Testing Strategy Decision Framework

Before writing a single test file, developers and AI agents must analyze the change profile using this framework:

```text
Identify Risk ──> Determine Business Impact ──> Assess Change Frequency ──> Select Target Test Level
```

### Risk & Priority Zones
Prioritize testing coverage on resources matching these features:
- **Core Authentication & Authorization Boundaries**: Login flows, token handshakes, role gates.
- **Financial Operations & Calculations**: Ledger writes, payment calculations, refunds.
- **Data Integrity Constraints**: Unique schema pivots, foreign keys, transaction rollbacks.
- **External Integrations**: Third-party API clients, payment webhook handlers.

---

## 3. The Test Pyramid

We balance confidence against execution cost by structuring our test suites according to the Test Pyramid:

```text
    / \
   /   \      End-to-End (E2E) (Full user journeys, slow execution, high cost)
  / E2E \
 /-------\
/  INTEG  \   Integration (Database transactions, API routers, network fakes)
/---------\
/  UNIT   \  Unit (Pure business logic, math calculations, rapid execution)
/___________\
```

### 1. Unit Tests
- **Purpose**: Validate single-responsibility business logic, pure algorithms, and math helpers in complete isolation.
- **No Side-Effects**: Unit tests must not access network sockets, read local databases, or manipulate file systems.
- **What Belongs Here**:
  - Currency rounding converters (e.g., converting dollars to cents).
  - Regular expression validations (e.g., checking password strength).
  - Business rules wrapped inside pure domain actions.

### 2. Integration Tests
- **Purpose**: Verify that multiple system layers work together (e.g., database tables mapping through models, API controllers checking policy files).
- **When to Use**: Integration tests provide higher value than unit tests when verifying database triggers, transaction rollbacks, API contract schemas, and network payload maps.

### 3. End-to-End (E2E) Tests
- **Purpose**: Validate full user journeys across the entire application interface (browser or mobile emulator to database and back).
- **Limits**: E2E tests are slow and expensive to maintain. Restrict E2E tests to core business paths (e.g., "User logs in, adds item to cart, pays, and views invoice").

---

## 4. Test Types & Scopes

- **Functional Testing**: Verifying that the system satisfies business logic rules (does it behave according to specification?).
- **Regression Testing**: Writing assertions around previously fixed bugs to prevent regression issues from returning.
- **Smoke Testing**: Lightweight integration tests run immediately post-deployment to verify that core endpoints are reachable (e.g., returning `200 OK` on `/health`).
- **Performance Testing**:
  - *Load testing*: Measuring system behavior under expected normal user loads.
  - *Stress testing*: Identifying the point at which the application crashes under extreme load.
  - *Endurance testing*: Checking for memory leaks by running sustained traffic over long periods.
- **Security Testing**: Verifying authorization policies, validation rules, and input injection prevention.
- **Contract Testing**: Ensuring API schemas match frontend and mobile client consumer expectations (preventing payload serialization breakages during backend deploys).

---

## 5. Test Design & Mocking Strategy

### Test Design Rules
- **Independence**: Every test must run independently. No test can depend on the execution results or database states of a preceding test.
- **Determinism**: Eliminate flaky tests. Avoid using dynamic clocks (`now()`), random variables, or network dependencies that can cause tests to fail randomly.
- **Clean Setups**: Avoid writing hundreds of lines of mock configurations. If your test setup is complex, your class is likely handling too many responsibilities.

### Mocking Strategy Guidelines
*Mocking* is replacing a real dependency with a controlled test double to isolate the code under test.

- **Mock External APIs**: Always mock third-party network APIs (e.g., Stripe, SendGrid, Twilio) using HTTP fakes. Never make real outgoing network calls during tests.
- **Do Not Mock Databases**: Do not mock your ORM database layer (active record models). Mocking model classes creates fragile, unmaintainable mocks. Use an in-memory database or run tests inside transaction rollbacks.
- **Mock Time Dynamically**: When testing temporal logic (e.g., checking if a subscription expires after 30 days), freeze or travel time using framework time helper utilities:
  ```php
  // Good: Freezing time guarantees deterministic assertions
  $this->travelTo(now()->addDays(30));
  ```

---

## 6. Test Data & Database Testing

### Test Data Management
- **Factories Only**: Use Model Factories to generate relational test data dynamically. Never hardcode data parameters inside seed files or SQL scripts.
- **Realistic Values**: Avoid using placeholder names (e.g., `test1`, `foo`, `bar`). Enforce realistic names, emails, and numbers using mock data generators (like Faker) to catch parsing edge cases.

### Database Testing Protocols
- **Transaction Isolation**: Wrap integration tests that interact with databases inside database transactions. Roll back all database changes automatically at the end of each test to prevent state pollution.
- **Constraint Testing**: Write assertions to verify that unique constraints and foreign key policies are enforced in the database:
  ```php
  // Good: Verify database prevents duplicate email signups
  expect(fn() => User::factory()->create(['email' => 'duplicate@example.com']))
      ->toThrow(QueryException::class);
  ```

---

## 7. API, Frontend & Mobile Testing

### API Testing
- **Request Validation**: Test that inputs are validated correctly (e.g., payload size bounds, field rules).
- **Authorization & Authentication**: Assert that routes return `401 Unauthorized` for missing tokens, and `403 Forbidden` for users lacking required permissions.
- **Response Schemas**: Verify that the JSON response structure matches your documentation schemas.

### Frontend Testing (Vue / Nuxt)
- **Component Behaviors**: Test component outputs based on input properties and user interactions.
- **Do Not Test Framework Internals**: Do not assert component internal states or private methods. Focus on the rendered DOM outputs and emitted events.
- **Accessibility (A11y) Audits**: Run automated contrast and semantic markup validation checks (e.g., using `axe-core`) on core components.

### Mobile Testing (Flutter)
- **Widget Testing**: Validate user interaction loops (e.g., tapping a button changes page states).
- **Platform Mocking**: Mock platform-specific channels (such as secure keychains or native hardware access layers) before executing tests.

---

## 8. Testing Financial Systems

Financial operations require zero-tolerance testing protocols.

### Financial Testing Rules
- **Decimal Precision Audits**: Assert calculation outcomes using strict value comparisons to verify that decimals are not rounded incorrectly.
- **Double-Entry Balance Verification**: Assert that the sum of debit and credit ledger records equals the balance checkpoint.
- **Race Condition Simulations**: Write concurrent integration tests to verify that balance checks cannot be bypassed during double-spending attempts:
  ```php
  // Verify that balance cannot go negative during concurrent withdraw requests
  ```
- **Duplicate Request Guards**: Test idempotency keys on payment routes to verify that repeating an API request does not charge the customer twice.

---

## 9. Legacy Code Testing

When working on legacy systems that lack tests, follow this incremental coverage workflow:

```text
Identify Hot Path ──> Write Characterization Tests ──> Apply Refactor ──> Verify Against Baseline
```

1. **Characterization Testing**: Before modifying a legacy file, write tests that assert its *current* behavior (even if it has bugs). This establishes a safety baseline.
2. **Do Not Rewrite for Test Coverage**: Do not rewrite working legacy systems just to add tests. Refactor only when you are actively modifying the file to add a feature or fix a bug.
3. **Risk-Based Testing**: Prioritize writing integration tests for legacy paths that touch database transactions, authentication, or payment processing.

---

## 10. Continuous Integration (CI) Testing

- **Pipeline Automation**: Run the entire unit, integration, and contract test suites on every pull request.
- **Fail Fast**: Run fast unit tests first in your CI configuration. If a unit test fails, terminate the build pipeline instantly to save resource costs.
- **Parallelization**: Configure CI runners to run tests in parallel using isolated test databases.
- **Flaky Test Isolation**: If a test fails randomly (flaky test), do not ignore it. Fix it immediately, or disable it temporarily while documenting a fix ticket.

---

## 11. Testing Anti-Patterns

- **Testing Mocks Instead of Behavior**: Setting up complex mocks that duplicate the application code's implementation details, causing tests to pass even when the code is broken.
- **Slow Integration Suites**: Running heavy database queries or network fakes inside loops in unit test folders, inflating test suite execution times.
- **Snapshot Overuse**: Relying on snapshot testing for dynamic, frequently changing pages, leading developers to regenerate snapshots blindly without verifying layout changes.
- **Obsession with Coverage Metrics**: Enforcing a strict 100% coverage rule, which incentivizes developers to write low-quality tests for code that doesn't need testing.

---

## 12. Decision Matrices

Use these matrices to identify the correct testing approach based on project context.

### Matrix 1: Unit vs. Integration Test
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Pure data translations, regex validators, calculations | **Unit Test** | Rapid execution speed (milliseconds), zero side-effects. |
| Model interactions, SQL query validation, API routes gates | **Integration Test** | Verifies database operations, middleware, and framework logic. |

### Matrix 2: Mock vs. Real Dependency
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Third-party external endpoints, payment APIs | **Mock (Fake)** | Prevents real network overhead, avoids charges, and ensures speed. |
| Local database tables, system config cache stores | **Real Dependency** | Mocking database models leads to fragile, unrealistic test setups. |

### Matrix 3: Manual vs. Automated Test
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Core business features, regression validation, security checks | **Automated Test** | Fast verification on every commit, preventing regressions. |
| Visual layouts rendering checks, initial design reviews | **Manual Test** | Low setup overhead for highly subjective, fluid UI reviews. |

### Matrix 4: E2E vs. Feature Test
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Multi-role workflows, complex browser actions (drag-and-drop) | **E2E Test** | Validates the actual browser rendering and full system integration. |
| Single API resource validation, request inputs validations | **Feature Test** | Fast execution; does not require browser boot overhead. |

### Matrix 5: High Coverage vs. Critical Coverage
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Standard CRUD applications with low complexity | **Critical Coverage** | Focuses testing resources on security borders, payments, and data safety. |
| Reusable system utility packages and developer tools | **High Coverage** | Ensures that utilities are robust across all input variations. |

### Matrix 6: Snapshot Testing vs. Behaviour Testing
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Fast UI regression checks on complex static components | **Snapshot Testing** | Catches unexpected markup changes automatically. |
| Component interactive states, user inputs transitions | **Behaviour Testing** | Verifies actual widget functionality and state modifications. |

---

## 13. AI Testing Rules

AI agents modifying or writing code in this repository must follow these rules:

1. **Write Accompanying Tests**: Every new feature or bug fix must include corresponding test cases verifying the change behavior.
2. **Never Suggest Hardcoded Seeds**: Do not write hardcoded arrays inside tests. Suggest using factories for all database record generations.
3. **No Flaky Clocks**: Enforce time freezing or time traveling helpers when writing tests that check dates or expiration times.
4. **Mock External Outgoing Calls**: Ensure all HTTP, Mail, and Queue facades are faked at the start of feature tests.
5. **Never Delete Failing Tests**: If a test fails after your modification, do not delete it. Analyze if your change introduces a regression or breaks backward compatibility.

---

## 14. Testing Review Checklist

Use this checklist during code review to evaluate testing implementations.

### Quality & Speed
- [ ] Do tests run independently of each other (no state sharing)?
- [ ] Are external network calls properly mocked (no real API requests)?
- [ ] Is the test suite fast (no slow sleeps or timeouts)?

### Data & Isolation
- [ ] Are test records generated dynamically using factories (no hardcoded seeds)?
- [ ] Are database tests run inside isolated transactions (rolled back on finish)?

### Business Logic & Edge Cases
- [ ] Do tests verify validation constraints and error states (negative values, empty limits)?
- [ ] For financial calculations, are values asserted using exact calculations (no loose comparisons)?
- [ ] Are race conditions and duplicate operations guarded against?

### Security & CI
- [ ] Are authentication failures and authorization permissions verified?
- [ ] Do all tests run successfully on your local runner and the CI pipeline?

---

## References
- Universal Naming Rules: [05-universal-coding-standards.md](05-universal-coding-standards.md)
- Database Optimizations: [06-database-engineering-standard.md](06-database-engineering-standard.md)
- API Testing Controls: [07-api-engineering-standard.md](07-api-engineering-standard.md)
- Security Testing and Threat Modeling: [09-security-testing-and-threat-modeling.md](09-security-testing-and-threat-modeling.md)
