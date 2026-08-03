# PromptEngine Reusable AI Prompts Library

Welcome to the PromptEngine Prompts Library. This directory contains production-ready, AI-agnostic prompts designed to orchestrate software development workflows with coding assistants (such as Cursor, Windsurf, Claude, ChatGPT, and Codex).

---

## Prompt Index

1. **[01. New Project Bootstrap](01-new-project.md)**: Prompt for conducting requirements discovery and initializing greenfield codebases.
2. **[02. Existing Project Bootstrap](02-existing-project.md)**: Prompt for scanning, reverse-engineering, and onboarding existing repositories.
3. **[03. Migrate Existing Project](03-migrate-existing-project.md)**: Prompt for legacy migrations and consolidating rules.
4. **[04. Add Feature](04-add-feature.md)**: Prompt for feature implementation planning and coding.
5. **[05. Bug Fix](05-bug-fix.md)**: Prompt for isolating bugs and generating minimal targeted patches.
6. **[06. Refactor Code](06-refactor.md)**: Prompt for safe step-by-step refactoring.
7. **[07. Architecture Change](07-architecture-change.md)**: Prompt for schema updates and ADR index records.
8. **[08. Documentation Sync](08-documentation-sync.md)**: Prompt for aligning markdown specifications with the live code.
9. **[09. Project Code Review](09-project-review.md)**: Prompt for security, performance, and UI reviews.
10. **[10. Release Readiness](10-release.md)**: Prompt for production deployment audits and release notes.

---

## How to Configure Prompt Parameters

Prompts in this library contain bracketed placeholder variables (e.g. `{STACK}` or `{FEATURE_DESCRIPTION}`). 

Before copy-and-pasting a prompt into your AI tool, replace the placeholders with specific values representing your project constraints.

### Standard Parameters
- `{PROJECT_NAME}`: The name of the target application.
- `{STACK}`: Core languages and frameworks (e.g. Laravel 10, Vue 3, MySQL).
- `{PATH}`: The target directory or absolute file path to edit.
- `{FEATURE_DESCRIPTION}`: A clear sentence describing the new user flow.
- `{SYMPTOMS}`: The unexpected behavior or logs from a bug.
