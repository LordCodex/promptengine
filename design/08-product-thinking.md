---
document_id: design-product-thinking
title: Product Thinking and User Flow Standard
ecosystem: cross-cutting
dependencies:
  - design-ui-ux-philosophy
  - design-systems
  - design-visual-quality
  - core-architecture-and-simplicity
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Product Thinking and User Flow Standard

## Inheritance
This document inherits from and extends the [UI/UX Philosophy](00-ui-ux-philosophy.md). That document establishes *design principles*. This document establishes *product engineering process* — the structured thinking an AI agent or developer must complete **before** creating any user-facing feature.

---

## 1. The Correct Starting Question

**Do not start with:** "What components should I create?"

**Start with:** "What problem are we solving?"

Every feature implementation must begin by understanding the problem, not the solution. A developer or agent who skips this step will build a technically correct but productively useless interface.

The implementation starts only after the product problem is fully understood.

---

## 2. Mandatory Pre-Implementation Declaration

Before implementing any user-facing feature, an AI agent must produce this declaration. It is not optional. It is the gate between requirement and code.

```markdown
## Feature: [Feature Name]

### 1. User Problem
[What specific problem does this solve? Whose problem is it? How do they currently work around it?]

### 2. User Flow
[Entry point → Steps → Success state — written as a numbered sequence of user actions and system responses]

### 3. Required Screens or Components
[List only what is necessary. Justify each one.]

### 4. Existing Components to Reuse
[What already exists in the project that covers part of this? Where are the files?]

### 5. Edge Cases
[What can go wrong? What happens at the boundaries? Invalid input, empty data, network failure, permission denial.]

### 6. Accessibility Considerations
[Which interactions require keyboard support? Which dynamic states require ARIA announcements?]

### 7. Responsive Considerations
[How does this feature behave on mobile? What changes at small viewports?]
```

If any section cannot be answered, **stop and clarify** before writing code.

---

## 3. User Persona Framework

Different users have different goals, contexts, and constraints. An interface designed for the wrong user type creates friction and fails adoption.

Before designing a workflow, identify the user:

| Dimension | Questions to Answer |
| :--- | :--- |
| **Role** | What is the user's job or position? Are they a customer, admin, operator, or superuser? |
| **Permissions** | What actions are they authorized to perform? What is restricted? |
| **Technical ability** | Are they a power user, a casual user, or a first-time user? |
| **Frequency** | How often do they use this feature? Daily? Weekly? Once at onboarding? |
| **Context** | Where are they when they use this? At a desk? On a mobile device? In a hurry? |
| **Primary frustration** | What makes this task harder than it needs to be today? |

### Design Priority by User Type

The user type determines what the interface must optimize for:

| User Type | Optimize For |
| :--- | :--- |
| **Daily power user** (admin, operator, back-office) | Speed, keyboard efficiency, data density, shortcuts |
| **Infrequent expert user** (manager, reviewer) | Clarity, summary, confidence, low learning curve |
| **First-time customer** | Guidance, trust signals, simplicity, progressive disclosure |
| **Technical user** (developer, analyst) | Data access, export, precision, low abstraction |
| **High-stakes user** (finance, compliance, medical) | Confirmation, audit trail, error prevention, data accuracy |

Do not apply the same interface pattern to different user types. A daily admin workflow must not look like a first-time onboarding flow.

---

## 4. User Journey Mapping

Every major feature requires a journey map before layout decisions are made. The map answers: what does the user actually experience from start to finish?

### 4.1 Entry Point

How does the user reach this feature?

- **Dashboard** — they navigate to it proactively.
- **Navigation** — it is always accessible.
- **Notification / Alert** — they are pulled into it reactively.
- **Search result** — they discover it while looking for something else.
- **Deep link / External URL** — they arrive from email, SMS, or another system.
- **Inline action** — they trigger it from within another task.

The entry point determines what context the user carries into the feature. A user arriving from a notification already knows something is wrong. A user arriving from navigation is exploring.

### 4.2 Main Flow

Define the primary sequence as a numbered list of alternating user actions and system responses:

```text
1. User: [action]
   System: [response]

2. User: [action]
   System: [response]

3. ...
```

Rules:
- Every user action must produce a visible system response.
- No step should leave the user uncertain about what to do next.
- Identify which step is the **point of no return** (after which the action cannot be undone).

### 4.3 Success Path

Define what success looks like from the user's perspective:

- What confirmation message or visual change signals completion?
- What information must be displayed after success (reference number, timestamp, next step)?
- Where does the user go next? Do they need to be redirected, or does the success state appear inline?
- Does the user need to perform another action immediately (pay → download receipt, sign up → verify email)?

### 4.4 Failure Paths

Every journey has failure modes. Design for them before designing the happy path.

