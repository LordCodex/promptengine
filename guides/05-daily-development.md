# 05. Daily Development Workflow

This guide details the step-by-step procedures for managing daily engineering tasks with AI coding assistants under PromptEngine. It explains how to coordinate planning, coding, and review cycles.

---

## 1. Adding a Feature

When adding a feature, follow these sequential steps:

### Step 1: Requirements Check
Copy the **Add Feature Prompt** from `prompts/04-add-feature.md` to trigger feature planning. The AI must inspect `docs/PRD.md` to read the feature requirements, scopes, and target user acceptance criteria.

### Step 2: Write the Implementation Plan
The AI will output a plan detailing:
- The database columns or tables to add.
- The classes or actions to create/modify.
- The unit and integration tests to write.
- Potential backward compatibility risks.
*Do not write code until you review and approve this plan.*

### Step 3: Implement & Test
The AI will implement the code and match the surrounding conventions. Once code is written, run the Pest, Vitest, or target test runner to ensure all tests pass.

### Step 4: Documentation Synchronization
The AI will update `docs/API.md` (if endpoints were modified) and `docs/Database.md` (if schemas changed), and log completed items in `docs/Progress.md`.

---

## 2. Fixing a Bug

To troubleshoot and fix bugs:

### Step 1: Root Cause Analysis
Use the **Bug Fix Prompt** from `prompts/05-bug-fix.md`. Provide symptoms and reproduction steps. The AI will scan the target subsystems to trace the bug.

### Step 2: Minimal Target Patch
The AI must propose the **smallest possible change** that corrects the behavior, avoiding unrelated refactoring or code formatting.

### Step 3: Regression Prevention
Write a regression test verifying the bug is fixed, and update `docs/Troubleshooting.md` if the bug represents a recurrent diagnostic symptom.

---

## 3. Refactoring Code

When refactoring legacy modules or complex structures:

### Step 1: Characterization Baseline
Use the **Refactor Prompt** from `prompts/06-refactor.md`. Ensure existing tests pass before starting.

### Step 2: Step-by-Step Refactoring
The AI will refactor the target class step-by-step (e.g. extracting methods, flattening nested control blocks, removing dead dependencies) without changing database schemas or return payloads.

### Step 3: Validation
Run tests after every single extraction step to verify that external behavior is fully preserved.

---

## 4. Architecture Changes

When executing major changes (e.g., changing authentication, adding cache, introducing a new database model):

### Step 1: Architecture Decision Record (ADR)
Use the **Architecture Change Prompt** from `prompts/07-architecture-change.md`. The AI will plan the schema change, draft an ADR in `docs/Decisions.md` (detailing alternatives and trade-offs), and write SQL DDL migrations.

### Step 2: Update project constitution
If the database engine, authentication method, or a major library has changed, update `AGENTS.md` Section 2.2 and Section 2.5 to reflect the change.

---

## 5. Documentation Updates

All code pull requests are considered incomplete until documentation matches the implementation. 
- Use the **Documentation Sync Prompt** from `prompts/08-documentation-sync.md`.
- Keep API schemas, business rule math, and environment variables fully synchronized.

---

## 6. Pre-Release Auditing

Before merging feature branches to `main` or preparing a production release:
1. Use the **Release Prompt** from `prompts/10-release.md` to trigger a production readiness review.
2. The AI will audit multi-stage Docker configurations, secrets handling (ensuring no keys are in code), and database migration lock times.
3. Update `docs/Progress.md` and draft release notes.
