---
document_id: checklists-security-review
title: Security Review Checklist
ecosystem: cross-cutting
dependencies:
  - core-security-engineering-standard
  - core-api-engineering-standard
  - core-database-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Security Review Checklist

Apply this checklist during code review and before deploying any endpoint that handles user input, authentication, authorization, or financial operations.

---

## The Three Questions (Every Endpoint, Every Time)

| # | Question | Resolved? |
| :--- | :--- | :--- |
| 1 | **Who sent this?** → Authentication configured? | `[ ]` |
| 2 | **Are they allowed?** → Authorization enforced? | `[ ]` |
| 3 | **Is the data safe?** → Input validated and output escaped? | `[ ]` |

If any answer is NO, do not ship the endpoint.

---

## Authentication & Session Management
- [ ] Every protected route has authentication middleware explicitly applied.
- [ ] Sessions are regenerated after login to prevent session fixation.
- [ ] Session cookies are configured as `HttpOnly`, `Secure`, `SameSite=Lax` or `Strict`.
- [ ] Password hashing uses bcrypt or Argon2id — never MD5, SHA1, or plain storage.
- [ ] Token comparison uses constant-time functions (`hash_equals()`, `timingSafeEqual()`) — never `==` or `===`.
- [ ] Token generation uses cryptographically secure functions (`random_bytes()`, `crypto.randomUUID()`) — never `rand()`, `mt_rand()`, or `Math.random()`.
- [ ] OTP and magic link tokens are single-use and expire after a defined window.
- [ ] The login error message is the same whether the email does not exist or the password is wrong (no account enumeration).

---

## Authorization
- [ ] Every resource has an ownership check (not just a role check).
- [ ] IDOR is prevented — a user cannot access another user's data by changing an ID in the URL or request body.
- [ ] BOPLA is addressed — even if a user owns the resource, they cannot modify privilege-escalation fields (`role`, `balance`, `is_admin`, `status`, `permissions`).
- [ ] Field-level write permissions are validated per role in a centralized policy.
- [ ] Multi-step workflows cannot be skipped by hitting a later endpoint directly.

---

## Input Validation
- [ ] All input is validated server-side using a strict whitelist — allowed fields, types, lengths, formats, and ranges.
- [ ] Duplicate parameters in query strings and JSON bodies are rejected (HTTP Parameter Pollution).
- [ ] `$request->all()` or equivalent is never passed directly into database inserts or updates.
- [ ] Mass assignment is protected — `$fillable` is defined explicitly; `$guarded = []` is never used.
- [ ] Quantity, price, and amount fields have server-side min and max limits.
- [ ] Prices and amounts are recalculated server-side — never trusted from the frontend.
- [ ] File uploads are validated by MIME type (not just extension), renamed server-side, and stored outside `public/`.

---

## Output & XSS Prevention
- [ ] All user inputs are escaped in the correct context before rendering (HTML, attribute, JavaScript).
- [ ] User-facing messages are inserted via DOM text APIs (`textContent`, `innerText`) — never `innerHTML` or string interpolation.
- [ ] Non-2xx async responses define explicit error handling — no raw server errors rendered to the browser.
- [ ] CSV/spreadsheet exports sanitize cells beginning with `=`, `+`, `-`, or `@`.

---

## Database & SQL
- [ ] All queries use parameterized statements — no string concatenation of user input into SQL.
- [ ] Financial operations use `SELECT ... FOR UPDATE` (pessimistic lock) before reading and writing balances.
- [ ] Financial mutations are wrapped in a database transaction.
- [ ] Idempotency keys are implemented for payment operations to prevent duplicate charges.

---

## Race Conditions
- [ ] Balance checks, booking availability, and inventory checks acquire a lock before reading and writing.
- [ ] Two simultaneous requests cannot both pass the same availability check.
- [ ] Refund and coupon operations are idempotent.
- [ ] Transactions are kept short — no HTTP calls, emails, or queue dispatches inside a database lock.

---

## Sensitive Data
- [ ] Passwords, tokens, OTPs, card numbers, and secrets are never logged.
- [ ] Secrets and credentials are in `.env` only — never hardcoded in source.
- [ ] Sensitive database fields are encrypted at rest.
- [ ] Public-facing IDs use UUIDs or equivalent — never expose sequential auto-increment IDs.
- [ ] API responses return only the fields the client needs — no full model objects.

---

## File Parsing
- [ ] XML parsing has external entity loading disabled (XXE prevention).
- [ ] Archive decompression has a maximum size limit (zip bomb prevention).
- [ ] File uploads are served through an authenticated controller, not via a direct public URL.

---

## Rate Limiting
- [ ] Sensitive endpoints are rate-limited.
- [ ] Rate limiting is user/account-based — not IP-only.
- [ ] Login, password reset, OTP, and financial action endpoints are rate-limited.

---

## Infrastructure & Headers
- [ ] Security headers are configured: `X-Frame-Options`, `X-Content-Type-Options`, `Strict-Transport-Security`, `Content-Security-Policy`, `Referrer-Policy`.
- [ ] `APP_DEBUG=false` in all production environments.
- [ ] Stack traces are never exposed to the client.
- [ ] CORS specifies exact allowed origins — never wildcard (`*`) in production.
- [ ] Webhook signatures are verified before processing any webhook payload.
- [ ] SPF, DKIM, and DMARC DNS records are configured on all email-sending domains.
- [ ] DNS subdomains are audited — no stale records pointing to decommissioned services.

---

## References
- Security Engineering Standard: [core/08-security-engineering-standard.md](../core/08-security-engineering-standard.md)
- Security Testing & Threat Modeling: [core/09-security-testing-and-threat-modeling.md](../core/09-security-testing-and-threat-modeling.md)
- API Engineering Standard: [core/07-api-engineering-standard.md](../core/07-api-engineering-standard.md)
