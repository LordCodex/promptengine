---
document_id: design-resources
title: Design Resources and Inspiration Guide
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Design Resources and Inspiration Guide

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md). It curates resources for design inspiration, learning, tooling, and pattern research — with clear rules for how inspiration may and may not be used.

---

## 0.1. Design Research Workflow

Before designing any medium or large user interface:
1. **Understand the product**: What is the core value? What are the key objects and behaviors?
2. **Understand the target users**: Identify user roles, technical capabilities, context, and usage frequency (see [design/08-product-thinking.md](08-product-thinking.md)).
3. **Understand the business goal**: What constitutes completion? What conversion, retention, or speed metric matters?
4. **Research similar products**: Identify competitors and market benchmarks to map user expectations.
5. **Identify common interaction patterns**: Notice how industry standards organize navigation, headers, forms, and data.
6. **Combine the strongest ideas**: Extract layout properties, hierarchy elements, and user states.
7. **Adapt to the design system**: Implement the synthesized solution using the project's design tokens and components.

Never start coding or rendering layouts without completing this workflow.

---

## 0.2. Product Benchmarking

When researching user expectations, analyze mature products solving similar problems in related domains:

- **Financial Services / SaaS Billing** → *Stripe*, *Wise*, *Revolut* (Optimize for: absolute clarity, fee transparency, payment state confidence, transaction feedback).
- **Project Management / Dense Dashboards** → *Linear*, *Notion*, *Jira* (Optimize for: keyboard shortcuts, high data density, clean filtering, clean navigation states).
- **E-Commerce / Catalogs** → *Shopify*, *Amazon*, *Gumroad* (Optimize for: clear product imagery, readable pricing details, trust indicators, frictionless checkouts).
- **Collaboration / Real-Time Messaging** → *Slack*, *Discord* (Optimize for: channel layouts, keyboard shortcuts, immediate message updates, persistent indicators).
- **Digital Content / Reading Interfaces** → *Kindle*, *Apple Books*, *Kobo* (Optimize for: typography layout limits, high contrast read modes, clean navigation focus).

**Rule**: Study the interaction patterns and information architecture of these products — never copy their branding, colors, or visual style.

---

## 1. Rules for Using Design Inspiration

Design inspiration must be used for **analysis and learning**, not copying.

### Permitted Use
- Study layout patterns and information hierarchy decisions.
- Analyze how successful products organize navigation and primary actions.
- Observe interaction patterns (how modals are triggered, how forms are structured).
- Understand how spacing, typography, and color combine to create hierarchy.
- Identify component patterns that solve the same problem in the current project.

### Prohibited Use
- Directly copying a design element, color palette, or layout and presenting it as original.
- Reproducing a competitor's visual language without a documented business reason.
- Using premium marketplace designs or Figma templates without licensing compliance.
- Adopting a trending visual style that does not fit the user's context (e.g., gaming UI for a healthcare app).

---

## 2. Design Inspiration Sources

### UI Pattern Libraries
Sources for studying real product interfaces and interaction patterns:

