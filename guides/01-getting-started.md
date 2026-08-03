# 01. Getting Started with PromptEngine

Welcome to PromptEngine! This guide introduces you to the core principles and structure of PromptEngine from a developer's perspective. It assumes you are new to the framework and provides a quick conceptual model.

---

## What is PromptEngine?

PromptEngine is a structured, machine-readable software engineering knowledge base. Think of it as an **Engineering Operating System** for your codebase. It defines how software should be planned, structured, coded, tested, and reviewed, making those rules directly understandable and enforceable by AI coding assistants (like Cursor, Claude Code, Windsurf, or Codex).

By combining reusable coding standards with a local Project Knowledge System, PromptEngine ensures you write high-quality code that conforms to your architectural standards without having to manually remind the AI of your guidelines in every chat.

---

## The Problem: Volatile Chat Memory

If you have ever pair-programmed with an AI coding helper, you have likely run into these issues:
1. **Context Erosion**: After a few thousand tokens of chat, the AI starts forgetting rules set at the beginning of the conversation.
2. **Framework Scaffolding Over-Reliance**: The AI blindly writes code using framework defaults, ignoring your custom patterns.
3. **Redundant Code Generation**: The AI generates duplicate helpers or models because it does not realize they already exist.
4. **Session Loss**: Starting a new chat session clears all previous alignment, requiring you to explain your project specifications all over again.

---

## The Solution: Repository as the Source of Truth

PromptEngine solves this by shifting project rules and specifications out of volatile chat memory and into the repository itself. 

It splits engineering knowledge into three distinct layers:

```text
┌────────────────────────────────────────────────────────┐
│ 1. Core Engineering Standards                          │
│    (Agnostic playbooks: security, testing, performance) │
└───────────────────────────┬────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│ 2. Technology Stack Playbooks                          │
│    (Conventions: PHP/Laravel, JS/Vue, Dart/Flutter)    │
└───────────────────────────┬────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────┐
│ 3. Project Knowledge System (Root Level)               │
│    ├─ AGENTS.md (AI Constitution: Tech stack & rules)   │
│    └─ docs/ (PRD, API.md, Database.md, Decisions.md)   │
└────────────────────────────────────────────────────────┘
```

1. **Core Standards**: General, language-agnostic rules for clean code, planning, and security (stored in `/core/`).
2. **Stack Playbooks**: Casing, structure, and design conventions specific to your programming language or framework (stored in `/stacks/`).
3. **Project Knowledge System**: Project-specific requirements, database schemas, API routes, and guidelines (stored in the project's root `AGENTS.md` and `/docs/` folder).

---

## How It Works in Practice

When you open your project in an AI-powered IDE or command-line tool, the AI is configured to execute the **5-Step Entry Rule**:

1. **Read AGENTS.md**: The AI loads the project constitution in your root directory.
2. **Read specs in docs/**: The AI reads the requirements (`docs/PRD.md`), schema (`docs/Database.md`), or API (`docs/API.md`) related to your request.
3. **Query Playbooks**: The AI looks up standard PromptEngine playbooks mapped to the active stack.
4. **Determine Workflow**: The AI selects the appropriate workflow (e.g. adding a feature, fixing a bug).
5. **Plan & Execute**: The AI writes a planning block in its thinking space first, gets your approval, and implements the change.

---

## Next Steps

To begin using PromptEngine:
- If starting a brand-new project, follow the **[Greenfield New Projects Guide](02-new-project.md)**.
- If adopting PromptEngine in a pre-existing codebase, follow the **[Active Existing Projects Guide](03-existing-project.md)**.
