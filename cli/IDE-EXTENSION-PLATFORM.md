# IDE Extension Platform

PromptEngine IDE integrations use stable public boundaries and do not import internal Go packages.

## Architecture

```text
IDE Extension -> PromptEngine CLI -> PromptEngine Engines
```

The first implementation is a VS Code extension in `extensions/vscode`. The communication layer is intentionally isolated so future JetBrains and LSP integrations can use the same command surface or move to a local service transport later.

## VS Code Commands

- `PromptEngine: Analyze Project`
- `PromptEngine: Generate Context`
- `PromptEngine: Generate AI Prompt`
- `PromptEngine: Run Workflow`
- `PromptEngine: Check Health`
- `PromptEngine: Sync Documentation`
- `PromptEngine: Selected Code to Context`
- `PromptEngine: Analyze Current File`

## Configuration

- `promptengine.path`: PromptEngine CLI path.
- `promptengine.configPath`: optional project config path.
- `promptengine.preferredAIClient`: `claude`, `codex`, `chatgpt`, `cursor`, `windsurf`, or `generic`.
- `promptengine.outputFormat`: `markdown`, `json`, or `yaml`.
- `promptengine.contextLimitBytes`: maximum context package size.

## Public SDK

The public Go SDK boundary lives in `pkg/sdk`. It exposes editor-safe operations:

- project analysis
- context generation and export
- prompt generation
- workflow execution
- documentation sync
- health checks

The SDK shells out to the CLI by default, preserving the architecture and avoiding direct dependencies on internal packages.
