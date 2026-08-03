# Project Instructions (AGENTS.md Template)

<!-- 
INSTRUCTIONS FOR THE HUMAN DEVELOPER:
1. Copy this file into your project root and rename it to `AGENTS.md`.
2. If PromptEngine is located in a nested folder or external shared location, update the file path references in Step 1 and Step 2 below.
-->

This project uses PromptEngine to enforce strict, high-quality software engineering standards.

Before performing any engineering task, the AI assistant must follow these steps:

1. **Read the Bootstrap Entry Point**:
   <!-- 
   If PromptEngine is located in an external directory (e.g. "../PromptEngine/"), 
   update this path to: ../PromptEngine/ai/bootstrap.md 
   -->
   Read and analyze [promptengine/ai/bootstrap.md](promptengine/ai/bootstrap.md) to understand operating principles and priority hierarchies.

2. **Read the Project Bootstrap Standard**:
   <!-- 
   If PromptEngine is located in an external directory (e.g. "../PromptEngine/"), 
   update this path to: ../PromptEngine/project/01-project-bootstrap-standard.md 
   -->
   Read and follow the workflows in [promptengine/project/01-project-bootstrap-standard.md](promptengine/project/01-project-bootstrap-standard.md) to adopt or initialize the project state.

3. **Read the Index Manifest**:
   <!-- 
   If PromptEngine is located in an external directory, 
   update this path to: ../PromptEngine/playbook-manifest.json 
   -->
   Read and parse [promptengine/playbook-manifest.json](promptengine/playbook-manifest.json) to locate the relevant standards mapping for the current task.

4. **Identify & Load Required Playbooks**:
   Determine which specific playbooks match the task classification. Load **only** the playbooks required for the current task to conserve context space. Do not load unrelated playbooks.

5. **Load Project-Specific Context**:
   After loading PromptEngine standards, read local project documentation to understand requirements and invariants:
   - [docs/PRD.md](docs/PRD.md) (Product requirements)
   - [docs/Architecture.md](docs/Architecture.md) (Service layouts, structural bounds)
   - [docs/Database.md](docs/Database.md) (Project schemas and keys)
   - [docs/API.md](docs/API.md) (Routing contracts and payload examples)
   - [docs/BusinessRules.md](docs/BusinessRules.md) (Core calculations, domain workflows)

6. **Priority Order of Guidance**:
   Always follow project-specific guidelines and user instructions before generic PromptEngine rules. Keep all applicable rules active throughout the coding process.

---

## Project-Specific Rules

<!-- 
Write any rules unique to this project below. 
Example: "This project uses MongoDB instead of MySQL, so database rules should target document mapping."
-->

- **Technology Stack**: [Add your frameworks, e.g. Laravel, React, Vue]
- **Special Configurations**: [Add database settings, caching configs, or hosting specifics]
- **Security Constraints**: [Add custom access control details]

---

## Best Practices

To maintain a clean and scalable developer setup:
- **Keep AGENTS.md Concise**: Keep this instructions file brief and focused on onboarding. Place detailed specifications in `docs/`.
- **Project Rules belong in AGENTS.md or docs/**: Put project-specific parameters (like database schemas or specific API keys) inside project docs.
- **Keep Reusable Standards in PromptEngine**: Keep all language conventions, generic security checklists, and deployment standards inside the shared PromptEngine repository.
- **Avoid Duplicating Rules**: Never copy-paste PromptEngine guidelines into project documentation or `AGENTS.md`.
- **Update PromptEngine**: When introducing new language stacks or universal standards, update PromptEngine directly instead of spreading rules across individual project files.
- **Read Minimized Context**: Ensure the AI assistant loads only the playbooks required for the current task.
- **Conflict Strategy**: Prefer local project documentation whenever business rules conflict with generic PromptEngine standards.
