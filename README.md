# AI Engineering Playbook

This repository serves as the single source of truth for software engineering decisions, architectural patterns, and coding standards. It is designed to act as an **Engineering Operating System** for both human developers and AI coding agents.

---

## Core Objectives
1. **Mental Modeling**: Force analysis and planning *before* code generation.
2. **Ecosystem-First Coherence**: Group architectural standards by language and framework hierarchies.
3. **Safe Legacy Modernization**: Standards for safe refactoring and risk reduction over a 10-year system lifecycle.
4. **Token-Optimized Discoverability**: A directory structure and indexing manifest designed to minimize agent context bloat.

---

## Playbook Structure Map

```text
ai-engineering-playbook/
├── README.md                            # Main index and onboarding guide (This file)
├── CONTRIBUTING.md                      # Guide to contributing new rules or modifying standards
├── LICENSE                              # Open-source licensing
├── playbook-manifest.json              # Machine-readable indexing file for AI agent directory routing
├── core/                                # Language-agnostic engineering principles
│   ├── 00-engineering-philosophy.md     # Core Values: Correctness, simplicity, trade-offs matrix
│   ├── 01-thinking-and-planning.md      # The Thinking Loop: Dissecting requirements
│   ├── 02-architecture-and-simplicity.md # Decoupling layers, cohesion, vertical slices, overengineering
│   ├── 03-data-and-api-modeling.md      # Schema patterns, REST and OpenAPI
│   ├── 04-testing-philosophy.md         # Testing hierarchies: pyramid, unit, integration
│   └── 05-universal-coding-standards.md # Standard naming, readability, class/method interfaces
├── environment/                         # Runtime environments, tools, and local configurations
│   ├── 01-local-dev-standards.md        # Containerization, runtime setups, host config
│   ├── 02-dependency-hygiene.md         # Auditing, ranges, and package integrity
│   └── 03-ci-cd-pipelines.md            # automated validations, validation check runners
├── legacy/                              # Incremental refactoring and system modernisation
│   ├── 01-safe-refactoring.md           # Strangler pattern, Adapters, Branch-by-Abstraction
│   ├── 02-backward-compatibility.md     # Expand-and-contract database migrations
│   └── 03-deployment-risk-reduction.md  # Canary verifications, rollback routines
├── security/                            # [Placeholder] Security Engineering Guidelines
│   ├── README.md                        # Overview of secure coding practices
│   └── ...                              # Details on authorization, encryption, and secrets
├── performance/                         # [Placeholder] Performance Engineering Guidelines
│   ├── README.md                        # Overview of performance testing
│   └── ...                              # Details on caching, database indexes, and concurrency
├── stacks/                              # Ecosystem-first language and framework conventions
│   ├── php-laravel/                     # PHP & Laravel Ecosystem
│   │   ├── php-conventions.md           # Strict types, error/exception boundaries, PSRs
│   │   ├── laravel-routing.md           # Routing layout, middleware, form request validations
│   │   ├── laravel-data.md              # Eloquent schemas, model scopes, relations
│   │   ├── laravel-logic.md             # Service classes, Actions pattern, queue dispatches
│   │   └── laravel-testing.md           # Feature tests, DB transactional states, Pest
│   ├── dart-flutter/                    # Dart & Flutter Ecosystem
│   │   ├── dart-conventions.md          # sound null safety, async/await, streams
│   │   ├── flutter-widgets.md           # Widget lifecycles, build context optimizations
│   │   ├── flutter-state.md             # State Management patterns (BLoC, Riverpod)
│   │   ├── flutter-architecture.md      # Layer hierarchies, GoRouter navigation, native modules
│   │   └── flutter-testing.md           # Widget/integration testing, golden tests
│   └── js-ts-vue-nuxt/                  # JS/TS & Vue Ecosystem
│       ├── js-ts-conventions.md         # ESM structure, TS typing boundaries, async boundaries
│       ├── vue-components.md            # Script Setup, Composition API, Pinia state integration
│       └── nuxt-ssr.md                  # SSR hydration lifecycle, Server-only routes
├── bridges/                             # Inter-stack integration architectures (The "Glue")
│   ├── laravel-inertia-vue.md           # Inertia configuration, state share middlewares
│   ├── laravel-api-flutter.md           # REST API conventions, contract testing, JWT lifecycle
│   └── cross-origin-sharing.md          # CORS policies, network routing rules
├── examples/                            # Validated, runnable example projects (Sandboxes)
│   ├── php-laravel-sandbox/             # Minimal Laravel execution reference
│   ├── dart-flutter-sandbox/            # Minimal Flutter widget/logic testing reference
│   └── js-vue-sandbox/                  # Minimal Vue 3 Composition layout reference
├── project/                             # Project Knowledge System
│   ├── README.md                        # Overview of project documentation categories
│   └── 01-project-bootstrap-standard.md # Workflows to adopt and initialize codebases
└── ai/                                  # AI Agent System Instructions
    ├── agent-prompts.md                 # Markdown instructions to inject into external AI models
    └── agent-rules.json                 # Path-to-rule mapping configuration for automated agent loading
```

