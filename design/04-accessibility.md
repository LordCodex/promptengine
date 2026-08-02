---
document_id: design-accessibility
title: Accessibility Engineering Standard
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Accessibility Engineering Standard

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md), the [Design Systems Standard](01-design-systems.md), and the [Universal Coding Standards](../core/05-universal-coding-standards.md). It defines the accessibility requirements for all UI work across the playbook.

---

## 1. Accessibility Philosophy

Accessibility is not an optional enhancement. It is a baseline requirement.

An inaccessible interface excludes users. Depending on the context:
- **Keyboard users** include people with motor disabilities, power users, and keyboard-navigation preferences.
- **Screen reader users** include people who are blind or have severe visual impairments.
- **Users with cognitive disabilities** benefit from clear language, predictable behavior, and reduced motion.
- **Users with color blindness** rely on shape and text, not color alone, to interpret information.

> [!IMPORTANT]
> Accessible design is better design for everyone. Focus states improve usability for all keyboard users. Clear error messages help all users. Sufficient color contrast helps users in bright sunlight.

**Target**: WCAG 2.2 Level AA compliance as the minimum standard for all production interfaces.

---

## 1.5. WCAG Core Principles

All interfaces must satisfy the four core WCAG principles (POUR):

### A. Perceivable
Users must be able to perceive the information presented.
- **Text alternatives**: Provide descriptive text/alt text for non-text items.
- **Captions**: Provide descriptions for video and audio content.
- **Visual hierarchy**: Ensure headings and spacing make layout scanning natural.
- **Contrast**: Enforce sufficient visual contrast between text/components and background layers.

### B. Operable
Users must be able to operate the interface.
- **Keyboard navigation**: The entire application must be operable without a mouse.
- **Focus management**: Focus outline, focus indicators, and focus boundaries must be clearly visible and correctly ordered.
- **Touch targets**: Provide target spacing and dimensions (minimum 44x44px) that prevent accidental clicks on small screens.
- **No dangerous interactions**: Avoid triggering actions on hover or timing constraints that trap users.

### C. Understandable
Users must be able to understand the information and operation of the user interface.
- **Explicit actions**: Labels and controls must describe what they trigger or where they lead.
- **Clear errors**: Help text and input errors must explain what is wrong and how to fix it.
- **Predictable behavior**: UI layouts and action responses must align with platform conventions.

### D. Robust
Content must be robust enough that it can be interpreted reliably by a wide variety of user agents, including assistive technologies.
- **Compatibility**: Ensure semantic HTML parses cleanly across different browsers, devices, and screen readers without structural syntax failures.

---

## 2. Native HTML Semantics First

The most effective accessibility technique is using the correct HTML element for each purpose. Native HTML elements carry built-in accessibility semantics, keyboard behavior, and ARIA roles — for free.

### Element Selection Rules

| UI Requirement | Correct Element | Never Use |
| :--- | :--- | :--- |
| Clickable action | `<button>` | `<div onclick>`, `<span onclick>` |
| Navigation link | `<a href>` | `<div onclick>` for navigation |
| Form input | `<input>`, `<textarea>`, `<select>` | `<div contenteditable>` |
| Page heading | `<h1>`–`<h6>` | Styled `<div>` or `<p>` with large font |
| Navigation landmark | `<nav>` | `<div class="nav">` |
| Main content | `<main>` | `<div id="main">` |
| Article | `<article>` | Unlabeled `<div>` |
| Tabular data | `<table>` with `<th>` | `<div>` grid layout |
| Form grouping | `<fieldset>` + `<legend>` | Unlabeled `<div>` groups |

### The ARIA Rule
Do not add ARIA attributes to elements that already communicate their role semantically.

```html
<!-- Wrong: redundant ARIA on native element -->
<button role="button" aria-pressed="false">Save</button>

<!-- Correct: native button handles role automatically -->
<button>Save</button>
```

> [!NOTE]
> The first rule of ARIA is: **do not use ARIA if native HTML can do the job.** ARIA supplements; it does not replace semantic HTML.

---

## 2.5. Buttons and Links

Always choose the correct element based on the intent of the user interaction:

- **Button (`<button>`)**: Use for actions that trigger changes or workflows on the current page.
  - *Examples*: Submitting a form, opening a modal dialog, toggling a dropdown, deleting an item, saving a draft.
