---
document_id: core-performance-engineering-standard
title: Performance Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-security-testing-and-threat-modeling
  - stacks-php-conventions
  - stacks-laravel-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Performance Engineering Standard

## Purpose & Inheritance
This document defines the core standards for profiling, measuring, and optimizing software performance. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), the [Database Engineering Standard](06-database-engineering-standard.md), and the [API Engineering Standard](07-api-engineering-standard.md). It outlines practical guidelines for resource utilization across backend frameworks (PHP/Laravel), frontends (Vue/Nuxt), mobile applications (Flutter/Dart), and server infrastructure.

---

## 1. Performance Philosophy

Performance is the direct result of disciplined engineering decisions. We reject premature optimization (fine-tuning code blocks that do not impact user response times) as well as ignoring efficiency (writing sloppy code and relying on vertical hardware scaling).

### Core Optimization Priorities
1. **User Experience First**: Optimization targets must focus on metrics that directly impact user perception: page load speed, interaction response times, and API query completion latency.
2. **Never Guess, Measure First**: Do not refactor code based on assumptions of what is slow. Optimize only after capturing concrete execution metrics (CPU, memory, SQL execution traces).
3. **Optimize the Correct Layer**: Address architecture and database query flaws first. Micro-optimizing individual language algorithms (e.g., swapping double quotes for single quotes in PHP) yields negligible gains if the system is blocked by a slow, unindexed database query.
4. **Maintain Readability**: Code readability and maintainability must not be sacrificed for micro-second gains unless the target block is a verified high-frequency hot path (executed millions of times per minute).

---

## 2. Performance Engineering Process

We enforce a repeatable, five-step workflow for all performance interventions:

```text
Define SLA Goals ──> Measure & Profile ──> Identify Bottlenecks ──> Apply Optimizations ──> Verify & Compare
```

### Step 1: Define Requirements
Establish Service Level Agreements (SLAs) and targets:
- **Throughput**: Target requests per second (RPS) the endpoint must support.
- **SLA Bounds**: Response latency thresholds (e.g., $P_{95} < 200\text{ms}$, $P_{99} < 500\text{ms}$).
- **Resource Constraints**: Maximum CPU, memory, and database connection limits.

### Step 2: Measure
Capture baseline metrics using production-equivalent datasets:
- **Runtimes**: Memory footprint, CPU utilization, garbage collection intervals.
- **Data Layers**: Number of SQL queries executed per HTTP request, average query latency.
- **Network bounds**: Time to First Byte (TTFB), payload sizes, DNS resolution latency.

### Step 3: Identify Bottlenecks
Locate the primary execution blockers:
- Analyze execution flame graphs to find slow functions.
- Run `EXPLAIN` query plans to catch unindexed database scans.
- Scan for memory leaks (e.g., retained references in singleton classes or background loops).

### Step 4: Optimize
Execute modifications in hierarchical order of impact:
1. **Architecture Optimizations**: Introduce async processing, database replicas, or cache layers.
2. **Database Optimizations**: Add indexes, fix N+1 queries, introduce cursor pagination.
3. **Algorithm Optimizations**: Implement generators, optimize loop nesting levels.
4. **Network Optimizations**: Enable Brotli compression, decrease payload footprints.
5. **Micro-Optimizations**: Optimize runtime configurations (OPcache buffer sizes).

### Step 5: Verify
- Rerun benchmarks under the exact same conditions.
- Compare memory usage and response times to ensure no regressions were introduced.

---

## 3. Application & Runtime Profiling

### PHP Performance Optimizations
- **OPcache Configuration (Production)**: Compile bytecode to memory to avoid disk-read overhead:
  ```ini
  opcache.enable=1
  opcache.validate_timestamps=0 ; Never check file changes on disk in production
  opcache.memory_consumption=256
  opcache.interned_strings_buffer=16
  ```
- **Composer Optimizations**: Always generate optimized, class-map autoloader files for production:
  ```bash
  composer dump-autoload --optimize --no-dev --classmap-authoritative
  ```
- **Garbage Collection (GC) footprint**: Disable garbage collection inside heavy, multi-million record background CLI import jobs to reduce CPU overhead. Call `gc_collect_cycles()` manually at safe batch intervals.

### Laravel Performance Optimizations
- **Compile Caches**: Cache configurations, routes, and views in your production deployment pipeline:
  ```bash
  php artisan config:cache
  php artisan route:cache
  php artisan view:cache
  ```
- **Query Logging Checks**: Ensure raw query listeners and debug bars (e.g., Laravel Debugbar) are completely disabled in production environments. Storing SQL trace statements in memory can lead to Out of Memory (OOM) crashes.
- **Horizon Queue Monitoring**: Monitor worker health and queue latency using Laravel Horizon. If a queue's latency increases, scale up background workers.

