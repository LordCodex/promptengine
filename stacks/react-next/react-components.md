---
document_id: stacks-react-components
title: React Server Components and Render Architecture
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# React Server Components and Render Architecture

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines component rendering architectures specific to React Server Components (RSC) and Client Components.

---

## 1. Server Components vs. Client Components

To build performant React architectures:
- **Server Components (RSC) by Default**: All components are Server Components by default. Use Server Components to fetch data from database tables or third-party APIs directly, reducing client-side bundle sizes.
- **Client Components on Demand**: Use `'use client'` at the absolute top of the file only if the component requires:
  - Interactive hooks: `useState`, `useReducer`, `useEffect`.
  - Event listeners: `onClick`, `onChange`.
  - Browser-specific APIs (e.g. `window`, `localStorage`).
- **Minimize Client Boundaries**: Keep Client Components leaf-nodes. Do not make root layout files Client Components simply because a small sub-element needs state.

---

## 2. Suspense and Streaming Boundaries

- **Granular Loading States**: Wrap slow-loading Server Components in `<Suspense fallback={<Skeleton />}>` boundaries to allow HTML streaming.
- **Streaming Response**: Let high-priority content (e.g. main page layouts) render immediately while slower data feeds resolve asynchronously behind Suspense gates.

---

## 3. Composing Layouts

- Place layout rules inside `layout.tsx` wrappers.
- Pass children dynamic variables as React nodes to layout wrappers to allow nesting without complete client redraws.

---

## Review Checklist

Verify rendering structures against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] Default components do not use `'use client'` unless interactive APIs are strictly required.
- [ ] Slower components are isolated behind `<Suspense>` loaders.
- [ ] Server Components are used for data fetching boundaries.
