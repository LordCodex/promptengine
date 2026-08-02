---
document_id: bridge-cors-policies
title: Cross-Origin Resource Sharing (CORS) Standards
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Cross-Origin Resource Sharing (CORS) Standards

## Purpose
This document defines standards for configuring CORS filters on backend APIs to allow safe communication with authorized browser clients while blocking unauthorized cross-origin requests.

## Scope
Applies to routing configurations, middleware parameters, and proxy integrations.

---

## Directives

### 1. Enforce Origin Whitelists
- **Rule**: Never expose wildcard parameters (`*`) on endpoints that process session credentials or private user authentication headers.
- **Action**: Explicitly parse authorized frontend clients domains from configuration variables:
  ```php
  // config/cors.php
  'allowed_origins' => explode(',', env('ALLOWED_ORIGINS', 'https://app.example.com')),
  ```

### 2. Allowed Header Mappings
- Restrict HTTP headers to the absolute minimum required for operations:
  ```json
  Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With, X-Inertia
  ```
- **X-Inertia Header**: For Inertia integrations, ensure the `X-Inertia` header is added to the allowed header array to prevent browser redirection blockages.

### 3. Mobile App Client Exemptions
- Mobile apps (Flutter/Dart clients) run in environments that do not respect browser-level CORS parameters directly.
- **Rule**: Do not bypass CORS configurations for web clients just because mobile client calls operate cleanly. CORS configs must remain strictly enabled for web browsers.

---

## Common Mistakes & Anti-Patterns
- **Wildcard Allow with Credentials**: Setting `Access-Control-Allow-Origin: *` while simultaneously setting `Access-Control-Allow-Credentials: true`. Most modern browsers block this configuration, preventing connections completely.
- **Reflecting Origin Header**: Dynamically reading the incoming request `Origin` header and writing it back directly into the `Access-Control-Allow-Origin` response header. This effectively disables CORS protection.
- **Pre-flight Request Failures**: Failing to handle OPTIONS requests in routers, causing browsers to abort actual POST/GET queries after pre-flight request timeout.

---

## References
- Authorization policies: [security/02-auth-and-permissions.md](file:///Users/kodexkode/Documents/workspace/promptengine/security/02-auth-and-permissions.md)
- REST endpoints security: [security/04-api-and-infra-security.md](file:///Users/kodexkode/Documents/workspace/promptengine/security/04-api-and-infra-security.md)
