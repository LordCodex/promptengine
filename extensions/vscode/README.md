# PromptEngine VS Code Extension

This extension is the first IDE integration layer for PromptEngine. It communicates with the stable PromptEngine CLI and does not import Go internals or require AI provider API keys.

## Commands

- `PromptEngine: Analyze Project`
- `PromptEngine: Generate Context`
- `PromptEngine: Generate AI Prompt`
- `PromptEngine: Run Workflow`
- `PromptEngine: Check Health`
- `PromptEngine: Sync Documentation`
- `PromptEngine: Selected Code to Context`
- `PromptEngine: Analyze Current File`

## Configuration

- `promptengine.path`: CLI binary path.
- `promptengine.configPath`: optional project config path.
- `promptengine.preferredAIClient`: `claude`, `codex`, `chatgpt`, `cursor`, `windsurf`, or `generic`.
- `promptengine.outputFormat`: `markdown`, `json`, or `yaml`.
- `promptengine.contextLimitBytes`: maximum context size.

## Architecture

VS Code runs PromptEngine through the CLI today:

```text
VS Code Extension -> PromptEngine CLI -> PromptEngine Engines
```

The client boundary is isolated so a future local service or LSP transport can replace CLI execution without changing command behavior.
