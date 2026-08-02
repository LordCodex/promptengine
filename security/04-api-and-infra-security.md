---
document_id: security-api-and-infra
title: API Security and Infrastructure Defense
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# API Security and Infrastructure Defense

## Purpose
This document defines parameters for shielding APIs from brute force requests, restricting cross-origin resource sharing, and auditing infrastructure dependencies.

## Scope
Applies to routing firewalls, CORS rules, and server configurations.

---

## Directives

### 1. API Rate Limiting
- **Rule**: Every public API route must implement rate limiting middleware to prevent Denial of Service (DoS) and authentication brute-force attacks.
- **Standards**:
  - Auth endpoints (`/api/login`, `/api/register`): Limit to **5 attempts per minute** per IP address.
  - General data endpoints: Limit to **60 requests per minute** per authenticated user token.
  - Failures must return HTTP status `429 Too Many Requests`.

### 2. CORS (Cross-Origin Resource Sharing) Policies
- **Rule**: Never set CORS configuration headers to wildcard allow patterns (`Access-Control-Allow-Origin: *`) on routes that process private user sessions or stateful cookies.
- **Vue / Inertia Web Clients**: Explicitly define the matching domains in host environment configuration variables (`ALLOWED_ORIGINS=https://app.example.com`).
- **Flutter Mobile Client Exceptions**: Mobile apps do not respect CORS flags directly, but backend endpoints must authenticate mobile requests using stateless JWT tokens instead of stateful browser cookies (refer to [security/02-auth-and-permissions.md](file:///Users/kodexkode/Documents/workspace/promptengine/security/02-auth-and-permissions.md)).

### 3. Software Bill of Materials (SBOM) and Scans
- **Rule**: Implement automated static analysis checks in CI workflows to verify dependency integrity (refer to [environment/03-ci-cd-pipelines.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/03-ci-cd-pipelines.md)).
- Fail execution when outdated dependencies with published CVE records are checked in.

---

## Common Mistakes & Anti-Patterns
- **Wildcard CORS for APIs**: Setting `Access-Control-Allow-Origin: *` to solve a CORS compile error on your Vue local environment instead of defining a specific local subdomain.
- **Uncapped API Fetches**: Providing queries that lack page/size parameters (`/api/users`), allowing clients to request millions of database records in a single call.
- **No IP Tracking**: Configuring backend rate-limits based solely on user IDs, allowing unauthorized scripts to launch brute-force attacks against login forms with no triggers.

---

## References
- REST endpoints: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)
- Flutter API integrations: [bridges/laravel-api-flutter.md](file:///Users/kodexkode/Documents/workspace/promptengine/bridges/laravel-api-flutter.md)
