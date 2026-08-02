---
document_id: stacks-nuxt-ssr
title: Nuxt 3 SSR and Server Routes
ecosystem: js-ts-vue-nuxt
target_versions:
  nuxt: "^3.0"
dependencies:
  - core-universal-coding-standards
  - stacks-js-ts-conventions
  - stacks-vue-components
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Nuxt 3 SSR and Server Routes

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for class scoping, module separation, and async error boundaries. This page specifies only Nuxt-specific hydration hooks.

## Purpose
This document defines coding conventions for Nuxt 3, handling Server-Side Rendering (SSR) states, resolving hydration mismatches, and structuring server api routes.

## Scope
Applies to Nuxt pages, layouts, server middleware, and routes under `/server/api/`.

---

## Directives

### 1. SSR-Friendly Hydration
- **Rule**: Avoid running browser-specific API code (e.g. `window`, `document`, `localStorage` references) directly inside the initialization block of script setups.
- **Why**: These elements do not exist on the server node during initial SSR render, causing compilation exceptions or hydration mismatches.
- **Code Syntax**:
  ```typescript
  import { onMounted, ref } from 'vue';

  const userTheme = ref('light');

  onMounted(() => {
    // browser APIs are safe to call only inside onMounted
    userTheme.value = localStorage.getItem('theme') || 'light';
  });
  ```

### 2. Server-Only Routes (/server/api)
- Use server-only endpoint files under `/server/api/` (using Nitro) to securely access backend databases or private credentials without exposing code to client bundles.
- **Naming Rule**: Suffix API routes with `.get.ts`, `.post.ts`, `.put.ts`, `.delete.ts` to enforce strict HTTP method mappings.
- **Code Syntax**:
  ```typescript
  // server/api/users/[id].get.ts
  export default defineEventHandler(async (event) => {
    const userId = getRouterParam(event, 'id');
    const config = useRuntimeConfig(event);

    // nitro fetches secure data
    const user = await $fetch(`https://api.external.com/v1/users/${userId}`, {
      headers: { Authorization: `Bearer ${config.apiSecret}` }
    });

    return { data: user };
  });
  ```

### 3. State Hydration with useState
- **Rule**: Do not use global ref variables declared outside components for sharing state between requests in SSR. This leaks state across separate concurrent users.
- **Good**: Use Nuxt's `useState()` helper to manage shared server-safe component states:
  ```typescript
  const activeTab = useState('active_tab', () => 'overview');
  ```

---

## Common Mistakes & Anti-Patterns
- **Hydration Mismatch**: Rendering dynamic states that differ between the server and the client (e.g., using `new Date()` directly in templates instead of formatting it inside `onMounted`).
- **Global Memory Leaks**: Defining global variables outside component cycles in server middleware or composables.
- **Direct Database in Client**: Exposing backend config tokens or executing raw DB integrations inside standard client components.

---

## References
- Client routing setups: [vue-components.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/js-ts-vue-nuxt/vue-components.md)
- Performance and thread boundaries: [performance/04-concurrency-and-async.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/04-concurrency-and-async.md)
