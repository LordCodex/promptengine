---
document_id: project-documentation-standard
title: Project Documentation Standard
ecosystem: cross-cutting
dependencies:
  - project-readme
  - project-bootstrap-standard
audience: [human, agent]
last_reviewed: 2026-08-03
---

# Project Documentation Standard

## Playbook Metadata
- **Purpose**: Establishes universal standards for creating, reviewing, updating, and maintaining project-specific documentation throughout a project's lifecycle.
- **Scope**: All project specifications, API designs, databases schemas, business calculations, and deployment scripts.
- **When to Read**: When preparing feature designs, executing pull requests, or updating codebase documentation.
- **Related Playbooks**: [Project Overview](README.md), [Project Bootstrap Standard](01-project-bootstrap-standard.md), [Documentation Engineering Standard](../core/19-documentation-engineering-standard.md).
- **Canonical Source**: This is the canonical document for project documentation rules, update workflows, and context retention hierarchies.
- **Version**: 1.0.0
- **Last Reviewed**: 2026-08-03

---

## 1. Purpose

The **Project Documentation Standard** defines the rules governing how project knowledge is captured and maintained. It exists to guarantee that the context of a software application survives:
- **AI Memory Expiration**: Session limits, chat splits, and token clears.
- **AI Model Switches**: Porting codebases between different LLMs or coding agents.
- **Developer Turnover**: Human handovers and onboarding processes.
- **Time**: Restoring focus on historical parts of a system years after their implementation.

By committing and version-controlling documentation alongside the code, project knowledge is preserved in a permanent, searchable, and auditable format.

---

## 2. Documentation Philosophy

- **How to Build vs. What is Built**:
  - **Engineering Standards** (e.g. PromptEngine): Answer *"How should software be built?"* (security, naming, test coverage, code architecture).
  - **Project Documentation** (e.g. `docs/`): Answers *"What is this software?"* (business rules, schemas, routes, domain calculations).
- **Living Documentation**: Documentation must not be a static snapshot. It must evolve incrementally with the codebase. A feature implementation is incomplete until its corresponding documentation has been updated and reviewed.

---

## 3. Core Principles

- **Codebase Integration**: Project documentation is code. It resides inside the repository, is tracked in Git, and changes with the code.
- **Separation of Concerns**: Documentation explains **intent**, **decisions**, and **business rules**. It should not duplicate source code (e.g., repeating a PHP function line-by-line).
- **Auditability**: Changes to requirements or architectural logic must be visible in the git history of the markdown documents.
- **Model Agnosticism**: Written clearly and concisely, structured so both human developers and LLMs can extract facts instantly without navigating conversational clutter.

---

## 4. Source of Truth Hierarchy

When code implementation, documentation, and agent assumptions conflict, the AI must follow this strict priority hierarchy:

```text
1. Explicit Developer Instructions (Highest Priority)
2. Project-Specific AI Constitution (AGENTS.md)
3. Approved Project Documentation (docs/ folder)
4. Repository Implementation (Source code)
5. Infrastructure Configuration (Dockerfiles, CI files)
6. Automated Tests (Unit and integration suites)
7. Historical Decisions (Committed ADR logs)
8. AI Inference (Codebase scans and logs analysis) (Lowest Priority)
```

### Conflict Resolution Heuristics
- If the **source code** behaves differently than the **approved documentation**, the AI must flag this discrepancy immediately before making changes.
- If **developer instructions** contradict committed documents, the developer instructions win. The AI should offer to update the documentation to reflect the new direction.
- The AI must **never** overwrite a higher-priority source based on lower-priority inferences.

---

## 5. Documentation Lifecycle

Every documentation file moves through six sequential stages:

```text
[Creation] ──> [Review] ──> [Approval] ──> [Maintenance] ──> [Archive] ──> [Replacement]
```

1. **Creation**: Drafted during project bootstrap or when designing a major new feature (RFC/design document).
2. **Review**: Checked by peers (humans and AI) for accuracy, format, and consistency.
3. **Approval**: Merged into the `main` branch, designating it as the active authority.
4. **Maintenance**: Updated incrementally as features are refined or hotfixed.
5. **Archive**: Marked as historical when features are deprecated (preserving ADR logs).
6. **Replacement**: Replaced by newer standards when a system is modernized or refactored.

---

## 6. Required Project Documentation

Every repository must maintain the project-specific AI Constitution in the project root, along with the ten core documents under the local `docs/` or `.agents/` folder:

### Root Onboarding File
| Document | Purpose | Owner | Update Frequency | Primary Audience |
| :--- | :--- | :--- | :--- | :--- |
| **AGENTS.md** | Project AI Constitution (tech stack, rules, constraints) | Lead Architect | On stack/exception change | AI Agents, Developers |

