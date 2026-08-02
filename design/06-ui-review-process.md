---
document_id: design-ui-review-process
title: UI Review Process
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
  - design-accessibility
  - design-visual-quality
  - design-responsive-design
audience: [human, agent]
last_reviewed: 2026-08-01
---

# UI Review Process

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md) and synthesizes the review requirements from the [Design Systems Standard](01-design-systems.md), the [Accessibility Standard](04-accessibility.md), the [Visual Quality Standard](05-visual-quality.md), and the [Responsive Design Standard](03-responsive-design.md). It defines the self-review and peer-review process for all frontend implementations.

---

## 1. Review Philosophy

A UI review is not a style debate. It is a structured quality gate that verifies the implementation meets five requirements:

1. **UX**: It solves the user's problem correctly.
2. **Design**: It follows the project's visual system.
3. **Accessibility**: It works for all users.
4. **Responsiveness**: It functions at all viewport sizes.
5. **Maintainability**: It can be extended without rewriting.

Any implementation that fails any of these five is incomplete, regardless of how it looks in a single browser at a single viewport.

---

## 2. AI Agent Self-Review Protocol

Before submitting any frontend output, an AI agent must run this self-review sequence. Do not skip layers. Document failures explicitly.

```text
Layer 1: UX Review
     ↓
Layer 2: Design System Review  
     ↓
Layer 3: Accessibility Review
     ↓
Layer 4: Responsive Review
     ↓
Layer 5: Performance Review
     ↓
Layer 6: Maintainability Review
```

---

## 3. Layer 1 — UX Review

**Goal**: Does the implementation solve the user's problem efficiently?

| Question | Pass Condition |
| :--- | :--- |
| Does this solve the user's stated task? | Implementation addresses the core requirement, not a literal interpretation of it |
| Is the user flow complete and logical? | No step requires the user to guess the next action |
| Are unnecessary steps eliminated? | No extra clicks, forms, or confirmations without a security/safety reason |
| Is the most important information immediately visible? | User does not need to scroll to find primary content on load |
| Are all states covered? | Loading, error, empty, and success are visible in the implementation |
| Is destructive action friction present? | Delete, cancel, and irreversible actions have confirmation |

**Failure actions**: If the implementation addresses the wrong problem, stop. Do not polish incorrect UX.

---

## 4. Layer 2 — Design System Review

**Goal**: Is the implementation consistent with the project's visual system?

| Question | Pass Condition |
| :--- | :--- |
| Are all colors from the design token system? | No hardcoded hex values or raw `rgb()` values |
| Are all spacing values from the spacing scale? | No arbitrary pixel values (`margin: 13px`, `padding: 7px`) |
| Is the typography from the type scale? | No arbitrary `font-size` values; heading levels are semantic |
| Are existing components reused? | No parallel implementation of a component that already exists |
| Does it match the project's visual language? | Not a different visual style introduced to this screen |
| Is one icon library used consistently? | No mixing of icon styles or sources |

**Common failures to check actively:**
- Inline styles applied instead of design tokens.
- New colors introduced without adding them to the token system.
- Spacing that "looks right" but uses arbitrary values.
- New components built without checking for existing equivalents.

---

## 5. Layer 3 — Accessibility Review

**Goal**: Can all users use this interface?

Run this in order:

### 5a. Structural Check (no tools required)
- [ ] Is the correct HTML element used for every interactive component? (Button for buttons, `<a>` for links)
- [ ] Is there exactly one `<h1>`? Are heading levels sequential?
- [ ] Does every image have an `alt` attribute?
- [ ] Does every input have a linked `<label>`?
- [ ] Are error messages linked to inputs via `aria-describedby`?

### 5b. Keyboard Test (keyboard only — no mouse)
- Tab through the entire implementation.
- [ ] Is every interactive element reachable?
- [ ] Is the focus indicator visible at every stop?
- [ ] Does every button activate on Enter and Space?
- [ ] Do modals trap focus and return focus on close?
- [ ] Can the user complete the primary task without using the mouse?

