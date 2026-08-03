# 09. Troubleshooting Guide

This guide details resolutions for common issues encountered when integrating PromptEngine with AI coding assistants.

---

## 1. Issue: AI Coding Assistant is Ignoring Playbooks or Rules

### Symptom
The AI generates code that violates naming conventions, ignores directory boundaries, or bypasses security gates.

### Resolution
1. **Reset Chat Session**: The conversation history may be too long, saturating the model's context window. Start a fresh chat.
2. **Re-establish Entry Context**: In your first message, explicitly command: *"Before planning, read the project constitution (`AGENTS.md`) and the PromptEngine bootstrap entry point (`promptengine/ai/bootstrap.md`)."*
3. **Targeted Playbook Injections**: Explicitly mention the specific playbook in chat (e.g. *"Conform strictly to stacks/php-laravel/laravel-logic.md"*).

---

## 2. Issue: AGENTS.md or project documentation is out of sync with code

### Symptom
Database schema modifications or routing edits are made in code, but the documentation under `docs/` or `AGENTS.md` remains unchanged.

### Resolution
1. **Rejection Policy**: Reject any pull requests or code outputs from the AI that do not include updates to the corresponding documentation.
2. **Documentation Generation Run**: Feed the modified source files (e.g. database migrations or controller routes) to the AI and run the **Documentation Sync Prompt** from `prompts/08-documentation-sync.md`. Command it to reconcile the docs.
3. **Run CLI Checks**: Run the future CLI command `promptengine doctor` (specified in [doctor spec](../cli/specs/doctor.md)) to automatically validate link targets, or `promptengine analyze` (specified in [analyze spec](../cli/specs/analyze.md)) to map local code diff changes to recommended documentation updates.

---

## 3. Issue: Migration and reverse-engineering errors in Existing Projects

### Symptom
When bootstrapping an existing repository, the AI makes incorrect assumptions about active technology stacks or database relationships.

### Resolution
1. **Clarify Constraints**: Check `AGENTS.md` Section 2.5 (Approved Exceptions). If the AI reverse-engineered a pattern incorrectly (e.g. assuming UUIDv7 instead of incrementing integers), manually log this as an exception.
2. **Audit Lock Files**: Ensure dependency lock files (`composer.lock`, `package-lock.json`, `pubspec.lock`) are committed and readable. If they are ignored by git, the AI will default to speculative estimates.

---

## 4. Issue: AI ignores workflows (Planning before coding)

### Symptom
The AI begins modifying multiple files immediately without outlining an implementation plan or obtaining approval.

### Resolution
1. **Interrupt Execution**: Stop the AI's generation thread immediately.
2. **Enforce the Workflow**: Remind the AI of its rules: *"Stop. Under the PromptEngine Core Rules, you must write a detailed implementation plan and obtain my explicit approval before modifying files. Present your plan first."*
3. **System Instructions Check**: Double-check that your IDE's system instructions include the planning directive.

---

## 5. Issue: Documentation Conflicts

### Symptom
The generated specifications (e.g. `docs/PRD.md`) contradict what is implemented in code or what was stated in historical ADR decisions.

### Resolution
1. **Apply Source of Truth Hierarchy**:
   - Explicit instructions override documents.
   - Approved documents (`docs/` & `AGENTS.md`) override source code.
   - Source code overrides historical Decisions logs.
2. **Reconcile**: Command the AI: *"The implementation in file X contradicts the rule in docs/Y. Reconcile this discrepancy immediately. If the code is correct, update the document. If the document is correct, refactor the code."*
