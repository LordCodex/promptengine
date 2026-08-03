# 06. Refactor Prompt

---

## Purpose
Instructs the AI to safely refactor a code module to improve design, readability, or testability while preserving external behaviors.

## When to use
Use when flattening complex nested structures, separating responsibilities, or modernizing older classes.

## Example
Extracting validation logic from a controller to a dedicated action class.

---

## Copy-and-Paste Prompt

```markdown
I want to refactor the following class/module to improve its internal structure.

Target Path: {PATH}
Current Test Coverage State: {TESTS_STATE}

Follow the Safe Refactoring Workflow:
1. Verify that all automated tests currently pass for this module.
2. Outline a refactoring plan in your thinking block. Propose small, step-by-step extractions (e.g. Early Exits, flattening nested if blocks, extracting services).
3. Do not modify public API URLs, payload keys, database fields, or validation parameters.
4. Execute the refactoring incrementally, running tests after each atomic change.
5. Wait for my approval of the plan before writing code edits.
```

---

## Expected AI Behaviour
1. The AI reviews the target code.
2. It outputs a step-by-step refactoring checklist showing what variables to extract, what loops to simplify, and what early exits to introduce.
3. It waits for your approval before starting.
4. It edits the files in small increments.

## Common Mistakes
- **Bulk rewriting**: The AI changing 500 lines of code in one prompt. If it does this, reject and remind it: *Step-by-step refactoring only, running tests after each change.*
- **Changing behavior**: Bypassing custom business logic. Check that return payloads and validation limits are unchanged.
