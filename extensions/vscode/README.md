# PromptEngine VS Code Extension

This extension connects VS Code to the local PromptEngine CLI. It does not import Go internals and does not require AI provider API keys.

## Daily AI workflow

1. Open a project that uses PromptEngine.
2. Set `promptengine.preferredAIClient` to `claude` or `codex`.
3. Run `PromptEngine: Generate AI Prompt`.
4. Describe the task in normal language.
5. PromptEngine analyzes the project, selects relevant context, builds the agent-specific prompt package, and hands it to the selected AI extension.

### Claude

When the installed Claude extension exposes `claude-vscode.editor.open`, PromptEngine opens Claude with the generated prompt pre-filled in the composer. Review the prompt and press Enter to send it.

### Codex

When the installed Codex extension exposes `chatgpt.addToThread`, PromptEngine adds the generated prompt to the active Codex thread as editor context and focuses the Codex view. Review the prepared task and send it.

### Fallback

If a selected AI client does not expose a compatible VS Code command, PromptEngine copies the generated prompt to the clipboard instead of failing.

No PromptEngine API key is required for this workflow.

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

VS Code runs PromptEngine through the local CLI:

```text
Developer task
    -> PromptEngine VS Code Extension
    -> PromptEngine CLI
    -> Discovery + Context + Personal Profile + Standards
    -> Generated AI Prompt
    -> Claude/Codex handoff
```

The client boundary is isolated so a future local service or LSP transport can replace CLI execution without changing the PromptEngine engines.
