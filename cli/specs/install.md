# CLI Command Specification: `install`

## Purpose

Install a local PromptEngine package or plugin from a YAML or JSON manifest.

## Inputs

- Optional argument: `[id]`, the package identifier.
- Required flag: `--manifest`, the local manifest path.

## Outputs

Installed package files, metadata, and state are written under:

```text
.promptengine/plugins/<id>/
```

## Workflow

1. Read and validate the local manifest.
2. Validate compatibility, permissions, and declared source files.
3. Copy declared files and persist plugin metadata.
4. Exit with code `1` when validation or installation fails.

Example:

```bash
promptengine install --manifest ./plugins/company/plugin.yaml
promptengine plugin list
```

Uninstalling removes the installed package state:

```bash
promptengine plugin remove company
```

## Related Documentation

- [Plugin Architecture Guide](../CLI-Architecture.md)
