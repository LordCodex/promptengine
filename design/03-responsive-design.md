---
document_id: design-responsive-design
title: Responsive Design Standard
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Responsive Design Standard

## Inheritance
This document inherits from the [UI/UX Philosophy](00-ui-ux-philosophy.md) and the [Design Systems Standard](01-design-systems.md). It defines the standards for building interfaces that work correctly across all viewport sizes, input methods, and device capabilities.

---

## 1. Responsive Design Philosophy

Responsive design is not "make it fit on a phone". It is a fundamental approach to layout that acknowledges that users access interfaces from an unpredictable range of devices, screen sizes, network conditions, and contexts.

**The design for a small screen is not a reduced version of the desktop design. It is a different priority ordering of the same content.**

---

## 2. Viewport Strategy

### Mobile-First by Default
Write CSS from the smallest viewport outward. This means:

- Base styles apply to all viewports (mobile).
- Media queries add complexity for larger screens, they do not strip it away.

```css
/* Mobile first — base styles apply everywhere */
.card-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-4);
}

/* Tablet and above */
@media (min-width: 640px) {
  .card-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* Desktop and above */
@media (min-width: 1024px) {
  .card-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

**Exception**: If the project is explicitly a desktop tool (IDE, analytics dashboard, CAD tool), design desktop-first and treat mobile as a degraded read-only or limited-function experience. Document this explicitly.

---

## 3. Standard Breakpoint System

Use a consistent named breakpoint system. Never use arbitrary pixel values in media queries.

```css
/* Breakpoint token conventions */
/* xs:  < 480px  — Small phones */
/* sm:  ≥ 480px  — Large phones */
/* md:  ≥ 768px  — Tablets, landscape phones */
/* lg:  ≥ 1024px — Small desktops, landscape tablets */
/* xl:  ≥ 1280px — Standard desktops */
/* 2xl: ≥ 1536px — Large desktops, wide monitors */
```

For Tailwind projects, these map directly to the default `sm`, `md`, `lg`, `xl`, `2xl` breakpoints. Configure them once in `tailwind.config.js` and never hardcode pixel values in templates.

### Breakpoint Usage Rules
- Do not create layouts that only break at `1px` increments. Use the named scale.
- Test at the following specific viewport widths: 320px, 375px, 430px (phones), 768px (tablet portrait), 1024px (tablet landscape / small desktop), 1280px (desktop), 1440px (wide desktop).
- The interface must function — not just display — at all standard breakpoints.

---

## 4. Layout Patterns by Viewport

### Navigation
| Viewport | Pattern |
| :--- | :--- |
| Mobile | Bottom tab bar (if ≤5 items) or hamburger menu drawer |
| Tablet | Collapsible sidebar or top navigation bar |
| Desktop | Persistent sidebar or expanded top navigation bar |

Do not use hamburger menus on desktop. Do not use persistent sidebars on mobile unless the navigation depth requires it.

### Data Tables
| Viewport | Pattern |
| :--- | :--- |
| Mobile | Priority columns only visible; remaining columns accessible via row expand or horizontal scroll |
| Tablet | Most columns visible with responsive column hiding |
| Desktop | Full table with all columns |

Never hide primary data or the primary action column on any viewport.

### Forms
| Viewport | Pattern |
| :--- | :--- |
| Mobile | Single-column layout; full-width inputs |
| Tablet | Two-column layout for closely related fields (e.g., First Name / Last Name) |
| Desktop | Multi-column layout with clear grouping |

Never force multi-column form layouts on mobile.

### Cards and Grids
| Viewport | Columns |
| :--- | :--- |
| Mobile | 1 column |
| Tablet | 2 columns |
| Desktop | 3–4 columns |

Do not exceed 4 columns for card grids. Reduce column count before reducing card content.

---

## 5. Container Queries

Prefer **container queries** over media queries for component-level responsive behavior. Container queries respond to the component's parent container width, not the viewport width.

```css
.card-container {
  container-type: inline-size;
  container-name: card;
}