| Source | What It Provides |
| :--- | :--- |
| **Mobbin** ([mobbin.com](https://mobbin.com)) | Mobile and web UI screenshots from real apps, searchable by pattern |
| **SaaSFrame** ([saasframe.io](https://saasframe.io)) | SaaS product UI patterns — onboarding, pricing, dashboards, emails |
| **Pageflows** ([pageflows.com](https://pageflows.com)) | Recorded user flows from real products (onboarding, checkout, etc.) |
| **Screenlane** ([screenlane.com](https://screenlane.com)) | Curated mobile UI screenshots by component type |
| **UI Sources** ([uisources.com](https://uisources.com)) | iOS and Android patterns from real apps |
| **UI Garage** ([uigarage.net](https://uigarage.net)) | Categorized interface design patterns and visual assets |

### Design Community & Galleries
Sources for design quality, trends, and emerging patterns:

| Source | What It Provides |
| :--- | :--- |
| **Dribbble** ([dribbble.com](https://dribbble.com)) | Design portfolio work — evaluate critically for production viability |
| **Behance** ([behance.net](https://behance.net)) | Detailed case studies including process, rationale, and outcomes |
| **Awwwards** ([awwwards.com](https://awwwards.com)) | Award-winning web design — experimental and innovative |
| **Figma Community** ([figma.com/community](https://figma.com/community)) | Free design files, UI kits, and templates |
| **Land-book** ([land-book.com](https://land-book.com)) | Curation of top landing pages, cataloged by category and element |
| **Lapa Ninja** ([lapa.ninja](https://www.lapa.ninja)) | Landing page gallery with focus on SaaS and product designs |
| **One Page Love** ([onepagelove.com](https://onepagelove.com)) | Specialized curation of single-page websites and templates |
| **Godly** ([godly.website](https://godly.website)) | Gallery highlighting interactions, animation flow, and visual effects |

### Design System References
Study these mature systems to borrow their structural principles, accessibility rules, and state patterns rather than copy their appearances:

| System | Primary Context | Focus |
| :--- | :--- | :--- |
| **Material Design 3** | Android / Web / Flutter | Interactive material states, fluid layouts, color schemes |
| **Apple HIG** | iOS / macOS / watchOS | Platform conventions, layouts, typography, safe areas |
| **Microsoft Fluent** | Windows / Web | Depth, light, material textures, desktop hierarchy |
| **IBM Carbon** | Enterprise SaaS / Admin | Data density, structured grids, extreme validation rules |
| **Shopify Polaris** | E-commerce Merchant | Retail management workflows, forms, search patterns |
| **Atlassian Design** | Collaboration / Jira | Multi-team workflows, boards, status indications |
| **Ant Design** | Enterprise / React | Comprehensive grid-based components, admin panels |
| **Chakra UI / Mantine** | Web (React/TypeScript) | Composable component configurations, theme tokens |
| **Radix UI** | Accessible Web (React) | Headless primitives with complete keyboard/ARIA support |
| **shadcn/ui** | Composable Web | CSS variable themed configurations |

---

## 3. Component and Styling References

When building interfaces, prefer using proven component patterns and tokens from these ecosystems rather than inventing new layouts.

### 3a. Tailwind Ecosystem
- **Tailwind UI** ([tailwindui.com](https://tailwindui.com)) — Official composable layout and component templates.
- **shadcn/ui** ([ui.shadcn.com](https://ui.shadcn.com)) — Accessible, styled primitive blocks.
- **Flowbite** ([flowbite.com](https://flowbite.com)) — Composed utility classes for standard components.
- **Preline UI** ([preline.co](https://preline.co)) — Prebuilt Tailwind component blocks.
- **DaisyUI** ([daisyui.com](https://daisyui.com)) — Semantic class selectors for Tailwind.
- **HyperUI** ([hyperui.dev](https://hyperui.dev)) — Open-source Tailwind component blocks.
- **Meraki UI** ([merakiui.com](https://merakiui.com)) — Responsive Tailwind UI component blocks.
- **TailGrids** ([tailgrids.com](https://tailgrids.com)) — Grid structures and layout blocks.
- **Float UI** ([floatui.com](https://floatui.com)) — Composed forms, tables, and dashboards.
- **Kutty** ([kutty.netlify.app](https://kutty.netlify.app)) — Prebuilt Tailwind elements.

### 3b. Bootstrap Ecosystem
- **Bootstrap Examples** ([getbootstrap.com/docs/5.3/examples/](https://getbootstrap.com/docs/5.3/examples/)) — Official layout starters.
- **MDBootstrap** ([mdbootstrap.com](https://mdbootstrap.com)) — Material style Bootstrap blocks.
- **CoreUI / AdminLTE** — Production-ready admin dashboard layouts.
- **Tabler** ([tabler.io](https://tabler.io)) — Highly modular, clean Bootstrap 5 admin layouts.
- **Sneat** — Modern admin templates.
- **Volt** — Bootstrap 5 admin dashboard theme.
- **Bootswatch** ([bootswatch.com](https://bootswatch.com)) — Customized theme sheets.
- **Start Bootstrap** ([startbootstrap.com](https://startbootstrap.com)) — Public templates.

### 3c. Pure CSS
- **Pico CSS** ([picocss.com](https://picocss.com)) — Minimal classless semantic styling.
- **Open Props** ([open-props.style](https://open-props.style)) — Design tokens packaged as CSS variables.
- **Water.css / Classless CSS** — Write raw semantic HTML; style sheets render them cleanly.
- **Modern Normalize** ([github.com/sindresorhus/modern-normalize](https://github.com/sindresorhus/modern-normalize)) — Reset styles.
- **Every Layout** ([every-layout.dev](https://every-layout.dev)) — Layout primitives without media query hacks.
- **Layout Land** ([layout.land](https://layout.land)) — Grid and Flexbox principles by Jen Simmons.

### 3d. Vue / Nuxt / Flutter Framework Libraries
- **Vue / Nuxt**: *PrimeVue* ([primevue.org](https://primevue.org)), *Vuetify* ([vuetifyjs.com](https://vuetifyjs.com)), *Naive UI* ([naiveui.com](https://www.naiveui.com)), *Element Plus* ([element-plus.org](https://element-plus.org)).
- **Flutter**: ThemeData-scoped native Material and Cupertino widget sets.

---

## 4. Typography Resources

### Font Selection
| Resource | Purpose |
| :--- | :--- |
| **Google Fonts** ([fonts.google.com](https://fonts.google.com)) | Free, high-quality web fonts with performance optimization |
| **Font Share** ([fontshare.com](https://www.fontshare.com)) | Free professional fonts from Indian Type Foundry |
| **Bunny Fonts** ([fonts.bunny.net](https://fonts.bunny.net)) | GDPR-compliant alternative to Google Fonts |
| **Type Scale** ([typescale.com](https://typescale.com)) | Generate and visualize typographic scale ratios |

### Recommended Starting Typefaces (Free)
| Typeface | Character | Best For |
| :--- | :--- | :--- |
| **Inter** | Geometric, neutral | SaaS, admin, dashboards |
| **Plus Jakarta Sans** | Friendly, modern | Consumer apps, marketing |
| **DM Sans** | Clean, legible | Content sites, documentation |
| **Figtree** | Rounded, approachable | Friendly apps, onboarding |
| **Geist** | Technical, monospace-inspired | Developer tools, code-adjacent |
| **Lato** | Classic, neutral | General purpose |
| **Source Sans 3** | Readable at small sizes | Data-heavy interfaces |

---

## 5. Icon Libraries

| Library | Style | License | Best For |
| :--- | :--- | :--- | :--- |
| **Lucide** ([lucide.dev](https://lucide.dev)) | Outline, consistent | MIT | Clean modern UI |
| **Heroicons** ([heroicons.com](https://heroicons.com)) | Outline + Solid | MIT | Tailwind/shadcn projects |
| **Phosphor Icons** ([phosphoricons.com](https://phosphoricons.com)) | Multiple weights | MIT | Design systems needing flexibility |
| **Material Symbols** ([fonts.google.com/icons](https://fonts.google.com/icons)) | Variable font | Apache 2.0 | Flutter, Material Design |
| **Bootstrap Icons** ([icons.getbootstrap.com](https://icons.getbootstrap.com)) | Outline | MIT | Bootstrap projects |
| **Tabler Icons** ([tabler.io/icons](https://tabler.io/icons)) | Outline, comprehensive | MIT | Admin UIs |
| **Remix Icon** ([remixicon.com](https://remixicon.com)) | Dual-tone options | Apache 2.0 | General purpose |

---

## 6. Color Resources

### Color Palette Generation
| Tool | Purpose |
| :--- | :--- |
| **Coolors** ([coolors.co](https://coolors.co)) | Generate and explore color palettes |
| **Realtime Colors** ([realtimecolors.com](https://www.realtimecolors.com)) | Visualize palettes on real UI layouts |
| **oklch.com** ([oklch.com](https://oklch.com)) | Generate perceptually uniform colors using OKLCH color space |
| **Accessible Palette** ([accessiblepalette.com](https://accessiblepalette.com)) | Build palettes that meet WCAG contrast requirements |

### Color Contrast Checking
| Tool | Purpose |
| :--- | :--- |
| **WebAIM Contrast Checker** ([webaim.org/resources/contrastchecker](https://webaim.org/resources/contrastchecker)) | Verify WCAG contrast ratios |
| **Who Can Use** ([whocanuse.com](https://whocanuse.com)) | Visualize how color combinations work for different vision types |
| **Colorblindly** (Chrome extension) | Simulate color blindness in real browser sessions |

---

## 7. Accessibility Tools

| Tool | Type | Purpose |
| :--- | :--- | :--- |
| **axe DevTools** ([deque.com/axe](https://www.deque.com/axe/)) | Browser extension | Automated accessibility testing |
| **Lighthouse** (Chrome DevTools) | Built-in | Automated accessibility + performance audit |
| **WAVE** ([wave.webaim.org](https://wave.webaim.org)) | Browser extension | Visual accessibility annotations |
| **NVDA** ([nvaccess.org](https://www.nvaccess.org)) | Screen reader (Windows) | Free screen reader for testing |
| **VoiceOver** (macOS/iOS built-in) | Screen reader | Native macOS and iOS screen reader |
| **ARC Toolkit** | Browser extension | Detailed structure and ARIA analysis |

---

## 8. Design Decision Documentation

When a non-obvious design decision is made, document it in the component or the project's design decisions log using this format:

```markdown
## Decision: [Short title]

**Date**: YYYY-MM-DD  
**Context**: What situation required this decision?  
**Options Considered**:
- Option A — [summary and tradeoff]
- Option B — [summary and tradeoff]

**Decision**: [Which option was chosen and why]  
**Consequences**: [What this means for the future — what it makes easier or harder]
```

This is especially important for:
- Choosing a custom component over a library component.
- Deviating from the project's established visual pattern.
- Adopting a new icon library or typeface.
- Establishing a new layout pattern.

---

## 8.5. Framework-Aware UI Patterns

Always align visual designs with the architectural and structural conventions of the target frontend or mobile framework. Never force the design style of one environment onto another.

- **Laravel Blade**: Optimize for server-side rendered layouts, traditional document flows, and session validation alerts.
- **Vue / Nuxt**: Optimize for composable component trees, client-side dynamic states, and Server-Side Rendering (SSR) page hydration constraints.
- **React / Next.js**: Maintain strict division between Server Components and Client Components; design state and loading hooks to align with React rendering lifecycles.
- **Flutter**: Ensure layout structures adapt to platform-consistent widgets (Material Design for Android/Cupertino for iOS), using widgets that scale smoothly inside device viewports.

---

## 9. Design Synthesis & Originality

A production user interface must never be a direct copy of a single inspiration source. Instead, construct your design through synthesis:
1. **Analyze multiple references**: Review 3–5 different benchmarks to map alternative solutions.
2. **Extract distinct sub-elements**:
   - Navigation configurations (header style vs sidebar layout).
   - Spacing density and margin sizing.
   - Information priorities (e.g. data hierarchy on cards).
   - Form and input layouts.
   - Tables and list arrangements.
   - User states (loading, empty, errors).
   - Search controls and filters.
   - Mobile-specific interactions.
3. **Coalesce into a unified experience**: Combine these features into a single interface that reflects the project's visual system and design tokens. The final product should be recognizable only as a polished solution with its own identity.

---

## 10. Design Consistency & Modern Characteristics

### Modern UI Characteristics
Production layouts should feel **calm, spacious, fast, focused, consistent, and trustworthy**.
- Avoid excessive shadows or high-blur elevations.
- Avoid dynamic gradients applied for purely decorative purposes.
- Avoid over-rounding all elements without visual hierarchy.
- Avoid adding features or visual decorations that have no product purpose.

### Consistency Scales
Ensure all components inside the application enforce:
- **One spacing scale** (e.g., standard pixel gutters).
- **One typography scale** (matching heading hierarchies).
- **One radius system** (uniform corner styles).
- **One elevation system** (shadow/depth logic).
- **One icon style** (matching stroke weights and families).
- **One button and form style**.

---

## 11. Research Before Invention

Do not write custom code to solve common interface problems unless explicitly required. Before building, check if a mature solution exists in your library or framework ecosystem. Do not reinvent:
- Complex tabular grids and pagination models.
- Date and range picker calendars.
- Dropdown select and autocomplete inputs.
- Multi-select controls.
- Rich text editors.
- Advanced drag-and-drop file uploaders.
- Command palette navigation widgets.

---

## 12. UI Research Web Browsing Rule

> [!IMPORTANT]
> When web access is available and the user has not prohibited browsing, you must research current design patterns and official design system documentation before proposing novel UI patterns. Prefer official documentation over screenshots or trend sites for implementation details.

---

## Design Inspiration Review Checklist

Before completing UI design synthesis, verify the interface against this checklist:
- [ ] **Product competitors studied**: Analyzed mature solutions to map industry standards and user expectations.
- [ ] **Multiple references considered**: Evaluated 3+ benchmarks rather than copying or imitating a single source.
- [ ] **No direct copying**: The design is synthesized from principles; visual identity is original to the project.
- [ ] **Components consistent**: Spacing, borders, colors, and controls follow a unified style guide.
- [ ] **Design system respected**: Design tokens are used for all colors, fonts, spacing, and sizing variables.
- [ ] **Framework conventions followed**: UI layout structures fit the selected framework (Blade, Vue, React, Flutter).
- [ ] **Mature interaction patterns used**: Relies on standard patterns for navigation, tables, and dropdowns rather than reinventing them.
- [ ] **Modern visual hierarchy**: High-priority information is visually dominant; no competing primary controls.
- [ ] **Responsive layouts**: Viewport adaptions maintain visual priority; touch target sizes are verified.
- [ ] **Accessibility maintained**: Focus states, markup semantics, and contrast ratios are built-in.
- [ ] **User goals prioritized**: Visual aesthetics serve to reduce cognitive load and facilitate task completion.

---

## References
- UI/UX Philosophy: [design/00-ui-ux-philosophy.md](00-ui-ux-philosophy.md)
- Design Systems: [design/01-design-systems.md](01-design-systems.md)
- Component Libraries: [design/02-component-libraries.md](02-component-libraries.md)
- Accessibility Standard: [design/04-accessibility.md](04-accessibility.md)
