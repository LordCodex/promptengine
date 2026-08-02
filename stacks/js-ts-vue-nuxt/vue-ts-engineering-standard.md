---
document_id: stacks-vue-ts-engineering-standard
title: Vue 3 and TypeScript Engineering Standard
ecosystem: js-ts-vue-nuxt
target_versions:
  vue: "^3.4"
  typescript: "^5.0"
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-security-engineering-standard
  - core-performance-engineering-standard
  - stacks-js-ts-conventions
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Vue 3 and TypeScript Engineering Standard

## Purpose & Inheritance
This document defines the core standards for Vue 3 frontend development using TypeScript. It inherits from and extends the [Universal Coding Standards](../../core/05-universal-coding-standards.md), the [Architecture Standards](../../core/02-architecture-and-simplicity.md), the [Security Engineering Standard](../../core/08-security-engineering-standard.md), and the [Performance Engineering Standard](../../core/10-performance-engineering-standard.md). It establishes implementation protocols for Single File Components (SFCs), script setup, typing architectures, Pinia state stores, and browser rendering optimization.

---

## 1. Frontend Philosophy

Modern frontend engineering treats the browser application as a **fully realized distributed software client, not a simple presentation UI layer**.
- **Data Architecture First**: Codebases must prioritize modeling data flows, validation boundaries, and runtime states before writing CSS styling and visual components.
- **Enforce Component Boundaries**: UI components should be pure visual interfaces. Extract business calculations, network integrations, and state modifications to dedicated Composables or Pinia stores.
- **Reject Abstraction Bloat**: Do not wrap standard HTML tags in customized components unless they introduce reusable behaviors or semantic layouts (avoid creating a `<CustomButton>` wrapper that just proxies standard click events).

---

## 2. TypeScript Standards & Type Safety

Strict type safety protects applications from runtime `undefined` errors and structural inconsistencies.

### Strict Compiler Configuration
All projects must enable `"strict": true` in `tsconfig.json` configurations. This activates:
- `noImplicitAny`
- `strictNullChecks`
- `strictFunctionTypes`
- `noUnusedLocals`

### Typing Standards
- **Use Interfaces for Object Schemas**: Define data shapes and component structures using `interface`. Interfaces support extension (`extends`) and merger declarations.
- **Use Types for Unions & Utilities**: Use `type` for primitive aliases, union pairings, and mapping utilities.
- **Avoid Enums; Use Union Types**: Enforce string union types (`'active' | 'inactive'`) or readonly constant objects over native TypeScript `enum`. Union types generate zero compilation code output, reducing bundle footprint.
- **Generics**: Enforce generics to build reusable data structures (e.g., standard API response envelopes, collection wrappers).
- **Type Guards**: Implement explicit type guard functions (`isUser(data: unknown): data is User`) to validate dynamic payloads securely.
- **Never Use Any; Enforce Unknown**: If a variable's type is dynamic or unverified (e.g., external API payload), declare it as `unknown`. Force developers and compilers to run type guards before accessing properties:

```typescript
// Good: Type checking required before usage
function processPayload(payload: unknown) {
    if (isInvoice(payload)) {
        console.log(payload.amount_cents); // Compiler is satisfied
    }
}
```

### TypeScript Decision Guidance Matrix
- **Interface vs. Type**: Use `interface` for object models and public API definitions. Use `type` for unions, intersections, and mapping helpers.
- **Enum vs. Union**: Use string union types (`'pending' | 'success'`) to avoid compiler code generation. Use constants objects if mapping key-values is required.
- **Any vs. Unknown**: Use `unknown` for any external payload or untrusted input. Never use `any` as it disables type checking.
- **Generic vs. Duplicate Code**: Create Generic parameters when the underlying logic is identical across 3+ structures (e.g., API adapters).
- **Explicit Types vs. Inference**: Rely on compiler inference for simple initializations (`const count = ref(0)`). Enforce explicit type annotations for complex return types and function arguments.

---

## 3. Vue 3 Application Architecture