@container card (min-width: 400px) {
  .card-body {
    display: flex;
    gap: var(--space-4);
  }
}
```

### When to Use Container Queries
- Components that appear in multiple layout contexts (sidebar, main content, modal).
- Components whose layout depends on their container, not the page.

### When to Use Media Queries
- Page-level layout changes (column counts, sidebar visibility, navigation patterns).
- Global adjustments (font scaling, body padding).

---

## 6. Touch and Gesture Interaction

Mobile users interact via touch targets and gestures. This requires different hit target boundaries and fallback interaction states.

### Touch Target Rules
- **Minimum Target Size**: **44×44px** (WCAG 2.5.5 requirement).
- **Control Spacing**: Do not place interactive touch targets closer than **8px** apart on mobile.
- **Hover Dependency**: Hover interactions must never be the only way to trigger functionality on touch viewports. Ensure dropdowns and tooltips open on tap events.

### Touch Gestures & Discoverable Alternatives
Provide gesture integrations (swipes, pull-to-refresh, long-press, drag-and-drop) to make the UI feel smooth, but **every gesture must have a discoverable text or pointer alternative**:
- *Swipe-to-delete* must have a visible "Delete" button option or toggle icon.
- *Pull-to-refresh* must have a standard manual "Refresh" icon link or button fallback.
- *Long-press* actions must be triggerable via standard double-tap or an options menu icon button.

### Input Method Considerations
```css
/* Remove hover-only behavior on touch devices */
@media (hover: none) {
  .dropdown-trigger:hover .dropdown-menu {
    /* Do not rely on hover for dropdown visibility on touch */
    display: none;
  }
}
```

---

## 6.5. Device Features & Safe Integration
Ensure UI logic adapts gracefully to device capabilities:
- **Feature Check gates**: Before calling camera API, clipboard integrations, share sheets, or native file pickers, verify target compatibility.
- **No hardware requirements**: Never make a native hardware API a mandatory blocker. Provide standard fallback options (e.g. manual text input copy box instead of write-to-clipboard API on failure).

---

## 6.6. Device Safe Areas & Display Notches
Ensure UI elements are not obscured by operating system overlays or screen notches:
- **Rounded Corners & Notches**: Avoid positioning interactive controls at the extreme top/bottom edges of the screen where notches, the iPhone Dynamic Island, or rounded corners clip content.
- **System Navigation Bars**: Respect bottom system bars (e.g., Apple home indicator line, Android system buttons).
- **CSS Safe Area Tokens**: Use environment variables to set layout paddings:

```css
.app-header {
  padding-top: max(var(--space-4), env(safe-area-inset-top));
}

.bottom-actions-bar {
  padding-bottom: max(var(--space-4), env(safe-area-inset-bottom));
  padding-left: env(safe-area-inset-left);
  padding-right: env(safe-area-inset-right);
}
```

---

## 6.7. Screen Orientation Support
- Ensure all layouts adapt to both **Portrait** and **Landscape** orientations.
- Never force users into a single orientation unless the product requires it (e.g. full-screen immersive games).
- Use flexible container widths and layout parameters to allow dynamic recalculation when orientation toggles.

---

## 7. Fluid Typography and Spacing

Avoid fixed type sizes that are too small on mobile or too large on wide screens.

### Fluid Typography with `clamp()`
```css
/* Font size scales fluidly between 16px (mobile) and 20px (desktop) */
--font-size-base: clamp(1rem, 0.875rem + 0.625vw, 1.25rem);

