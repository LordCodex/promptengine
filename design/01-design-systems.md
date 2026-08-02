---
document_id: design-systems
title: Design Systems Standard
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Design Systems Standard

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md) and the [Universal Coding Standards](../core/05-universal-coding-standards.md). It defines how design systems are discovered, established, and maintained in projects across the playbook ecosystem.

---

## 1. What Is a Design System

A design system is a single source of truth for all visual and interaction decisions in a product. It consists of:

- **Design Tokens**: Named variables for colors, spacing, typography, shadows, border radii.
- **Component Library**: A set of reusable UI components built on top of the tokens.
- **Patterns**: Documented solutions for common UI problems (forms, navigation, data display).
- **Principles**: The product's design philosophy and voice.

A design system is not a collection of screenshots. It is a living, versioned engineering artifact.

---

## 2. Design System First Protocol

Before creating any UI component or writing any styling code, run this discovery protocol:

```text
1. Does this project already have a design system?
2. Is there an existing component library (internal or third-party)?
3. Are design tokens defined (CSS variables, Tailwind config, theme file)?
4. Are typography rules already established?
5. Are spacing and sizing scales already defined?
6. Are color roles already named and documented?
7. Are existing patterns documented and available?
```

If the answer to any question is **yes**, use what exists. Do not create parallel definitions.

If the answer to all questions is **no**, establish the design system first before building components.

---

## 3. Design Tokens

Design tokens are the foundation of a consistent UI. Every value that repeats across the design must be a token.

### Required Token Categories

#### Color Tokens
```css
/* Semantic color roles — not raw values */
--color-background:       #ffffff;
--color-surface:          #f8f9fa;
--color-border:           #dee2e6;
--color-text-primary:     #212529;
--color-text-secondary:   #6c757d;
--color-text-disabled:    #adb5bd;

/* Intent colors */
--color-primary:          #0d6efd;
--color-primary-hover:    #0b5ed7;
--color-success:          #198754;
--color-warning:          #ffc107;
--color-danger:           #dc3545;
--color-info:             #0dcaf0;

/* Interactive states */
--color-focus-ring:       rgba(13, 110, 253, 0.25);
```

> [!IMPORTANT]
> Define colors as **semantic roles**, not raw values. A button background should reference `--color-primary`, not `#0d6efd`. This enables theming and dark mode without touching component code.

#### Spacing Tokens
Use a consistent spacing scale. The most common approach is a base-4 or base-8 scale:

```css
--space-1:   4px;
--space-2:   8px;
--space-3:   12px;
--space-4:   16px;
--space-6:   24px;
--space-8:   32px;
--space-12:  48px;
--space-16:  64px;
--space-24:  96px;
```

Never use arbitrary pixel values in component code. Always use a spacing token.

#### Typography Tokens
```css
--font-family-base:       'Inter', system-ui, -apple-system, sans-serif;
--font-family-mono:       'Fira Code', 'Cascadia Code', monospace;

--font-size-xs:           0.75rem;   /* 12px */
--font-size-sm:           0.875rem;  /* 14px */
--font-size-base:         1rem;      /* 16px */
--font-size-lg:           1.125rem;  /* 18px */
--font-size-xl:           1.25rem;   /* 20px */
--font-size-2xl:          1.5rem;    /* 24px */
--font-size-3xl:          1.875rem;  /* 30px */
--font-size-4xl:          2.25rem;   /* 36px */

--font-weight-normal:     400;
--font-weight-medium:     500;
--font-weight-semibold:   600;
--font-weight-bold:       700;

--line-height-tight:      1.25;
--line-height-normal:     1.5;
--line-height-relaxed:    1.75;
```

#### Border & Shape Tokens
```css
--radius-sm:   4px;
--radius-md:   8px;
--radius-lg:   12px;
--radius-xl:   16px;
--radius-full: 9999px;

--border-width: 1px;
--border-color: var(--color-border);
```

#### Shadow Tokens
```css
--shadow-sm:  0 1px 2px 0 rgba(0, 0, 0, 0.05);
--shadow-md:  0 4px 6px -1px rgba(0, 0, 0, 0.1);
--shadow-lg:  0 10px 15px -3px rgba(0, 0, 0, 0.1);
--shadow-xl:  0 20px 25px -5px rgba(0, 0, 0, 0.1);
```

---

## 4. Typography System

### Scale Principle
A typography scale defines the visual hierarchy of a page. Without a deliberate scale, type sizes are arbitrary and hierarchy is unclear.

### Required Type Levels

| Level | Token | Use Case |
| :--- | :--- | :--- |
| Display | `--font-size-4xl` / `3xl` | Page heroes, landing H1 |
| H1 | `--font-size-3xl` | Page primary title (one per page) |
| H2 | `--font-size-2xl` | Section headings |
| H3 | `--font-size-xl` | Sub-section headings |
| Body | `--font-size-base` | All running text |
| Small | `--font-size-sm` | Captions, metadata, helper text |
| Micro | `--font-size-xs` | Badges, timestamps, labels |

