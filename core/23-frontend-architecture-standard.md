---
document_id: core-frontend-architecture
title: Frontend Architecture Standard
ecosystem: cross-cutting
dependencies:
  - core-architecture-and-simplicity
  - core-universal-coding-standards
  - core-security-engineering-standard
  - core-performance-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Frontend Architecture Standard

## Inheritance
This document inherits from and extends the [Architecture and Project Structure Standard](02-architecture-and-simplicity.md) and the [Universal Coding Standards](05-universal-coding-standards.md). It defines framework-agnostic architectural boundaries and rules for client-side web and mobile applications, ensuring scalability, reliability, and security.

---

## 1. General Frontend Principles

Every frontend application must maintain a strict separation of concerns across four key layers:

| Layer | Responsibility | Key Rule |
| :--- | :--- | :--- |
| **Presentation** | Rendering UI, visual state representation, handling direct user interactions. | Purely declarative; zero business logic or API communication. |
| **Application Logic** | Managing user workflows, coordinating state transitions, and applying client-side business rules. | Framework-specific logic is encapsulated in hooks, composables, or controllers. |
| **Data Layer** | API communication, local database management, caching, data normalization/transformation. | Unifies client-server synchronization; acts as the single data source. |
| **Shared Utilities** | Stateless helper functions, validators, formatting routines, mathematical operations. | Pure functions; zero dependency on framework lifecycles. |

---

## 2. Component Separation & Decoupling

A component must not become a complete application. To prevent maintenance rot:

- **Do not put everything in components**: Keep components small, focused, and single-purpose.
- **No API calls in UI templates**: A component should not directly trigger HTTP requests. Delegate fetching to the data layer or application logic.
- **Separate formatting from presentation**: Do not place date formatting, currency calculations, or string mutations directly inside markup templates. Use shared utilities or filters.
- **No duplicate state handling**: If two components require the same dynamic state, lift the state to their nearest common parent or a dedicated store. Never synchronize state manually across sibling components.

---

## 3. Frontend File Organization

File layouts must match project complexity. Do not create folders without purpose or over-structure small codebases.

### Standard Production Structure

```text
src/
├── api/            — Central API client, endpoint definitions, interceptors
├── assets/         — Static assets (images, global icons, fonts)
├── components/     — Shared visual components (base UI and composed blocks)
├── composables/    — Vue composition functions (or hooks/ for React)
├── constants/      — System-wide read-only config/values
├── features/       — Vertical slices of domain-specific screens and logic
├── layouts/        — Route structural templates (dashboard, auth, print)
├── pages/          — Entry points linked directly to routes
├── services/       — Business rules and external client integrations
├── stores/         — Shared global/server state instances
├── types/          — TypeScript interfaces and type definitions
└── utils/          — Pure helper functions and validators
```

---

## 4. State Management

Choose state scopes based on real architectural needs. Do not introduce global state solutions by default.

### State Scopes

- **Local State**: For component-specific interaction state (e.g. `isOpen`, input validation states, hover status).
- **Shared State**: For cross-component configurations, global settings, and metadata (e.g. user authentication, user theme choice, active locale).
- **Server State**: For caching remote API data and resources.

### State Rules
- **Single Source of Truth**: Never duplicate server-retrieved data into multiple global stores. Keep a single cache.
- **No Global Store Overuse**: Do not use global state as a quick shortcut to bypass clean component prop/event architecture.
- **Predictable Updates**: State updates must be explicit and traceable. Mutate state via named actions or mutations; never mutate global state directly inside a view.

---

## 5. Data Fetching and API Communication

### Data Fetching Requirements
All network data operations must explicitly support and handle:
- **Loading states**: Ensure visual skeletons or spinners prevent interaction while fetching.
- **Error handling**: Catch all network failures and translate them into friendly, readable error UI.
- **Retry behavior**: Build in back-off retry logic for critical background fetches.
- **Caching**: Prevent redundant requests by matching request keys against cached data.
- **Pagination**: Implement cursor or keyset page boundaries to protect client memory from huge payloads.
- **Request cancellation**: Abort pending requests when a user navigates away from the active screen.

### Centralized API Client
- Every application must route communication through a centralized API client instance (e.g., Axios client, Fetch wrapper).
- Configure authentication, authorization headers, and token refreshes centrally.
- Normalise error payloads into a consistent client-side error format using response interceptors.
- **Never allow individual components to instantiate their own HTTP configurations.**

---

## 6. Performance Rules

### Rendering Optimization
- **Prevent Unnecessary Re-renders**: Use memoization, keys, and shallow comparisons to prevent duplicate paint calculations.
- **Limit Reactive State**: Keep reactivity chains short. Do not wrap large, static data arrays in deeply tracked reactive structures.
- **Manage Large Component Trees**: Split deep trees into clean, lazy-rendered blocks.

### Asset Management
- **Optimize Assets**: Compress images, specify fallback system fonts, and use SVG for all vector items.
- **Lazy Loading**: Use code-splitting and dynamic imports for non-critical pages and heavy modal components.

---

## 7. Routing and Navigation

