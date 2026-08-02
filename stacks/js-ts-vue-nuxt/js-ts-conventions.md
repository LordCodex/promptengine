---
document_id: stacks-js-ts-conventions
title: Modern JavaScript and TypeScript Conventions
ecosystem: js-ts-vue-nuxt
target_versions:
  typescript: ">=5.0"
  node: ">=18.0"
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Modern JavaScript and TypeScript Conventions

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for module layouts, variable scopes, and async safety. This page specifies JavaScript and TypeScript syntax types.

## Purpose
This document defines standards for writing ES modules, establishing TypeScript typing boundaries, and configuring async handler operations.

## Scope
Applies to client-side Vue files, Nuxt backend API routes, and general TS configurations.

---

## Directives

### 1. Modern ES Modules (ESM)
- **Standard**: Always use standard `import`/`export` syntax. Avoid legacy CommonJS `require()` bindings.
- **Rule**: Keep imports clean. Group third-party package imports at the top, followed by local components or helper modules.

### 2. Strict TypeScript Typing
- **Standard**: Configure `tsconfig.json` with strict mode enabled:
  ```json
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true
  }
  ```
- **Rule 1**: Never use `any` as a type definition. If a parameter type is unknown, use `unknown` and perform type checking before accessing fields.
- **Rule 2**: Use `interface` for structural object shapes and public APIs. Use `type` for unions, intersections, or utility aliases.
- **Code Syntax**:
  ```typescript
  interface UserProfile {
    readonly id: string;
    email: string;
    displayName: string;
    role: 'admin' | 'marketer' | 'author';
  }
  ```

### 3. Asynchronous Error Boundaries
- **Rule**: Every `await` statement must operate inside a `try/catch` block or be chained with a global fallback block to prevent silent request failures.
- **Code Syntax**:
  ```typescript
  async function fetchUserData(userId: string): Promise<UserProfile> {
    try {
      const response = await api.get<UserProfile>(`/users/${userId}`);
      return response.data;
    } catch (error) {
      logger.error('Failed to fetch user', { userId, error });
      throw new FetchEntityException('User lookup failed');
    }
  }
  ```

---

## Common Mistakes & Anti-Patterns
- **The any Escape Hatch**: Using `any` typing blocks to silence TS compiler errors instead of writing proper interfaces or dynamic validation assertions.
- **Floating Promises**: Dispatched asynchronous calls that run without `await` or `.catch()` callbacks, leading to uncaught event loops failures.
- **Mutating Readonly Parameters**: Direct mutation of values passed as function parameters instead of returning a new, clean object copy.

---

## References
- Environment validation checks: [environment/03-ci-cd-pipelines.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/03-ci-cd-pipelines.md)
- Vue Component integrations: [vue-components.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/js-ts-vue-nuxt/vue-components.md)