- **Link (`<a>`)**: Use for navigation that changes the user's location (changing the URL).
  - *Examples*: Navigating to a different route/page, visiting an external resource, skipping to a content anchor.

### Rules
- **Never use clickable `<div>` or `<span>` elements as buttons.** Doing so strips away native keyboard operability (`Enter`, `Space`) and screen reader role announcements unless heavily overridden (which adds unnecessary complexity).
- **Never use links for state-changing actions.** If the interaction modifies data, triggers a popup, or executes logic without navigation, it must be a button.

---

## 3. Keyboard Navigation

Every interactive element must be operable by keyboard without a mouse.

### Required Keyboard Behaviors

| Component | Required Keyboard Behavior |
| :--- | :--- |
| **Button** | `Enter` and `Space` activate it |
| **Link** | `Enter` activates it |
| **Checkbox** | `Space` toggles it |
| **Radio group** | Arrow keys move between options |
| **Dropdown/Select** | `Enter`/`Space` opens; arrows navigate; `Escape` closes |
| **Modal/Dialog** | Focus moves inside on open; `Escape` closes; focus returns to trigger on close |
| **Tab panel** | Arrow keys navigate tabs; `Enter`/`Space` activates |
| **Combobox** | Arrow keys navigate suggestions; `Enter` selects; `Escape` closes |
| **Date picker** | Arrow keys navigate dates; `Enter` selects; `Escape` closes |

### Focus Management Rules
- Do not remove the focus ring without replacing it with a clearly visible custom focus indicator.
- Never use `outline: none` or `outline: 0` without a custom replacement.
- When a modal opens, move focus to the first interactive element inside the modal.
- When a modal closes, return focus to the element that triggered it.
- Focus trapping (keeping focus inside a modal while it is open) is required for modal dialogs.

```css
/* Correct: visible focus ring using project focus token */
:focus-visible {
  outline: 2px solid var(--color-focus-ring);
  outline-offset: 2px;
}

/* Wrong: removing focus ring entirely */
:focus {
  outline: none;
}
```

---

## 4. Color and Contrast

### Contrast Requirements (WCAG AA)

| Content Type | Minimum Ratio |
| :--- | :--- |
| Body text (< 18px normal or < 14px bold) | **4.5:1** |
| Large text (≥ 18px normal or ≥ 14px bold) | **3:1** |
| UI components (buttons, inputs, icons) | **3:1** |
| Focus indicators | **3:1** against adjacent colors |

### Color Usage Rules
- Never use color as the **only** way to convey information. Always pair color with a shape, icon, text label, or pattern.
  - Wrong: "Required fields are shown in red."
  - Correct: "Required fields are marked with a red asterisk (*) and the label 'Required'."
- Error states must use text, not just a red border.
- Success states must use text or an icon, not just a green color.
- Design systems must define color tokens that automatically meet contrast requirements in both light and dark mode.

### Color Blindness Considerations
- Red/green color blindness (deuteranopia) is the most common form. Never rely solely on red vs. green to convey state.
- Use shape, pattern, or icon paired with color for status indicators.
- Test designs with a color blindness simulation tool (e.g., Colorblindly browser extension, Figma Color Blind plugin).

---

## 5. Forms and Inputs

Forms are among the most interaction-heavy and accessibility-critical parts of any interface.

### Label Rules
- Every input must have a `<label>` that references it via `for`/`id` linkage.
- Never use `placeholder` as a replacement for a label. Placeholders disappear on input and fail contrast requirements.
- If a label must be visually hidden, use a visually-hidden CSS class — do not use `display: none` or `visibility: hidden` (which hides it from screen readers too).

```html
<!-- Correct: linked label and input -->
<label for="email">Email address</label>
<input type="email" id="email" name="email" required>

<!-- Wrong: placeholder-only "label" -->
<input type="email" placeholder="Email address">
```

### Validation and Error Messages
- Validation errors must be announced to screen readers. Use `aria-describedby` to link error messages to inputs.
- Place error messages adjacent to the invalid field, not only at the top of the form.
- When validation fails on submit, move focus to the first invalid field.
- Error messages must explain what went wrong and how to fix it — not just "This field is invalid."