- **Authorization Guards**: Enforce authentication and role permissions inside router navigation guards before mounting pages.
- **UI is Not Security**: Client-side route blocking is for user experience only. **Security decisions must always be enforced server-side.**
- **Error Boundaries**: Configure route boundaries to handle 404 (Not Found), 403 (Unauthorized), and 500 (Internal Server Error) states gracefully.

---

## 8. Form Standards

Every form implementation must cleanly manage its state lifecycle:
- **Validation**: Enforce instant visual validation indicators during user input (e.g. email format checks on blur).
- **Processing states**: Disable the submit button and all inputs while submission is in progress.
- **Failure feedback**: Display server validation errors next to their specific input fields.
- **Never trust client-side validation for security**: Frontend checks are for convenience; backend validation is mandatory.

---

## 9. Framework-Specific Architectures

### 9a. Vue 3 / Nuxt
- **Prefer**:
  - Vue 3 Composition API with `<script setup lang="ts">`.
  - Composables (`composables/`) to encapsulate reusable stateful logic.
  - Pinia setup stores for global shared state.
  - VueUse for standardized helper utilities where appropriate.
- **Avoid**:
  - Mixing Options API with Composition API in new components.
  - Writing complex business logic inside visual component templates.
  - Massive, multi-purpose single-file components.

### 9b. React / Next.js
- **Prefer**:
  - Small, functional components with explicit, typed props.
  - Custom React Hooks to decouple presentation from side-effects.
  - Strict Server Component / Client Component boundaries (`"use client"`).
- **Avoid**:
  - Excessive prop drilling (use context or composition to pass data down).
  - Unnecessary global context state for items that should be local.

### 9c. Flutter / Dart
- **Prefer**:
  - Declaring widgets with `const` constructors to maximize compilation and render optimization.
  - Splitting widgets into distinct feature-specific directories.
  - Clean state management architectures (Riverpod, Bloc/Cubit, or Provider depending on complexity).
- **Avoid**:
  - Putting API calls or database operations inside widget build methods.
  - Repeating custom widget trees inside screens without extracting them.

### 9d. Pure JavaScript Applications
- **Prefer**:
  - ES6 modules (`import`/`export`) to maintain boundaries.
  - Native browser APIs over external utility packages unless strictly required.
  - Small, single-responsibility files.
- **Avoid**:
  - Pulling in complex framework wrappers for simple, lightweight applications.

---

## 10. TypeScript Guidelines

- **Explicit Boundaries**: Declare explicit types or interfaces for all API request/response payloads and component interfaces.
- **Avoid `any`**: The use of `any` is prohibited. Use `unknown` or specify explicit types.
- **Simple Typings**: Avoid over-engineered generic types that make the codebase difficult for other developers to read.

---

## 11. Error Handling & Security Gates

### Error Safety
- Catch and handle network timeouts, API auth failures, database errors, and empty responses.
- **Never assume an API call succeeds.** Provide recovery fallback UI.

### Security
- **No Secrets in Frontend**: Secrets, keys, and private passwords must never be compiled into the frontend build. Use server-side proxy routes for protected API calls.
- **Prevent XSS**: Sanitize all dynamic string variables before inserting them into HTML templates. Use native text binding.
- **Cookies**: Use `HttpOnly` and `Secure` cookie-based authentication schemas where applicable.

---

## 12. Accessibility Engineering

Accessibility is part of frontend architecture. All interactive elements must:
- Use semantic HTML elements first (e.g. `<button>` for clickable actions).
- Support tab-key keyboard focus and clear focus states.
- Support screen reader text labels and standard ARIA attributes where custom UI is necessary.

---

## Review Checklist

Before declaring any frontend work complete, verify against this checklist:
- [ ] **Logic Separation**: Is all application and business logic decoupled from UI presentation files?
- [ ] **Component Focus**: Do all components have a single, clean responsibility?
- [ ] **State Clarity**: Is the state scoped correctly (local vs. shared vs. server)?
- [ ] **API Organization**: Are all API endpoints routed through a centralized client wrapper?
- [ ] **Error Fallbacks**: Do friendly, human-readable states handle network and validation errors?
- [ ] **Loading States**: Are visual placeholders or skeletons displayed during all async processes?
- [ ] **Performance**: Are code-splitting, lazy loading, and rendering optimizations implemented?
- [ ] **Security Checked**: Are there zero hardcoded API keys or secrets in the client-side files?
- [ ] **Patterns Followed**: Does the layout align with the project's standard structure?
- [ ] **Simplicity**: Is this the simplest architectural approach that solves the problem?

---

## References
- Architecture & Simplicity: [core/02-architecture-and-simplicity.md](02-architecture-and-simplicity.md)
- Universal Naming: [core/05-universal-coding-standards.md](05-universal-coding-standards.md)
- Security Engineering: [core/08-security-engineering-standard.md](08-security-engineering-standard.md)
- Performance Engineering: [core/10-performance-engineering-standard.md](10-performance-engineering-standard.md)
- Design Governance: [design/09-design-system-governance.md](../design/09-design-system-governance.md)
