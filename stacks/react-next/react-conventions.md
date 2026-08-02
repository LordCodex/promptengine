---
document_id: stacks-react-conventions
title: React Conventions and Component Architecture
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# React Conventions and Component Architecture

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines conventions for React projects in this codebase, avoiding general React tutorials.

---

## 1. Functional Components & Declarative Style

- **Functional Components by Default**: Declare all components using standard `const Component = () => {}` functional forms or traditional function declarations. Never use legacy ES6 class components.
- **Strict Declarative UI**: Keep rendering pure. Do not write side-effects inside components' render phases. Component return blocks must only reflect state projections.

---

## 2. Hooks Best Practices

- **The Rules of Hooks**: Enforce eslint rules for hooks. Never call hooks conditionally or inside loops.
- **Custom Hooks for Logic**: Separate complex state and asynchronous operations into dedicated custom hooks (`useAuth`, `useCart`). Leave components focused solely on rendering visual structure.
- **Hook Dependencies**: Maintain accurate dependency arrays for `useEffect`, `useCallback`, and `useMemo`. Never bypass lint warnings; resolve correct dependency updates instead.

---

## 3. Directory Layout Conventions

Organize React projects using this structure:
```text
src/
├── components/
│   ├── base/        # Headless or low-dependency atomic components (Buttons, Inputs)
│   ├── layout/      # Layout containers (Sidebar, Header, Footer)
│   └── feature/     # Feature-bound domain components (CartList, PaymentForm)
├── hooks/           # Globally shared custom hooks
└── utils/           # Helper scripts and formatters
```

---

## Review Checklist

Verify component structures against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] No class-based components used.
- [ ] Hooks dependencies are exhaustively specified.
- [ ] Complex component state extracted to custom hooks.
- [ ] Directory architecture matches base/layout/feature groupings.