### 5c. Contrast Check (automated tool)
- [ ] Does all body text meet 4.5:1 contrast ratio?
- [ ] Do buttons, inputs, and icons meet 3:1 against adjacent backgrounds?
- [ ] Does the focus ring meet 3:1 against adjacent colors?

### 5d. Screen Reader Spot Check (one screen reader)
- Navigate through primary interactive elements.
- [ ] Are button and link purposes announced (not just "button")?
- [ ] Are form errors announced when validation fails?
- [ ] Are dynamic updates (async content, status messages) announced via `aria-live`?

**Pass threshold**: All 5a structural checks pass. All 5b keyboard interactions work. Contrast ratios verified. At least one screen reader test completed.

---

## 6. Layer 4 — Responsive Review

**Goal**: Does the implementation function at all relevant viewports?

Test at these viewports. Do not only check desktop.

| Viewport | Width | Test Requirement |
| :--- | :--- | :--- |
| Small phone | 320px | Content visible, no horizontal overflow |
| Standard phone | 375px | Layout correct, touch targets accessible |
| Large phone | 430px | Layout correct |
| Tablet portrait | 768px | Layout transitions correctly |
| Small desktop | 1024px | Full layout present |
| Standard desktop | 1280px | Primary design viewport |

**Failure indicators:**
- Horizontal scrollbar appears on any viewport below 1280px.
- Content overflows its container.
- Text truncated without an accessible way to read the full content.
- Interactive elements overlap.
- Primary action is hidden or inaccessible at mobile viewports.
- Touch targets smaller than 44px.

---

## 7. Layer 5 — Performance Review

**Goal**: Does the implementation load and respond without unnecessary cost?

| Check | Pass Condition |
| :--- | :--- |
| Image sizes | Images use `srcset`; no desktop-resolution images served to mobile |
| Bundle size | No full library imported for a single utility function |
| Animation performance | Transitions use `transform`/`opacity` only; no layout-triggering animation |
| Lazy loading | Below-fold images and non-critical components use lazy loading |
| Heavy operations | No expensive computation in the render path |
| External fonts | Only necessary weights loaded; `font-display: swap` applied |

---

## 8. Layer 6 — Maintainability Review

**Goal**: Can another developer extend this in 6 months without difficulty?

| Check | Pass Condition |
| :--- | :--- |
| Component placement | File is in the correct directory per project conventions |
| Props typing | All props are typed with explicit interfaces |
| No duplication | No component or logic that already exists was recreated |
| No magic values | No hardcoded colors, sizes, or strings that should be design tokens or constants |
| State clarity | Component state is readable — no cryptic flag variables |
| No excessive abstraction | Component does not use 3 layers of abstraction for a simple display task |
| Reusability | The component is generic enough to be used in other contexts (if it is in a shared directory) |

---

## 9. Peer Review Criteria

Human reviewers checking a frontend implementation must verify the same 6 layers and additionally check:

- **Does the component match the provided design reference?** (if one exists)
- **Was the Design System First Protocol followed?** Confirm no parallel component implementations were created.
- **Was an existing library component available that was not used?** Justify custom implementations.
- **Is the implementation aligned with the project's established code style?**

---

## 10. Common UI Review Failures

These are the most frequently found failures during UI review:

### Design Failures
- Arbitrary pixel values instead of spacing tokens.
- Hardcoded color hex values instead of design tokens.
- New components built without checking for existing equivalents.
- Type sizes that do not follow the scale.

### Accessibility Failures
- `outline: none` with no focus ring replacement.
- Labels missing on form inputs.
- Color as the only state indicator.
- Modals that do not trap focus.
- Images missing `alt` text.

### Responsive Failures
- Fixed pixel widths that break on small viewports.
- Desktop-only layouts untested at 375px.
- Touch targets smaller than 44px.
- Horizontal scroll on mobile.

