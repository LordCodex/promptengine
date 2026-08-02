---
document_id: design-component-libraries
title: Component Library Awareness Standard
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
  - core-architecture-and-simplicity
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Component Library Awareness Standard

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md) and the [Design Systems Standard](01-design-systems.md). It defines the decision framework for selecting, evaluating, and using UI component libraries, and the rules for building custom components only when necessary.

---

## 1. Core Rule: Evaluate Before Building

Before writing a single line of component code, confirm:

1. Does the project already have this component in its existing library?
2. Does the project's chosen component library provide this component?
3. Can an existing open-source component be composed or extended to meet the requirement?
4. Is the custom-build justified by a specific constraint (performance, branding, accessibility) that existing solutions cannot meet?

**If the answer to any of the first three questions is yes, do not build a custom component.**

Building custom versions of solved components (buttons, modals, dropdowns, date pickers) consumes engineering time, introduces bugs that library authors have already fixed, and creates maintenance burden.

---

## 2. Components That Must Not Be Built Custom (Without Justification)

The following components have mature, well-tested, accessible library implementations. Building them from scratch requires explicit architectural justification:

| Component | Complexity Risk | Why Libraries Excel |
| :--- | :--- | :--- |
| **Button** | Low | Requires focus ring, disabled state, loading state, keyboard activation |
| **Modal/Dialog** | High | Requires focus trap, scroll lock, escape key, ARIA roles, backdrop |
| **Dropdown/Select** | High | Requires keyboard nav, search filtering, ARIA combobox role, portal |
| **Date Picker** | Very High | Requires keyboard navigation, locale support, range selection, accessibility |
| **Data Table** | High | Requires sorting, pagination, keyboard navigation, responsive behaviour |
| **Toast/Notification** | Medium | Requires ARIA live regions, animation queues, z-index management |
| **Tooltip** | Medium | Requires pointer event handling, position calculation, ARIA |
| **Autocomplete** | High | Requires debounce, ARIA combobox, keyboard navigation, async loading |
| **Form Fields** | Medium | Requires validation states, ARIA `aria-describedby`, error announcement |
| **Navigation (Sidebar, Nav)** | Medium | Requires keyboard nav, active state, responsive collapse |

---

## 3. CSS Ecosystem — Evaluation by Context

### 3a. Vanilla CSS Projects

Use a structured CSS architecture. Do not write flat, uncontrolled CSS files.

**Evaluate these resources first:**

| Resource | Purpose | Best For |
| :--- | :--- | :--- |
| **Open Props** | CSS custom properties library (tokens, utilities) | Any project needing a token system fast |
| **Every Layout** | Layout primitives using intrinsic CSS | Complex, responsive layouts without media query bloat |
| **Pico CSS** | Semantic, class-light CSS framework | Content sites, forms, minimal UI |

**CSS Architecture Rules:**
- Organize CSS by component, not by page.
- All repeating values must be CSS custom properties.
- Do not write media queries for individual component properties; use container queries where supported.
- CSS files must not exceed 200 lines without explicit justification. Large CSS files are a signal of missing abstraction.
- Never use `!important` to solve specificity problems. Fix the selector hierarchy instead.

---

### 3b. Tailwind CSS Projects

Tailwind is a utility-first framework. Use it correctly or it creates unmaintainable code.

**Evaluate these component layers first:**

| Library | Purpose | When to Use |
| :--- | :--- | :--- |
| **shadcn/ui** | Copy-paste accessible components (Radix + Tailwind) | React/Next.js projects needing production-ready components |
| **DaisyUI** | Semantic component classes for Tailwind | Rapid builds, prototypes, projects not needing deep customization |
| **Headless UI** | Unstyled accessible behavior components | When you need full visual control with correct accessibility behavior |
| **Radix UI** | Unstyled primitives with strong accessibility | Complex interactive components, design systems |

**Tailwind Usage Rules:**
- Extract repeated utility class sets into component classes using `@apply` or dedicated component files.
- Never write Tailwind utility class strings longer than 10–12 classes inline. If a string exceeds this, it is a component that needs a name.
- Configure the Tailwind theme file to match the project design tokens. Do not override colors, spacing, or typography ad-hoc in templates.
- Avoid using arbitrary values (`w-[347px]`, `text-[13px]`) except for one-off constraints that genuinely cannot be expressed in the scale.

---

### 3c. Bootstrap Projects

Bootstrap provides a complete component system and grid. Use it as intended — not as a starting point for reinventing its components.

**Bootstrap Usage Rules:**
- Customize Bootstrap via Sass variables before compilation, not by overriding classes after the fact.
- Use Bootstrap's utility classes (margin, padding, display, flex) before writing custom CSS for the same purpose.
- Never produce default unstyled Bootstrap output in production. Customize the color palette, typography, and border radius at minimum.
- Use Bootstrap Icons when the project uses Bootstrap to avoid importing a second icon library.
- Do not mix Bootstrap's grid with a custom CSS grid unless the project explicitly documents a migration path.

