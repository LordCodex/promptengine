# 05. Bug Fix Prompt

---

## Purpose
Instructs the AI to isolate the root cause of an application bug, generate a minimal patch, and write a reproduction test.

## When to use
Use when you want to troubleshoot unexpected behaviors, logic errors, or system crashes.

## Example
Fixing a null pointer exception on user login callback.

---

## Copy-and-Paste Prompt

```markdown
I want you to investigate and fix a bug in our application.

Target Stack: {STACK}
Target Path/Module: {PATH}
Symptoms: {SYMPTOMS}
Reproduction Steps: {REPRODUCTION_STEPS}

Before writing any code:
1. Scan the target codebase files to locate the root cause.
2. Read the project constitution (`AGENTS.md`) and the relevant specs under `docs/` (such as `docs/Troubleshooting.md` or `docs/API.md`).
3. Outlines the root cause diagnostics and proposed bugfix. Keep the patch as minimal as possible to avoid side effects.
4. Draft a test case that reproduces the bug (fails before the patch, passes after the patch).
5. Wait for my confirmation before modifying the codebase files.
```

---

## Expected AI Behaviour
1. The AI scans the target files to trace inputs/outputs.
2. It outputs an explanation of why the bug occurs.
3. It drafts a minimal code fix and a reproduction test case.
4. It waits for your approval before writing the patch to disk.

## Common Mistakes
- **Applying speculative fixes**: Letting the AI generate guesses about why a config is broken. Force it to explain the evidence in code before applying patches.
- **Unrelated modifications**: The AI refactoring adjacent classes while applying a fix. Enforce: *"Minimal patch only."*
