---
document_id: design-system-governance
title: Design System Governance Standard
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
  - design-component-libraries
  - core-architecture-and-simplicity
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Design System Governance Standard

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md), the [Design Systems Standard](01-design-systems.md), and the [Component Library Awareness Standard](02-component-libraries.md). Those documents define *what* a design system is and *which libraries to evaluate*. This document defines *how a design system is governed* — how components are organized, promoted, documented, and kept from proliferating into chaos.

---

## 1. UI as a System, Not Isolated Screens

The most important governance principle is also the most frequently violated:

**Treat the UI as a system. Never build isolated screens.**

Every component built for one screen is a candidate for reuse elsewhere. Every token defined for one section must be consistent with every other section. Every interaction pattern introduced must match or deliberately extend the patterns that already exist.

An AI agent or developer who creates a screen in isolation — without reviewing what components already exist, without using the token system, without considering how this screen fits into the broader product — is building technical debt, not a feature.

---

## 2. Component Hierarchy and Ownership

Components are not all equal. They have different scopes of reuse, different levels of abstraction, and different ownership models. The hierarchy defines who owns a component and where it lives.

### Standard Hierarchy

```text
components/
│
├── base/           — Foundational UI primitives
│   ├── Button/
│   ├── Input/
│   ├── Modal/
│   ├── Badge/
│   └── Typography/
│
├── common/         — Composed reusable blocks
│   ├── SearchBar/
│   ├── Pagination/
│   ├── DataTable/
│   ├── EmptyState/
│   └── LoadingSpinner/
│
├── features/       — Feature-specific components
│   ├── billing/
│   ├── auth/
│   └── dashboard/
│
└── layouts/        — Page structure wrappers
    ├── DashboardLayout/
    ├── AuthLayout/
    └── PublicLayout/
```

Adapt this structure to the project's framework and conventions. Do not force it on a project that does not need it. A small project with 10 components does not need 4 directory levels.

### Ownership Rules by Level

