---
document_id: checklists-ui-quality-review
title: UI Quality Review Checklist
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
  - design-component-libraries
  - design-responsive-design
  - design-accessibility
  - design-visual-quality
  - design-ui-review-process
audience: [human, agent]
last_reviewed: 2026-08-01
---

# UI Quality Review Checklist

Use this checklist after implementing any frontend feature, page, or component. Run every section in order. A frontend implementation is **not complete** until every applicable item passes.

> [!IMPORTANT]
> This checklist is a gate, not a suggestion. An AI agent or developer must not mark frontend work as done without completing this review. Failures must be fixed before submission, not flagged as follow-up tasks.

For the full explanation of each category, see [design/06-ui-review-process.md](../design/06-ui-review-process.md).

---

## Section 1 — User Experience

### Primary Goal
- [ ] Is the primary user goal immediately obvious without explanation?
- [ ] Can a new user understand what to do without reading instructions?
- [ ] Is the most important action visually prominent on the screen?

### User Flow
- [ ] Is the natural sequence of actions reflected in the screen sequence?
- [ ] Are unnecessary steps removed? (No extra clicks or confirmations without a safety reason)
- [ ] Is feedback provided after every user action (submit, save, delete, send)?
- [ ] Do destructive or irreversible actions have a confirmation step?

### States — All must be implemented before work is complete
- [ ] **Loading state**: A visible indicator (skeleton or spinner) appears while content is fetching.
- [ ] **Empty state**: An explanatory message appears when there is no data. It tells the user why and what they can do.
- [ ] **Error state**: A plain-language message explains what went wrong and provides a recovery action where possible.
- [ ] **Success state**: The user receives confirmation that their action completed.
- [ ] **Disabled state**: Disabled elements are visually distinct and — where helpful — explain why they are disabled.

---

## Section 2 — AI-Generated Pattern Detection

Review the implementation for the following generic AI-generated patterns. Any item found must be justified with an explicit design reason or replaced.

- [ ] **No random gradients** applied to backgrounds, buttons, or containers without a documented brand or depth purpose.
- [ ] **No excessive rounded containers** — `border-radius` is applied intentionally, not as a default to every element.
- [ ] **No card overuse** — cards are used only when content items are discrete, independent, and comparable. Not every content section is a card.
- [ ] **No fake statistics sections** — hero stat numbers (e.g., "10k+ users", "99% uptime") are real data, not decorative filler.
- [ ] **No decorative icons without purpose** — every icon either navigates, triggers an action, or provides meaningful context.
- [ ] **No generic dashboard layout copied into a non-dashboard context** — the layout matches the actual product type.
- [ ] **No unnecessary animations** — every animation communicates state change, relationship, direction, or feedback.
- [ ] **No copy-paste SaaS landing page structure** in an actual application screen (hero + features + pricing + CTA when not warranted).
- [ ] **No purple/blue generic SaaS palette** applied without reference to brand or product context.
- [ ] **No large empty whitespace** between elements that should be grouped together.

If a pattern from this list is present and intentional, document the reason in a comment or design decision note.

---

## Section 3 — Visual Design

### Layout and Spacing
- [ ] All spacing values use the project's spacing token scale — no arbitrary pixel values.
- [ ] Related elements are visually grouped by proximity (close spacing = related, large spacing = separate).
- [ ] All elements align to a consistent axis. No elements positioned by eye.
- [ ] No inconsistent or random margins between sibling elements.
- [ ] Section padding is consistent and uses the spacing scale.

### Typography
- [ ] Body text is at least 16px with a line height of 1.5 or greater.
- [ ] Headings and body text are clearly differentiated in size and weight.
- [ ] No more than 3–4 distinct type sizes are visible on a single screen.
- [ ] Paragraph line length does not exceed approximately 75 characters on wide viewports.
- [ ] Text is left-aligned (not justified) for running content.

### Colors
- [ ] All colors reference design tokens — no hardcoded hex or `rgb()` values in component code.
- [ ] Color is never the **only** way to communicate state or meaning (error, success, warning).
- [ ] Status colors (danger, success, warning) are used consistently across the implementation.
- [ ] The palette does not exceed 5 distinct hues in a single screen context without justification.

### Component Visual Language Consistency
Verify each component type present in the implementation follows the same visual language:
- [ ] **Buttons** — consistent size, shape, weight, and state styles.
- [ ] **Forms / Inputs** — consistent label placement, border style, and error presentation.
- [ ] **Cards** — consistent padding, border, shadow, and hover behavior.
- [ ] **Tables** — consistent header, row, and cell styling.
- [ ] **Navigation** — consistent active state, spacing, and icon alignment.
- [ ] **Modals** — consistent header, body, footer structure, and close mechanism.
- [ ] **Notifications / Toasts** — consistent position, color use, and dismiss behavior.

---

## Section 4 — Accessibility

