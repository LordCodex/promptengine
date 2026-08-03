---
document_id: workflows-ui-review
title: UI Review & Verification Workflow
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-03
---

# UI Review & Verification Workflow

This document details the review protocol to guarantee that all frontend elements, layouts, and animations comply with established UX, accessibility, and visual guidelines.

---

## 1. Visual Quality Assessment
- Verify typography scaling, color contrast, and spacing systems conform to design tokens (refer to [Design Systems](../design/01-design-systems.md)).
- Check all component states (loading, empty, error, disabled) as defined in [Visual Quality](../design/05-visual-quality.md).
- Ensure animations are subtle and do not delay core user workflows.

---

## 2. Accessibility (a11y) Verification
- Perform keyboard navigation audits to ensure focus states are logical and visible (refer to [Accessibility](../design/04-accessibility.md)).
- Check for correct semantic HTML elements and ARIA tags.
- Run DevTools AXE audits or voiceover checks to guarantee accessibility compliance.

---

## 3. Responsive and Device Auditing
- Verify layout structures wrap and scale correctly across device break points (mobile, tablet, desktop) as defined in [Responsive Design](../design/03-responsive-design.md).
- Ensure touch targets are at least `44x44px` to support physical interactive inputs.

---

## 4. UI Gate Sign-off
- Perform a self-review using the [UI Quality Review Checklist](../checklists/04-ui-quality-review.md).
- Document and resolve any deviation from Figma specifications before requesting peer reviews.
