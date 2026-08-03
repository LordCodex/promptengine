---
document_id: core-reusable-ai-prompt-template-library
title: Reusable AI Prompt Template Library
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-testing-engineering-standard
  - core-git-and-collaboration-standard
  - core-cicd-and-deployment-standard
  - core-infrastructure-and-devops-standard
  - core-legacy-modernization-and-refactoring-standard
  - core-refactoring-standards-and-safe-migration-workflow
  - core-code-review-engineering-standard
  - core-documentation-engineering-standard
  - core-ai-agent-engineering-workflow-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Reusable AI Prompt Template Library

## Purpose & Inheritance
This library defines the standard reusable templates for prompting AI coding agents and assistants (such as Codex, Claude Code, Cursor, and Gemini). It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md) and the [AI Agent Engineering Workflow Standard](20-ai-agent-engineering-workflow-standard.md). These templates enforce consistent quality across project analysis, feature additions, bug fixes, refactoring, security reviews, and deployment prep tasks.

---

## 1. Prompt Design Principles

To maximize AI code quality, all custom prompts must follow this five-part structure:
- **Context**: Define the current state, tech stack versions, framework patterns, and configurations.
- **Goal**: Define the exact problem to be solved.
- **Constraints**: Set bounds (e.g. no new external packages, lock specific files, preserve backward compatibility, performance targets).
- **Expected Output**: Define the required files, code structures, diff formats, and explanations.
- **Validation**: Define success checks (e.g. compile requirements, specific tests execution, latency profile).

---

## 2. Prompts Library Integration

To prevent rule duplication and maintain a single source of truth, all reusable AI prompts have been consolidated into the top-level **[prompts/](../prompts/README.md)** directory. 

Developers and AI agents must reference the specific prompt files directly:

### Project Initialization & Setup
- **[Bootstrap New Project](../prompts/01-new-project.md)**: Prompt for conducting requirements discovery and initializing greenfield codebases.
- **[Scan Existing Project](../prompts/02-existing-project.md)**: Prompt for scanning, reverse-engineering, and onboarding existing repositories.
- **[Migrate Legacy Systems](../prompts/03-migrate-existing-project.md)**: Prompt for legacy migrations and consolidating rules.

### Daily Development
- **[Add Feature](../prompts/04-add-feature.md)**: Prompt for feature implementation planning and coding.
- **[Bug Fix](../prompts/05-bug-fix.md)**: Prompt for isolating bugs and generating minimal targeted patches.
- **[Safe Refactoring](../prompts/06-refactor.md)**: Prompt for safe step-by-step refactoring.
- **[Architecture Change](../prompts/07-architecture-change.md)**: Prompt for database schema updates and module layers changes.
- **[Documentation Sync](../prompts/08-documentation-sync.md)**: Prompt for aligning markdown specifications with the live code.

### Reviews & Audits
- **[Project Code Review](../prompts/09-project-review.md)**: Prompt for security, performance, and UI reviews.
- **[Release Readiness](../prompts/10-release.md)**: Prompt for production deployment audits and release notes.

---

## 3. Prompt Variable Placeholders

When executing prompt templates, inject parameters using this standardized syntax:
- `Technology`: Language and framework version (e.g. PHP 8.3, Laravel 11.x).
- `Path`: Absolute file system path (e.g. `[RegisterUserAction.php](file:///app/Actions/RegisterUserAction.php)`).
- `Task`: Concise goal description.
- `Constraints`: Strict architectural rules (e.g., "no raw database queries").
- `Validation`: Verification steps (e.g., "run `composer test`").

---

## 4. Prompt Usage & Context Efficiency Guide

To optimize context window usage and avoid AI response decay:
- **Use Task-Focused Micro-Prompts**: Do not compile all prompts into a single master request. Choose the smallest relevant prompt for the task.
- **Provide Targeted Context Files**: Pass only the specific files relevant to the prompt target.
- **Do Not Repost Baseline Playbooks**: AI agents should fetch general rules from local playbook references rather than having them copied directly into every prompt message.

---

## References
- Simplicity Standards: [02-architecture-and-simplicity.md](02-architecture-and-simplicity.md)
- Universal Naming Standards: [05-universal-coding-standards.md](05-universal-coding-standards.md)
- Testing & QA Harnesses: [11-testing-engineering-standard.md](11-testing-engineering-standard.md)
- AI agent workflow parameters: [20-ai-agent-engineering-workflow-standard.md](20-ai-agent-engineering-workflow-standard.md)
- User Guides: [guides/README.md](../guides/README.md)
