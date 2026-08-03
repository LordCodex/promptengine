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

PromptEngine is designed to be reusable across multiple software projects. It acts as a centralized engineering knowledge base that you can integrate directly into your codebase or keep in a shared location.

## Quick Start

To use PromptEngine in your projects, follow this workflow:

1. **Add PromptEngine**: Add PromptEngine to your project repository (or keep it in a shared directory).
2. **Create AGENTS.md**: Create an `AGENTS.md` file in your project root.
3. **Point AGENTS.md to PromptEngine**: Write instructions in `AGENTS.md` pointing to your PromptEngine location.
4. **Start Your Assistant**: Start your AI coding assistant or agent.
5. **Read AGENTS.md**: The assistant reads `AGENTS.md` automatically on startup.
6. **Load PromptEngine**: The assistant loads PromptEngine from the path defined.
7. **Select Playbooks**: The assistant selects only the playbooks relevant to the current task using `playbook-manifest.json`.
8. **Begin Engineering**: Begin implementing your features with the assistant following the mapped playbooks.

### Example Project Structure (Nested)

When PromptEngine is nested directly within your project repository:

```text
Project/
├── AGENTS.md
├── app/
├── docs/
├── routes/
└── promptengine/      # PromptEngine repository clone
```

### Example Project Structure (Shared/External)

If you have multiple projects and want to keep a single, shared PromptEngine in a parent folder:

```text
Workspace/
├── PromptEngine/      # Shared PromptEngine repository
├── Project-A/
│   ├── AGENTS.md      # Points to ../PromptEngine/
│   └── ...
└── Project-B/
    ├── AGENTS.md      # Points to ../PromptEngine/
    └── ...
```

In this case, adjust the paths inside `AGENTS.md` to reference `../PromptEngine/`.

## Recommended Workflow

1. **Keep Global Standards Centralized**: Keep all reusable, language-agnostic coding standards and stack playbooks inside PromptEngine.
2. **Declare Local Rules Locally**: Put project-specific parameters (e.g. database schemas, business rules, product requirements) inside the project's `docs/` folder or `AGENTS.md`.
3. **Automate Agent Routing**: Let the AI assistant programmatically query `playbook-manifest.json` and load the minimum necessary documents for the current task.
4. **Update PromptEngine**: Rather than duplicating rules across multiple projects, update the centralized PromptEngine repository and pull changes to keep configurations synchronized.

