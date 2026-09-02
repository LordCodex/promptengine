---
document_id: bridge-laravel-inertia-vue
title: Laravel Inertia Vue Integration Bridge
ecosystem: bridge
audience: [human, agent]
last_reviewed: 2026-09-02
---

# Laravel + Inertia + Vue Integration Bridge

## Purpose

This bridge contains only the rules that matter because Laravel, Inertia, and Vue are used together. It is not a second copy of the framework standards.

Resolve the authoritative sources first:

1. universal engineering rules;
2. PHP implementation rules;
3. Laravel rules, including `docs/INERTIA.md`;
4. Vue rules;
5. this bridge;
6. project-specific instructions and decisions.

PromptEngine should select task-relevant documents from those sources rather than concatenating every repository.

## Authority boundaries

Laravel remains authoritative for:

- routes and middleware;
- authentication/session state;
- authorization/policies;
- validation;
- transactions and persistence;
- queues/jobs/events;
- server-side business invariants.

Inertia owns the integration contract between Laravel responses and Vue pages:

- page/component resolution;
- page props and shared props;
- navigation and redirects;
- form/validation-error transport;
- partial/deferred/lazy data mechanisms supported by the installed Inertia version;
- optional SSR integration.

Vue owns presentation implementation:

- components;
- composables;
- reactivity;
- local/client state;
- frontend accessibility and rendering behavior;
- Vue tests.

A frontend permission check, hidden button, page prop, or route guard is UX only. Laravel authorization remains mandatory.

## Shared and page props

Keep globally shared props small and genuinely global. Feature-specific data belongs on the page/feature that needs it.

Treat props as a public presentation contract:

- send only data required by the UI;
- explicitly shape sensitive or wide model data;
- avoid accidental relationship expansion and N+1 queries;
- avoid secrets, internal-only columns, or unrelated user data;
- use DTO/resource-like shapes when they improve safety or stability, but do not require one wrapper type for every prop.

Use deferred, lazy, optional, or partial-reload mechanisms only when supported by the installed Inertia version and when they materially reduce unnecessary work. They do not excuse inefficient queries.

## Forms and validation

Prefer the Inertia form primitives supported by the installed adapter when they simplify processing state, errors, resets, and navigation.

Laravel validation remains authoritative. Normal Inertia validation follows Laravel redirect/session error semantics rather than requiring a JSON `422` API contract.

Client-side validation may improve UX, but it must not replace server validation, authorization, idempotency, or transaction rules.

Duplicate-sensitive submissions such as payments, checkout, orders, or one-time actions need server-side idempotency/deduplication where required; disabling a Vue button is not enough.

## Navigation and routing

Laravel routes remain the server routing source of truth. Do not add a second full SPA router merely to duplicate Laravel navigation.

Use Inertia-aware navigation for internal Inertia page transitions when appropriate. A normal anchor remains correct for external links, downloads, non-Inertia destinations, or cases where a full browser navigation is intentionally required.

Route-generation helpers such as Wayfinder, Ziggy, generated TypeScript helpers, or project-specific route utilities may be used when deliberately selected. This bridge does not mandate one helper.

## State ownership

Do not copy every server prop into a Vue global store.

Prefer:

- page props for server-owned page state;
- URL/query state for shareable navigation/filter state;
- component/composable state for local interaction;
- Pinia only for justified durable cross-page client state.

When account, tenant, session, or permissions change, clear or refresh stale client state so privileged data or capabilities are not retained incorrectly.

## Feature boundaries

Keep backend and frontend feature ownership aligned where practical:

```text
Laravel feature -> same feature / Shared
Vue feature     -> same feature / Shared
Shared          -> feature-neutral only
```

Do not import another feature's private Vue components, composables, stores, or services merely because the bundler permits it. Do not use Inertia props/events as a hidden substitute for an intentional cross-feature server boundary.

## Performance and failure behavior

Review the complete request path:

- database query count/N+1;
- serialized prop size;
- shared prop size;
- partial/deferred request behavior;
- duplicate requests;
- Vue render/reactivity cost;
- bundle/chunk size;
- external dependency failures;
- session expiry and authorization changes.

Do not require cursor pagination, queues, Redis, SSR, or deferred props by default. Choose them from actual requirements and measured cost.

## SSR

Inertia SSR is optional. Enable it only when SEO, first-render, or product requirements justify the extra Node/runtime deployment and operational surface.

If enabled, ensure request-specific user state cannot leak across SSR requests and that the SSR process is supervised, observable, deployable, and restartable.

## Testing

For meaningful Laravel + Inertia + Vue flows, combine the appropriate layers:

- Laravel feature tests for route access, policies, validation, redirects/session behavior, Inertia component names, and important prop contracts;
- Vue tests for component behavior and client interaction;
- critical-flow tests for expired sessions, forbidden access, validation failures, account/tenant switching, mutation failure, and duplicate-sensitive submissions where relevant.

Frontend tests do not prove backend authorization.

## Version awareness

Inspect the installed Laravel, Inertia server adapter, `@inertiajs/vue3`, and Vue versions before using version-specific APIs. Do not copy old examples such as deprecated lazy-prop APIs blindly.

## Bridge checklist

Before completing Laravel + Inertia + Vue work, verify:

- server authority remains in Laravel;
- props are minimal and safe;
- validation/authorization are enforced server-side;
- navigation uses the correct Inertia/full-browser behavior;
- state is owned at the narrowest correct layer;
- feature boundaries are preserved;
- no route-helper package is assumed without project adoption;
- version-specific Inertia APIs match installed versions;
- performance/failure behavior was reviewed;
- tests/checks are reported only if they actually ran.
