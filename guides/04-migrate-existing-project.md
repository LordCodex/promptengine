# 04. Migrating Existing Projects to PromptEngine

This guide outlines the step-by-step migration path for developers moving an active software project under PromptEngine governance. It provides practical instructions for audit reconciliation and legacy rules cleanup.

---

## Migration Steps

Follow this workflow sequence to migrate your codebase to PromptEngine:

```text
[1. Clone PromptEngine] ──> [2. Run Migration Prompt] ──> [3. Review AGENTS.md]
                                                                │
   [6. Start Coding]    <── [5. Reconcile Delta]       <── [4. Generate Specs] ┘
```

### Step 1: Add PromptEngine to your Workspace
Clone PromptEngine into your repository or in an external parent folder. Ensure `playbook-manifest.json` and standard stacks are accessible.

### Step 2: Run the Migration Audit Prompt
Start your AI assistant and copy the **Migrate Existing Project Prompt** from `prompts/03-migrate-existing-project.md`. This prompt instructs the AI to scan package managers, directories, database schemas, and API endpoints, and to identify any legacy developer guidelines or rule overrides.

### Step 3: Review and Refine AGENTS.md
The AI will generate `AGENTS.md` in your project root. Review Section 2 (Project Constitution) to ensure that the detected backend, database, authentication approach, and deployment details are correct. 
*If your project has legacy constraints (e.g. database tables that cannot be refactored or custom coding styles), ensure they are explicitly logged under Section 2.3 (Operating Constraints) or Section 2.5 (Approved Exceptions).*

### Step 4: Generate or Merge core specifications under docs/
The AI will generate the 10 core documents under `docs/`.
- If you have pre-existing markdown guides or text specifications, instruct the AI to **merge** them into the PromptEngine template structures (`docs/PRD.md`, `docs/Architecture.md`, etc.) rather than overwriting them, preserving historical context.
- Mark any untested or legacy parts of the system as *"Legacy/Inferred"* in `docs/Progress.md` and `docs/Decisions.md` to prevent future AI agents from refactoring them blindly.

### Step 5: Reconcile Code vs. Documentation Deltas
Run a diff audit. If the generated documentation conflicts with live code behavior:
- If the code is correct, adjust the document.
- If the code behavior is a legacy bug or technical debt, log it as a technical debt item in `docs/Progress.md` and maintain the documentation as the target "to-be" state.

### Step 6: Commit and Onboard
Once aligned, commit `AGENTS.md` and the `docs/` directory to git. From this point on, configure your IDE assistant (Cursor, Claude Code, etc.) to use `AGENTS.md` as its entry point.

---

## Legacy Rules Cleanup

When migrating, you might have legacy instructions files (like `instructions.txt` or `prompt.md` files) littered in your directories.
- **Rule**: Consolidate all project-specific guidelines into the newly generated `AGENTS.md` (Project Constitution) or `docs/BusinessRules.md`. 
- **Action**: Delete the old separate rule text files to prevent AI context window saturation and rule duplication. Keep the repository clean.
