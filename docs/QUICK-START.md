# Quick Start Guide

Welcome to PromptEngine! This guide will walk you through setting up PromptEngine and using it to govern your repository coding standards and generate context-optimized prompts for AI assistants.

---

## 1. Initialise the Workspace

Run the `init` command at the root of your project:

```bash
promptengine init
```

This bootstraps:
- `.promptengine/` configuration directory
- `playbook-manifest.json` constitution file
- `docs/` templates (scaffolds: `Architecture.md`, `Database.md`, `API.md`, `Decisions.md`)

---

## 2. Scan Workspace Stack

Run `scan` to discover the languages, frameworks, and tools in your repository:

```bash
promptengine scan
```

This runs deterministic detection checks to analyze your code without AI overhead.

---

## 3. Verify Setup Health

Run `doctor` to identify missing standards, broken references, or config drifts:

```bash
promptengine doctor
```

To attempt automated repairs:
```bash
promptengine doctor --fix
```

To view a detailed score metric breakdown:
```bash
promptengine health
```

---

## 4. Run Code and Standards Review

Scan code modifications or directories against standards:

```bash
promptengine review --path .
```

---

## 5. Fetch Task Context

When pairing with an AI coding assistant (like Claude, ChatGPT, or Cursor), request only the minimal context files necessary for a specific engineering task to reduce tokens:

```bash
promptengine context --task bug-fix
```

---

## 6. Output AI Prompts

Generate complete, context-injected markdown prompts ready to be copied into your AI assistant:

```bash
promptengine prompt --workflow bug-fix --provider claude
```
