---
document_id: stacks-react-testing
title: React Testing Conventions
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# React Testing Conventions

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines component testing standards.

---

## 1. Testing Behaviors Over Implementation

- **React Testing Library (RTL)**: Test component outputs from the user's perspective. Rely on `screen.getByRole()` or `screen.getByText()` query hooks rather than verifying internal component state.
- **Simulate User Clicks**: Use `@testing-library/user-event` to simulate pointer interactions, typing, and selections rather than calling trigger handlers directly.

---

## 2. Mocking Strategies

- **MSW (Mock Service Worker)**: Intercept client-side fetch networks at the network layer. Avoid mock wrappers around fetch variables directly.
- **Custom Hooks Mocks**: Mock hooks output structures (`useAuth`) when testing visual components in isolation, ensuring state paths are covered.

---

## 3. Testing Async Operations and Hooks

- **Async Actions**: Wrap assertions querying elements appearing after network requests inside `await screen.findByRole()` or `waitFor()`.
- **Render Custom Hooks**: Use `renderHook()` to test standalone helper hooks in isolation.

---

## Review Checklist

Verify test suites against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] React Testing Library query hooks match role mappings.
- [ ] User actions are triggered via `user-event` APIs.
- [ ] Mock Service Worker manages REST endpoint mocks.
