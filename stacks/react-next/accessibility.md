---
document_id: stacks-react-accessibility
title: React and Next.js Accessibility Standards
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# React and Next.js Accessibility Standards

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines accessibility standards for React and Next.js projects.

---

## 1. Semantic Elements and Interactive Roles

- **Use Semantic HTML Elements**: Prioritize standard elements (`<button>`, `<a href>`, `<nav>`) over custom styled `<div>` layout blocks.
- **Dynamic Aria States**: Update `aria-expanded`, `aria-selected`, and `aria-invalid` values dynamically based on component state variables.
- **Focus Management**:
  - Implement focus trap boundaries inside modal dialog blocks.
  - Return focus to the trigger element when overlays close.

---

## 2. Headless primitives & Components

- **Headless libraries**: When building complex layouts (toggles, accordions, select menus), build on accessible primitive layers (e.g. Radix UI primitives) instead of custom markup solutions.
- **Image Alts**: Enforce descriptive alt tags on all Next.js `<Image>` tags.

---

## Review Checklist

Verify accessibility compliance against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] Interactive elements maintain keyboard navigation.
- [ ] Modals traps focus inside active boundaries.
- [ ] All images declare descriptive alt text attributes.
