# CLI Configuration (CLI-Configuration.md)

This document specifies the format, location, and parameters of the PromptEngine CLI configuration files.

---

## 1. Configuration File Locations

The CLI reads configuration settings from two levels:

1. **Global Configuration**:
   - Location: `~/.promptengine/config.json`
   - Purpose: User-specific preferences (default AI provider, editor commands, global plugin registries, custom local PromptEngine repository path).
2. **Local Project Configuration**:
   - Location: `[project-root]/.promptengine.json`
   - Purpose: Project-specific settings overrides (active stack parameters, custom documentation path mappings, specific plugin bindings, exception registers).

---

## 2. Configuration Schema (JSON Format)

The local `.promptengine.json` file uses the following schema definition:

```json
{
  "$schema": "https://promptengine.dev/schemas/config.v1.json",
  "project": {
    "name": "MySaaSPlatform",
    "version": "1.0.0",
    "description": "SaaS Billing and user dashboard",
    "promptengine_path": "./promptengine"
  },
  "docs": {
    "root_dir": "docs",
    "agents_constitution": "AGENTS.md",
    "api_spec": "docs/API.md",
    "database_spec": "docs/Database.md",
    "business_rules": "docs/BusinessRules.md"
  },
  "stack": {
    "detection_overrides": {
      "language": "php",
      "language_version": "8.3",
      "framework": "laravel",
      "framework_version": "11.x",
      "database": "postgresql",
      "caching": "redis",
      "testing": "pest"
    }
  },
  "ai": {
    "default_provider": "claude",
    "model_preferences": {
      "claude": "claude-3-5-sonnet",
      "chatgpt": "gpt-4o",
      "gemini": "gemini-1.5-pro"
    },
    "token_optimization": {
      "compress_context": true,
      "max_context_tokens": 16384,
      "exclude_paths": [
        "node_modules/**",
        "vendor/**",
        "public/build/**"
      ]
    }
  },
  "plugins": [
    {
      "id": "laravel-logic-standards",
      "enabled": true,
      "source": "npm",
      "version": "^1.0.0"
    }
  ],
  "cli": {
    "colors": true,
    "interactive": true,
    "quiet_mode": false,
    "verbose": false
  }
}
```

---

## 3. Configuration Parameter Details

### `project`
- `promptengine_path`: Path pointing to the location of the PromptEngine folder. Can be local (relative path) or global (absolute path).

### `docs`
- Maps path overrides. Allows renaming `docs/` or moving documentation files to, for example, `.agents/` or `.github/wiki/` directories.

### `stack`
- `detection_overrides`: If the Project Detection Engine makes incorrect assumptions or cannot parse a custom compilation framework, developers specify stack versions here to lock down playbooks loading.

### `ai`
- `default_provider`: The target target client. Controls prompt construction headers formatting.
- `compress_context`: When enabled, CLI truncates history and extracts only files structure to save token limits.

### `plugins`
- An array containing enabled extensions that load custom playbooks or stack-specific conventions rules.
