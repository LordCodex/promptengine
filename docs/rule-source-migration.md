# Rule Source Migration

PromptEngine is transitioning from owning copies of framework engineering standards to orchestrating dedicated authoritative rule repositories.

This is a preservation-first migration, not a cleanup pass.

## Non-negotiable invariant

Never remove an existing PromptEngine rule merely because a similar rule exists in an authoritative repository. Compare it semantically first. Preserve implementation value, move genuinely missing guidance to the correct owner, and replace only zero-value duplication with an explicit reference.

## Classification model

Every existing rule or rule section must be classified before removal or replacement:

- `KEEP` — PromptEngine-specific orchestration, bridge, workflow, project-knowledge, or implementation guidance still belongs here.
- `ADD` — useful PromptEngine guidance is missing from the authoritative owner; add it there before changing the local source.
- `MERGE` — both sources contain useful complementary guidance; merge into the correct owner and keep only PromptEngine-specific glue locally.
- `UPDATE` — the concept remains useful but the local wording/API/version assumptions are stale or overly absolute.
- `REFERENCE` — the authoritative source fully owns the concept and the PromptEngine copy adds no implementation value.
- `CONFLICT` — the sources disagree on security, correctness, architecture, framework behavior, or ownership. Stop and review rather than silently choosing one.

## Ownership model

PromptEngine owns detection, rule resolution, task-context selection, cross-stack bridges, workflows, prompts, project knowledge, CLI behavior, and agent adapters.

Dedicated repositories own engineering rules for their scope. Their pinned locations are recorded in `sources/rules-sources.yaml`.

Project-local `AGENTS.md`, project documentation, and explicit project decisions remain the final specialization layer. They may specialize implementation details but must not silently weaken inherited security, correctness, reliability, or data-integrity requirements.

## First migration slice: Laravel + Inertia + Vue

The first end-to-end profile is `sources/profiles/laravel-inertia-vue.yaml`.

### `bridges/laravel-inertia-vue-engineering-standard.md`

Current classification: `MERGE` + `UPDATE`, then retain as bridge-only guidance.

Useful concepts already preserved by the newer Laravel Inertia standard include:

- Laravel remains the server authority.
- frontend visibility is not authorization.
- Laravel routes remain the normal navigation source of truth.
- shared props should stay small.
- page props are a public presentation contract.
- backend validation remains authoritative.
- Inertia forms should use supported Inertia form primitives.
- backend duplicate protection is required for duplicate-sensitive mutations.
- partial/lazy/deferred data should reduce real work rather than hide inefficient queries.
- backend and frontend tests prove different things.

Local bridge guidance that still has value should remain only when it explains cross-stack interaction rather than re-stating Laravel or Vue rules.

The following local statements require `UPDATE` before they can remain normative:

- pinned Laravel 11 / Inertia 1 / Vue 3.4 assumptions;
- describing Inertia with version-specific APIs without first checking installed adapter versions;
- requiring Ziggy for every project instead of allowing the route helper selected by the project;
- requiring API Resources for every prop instead of requiring explicit safe shaping and using resources/DTO-like shapes when useful;
- requiring cursor pagination everywhere instead of choosing pagination from dataset and access-pattern requirements;
- requiring lazy props everywhere for slow data instead of measuring and using the current installed Inertia lazy/deferred/partial mechanisms appropriately;
- treating every slow third-party operation, including payment work, as automatically queueable without considering synchronous confirmation, idempotency, transaction ordering, and user-visible correctness;
- describing frontend disabled/loading state as duplicate protection; it is UX only.

### `stacks/js-ts-vue-nuxt/vue-ts-engineering-standard.md`

Current classification: `MERGE` + `UPDATE` + eventual `REFERENCE` for concepts fully owned by `vue-ai-rules`.

Useful material already represented in the authoritative Vue standard includes strict TypeScript, Vue 3 SFCs, Composition API, `<script setup lang="ts">`, immutable props, typed props/emits, local state first, feature-owned Pinia stores when justified, composables, runtime validation at untrusted boundaries, feature architecture, security, accessibility, performance, and testing.

Local rules that should not be copied upstream mechanically include arbitrary thresholds such as mandatory component splitting at a fixed line count or a fixed maximum prop count. Review responsibility and cohesion instead of enforcing magic numbers.

The old statement that TypeScript `strict: true` automatically enables `noUnusedLocals` must not be preserved as a technical claim; that compiler option is independent.

### PromptEngine core and domain Markdown

No core, security, performance, workflow, checklist, guide, prompt, or project Markdown is deleted in this migration slice. Those files remain source material until individually classified against the authoritative repositories. PromptEngine-specific operational guidance is expected to remain here.

## Resolution rule

For a detected Laravel + Inertia + Vue project, the resolver should conceptually load:

1. project-local constitution and explicit decisions;
2. the Laravel + Inertia + Vue profile;
3. universal rules required by the task;
4. PHP rules required by the task;
5. Laravel rules required by the task;
6. Laravel Inertia integration rules required by the task;
7. Vue rules required by the task;
8. PromptEngine bridge/workflow/checklist material only when it adds cross-stack or task-specific value.

The resolver must not concatenate whole repositories by default.

## Migration completion criteria

A PromptEngine Markdown source may be converted to a reference only when:

1. its sections have been semantically compared;
2. missing useful guidance has been merged into the proper owner;
3. stale or conflicting guidance has been corrected or retired explicitly;
4. PromptEngine-specific bridge/orchestration value has been retained;
5. links/manifest mappings resolve to the replacement source;
6. relevant resolver/context tests cover the new path;
7. no test or CI result is claimed unless actually executed.
