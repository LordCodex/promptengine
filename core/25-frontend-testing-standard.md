---
document_id: core-frontend-testing
title: Frontend Testing Standard
ecosystem: cross-cutting
dependencies:
  - core-testing-philosophy
  - core-frontend-architecture
  - core-frontend-security
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Frontend Testing Standard

## Purpose & Inheritance
This document defines client-side testing methodologies, verification boundaries, and tooling guidelines. It inherits from and extends the [Testing Philosophy and Patterns](04-testing-philosophy.md) and the [Frontend Architecture Standard](23-frontend-architecture-standard.md), ensuring that frontend code is regression-resistant, maintainable, and verifiable.

---

## 1. The Core Testing Principle

**Test what users do, not implementation details.**

A test suite must assert user-facing behavior over private, code-level execution details.
- **Bad**: Asserting that a component's internal helper method `parseCouponCode()` is called.
- **Good**: Asserting that a user typing an invalid coupon code sees the inline validation error message `"Invalid coupon code"`.

Testing internal class properties or private component methods makes refactoring difficult and creates brittle test suites. Test interactions and DOM output changes instead.

---

## 2. Frontend Test Pyramid

We balance testing confidence against execution time and maintenance cost by structuring our client-side test suites into four layers:

```text
       /\
      /  \      End-to-End (E2E) (5%) - Critical user journeys (Playwright, Cypress)
     /----\
    /      \    Integration Tests (15%) - Multi-component flows & mock API interactions
   /--------\
  /          \  Component Tests (50%) - UI state changes, forms, validation (RTL, VTU, Widget tests)
 /____________\ Unit Tests (30%) - Stateless helpers, data formatting, business math (Vitest, Jest)
```

### 2a. Unit Tests
- **Focus**: Pure functions, stateless utility helper code, data transformations, and input validation algorithms.
- **Examples**: Currency formatter calculations, relative date parsing, permission checker filters.
- **Constraint**: No browser dependencies, no API calls, and zero external I/O. Must execute in milliseconds.

### 2b. Component Tests
- **Focus**: Isolated component rendering, user click events, prop variation behaviors, local state modifications, and component accessibility attributes.
- **Examples**: Form field validity, modal open/close states, input character counting, custom dropdown search filters.
- **Constraint**: Mount only the target component and stub heavy downstream children. Mock API clients.

### 2c. Integration Tests
- **Focus**: Coordination of multiple component layers, client-side route transitions, global store state changes, and simulated API payloads.
- **Examples**: Add-to-cart workflow, dashboard filter updates, multi-step checkout processes.
- **Constraint**: Use mock service workers (e.g. MSW) to intercept and simulate network calls rather than hitting real staging servers.

### 2d. End-to-End (E2E) Tests
- **Focus**: Critical user journeys verifying correct integration between frontend client, backend API databases, and external third-party services.
- **Examples**: Complete registration pipelines, login-to-purchase success flow, password recovery verification.
- **Constraint**: Restrict E2E tests to high-risk business flows to prevent test execution bloat.

---

## 3. Test Cases for Critical User Flows

All critical application paths must have automated test cases covering both success and failure vectors:

### 3a. Authentication
- [ ] **Successful Login**: Verify valid inputs redirect to the intended dashboard.
- [ ] **Invalid Credentials**: Verify dynamic error states announce failed password attempts.
- [ ] **Session Expiration**: Verify UI redirects to login when requests return a `401 Unauthorized` response.
- [ ] **Route Guards**: Verify direct navigation to protected paths is blocked for guest sessions.

### 3b. Forms and Input
- [ ] **Valid Submission**: Verify submit button triggers action when all required fields are filled.
- [ ] **Invalid Input**: Verify inputting illegal formats blocks submission and shows error states.
- [ ] **Required Fields**: Verify blank required fields display visual error warnings.
- [ ] **Server Validation**: Verify API validation error payloads are rendered next to the correct form inputs.
- [ ] **Network Failures**: Verify form data is preserved and a retry message appears on network loss.

### 3c. Payments
- [ ] **Amount Display**: Verify correct billing amounts and currency symbols are rendered.
- [ ] **Success Flow**: Verify checkout completion displays a transaction success message.
- [ ] **Failure Flow**: Verify decline errors (insufficient funds, invalid code) display fallback recovery actions.
- [ ] **Double Submit Prevention**: Verify the submit button is disabled during process execution to prevent duplicate charges.

> [!CAUTION]
> **Never test with real payment credentials.** Use mock tokens, sandbox environments, or client-side stubs.

### 3d. Data Loading and States
- [ ] **Loading State**: Verify skeleton views or loading spinners display during data retrieval.
- [ ] **Empty State**: Verify an informative empty message appears when API responses return empty lists.
- [ ] **Success State**: Verify data renders correctly when the payload arrives.
- [ ] **Failed State**: Verify an error message and a retry button appear on API request failure.

---

## 4. Mocking Rules

- **Mock External Boundaries**: Mock API gateways, third-party payment providers, and tracking scripts (e.g. Stripe, Google Analytics).
- **Do Not Over-Mock**: Do not mock internal child components, stores, or state managers unless they perform network fetches or CPU-intensive processes.
- **Use Network Interception**: Use libraries like Mock Service Worker (MSW) for web testing to mock network communication at the network layer rather than mocking HTTP client classes.

---

## 5. API and Error Testing

