---
document_id: decision-guides-state-management
title: State Management Decision Guide
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-03
---

# State Management Decision Guide

This guide establishes the rules for selecting and structuring application state solutions in React/Next.js, Vue/Nuxt, and Flutter ecosystems.

---

## 1. Ecosystem Selection Table

| Framework | Local State (Component) | Global Client State | Server State (Caching) |
| :--- | :--- | :--- | :--- |
| **React / Next.js** | `useState` | **Zustand** (or Context for simple props) | **TanStack Query** (React Query) |
| **Vue / Nuxt** | `ref`, `reactive` | **Pinia** | **useFetch** / **TanStack Query** |
| **Flutter** | `StatefulWidget` / `useState` | **Riverpod** or **BLoC** | Cached repository classes |

---

## 2. Selection Framework

### Local State Boundary
- **Standard**: Keep state as close to the component that needs it as possible.
- **Rule**: Do not add state to global stores (Pinia/Zustand) if the state is only used by a single component and its direct children.

### Server State vs. Client State
- **Rule**: Do not duplicate API data in global stores. Use server-caching mechanisms (e.g. TanStack Query or Pinia's data cache) to manage fetched payloads.
- Use global client stores only for UI-specific persistent flags (e.g. dark mode toggle, side navigation expanded state) and authenticated session attributes.

---

## 3. Reference Implementations
- React State Guidelines: [React State Management](../stacks/react-next/react-state-management.md)
- Vue State Guidelines: [Vue Component Rules](../stacks/js-ts-vue-nuxt/vue-components.md)
- Flutter State Guidelines: [Flutter State Management](../stacks/dart-flutter/flutter-state.md)