### UX Failures
- Missing loading state on async operations.
- Missing empty state on data lists.
- Generic or missing error messages.
- Destructive actions with no confirmation.

---

## 10.5. Modern UI Evaluation & Polish

Do not approve a user interface simply because it renders, contains the correct component types, or meets the literal text of the requirements. Assess its overall visual execution against these parameters:

### A. Visual Hierarchy and Priority
- **Primary Actions**: The user must immediately identify the most important action. Avoid placing competing primary buttons next to each other.
- **Content Weight**: Crucial info must have better placement, size, weight, or spacing contrast. Do not let secondary or decorative meta-information shout louder than primary data.

### B. Human Design vs. AI-Generated Patterns
- **Human-centered UI** feels intentional, calm, clear, and predictable.
- **Generated UI** feels busy, over-designed, template-driven, and generic.
- **Strictly reject**: Random background gradients, excessive card components, dashboard widgets added without a product analytics task, and animations that do not communicate state changes.

### C. Application-Type Optimization
- **Data-Dense/Admin Apps**: Prioritize efficiency, quick access, scanning speed, columns alignment, and keyboard operations. Avoid wrapping simple tables in nested cards or adding large empty margins.
- **Customer-Facing/Consumer Apps**: Prioritize trust, security indicators, visible pricing, and explicit validation confirmations (e.g. secure transaction checkmarks in payment checkout screens).

### D. Design Reference Usage
When seeking inspiration from platforms like the Figma Community, Behance, Dribbble, or Awwwards:
- Analyze layout structures, user flows, and interaction patterns.
- **Never copy design assets or styles blindly.** Adapt patterns to fit the current project's design tokens and guidelines.

---

## 10.6. UI Review Scorecard

Before completing UI work, verify the implementation passes this review scorecard:
- [ ] **User goal is obvious**: A new user instantly understands the screen's main purpose.
- [ ] **Primary action is clear**: No competing primary actions; the next logical step is visually dominant.
- [ ] **Layout hierarchy is strong**: Visual layout guides the eye naturally to high-priority information.
- [ ] **Components are consistent**: Similar components look and behave consistently.
- [ ] **Typography is readable**: Hierarchy uses a structured scale; line lengths are bounded.
- [ ] **Colors have purpose**: Colors carry semantic meaning (success, warning, error) rather than raw decoration.
- [ ] **Responsive behavior works**: Content priority adjusts correctly without shrinking layout ratios.
- [ ] **Accessibility considered**: Operable by keyboard; color contrast and focus rings meet standard checks.
- [ ] **Empty/error states exist**: UI handles loading, data absence, and network errors gracefully.
- [ ] **Interface feels intentional**: The design is clean, readable, and free of generic AI-generated templates.

---

## 11. Definition of Done for Frontend Work

A frontend implementation is complete when:

- [ ] All 6 self-review layers pass.
- [ ] The implementation is tested at the minimum viewport set (320px, 375px, 768px, 1024px, 1280px).
- [ ] Automated accessibility audit (axe or Lighthouse) shows zero violations.
- [ ] Keyboard-only navigation works for all primary interactions.
- [ ] All component states are implemented: default, loading, empty, error, success (where applicable).
- [ ] No hardcoded design values (colors, spacing, fonts) outside the design token system.
- [ ] No duplicate component or logic that already exists in the project.
- [ ] Performance: no unnecessary imports, no layout-triggering animations.

---

## References
- UI/UX Philosophy: [design/00-ui-ux-philosophy.md](00-ui-ux-philosophy.md)
- Design Systems: [design/01-design-systems.md](01-design-systems.md)
- Accessibility Standard: [design/04-accessibility.md](04-accessibility.md)
- Visual Quality: [design/05-visual-quality.md](05-visual-quality.md)
- Responsive Design: [design/03-responsive-design.md](03-responsive-design.md)
- Code Review Standard: [core/18-code-review-engineering-standard.md](../core/18-code-review-engineering-standard.md)