/* Heading scales between 24px and 40px */
--font-size-display: clamp(1.5rem, 1.2rem + 1.5vw, 2.5rem);
```

### Fluid Spacing
```css
/* Section padding scales with viewport */
--section-padding: clamp(var(--space-8), 5vw, var(--space-16));
```

Only use fluid values for layout-level spacing and display-level typography. Component-level spacing should use fixed tokens from the spacing scale.

---

## 8. Responsive Images

Image handling is one of the most common responsive design failures.

### Rules
- Never serve desktop-resolution images to mobile viewports.
- Use `srcset` and `sizes` attributes for raster images.
- Use `width` and `height` attributes on `<img>` tags to prevent layout shift (CLS).
- Use `loading="lazy"` for images below the fold.
- Use SVG for icons, illustrations, and logos — they scale without quality loss.
- Use modern formats: WebP or AVIF as primary with JPEG/PNG fallback.

```html
<img
  src="photo-800.jpg"
  srcset="photo-400.jpg 400w, photo-800.jpg 800w, photo-1200.jpg 1200w"
  sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
  width="800"
  height="600"
  loading="lazy"
  alt="Descriptive text about the image"
>
```

---

## 9. Performance at Small Viewports

Mobile devices often have slower network connections and less powerful CPUs than the developer's machine.

- Do not load desktop-only assets on mobile. Use conditional loading.
- Avoid large JavaScript bundles that block the main thread on load.
- Defer non-critical styles and scripts.
- Test on a real device under simulated slow network conditions (Chrome DevTools throttling) before shipping.

---

## 10. Responsive Framework Guidance

- **Vue / Nuxt**: Use clean responsive template blocks. Rely on media query utility hooks or Tailwind variants dynamically rather than maintaining parallel mobile/desktop templates.
- **React / Next.js**: Prefer dynamic adaptive rendering using CSS hooks or layout wrappers over duplicating entire layout structures. Keep props and data flows synchronized.
- **Flutter**: Ensure layout constraints scale dynamically using widgets like `LayoutBuilder`, `MediaQuery`, or `AspectRatio`. Respect Material and Cupertino platform design spacing conventions.

---

## 11. AI UX Review Questions

Prior to layout verification, evaluate the interface using these questions:
1. *Does this interface feel designed for mobile from the ground up, rather than a shrunk-down desktop screen?*
2. *Can all interactive elements and primary buttons be navigated and activated with one hand where practical?*
3. *Is the most important context and information immediately visible in the active mobile viewport?*
4. *Is the vertical scrolling height reasonable, or does it force excessive navigation to complete the task?*
5. *Does the layout contain any edge-to-edge text overlaps or text shifts that would frustrate a phone user?*

---

## Review Checklist

Before completing UI work, verify the layout against this responsive checklist:
- [ ] **Mobile-first layout**: CSS styling is written mobile-first, progressively enhancing as viewports grow larger.
- [ ] **Flexible grids**: Grid columns and Flexbox items use relative units (percentages, fractional units) rather than fixed pixel layouts.
- [ ] **Responsive typography**: Heading sizes and text blocks scale smoothly using fluid typography variables (`clamp()`).
- [ ] **Comfortable touch targets**: Interactive inputs and buttons are at least 44×44px with a minimum 8px spacing clearance.
- [ ] **Adaptive navigation**: Menus adapt correctly (e.g. collapsing sidebar to bottom navigation tab bar on mobile viewports).
- [ ] **Responsive tables**: Tables wrap column sizes or display as layout-friendly cards on small screens.
- [ ] **Responsive forms**: Inputs display full-width on mobile; standard keyboards match expected input formats.
- [ ] **Optimized images**: Vector SVGs are used for UI controls; raster images use `srcset` and include dimensions.
- [ ] **Accessible interactions**: Keyboard inputs, zoom levels, and focus boundaries are fully supported at all breakpoints.
- [ ] **Performance considered**: Mobile CPU, RAM, and network throttling profiles have been tested with zero layout freezes.

---

## References
- Design Systems (spacing tokens): [design/01-design-systems.md](01-design-systems.md)
- Accessibility (touch targets): [design/04-accessibility.md](04-accessibility.md)
- Visual Quality: [design/05-visual-quality.md](05-visual-quality.md)
- Every Layout (intrinsic CSS patterns): [https://every-layout.dev](https://every-layout.dev)
