# Threat Modeling Standard

This document teaches AI coding agents and developers to **think like an attacker** before implementing software features.

---

## 1. Asset Mapping

Identify the high-value targets (assets) in the active application scope:
- **User Data**: Passwords, PII, profile records.
- **Money & Balances**: Wallet ledger balances, payment credentials.
- **Tokens & Credentials**: JWTs, cookies, session identifiers.
- **Secrets**: API keys, database connection strings, KMS encryption configurations.
- **Files**: Storage objects, user-uploaded documents.
- **Admin Actions**: Setting roles, altering configuration settings.

---

## 2. Entry Points & Attack Surface

Document and inspect every point where input enters the system:
- **HTTP Endpoints & REST/GraphQL APIs**: Parameter vectors.
- **Upload Endpoints**: File payload parses.
- **Webhooks**: Unauthenticated inbound payloads (e.g., Stripe hooks).
- **OAuth Callback endpoints**: Dynamic redirect rules.
- **CLI Commands**: Local shell injections.
- **Queues & Cron Jobs**: Asynchronous payload manipulation.
- **Admin Dashboards**: Escalation access checks.

---

## 3. Trust Boundaries

Trace data transitions from untrusted zones to trusted zones:
- Identify exact trust transitions:
  `Browser Client` → `API Gateway / Route Handler` → `Application Database` → `Third-Party Services`.
- Never trust client inputs. Always assume input arriving at a boundary is hostile.

---

## 4. STRIDE Threat Model

Evaluate every target feature against the STRIDE threat matrix:

| Category | Threat Definition | Attack Vector Example | Mitigation Control |
| :--- | :--- | :--- | :--- |
| **S**poofing | Impersonating a user or service | Forging webhook payload headers | Cryptographic signature validation |
| **T**ampering | Modifying code or data | Direct path traversal in uploads | Parameterized paths, UUID identifiers |
| **R**epudiation | Denying an action occurred | Deleting database ledger rows | Immutable logging, append-only logs |
| **I**nformation | Leaking private information | Exposing stack traces on error | Global error handlers sanitization |
| **D**enial of Service | Exhausting system resources | Sending huge payload sizes | Rate limiting, max input bounds |
| **E**levation | Gaining unauthorized privileges | Altering user role variables | Strict server-side RBAC checks |

---

## 5. Feature Abuse Scenarios (Examples)

Do not limit evaluation to "how the code should run." Audit how it can be abused:
- **Duplicate Payments / Race Conditions**: Initiating multiple parallel withdrawal requests to double-spend wallet balances.
- **Coupon / Discount Abuse**: Applying discount strings concurrently to bypass quantity limits.
- **Fake Webhooks**: Sending mock payment success webhooks to unlock paid content.
- **Infinite Retries**: Triggering queue loops that consume computing resource cycles.
- **Account Enumeration**: Querying login/signup routes to harvest user registers.

---

## 6. Financial Threats Checklist

For ledger, wallet, and checkout systems, verify:
- [ ] **Replay protection**: Timestamps or nonces prevent duplicate payload submissions.
- [ ] **Race condition safety**: Pessimistic/optimistic row locks prevent double-spending.
- [ ] **Currency integrity**: Verify payment currency match constraints on transaction logs.
- [ ] **Rounding calculations**: Ensure fractional pennies/cents are rounded deterministically without leakage.
- [ ] **Negative parameters block**: Prevent inputting negative payment amounts to trigger balance increases.
- [ ] **Overflow checks**: Use overflow-safe datatypes or validation gates.

---

## 7. Required AI Threat Review Questions

Prior to writing code for high-risk features, answer these questions:
1. *What parameter values can be spoofed or manipulated?*
2. *Can this API payload be replayed or submitted multiple times?*
3. *Are ID values guessable, sequential, or enumerable?*
4. *Can validation gates be bypassed by altering HTTP headers or request formats?*
5. *Can multiple concurrent requests create a race condition or balance split?*
6. *What happens if an external dependency fails or returns error statuses?*
7. *What happens if an attacker sends malformed or huge payloads?*
