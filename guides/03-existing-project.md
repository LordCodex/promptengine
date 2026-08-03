# 03. Active Existing Projects Guide

This guide explains how to adopt PromptEngine in an active, pre-existing codebase. It covers two scenarios: when the codebase has never used PromptEngine before, and when it already contains a generated constitution.

---

## 1. Scenario A: Core Project Has No AGENTS.md

If the repository does not contain an `AGENTS.md` file, you need to trigger the **[Existing Project Bootstrap Workflow](../workflows/ExistingProjectBootstrap.md)** to reverse-engineer the project configuration.

### Step 1: Clone PromptEngine in the Workspace
Add the PromptEngine directory clone nested inside your repository or in an external parent folder (see [New Project Guide](02-new-project.md)).

### Step 2: Run the Existing Project Bootstrap Prompt
Start your AI assistant and copy the **Existing Project Bootstrapper Prompt** from `prompts/02-existing-project.md`. This prompt commands the AI to perform a detailed workspace audit.

### Step 3: Workspace Scan & Reverse Engineering
The AI will scan your code on disk to map out the system boundaries:
1. **Dependency Analysis**: It inspects package manifests (`package.json`, `composer.json`, `pubspec.yaml`) to capture technologies, frameworks, and library versions.
2. **Directory Mapping**: It traces folder paths to identify layering patterns (e.g. controllers, services, actions).
3. **Database Audit**: It reads migration files and schemas to capture tables, relationships, and constraints.
4. **API Route Mapping**: It audits routes files to capture endpoints and query/payload schemas.

### Step 4: Auto-Generating AGENTS.md & Specs
Based on the code audit, the AI will generate:
- **`AGENTS.md`**: Fill Section 2 with the reverse-engineered parameters (tech stack, databases, constraints, coding style).
- **`docs/` Folder**: Create the 10 core documents mapping the *live codebase* state. The AI will populate `PRD.md` with active user flows, `Architecture.md` with current modules, and `Database.md` with database tables on disk.

---

## 2. Scenario B: Project Already Has AGENTS.md

If you or another developer has already bootstrapped the repository, `AGENTS.md` exists in the root. 

In this case:
1. **Automatic Alignment**: The AI coding assistant is configured to automatically read `AGENTS.md` as its first step.
2. **Ecosystem Matching**: It parses Section 2 to identify that the project is, for example, a *Laravel + Vue 3 Monolith*, and automatically loads matching stacks rules from PromptEngine.
3. **Task Coordination**: All instructions given in chat will be interpreted through the lens of the constraints and conventions logged in `AGENTS.md` and `docs/`.

---

## 3. Keeping Documentation Synchronized

The primary rule of PromptEngine is that **documentation is code**. 

When the AI implements a feature or refactors a subsystem:
- If a pull request modifies database columns, the AI must update `docs/Database.md` in the same commit.
- If a new route is added, the AI must document it in `docs/API.md`.
- If an architectural constraint is broken or bypassed (e.g. using auto-increment IDs for a temporary table), the AI must record this exception in `docs/Decisions.md` (or decisions log) and update `AGENTS.md` Section 2.5.
- The developer must enforce documentation synchronization during code reviews before merging pull requests.
