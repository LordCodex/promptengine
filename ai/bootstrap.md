---
document_id: ai-bootstrap
title: AI Agent Bootstrap Guide
ecosystem: cross-cutting
audience: [agent]
last_reviewed: 2026-08-03
---

# PromptEngine AI Agent Bootstrap Guide

This document is the mandatory entry point and bootstrapping configuration for every AI coding agent, autonomous developer, or external AI assistant interacting with this repository. Read this document before executing any file reads or codebase modifications.

---

## 1. Purpose of PromptEngine

PromptEngine is a structured, machine-readable software engineering knowledge base. It functions as an **Engineering Operating System** for the repository, providing:
- Reusable, language-agnostic core engineering standards.
- Stack-specific coding conventions (Laravel/PHP, Flutter/Dart, Vue/Nuxt, React/Next.js).
- Verification checklists, procedural workflows, and architectural decision trees.

This framework is design-compatible with any LLM, programming assistant, or autonomous agent loop.

---

## 2. AI Operating Principles

Every AI coding assistant must adhere to these execution boundaries:
1. **Read This Document First**: Establish context constraints and priority boundaries.
2. **Read the Indexing Manifest Second**: Query [playbook-manifest.json](../playbook-manifest.json) using keywords matching the active task.
3. **Load Minimal Context**: Load only the specific playbooks and conventions relevant to the active task to conserve token footprint. Do not load the entire repository.
4. **Follow Project-First Patterns**: Identify and reuse established architectural patterns inside the workspace before writing custom utility code.
5. **Enforce Constraints Continuously**: Keep rules and error envelopes active in memory throughout the editing lifecycle.

---

## 3. Task Classification

Before planning, classify the user request into one of these target categories to narrow down standard requirements:

- **New Feature**: Requires API/data modeling, testing planning, and progressive implementation.
- **Bug Fix**: Requires isolation, reproduction test cases, and a targeted patch.
- **Refactor**: Requires behavior preservation, complexity reduction, and test safety nets.
- **Performance Optimization**: Requires baseline profiling, query checks, or cache integrations.
- **Security Review**: Requires checking input boundaries, cookies, headers, or authorization guards.
- **UI Implementation**: Requires checking layouts, typography, responsiveness, and accessibility targets.
- **Documentation**: Requires updating ADR logs, README files, or API schema declarations.
- **Testing**: Requires adding unit, integration, golden, or E2E tests.
- **DevOps / CI-CD**: Requires modifying Dockerfiles, GitHub workflows, or rollback scripts.

---

## 4. Selecting Playbooks

Optimize context loading by selecting the minimum subset of files required for the task. 

- Match task classifications and keywords against the tags in [playbook-manifest.json](../playbook-manifest.json).
- **Core Rules**: Always load [05-universal-coding-standards.md](../core/05-universal-coding-standards.md) for standard conventions.
- **Tech Stack Rules**: Load the matching stack folder (e.g. `stacks/react-next/*` if working on Next.js).
- **Process Checklists**: Load the matching checklist (e.g. `checklists/02-security-review-checklist.md` for auth endpoint modifications).

---

## 5. Priority Order

When guidelines or instructions conflict, enforce the following order of precedence:
1. **User Instructions**: The explicit requests provided in the current chat prompt.
2. **Project-Specific Rules**: Project-level instructions defined in `AGENTS.md` or `.agents/` configurations.
3. **PromptEngine Stack Playbooks**: Language/framework-specific files located under `/stacks/`.
4. **PromptEngine Core Standards**: Agnostic engineering principles located under `/core/`.
5. **Language Conventions**: Standard coding guidelines (e.g. PSR-12, Effective Dart, PEP 8).
6. **General Best Practices**: General engineering heuristics (e.g. DRY, SOLID, clean code).

---

## 6. Working Process

Enforce the following workflow loop during execution:

```text
[Understand Request] ──> [Select Playbooks] ──> [Analyze Project Patterns]
                                                       │
   [Verify & Test]   <──  [Implement]   <── [Write Implementation Plan] ┘
```

1. **Understand**: Dissect inputs, outputs, constraints, and prerequisites.
2. **Select**: Query the manifest and load only the minimum rule files.
3. **Analyze**: Search the workspace for existing folders, classes, and tests matching the target pattern.
4. **Plan**: Write a concise step-by-step implementation plan (or edit task logs).
5. **Implement**: Write simple, typed, senior-quality code without debug lines or placeholder code.
6. **Verify**: Execute test suites and linters to confirm correctness.

---

## 7. Repository Usage

PromptEngine is modular. **Never read every document in the repository.** Instead, use tools (e.g. search, manifest query) to isolate and load target chapters. Conserving tokens improves response accuracy and minimizes logic errors.

---

## 8. Updating PromptEngine

When updating or adding rules to the PromptEngine repository:
- **Zero Duplication**: Do not repeat principles across multiple files. Add universal rules to `/core/` and refer to them from `/stacks/`.
- **Extend Existing Files**: Prefer adding sections to existing playbooks rather than creating new markdown files.
- **Valid Links Only**: All internal documentation links must use clean, relative paths (e.g. `[thinking.md](../core/01-thinking-and-planning.md)`). Never use absolute URL hosts (such as github.com) for local files.
- **Metadata Frontmatter**: Always include YAML metadata blocks at the top of new playbook documents.

---

## 9. General Principles

Always write code that prioritizes:
- **Simplicity**: Code should be boring and cheap to change. Reject unnecessary abstractions.
- **Correctness & Security**: Never trust inputs. Enforce authorization checks server-side on every request.
- **Performance**: Optimize database queries and prevent N+1 issues before applying caching wrappers.
- **Long-term Maintainability**: Avoid quick hacks. Leave code cleaner than you found it.
