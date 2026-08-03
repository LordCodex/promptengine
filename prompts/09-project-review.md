# 09. Project Review Prompt

---

## Purpose
Instructs the AI to perform a comprehensive code review on a git diff or file, focusing on security boundaries, performance bottlenecks, and accessibility (a11y) violations.

## When to use
Use during pull request checks or code audit cycles before staging features.

## Example
Reviewing a new user registration route and view controller.

---

## Copy-and-Paste Prompt

```markdown
I want you to perform a security, performance, and accessibility code review.

Target Files / Git Diff: {REVIEW_SOURCE}
Target Stack: {STACK}

Please perform the review following these rules:
1. **Security Audit**: Check for missing authorization gates (RBAC/Policies), unescaped raw outputs, parameter validation weaknesses, and secrets leakage. Apply the Three Questions (Who sent this? Are they allowed? Is the data safe?).
2. **Performance Audit**: Check for N+1 database queries, un-indexed search columns, missing caching layers, or memory-heavy array loading.
3. **Accessibility Audit**: Check for semantic HTML tags, keyboard focus management, ARIA labels, and color contrasts.
4. **Severity Classification**: Categorize all findings using severity prefixes:
   - `[BLOCK]` (Critical flaw requiring immediate remediation).
   - `[IMPORTANT]` (Important recommendation to fix soon).
   - `[SUGGESTION]` (Optional style or cleanup improvement).
5. For every finding, provide secure or performant code replacement proposals.
```

---

## Expected AI Behaviour
1. The AI evaluates the target code.
2. It outputs a structured review report divided into Security, Performance, and Accessibility sections.
3. It labels issues with `[BLOCK]`, `[IMPORTANT]`, or `[SUGGESTION]`.
4. It provides code snippets showing how to fix each issue.

## Common Mistakes
- **Subjective style critiques**: The AI commenting on braces or formatting styles. Force it to focus on core standards: *Only technical, standard-based findings.*
- **Ungrounded warnings**: Suggesting issues that do not apply to the project context.
