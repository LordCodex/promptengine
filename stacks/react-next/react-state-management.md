---
document_id: stacks-react-state-management
title: React State Management and Fetching Standards
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# React State Management and Fetching Standards

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines rules for state storage, client fetching, and cache boundaries.

---

## 1. Local State vs. Global State

- **Localize State by Default**: Keep component states local (`useState`) whenever possible. Do not lift state globally unless multiple independent components must share reading/writing access.
- **React Context**: Use React Context only for static settings that rarely change (e.g. active color themes, localized translations). Avoid Context for high-frequency updates to prevent parent-tree re-render performance issues.

---

## 2. Server Cache & Fetching: TanStack Query (React Query)

Prefer **TanStack Query** for client-side API synchronization, caching, and background loading.
- **Client Cache**: Use query keys systematically (e.g. `['users', userId]`).
- **Separation of Concerns**: Wrap query/mutation hooks into custom hooks (`useUser`, `useUpdateUser`) to keep components clean.
- **Stale Times**: Configure explicit `staleTime` and `gcTime` limits to prevent redundant HTTP requests.

---

## 3. Global Store Selection Criteria (Zustand vs. Redux)

Choose global state managers according to architectural complexity:
- **Zustand (Recommended Default)**: Use for lightweight, composable global state stores. Best for UI states, cart collections, and simple client configurations.
- **Redux Toolkit**: Use only in legacy enterprise environments with complex transaction flows, deep dev-tools trace rules, or strict multi-reducer middleware structures.

---

## Review Checklist

Verify state configurations against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] No high-frequency state properties are bound to raw React Context elements.
- [ ] Custom hooks wrap all client fetch queries using TanStack Query.
- [ ] Zustand is selected as the default lightweight global state manager.
