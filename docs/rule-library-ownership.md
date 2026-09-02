# PromptEngine rule library ownership

PromptEngine is the rule orchestrator, not the long-term canonical owner of every engineering standard bundled in its original Markdown library.

The machine-readable classification is `sources/migration-audit.yaml`.

## Governing rule

Existing PromptEngine Markdown is preserved until its useful knowledge has been moved to the correct authoritative owner or intentionally retained as PromptEngine-owned integration/orchestration guidance.

Use these outcomes:

- **KEEP** — PromptEngine is still the correct owner, or there is no authoritative specialized repository yet.
- **ADD** — PromptEngine contains a useful topic that the authoritative owner does not yet have; add it upstream before converting the local copy to a reference.
- **MERGE** — both copies contain useful guidance; semantically merge missing value upstream first.
- **UPDATE** — the PromptEngine document is still useful but contains stale or over-prescriptive guidance; rewrite it around the current ownership boundary.
- **REFERENCE** — the authoritative source fully owns the rule and the PromptEngine copy adds no implementation/integration value.
- **CONFLICT** — stop automatic migration until the disagreement is resolved explicitly.

No file becomes `REFERENCE` merely because its title overlaps an upstream document.

## Current ownership boundary

Authoritative repositories own:

```text
engineering-ai-rules     universal engineering invariants
php-broilerplate         PHP implementation rules
laravel-ai-rules         Laravel / Blade / Livewire / Inertia-Laravel rules
vue-ai-rules             Vue implementation rules
react-ai-rules           React implementation rules
```

PromptEngine owns:

```text
detection
rule-source resolution and synchronization
task-relevant context selection
cross-stack bridges
workflows and prompts
project bootstrap and project knowledge
AI-adapter behavior
CLI behavior
decision-support/checklist material that is specifically about PromptEngine execution
```

Project-specific instructions remain highest precedence for project decisions, provided they do not silently weaken inherited security, correctness, reliability, or data-integrity invariants.

## Material intentionally retained for now

Nuxt, Next.js, Dart, and Flutter playbooks remain active bundled fallbacks because dedicated authoritative repositories for those specializations are not registered yet. Their existence must not cause PromptEngine to pretend a Vue, React, or universal repository fully covers framework-specific Nuxt/Next/Flutter behavior.

The design library also remains PromptEngine-owned until a separate design-system engineering standard is deliberately introduced.

## Cross-stack bridges

Bridge documents are not duplicates merely because they mention framework rules. Their legitimate role is to explain the interaction boundary between stacks.

For example, the Laravel + Inertia + Vue bridge should eventually contain only integration-specific guidance such as:

```text
Laravel server authority
        +
Inertia page/prop/form/navigation boundary
        +
Vue presentation implementation
        +
PromptEngine feature/shared ownership model
```

It should reference the current Laravel and Vue authoritative sources for framework policy rather than freezing old framework versions or duplicating entire standards.

## Migration sequence

1. Classify the old document.
2. Compare it semantically with its authoritative target.
3. Correct stale or overly absolute guidance instead of copying it upstream.
4. Add or merge genuinely missing value upstream.
5. Pin PromptEngine to the new authoritative commit.
6. Verify rule synchronization and stack resolution.
7. Only then convert exact duplicate local material into a small reference, or remove it from active playbook selection while retaining history.
8. Run PromptEngine CI and context-resolution tests before merge.

This sequence applies to every migration batch.