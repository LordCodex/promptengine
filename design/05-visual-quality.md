---
document_id: design-visual-quality
title: Visual Quality Standard
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
  - design-accessibility
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Visual Quality Standard

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md), the [Design Systems Standard](01-design-systems.md), and the [Accessibility Standard](04-accessibility.md). It defines the visual quality requirements that all implemented interfaces must meet before being considered production-ready.

---

## 1. Visual Quality Is Intentional Design

Visual quality is not about using the latest trends or adding more visual elements. It is about making every visual decision serve a purpose.

A high-quality interface communicates:
- **Hierarchy**: What is most important?
- **Relationship**: What belongs together?
- **State**: What is the current system condition?
- **Affordance**: What can be interacted with?
- **Brand**: What is the personality of this product?

If a visual element does not communicate at least one of the above, it is noise.

---

## 2. Typography Quality

Typography carries the largest proportion of information in most interfaces. Poor typography destroys readability regardless of how good the layout is.

### Font Selection Rules
- Use a typeface designed for screen legibility, not print. System font stacks are a valid and performant choice.
- For web-loaded fonts, use a maximum of 2 typefaces per project: one for body/UI, one for headings/display (if needed).
- Never use decorative or display fonts for body text.
- Load only the weights and styles needed: typically 400, 500, 600, 700. Avoid loading every weight.

**Recommended web-safe starting points:**
```css
/* System font stack — zero load time, native legibility */
font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;

/* Google Fonts — legible, widely tested */
font-family: 'Inter', system-ui, sans-serif;
font-family: 'Plus Jakarta Sans', system-ui, sans-serif;
font-family: 'DM Sans', system-ui, sans-serif;
```

### Readability Rules
- Body text: minimum 16px. Never below 14px for running content.
- Line height for body text: 1.5 to 1.75.
- Paragraph line length: 60–75 characters (approximately 38–45em). Text that spans 90%+ of a widescreen monitor is unreadable.
- Letter spacing: do not increase letter spacing on lowercase body text. Only apply it to uppercase labels or headings.
- Do not set all-caps text below 12px. It becomes unreadable at small sizes.
- Do not justify text. Left-aligned text is easiest to read in left-to-right languages.

### Hierarchy Rules
- A maximum of 3–4 distinct type sizes should be visible on any single screen.
- More than 4 sizes signals a broken hierarchy. Consolidate.
- Size alone is not enough hierarchy. Combine: size + weight + color + spacing.
- The most important text on the screen should be immediately identifiable without scanning.

---

## 3. Color Quality

### Palette Discipline
- Define a palette before building. Do not add colors ad-hoc.
- A minimal production palette:
  - 2–3 neutrals (background, surface, border)
  - 1 primary (brand action color)
  - 1 danger (destructive actions, errors)
  - 1 success (confirmations)
  - 1 warning (caution states)
- Avoid using more than 5 distinct hues in a single product interface.

### Color as Communication
Colors carry meaning. Use them consistently:

| Color Role | Meaning | Never Use It For |
| :--- | :--- | :--- |
| Primary | The main action, active selection | Decorative backgrounds |
| Danger | Destructive actions, errors, alerts | Warnings or neutral states |
| Success | Confirmations, completion | Ongoing active states |
| Warning | Caution, potential issues | Errors or success |
| Neutral | Structure, backgrounds, borders | Communicating state |

### Anti-Patterns
- **Color overload**: Using 6+ distinct colors on a single screen for no semantic reason.
- **Meaningless gradients**: Gradients applied to backgrounds that do not represent motion, depth, or brand.
- **Vibrating color pairs**: High-saturation complementary colors placed adjacent (red/green, blue/orange). They create visual vibration and are inaccessible.
- **Low-contrast decorative text**: Light gray text on white for visual refinement that sacrifices legibility.

---

## 4. Spacing and Alignment Quality

### Alignment Discipline
- Use a consistent alignment axis. On most interfaces this is a left alignment with a consistent left gutter.
- Do not center-align long text or multi-element sections. Center alignment is for short labels, headings, and single elements.
- All related elements align to the same axis. Misaligned elements signal disorganization.
- Use a grid or layout system — do not position elements by eye.

### Spacing Intention
Every spacing decision must be intentional:
- **Close spacing (4–8px)**: Items that are the same type and belong together (list items, form field + its label).
- **Medium spacing (16–24px)**: Groups of related items separated from other groups.
- **Large spacing (32–64px)**: Major sections with distinct content purposes.
- **Section breaks (64–96px)**: Full page sections on marketing or landing pages.