| Level | Scope | Ownership | When to Create |
| :--- | :--- | :--- | :--- |
| **base/** | Project-wide | Design system team / senior engineers | When a primitive is used in 3+ places; when it wraps an accessibility concern |
| **common/** | Multiple features | Any engineer, reviewed by senior | When a composed pattern appears in 2+ features |
| **features/** | One feature | Feature engineer | When a component is specific to a single domain and has no reuse value |
| **layouts/** | Entire page | Senior / architect | When a page structure pattern is shared across multiple routes |

### Ownership Rules
- A component in `base/` may be used by any component at any level.
- A component in `common/` may be used by any feature or layout — but never in `base/`.
- A component in `features/feature-x/` must not be used in `features/feature-y/`. If reuse is needed, promote it to `common/`.
- A component in `layouts/` wraps page structure only — it must not contain feature-specific business logic.

---

## 3. When to Create a Component

The most common failures are: creating too many components (fragmentation) and not creating enough (copy-paste duplication).

### Create a Component When
- The same UI pattern appears in **2 or more places** in the codebase.
- The component encapsulates **meaningful behaviour** (validation, state management, keyboard handling).
- The component represents a **clear, nameable concept** (UserAvatar, PriceDisplay, StatusBadge).
- Extracting it will **improve consistency** — future instances use the component instead of reimplementing the pattern.

### Do Not Create a Component For
- A **one-time wrapper** around a single line of markup with no reuse potential.
- A **trivial text block** with no behaviour, state, or styling variation.
- An **isolated element** that will never logically appear outside its current context.
- A component that **only wraps another component** without adding any props, behavior, or semantic meaning.

### The Rule of Three for Component Promotion
1. **First occurrence** — implement inline. Do not extract yet.
2. **Second occurrence** — note the duplication. Extract if the pattern is identical or nearly identical.
3. **Third occurrence** — extraction is mandatory. Three inline implementations of the same pattern is a design system violation.

---

## 4. Component Responsibility

Every component must have one clear, single responsibility.

A component should not:
- Manage unrelated business logic.
- Fetch unrelated data.
- Control unrelated components.
- Contain multiple unrelated workflows.

Always separate:
- **Presentation**: How the component renders visually (styles, layouts, accessibility wrappers).
- **State Management**: Tracking component-specific interaction state (e.g., `isOpen`, `activeTab`).
- **Business Logic**: Domain rules and calculations (e.g., pricing calculations, permission checks).
- **Data Fetching**: Retrieving data from APIs or local storage.

---

## 5. Component API Design

A component's API is a public contract. Once other parts of the codebase depend on it, changing that API is a breaking change. Design it deliberately from the start.

### Prop Naming Rules
- Use descriptive, unambiguous prop names that read naturally in context.
- Props should describe **what the component is or does**, not implementation details:
  - `variant="primary"` not `isPrimary`
  - `size="lg"` not `isLarge`
  - `isLoading` not `showSpinner`
  - `isDisabled` not `disabled` (unless the framework convention requires it)
- Boolean props should represent states, not configuration: `isLoading`, `isDisabled`, `isOpen` — not `useModal`, `enableShadow`, `withBorder`.

### Avoid Boolean Prop Explosion

The following anti-pattern creates an unmaintainable API:

```html
<!-- Wrong: boolean prop explosion -->
<Button primary big rounded blue shadow leftIcon disabled loading />
```

Every boolean combination becomes a hidden variant that is never tested. Use explicit variant and size props instead:

```html
<!-- Correct: explicit, readable, testable -->
<Button variant="primary" size="lg" :isLoading="submitting" :isDisabled="!isValid">
  Save Changes
</Button>
```

### Prop Design Rules
- Define every prop with an explicit type and — where applicable — a set of allowed values (TypeScript enum or union type).
- Provide sensible defaults for all optional props.
- The maximum number of props a component should accept before architectural review is **12**. More than 12 props is a signal the component is doing too much.
- Never use prop names that conflict with native HTML attributes unless the component explicitly wraps that element.

---

## 6. Variant Management

When a component needs to express visual or behavioral variations, always prefer **variants over parallel components**.

### Variants Over Parallel Components

| Anti-pattern (parallel components) | Correct approach (variants) |
| :--- | :--- |
| `PrimaryButton`, `SecondaryButton`, `DangerButton` | `<Button variant="primary\|secondary\|danger">` |
| `LargeInput`, `SmallInput`, `InlineInput` | `<Input size="sm\|md\|lg" display="block\|inline">` |
| `DashboardCard`, `ProfileCard`, `SummaryCard` | `<Card variant="dashboard\|profile\|summary">` |
| `AdminTable`, `CustomerTable` | `<DataTable :columns="adminCols\|customerCols">` |

Parallel components diverge immediately. The same bug must be fixed in multiple places. Variants stay synchronized because they share a single implementation.

### When Parallel Components Are Justified
A separate component (not a variant) is justified only when:
- The internal structure is fundamentally different (not just styled differently).
- The props interface is completely different.
- The component would require more than 30% of its code to be conditional on a variant flag.

Document the justification in the component file.

### Variant Implementation Approaches

Prefer these techniques, in order of preference:

1. **Variant prop + CSS class map** — the component applies a class based on the variant value.
2. **Slots** — the parent controls content layout; the component controls wrapping structure.
3. **Composition** — combine base components to produce higher-level components.
4. **`as` prop / polymorphic component** — the component renders as different HTML elements based on context.

---

## 7. Slots and Composition

Slots allow a component to wrap arbitrary content without needing to know what that content is. Use them when:
- The parent must control content layout or content type.
- The component is a structural wrapper (card, modal, dropdown, layout).
- The content varies significantly between uses.

```html
<!-- Correct: slot-based card that controls structure, not content -->
<BaseCard>
  <template #header>
    <h3>Order Summary</h3>
  </template>
  <template #body>
    <OrderLineItems :items="order.items" />
  </template>
  <template #footer>
    <Button variant="primary">Pay Now</Button>
  </template>
</BaseCard>
```

**Named slots** are required when a component has more than one distinct content area. Anonymous slots are acceptable only when there is exactly one content area.

---

## 8. Component Documentation Standard

Documentation burden must match the component's scope. Do not document every private component. Do document every component in `base/` and `common/`.

### What to Document

Every shared component (`base/` and `common/`) must have a documentation block at the top of the file:

```typescript
/**
 * BaseButton
 *
 * Purpose:
 *   The primary interactive action element. Used for all clickable actions
 *   that do not navigate to a new URL. For navigation, use <BaseLink>.
 *
 * Variants:
 *   - primary   — Main call-to-action. Use once per screen section.
 *   - secondary — Supporting action. Lower visual weight than primary.
 *   - danger    — Destructive actions (delete, revoke, cancel subscription).
 *   - ghost     — Tertiary actions with minimal visual presence.
 *
 * Sizes:
 *   - sm  — Use inline with text or in tight spaces.
 *   - md  — Default size for most contexts.
 *   - lg  — Use in hero sections or primary page actions.
 *
 * Props:
 *   @prop {string}  variant    — Visual style. One of: primary | secondary | danger | ghost
 *   @prop {string}  size       — Size scale. One of: sm | md | lg. Default: md
 *   @prop {boolean} isLoading  — Replaces label with a spinner; disables interaction.
 *   @prop {boolean} isDisabled — Prevents interaction; reduces visual weight.
 *
 * Accessibility:
 *   - Renders as <button type="button"> by default.
 *   - Use type="submit" only inside a <form>.
 *   - Loading state sets aria-busy="true" and aria-disabled="true".
 *   - Disabled state sets aria-disabled="true"; native `disabled` is also applied.
 *
 * Usage:
 *   <BaseButton variant="primary" size="md" :isLoading="submitting">
 *     Save Changes
 *   </BaseButton>
 */
