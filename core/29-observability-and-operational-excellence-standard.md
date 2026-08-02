---
document_id: core-observability-and-operational-excellence
title: Observability, Monitoring, Logging, and Operational Excellence Standard
ecosystem: cross-cutting
dependencies:
  - core-frontend-architecture
  - core-frontend-security
  - core-frontend-performance
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Observability, Monitoring, Logging, and Operational Excellence Standard

## Purpose & Inheritance
This document defines core standards for system observability, logging hygiene, metric telemetry, distributed tracing, alerting policies, and production operational security. It inherits from the [Frontend Architecture Standard](23-frontend-architecture-standard.md), the [Frontend Security Standard](24-frontend-security-and-privacy-hardening-standard.md), and the [Frontend Performance Standard](26-frontend-performance-engineering-standard.md), ensuring that applications are healthy, measurable, and simple for operators to diagnose and maintain in production.

---

## 1. The Core Observability Principle

**Software cannot be safely operated if it cannot be observed.**

An application is not production-ready until operators can inspect its runtime behavior, identify bottlenecks, trace request paths across system boundaries, and locate the root cause of failures.

### The Three Pillars
Observability requires the coordination of three primary inputs:
- **Logs**: High-cardinality, contextual descriptions of discrete system events.
- **Metrics**: Aggregated numeric data points measuring system throughput, latencies, and resource utilization.
- **Traces**: End-to-end transaction paths tracking request flows through services.

---

## 2. Structured Logging & Hygiene

### Structured Logging
All application logs must be structured (typically formatted as JSON objects) to allow search engine filtering. Avoid free-form, plain text log strings.

Every log record should capture:
- **Timestamp**: Pinned to ISO 8601 in UTC.
- **Correlation & Request IDs**: Unique string identifiers traced across microservices.
- **User Context**: Anonymized or database User IDs (never log email addresses or raw names).
- **Service Info**: Application name, container version, and active environment (`production`, `staging`).
- **Operation**: Active route, controller class method, or queue job namespace.
- **Severity**: Appropriate log level.
- **Error Context**: Error codes, validation exception mappings, and sanitized messages.

### Log Levels
Use log levels consistently. Never log everything under the same level:
- **DEBUG**: Low-level diagnostics; disabled in production environments.
- **INFO**: Standard system milestones (e.g., job completed, system booted).
- **WARNING**: Unexpected events that are automatically recovered (e.g., database connection retry success, fallback API call).
- **ERROR**: Failures that block the active request or background job but do not crash the service.
- **CRITICAL**: Major system-threatening events requiring immediate operator pager alerts (e.g., primary database unavailable, storage disk depleted, critical webhook signature failures).

---

## 3. Sensitive Data Redaction

Logging platforms are primary vectors for accidental data leaks.

### Prohibited Log Parameters
Never write the following data points to logs, metrics, trace attributes, or dashboard panels:
- Passwords, recovery pins, and secure OTP codes.
- API keys, private tokens, webhook secrets, and database credentials.
- Credit card numbers, expiration dates, CVVs, or payment codes.
- Session tokens, JWT strings, and Authorization HTTP headers.
- Personally Identifiable Information (PII) including social security numbers, government IDs, and unredacted medical fields.

### Redaction Enforcers
- Configure middleware or logging adapters to scrub common parameter names (e.g. `password`, `token`, `card_number`) using sanitization scripts.
- Mask or redact values server-side before writing payloads to storage output streams.

---

## 4. Request Correlation & Distributed Tracing

### Request Correlation
- Every incoming HTTP request or queue job initiation must generate a unique `Request-ID` and `Correlation-ID`.
- Propagate correlation headers to all downstream services, external API requests, database queries, and background queues.
- Ensure all log records generated during a single request workflow inherit the same correlation identifier to simplify search grouping.

### Distributed Tracing
Trace critical transaction paths across system boundaries:
```text
[ Incoming API Gateway ]
          ↓ (injects Trace ID)
[ Controller Handler ]
          ↓ (propagates Trace ID)
[ Business Service / Domain Model ]
          ↓
[ Database Query / External API Call ]
```
Distribute trace context across HTTP headers using standard formats (e.g., W3C Trace Context or Jaeger headers).

---

## 5. Health Checks & Diagnostics

Expose dedicated health API endpoints to orchestrate automated scaling, load balancer routing, and container health checks:
- **Liveness Endpoints (`/health/live`)**: Confirms the runtime container is executing. Returns a lightweight `200 OK` response without querying databases or third-party APIs.
- **Readiness Endpoints (`/health/ready`)**: Confirms the container is ready to accept user traffic. Verifies database connectivity, memory capacity, and active cache availability.
- **Dependency Health**: Periodically check key external integrations in the background; do not execute slow, synchronous dependency checks on every load balancer readiness probe.

---

## 6. Metric Collection Boundaries

Gather operational metrics to capture trends:
- **System Metrics**: CPU load, RAM usage, storage consumption, and socket connections.
- **Application Latencies**: P95 and P99 HTTP response times, database query execution durations, and external API latency.
- **Operational Volumes**: Total request count, error rates (5xx counts), queue depth backlogs, and background job processing times.
- **Caches**: Cache hit-to-miss ratios.