---

## Instructions for Human Developers

1. **Philosophy First**: Read [00-engineering-philosophy.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/00-engineering-philosophy.md) to understand the core trade-off valuation framework.
2. **Plan First**: Before making any major architecture modification or writing code, read [01-thinking-and-planning.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/01-thinking-and-planning.md). Follow the "Thinking Loop" methodology.
3. **Universal Code Standards**: Conhabit with the global, framework-agnostic rules defined in [05-universal-coding-standards.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/05-universal-coding-standards.md). Class separation, naming, and error guidelines defined there are inherited globally.
4. **Context-Specific Reference**: When writing code in a specific language ecosystem (e.g. Flutter), refer to the respective stack guide: [dart-conventions.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/dart-conventions.md) and [flutter-widgets.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/dart-flutter/flutter-widgets.md). Do not duplicate core conventions; reference them.
5. **Backward Compatibility**: If modifying legacy code, consult [01-safe-refactoring.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/01-safe-refactoring.md) before writing implementations.

---

## Instructions for AI Coding Agents

1. **Pre-Flight Indexing**: Locate the `playbook-manifest.json` file in the root. Query it using the keywords of your assigned task (e.g., "Laravel routing", "Inertia Vue sharing") to find the minimal required rule files.
2. **Minimal Context Loading**: Load only the files specified by the manifest search. Avoid reading the entire repository to conserve context tokens.
3. **Execution Plan**: Always write a detailed implementation plan in your thinking block or scratchpad based on [01-thinking-and-planning.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/01-thinking-and-planning.md) before executing filesystem changes.

---

# Using PromptEngine in Your Projects

PromptEngine acts as a centralized engineering knowledge base that integrates directly into your workspace. By utilizing our structured playbooks, templates, workflows, and prompts, you can establish an automated, high-quality development ecosystem.

---

## Quick Start Navigation

- **Developer Guides**: Explore step-by-step instructions in the [Guides Index](guides/README.md).
  - [Getting Started](guides/01-getting-started.md) | [New Project Setup](guides/02-new-project.md) | [Existing Project Onboarding](guides/03-existing-project.md) | [Legacy Migration](guides/04-migrate-existing-project.md) | [Daily Development](guides/05-daily-development.md) | [Troubleshooting](guides/09-troubleshooting.md)
- **AI Prompts Library**: Access copy-and-paste ready prompts in the [Prompts Index](prompts/README.md).
  - [Bootstrap New Project](prompts/01-new-project.md) | [Scan Existing Project](prompts/02-existing-project.md) | [Migrate Legacy Systems](prompts/03-migrate-existing-project.md) | [Add Feature](prompts/04-add-feature.md) | [Bug Fix](prompts/05-bug-fix.md) | [Safe Refactoring](prompts/06-refactor.md) | [Reviews](prompts/09-project-review.md)
- **CLI Foundation Specs**: Explore the design, roadmap, and specifications for the future production CLI in the [CLI Foundation Index](cli/README.md).
  - [CLI Spec](cli/CLI-Spec.md) | [Architecture](cli/CLI-Architecture.md) | [Configuration](cli/CLI-Configuration.md) | [Command Reference](cli/CLI-Command-Reference.md) | [Roadmap](cli/CLI-Roadmap.md)