```html
<!-- Correct: input linked to its error message -->
<label for="email">Email address</label>
<input
  type="email"
  id="email"
  name="email"
  aria-describedby="email-error"
  aria-invalid="true"
>
<p id="email-error" role="alert">
  Enter a valid email address (e.g., user@example.com).
</p>
```

### Required Fields
- Mark required fields with both a visual indicator and a text label (not just an asterisk).
- Include a legend: "Fields marked with * are required."
- Use the `required` attribute on inputs for native browser validation support.

---

## 6. Images and Media

### Alt Text Rules
- Every `<img>` must have an `alt` attribute.
- Decorative images that convey no information: use `alt=""` (empty string) so screen readers ignore them.
- Informative images: write descriptive alt text that conveys the same information the image conveys.
- Images of text: the alt text must contain the same text shown in the image.
- Complex images (charts, graphs): provide a short alt text and a full text description nearby or linked.

```html
<!-- Informative image -->
<img src="error-diagram.svg" alt="Diagram showing the request flow when a 500 error occurs">

<!-- Decorative image -->
<img src="background-texture.png" alt="">

<!-- Icon button — the alt text describes the action -->
<button>
  <img src="icons/save.svg" alt="Save document">
</button>
```

### Video and Audio
- Provide captions for all video content.
- Provide transcripts for audio content.
- Video players must be keyboard-operable.
- Do not autoplay video or audio with sound without user consent.

---

## 7. Motion and Animation

Some users experience nausea, dizziness, or seizures triggered by motion and animation.

### Reduced Motion Rules
- All animations and transitions must respect the `prefers-reduced-motion` media query.
- This applies to: page transitions, hover animations, loading spinners, scroll animations, parallax effects.

```css
/* Apply animations by default */
.card {
  transition: transform 200ms ease, box-shadow 200ms ease;
}

/* Disable or reduce animations for users who prefer it */
@media (prefers-reduced-motion: reduce) {
  .card {
    transition: none;
  }
}
```

### Animation Content Rules
- Nothing on the page should flash more than 3 times per second (seizure prevention).
- Auto-playing content must be stoppable.
- Infinite loading spinners are acceptable; large, looping background videos require a pause control.

---

## 8. ARIA Usage Guide

ARIA (Accessible Rich Internet Applications) supplements native HTML semantics for custom components.

### When to Use ARIA
- To describe dynamic content regions that update without a page reload: `aria-live`.
- To describe complex custom widget roles: `role="tablist"`, `role="tree"`, `role="dialog"`.
- To communicate relationship between elements that HTML cannot express: `aria-describedby`, `aria-controls`, `aria-owns`.
- To communicate state not expressible in HTML: `aria-expanded`, `aria-selected`, `aria-invalid`.

### When NOT to Use ARIA
- On native HTML elements that already convey the role: `<button>`, `<input>`, `<a>`, `<nav>`.
- To add `role="button"` to a `<div>` — use a `<button>` element instead.
- As a substitute for correct HTML structure.

### Common ARIA Patterns

```html
<!-- Modal dialog -->
<div
  role="dialog"
  aria-modal="true"
  aria-labelledby="modal-title"
  aria-describedby="modal-description"
>
  <h2 id="modal-title">Confirm deletion</h2>
  <p id="modal-description">This action cannot be undone.</p>
</div>

<!-- Live region for async status -->
<div aria-live="polite" aria-atomic="true">
  <!-- Screen readers announce when this content changes -->
  <p>{{ statusMessage }}</p>
</div>

<!-- Icon-only button -->
<button aria-label="Close dialog">
  <svg aria-hidden="true"><!-- icon svg --></svg>
</button>
```

---

## 9. Accessibility Testing Protocol

Testing must be multi-layered. Automated tools catch approximately 30–40% of accessibility issues. Manual testing is required.

### Layer 1: Automated Testing
- Run **axe-core** (browser extension or CI integration) on every page.
- Run **Lighthouse** accessibility audit in Chrome DevTools.
- Fix all automated tool violations before manual testing.

### Layer 2: Keyboard-Only Testing
- Navigate the entire page using only Tab, Shift+Tab, Enter, Space, and arrow keys.
- Confirm every interactive element is reachable and operable.
- Confirm focus is always visible.
- Confirm modals trap focus correctly.
- Confirm focus returns correctly after modal close.