### Common Spacing Failures
- **Random margins**: Elements that have different margins for no reason. Apply the spacing scale uniformly.
- **Padding asymmetry**: Inconsistent vertical/horizontal padding inside cards and panels.
- **Touch proximity**: Interactive elements placed so close that users tap the wrong one.
- **Empty whitespace**: Large sections of empty space between elements that should be grouped.
- **No whitespace**: Everything crammed together with no breathing room.

---

## 5. Component States

Every interactive component must handle every possible state it can enter. Shipping a component without considering all its states is an incomplete implementation.

### Required States by Component Type

| Component | Required States |
| :--- | :--- |
| **Button** | Default, Hover, Active (pressed), Focus, Loading, Disabled |
| **Input** | Default, Focus, Filled, Invalid (error), Disabled, Read-only |
| **Form** | Default, Validation loading, Validation error, Submission success |
| **Data table / List** | Loading skeleton, Empty, Error, Populated |
| **Card** | Default, Hover (if interactive), Loading skeleton, Error |
| **Modal** | Open, Closing animation, Error state (if action fails) |
| **Dropdown** | Closed, Open, Selected, Empty (no options), Loading (async) |
| **Page** | Loading, Error, Empty, Content |

### State Design Rules

#### Loading State
- Use skeleton screens (content-shaped placeholders) instead of generic spinners for content that takes more than 200ms to load.
- Spinners are appropriate for actions (button submits, saving) — not for page content loading.
- Never leave a blank area with no loading indicator.

#### Empty State
- Explain why there is nothing here.
- Tell the user what they can do about it.
- If an action can remedy the empty state (e.g., "No projects yet — Create your first project"), provide the action inline.
- Never display a blank table, empty list, or empty card grid without explanation.

#### Error State
- Use plain language. "Something went wrong. Please try again." is acceptable. A stack trace is not.
- Provide a recovery action when possible (retry button, link to support).
- Distinguish between user errors (fix their input), system errors (retry or contact support), and network errors (check connection).

#### Disabled State
- Reduce opacity or use the disabled color token. Never hide disabled elements.
- If a user might wonder why something is disabled, add a tooltip or explanatory label.
- Disabled elements must still be visually identifiable as the component type they represent.

---

## 6. Iconography Quality

Icons communicate at a glance. Poor icon usage creates confusion.

### Icon Rules
- Every icon used alone (without a text label) must have an accessible `aria-label` or tooltip.
- Never use an icon that is ambiguous for the action it represents. When in doubt, add a text label.
- Maintain consistent icon size within a context. Do not mix 16px and 24px icons on the same line.
- Use filled variants for active/selected states and outline variants for default states, or vice versa — but never mix the convention within a product.
- Do not use icons decoratively in body text.

### Icon Sizing
| Context | Size |
| :--- | :--- |
| Inline with body text | 16px |
| Buttons, form elements | 20px |
| Navigation, tabs | 20–24px |
| Empty state illustration | 48–64px |
| Hero / large display | SVG, scalable |

---

## 7. Animation and Motion Quality

Animation must earn its place. Every motion should communicate something: state change, relationship, direction, or feedback. Motion must support the product — not become the product.

### When Animation Adds Value
Use subtle transitions and animations for:
- Button feedback and hover/active states.
- Expanding or collapsing sections (accordions).
- Skeletons and progress bar loaders.
- Modal dialog appearance and drawer slide-outs.
- Toast notifications entering/leaving.
- Drag-and-drop actions, list item sorting, and filter changes.
- Page transitions that help preserve context.

### When Animation Adds Noise
Avoid animating:
- Static reading text.
- Every individual card container in a grid.
- Every icon inside navigation lines.
- Page elements loading sequentially (jarring staggered fades).
- Heavy decorative parallax layers or screen-filling transitions.
Too much motion reduces usability and causes cognitive overload.

### Micro-interactions & Immediate Feedback
State changes should trigger immediate visual feedback:
- **Button presses / toggles**: Instant scale-down or background shift to confirm activation.
- **Checkbox selection**: Crisp, instant checked indicators.
- **Successful saves**: Simple transient checkmarks or toast animations.
- **File upload**: Real-time progress bars, not looping spinners.
- **Copy-to-clipboard**: Instant confirmation label change or tooltip.

### Component-Specific Motion Behaviors

#### Page Transitions
- Keep transitions fast, subtle, and consistent.
- Use simple sliding or opacity fades. Avoid dramatic zooms, structural rotations, or heavy parallax shifts.

#### Modals and Drawers
- Draw attention to the modal focus area using a quick slide-up or fade-in transition.
- Focus must move inside the modal immediately on open and return to the trigger element on close.

#### Dropdowns and Menus
- Dropdowns must open instantly to keep interactions crisp.
- Motion must never interfere with pointer clicks or keyboard selections.

#### Tables and Data Lists
- **Do not animate entire tables or large datasets.**
- Limit motion to subtle, isolated row insertions, column sorting indicators, and inline cell status updates.

