# Plugin Development Guide

PromptEngine is built to be highly extensible. This guide outlines how to write custom plugins to contribute custom detectors, templates, validators, auto-fixes, and reviewers.

---

## 1. Plugin Directory Layout

Plugins are located under `.promptengine/plugins/[plugin-id]/`. Each plugin must contain a `plugin.json` manifest:

```json
{
  "id": "my-custom-plugin",
  "version": "1.0.0",
  "description": "Contributes custom quality checkers and auto-fixes",
  "entry_point": "main.go",
  "capabilities": ["detector", "validator", "fix"]
}
```

---

## 2. Contribution Interfaces

Plugins interact with the PromptEngine core using Go interfaces.

### A. Custom Project Detectors
Extend discovery stage analysis by implementing `Stage`:

```go
type Stage interface {
	Name() string
	Execute(ctx context.Context, fs filesystem.FileSystem, pm *ProjectModel) error
}
```

### B. Custom Validators
Register workspace validators under `internal/domain/quality/validation/`:

```go
type Validator interface {
	ID() string
	Validate(fs filesystem.FileSystem) ([]Finding, error)
}
```

### C. Custom Auto-Fix Actions
Contribute reversible, safe automated fixes under `internal/domain/quality/fix/`:

```go
type RepairAction interface {
	ID() string
	Description() string
	Safety() FixSafety
	Preview(fs filesystem.FileSystem) (string, error)
	Apply(fs filesystem.FileSystem) error
	Rollback(fs filesystem.FileSystem) error
}
```

---

## 3. Registering your Plugin

During bootstrap, the Plugin Engine walks `.promptengine/plugins/` directories, parses the manifests, compiles the Go plugins or references them dynamically, and calls the respective `Register` hooks.

For local registration:
```bash
# Verify plugin loading
promptengine plugins list
```