### Layer 3: Screen Reader Testing
Test with at minimum one screen reader + browser combination:
- **NVDA + Firefox** (Windows, free)
- **JAWS + Chrome** (Windows, most enterprise usage)
- **VoiceOver + Safari** (macOS/iOS, built-in)
- **TalkBack + Chrome** (Android, built-in)

Verify that interactive elements announce their role, name, and state correctly. Verify that dynamic updates (form errors, async content) are announced.

### Layer 4: Color and Contrast
- Test all text/background combinations with a contrast checker (e.g., WebAIM Contrast Checker).
- Simulate color blindness with Colorblindly or equivalent tool.

---

## 10. Accessibility Framework Guidance

Follow these framework-specific patterns to ensure accessibility constraints are met natively:

### 10a. Vue / Nuxt
- **Template Semantics**: Enforce standard HTML elements inside Vue templates rather than wrapping custom tags without semantic mapping.
- **Accessible Component Patterns**: Use proven libraries (e.g. Headless UI Vue) to provide accessibility wrappers for tabs, lists, and dropdowns.
- **Keyboard Event Handling**: Utilize Vue key modifiers (`@keydown.enter`, `@keydown.space`, `@keydown.esc`) to configure keyboard support on interactive inputs.

### 10b. React / Next.js
- **Semantic JSX**: Use JSX tags that compile to valid semantic HTML elements (`<main>`, `<section>`).
- **Accessible Component Libraries**: Prefer Radix UI or React Aria primitives when building custom component designs to ensure ARIA role management is solved correctly.
- **Proper ARIA Usage**: Ensure ARIA attributes map correctly to dynamic React state triggers (e.g. `aria-expanded={isOpen}`).

### 10c. Flutter
- **Semantics Widgets**: Wrap custom visual configurations in `Semantics` widgets to supply names, values, and hints to screen readers.
- **Screen Reader Support**: Test application accessibility with TalkBack (Android) and VoiceOver (iOS).
- **Touch Target Sizes**: Enforce minimum touch sizes on all GestureDetector widgets to meet spacing rules.

---

## Review Checklist

Before completing UI work, verify against this accessibility checklist:
- [ ] **Works without a mouse**: The entire primary task and all controls are operable using keyboard navigation alone.
- [ ] **Semantic HTML used**: Correct markup elements (`<button>`, `<nav>`, `<main>`, `<header>`, `<footer>`, `<form>`, `<label>`, `<table>`) are used instead of custom styled `<div>` / `<span>` elements.
- [ ] **Keyboard navigation works**: Tab, Enter, Space, and Escape behaviors match platform interaction standards.
- [ ] **Focus states exist**: A visible custom focus ring is present on every interactive element; focus trapping is active for modals.
- [ ] **Forms have labels**: Every input has an explicitly linked visible `<label>` (no placeholder-only configurations).
- [ ] **Errors are understandable**: Error notifications explain what is wrong and how to fix it; error states are announced to assistive technologies.
- [ ] **Color is not the only indicator**: Errors, status updates, and links use text, icons, or visual markers in addition to color variations.
- [ ] **Screen readers receive useful information**: Images have descriptive `alt` tags (or `alt=""` for decorative assets); dynamic sections use ARIA live regions.
- [ ] **Components handle accessibility internally**: Reusable UI components include built-in keyboard, focus, and state aria roles out of the box.
- [ ] **Reduced motion considered**: All transitions and animation effects respect the `prefers-reduced-motion` user query.
- [ ] **Contrast considered**: Body text meets the WCAG 4.5:1 ratio; interactive components and large headings meet the 3:1 ratio.

---

## References
- Responsive Design (touch targets): [design/03-responsive-design.md](03-responsive-design.md)
- Visual Quality (color contrast): [design/05-visual-quality.md](05-visual-quality.md)
- WCAG 2.2: [https://www.w3.org/WAI/WCAG22/quickref/](https://www.w3.org/WAI/WCAG22/quickref/)
- ARIA Authoring Practices: [https://www.w3.org/WAI/ARIA/apg/](https://www.w3.org/WAI/ARIA/apg/)
- axe-core: [https://github.com/dequelabs/axe-core](https://github.com/dequelabs/axe-core)
