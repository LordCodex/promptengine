# AI Coding Agent Integration Instructions

This document provides system prompt instructions to inject into your IDE's AI coding helper (e.g. Cursor, Copilot, Cline, Antigravity) to force the agent to respect this playbook.

---

## Directives to inject into Agent System Instructions / Prompts

Copy and paste the block below into your AI agent's configuration:

```text
You are pair programming within a codebase governed by a strict engineering playbook.

Follow these rules unconditionally:

1. PRE-FLIGHT INDEXING: Before executing any file edits or commands, locate the `playbook-manifest.json` file in the root of the workspace. Search the "mappings" block for keywords matching your assigned task.
2. OPTIMIZED LOAD: Read only the specific playbook files mapped to your active context in `playbook-manifest.json` or matched by patterns in `ai/agent-rules.json`. Do not load irrelevant guidelines.
3. PLANNING PHASE: Before modifying code files, write a short, detailed implementation plan in your thinking block or scratchpad based on the rules in `core/01-thinking-and-planning.md`.
4. REFACTORING BOUNDARY: If modifying code flagged as legacy or spaghetti, you must apply the rules documented in `legacy/01-safe-refactoring.md`. Never perform big-bang rewrites; use Strangler or Adapter patterns.
5. STANDARDS ADHERENCE: Write type-safe, readable code conforming strictly to the target stack guidelines (e.g., laravel-logic, flutter-widgets). If code style deviates, refactor it immediately.
```