### Structural
- [ ] The correct HTML element is used for every interactive component.
  - Clickable actions use `<button>`, not `<div onclick>` or `<span onclick>`.
  - Navigation links use `<a href>`, not `<div onclick>`.
  - Forms use `<input>`, `<select>`, `<textarea>` — not `<div contenteditable>`.
- [ ] There is exactly one `<h1>` per page.
- [ ] Heading levels are sequential (no skipping from `h1` to `h4`).
- [ ] Landmark regions are present: `<main>`, `<nav>`, `<header>`, `<footer>` where applicable.

### Forms
- [ ] Every input has a `<label>` linked via `for`/`id` — not placeholder-only.
- [ ] Required fields are marked with both a visual indicator and a text label.
- [ ] Validation error messages are linked to inputs via `aria-describedby`.
- [ ] When form submission fails, focus moves to the first invalid field.
- [ ] Error messages explain what is wrong and how to fix it — not just "This field is invalid."

### Images
- [ ] All `<img>` elements have an `alt` attribute.
- [ ] Informative images have descriptive alt text.
- [ ] Decorative images use `alt=""`.
- [ ] Icon-only buttons have an `aria-label` describing the action.

### Keyboard Navigation
- [ ] Tab through the entire implementation using only the keyboard (no mouse).
- [ ] Every interactive element is reachable by Tab.
- [ ] A visible focus ring is present at every interactive element.
- [ ] `outline: none` or `outline: 0` is **not** used without a visible custom focus replacement.
- [ ] Buttons activate on both `Enter` and `Space`.
- [ ] Modals trap focus internally and return focus to the trigger element on close.
- [ ] The primary task can be completed without using the mouse.

### Color and Contrast
- [ ] All body text meets WCAG AA 4.5:1 contrast ratio against its background.
- [ ] All large text (≥ 18px normal or ≥ 14px bold) meets 3:1 contrast ratio.
- [ ] UI components (buttons, inputs, focus rings) meet 3:1 contrast ratio.
- [ ] The focus indicator meets 3:1 contrast against adjacent colors.

### Motion
- [ ] All animations respect `prefers-reduced-motion: reduce`.
- [ ] Nothing on the page flashes more than 3 times per second.

---

## Section 5 — Responsive Design

Test at each of these viewports. Check the boxes only after testing — not by visual inspection at desktop alone.

| Viewport | Width | Status |
| :--- | :--- | :--- |
| Small phone | 320px | `[ ]` |
| Standard phone | 375px | `[ ]` |
| Large phone | 430px | `[ ]` |
| Tablet portrait | 768px | `[ ]` |
| Small desktop | 1024px | `[ ]` |
| Standard desktop | 1280px | `[ ]` |

For each tested viewport, verify:
- [ ] No horizontal scrollbar is present.
- [ ] No content overflows its container.
- [ ] Text is readable — no truncation without an accessible full-text alternative.
- [ ] Interactive elements do not overlap.
- [ ] The primary action is visible and accessible without scrolling.
- [ ] All touch targets are at least 44×44px.
- [ ] Hover-only interactions have a tap-accessible alternative.

### Navigation Responsiveness
- [ ] Navigation collapses or adapts correctly on mobile (hamburger menu, bottom tab bar, or equivalent).
- [ ] Desktop navigation is not a hamburger menu.
- [ ] Persistent sidebars on desktop collapse or hide on mobile.

---

## Section 6 — Frontend Performance

### Images and Media
- [ ] Images use `srcset` and `sizes` for responsive resolution.
- [ ] Images have `width` and `height` attributes to prevent layout shift (CLS).
- [ ] Images below the fold use `loading="lazy"`.
- [ ] Modern image formats are used (WebP or AVIF) with fallbacks.
- [ ] SVG is used for icons, logos, and illustrations.

### JavaScript and Bundles
- [ ] No full library is imported for a single function (e.g., importing all of Lodash for `_.cloneDeep`).
- [ ] Non-critical components and routes use lazy loading / code splitting.
- [ ] No expensive computation runs in the render path or on every re-render.
- [ ] No duplicate logic exists that is already handled by an existing utility or composable.

### CSS and Animation
- [ ] All transitions and animations use `transform` and `opacity` — not layout-triggering properties (`width`, `height`, `top`, `left`, `margin`).
- [ ] `backdrop-filter` and `filter` effects are not applied to frequently repainted elements.
- [ ] `will-change` is removed after animation completes — it is not applied statically to all animated elements.

### Fonts
- [ ] Only the font weights and styles actually used are loaded.
- [ ] Web fonts use `font-display: swap` to prevent invisible text during load.

---

## Section 7 — Design System Compliance

### Before Adding a New Component — Answer All Four Questions
- [ ] Does an existing component in the project already solve this?
- [ ] Can an existing component be extended rather than replaced?
- [ ] Should this be a shared component (used in multiple places) or a page-specific component?
- [ ] Is the new component placed in the correct directory per project conventions?

### Token Compliance
- [ ] No hardcoded color values in component code.
- [ ] No hardcoded spacing values (`margin: 13px`, `padding: 7px`).
- [ ] No hardcoded font sizes outside the typography token scale.
- [ ] No hardcoded shadow or border radius values outside the token system.

