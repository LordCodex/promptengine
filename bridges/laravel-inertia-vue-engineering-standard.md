---
document_id: bridges-laravel-inertia-vue-engineering-standard
title: Laravel + Inertia + Vue Engineering Bridge (Compatibility Entry)
ecosystem: bridge
audience: [human, agent]
last_reviewed: 2026-09-02
---

# Laravel + Inertia + Vue Engineering Bridge

This document is retained as a compatibility entry for existing PromptEngine manifests and references.

The previous version duplicated Laravel, Inertia, Vue, database, queue, security, testing, routing, and performance policy. Those concerns now have authoritative owners and must not be maintained independently here.

## Authoritative resolution

Resolve this stack as:

```text
Universal engineering rules
    -> PHP rules
        -> Laravel rules
            -> Laravel Inertia rules

Universal engineering rules
    -> Vue rules

Laravel/Inertia + Vue
    -> PromptEngine integration bridge

Project-specific rules
    -> highest project-specific specialization
```

Use:

- `LordCodex/engineering-ai-rules` for universal invariants;
- `LordCodex/php-broilerplate` for PHP implementation rules;
- `LordCodex/laravel-ai-rules` for Laravel and Inertia server/integration rules;
- `LordCodex/vue-ai-rules` for Vue implementation rules;
- `bridges/laravel-inertia-vue.md` for cross-stack integration guidance that is not owned by one framework alone.

Pinned source commits are defined in `sources/rules-sources.yaml`; do not duplicate those SHAs here.

## Explicitly retired prescriptions

Do not infer any of these from the historical version of this document:

- Ziggy is not mandatory. Use the route-generation approach deliberately selected by the project.
- API Resources are not mandatory for every Inertia prop. Props must be explicitly and safely shaped; resources/DTOs are tools when useful.
- Cursor pagination is not a blanket default. Choose pagination from data size, ordering, navigation semantics, and query cost.
- Events/listeners are not required merely to separate code. Use direct application calls for synchronous required work and events when decoupled reaction semantics are appropriate.
- Queue work from actual latency/durability/retry requirements, not an arbitrary millisecond threshold.
- Redis is not a mandatory production cache/queue choice.
- Inertia lazy/deferred APIs are version-sensitive. Inspect installed versions and current adapter APIs.
- The frontend is not forbidden from all application logic; Vue may own presentation/application interaction logic, but server business invariants, authorization, and persistence authority remain in Laravel.
- Internal Inertia navigation should use Inertia-aware navigation when appropriate, but ordinary anchors remain correct for external links, downloads, and intentional full navigations.

## Compatibility behavior

PromptEngine may continue to select this document for older manifest IDs, but agents should immediately follow `bridges/laravel-inertia-vue.md` and the authoritative sources above. Do not copy framework policy back into this compatibility entry.
