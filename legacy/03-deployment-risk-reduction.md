---
document_id: legacy-risk-reduction
title: Rollback Planning and Deployment Risk Reduction
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Rollback Planning and Deployment Risk Reduction

## Purpose
This document defines procedures for releasing code safely, configuring feature flags, and establishing rollback pathways to protect production stability.

## Scope
Applies to deployments, feature flags integration, logging standards, and release pipeline operations.

---

## Feature Flags for Risk Isolation

When releasing high-impact business logic changes (e.g. changing the transaction checkout routing, or updating the auth provider):

- **Rule 1**: Wrap new execution paths in feature flags. Never deploy complex changes as hot paths directly.
- **Rule 2**: Keep flag checks simple and performant. Use a lightweight local check caching mechanism (e.g. Redis cache or memory variable) to prevent database call overhead on every request.

```php
if ($featureFlags->isEnabled('new-payment-gateway', $user)) {
    return $this->newPaymentService->charge($amount);
}

return $this->legacyPaymentService->charge($amount);
```

- **Rule 3**: Decommission flags within 30 days of 100% rollout to prevent code path clutter and technical debt.

---

## Rollback Planning Checklist

Every deployment plan containing database schema modifications or architectural shifts must have an explicit rollback runbook:

1. **Verify Backwards Compatibility**: Ensure that reverting to the previous git version does not cause the application to crash under the upgraded database schema (refer to [legacy/02-backward-compatibility.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/02-backward-compatibility.md)).
2. **Step-by-Step Reversion**:
   - *How is code reverted?* (e.g. `git revert` or redeploying the previous stable container tag).
   - *How are feature flags switched off?* Document the dashboard toggle or console command command to disable the code path.
3. **Database Fallbacks**: If the database schema has been migrated, verify that any new fields have nullable states or default bindings to prevent insertion errors when legacy code executes queries.

---

## Self-Verifying Logs & Metrics
- **Log Feature Health**: When introducing new logic blocks, emit structured log events containing the feature flag status, environment tags, and transaction identifiers:
  ```json
  {
    "event": "transaction_processed",
    "gateway": "stripe_v2",
    "status": "success",
    "feature_flag": "new-payment-gateway",
    "user_id": 9845
  }
  ```
- **Error Spikes**: Monitor error tracking channels (Sentry, Crashlytics, or internal log monitors). If error rates spike by more than 1% post-deployment, immediately toggle off the feature flag or execute the rollback runbook.

---

## Common Mistakes & Anti-Patterns
- **Deploying and Leaving**: Pushing a major release to production and logging off immediately without monitoring logs or checking system dashboards for 15 minutes.
- **Irreversible DB Migration**: Running a destructive database migration (e.g. deleting a column) that makes code rollbacks impossible without data restoration.
- **Flag Creep**: Retaining dead feature flags in the codebase for years, leading to cognitive fatigue when developers try to audit execution code routes.

---

## References
- Database migration safety: [legacy/02-backward-compatibility.md](file:///Users/kodexkode/Documents/workspace/promptengine/legacy/02-backward-compatibility.md)
- Safe code modularity: [core/02-architecture-and-simplicity.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/02-architecture-and-simplicity.md)