```

### What Not to Document
- Private, page-specific components with no reuse potential.
- Wrapper components that are obviously named and require no explanation.
- Components whose entire implementation is 10 lines.

---

## 9. Design System Promotion Process

A pattern becomes part of the design system when it meets the promotion criteria. This is a deliberate, reviewed process — not something any engineer does unilaterally.

### Promotion Criteria
A pattern is ready for promotion to `base/` or `common/` when:
1. It appears in **3 or more independent places** in the codebase.
2. It represents a **stable, named concept** that the team uses in conversation.
3. It has **no known upcoming redesign** that would make the extraction premature.
4. Its **props interface is stable** — no major changes expected.

### Promotion Steps
1. Create the component in `common/` or `base/` with full documentation.
2. Write tests: success state, error state, keyboard interaction (minimum).
3. Replace all existing inline implementations with the promoted component.
4. Submit for review — design system changes require senior review.
5. Communicate the addition to the team.

### Deprecation Process
When a component is superseded:
1. Mark it as deprecated in the documentation block with a migration note.
2. Do not remove it in the same PR that introduces the replacement.
3. Migrate all usages in a subsequent PR.
4. Remove the deprecated component only after all usages are migrated.

---

## 10. Preventing Component Chaos

The following anti-patterns indicate a design system in decay. During code review and system audits, you must proactively reject:
- **Duplicate components** (parallel implementations of the same UI pattern).
- **Slightly different versions of the same component** (components that vary only by visual attributes that should be handled as variants).
- **Random naming** (different terms or suffixes used for similar component patterns).
- **Components created only to wrap another component** (wrappers that add no semantic, accessibility, or layout value).
- **Components with unclear ownership** (violating hierarchy boundary and dependency rules).

### Anti-Pattern 1: Parallel Component Proliferation
```text
❌ PrimaryButton.vue
❌ BlueButton.vue
❌ LargeButton.vue
❌ DashboardButton.vue
❌ CustomerButton.vue

✅ Button.vue with variant + size props
```

**Detection**: If two component names differ only by a color, size, context, or page name — they should be one component with variants.

### Anti-Pattern 2: Copy-Pasted Markup
```text
❌ The same 30 lines of template code in 4 different files.
✅ One component, 4 usages.
```

**Detection**: Any markup block that appears more than once is a component candidate.

### Anti-Pattern 3: Wrapper Components That Add Nothing
```vue
<!-- Wrong: this component adds no value -->
<template>
  <div class="wrapper">
    <slot />
  </div>
</template>
```

**Rule**: A wrapper component must add at least one of: semantic meaning, accessibility behavior, consistent styling, state management.

### Anti-Pattern 4: Unclear Ownership
```text
❌ A feature-specific component imported directly into base/
❌ A base component that contains feature-specific business logic
❌ A layout component that fetches and manages feature data
```

**Rule**: Components must not depend on components at a higher specificity level than their own.

### Anti-Pattern 5: Random Naming
```text
❌ CardBox.vue, BoxCard.vue, CardContainer.vue (three different names for the same pattern)
✅ Card.vue (one name, consistently used everywhere)
```

**Rule**: A design system concept gets exactly one name. That name is used consistently across components, documentation, design files, and team conversation.

---

## 11. Theming Governance

When a project supports multiple themes (light/dark, brand customization, white-label), the token system is the only correct mechanism. Duplicate style sheets are not acceptable.

### Theme Rules
- All visual variation between themes must be expressible by changing token values.
- If a theme requires structural changes (different component layouts, different slot content), it is not a theme — it is a different product configuration.
- Token overrides must be scoped at the `:root` level or via a `[data-theme]` attribute. Never override tokens inside component selectors.
- Brand customization (logo, colors, fonts) must be achievable by changing a theme config file — not by editing component files.

### Theme Token Scoping
```css
/* Base theme — always defined */
:root {
  --color-primary: #0d6efd;
  --font-family-base: 'Inter', system-ui, sans-serif;
}

/* Dark mode override */
[data-theme="dark"] {
  --color-background: #0d0d0d;
  --color-surface:    #1a1a1a;
  --color-text-primary: #f8f9fa;
}