#### Form Inputs
- Provide immediate, subtle transitions on input focus (e.g. border color changes).
- Animate validation state changes (error shake or fade-in validation labels) dynamically.

### Mobile Device Considerations
- Design for touch target feedback (visual active states on touch).
- Match gesture expectations (swipe to dismiss drawers, pull-to-refresh).
- Keep animations simple and performant on lower-powered mobile processors.

### Easing & Duration Guidelines
Use consistent duration and easing rules across the application:

| Motion Type | Duration | Easing |
| :--- | :--- | :--- |
| **Micro-interactions** (hover, toggle, press) | 80–150ms | `ease-out` (decelerates to finish) |
| **Component transitions** (dropdowns, accordions) | 150–200ms | `ease-in-out` (smooth enter/exit) |
| **Overlay transitions** (modals, drawers) | 200–250ms | `ease-out` on enter, `ease-in` on exit |
| **Page transitions** | 200–300ms | `ease-in-out` |
| **Loading loops** (skeletons, progress bars) | 1000–1500ms | `linear` |

*Never use `linear` easing for interactive transitions.*

### Framework Motion Guidance
- **Vue / Nuxt**: Prefer native `<Transition>` and `<TransitionGroup>` wrappers. Avoid importing external motion packages for standard transitions.
- **React / Next.js**: Prefer standard CSS transitions. Use libraries like Framer Motion only when complex SVG path morphing or physics-based gestures are strictly required.
- **Flutter**: Use built-in animation widgets (`AnimatedContainer`, `Hero`) and ThemeData-scoped platform motion.

### AI Motion Review Queries
Before deploying any animation or transition, answer these four questions:
1. *Does this animation communicate a useful state change or user action?*
2. *Does it improve the user's understanding of the layout?*
3. *Does the duration slow down the user's workflow?*
4. *Can it be removed without harming usability?*

If the animation slows down the user or serves only as decoration, **remove it**.

---

## 8. UI Performance (Visual)

Visual performance is part of quality. An interface that is beautiful but slow is a poor interface.

### Rules
- Do not use heavy CSS filter effects (`backdrop-filter`, `blur`) on frequently repainted elements.
- Do not animate properties that trigger layout recalculation (`width`, `height`, `top`, `left`, `margin`). Animate `transform` and `opacity` instead.
- Do not load web fonts that are not used within the first 3 seconds of the page.
- Use `will-change: transform` only on elements that are actively animating — remove it after animation completes.
- Images that are not visible on initial load must be lazy-loaded.

---

## Review Checklist

### Typography
- [ ] Is body text at least 16px with a line height of 1.5+?
- [ ] Is maximum paragraph width limited to approximately 75 characters?
- [ ] Is the type hierarchy visible with no more than 3–4 distinct sizes on a screen?
- [ ] Is text left-aligned (not justified)?

### Color
- [ ] Is the palette defined with semantic role tokens (not raw values)?
- [ ] Is color never the only way to communicate state or information?
- [ ] Does all text meet WCAG AA contrast ratios?

### Spacing
- [ ] Does all spacing use the project's spacing token scale?
- [ ] Are related elements visually grouped by proximity?
- [ ] Are there no random or inconsistent margins?

### Component States
- [ ] Do all interactive components handle: default, hover, focus, loading, error, disabled?
- [ ] Are empty states informative with a clear action?
- [ ] Are loading states visible (skeleton or spinner)?
- [ ] Do error states explain what happened and what to do?

### Animation
- [ ] **Motion communicates purpose**: The animation explains a state change or guides user attention.
- [ ] **Hover feedback exists**: All interactive components display clear visual hover states.
- [ ] **Loading states exist**: Async activities display skeleton screens or progress indicators.
- [ ] **Success states exist**: Form completions and actions display visual success indicator confirmation.
- [ ] **Error states exist**: Validation failures or network timeouts trigger error warning animations (e.g. shake).
- [ ] **Reduced motion supported**: Motion transitions respect user `prefers-reduced-motion` preferences.
- [ ] **Mobile interaction considered**: Touch feedback active; animations are performant on low-powered devices.
- [ ] **Performance preserved**: Animation targets are limited to `transform` and `opacity` to prevent layout recalculations.
- [ ] **Consistent timing**: Durations and easing functions adhere to the project's consistent scales.
- [ ] **Animations are subtle**: The visual presence of transitions is low-profile and does not delay user workflows.

---

## References
- Design Systems (tokens): [design/01-design-systems.md](01-design-systems.md)
- Accessibility (contrast, motion): [design/04-accessibility.md](04-accessibility.md)
- UI Review Process: [design/06-ui-review-process.md](06-ui-review-process.md)
