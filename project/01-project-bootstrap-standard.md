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
- Generate the initial project context templates under a local `/docs/` folder (such as `docs/PRD.md`, `docs/Architecture.md`, and `docs/Database.md`).
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

### Step 3: Identify Documentation Gaps
- Compare the existing code implementation against the `/docs/` folder.
- Highlight missing schemas, undocumented APIs, or obsolete requirements.

### Step 4: Infer vs. Confirm Knowledge
- **Inferred Knowledge**: Code behavior you reverse-engineered from scanning files.
- **Confirmed Knowledge**: Rules explicitly stated in documentation, requirement files, or user instructions.
- **Rule**: Clearly distinguish assumptions from confirmed facts in your planning reports.

---

## 4. Documentation Initialization Lifecycle

- **Create Once**: Generate the core documentation files (`PRD.md`, `Architecture.md`, `Database.md`, `API.md`) on startup.
- **Human Audit**: Present the initialized documentation to human developers for manual review and staging confirmation.
- **Incremental Maintenance**: Do not let documentation go stale. Update API files and business rules incrementally during active feature pull requests.

---

## 5. Repository Analysis Workflow

AI coding assistants must execute this 8-step workflow on boot:

```text
[1. Read AGENTS.md] ──> [2. Read Bootstrap] ──> [3. Read Local Docs]
                                                      │
 [6. Ask Developer] <──  [5. Audit Code]    <── [4. Identify Stack] ┘
         │
 [7. Stage Docs]    ──> [8. Begin Plan]
```

1. **Step 1 — Read AGENTS.md**: Read the onboarding instruction file in the project root.
2. **Step 2 — Read PromptEngine Bootstrap**: Read `promptengine/ai/bootstrap.md` to load rules and manifest priorities.
3. **Step 3 — Read Local Docs**: Scan the local `docs/` or `.agents/` folder for product requirements, databases, and API schemas.
4. **Step 4 — Identify Stack**: Identify dependencies and framework targets from lock files.
5. **Step 5 — Audit Code**: Scan directories to confirm database models, routes, and controllers align with documentation.
6. **Step 6 — Detect Gaps**: Highlight missing information or contradictions.
7. **Step 7 — Ask & Verify**: Compile assumptions and ask the developer to confirm design decisions.
8. **Step 8 — Begin Planning**: Draft your implementation plan inside the conversation log.

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
