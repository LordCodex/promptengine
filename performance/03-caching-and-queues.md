---
document_id: performance-caching-queues
title: Caching and Background Queue Workers
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Caching and Background Queue Workers

## Purpose
This document establishes rules for caching data in Redis and executing long-running operations asynchronously in queue workers to protect API response times.

## Scope
Applies to Redis caching logic, background queue worker configurations, and job payloads.

---

## Directives

### 1. Redis Caching Conventions
- **Cache Key Namespacing**: Always prefix cache keys with a structured namespace (`project:subsystem:id:property` $\rightarrow$ `playbook:user:124:profile`).
- **Define Time-To-Live (TTL)**: Never write keys with an infinite TTL unless they are handled by an automated invalidation script. Set a default TTL (e.g. 24 hours) to prevent stale cache memory bloat.
- **Cache Invalidation**: Use event-driven invalidation. When a database model updates, trigger an event listener that clears the specific cache keys immediately:
  ```php
  // Invalidating user cache on record save
  public function saved(User $user) {
      Cache::forget("playbook:user:{$user->id}:profile");
  }
  ```

### 2. Background Queue Standards
- **Keep Job Payloads Small**: Never pass large objects or instantiated database models to background jobs. Pass primitive identifiers (e.g., `user_id`, `invoice_id`). The queue worker should retrieve the active record from the database when it starts execution.
- **Why**: Serializing large database models into database/Redis queue tables consumes storage space and risks executing actions on stale data if the record changes before the worker starts processing.
- **Idempotency**: Verify that background jobs can execute multiple times without causing duplicate side effects (e.g., sending double emails or charging a credit card twice).

---

## Common Mistakes & Anti-Patterns
- **The Caching Stampede**: Setting a cache TTL so that thousands of concurrent users attempt to query the database simultaneously when the cache expires. Use locks or cache warming to prevent this.
- **Queued Model Serialization**: Passing a whole Eloquent model into a Laravel queue. If the model properties are updated after dispatching the job, the job executes with outdated values.
- **Slow Sync Hooks**: Sending emails or registering third-party tracking pixels directly within the HTTP request thread instead of dispatching them to background workers.

---

## References
- System Latency Budgets: [performance/README.md](file:///Users/kodexkode/Documents/workspace/promptengine/performance/README.md)
- Laravel Logic Layer: [stacks/php-laravel/laravel-logic.md](file:///Users/kodexkode/Documents/workspace/promptengine/stacks/php-laravel/laravel-logic.md)