---

## 4. Database Performance Rules

Database access is the most common bottleneck in web applications.

### Query Analysis
- **EXPLAIN Analysis**: Never write raw query loops. Use `EXPLAIN` to verify that queries use index scans instead of full table scans (seq scans).
- **Eager Loading**: Prevent N+1 queries by eager loading all relational dependencies:
  ```php
  // Good: Performs 2 queries instead of 101 queries
  $invoices = Invoice::with('customer')->limit(100)->get();
  ```

### Large Data Streaming
- **Avoid get() on Large Sets**: Do not fetch thousands of rows into memory using `get()`. Use `lazy()` (generators) or `chunkById()` to stream chunks:
  ```php
  // Good: Streams records in chunks of 500 to keep memory usage low
  Invoice::chunkById(500, function ($invoices) {
      foreach ($invoices as $invoice) {
          $this->process($invoice);
      }
  });
  ```
- **Cursor Pagination**: Enforce keyset/cursor pagination for large collection routes (e.g., `/api/v1/logs`) to prevent database slow offsets.

---

## 5. Caching Strategy & Cache Design

Caching reduces database read load, but incorrect invalidation introduces structural bugs.

### Cache Scoping Matrix

```text
HTTP Request
  ├── CDN Cache (Returns static assets instantly)
  └── Application Cache (Redis)
        ├── Found? ──> Return Cached JSON/Object
        └── Not Found? ──> Query Database ──> Write to Redis ──> Return Response
```

- **What to Cache**: Static configuration variables, computed analytics statistics, external API responses, and frequently accessed read-only profiles (e.g., countries, roles).
- **What NOT to Cache**: Rapidly changing transactional records (e.g., account ledger balances, real-time inventory levels, user security permissions) where stale values introduce business logic or security vulnerabilities.

### Cache Keys & Invalidation
- **Predictable Keys**: Namespace keys using a colon-separated schema: `{domain}:{identifier}:{attribute}` (e.g., `user:1024:profile`).
- **Explicit TTL (Time To Live)**: Always define a maximum TTL for cached items. Never write cached values without expiration limits.
- **Stampede Prevention**: Use cache locks or return placeholder values during cache rebuilding to prevent cache stampedes (where concurrent requests overload the database when a popular cache key expires).

---

## 6. Redis Integration

Redis is a high-performance, in-memory data store. We use Redis for caching, session storage, rate limiting, and queues.

### Redis Operational Rules
- **Connection Pools**: Configure connection pooling in high-concurrency environments to prevent socket exhaustion.
- **Key Eviction Policies**: Enforce a strict eviction policy on cache instances (e.g., `maxmemory-policy volatile-lru`) to prevent Redis from crashing when memory limits are reached.
- **Memory Footprint**: Do not store massive object payloads in Redis. Store only primary attributes or serialized DTO strings.

---

## 7. Queues & Asynchronous Processing

Offloading operations to background queues is essential to keep HTTP response times low.

### Queue Best Practices
- **Queueable Tasks**: Move operations that take longer than `50ms` (e.g., sending emails, generating PDFs, third-party API calls) off the HTTP thread.
- **Idempotency**: Background jobs must be idempotent. If a job is executed multiple times due to a network timeout, the application state must remain correct.
- **Dead Letter Queues (DLQ)**: Configure jobs with a maximum retry limit (e.g., `tries=3`) and route repeatedly failing jobs to a dead letter queue for analysis.

---

## 8. Concurrency & Locking

Concurrency increases system throughput, but introduces race conditions.

### Race Condition Protections
- **Database Transactions**: Enforce database transaction boundaries around multi-table updates.
- **Atomic Locks**: Use distributed locks (e.g., Redis locks) to prevent multiple processes from modifying the same resource concurrently:
  ```php
  $lock = Cache::lock('process-invoice-102', 10);
  
  if ($lock->get()) {
      // Safe execution context
      $lock->release();
  }
  ```

---

## 9. API & Frontend/Mobile Performance

### API Performance
- **Response Transformation**: Hydrate data payloads into lightweight DTOs or Resource maps. Avoid returning raw Eloquent models containing unnecessary schema metadata.
- **Payload Compression**: Enforce Brotli or Gzip compression for all JSON API responses.
- **API Caching**: Set HTTP caching headers (`Cache-Control: private, max-age=60`) for user-specific GET resources.

### Frontend Performance (Vue 3 / Nuxt 3)
- **Code Splitting**: Lazy load routes and components to keep bundle sizes small.
- **SSR Hydration Limits**: Keep server-side rendering (SSR) data payloads small. Avoid hydrating client states with massive data objects that are not rendered in the view.
- **Reactivity Optimization**: Avoid storing large, static datasets (e.g., parsed reports) in reactive state variables (`ref` or `reactive`). Use non-reactive objects (`shallowRef`) to avoid reactivity tracking overhead.

