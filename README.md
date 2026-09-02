# PromptEngine

PromptEngine is an engineering-rule **orchestrator** for human developers and AI coding agents. It detects a project's stack, resolves the applicable authoritative standards, selects task-relevant guidance, preserves project-specific instructions, and provides workflows/prompts/bridges without forcing every rule into every task.

PromptEngine is no longer the sole source of truth for framework engineering policy. Canonical standards are maintained in dedicated repositories and pinned in `sources/rules-sources.yaml`.

## Rule architecture

```text
LordCodex/engineering-ai-rules
        |
        +--> LordCodex/php-broilerplate
        |        |
        |        +--> LordCodex/laravel-ai-rules
        |
        +--> LordCodex/vue-ai-rules
        |
        +--> LordCodex/react-ai-rules

                PromptEngine
        detect -> resolve -> sync/select
                    |
                    +--> cross-stack bridges
                    +--> workflows/prompts
                    +--> project rules/knowledge
                    |
                    v
             task-specific context
```

The external repositories own their engineering rules. PromptEngine owns orchestration: stack detection, source resolution, synchronization/cache, task-context selection, project knowledge, workflows, prompts, AI adapters, and cross-stack integration guidance.

## Preservation policy

PromptEngine contains a large historical Markdown playbook library. It is not deleted merely because a newer authoritative repository exists.

Every legacy rule is handled semantically using:

- `KEEP` — PromptEngine still owns the material or no authoritative owner exists yet.
- `ADD` — useful missing knowledge should be added to the authoritative owner.
- `MERGE` — both sides contain useful knowledge that must be combined.
- `UPDATE` — PromptEngine material is useful but stale or overprescriptive and must be normalized.
- `REFERENCE` — the authoritative owner now contains the useful rule; PromptEngine keeps only routing/fallback value.
- `CONFLICT` — the standards disagree and require an explicit decision.

See `sources/migration-audit.yaml`.

The bundled library remains available as preservation/offline fallback. PromptEngine must not silently activate a partial authoritative source set: a profile uses synchronized authoritative sources only when all required sources and entrypoints are present; otherwise the bundled playbooks remain the safe fallback.

## Supported authoritative profiles

PromptEngine currently resolves authoritative profiles for:

- PHP
- Laravel
- Laravel + Livewire
- Laravel + Inertia + Vue
- Laravel + Inertia + React
- Vue
- React

Nuxt, Next.js, Dart, and Flutter material remains bundled fallback until dedicated authoritative ownership is established for those specializations.

For Laravel + Inertia + Vue, the intended responsibility model is:

```text
Universal invariants
        |
       PHP
        |
     Laravel -------- Vue
        |              |
     Inertia           |
        \              /
         integration bridge
                |
          project rules
                |
          task context
```

Laravel remains authoritative for server routing, validation, authentication/session behavior, authorization, transactions, persistence, queues, and server-side business invariants. Inertia rules own the server/client page-prop/navigation/form integration boundary. Vue owns Vue component/reactivity/composable/client-state implementation. Project-specific decisions specialize those layers without weakening security, correctness, reliability, or data integrity.

## CLI rule commands

PromptEngine exposes explicit rule-source commands:

```bash
promptengine rules resolve
promptengine rules status
promptengine rules sync
```

`rules resolve` detects the project and reports the most-specific matching profile and source inheritance.

`rules status` reports whether the pinned authoritative sources required by the detected profile are present locally and usable.

`rules sync` downloads the exact pinned source snapshots into `.promptengine/rules/<source>/<ref>/...`. It does not modify the authoritative repositories and does not overwrite project-specific instructions. Private GitHub repositories may be accessed with `GITHUB_TOKEN` or `GH_TOKEN`; PromptEngine does not persist those tokens.

After synchronization, context generation can use the authoritative source overlay. If a required snapshot/entrypoint is missing, PromptEngine falls back to its bundled playbook rather than mixing an incomplete authority set.

## Project-specific precedence

Project instructions are intentionally preserved as the most-specific project layer. A typical order is:

```text
universal invariants
-> language/framework standards
-> integration bridge where relevant
-> project AGENTS.md / manifest / architecture decisions
-> task-specific context
```

Project rules may specialize framework implementation but must not silently weaken universal/server security, authorization, correctness, data-integrity, or reliability guarantees.

## Minimal context loading

PromptEngine should not concatenate all Markdown from all matched repositories.

For each task it should:

1. detect the project stack;
2. resolve the most-specific profile;
3. identify applicable authoritative sources;
4. include only task-relevant capability documents/skills;
5. include an integration bridge only when the stack requires one;
6. include relevant project instructions/decisions;
7. stay within the requested context budget.

A Laravel migration task should not need Vue component guidance. A Laravel + Inertia + Vue form may require Laravel validation/authorization, Inertia form/error semantics, Vue component/form behavior, and the integration bridge.

## Cross-stack bridges

`bridges/` contains only integration guidance that has no single framework owner.

For example, `bridges/laravel-inertia-vue.md` explains the authority boundary between Laravel, Inertia, and Vue, prop/navigation/form integration, state ownership, version awareness, failure behavior, and combined testing. It must not become a duplicate Laravel or Vue engineering standard.

Compatibility bridge files may remain when older manifests refer to their IDs, but they should redirect to current authoritative sources rather than carrying a second framework policy.

## Bundled playbook library

PromptEngine still embeds its playbook library so an installed binary can operate without cloning this repository or requiring network access for every task.

Important areas include:

```text
core/                 cross-cutting historical/source playbooks
stacks/               stack-specific fallback/source material
security/             security fallback/source material
performance/          performance fallback/source material
legacy/               modernization guidance
bridges/              cross-stack integration guidance
workflows/            execution workflows
prompts/              reusable task prompts
project/              project-knowledge standards/templates
ai/                   AI system/bootstrap instructions
sources/              authoritative source registry/profiles/migration audit
```

Do not treat every bundled framework file as canonical after its migration has completed. `sources/rules-sources.yaml` and `sources/migration-audit.yaml` define ownership.

## AI coding-agent workflow

For an implementation task, agents should:

1. read project instructions (`AGENTS.md` or equivalent) and PromptEngine configuration/manifest when present;
2. detect/confirm the stack and relevant project constraints;
3. resolve the rule profile;
4. load only the task-relevant authoritative/fallback documents;
5. inspect the existing implementation and conventions;
6. plan the smallest coherent change;
7. preserve feature/module boundaries and public contracts;
8. review security, performance, reliability, failure behavior, and deployment consequences as relevant;
9. run the available checks when the environment permits;
10. report exactly what was and was not verified.

Never fabricate test, build, audit, benchmark, or CI results.

## Existing guides and prompts

The repository also includes onboarding and workflow material under `guides/`, `prompts/`, `workflows/`, `project/`, and `ai/`. These remain PromptEngine-owned orchestration/project-knowledge resources.

Some older guides may describe cloning PromptEngine as the primary standards source. The source resolver and `promptengine rules ...` commands are now the preferred architecture; older onboarding text should be interpreted in that context until each guide is refreshed.

## Development and validation

PromptEngine is implemented as a Go CLI. Changes should be made on focused branches and verified through the repository's available Go quality/test matrix, security workflow, and example validation. Do not claim those checks pass unless the relevant run completed successfully.

## Authoritative source registry

Current source ownership and pinned refs are defined in:

- `sources/rules-sources.yaml`
- `sources/profiles/*.yaml`
- `sources/migration-audit.yaml`

Update pins only after the corresponding semantic review/merge has completed. Synchronization is deliberately reproducible: a pinned commit identifies exactly which rules were reviewed.