---

### 3d. Vue / Nuxt Projects

**Evaluate these libraries before building custom components:**

| Library | Strengths | When to Consider |
| :--- | :--- | :--- |
| **PrimeVue** | Data-heavy components, tables, charts, forms | Admin UIs, dashboards, enterprise apps |
| **Vuetify** | Material Design alignment, broad component set | Projects aligned with Google/Android design language |
| **Naive UI** | TypeScript-first, composable | TypeScript-strict Vue projects |
| **Element Plus** | Form-heavy, config-rich | CMS-style, admin panels |
| **Headless UI (Vue)** | Accessible unstyled behaviors | Custom-styled, accessible components |

**Vue Component Rules:**
- Do not install more than one component library per project without explicit team decision.
- When a library component handles 80% of the use case, use it and extend it — do not rebuild it.
- Evaluate accessibility (keyboard nav, ARIA) and active maintenance (last release, issues) before committing.
- Library component APIs must be wrapped in project-specific components for portability (so switching libraries later does not require changes in every template).

---

### 3e. Flutter / Dart Projects

**Flutter's material library is the primary design system.** Evaluate it before building custom widgets.

| Source | Purpose | When to Use |
| :--- | :--- | :--- |
| **Material 3 Widgets** | Platform-consistent, theme-aware | Default for all Android/cross-platform apps |
| **Cupertino Widgets** | iOS-native appearance | iOS-specific screens or adaptive UIs |
| **flutter_animate** | Animation library | When explicit animation control is needed |
| **Lottie** | JSON animation playback | Complex, design-authored animations |

**Flutter Widget Rules:**
- Prefer `Material 3` widget set and `ThemeData` tokens over hardcoded styling.
- Always use `const` constructors for stateless widgets to prevent unnecessary rebuilds.
- Never set hardcoded `Color` values on widgets. Route all colors through `Theme.of(context).colorScheme`.
- Never hardcode text styles. Route all typography through `Theme.of(context).textTheme`.
- Custom widgets must accept a `style` or equivalent override pattern — never hardcode visual decisions inside widget implementations.

---

## 4. Icon Library Selection

Icons are a component category that requires deliberate selection.

### Rules
- Choose one icon library per project and use it exclusively.
- Prefer icon libraries that provide consistent weight, optical size, and style across all icons.
- Do not mix icon styles (filled + outlined in the same context) without a documented visual system reason.
- SVG icons are preferred over icon fonts for accessibility and performance.

### Common Libraries

| Library | Style | When to Use |
| :--- | :--- | :--- |
| **Lucide** | Outline, consistent weight | Clean modern UIs |
| **Heroicons** | Outline and solid variants | Tailwind/shadcn projects |
| **Phosphor Icons** | Multiple weights | Design systems needing flexibility |
| **Material Symbols** | Variable font, Google ecosystem | Flutter / Material Design projects |
| **Bootstrap Icons** | Outline | Bootstrap projects |
| **Feather Icons** | Simple outline | Minimal UIs |

---

## 5. Custom Component Build Rules

When a custom component is genuinely required, apply these engineering rules:

### Before Building
- Document why existing solutions are insufficient.
- Define the component's props interface and emitted events before implementation.
- Write the component's accessibility requirements (keyboard behaviour, ARIA roles, focus management) before writing the template.

### During Building
- The component must handle all required states: default, loading, empty, error, disabled, success (where applicable).
- Props must be typed and documented.
- The component must be focusable and operable by keyboard.
- The component must not hardcode visual values — use design tokens.
- The component must be placed in the correct directory (see [vue-components.md](../stacks/js-ts-vue-nuxt/vue-components.md) for location conventions).

### After Building
- Write at minimum one unit test covering the success state, one covering the error state, and one covering keyboard interaction.
- Document the component's API in the component file's top-level comment block.

---

## Review Checklist
- [ ] Has an existing library been evaluated before building a custom component?
- [ ] Is the component library selection documented with a stated reason?
- [ ] Are Tailwind class strings extracted into components when they exceed 10–12 utilities?
- [ ] Are custom CSS values references to design tokens, not hardcoded?
- [ ] Is only one icon library in use?
- [ ] Do custom components handle all required states (loading, error, empty, disabled)?
- [ ] Are custom component props typed and documented?

---

## References
- Design Systems: [design/01-design-systems.md](01-design-systems.md)
- Vue Component Discipline: [stacks/js-ts-vue-nuxt/vue-components.md](../stacks/js-ts-vue-nuxt/vue-components.md)
- Accessibility: [design/04-accessibility.md](04-accessibility.md)
- Design Resources: [design/07-design-resources.md](07-design-resources.md)
