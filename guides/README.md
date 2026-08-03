# PromptEngine Developer Guides

Welcome to the PromptEngine Developer Guides. This directory contains practical, step-by-step documentation designed to help developers use PromptEngine to govern software projects with AI coding assistants.

---

## Guide Index

1. **[01. Getting Started](01-getting-started.md)**: An introduction to the core concepts of PromptEngine as an "Engineering OS" for AI pair programming.
2. **[02. Greenfield New Projects](02-new-project.md)**: How to bootstrap a brand-new project with PromptEngine, conduct the discovery phase, and auto-generate documentation.
3. **[03. Active Existing Projects](03-existing-project.md)**: How to adopt PromptEngine in active, pre-existing codebases and reverse-engineer specifications.
4. **[04. Migrating Legacy Systems](04-migrate-existing-project.md)**: A step-by-step guide to migrating and modernizing older legacy codebases safely.
5. **[05. Daily Development Workflow](05-daily-development.md)**: Practical guidelines for daily development tasks (adding features, fixing bugs, refactoring, and releases).
6. **[06. Managing Project Documentation](06-project-documentation.md)**: Deep dive into the 10 core documentation files and how to keep them synchronized.
7. **[07. AI Prompt Library Integration](07-ai-prompt-library.md)**: How to leverage pre-written, reusable AI prompts to steer coding assistants.
8. **[08. Developer Best Practices](08-best-practices.md)**: Tips for token optimization, avoiding context loss, and integrating PromptEngine with ChatGPT, Claude, Cursor, Windsurf, and Codex.
9. **[09. Troubleshooting Guide](09-troubleshooting.md)**: How to remediate common issues such as AI ignoring rules, documentation drift, or context loss.
10. **[CLI Foundation Index](../cli/README.md)**: Explore the design, roadmap, and command specifications for the future production PromptEngine CLI.

---

## Core Philosophy

PromptEngine operates on the principle that **the codebase and its version-controlled documentation are the single source of truth**. 

Traditional chat histories are volatile and ephemeral. By storing project specifications, database structures, API routes, and engineering rules directly inside the repository, we ensure that:
- You can transition between different AI tools or LLM models without losing context.
- Onboarding new developers (human or AI) takes minutes instead of days.
- AI hallucinations are minimized because the model's instructions are grounded in concrete repository documentation.
