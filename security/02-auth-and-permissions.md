---
document_id: security-auth-and-permissions
title: Authentication and Authorization Standards
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Authentication and Authorization Standards

## Purpose
This document establishes rules for verifying user identities (Authentication) and checking permissions (Authorization) across web and mobile services.

## Scope
Applies to routing middleware, session tracking, JWT lifecycle, and role management.

---

## Directives

### 1. Authentication Patterns

We support two primary authentication methods:

- **Stateful Cookie Sessions (Web App)**:
  - Standard for Laravel + Vue/Inertia setups.
  - **Security Rule**: Cookies must be configured with flags: `HttpOnly`, `Secure`, `SameSite=Lax` to prevent XSS session thefts and CSRF attacks.
- **Stateless Tokens (Mobile App / APIs)**:
  - Standard for Flutter communicating with Laravel APIs.
  - **Security Rule**: Access tokens (JWT/Laravel Sanctum tokens) must have short lifetimes (under 60 minutes) and use refresh tokens for renewal. Store tokens securely on mobile devices using Keychain (iOS) or Keystore (Android).

### 2. Authorization Rules (Access Control)
- **Always Default to Deny**: If a user role is undefined, default their permissions to blocked.
- **RBAC vs. ABAC**:
  - Use Role-Based Access Control (RBAC) for generic dashboard groupings (e.g. `admin`, `marketer`, `author`).
  - Use Attribute-Based Access Control (ABAC) or Laravel Policies for resource-level authorization (e.g., "A marketer can only view earnings that belong to their marketer ID").

```php
// Safe policy authorization check
public function view(User $user, Earning $earning): bool {
    return $user->id === $earning->marketer_id;
}
```

---

## Common Mistakes & Anti-Patterns
- **ID-Only Authorization**: Checking if a user is logged in, but forgetting to check if the requested resource ID belongs to them (broken object level authorization).
- **Frontend Authorization Guarding**: Hiding a button in Vue/Flutter but forgetting to verify the respective permission in the backend HTTP controller.
- **Static Credentials**: Storing passwords in database columns as plain text instead of using strong hashing algorithms like `bcrypt` or `argon2`.

---

## References
- REST API conventions: [core/03-data-and-api-modeling.md](file:///Users/kodexkode/Documents/workspace/promptengine/core/03-data-and-api-modeling.md)
- Laravel API bridge: [bridges/laravel-api-flutter.md](file:///Users/kodexkode/Documents/workspace/promptengine/bridges/laravel-api-flutter.md)
