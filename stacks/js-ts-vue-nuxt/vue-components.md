---
document_id: stacks-vue-components
title: Vue 3 Composition API and Components
ecosystem: js-ts-vue-nuxt
target_versions:
  vue: "^3.3"
dependencies:
  - core-universal-coding-standards
  - stacks-js-ts-conventions
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Vue 3 Composition API and Components

## Inheritance
This document inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md). Refer to the universal standards for component naming, state mutation control, and memory cleanup parameters. This page specifies Vue-specific reactivity hooks.

## Purpose
This document defines patterns for Vue 3 client-side components, reactive structures, Pinia state stores, and lifecycle hooks.

## Scope
Applies to Vue single-file components (`.vue` files) and Pinia configurations.

---

## Directives

### 1. Script Setup & TypeScript
- **Rule**: Every Vue component must use the `<script setup lang="ts">` tag. Do not use the legacy Options API (`data()`, `methods`).
- **Why**: Keeps code concise, compiles with high performance, and integrates naturally with static TypeScript compiler checks.

### 2. Reactivity Guidelines
- **Use ref by default**: Declare reactive variables using `ref()` instead of `reactive()`.
- **Why**: `ref` variables require `.value` in scripts, making their reactive tracking explicit. `reactive` variables lose reactivity when destructured.
- **Exceptions**: Use `reactive()` only when managing cohesive, low-overhead states (e.g. form fields).

### 3. Component Interfacing
- **Strict Interface Typing**: Define props and emits using typed TypeScript configurations:
  ```html
  <script setup lang="ts">
  interface Props {
    title: string;
    isActive?: boolean;
  }

  const props = withDefaults(defineProps<Props>(), {
    isActive: false
  });

  const emit = defineEmits<{
    (e: 'update:status', value: boolean): void;
  }>();
  </script>
  ```

### 4. Pinia State Management
- Use the **Setup Store** syntax instead of Options Store syntax. Setup stores feel natural to composition structures.
- **Code Syntax**:
  ```typescript
  import { ref, computed } from 'vue';
  import { defineStore } from 'pinia';

  export const useUserStore = defineStore('user', () => {
    const profile = ref<UserProfile | null>(null);

    const isLoggedIn = computed(() => profile.value !== null);

    function setProfile(data: UserProfile) {
      profile.value = data;
    }

    return { profile, isLoggedIn, setProfile };
  });
  ```

---

## Component Structure & Discipline

### 5. When to Extract a Component
- Extract a UI block into a reusable component when it **appears in more than one place** in the project.
- Extract even if used once when the block is **complex enough to deserve its own name**.
- Never copy and paste markup \u2014 if you are copying a block, you need a component.

### 6. Component Naming Conventions
- Use **PascalCase** for all component file names and component registrations.
- Names must describe **what the component is**, not what it looks like:
  - `UserProfileCard` (correct) vs. `BigBlueCard` (wrong)
- Choose **one prefix** for shared/base components and use it consistently across the project \u2014 never mix:
  - `BaseButton`, `BaseInput`, `BaseModal` (using `Base` prefix)
- Page-level components get the **`Page` suffix**: `DashboardPage`, `LoginPage`.
- Never give two components names that could be confused with each other.

### 7. Component Responsibility
- Every component does **one thing only**: display, form, layout, or feedback \u2014 never mixed responsibilities in one component.
- Never hardcode content that should come through props.
- Never hardcode styles that should be configurable through props or slots.
- If a component needs **more than one screen of code** to read, it is doing too much. Split it.

### 8. Component Location
- Shared components used across multiple pages → `components/shared/` or `components/common/`
- Page-specific components → `components/[page-name]/`
- Base/UI components → `components/base/` or `components/ui/`
- Never scatter components in the root of `components/` without a folder structure.

### 9. Props Discipline
- Define every prop **explicitly** with its type and whether it is required (using TypeScript interfaces, as shown in Section 3).
- Never use `$parent` or `$root` to pass data down \u2014 use props.
- **Never mutate props** inside a component \u2014 emit an event back to the parent instead.
- Provide **sensible default values** for optional props using `withDefaults`.

### 10. Slots
- Use slots when a component needs to wrap **arbitrary content** it does not control.
- Use **named slots** for components that have multiple distinct content areas.
- Never hardcode content inside a component that the parent should control.

### 11. Reactivity Discipline
- Avoid **unnecessary reactivity** \u2014 not every variable needs to be wrapped in `ref()` or `reactive()`.
- Declare variables as plain constants when their value does not need to trigger re-renders.
- Keep reactive state as shallow as possible. Deeply nested reactive objects are harder to track and debug.

### 12. Composables and Utils
- Logic shared between two or more components belongs in a **composable**, not duplicated inside each component.
- `composables/` for stateful reusable logic (e.g., `useAuth`, `useCart`).
- `utils/` for pure functions with no state (e.g., `formatCurrency`, `truncateText`, `parseDate`).
- Always check `composables/` and `utils/` before writing new logic inside a component.
- Never write the same formatting, calculation, or helper function twice.

---

## Common Mistakes & Anti-Patterns
- **Reactivity Loss**: Destructuring a reactive object (e.g. `const { profile } = store;`) which breaks reactivity connection. Use `storeToRefs(store)` instead.
- **Dangling Timers**: Registering `setInterval` or `addEventListener` on window hooks without removing them on component unmount, causing client memory leaks (refer to [performance/01-profiling-and-benchmarks.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/01-profiling-and-benchmarks.md)).
- **Direct Prop Mutations**: Mutating values received in `props` directly within child components instead of emitting events back to the parent.
- **Unnamed Components**: Anonymous or inline component definitions with no descriptive name in the codebase.
- **Root-Level Scattering**: Placing components directly in `components/` root without any folder grouping as the project grows.
- **Over-Reactive State**: Wrapping static lookup tables or computed reference data in `ref()`/`reactive()` when they never change at runtime.

---

## References
- Caching local states: [performance/01-profiling-and-benchmarks.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/01-profiling-and-benchmarks.md)
- Inertia client integration: [bridges/laravel-inertia-vue.md](file:///Users/kodexkode/Documents/workspace/promptengine/bridges/laravel-inertia-vue.md)