---

## 1. Creating a New Project

When initiating a greenfield project:
1. Create a root directory for your project. Clone the PromptEngine repository inside it (or in a shared parent directory).
2. Start your AI assistant (e.g., Cursor, Claude Code) and execute the **[New Project Bootstrap Prompt](prompts/01-new-project.md)**.
3. The AI will conduct a brief requirements interview and automatically generate:
   - **`AGENTS.md`** in your project root (your project's AI Constitution).
   - **`docs/` Folder** containing the 10 core documentation specifications.
4. For detailed guidelines, read the **[New Project Setup Guide](guides/02-new-project.md)**.

---

## 2. Migrating an Existing Project

If you want to onboard an existing active codebase:
1. Clone PromptEngine nested or shared in your workspace.
2. Start your AI assistant and execute the **[Existing Project Bootstrap Prompt](prompts/02-existing-project.md)** (or the **[Migration Prompt](prompts/03-migrate-existing-project.md)** if you have pre-existing custom instruction files).
3. The AI will scan package files and source code directories on disk to reverse-engineer your database schema, routing boundaries, and stack conventions.
4. The AI will automatically generate your root `AGENTS.md` constitution and the 10 specifications in the `docs/` folder.
5. For detailed steps, consult the **[Existing Project Onboarding Guide](guides/03-existing-project.md)** and the **[Legacy Migration Guide](guides/04-migrate-existing-project.md)**.

---

## 3. Daily Development Workflow

For daily development tasks, always plan and coordinate with the AI using our structured workflows:
- **Adding Features**: Use the **[Add Feature Prompt](prompts/04-add-feature.md)** to check `docs/PRD.md` and write an implementation plan. Follow the [Feature Workflow](workflows/01-feature-implementation.md).
- **Fixing Bugs**: Use the **[Bug Fix Prompt](prompts/05-bug-fix.md)** to isolate issues and write reproduction tests.
- **Refactoring Code**: Use the **[Safe Refactoring Prompt](prompts/06-refactor.md)** to simplify logic step-by-step while preserving behaviour. Follow the [Modernization Workflow](workflows/02-legacy-modernization.md).
- **Architecture Updates**: Use the **[Architecture Change Prompt](prompts/07-architecture-change.md)** to log decisions in `docs/Decisions.md` (ADR log).
- **Reviewing Code**: Use the **[Project Review Prompt](prompts/09-project-review.md)** to audit security, performance, and accessibility.
- See the **[Daily Development Guide](guides/05-daily-development.md)** for a complete walkthrough.

---

## 4. The Project Knowledge System

PromptEngine enforces a structured documentation model consisting of:
- **`AGENTS.md`**: Placed in the project root. It serves as the project's AI Constitution. It is automatically generated and contains Section 1 (PromptEngine Core Rules) and Section 2 (Project Constitution detailing tech stack, auth, databases, constraints, and exceptions).
- **`/docs/` Specifications**: The 10 core markdown documents (`PRD.md`, `Architecture.md`, `Database.md`, `API.md`, `Progress.md`, `Roadmap.md`, `Decisions.md`, `Deployment.md`, `Troubleshooting.md`).
- Read the **[Managing Project Documentation Guide](guides/06-project-documentation.md)** to understand how to keep specs synchronized.

---

## 5. PromptEngine Workflows & Best Practices

To prevent AI context loss, minimize token consumption, and resolve ignoring behaviors:
- Always enforce the **5-Step Entry Rule**: The AI must read `AGENTS.md` first, check `docs/` specs, load standards playbooks, select the correct workflow, and plan before coding.
- Keep chat threads short. Start a fresh chat thread once a task is completed.
- Learn token optimization tips and integration instructions for Cursor, Claude Code, Windsurf, ChatGPT, and Codex in the **[Developer Best Practices Guide](guides/08-best-practices.md)**.
- For unresolved issues, view the **[Troubleshooting Guide](guides/09-troubleshooting.md)**.