---

## 7. Incident Management & Error Tracking

### Centralized Error Reporting
- Integrate centralized error reporting tools (e.g., Sentry, Bugsnag, Rollbar) to capture exceptions.
- Gather debugging metadata: stack traces, request route, active release commit SHA, and environment config variables.
- Group duplicate errors to prevent alert fatigue.

### Deployment Visibility
Record metadata on every production build deployment:
- Commit SHA, version, build number, and deployment timestamp.
- Integrate deployment tracking with error reporting tools to identify which release introduced a new exception.

---

## 8. Feature Flag Governance

Use feature flags to separate code deployment from feature release:
- **Short Lifecycles**: Feature flags are temporary tools for risky rollouts. Define an owner and an expiration target (typically within 4–6 weeks) for every flag.
- **Clean Up Code**: Remove flag checks and conditional branches from templates and controllers immediately after a successful rollout.
- **Do Not Use for Configuration**: Never use feature flags as a permanent configuration engine.

---

## 9. Background Job and Queue Monitoring

Background tasks must be measured to identify failures:
- **Telemetry**: Track processing durations, queue backlog depths, job retry frequencies, and dead-letter queue (DLQ) size.
- **Dead-Letter Queues**: Configure DLQs to capture jobs that exceed maximum retry counts. Never discard failed jobs silently.
- **Progress Reporting**: Long-running asynchronous tasks must report execution progress to prevent operators from assuming a task is frozen.

---

## 10. Audit Logging (Business Compliance)

Audit logs track critical business actions and are structurally separate from system diagnostics:
- **Scope**: Record security-sensitive activities: login attempts, role/permission updates, password changes, financial transactions, and configuration updates.
- **Payload**: Log *Who* acted, *What* changed (original vs new values), *When* it occurred, and *From Where* (IP address and client type).
- **Immutability**: Audit records must be write-once, read-many. Store audit logs in dedicated, tamper-resistant tables or write them to append-only log locations.

---

## 11. Alerting Strategy

Avoid alert fatigue by defining alerts based on actionable criteria:
- **Critical Alerts (Pager)**: High error rate thresholds (5xx counts exceeding 1% of total traffic), primary database unavailable, background queues backed up past SLA thresholds.
- **Warning Alerts (Slack/Email)**: Disk usage exceeding 80%, certificate expirations within 30 days, minor webhook delivery failures.
- **Alert Rule**: Only trigger critical pager alerts if immediate operator action is required.

---

## 12. Operational Tooling & Dashboard Security

- **Least Privilege Access**: Apply the principle of least privilege to SRE dashboards, logging search tools, and cloud provider consoles.
- **No Secrets in Tools**: Ensure no secrets, private credentials, or raw client records are exposed in dashboard panels, build logs, or CI/CD pipelines.

---

## 13. Framework Specific Guidelines

### 13a. Laravel
- Use Laravel's logging configuration to route output to JSON log processors.
- Monitor Redis queues and dead-letter pipelines using Horizon.
- Enable Telescope only in development environments; disable it in production to prevent CPU and storage exhaustion.

### 13b. Vue / Nuxt
- Capture client-side exceptions and route them to your error tracking platform (e.g. Sentry Vue SDK).
- Capture performance telemetry (TTFB, LCP, INP) using browser metrics APIs.
- Mask all user input strings before sending client crash payloads.

### 13c. Flutter
- Log crashes and App Not Responding (ANR) exceptions.
- Gather screen navigation times and network latency metrics while respecting user privacy rules and device permissions.

---

## Review Checklist

Before declaring a feature ready for production, verify against this checklist:
- [ ] **Structured logging implemented**: Are all logs written in JSON format with correlation/request IDs?
- [ ] **Sensitive data protected**: Are passwords, access tokens, PII, and credentials redacted or masked?
- [ ] **Request correlation available**: Do headers propagate correlation and trace identifiers to all services?
- [ ] **Health endpoints considered**: Are liveness (`/health/live`) and readiness (`/health/ready`) endpoints configured?
- [ ] **Metrics identified**: Are request latency, error rates, queue depths, and DB queries tracked?
- [ ] **Error tracking configured**: Is centralized exception reporting active with release tracking?
- [ ] **Audit logging reviewed**: Are critical business and security actions saved in immutable audit records?
- [ ] **Feature flags justified**: Do flags have designated owners, expiration dates, and rollback plans?
- [ ] **Alerts defined**: Are pager notifications configured only for actionable system failures?
- [ ] **Dashboards considered**: Do metric dashboards support rapid diagnosis?
- [ ] **Deployment version identifiable**: Do build files compile with the active commit SHA and timestamp?

---

## References
- Frontend Security: [core/24-frontend-security-and-privacy-hardening-standard.md](core/24-frontend-security-and-privacy-hardening-standard.md)
- Frontend Performance: [core/26-frontend-performance-engineering-standard.md](core/26-frontend-performance-engineering-standard.md)
- OpenTelemetry: [https://opentelemetry.io](https://opentelemetry.io)
- Sentry Documentation: [https://docs.sentry.io](https://docs.sentry.io)
