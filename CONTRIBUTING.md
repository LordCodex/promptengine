# Playbook Contribution Guidelines

This document details how to contribute new standards, architectural decisions, and coding rules to the playbook repository.

---

## Philosophy of Contribution

1. **Keep Rules Actionable**: Rules must be enforceable. Avoid vague philosophy; explain **Why**, **When**, **When Not**, and provide **Before/After** code examples.
2. **De-duplicate Aggressively**: If a principle applies across all stacks, place it under `/core/`. Do not write duplicate explanations in language-specific files.
3. **Minimize Token Footprint**: Keep documents highly focused. Average document size should not exceed 250 lines (~1,500 tokens).

---

## Standard Metadata Frontmatter

Every markdown file must start with a YAML frontmatter block to allow the indexing engine and AI agents to parse compatibility:

```yaml
---
document_id: [unique-kebab-case-string]
title: [Short Descriptive Title]
ecosystem: [php-laravel | dart-flutter | js-ts-vue-nuxt | cross-cutting]
target_versions:
  [tech-name]: "[semver range]"
dependencies:
  - [dependency-document-id]
audience: [human, agent]
last_reviewed: YYYY-MM-DD
---
```

---

## Sandbox Integration Rules

When contributing code examples:
1. Do not rely entirely on inline markdown code blocks for complex patterns.
2. Add a fully functioning, clean implementation in the `/examples/` directory.
3. Verify that the example project compiles and passes lint stages locally before proposing a merge.
4. Run standard ecosystem linter validations:
   - PHP: `vendor/bin/phpstan analyse`
   - Dart/Flutter: `flutter analyze`
   - Node/Vue: `npm run lint`

---

## Link Validity & Coherence

Ensure that all cross-references use relative file links (e.g. `[thinking.md](core/01-thinking-and-planning.md)`). Do not use absolute host links (such as github.com or gitlab.com) for internal file paths, as these links break when read by offline local agents.
