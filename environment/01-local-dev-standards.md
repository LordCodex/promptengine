---
document_id: env-local-dev-standards
title: Local Development Environment Standards
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Local Development Environment Standards

## Purpose
This document establishes conventions for running local developer environments, ensuring identical runtime versions across developer machines and preventing local credential exposure.

## Scope
Applies to local database configuration, host port routing, compiler settings, and environment files (`.env`).

---

## Environment Isolation

All developers and AI agents must run projects using one of the two standard isolation mechanisms:

### Option A: Containerized (Docker Compose)
- **Standard**: Every project must contain a root-level `docker-compose.yml` specifying exact runtime images and database tags.
- **Rules**:
  - Never use the `latest` tag for databases. Specify exact minor tags (e.g., `postgres:15.4-alpine`, `redis:7.0-alpine`).
  - Pin the PHP or Node compiler container to match production (e.g. `php:8.2-fpm-alpine`).

### Option B: Native macOS (ServBay / Homebrew)
- For projects requiring high-performance disk access (e.g. large legacy monoliths), ServBay or pinned Homebrew services are authorized.
- **Rules**:
  - Keep configuration files (like `php.ini` or virtual host declarations) version-controlled under a `.dev-config/` project sub-directory.

---

## Port Assignment Matrix
To prevent local system port clashes across multiple running client projects:

| Port | Standard Usage | Clash Resolution Strategy |
| :--- | :--- | :--- |
| `80` / `443` | Host Reverse Proxy (Nginx/ServBay) | Route subdomains dynamically (`http://project-a.test`) |
| `3306` | MySQL Core Database | Map container port to `33061` if running multiple MySQLs |
| `5432` | PostgreSQL Core Database | Map container port to `54321` |
| `6379` | Redis Cache Instance | Map container port to `63791` |
| `8000` | Laravel Dev Server (`artisan serve`) | Bind to `8001`, `8002` sequentially |
| `3000` / `5173` | Vue / Vite / Nuxt Dev Server | Enable automated port increments in configuration |

---

## Environment Variables (.env) Security
- **Rule 1**: Never commit `.env` files to git repositories. Always ignore them via `.gitignore` (refer to [.gitignore](file:///Users/kodexkode/Documents/workspace/promptengine/.gitignore)).
- **Rule 2**: Maintain a `.env.example` file in the repository root. This file must contain all valid keys with safe, local-only defaults or dummy placeholder strings (`STRIPE_KEY=your_stripe_key_here`).
- **Rule 3**: Never put production secrets or API keys in `.env.example`.

---

## Common Mistakes & Anti-Patterns
- **Global Tooling Dependencies**: Running terminal scripts using a globally installed PHP/Node version that differs from the project-locked composer/package version.
- **State Leakage**: Failing to mount database storage to a named Docker volume, causing local test data to vanish every time containers stop.
- **Hardcoded Localhost**: hardcoding `127.0.0.1` inside server code instead of referencing environmental dynamic host properties (`DB_HOST`).

---

## References
- Auditing dependencies: [02-dependency-hygiene.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/02-dependency-hygiene.md)
- Secret storage rules: [security/03-secrets-management.md](file:///Users/kodexkode/Documents/workspace/promptengine/security/README.md)