/* Brand override (white-label) */
[data-brand="acme"] {
  --color-primary:       #e63946;
  --font-family-base:    'DM Sans', system-ui, sans-serif;
}
```

---

## 12. Frontend Stack Guidance

Follow these stack-specific recommendations to ensure governance rules integrate cleanly with the project's chosen technologies:

### 12a. Pure CSS
- **Prefer**:
  - CSS variables for all design tokens (colors, spacing, typography sizes).
  - Explicit theme-scoped variable updates (`[data-theme="dark"]`).
  - Organized, modular styles grouped by component name.
  - Component-level styling ownership.
- **Avoid**:
  - Large global CSS files exceeding 200 lines.
  - Repeated raw color/padding/margin values throughout the styles.

### 12b. Tailwind CSS
- **Prefer**:
  - Referencing tokens configured in `tailwind.config.js` or `tailwind.config.ts`.
  - Component variants using custom composition structures.
  - Composing utility classes cleanly using library helpers (like `tailwind-merge` or `clsx`).
- **Avoid**:
  - Copying huge, unreadable utility class strings (12+ classes) repeatedly. Extract them.
  - Random arbitrary values (`w-[347px]`, `text-[13px]`) without product layout justification.

### 12c. Bootstrap
- **Prefer**:
  - Overriding Bootstrap standard variables in customized Sass theme sheets.
  - Native utility classes (e.g. `m-3`, `p-2`, `d-flex`) instead of duplicate custom CSS.
  - Custom themed templates that do not resemble default unstyled Bootstrap layouts.
- **Avoid**:
  - Fighting the framework selectors using `!important` overrides.
  - Introducing layout systems that contradict Bootstrap's grid without clear migration paths.

### 12d. Vue / Nuxt
- **Prefer**:
  - Explicit prop typing and boundary definitions.
  - Composables (`composables/`) for sharing stateful component logic.
  - Event-driven props/events communication (`defineEmits`) between parents and children.
- **Avoid**:
  - Components directly modifying/accessing unrelated global state or parent references ($parent, $root).
  - Unclear or hidden state mutation side effects.

### 12e. Flutter
- **Prefer**:
  - Reusable platform-consistent widgets.
  - Visual values derived strictly from `ThemeData` (e.g., `Theme.of(context).colorScheme`).
  - Proper Material Design 3 and Cupertino system defaults.
- **Avoid**:
  - Repeating custom widget trees directly inside screens without extract refactoring.

---

## Component Creation Decision Tree

Use this decision tree before creating any new UI component:

```text
Does a component already exist for this?
       │
       ├─ YES → Use it. Extend via props/slots if needed.
       │
       └─ NO → Does this appear more than once in the project?
                    │
                    ├─ NO → Implement inline. Do not extract yet.
                    │
                    └─ YES → Is there an existing component close to this?
                                  │
                                  ├─ YES → Add a variant or slot. Do not create a parallel component.
                                  │
                                  └─ NO → Create a new component.
                                              │
                                              ├─ Will it be used across multiple features?
                                              │   → Place in common/ with full documentation.
                                              │
                                              ├─ Is it a foundational primitive?
                                              │   → Place in base/ with full documentation + tests.
                                              │
                                              └─ Is it specific to one feature?
                                                  → Place in features/[feature-name]/ without promotion.
```

---

## Review Checklist

Before adding or checking in any new component, verify the following:
- [ ] **Does this already exist?** (Check existing component library/common folders to prevent duplicate components)
- [ ] **Is this the correct ownership level?** (Is it appropriately placed in `base/`, `common/`, `features/`, or `layouts/`?)
- [ ] **Is the API simple?** (API has simple inputs, clean composition, and is easy to use)
- [ ] **Are props understandable?** (Prop names are descriptive; avoids boolean prop explosion)
- [ ] **Are variants better than separate components?** (Variants, slots, or composition are used instead of parallel components)
- [ ] **Is accessibility handled?** (Keyboard support, focus states, screen reader semantics, and accessible labels are built in)
- [ ] **Is styling consistent?** (Uses design tokens for color/spacing; fits visual system)
- [ ] **Will another developer understand it?** (Naming is intuitive, code is documented where complex, and structure is clear)

---

## References
- Design Systems (tokens, token maintenance): [design/01-design-systems.md](01-design-systems.md)
- Component Libraries (evaluate before building, custom build rules): [design/02-component-libraries.md](02-component-libraries.md)
- Vue Component Discipline (naming, location, props): [stacks/js-ts-vue-nuxt/vue-components.md](../stacks/js-ts-vue-nuxt/vue-components.md)
- Accessibility (keyboard, ARIA in components): [design/04-accessibility.md](04-accessibility.md)
- UI Review Process: [design/06-ui-review-process.md](06-ui-review-process.md)
- Architecture and Simplicity: [core/02-architecture-and-simplicity.md](../core/02-architecture-and-simplicity.md)