### Typography Rules
- One `<h1>` per page — no exceptions.
- Heading levels must be sequential — do not skip from `h1` to `h4`.
- Do not use heading elements for visual size. Use heading elements for semantic structure; control size with CSS.
- Body text minimum: 16px (1rem). Never set body text below 14px.
- Line height for body text: 1.5–1.75.
- Maximum line length for readable text: 65–75 characters (approximately 38em).

---

## 5. Spacing System

### The Golden Rule of Spacing
Spacing creates relationships. Items that are close together appear related. Items that are far apart appear unrelated.

Use spacing to communicate grouping — not to fill whitespace.

### Spacing Application Rules
- Use the spacing token scale exclusively. Never hardcode arbitrary pixel values.
- Spacing between related elements: 4–8px.
- Spacing between grouped sections: 16–24px.
- Spacing between distinct sections: 32–64px.
- Padding inside interactive components: minimum 8px vertical, 12px horizontal.
- Touch targets on mobile: minimum 44×44px.

---

## 6. Dark Mode

Dark mode is an accessibility requirement for many users, not a visual trend.

### Dark Mode Token Strategy
```css
:root {
  --color-background:    #ffffff;
  --color-surface:       #f8f9fa;
  --color-text-primary:  #212529;
}

@media (prefers-color-scheme: dark) {
  :root {
    --color-background:    #0d0d0d;
    --color-surface:       #1a1a1a;
    --color-text-primary:  #f8f9fa;
  }
}
```

If the project uses JavaScript-controlled dark mode:
```css
[data-theme="dark"] {
  --color-background:    #0d0d0d;
  --color-surface:       #1a1a1a;
  --color-text-primary:  #f8f9fa;
}
```

### Dark Mode Rules
- Design tokens, not hardcoded values, make dark mode possible. If values are hardcoded, dark mode requires rewriting every component.
- Test color contrast in both light and dark modes. WCAG AA minimum: 4.5:1 for body text, 3:1 for large text and UI components.
- Never use pure black (`#000000`) backgrounds or pure white (`#ffffff`) text in dark mode — it creates harsh contrast.

---

## 7. Existing Design Systems to Evaluate

Before building a custom design system, evaluate whether an established system fits the project:

| System | Best Fit | CSS Approach |
| :--- | :--- | :--- |
| **Material Design 3** | Flutter apps, Google ecosystem | Material You tokens |
| **shadcn/ui** | React/Next.js projects | Tailwind + CSS variables |
| **Radix UI** | React, accessibility-critical UI | Headless + custom CSS |
| **DaisyUI** | Rapid prototyping with Tailwind | Tailwind component classes |
| **PrimeVue** | Vue applications, data-heavy UIs | Component + theme system |
| **Vuetify** | Vue, Material Design alignment | Material component library |
| **Naive UI** | Vue, TypeScript-first | Component + theme overrides |
| **Open Props** | Vanilla CSS, any framework | CSS custom properties |
| **Pico CSS** | Minimal, semantic HTML sites | Class-light CSS |

Evaluate using:
1. Does the library's component set cover 80%+ of the project's needs?
2. Is it actively maintained with recent releases?
3. Does it support the project's accessibility requirements?
4. Can it be themed to match the project's design language?
5. What is the bundle size impact?

---

## 8. Design System Maintenance Rules

- Never modify design tokens component-by-component. Changes to tokens propagate automatically.
- Token additions are always safe. Token removal is a breaking change — treat it as one.
- Document every token with its intended use, what it affects, and what values are forbidden.
- When a pattern is implemented twice, extract it into the design system. The third implementation is a violation.
- Design system changes require the same level of review as API changes.

---

## Review Checklist
- [ ] Are design tokens defined for colors, spacing, typography, borders, and shadows?
- [ ] Are semantic color roles used (not raw values) in all component code?
- [ ] Is the typography scale defined and applied consistently?
- [ ] Is there no arbitrary pixel value outside the spacing token scale?
- [ ] Is dark mode supported through token definitions?
- [ ] Does the project use an appropriate existing design system or library?
- [ ] Are duplicate component implementations avoided?

---

## References
- Design System Governance: [design/09-design-system-governance.md](09-design-system-governance.md)
- Component Libraries: [design/02-component-libraries.md](02-component-libraries.md)
- Visual Quality: [design/05-visual-quality.md](05-visual-quality.md)
- Accessibility: [design/04-accessibility.md](04-accessibility.md)
- Open Props: [https://open-props.style](https://open-props.style)
- Material Design 3: [https://m3.material.io](https://m3.material.io)