### Component Duplication Check
- [ ] No component was built that already exists in the project's library.
- [ ] No markup was copy-pasted from another component — shared markup is extracted into a reusable component.

---

## Section 8 — Code Maintainability

### Component Structure
- [ ] Each component has a single, clear responsibility (display, form, layout, or feedback — not multiple).
- [ ] Component names are descriptive (describe what the component **is**, not what it looks like).
- [ ] Props are explicitly typed with documented interfaces.
- [ ] Props are not mutated inside a component — events are emitted instead.
- [ ] No `$parent` or `$root` is used to pass data downward.

### Styling
- [ ] No inline styles applied via `style=""` — component-level styles use the token system or scoped styles.
- [ ] No CSS class names used inconsistently (different naming conventions mixed in the same file).
- [ ] No duplicated CSS rules that could be a shared utility or token.

### State Management
- [ ] Component state variables are named clearly (no `flag`, `temp`, `data2`).
- [ ] Loading, empty, error, success, and disabled states are tracked explicitly, not inferred from data shape.
- [ ] No deeply nested reactive objects where shallow alternatives work.

### Readability
- [ ] Another developer can understand the component's purpose and behavior within 30 seconds of opening the file.
- [ ] No commented-out code blocks without an explicit, dated reason for preservation.
- [ ] No TODO comments in production-bound code.
- [ ] No debug output (`console.log`, `debugger`) left in submitted code.

---

## Section 9 — Figma / Design Handoff Compliance

*(Complete this section only when working from a design reference.)*

- [ ] Spacing is implemented from the design token system — not by measuring Figma pixel distances literally.
- [ ] Component types used match the design reference — no substituting a card where a list is designed.
- [ ] All interactive states shown in the design are implemented (hover, active, focus, disabled, loading).
- [ ] All states **not** shown in the design are defined and implemented (empty state, error state, loading state).
- [ ] Responsive behavior is implemented based on design intent — not only the artboard size shown.
- [ ] Icon assets match the project's icon library — not Figma-specific icons that differ from the codebase.

---

## Section 10 — Testing Tools Verification

*(Complete the rows applicable to the project setup.)*

### Accessibility Testing
| Tool | Run | Result |
| :--- | :--- | :--- |
| **axe-core** (browser extension or CI) | `[ ]` | |
| **Lighthouse** accessibility audit | `[ ]` | |
| **Manual keyboard navigation test** | `[ ]` | |
| **Screen reader spot check** (VoiceOver / NVDA) | `[ ]` | |

### Visual and Functional Testing
| Tool | Run | Result |
| :--- | :--- | :--- |
| **Playwright / Cypress** functional test | `[ ]` | |
| **Storybook** component state review | `[ ]` | |
| **Chromatic** visual regression | `[ ]` | |
| **Browser DevTools responsive mode** | `[ ]` | |

> [!NOTE]
> Do not introduce testing tools that are not already in the project without an architectural decision. Use what the project has. The goal is coverage — not tool accumulation.

---

## Final Approval Gate

A UI implementation is complete **only** when all of the following are true:

```
[ ] User goal is immediately clear
[ ] UX flow is logical and complete
[ ] No generic AI-generated patterns without justification
[ ] Design system tokens are used throughout
[ ] Existing components were reused before new ones were created
[ ] Responsive behaviour tested at 320px, 375px, 768px, 1024px, 1280px
[ ] Accessibility verified: semantic HTML, keyboard, contrast, labels
[ ] Loading state is implemented
[ ] Empty state is implemented
[ ] Error state is implemented
[ ] Success state is implemented where applicable
[ ] No hardcoded design values (colors, spacing, fonts) outside tokens
[ ] No duplicate component or logic
[ ] Images are optimized with srcset, alt, width/height
[ ] Animations use transform/opacity only
[ ] No debug output or TODO comments in submitted code
[ ] Another developer can understand the UI without explanation
```

> [!CAUTION]
> A frontend feature is **not complete** because it renders. It is complete when users can understand it, successfully complete their tasks with it, developers can maintain it, it follows the product's design language, and it performs well.

---

## References
- UI/UX Philosophy: [design/00-ui-ux-philosophy.md](../design/00-ui-ux-philosophy.md)
- Design Systems: [design/01-design-systems.md](../design/01-design-systems.md)
- Component Libraries: [design/02-component-libraries.md](../design/02-component-libraries.md)
- Responsive Design: [design/03-responsive-design.md](../design/03-responsive-design.md)
- Accessibility Standard: [design/04-accessibility.md](../design/04-accessibility.md)
- Visual Quality: [design/05-visual-quality.md](../design/05-visual-quality.md)
- UI Review Process: [design/06-ui-review-process.md](../design/06-ui-review-process.md)
- Feature Implementation Checklist: [checklists/01-feature-implementation-checklist.md](01-feature-implementation-checklist.md)