| Failure Type | Required UI Response |
| :--- | :--- |
| **Invalid input** | Inline validation message on the specific field; explain what is wrong and how to fix it |
| **Network failure** | Preserve entered data; offer a retry action; do not lose work |
| **Permission denial** | Explain why access is denied and who to contact, if applicable |
| **Missing or empty data** | Explain what is missing and what action can generate or supply it |
| **System error** | Plain-language explanation; reference number for support if applicable; do not expose technical details |
| **User cancellation** | Confirm intent for multi-step or destructive flows; preserve draft data where possible |
| **Timeout** | Warn before expiry; offer session extension where technically possible |

---

## 5. Information Architecture

Before choosing a layout, determine the information hierarchy. Layout must follow hierarchy — not the other way around.

### Information Priority Questions

Answer these before wireframing or coding:

1. **What must the user see first?** (Primary — immediately visible on load)
2. **What supports the primary information?** (Secondary — visible but subordinate)
3. **What is only needed on demand?** (Tertiary — accessible via expand, tab, or drill-down)
4. **What should never appear unless explicitly requested?** (Hidden — advanced settings, raw data, audit logs)

### Progressive Disclosure

Do not show everything at once. Surface only what the user needs at each step.

- A form with 20 fields is often 3 forms with 6–8 fields each, shown progressively.
- Advanced settings belong behind a collapsible "Advanced" section or a separate screen.
- Confirmation details (itemized breakdown, legal text) belong on a review step, not the initial action screen.
- Metadata and audit history belong on a detail view, not a list view.

### Grouping Rules

Group information that belongs together. Separate information that belongs apart.

- Fields that refer to the same subject (name, email, phone) belong in the same section.
- Sections with different subjects (personal details, billing details, preferences) belong in separate sections with clear headings.
- Actions (submit, cancel, save draft) must be spatially separated from informational content.
- Destructive actions (delete, revoke, archive) must be separated from constructive actions (save, create, publish).

---

## 6. Form Design Intelligence

Forms are the most common mechanism for user input. Poor form design is the most common source of abandonment and error.

### Before Submission

- Every label must describe what to enter, not just what the field is named in the database (`"Your email address"` not `"email"`).
- Helper text should explain format, constraints, or why the information is needed — not restate the label.
- Use the correct input type for the data:
  - `type="email"` for emails — triggers the correct mobile keyboard.
  - `type="tel"` for phone numbers.
  - `type="number"` for quantities.
  - `type="date"` or a date picker component for dates.