We organize components and business layers in feature-focused folders to limit code coupling.

### Project Structure (Feature-Based)
```text
src/
├── assets/           # Global styles and static image assets
├── core/             # Shared utilities, global router, and Axios clients
├── components/       # Shared UI primitives (Buttons, Inputs, Modals)
├── composables/      # Global composables (e.g., useLocalStorage)
├── features/         # Feature domains
│   ├── invoicing/
│   │   ├── components/  # Invoicing UI elements (InvoiceItem.vue)
│   │   ├── composables/ # Invoicing specific logic (useInvoice.ts)
│   │   ├── stores/      # Invoicing state store (invoiceStore.ts)
│   │   └── InvoicesPage.vue
```

### Script Setup Standard
Single File Components (SFCs) must utilize the **`<script setup lang="ts">`** Composition API pattern.

```vue
<!-- Good: Typed SFC with Script Setup Composition API -->
<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  customerId: number;
  amountCents: number;
}

const props = defineProps<Props>();

// Define strictly typed emitted events
const emit = defineEmits<{
  (e: 'pay', id: number): void;
}>();

const formattedAmount = computed(() => `$${(props.amountCents / 100).toFixed(2)}`);

function handlePayment() {
  emit('pay', props.customerId);
}
</script>

<template>
  <div class="invoice-row">
    <span>Amount: {{ formattedAmount }}</span>
    <button type="button" @click="handlePayment">Pay Now</button>
  </div>
</template>

<style scoped>
.invoice-row {
  display: flex;
  justify-content: space-between;
}
</style>
```

---

## 4. Component Design & Composition API

### Component Design Rules
- **Props are Immutable**: Never mutate props inside child components. Emitting events (`emit`) is the only valid way to notify parents of changes.
- **Use Slots for Layout Layouts**: Delegate template layout variations to slots (`<slot name="header" />`) instead of passing complex HTML markup strings through props.
- **Computed Over Watches**: Always use `computed` properties for calculating derived states. Use `watch` or `watchEffect` only to trigger side effects (e.g. running an API fetch on prop changes).

### Component Anti-Patterns
- **Huge components**: Components exceeding 300 lines of template and script logic must be broken down into smaller child components.
- **Too many props**: If a component requires more than 7 props, it is likely handling too many responsibilities. Pass a structured object interface instead.
- **Business Logic in Templates**: Never execute complex logic calculations or database formats inside HTML templates. Keep templates clean and assign formatting to computed properties.

---

## 5. Composables

Composables are stateful functions that capture reusable logic.

### Composable Conventions
- **Naming**: Always prefix composables with `use` (e.g., `usePaymentGateway`).
- **Return Type Standard**: Return properties as reactive references wrapped in a flat object, allowing easy destructuring in `<script setup>` contexts:
  ```typescript
  // return refs explicitly
  return {
      isLoading: readonly(isLoading),
      data: readonly(data),
      execute
  };
  ```
- **State Isolation**: When writing composables, ensure state references are created *inside* the function to prevent accidental state sharing between separate component instances:

```typescript
// Good: State isolated per instance
export function useCounter() {
    const count = ref(0); // Isolated
    const increment = () => count.value++;
    return { count, increment };
}
```

---

## 6. State Management & Pinia

Correct state management prevents duplicate server calls and synchronizes system attributes.

### Store Architecture (Pinia)
- **Local State Default**: Keep state local (`ref`) inside components. Move state to global Pinia stores only when the data is accessed by multiple distinct components across the router namespace.
- **Store Structure**: Group stores by feature domains (e.g., `invoicesStore.ts`). Use actions to modify state; never write direct state alterations outside store actions.

```typescript
// Good: Feature Store (Pinia Setup Store Pattern)
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useInvoiceStore = defineStore('invoice', () => {
    const invoices = ref<Invoice[]>([]);
    const isLoading = ref(false);

    const overdueInvoices = computed(() => 
        invoices.value.filter(inv => inv.status === 'overdue')
    );

    async function fetchInvoices() {
        isLoading.value = true;
        try {
            invoices.value = await api.getInvoices();
        } finally {
            isLoading.value = false;
        }
    }

    return { invoices, isLoading, overdueInvoices, fetchInvoices };
});
```

