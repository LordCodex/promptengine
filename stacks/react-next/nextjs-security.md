---
document_id: stacks-nextjs-security
title: Next.js Security Hardening Standard
ecosystem: react-next
dependencies:
  - core-universal-coding-standards
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Next.js Security Hardening Standard

## Inheritance & Constraints
This document inherits from the [Universal Coding Standards](../../core/05-universal-coding-standards.md). It outlines security guidelines specific to Next.js.

---

## 1. Cross-Site Scripting (XSS)

- **Safeguard HTML**: Avoid `dangerouslySetInnerHTML`. If raw HTML rendering is mandatory, validate and sanitize inputs server-side before execution using tools like `dompurify`.
- **Client Data Hydration**: Never pass database models directly to Client Components. Extract and expose only the specific fields needed to prevent accidental PII leakage.

---

## 2. Server Action Hardening

- **Authorization Boundaries**: Re-authenticate and re-authorize users inside every Server Action block. Do not assume client-side middleware routing checks are sufficient.
- **CSRF Mitigations**: Next.js automatically verifies CSRF headers for Server Actions. Do not bypass or manipulate custom action post headers.

---

## 3. Environment Variables Configuration

- **Expose Rules**: Only prefix variables with `NEXT_PUBLIC_` if they are intended to be visible to the public user browser.
- **Keep Secrets Server-Only**: Keep private API keys, credentials, and passwords un-prefixed. Ensure they remain inaccessible to the Client bundle.

---

## Review Checklist

Verify security boundaries against this checklist:
- [ ] Inherits correctly from Universal Coding Standards.
- [ ] No database models are leaked to Client Components.
- [ ] Server Actions validate authentication status explicitly.
- [ ] Private environment variables exclude public prefixes.
