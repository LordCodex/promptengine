# CLI Command Reference

This document describes the production PromptEngine CLI surface. Commands are thin presentation adapters and delegate to the domain engines.

## Global Flags

Every command supports:

- `--config`: path to project configuration.
- `--debug`: enable debug logging.
- `--json`: render machine-readable JSON.
- `--output text|json|yaml`: select output format.
- `--verbose`, `-v`: enable verbose output.

## Commands

| Command | Purpose | Engine |
| :--- | :--- | :--- |
| `init` | Initialize PromptEngine project files and run bootstrap workflow | Configuration, Manifest, Workflow |
| `migrate` | Migrate active configuration and validate project state | Configuration, Quality |
| `scan` | Analyze repository structure and stack | Discovery |
| `detect` | Alias for `scan` | Discovery |
| `context` | Build optimized context package | Context |
| `context export` | Export optimized context for a coding agent | Context, Agent Integration |
| `workflow` | Execute a registered workflow | Workflow |
| `docs generate` | Generate a document from templates or defaults | Documentation |
| `docs validate` | Validate existing documentation | Documentation |
| `docs sync` | Detect documentation drift from changed files | Documentation |
| `doctor` | Run project diagnostics | Quality |
| `review` | Run engineering review checks | Quality |
| `audit` | Run full quality audit | Quality |
| `health` | Calculate health score | Quality |
| `prompt` | Export an AI-ready prompt package for Claude, Codex, ChatGPT, or generic clients | Context, AI compiler |
| `ai` | Execute an AI request through a provider adapter | AI |
| `plugin list` | List registered plugins | Extensibility |
| `plugin install` | Install a local plugin manifest and declared files | Extensibility |
| `plugin remove` | Remove an installed plugin | Extensibility |
| `plugin enable` | Enable an installed plugin | Extensibility |
| `plugin disable` | Disable an installed plugin | Extensibility |
| `plugin health` | Run plugin diagnostics | Extensibility |
| `agents list` | List supported coding agent profiles | Agent Integration |
| `agents sync` | Generate or update agent instruction files | Agent Integration |
| `profile init` | Create a local personal developer profile | Personal Workflow |
| `profile show` | Show personal preferences | Personal Workflow |
| `task` | Generate a complete AI-ready task package from one request | Personal Workflow |
| `verify` | Run personal pre-commit quality verification | Quality, Documentation |
| `memory show` | Show local non-sensitive memory | Personal Workflow |
| `memory add` | Store a local non-sensitive workflow note | Personal Workflow |
| `insights` | Show detected project patterns and explainable recommendations | Intelligence |
| `decisions list` | List local architecture decisions | Intelligence |
| `decisions store` | Store a local architecture decision | Intelligence |
| `impact` | Analyze potential impact of current git changes | Intelligence |
| `config view` | Show active configuration | Configuration |
| `config set` | Persist a supported project configuration key | Configuration |
| `version` | Show semantic version and build metadata | Versioning |
| `completion` | Generate shell completions | System |

## IDE Integration

The VS Code extension in `extensions/vscode` exposes PromptEngine commands inside the editor while communicating through the CLI. It supports workspace detection, selected-code context generation, current-file context export, prompt generation for external AI clients, health checks, workflow execution, and documentation sync.

The public SDK boundary for future IDE integrations is `pkg/sdk`.

## Examples

```bash
promptengine init
promptengine scan --output yaml
promptengine context --task feature --budget small
promptengine context export --task feature --agent codex --format markdown
promptengine workflow --id feature-implementation
promptengine docs generate --doc architecture --overwrite
promptengine doctor --json
promptengine prompt --task bug_fix --request "Fix the payment retry bug"
promptengine prompt --task feature --client claude --format markdown --copy
promptengine prompt --task review --client codex --format json --out review-prompt.json
promptengine ai --provider ollama --prompt "Summarize this project"
promptengine plugin install --manifest ./plugins/company/plugin.yaml
promptengine agents sync --agent all
promptengine init --agents codex,claude,cursor,windsurf
promptengine profile init
promptengine task --template feature "Add subscription billing"
promptengine verify
promptengine memory add --key test-command --value "go test ./..."
promptengine insights
promptengine decisions store --title "Use UUID primary keys" --reason "Separate public identifiers" --affected models,apis
promptengine decisions list
promptengine impact
promptengine config set --key project.name --value Example
promptengine completion zsh
```

## Exit Codes

- `0`: success.
- `1`: general or command error.
- `2`: validation warning.
- `3`: health score block.
- `4`: configuration error.

## External AI Workflow

`promptengine prompt` is the default workflow for external AI clients. It builds an optimized context package, applies a client template, estimates context size, and writes a portable prompt file without calling an AI provider or requiring API credentials.

Supported prompt export formats:

- `markdown`: writes `TASK-prompt.md`.
- `text`: writes `TASK-prompt.txt`.
- `json`: writes `TASK-prompt.json`.

Supported client templates:

- `claude`
- `codex`
- `chatgpt`
- `generic`

Use `--copy` to copy the rendered prompt package to the system clipboard.

## Agent Integration

PromptEngine can generate instruction files for existing coding agents without integrating with private app UIs or requiring model API keys.

Supported built-in profiles:

- `codex`: writes `AGENTS.md`.
- `claude`: writes `CLAUDE.md`.
- `codex-md`: writes `CODEX.md`.
- `cursor`: writes `.cursor/rules/promptengine.md`.
- `windsurf`: writes `.windsurf/rules/promptengine.md`.

Custom profiles can be configured:

```yaml
agents:
  team-agent:
    instruction_file: .team/agent-instructions.md
    format: markdown
```

Use `promptengine agents sync --agent all` to generate or refresh instruction files when PromptEngine standards, workflows, or project context change.

## Personal Workflow

Personal workflow commands are local-only and designed for one developer:

- `profile init` creates `.promptengine/profile.yaml`.
- `task` detects task intent, uses discovery and git context, applies the personal profile, builds context, and exports an AI-ready prompt.
- `verify` reuses the Quality and Documentation platforms for a pre-commit check.
- `memory add` stores non-sensitive local notes such as common commands and project habits. Secrets, API keys, credentials, and tokens should not be stored.

## Local Intelligence

The intelligence layer is deterministic and explainable. It analyzes folder structure, naming conventions, known implementation patterns, local decisions, and git changes.

- `insights` reports detected patterns, stored decisions, and recommendations.
- `decisions store` records local architecture memory in `.promptengine/decisions.yaml`.
- `impact` uses current git changes to identify likely affected areas such as API responses, database relations, tests, documentation, authentication, and authorization.
