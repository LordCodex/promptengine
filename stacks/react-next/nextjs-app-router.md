---
document_id: stacks-nextjs-app-router
title: Next.js App Router Conventions
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Next.js App Router Conventions

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines standards for route hierarchies, page/layout files, and navigation components.

---

## 1. Directory Structure and File Conventions

Enforce directory-based routing in the `app/` folder:
- **`layout.tsx`**: Defines layout structures. Must receive children as react node props. Layouts do not re-render on navigation hooks.
- **`page.tsx`**: The unique visual route entry component.
- **`loading.tsx`**: Standard fallback skeleton screen displayed automatically during page load phases.
- **`error.tsx`**: Client-side error boundary catching runtime exceptions. Must declare `'use client'`.
- **`not-found.tsx`**: Global or route-level resource missing page.

---

## 2. Dynamic and Catch-All Routes

- **Dynamic Segments**: Use `[id]` folders for parametric paths. Read dynamic arguments through `params` parameters in pages or layouts.
- **Optional Catch-All**: Use `[[...slug]]` only when routing nested directories dynamically under a single page hierarchy.

---

## 3. Parallel and Intercepting Routes

- **Parallel Routes (`@slot`)**: Use to render multiple independent sub-pages simultaneously in the same layout (e.g. dashboards with side-panels).
- **Intercepting Routes (`(...)`)**: Use to load a page route dynamically inside the current context overlay (e.g. opening a photo detail modal while keeping the underlying feed visible).

---

## Review Checklist

Verify routing configurations against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] Directory routes contain unique page/layout files.
- [ ] `error.tsx` utilizes Client Components.
- [ ] Dynamic parameter reads are typed cleanly.
