# 01. New Project Bootstrap Prompt

---

## Purpose
Instructs the AI assistant to conduct a discovery interview and automatically generate the project's AI Constitution (`AGENTS.md`) and the 10 core documentation files under `docs/`.

## When to use
Use when starting a brand-new, greenfield software project from scratch in an empty directory.

## Example
Initializing a SaaS Billing platform using Laravel, Vue 3, and PostgreSQL.

---

## Copy-and-Paste Prompt

```markdown
You are going to initialize a brand-new project named "{PROJECT_NAME}".

Before writing any code or generating folders, you must execute the PromptEngine New Project Bootstrap workflow.

Follow these rules:
1. Conduct a brief discovery interview to identify my requirements. Keep questions grouped, concise, and focused on high-impact choices (architecture, frameworks, database, authentication, primary keys, and security).
2. Adhere strictly to the Discovery Efficiency Rules to minimize developer cognitive load.
3. Once requirements are gathered, generate:
   - `AGENTS.md` in the project root by copying and filling out Section 2 of `project/templates/AGENTS.template.md` with the confirmed tech stack, overview, and constraints.
   - A `docs/` folder containing the 10 core specifications (`PRD.md`, `Architecture.md`, `Database.md`, `API.md`, `Progress.md`, `Roadmap.md`, `Decisions.md`, `Deployment.md`, `Troubleshooting.md`) filled out with the interview specifications.
4. Stop and wait for my explicit approval before writing any codebase classes, configs, or routes.

Start the interview by presenting the first set of questions.
```

---

## Expected AI Behaviour
1. The AI will ask 3–5 high-level questions about the stack, authentication, and core business features.
2. It will not write code or scaffolding.
3. After the interview, it will create the `AGENTS.md` and the `docs/` specifications.
4. It will present clickable markdown links to the files and wait for your review.

## Common Mistakes
- **Allowing the AI to Scaffold Immediately**: Letting the AI run commands like `npm init` or `laravel new` before generating `AGENTS.md` and `docs/`.
- **Answering too many questions**: Answering an excessively long interview list. Command the AI to stick to the Discovery Efficiency Rules.
