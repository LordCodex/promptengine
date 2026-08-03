# Security Architecture Guide

This document details the security principles, boundaries, and hardening mechanisms implemented in PromptEngine.

---

## 1. Sandbox Boundaries & Path Traversal Prevention

PromptEngine interacts directly with your local codebase filesystem. To prevent malicious configurations or compromised plugins from accessing files outside your project root (e.g. reading private SSH keys in `~/.ssh` or system configs in `/etc`), all file operations strictly pass through path traversal guards in `internal/security/sandbox.go`:

```go
func ValidateSafePath(baseDir, targetPath string) (string, error)
```

Any attempt to reference paths outside the project root directory triggers a validation exception:

```
Error: path traversal detected: path '../../.ssh/id_rsa' goes outside base directory '.'
```

---

## 2. Command Injection Defenses

PromptEngine operates deterministically.
- All discovery stages use standard file readings and parsing algorithms instead of executing shell scripts.
- CLI subcommands use structured process execution vectors rather than raw shell pipelines (preventing execution of trailing commands via string injections).

---

## 3. Telemetry and User Consent

- **Zero-Tracking by Default**: Usage telemetry is strictly disabled by default. No logs are transmitted.
- **Opt-in Only**: Telemetry is only enabled if the user explicitly opts in via configuration files or environment overrides.
- **Anonymised Data**: When enabled, collected data is restricted to command names, run times, and platform architecture metrics. No file names, codebase contents, or private keys are ever collected or sent.
