---
document_id: security-secrets-management
title: Secrets and Configuration Management
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Secrets and Configuration Management

## Purpose
This document outlines rules for storing private keys, database passwords, and API credentials securely in code deployment pipelines.

## Scope
Applies to environment config files, deployment parameters, and runtime vaults.

---

## Directives

### 1. Separation of Configuration and Secrets
- **Configuration** defines non-sensitive parameters (e.g. `DB_PORT`, `APP_TIMEZONE`, `LOG_CHANNEL`). It is safe to commit default configuration values to standard version control.
- **Secrets** define access variables (e.g. database password, Stripe private token, SSH key). These must never be committed to code repositories.

### 2. Local Environment Isolation
- Always define secret values in a local, ignored `.env` file (refer to [environment/01-local-dev-standards.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/01-local-dev-standards.md)).
- Do not store fallback api credentials directly in code constants. Use env lookup with fallback triggers that throw runtime warnings:
  ```php
  // Safe env loading
  $stripeSecret = env('STRIPE_SECRET') ?? throw new ConfigurationException("Stripe token missing");
  ```

### 3. Production Secret Injection
- In production platforms (AWS, GCP, Heroku, Forge), inject secrets dynamically through:
  - Host environment variable containers.
  - KMS vault key store systems (e.g., HashiCorp Vault, AWS Secrets Manager).
- Never share production database keys with team members over messaging systems or email channels.

---

## Common Mistakes & Anti-Patterns
- **Committed SSH Keys**: Storing deployment public/private keys under `/storage/` or `/keys/` subfolders within the git tree.
- **Leaked API Keys in Pull Requests**: Temporarily hardcoding an active API token during local debugging and committing it to a remote pull request.
- **Shared Production Keys**: Using the production database credentials for local testing or staging environment containers.

---

## References
- Local environment parameters: [environment/01-local-dev-standards.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/01-local-dev-standards.md)
- git exclusion list: [.gitignore](file:///Users/kodexkode/Documents/workspace/promptengine/.gitignore)
