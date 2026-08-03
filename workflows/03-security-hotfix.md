---
document_id: workflows-security-hotfix
title: Security Hotfix Workflow
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-03
---

# Security Hotfix Workflow

This document defines the specialized fast-track workflow for developing, testing, and deploying patches for critical security vulnerabilities.

---

## 1. Vulnerability Isolation and Triage
- Identify the exact vulnerability vector (e.g. CSRF leak, SQL injection boundary, BOPLA failure).
- Perform threat modeling as described in [Security Testing and Threat Modeling](../core/09-security-testing-and-threat-modeling.md).
- Create a reproduction test case that exploits the vulnerability, demonstrating the bug on local staging environments.

---

## 2. Secure Patch Formulation
- Apply the **Three Questions Mnemonic** (Authentication -> Authorization -> Validation) to the target endpoint (refer to [Security Engineering Standard](../core/08-security-engineering-standard.md)).
- If input processing is compromised, apply strict type, boundary, and regex filters (refer to [Secure Coding Standards](../security/01-secure-coding.md)).
- Avoid introducing any new libraries or structural refactorings during a hotfix. Keep the patch footprint minimal.

---

## 3. Verification & Deployment Gateways
- Run the reproduction test to verify the patch completely blocks the exploit.
- Run the full test suite to guarantee no regressions were introduced.
- Perform a manual security sweep using the [Security Review Checklist](../checklists/02-security-review-checklist.md).
- Deploy the patch directly to production and verify that rate limiters and audit logs capture any post-deploy scans.