### API Interface Coverage
API tests must verify client robustness:
- Verify that request parameters (headers, search queries, payload bodies) match API specifications.
- Verify client handling of non-2xx responses (`400 Bad Request`, `403 Forbidden`, `422 Unprocessable Entity`, `500 Server Error`).
- Verify handler recovery on unexpected JSON body structures.

### Error Case Testing
Every user-facing feature must explicitly test client behavior during system errors:
- **Network loss**: Trigger offline conditions (`navigator.onLine = false`).
- **Server timeout**: Simulate request timeouts and verify user warning dialogs.
- **Permission errors**: Simulate `403` responses and verify access-denied panels.

---

## 6. Accessibility & Responsive Verification

### Accessibility (a11y) Assertions
Verify client compliance with accessibility guidelines programmatically:
- Assert that interactive components support keyboard focus (`Tab`) and execution (`Enter`, `Space`).
- Assert that focus is managed during modal popups (focus moves inside, return on close).
- Assert that all inputs have linked labels.
- Assert correct ARIA roles on custom elements (e.g., `aria-expanded` toggles correctly).

### Responsive Assertions
Test important application screens at different viewport sizes:
- **Mobile**: Test touch target clickability (minimum 44×44px), vertical scrolling, and collapsing menu buttons.
- **Tablet**: Test intermediate grid/layout transitions.
- **Desktop**: Test column rendering and content density.

---

## 7. Visual Regression Testing

Use visual screenshot comparisons to protect against silent CSS regression:
- Apply visual testing to **shared design components** (buttons, form inputs, alerts) and **critical structural pages** (checkout, registration, billing).
- Run visual verification across multiple viewports to detect media query layout breakage.
- **Avoid snapshotting everything**: DOM snapshot testing (`expect(markup).toMatchSnapshot()`) creates brittle test suites that break on trivial whitespace or markup changes. Use visual screenshot assertions (e.g. Playwright screenshots) instead.

---

## 8. Stack-Specific Testing Guides

### 8a. Vue 3 / Nuxt
- **Vitest**: Use as the test runner for fast unit and component tests.
- **Vue Test Utils (VTU)**: Use to mount components, trigger input interactions, and assert emitted custom events.
- **Mocking Pinia**: Use `@pinia/testing` to mock shared global store states easily.
- **Playwright**: Use for E2E user flow tests and responsive layout checks.

### 8b. React / Next.js
- **Jest / Vitest**: Use for fast test execution.
- **React Testing Library (RTL)**: Use for component testing. Query elements by accessible roles (e.g. `screen.getByRole('button', { name: /save/i })`) to enforce accessibility testing automatically.
- **Playwright**: Use for integration and E2E browser verification.

### 8c. Flutter / Dart
- **Unit Tests**: Test business logic classes, repositories, and state models.
- **Widget Tests**: Test layout rendering, user touch interactions, and widget state updates.
- **Integration Tests**: Verify end-to-end device behavior.
- **Visual Goldens**: Use golden image tests (`matchesGoldenFile`) to verify visual component consistency.

---

## 9. Test Quality and Selection Rules

### Avoid Brittle Selectors
Never query DOM elements using class names that change frequently or are auto-generated:
- **Bad**: `wrapper.find('.btn-primary-blue-large')`
- **Bad**: `wrapper.find('.sc-bdVaJa')` (styled-components auto-generated class)
- **Good**: Query by accessibility role: `screen.getByRole('button', { name: /submit/i })`
- **Good**: Use custom data attributes designated for testing: `wrapper.find('[data-testid="submit-button"]')`

### Test Intentional Edge Cases (AI requirements)
Before completing a task, verify the test suite handles these AI quality gates:
1. **What can fail?** (Network offline, invalid inputs, unauthorized states)
2. **What degrades user experience?** (Missing loading states, missing error panels)
3. **What security boundaries exist?** (Attempting unauthorized operations, route guards)
4. **What boundary inputs exist?** (Zero values, empty strings, character limit overflows)

---

## Test Review Checklist

Verify the frontend work against this testing checklist before shipping:
- [ ] **Critical flows tested**: Are registration, login, payment, and critical business flows verified by automated tests?
- [ ] **Loading states tested**: Do component tests confirm that skeletons or loading indicators appear during API fetches?
- [ ] **Error states tested**: Are network failures, validation errors, and timeout recoveries covered by assertions?
- [ ] **Empty states tested**: Do data list components verify correct messaging when response datasets are empty?
- [ ] **Forms tested**: Are valid inputs, invalid formats, and disabled submit states verified?
- [ ] **Accessibility considered**: Do tests verify keyboard navigation, focus management, and label presence?
- [ ] **Responsive behavior considered**: Are key layouts tested at multiple viewports?
- [ ] **API failures handled**: Do tests assert application resilience against non-2xx server payloads?
- [ ] **No brittle implementation tests**: Are all assertions testing user-visible behavior rather than private methods?
- [ ] **Tests explain user behavior**: Are test block names written from the user's perspective (e.g. `should show error when billing card fails`)?

---

## References
- Testing Philosophy: [core/04-testing-philosophy.md](04-testing-philosophy.md)
- Frontend Architecture: [core/23-frontend-architecture-standard.md](23-frontend-architecture-standard.md)
- Frontend Security: [core/24-frontend-security-and-privacy-hardening-standard.md](24-frontend-security-and-privacy-hardening-standard.md)
- Playwright: [https://playwright.dev](https://playwright.dev)
- MSW (Mock Service Worker): [https://mswjs.io](https://mswjs.io)
