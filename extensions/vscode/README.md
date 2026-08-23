# PromptEngine VS Code Extension

This extension connects VS Code to the local PromptEngine CLI. It does not import Go internals and does not require AI provider API keys.

## Daily AI workflow

PromptEngine uses a review-first workflow so generating a prompt never automatically consumes Claude or Codex usage.

1. Open a project that uses PromptEngine.
2. Set `promptengine.preferredAIClient` to `claude` or `codex`.
3. Run `PromptEngine: Generate AI Prompt for Review`.
4. Describe the task in normal language.
5. PromptEngine analyzes the project, selects relevant context, and generates an editable Markdown prompt.
6. Review the prompt carefully. Add, remove, or rewrite anything you want.
7. Only when you are satisfied, run `PromptEngine: Send Reviewed Prompt to AI`.

Generating and editing the prompt does not hand anything to Claude or Codex.

### Claude

When you explicitly run `PromptEngine: Send Reviewed Prompt to AI`, PromptEngine uses the reviewed editor content. If the installed Claude extension exposes `claude-vscode.editor.open`, it opens Claude with that reviewed prompt pre-filled in the composer. You still decide whether to send it.

### Codex

When you explicitly run `PromptEngine: Send Reviewed Prompt to AI`, PromptEngine uses the reviewed editor content. If the installed Codex extension exposes `chatgpt.addToThread`, it adds the reviewed prompt to the Codex thread as editor context and focuses Codex. You still decide whether to send it.

### Fallback

If a selected AI client does not expose a compatible VS Code command, PromptEngine copies the reviewed prompt to the clipboard instead of failing.

No PromptEngine API key is required for this workflow.

## Commands

- `PromptEngine: Analyze Project`
- `PromptEngine: Generate Context`
- `PromptEngine: Generate AI Prompt for Review`
- `PromptEngine: Send Reviewed Prompt to AI`
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
    -> Editable reviewed prompt
    -> Explicit user handoff
    -> Claude/Codex
```

The review gate is intentional: PromptEngine never automatically submits or hands off a newly generated prompt.

The client boundary is isolated so a future local service or LSP transport can replace CLI execution without changing the PromptEngine engines.
