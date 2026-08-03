---
document_id: project-bootstrap-standard
title: Project Bootstrap Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-thinking-and-planning
  - project-readme
audience: [human, agent]
last_reviewed: 2026-08-03
---

# Project Bootstrap Standard

## Playbook Metadata
- **Purpose**: Establishes standards and step-by-step workflows for AI coding agents to initialize and align their context with new or existing projects.
- **Scope**: Codebase audits, technology stack identification, initial documentation generation, and developer onboarding.
- **When to Read**: Immediately on starting a new development thread or adopting an unfamiliar repository.
- **Related Playbooks**: [Project Overview](README.md), [Thinking and Planning Before Coding](../core/01-thinking-and-planning.md), [AI Agent Engineering Workflow](../core/20-ai-agent-engineering-workflow-standard.md).
- **Canonical Source**: This is the canonical document for codebase bootstrapping and project onboarding rules.
- **Version**: 1.0.0
- **Last Reviewed**: 2026-08-03

---

## 1. Bootstrapping Philosophy

- **Repository as Source of Truth**: The codebase and committed documentation represent the absolute state of the project. Do not trust external chat logs or temporary session memory over repository contents.
- **Temporary Sessions**: Chat history is ephemeral. Your understanding of a project must be fully reconstructible from the files inside the repository.
- **Evidence-Based Engineering**: Always verify code rules, database schemas, and endpoints by inspecting actual code files. Never rely on speculative assumptions or memory.
- **Documentation Preservation**: If context is missing, record it in documentation immediately so future models and sessions can adopt the project seamlessly.

---

## 2. New Project Bootstrap

When starting a project from scratch (zero codebase lines):

### Step 1: Stack & Ecosystem Identification
- Review user requirements and define the target technology stack (backend, database, frontend, mobile layout frameworks).
- Enforce the stack conventions defined in the PromptEngine `/stacks/` folder.

### Step 2: Goal and Boundary Definition
- Outline the primary business objectives, user flows, and core features.
- Define what the project **will not do** to establish a clear architectural scope boundary.

### Step 3: Document Initialization
- Generate the project's AI Constitution (`AGENTS.md`) in the project root by filling the PromptEngine `AGENTS.template.md` with discovered stack, constraints, and overview.
- Generate the initial project context documents under a local `/docs/` folder (such as `docs/PRD.md`, `docs/Architecture.md`, and `docs/Database.md`).
- Present these documents to the developer for review and confirmation.

### Step 4: Avoid Speculation
- If a requirement is ambiguous or underspecified, ask the developer for clarification instead of guessing or writing placeholder rules.

---

## 3. Existing Project Bootstrap

When adopting an existing repository (already contains code):

### Step 1: Codebase Scanning
Perform a systematic scan of the repository structure to identify:
- **Folders Map**: Identify directories (`app/`, `src/`, `tests/`, `docs/`).
- **Dependency Manifests**: Read package files (`package.json`, `composer.json`, `pubspec.yaml`) to capture technologies, frameworks, and locked versions.
- **Coding Conventions**: Analyze existing codebase files to detect naming styles (e.g. PSR, camelCase, snake_case) and architectural layout structures.

### Step 2: Extract Business Rules and Workflows
- Scan the source code to locate core business calculations, state transitions, and integration points.
- Map the primary data models, database relationships, and API endpoints.

### Step 3: Identify Documentation & Constitution Gaps
- Compare the existing code implementation against the `/docs/` folder and `AGENTS.md`.
- Highlight missing schemas, undocumented APIs, obsolete requirements, or outdated stack definitions.

### Step 4: Auto-Generate or Update AGENTS.md
- If `AGENTS.md` is missing or stale, generate it in the project root by filling the `AGENTS.template.md` template, reverse-engineering the tech stack, database, and constraints from codebase scanning.

### Step 5: Infer vs. Confirm Knowledge
- **Inferred Knowledge**: Code behavior you reverse-engineered from scanning files.
- **Confirmed Knowledge**: Rules explicitly stated in documentation, requirement files, or user instructions.
- **Rule**: Clearly distinguish assumptions from confirmed facts in your planning reports.

---

## 4. Documentation & Constitution Initialization Lifecycle

- **Create Once**: Generate `AGENTS.md` in the project root and the core documentation files under `docs/` (`PRD.md`, `Architecture.md`, `Database.md`, `API.md`) on startup. The developer never manually writes `AGENTS.md` but audits it.
- **Human Audit**: Present the initialized constitution and documentation to human developers for manual review and staging confirmation.
- **Incremental Maintenance**: Do not let documentation or the project constitution go stale. Update `AGENTS.md` and API/database files incrementally during active feature pull requests when architectural paths or stacks evolve.

---

## 5. Repository Analysis Workflow (The 5-Step Entry Rule)

AI coding assistants must execute this 5-step entry workflow on boot:

```text
[1. Read AGENTS.md] ──> [2. Read docs/ Specs] ──> [3. Load Playbooks]
                                                        │
   [5. Understand & Plan] <── [4. Audit Code & Stack] ──┘
```

1. **Step 1 — Read AGENTS.md First**: Always read the `AGENTS.md` in the project root to load the project's AI Constitution (tech stack, exception logs, constraints). If starting a new bootstrap, read PromptEngine's own bootstrap files and manifest.
2. **Step 2 — Read docs/ Specs Second**: Read the relevant project requirements, architecture design, and database schema files under `docs/` (or `.agents/`).
3. **Step 3 — Load Playbooks Third**: Query the `playbook-manifest.json` and load the appropriate PromptEngine standard playbooks matching the stack.
4. **Step 4 — Audit Code & Stack Fourth**: Scan directories and dependency manifests to verify implementation aligns with `AGENTS.md` and `docs/`.
5. **Step 5 — Understand & Plan Fifth**: Build a complete mental model of the system. Propose architectural exceptions or updates to the Project Constitution if gaps are found, then draft the implementation plan.

---

## 6. Strict Execution Rules

- **Never Invent Business Rules**: If a business logic calculation (e.g., tax computation or discount threshold) is not documented or written in code, do not invent placeholder variables. Ask the developer.
- **Never Guess Relationships**: Do not guess database columns, foreign keys, or API endpoint payload schemes. Verify them in code or prompt for them.
- **Distinguish Assumptions**: Always declare assumptions in your thinking blocks using:
  > **[ASSUMPTION]**: We assume $X$ because of file $Y$.
- **No Conversation Memory Reliance**: Never say "as we discussed in the previous session." If an agreement was reached, commit it to a local markdown file in the repository.
- **Prefer Repository Evidence**: Treat the codebase as primary evidence. If a document conflicts with what is compiled on disk, raise a contradiction warning immediately.

---

## 7. Required Onboarding Deliverables

A successful project bootstrap process must output a structured summary containing:

1. **Project Understanding**: A brief sentence explaining the target app's purpose.
2. **Identified Stack**: Language and framework versions locked in packages.
3. **Identified Architecture**: Design patterns detected (e.g., Action Pattern, decoupled Next.js API layers).
4. **Documentation Gaps**: List of missing schemas or requirements.
5. **Risks & Concerns**: Security flaws, concurrency issues, or legacy blockages discovered during codebase scanning.
6. **Active Assumptions**: Assumptions made during the scan.
7. **Clarification Questions**: Targeted, high-value questions for the developer.