---

## 7. Data Fetching & Forms

### API Communication SRE
- **Request Cancellation**: Use `AbortController` in composables to cancel pending HTTP queries when a component is unmounted (preventing race conditions and memory leaks on slow networks).
- **Optimistic Updates**: For low-risk mutations (like liking a post), update the local state instantly while processing the backend API call in the background. Roll back the state if the API fails.
- **Unified Loading States**: Bind UI submit buttons directly to reactive loading states (`isLoading`) to prevent duplicate form submissions.

### Forms Architecture
- **Validation Libraries**: Use validation libraries (like VeeValidate or custom TypeScript validation models) to manage schema validation parameters.
- **Decoupled Form State**: Forms must build a mutable local state object separate from model schemas. Bind inputs to a local reactive model and emit a typed payload on submit.

---

## 8. Frontend Performance & Optimization

- **Virtual Lists on Large Datasets**: Never render lists exceeding 200 items using a standard `v-for` loop. Use virtual list libraries (e.g., `vue-virtual-scroller`) to render only the elements visible in the viewport.
- **shallowRef for Large Data Sets**: When storing large static API arrays that do not require deep nested reactivity checks, use `shallowRef` to bypass reactivity tracking overhead:
  ```typescript
  const largeDataset = shallowRef<MetricRecord[]>([]);
  ```
- **Code Splitting (Lazy Loading)**: Configure router definitions to lazy-load page bundles to keep the initial application bundle footprint small:
  ```typescript
  const InvoicesPage = () => import('./features/invoicing/InvoicesPage.vue');
  ```

---

## 9. Security Engineering (Frontend Specific)

The client application must protect user sessions and prevent script injection.

### Security Controls
- **Never Use `v-html` Directly**: Rendering raw HTML using `v-html` invites Cross-Site Scripting (XSS) attacks. If HTML rendering is required, sanitize the input using a proven sanitization library (such as DOMPurify).
- **Secure Token Storage**: Do not store authentication JWTs or credentials inside `localStorage` if they are susceptible to XSS. Store them in cookie configurations configured with `HttpOnly` and `Secure` parameters (managed by the backend gateway).
- **CSP Nonces**: Pass server-side generated nonces to dynamically loaded JavaScript packages.

---

## 10. Accessibility (A11y)

Building accessible interfaces is a core requirement, not an optional enhancement.

### A11y Implementation Guidelines
- **Semantic HTML**: Use native elements (`<button>`, `<nav>`, `<main>`, `<header>`) instead of styling `<div>` blocks with click listeners.
- **Keyboard Navigation**: Ensure all interactive elements can be focused and triggered using only the keyboard (`Tab`, `Space`, `Enter`).
- **Focus Management**: When opening modals or dialogues, trap focus inside the modal frame and return it to the trigger button when the modal closes.
- **ARIA Attributes**: Use explicit ARIA attributes (`aria-expanded`, `aria-label`, `aria-hidden`) to provide context to screen readers on custom components.

---

## 11. Testing Strategy

Verify behavior, not implementation configurations.

### Vitest & Vue Test Utils Examples
- **Component Tests**: Mount components, simulate user actions, and assert DOM output changes:

```typescript
import { mount } from '@vue/test-utils';
import { test, expect } from 'vitest';
import InvoiceRow from './InvoiceRow.vue';

test('InvoiceRow emits pay event when button clicked', async () => {
    const wrapper = mount(InvoiceRow, {
        props: { customerId: 102, amountCents: 5000 }
    });

    await wrapper.find('button').trigger('click');

    expect(wrapper.emitted().pay[0]).toEqual([102]);
});
```

---

## 12. Decision Matrices

Use these matrices to identify the correct frontend engineering decision based on project context.

