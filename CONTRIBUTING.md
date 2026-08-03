# Contributing to PromptEngine

Thank you for choosing to contribute to PromptEngine! This document outlines the standards, guidelines, and steps required to build, test, and submit contributions.

---

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please report any violations to **security@promptengine.dev**.

---

## Technical Stack & Code Standards

PromptEngine CLI and engines are built using Go. We adhere to Go community standards and best practices:
1. **Clean Code**: Keep packages focused, handle errors explicitly, and structure domain packages cleanly.
2. **Go Formatting**: Code must be formatted with `go fmt`. Run `make fmt` before committing.
3. **No Circular Dependencies**: Ensure package imports follow the strict architectural flow:
   - `cmd/` -> `internal/app/` -> `internal/domain/` & `internal/filesystem/` & `internal/ui/`
   - Domain packages must never import the `app` or `ui` layers.

---

## Playbook Contribution Standards

If you are contributing template documents, playbooks, or standards under `core/` or `guides/`:
1. **Metadata Frontmatter**: Every markdown document must start with a YAML frontmatter block:
   ```yaml
   ---
   document_id: unique-kebab-case-id
   title: Descriptive Title
   ecosystem: cross-cutting | php-laravel | dart-flutter | js-ts-vue-nuxt
   target_versions:
     go: ">=1.21"
   audience: human, agent
   last_reviewed: YYYY-MM-DD
   ---
   ```
2. **Deterministic Rules**: Avoid vague guidelines. Describe **Why**, **When**, and **Before/After** examples where relevant.
3. **Link Validity**: Always use relative markdown links (e.g. `[link](core/01-standard.md)`). Avoid absolute URLs.

---

## Local Development Workflow

We use a `Makefile` to automate common development tasks:

```bash
# Clone the repository
git clone https://github.com/LordCodex/promptengine.git
cd promptengine

# Format code
make fmt

# Run linters
make lint

# Run all unit tests
make test

# Build local binary
make build
```

---

## Pull Request Guidelines

1. **Self-Verification**: Ensure the build compiles cleanly and all unit tests pass.
2. **Documentation**: If your change modifies CLI commands, ensure configuration flags or new subcommands are updated in the docs.
3. **Descriptive Commits**: Use descriptive commit messages (e.g. `feat: add auto-fix action for manifest integrity`).
4. **Structured PR**: Fill out the provided Pull Request template completely.