### Core Project Specs
| Document | Purpose | Owner | Update Frequency | Primary Audience |
| :--- | :--- | :--- | :--- | :--- |
| **PRD** | Product requirements, user stories, features scope | Product Manager / Lead | Per feature cycle | Developers, AI |
| **Architecture** | System boundaries, services map, trust domains | Architect / Lead Dev | On structural change | Developers, SREs |
| **Business Rules** | Domain rules, formulas, status state machines | Business Analyst / Dev | On business logic edit | Developers, QA, AI |
| **Database** | Active schemas, keys, indices, relations | Database Lead / Dev | On migration write | Developers, AI |
| **API** | Endpoints, payloads, validations, response models | Backend Lead | On endpoint edit | Frontend, Mobile, AI |
| **Progress** | Active tasks list, bug logs, sprint tasks | Dev Team | Daily | Team, AI |
| **Roadmap** | Planned releases, technical debt, deprecations | Product Lead / Architect | Per quarter | Stakeholders, Team |
| **Decisions** | Architecture Decision Records (ADR) logs | Architect / Lead Dev | On major design choice | Developers, AI |
| **Deployment** | CI/CD parameters, env variables, rollbacks | DevOps / SRE | On pipeline adjustment | SREs, AI |
| **Troubleshooting** | Known error codes, diagnostic checklists | SRE / Support Lead | On resolving incident | SREs, Team, AI |

---

## 7. Documentation Update Rules

When modifying codebase implementations, the AI must follow these execution rules:

- **Determine Impact**: Before writing code, analyze whether the change affects requirements, endpoints, database schemas, or business formulas.
- **Section Isolation**: Update only the specific sections of the document that are affected by the code modification. Do not rewrite unrelated text.
- **No Silent Logic Changes**: Never silently modify a business rule or database relationship. If an adjustment is required, state it clearly in the planning stage.
- **Explain the "Why"**: Include a short changelog or note explaining the reason for the update (referencing a ticket or PR ID).

---

## 8. Documentation Quality Rules

Documentation must satisfy these ten criteria:
1. **Accurate**: Must match the exact state of the source code.
2. **Current**: Kept up to date; stale content must be corrected or archived.
3. **Reviewable**: Staged in markdown format to allow clear git diff reviews.
4. **Traceable**: References original PRDs, design tickets, or ADR choices.
5. **Searchable**: Uses clean headers and semantic tags for quick indexing.
6. **Consistent**: Uses the same domain terms across code, database, and docs.
7. **Non-redundant**: Points to a single canonical source for each topic.
8. **Developer-Friendly**: Uses code blocks, tables, and lists instead of prose.
9. **AI-Friendly**: Free of conversational fluff; presents structured facts.
10. **Future-Proof**: Employs relative paths rather than fragile absolute URLs.

---

## 9. Responsibilities

### AI Responsibilities
- **Context Scan**: Read all local documentation before planning any codebase changes.
- **Detect Staleness**: If code scans reveal that an endpoint or column has changed but the documentation is stale, alert the developer.
- **Incremental Staging**: Draft documentation updates alongside code updates in the same workspace transaction.
- **Preserve History**: Do not alter completed ADR records; add new records to document subsequent changes.

### Developer Responsibilities
- **Review Generated Docs**: Review documentation drafts created by the AI as carefully as the source code.
- **Authoritative Decisions**: Make all final calls on architecture and requirements; reject AI assumptions that drift from requirements.
- **Commit History**: Enforce documentation updates as a mandatory requirement for PR merges.

---

## 10. Documentation Anti-Patterns

- **Source Code Dumping**: Copying large blocks of source code into markdown files. (Anti-pattern: Code changes will instantly make documentation obsolete.)
- **Speculative Requirements**: Writing documentation for features that haven't been approved yet.
- **Tribal Knowledge**: Keeping business-critical formulas or security tokens undocumented in a developer's local environment.
- **Unreviewed AI Generates**: Allowing the AI to write extensive documentation without human validation.

---

## 11. Repository Rules

- The project root must contain the auto-generated `AGENTS.md` file.
- All other core project context files must reside in the `docs/` folder (or `.agents/` folder).
- Documentation files participate in pull requests.
- A pull request containing database changes or API updates **must** include corresponding documentation changes in the same commit.

---

## 12. AI Documentation Checklist

Before marking an engineering task as complete, the AI must verify:

```text
Did any of the following change?
[ ] Project Constitution / Tech Stack (AGENTS.md)
[ ] Product Requirements (PRD)
[ ] Architecture Layout
[ ] Business Calculations / Rules
[ ] Database Schemas / Indexes
[ ] API Routes / Payloads
[ ] Deployment Configurations
[ ] Architecture Decisions (ADR)
[ ] Sprint Progress Tasks
```

If **yes** to any of the above, update the corresponding documentation files.
If **no** changes occurred, output:
> "No documentation changes required."

---

## 13. Deliverables

Successful documentation management produces:
- **Updated Markdown Docs**: Precise, version-controlled documentation reflecting the exact state of the codebase.
- **Audit Logs**: Git commit history linking documentation updates to implementation code changes.
- **Active Task Progress**: Updated sprint/progress checklists tracking implementation status.