- Provide sensible defaults where the correct value can be inferred (country from location, currency from account, today's date for a "start date" field).
- Group related fields visually using `<fieldset>` + `<legend>` or clearly headed sections.
- Never put more than 8–10 fields in a single unsectioned form. If the form needs more, it needs structure.

### During Input

- Validate on blur (when the user leaves the field), not on every keystroke — keystroke validation is distracting.
- Exception: real-time validation is appropriate for password strength, username availability, and character count limits.
- Error messages must describe what is wrong and how to fix it — not just "Invalid input."
  - Wrong: `"Email is invalid."`
  - Correct: `"Enter an email address in the format: name@example.com"`
- Preserve all entered data on validation failure. Never clear a form on error.
- Do not disable the submit button until submission is attempted. Disabling it pre-emptively removes the user's ability to discover which fields are invalid.

### After Submission

- Provide immediate, unambiguous success feedback. Do not redirect silently.
- Tell the user what happens next:
  - "Your application has been submitted. You will receive a confirmation email within 5 minutes."
  - "Payment successful. Your receipt has been sent to you@example.com."
- If the next step requires user action (verify email, complete profile), provide the action immediately.
- If the form will be reviewed by a human, tell the user the expected timeline.

---

## 7. Data Display Intelligence

Choosing the wrong data display pattern forces users to work harder than they need to. Match the pattern to the nature of the data and the user's task.

### Pattern Selection Rules

| Pattern | Use When | Do Not Use When |
| :--- | :--- | :--- |
| **Table** | Users need to compare multiple attributes across many items; data is structured and uniform | Items have very different shapes; content is narrative or visual |
| **List** | Items are sequential or ranked; users scan and select one at a time | Users need to compare attributes across items |
| **Card grid** | Items are discrete, visually distinct, and comparable at a glance; items have a thumbnail or visual | Data is dense; users need to compare many fields; items number in the hundreds |
| **Feed / Timeline** | Items are chronological; each item is contextually self-contained | Items are structured data that benefits from column comparison |
| **Form** | The user is providing or editing structured input | The user is only reading data |
| **Detail view** | A single item requires full description with many attributes | Summary information is sufficient |
| **Chart / Graph** | Visualizing trends, distributions, or comparisons where pattern recognition adds value | Displaying a single number; the data is not better understood visually |

### Chart Discipline

Do not add charts because "dashboards have charts." Add charts only when:

- The visual representation reveals a pattern that a number or table cannot communicate.
- The user's task is to understand a trend, distribution, or comparison — not to retrieve a specific value.
- The data changes over time and the rate or direction of change matters.

If the user only needs to know the current value, show the number. A chart to display one value is decoration.

### Table Design Rules

- Always provide column headers with clear, non-abbreviated labels.
- Make columns sortable where the user's task involves finding items by value.
- For tables with more than 20 rows, provide pagination or infinite scroll — never a single unbounded list.
- Use right-alignment for numeric columns to enable visual comparison.
- Provide a clear empty state when the table has no rows.
- On mobile, reveal only the most critical columns; provide a row-expand or detail view for the rest.

---

## 8. Confirmation and Trust Patterns

Certain actions carry real-world consequences. The interface must communicate these consequences clearly and provide appropriate confidence before the user commits.

### Payment Confirmation
Before a payment is finalized, display:
- Exact amount with currency symbol.
- What is being paid for (description, quantity, period).
- Payment method being charged.
- A final "Confirm payment" action — never auto-submit.

After payment:
- Payment status (succeeded, pending, failed).
- Reference or transaction ID.
- Receipt or confirmation to the user's contact method.

### Destructive Action Confirmation
Before an irreversible action (delete, archive, revoke, ban):
- Name the specific item being affected — never "Are you sure?" without context.
- Explain the consequence: "This will permanently delete 47 orders and cannot be undone."
- For high-risk actions, require the user to type the resource name or "DELETE" to confirm.
- Provide a cancel path that is at least as prominent as the confirm path.

### Account and Security Changes
- Confirm email changes by sending a verification link to the **old** email address.
- Confirm password changes with the current password before accepting the new one.
- Display a security notification (in-app and email) after sensitive account changes.
- Show the timestamp and device for recent account activity.

---

## 9. Business Logic Awareness

The UI must represent the actual system behavior. It must not invent, assume, or simplify business rules that the backend enforces.

### Before Designing Any Feature

Understand:
- What are the validation rules? (Required fields, allowed values, format constraints)
- What permissions govern this action? (Who can see it, who can do it, who can approve it)
- What are the existing user expectations? (How does this currently work? What would change?)
- What are the downstream effects? (What happens in the system when this action is completed?)

### Rules
- If a business rule is unclear, ask before implementing. Do not guess.
- If the UI would require the user to work around a backend constraint, surface the constraint — do not hide it.
- If a field is read-only because of a business rule (e.g., invoice number is auto-generated), display it as read-only with an explanation — do not make it editable.
- Validation that happens server-side must be communicated on the form — do not make the user guess why their submission was rejected.

---

## 10. Consistency with the Existing Product

A new feature that introduces a completely different visual or interaction pattern creates dissonance. Users learn a product's mental model. Deviating from it without reason increases cognitive load.

### Before Creating a New Experience

Review:
- **Existing navigation** — Is there already a navigation pattern this feature belongs within?
- **Existing terminology** — What does the product call the things this feature involves? Use those words, not new ones.
- **Existing components** — Which components already exist that this feature can compose? (See [design/02-component-libraries.md](02-component-libraries.md))
- **Existing workflows** — Is there a similar workflow in the product this feature should mirror?
- **Existing patterns** — How does the product handle forms, tables, modals, and confirmations? Match those patterns.

### Rules
- Do not introduce a new interaction pattern when an existing one works.
- Do not introduce new terminology for concepts the product already names.
- Do not create a new navigation section when an existing section is the correct home.
- If a deviation from existing patterns is genuinely justified, document the reason explicitly.

---

## 11. Over-Design Prevention

The most common failure mode of AI-generated product design is adding complexity that was never requested.

**Do not add:**
- Features nobody requested ("while I'm here, I added a filter panel").
- Extra steps that add confirmation without reducing risk.
- Onboarding flows for features that are self-explanatory.
- Animations that serve no state-change communication purpose.
- Charts or visualizations for data that is fully communicated by a number.
- Settings and customization options that only 1% of users would ever change.

**The simplest experience that solves the problem is the correct answer.** Complexity must be earned by user need, not added by default.

---

## References
- UI/UX Philosophy (design principles): [design/00-ui-ux-philosophy.md](00-ui-ux-philosophy.md)
- Design Systems (tokens, typography, spacing): [design/01-design-systems.md](01-design-systems.md)
- Component Libraries (evaluate before building): [design/02-component-libraries.md](02-component-libraries.md)
- Responsive Design: [design/03-responsive-design.md](03-responsive-design.md)
- Accessibility: [design/04-accessibility.md](04-accessibility.md)
- Visual Quality (component states, data display): [design/05-visual-quality.md](05-visual-quality.md)
- UI Review Process (6-layer self-review): [design/06-ui-review-process.md](06-ui-review-process.md)
- UI Quality Checklist (gate): [checklists/04-ui-quality-review.md](../checklists/04-ui-quality-review.md)
