# 02. Greenfield New Projects Guide

This guide explains how to initialize a brand-new project from scratch using PromptEngine. It details where files are placed, how the AI discovery interview behaves, and how the baseline specifications are generated.

---

## 1. Project Directory Layout

When starting a new project, create your empty root directory. You can connect PromptEngine in one of two configurations:

### Option A: Nested Setup (Recommended for isolated repositories)
Place a clone of the PromptEngine folder directly inside your project:

```text
MyNewProject/             # Project root
├── app/                  # Application code folder
├── promptengine/         # Cloned PromptEngine directory
└── AGENTS.md             # Generated project constitution
```

### Option B: Shared Setup (Recommended for multiple local projects)
Keep a single PromptEngine instance in a shared parent directory:

```text
SharedWorkspace/
├── PromptEngine/         # Shared PromptEngine directory
└── MyNewProject/         # Project root
    ├── app/
    └── AGENTS.md         # Generated project constitution pointing to ../PromptEngine/
```

---

## 2. Bootstrapping Workflow

To initialize the project, follow the **[New Project Bootstrap Workflow](../workflows/NewProjectBootstrap.md)**. 

### Step 1: Trigger the AI Bootstrap
Start your AI assistant in your empty workspace and copy the **New Project Bootstrapper Prompt** from `prompts/01-new-project.md`. This prompt instructs the AI to load the PromptEngine templates and manifest.

### Step 2: The Discovery Interview
The AI will act as a technical product manager and conduct a short interview to gather your requirements. To avoid cognitive overload, the AI follows **Discovery Efficiency Rules**:
- It groups related questions together.
- It only asks high-impact questions that affect database schema, authentication, stack choices, and performance targets.
- It avoids granular details that can be resolved during coding.

### Step 3: Requirements Classification
After the interview, the AI isolates requirements in its thinking space:
- **Confirmed Requirements**: Directly requested by you.
- **Assumptions**: Inferred by stack defaults or framework patterns (requires your validation).
- **Open Questions**: Architectural gaps that must be answered before generation.

### Step 4: Generating AGENTS.md and the 10 Core Specs
Once the scope is finalized, the AI will automatically create:
1. **`AGENTS.md`** in your project root. The AI copies [AGENTS.template.md](../project/templates/AGENTS.template.md) and fills out Section 2 with your project description, tech stack, caching, auth, and database settings. This becomes the project's **AI Constitution**.
2. **`docs/` Directory** containing the 10 core documents (e.g. `PRD.md`, `Architecture.md`, `Database.md`, `API.md`) populated with the interview results.

---

## 3. Human Audit & Coding Approval

Before any code is generated:
1. **Critical Review**: Open the generated `AGENTS.md` and the documents under `docs/` (using the markdown file links provided by the AI). Check for incorrect stack versions, wrong table relationships, or invalid business logic.
2. **Corrections Loop**: Tell the AI to modify any sections that do not match your vision.
3. **Explicit Approval**: The AI is strictly barred from writing code files, migrations, or layouts until you provide explicit text approval (e.g. *"The documentation is approved, begin implementation"*).

Once approved, the AI will begin executing the implementation plan.
