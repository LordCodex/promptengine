# 07. Architecture Change Prompt

---

## Purpose
Instructs the AI to design and plan a major architectural or database schema change, draft an Architecture Decision Record (ADR), and write migrations.

## When to use
Use when changing database design, introducing caching layers, changing auth protocols, or modifying module structures.

## Example
Adding a Redis caching layer for the Product Catalog API.

---

## Copy-and-Paste Prompt

```markdown
I want to plan an architectural change in our system.

Change Goal: {CHANGE_GOAL}
Tech Stack: {STACK}

Before writing any code or database migration scripts:
1. Read `AGENTS.md` and `docs/Decisions.md` (or the decisions log).
2. Outline the design proposal, including:
   - Alternating options considered.
   - Specific trade-offs (performance, security, complexity).
   - SQL DDL migrations or framework schema code.
   - Pinned version libraries to add (if any).
3. Draft a new Architecture Decision Record (ADR) under `docs/Decisions.md` (or the decisions log).
4. If this change affects our technology stack or constraints, update Section 2 of `AGENTS.md`.
5. Wait for my approval before proceeding to implementation.
```

---

## Expected AI Behaviour
1. The AI reads decisions log and configuration files.
2. It outputs an ADR draft (Context, Decision, Consequences, Status: Proposed).
3. It outlines migration files and code adapters to write.
4. It waits for your approval of the ADR before writing migrations to disk.

## Common Mistakes
- **Undocumented schema evolution**: Modifying database columns in code without generating an ADR or updating `docs/Database.md`.
- **Ignoring soft/hard constraints**: Changing primary key strategies (e.g. using integers instead of UUIDs) without recording it as an approved exception in `AGENTS.md`.