### Matrix 1: Component vs. Composable
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| UI elements layout, rendering lists, styling widgets | **Component** | Handles HTML templates and visual layouts. |
| Reusable behavior, API fetch logic, key event handlers | **Composable** | Decouples business logic from rendering structures. |

### Matrix 2: Local State vs. Global State
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Component-specific inputs, toggles, loading states | **Local State (`ref`)** | Simple, fast, and self-contained within the component lifecycle. |
| User profile information, shared catalog caches | **Global State (Pinia)** | Synchronizes state across distinct routing targets. |

### Matrix 3: Pinia vs. Provide/Inject
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Features state accessed by nested child components in a hierarchy | **Provide/Inject** | Limits state exposure to the sub-tree context. |
| Cross-domain data shared across unrelated page layouts | **Pinia** | Centralizes store management and provides debugging footprints. |

### Matrix 4: Computed vs. Watch/WatchEffect
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Calculating formatted amounts, filtering lists, combining states | **Computed** | Automatically caches results and re-evaluates only when dependencies change. |
| Fetching API data on user ID change, writing to local storage | **Watch** | Executes asynchronous or procedural side-effects. |

### Matrix 5: Reusable Component vs. One-Off Component
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Elements appearing multiple times (Buttons, Inputs) | **Reusable** | Centralizes global design adjustments and keeps styling consistent. |
| Complex page-specific headers or custom widgets | **One-Off** | Avoids over-parameterization; keeps components easy to modify. |

### Matrix 6: CSS Utility (e.g., Tailwind) vs. Scope Component Styling
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Rapid layout prototyping, dashboard grids setups | **CSS Utility** | Speeds up development and eliminates custom CSS clutter. |
| UI Primitive libraries (Buttons, Modals) | **Scoped CSS** | Encapsulates component styles, preventing style leaks. |

---

## 13. AI Frontend Rules

AI agents modifying or writing frontend code in this repository must follow these rules:

1. **Verify Existing Primitives**: Before coding a new form field or input component, check if a primitive component already exists in the design system.
2. **Strict Props and Emits Typing**: Enforce explicit interfaces for all props and emits declarations. Do not use generic, untyped events.
3. **No raw v-html**: Do not suggest using `v-html` unless the input is processed through a sanitization helper.
4. **Use Computed for Formatting**: Suggest computed variables for formatting dates, currency, or filter structures.
5. **No Framework Internals Testing**: Ensure component tests verify rendered output behaviors (e.g., text displays, events emit) rather than checking internal script variables.

---

## 14. Frontend Review Checklist

Use this checklist during code review to evaluate Vue and TypeScript changes.

### Type Safety & TS
- [ ] Are all props and emits declared using explicit TypeScript interfaces?
- [ ] Is the use of `any` blocked (enforcing `unknown` or explicit types)?
- [ ] Are union types used instead of TypeScript `enums`?

### Components & Composition API
- [ ] Is script setup Composition API used exclusively?
- [ ] Are computed properties used for derived calculations (no watches for calculations)?
- [ ] Are component files kept under 300 lines (with business logic extracted)?

### State & Fetching
- [ ] Is global store data scoped correctly (no unnecessary global states in Pinia)?
- [ ] Are API requests aborted (`AbortController`) on component unmount?

### Security & A11y
- [ ] Is the use of `v-html` restricted and sanitized?
- [ ] Do all interactive elements support keyboard navigation and focus management?
- [ ] Are screen reader attributes (ARIA labels) declared on custom elements?

### Testing & Code Quality
- [ ] Do component tests verify output behavior (not script internals)?
- [ ] Are all resources lazy-loaded at the routing layer?

---

## References
- Universal Naming Rules: [core/05-universal-coding-standards.md](../../core/05-universal-coding-standards.md)
- Security Engineering: [core/08-security-engineering-standard.md](../../core/08-security-engineering-standard.md)
- Performance Engineering: [core/10-performance-engineering-standard.md](../../core/10-performance-engineering-standard.md)
- Nuxt SSR Boundaries: [nuxt-ssr.md](nuxt-ssr.md)