### Mobile Performance (Flutter / Dart)
- **Avoid Widget Rebuilds**: Keep widgets focused. Use `const` constructors where possible to prevent redundant widget tree rebuilds.
- **Image Optimization**: Never download raw, full-size images. Fetch optimized, resized assets matching target layout resolutions.
- **State Management footprint**: Use lightweight state libraries (e.g., Riverpod, Bloc) and scope rebuilds to the smallest widget nodes possible.

---

## 10. Infrastructure & Observability

### Scaling Strategy
- **Horizontal Scaling**: Add more application server instances (behind a load balancer) to scale CPU/memory bounds. Best for stateless web nodes.
- **Vertical Scaling**: Upgrade the CPU, RAM, or storage disk speeds (IOPS) of a single server. Best for database instances where clustering introduces high complexity.

### Observability Metrics
Enforce the collection of core metrics:
- **SLA Telemetry**: Track request rate, error rate, and duration (RED method).
- **APM Profiling**: Implement Application Performance Monitoring (e.g., OpenTelemetry, Datadog) to trace request lifecycles across services.
- **Slow Query Alerts**: Configure database engines to log queries taking longer than `100ms` and route alerts to engineering Slack/PagerDuty channels.

---

## 11. Performance Anti-Patterns

- **Premature Caching**: Adding cache layers to slow queries instead of optimizing query indexes. Caching slow queries hides the database problem and leads to data consistency issues.
- **Micro-Optimizing Clean Code**: Refactoring readable algorithms for micro-second gains at the cost of maintainability.
- **Overusing Queues**: Offloading trivial, lightweight operations (taking `<5ms`) to queues, introducing unnecessary infrastructure overhead and latency.
- **Ignoring Payload Sizes**: Returning hundreds of rows of un-paginated data to API consumers, saturating bandwidth and memory on both client and server.

---

## 12. Decision Matrices

Use these matrices to identify the correct performance decision based on project context.

### Matrix 1: Cache vs. No Cache
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Read-heavy queries, static configurations, computed values | **Cache** | Eliminates database network overhead for repetitive reads. |
| Financial ledgers, real-time transaction states, user privileges | **No Cache** | Guarantees data correctness and prevents authorization bypasses. |

### Matrix 2: Redis vs. Database
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Temporary rate limiting counters, sessions, low-value queues | **Redis** | In-memory speed; provides fast read/write operations for transient data. |
| Relational data, billing history, persistent logs | **Database** | Enforces data integrity, relationships, and ACID guarantees. |

### Matrix 3: Queue vs. Immediate Processing
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Background reports exports, sending welcome emails, bulk updates | **Queue** | Prevents blocking HTTP request threads, improving response speeds. |
| Deducting inventory balances, validating user credentials | **Immediate** | Confirms critical business constraints before returning HTTP replies. |

### Matrix 4: Vertical Scaling vs. Horizontal Scaling
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Database instances, systems requiring strict data consistency | **Vertical Scaling**| Simple to scale; avoids database clustering complexity. |
| Stateless web/API application servers | **Horizontal Scaling**| Cost-effective; increases system redundancy and availability. |

### Matrix 5: SQL Optimization vs. Caching
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Unindexed slow queries, N+1 query loops | **SQL Optimization** | Solves the root database bottleneck; prevents cache consistency issues. |
| Expensive queries that are already indexed but run frequently | **Caching** | Reduces CPU read load on the database engine. |

### Matrix 6: SSR vs. Client Rendering (Vue/Nuxt)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Public-facing marketing pages, e-commerce catalog pages | **SSR (Nuxt)** | Optimizes Search Engine Optimization (SEO) and First Contentful Paint (FCP). |
| Internal dashboard panels, authenticated SaaS work surfaces | **Client Rendering** | Reduces server CPU load; client handles view rendering directly. |

### Matrix 7: Pagination vs. Loading Everything
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Table lists containing more than 50 rows | **Pagination** | Limits memory consumption on both server and client layers. |
| Small, static dropdown config options (e.g., list of 10 statuses) | **Load Everything** | Saves query complexity and client-side code overhead. |

### Matrix 8: Optimization vs. Simplicity
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Core request workflow loops executing millions of times per day | **Optimization** | Reduces infrastructure costs and database connection saturation. |
| Standard CRUD routes, low-frequency internal administration panels | **Simplicity** | Minimizes development costs and keeps code readable for maintenance. |

---

## 13. AI Performance Directives

AI agents modifying code in this repository must follow these rules:

1. **Eager Load by Default**: When fetching relation dependencies in ORM models, ensure eager loading is configured (preventing N+1 queries).
2. **Never Suggest Caching Without Invalidation**: When introducing cache structures, always write corresponding invalidation logic (e.g., database observers or events).
3. **No All() or Get() on Collection Lists**: Enforce paginations or chunks on all collection database operations.
4. **Use Prepared Statements**: Do not suggest variable concatenation inside queries. Always use prepared statement parameters.
5. **No Micro-Optimizations**: Do not modify clean, readable code blocks for performance unless profiling logs verify that the block is a primary bottleneck.

---

## 14. Performance Review Checklist

Use this checklist during code review to evaluate performance modifications.

### Database & Queries
- [ ] Have you verified that eager loading is configured (no N+1 queries)?
- [ ] Have you run an `EXPLAIN` query plan to confirm index usage?
- [ ] Do large database updates process records in batches or chunks?

### Memory & Caching
- [ ] Are memory-heavy background import scripts streaming rows using generators?
- [ ] Do cached values specify an explicit TTL (Time To Live)?
- [ ] Have invalidation handlers been implemented for all cached structures?

### Queues & Concurrency
- [ ] Are operations taking longer than `50ms` offloaded to background queued jobs?
- [ ] Do queued jobs receive database IDs instead of fully instantiated models?
- [ ] Are race-condition prone balances updates protected by transactions or atomic locks?

### Frontend & Mobile
- [ ] Are expensive Vue state variables shallow-referenced to avoid reactivity tracking?
- [ ] Are routes and components configured to lazy load?
- [ ] Do Flutter layouts use `const` widget constructors to prevent redundant rebuilds?

### Infrastructure & Monitoring
- [ ] Are application route caching, config caching, and view caching commands configured in the deployment pipeline?
- [ ] Are slow SQL queries logged and monitored?

---

## 15. Pagination Strategy

Choosing the correct pagination strategy is a scalability decision, not a convenience decision.

| Strategy | Use Case | Limitation |
| :--- | :--- | :--- |
| **Offset Pagination** (`LIMIT x OFFSET y`) | Bounded interactive UI pages (user-facing tables, search results). | Performance degrades on large offsets; avoid for internal batch processing. |
| **Keyset Pagination** (`WHERE id > last_seen_id ORDER BY id`) | Internal batch jobs, background exports, cleanup workers, and API cursors. | Requires a stable, indexed sort key; cannot jump to arbitrary pages. |

### Rules
- Use **keyset pagination** for all internal batch processing, background workers, data exports, and cleanup jobs. Offset pagination scans all preceding rows on each page request, causing quadratic performance degradation at scale.
- Use **offset pagination** only for bounded, interactive UI pages when legacy behavior requires it, and when page sizes are small and fixed.
- Do not run a filtered-count query when no filter is applied if the total count is already known.
- Parse numeric identifiers as identifiers — do not cast indexed ID columns to strings for wildcard search.
- Every cleanup or bulk mutation must have a deterministic order and a maximum number of rows processed per run.

---

## 16. Dangerous Search Patterns

The following query patterns bypass indexes and should be treated as measurable scalability risks:

| Pattern | Risk | Mitigation |
| :--- | :--- | :--- |
| **Leading wildcard** (`LIKE '%term'`) | Full table scan on all rows. | Use full-text search indexes or prefix-only search. |
| **Concatenated column search** (`CONCAT(first_name, ' ', last_name) LIKE ?`) | Computed expression bypasses column index. | Index both columns separately; search each independently. |
| **Broad OR clauses** (`WHERE col1 = ? OR col2 = ? OR col3 = ?`) | Can prevent index usage depending on cardinality. | Use `UNION` or full-text search for multi-field searches. |
| **Large-text column search** (`WHERE body LIKE '%term%'`) | Full table scan on potentially large text fields. | Use a dedicated search index (e.g., Meilisearch, Elasticsearch). |

### Rules
- Do not remove required legacy search behavior without explicit approval. Instead, document the performance tradeoff and add bounded input length limits and page limits.
- When introducing a new search feature, confirm which index it uses before shipping.
- Treat any search pattern from the table above as a P2 performance risk requiring measurement before deployment at scale.

---

## References
- Database Optimizations: [06-database-engineering-standard.md](06-database-engineering-standard.md)
- API Optimization: [07-api-engineering-standard.md](07-api-engineering-standard.md)
- PHP OPcache Tuning: [stacks/php-laravel/php-conventions.md](../stacks/php-laravel/php-conventions.md)
- Laravel Routing Optimization: [stacks/php-laravel/laravel-routing.md](../stacks/php-laravel/laravel-routing.md)

