---
document_id: core-deployment-cicd-and-production-readiness
title: Deployment, CI/CD, Infrastructure, and Production Readiness Standard
ecosystem: cross-cutting
dependencies:
  - core-frontend-security
  - core-observability-and-operational-excellence
  - core-git-and-collaboration-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Deployment, CI/CD, Infrastructure, and Production Readiness Standard

## Purpose & Inheritance
This document defines core standards for automated deployment workflows, Continuous Integration (CI) and Continuous Delivery (CD) quality gates, Infrastructure as Code (IaC) hygiene, database migration safety, and post-deployment validation. It inherits from the [Frontend Security Standard](24-frontend-security-and-privacy-hardening-standard.md), the [Observability Standard](29-observability-and-operational-excellence-standard.md), and the [Git & Collaboration Standard](12-git-and-collaboration-standard.md).

---

## 1. Deployment Philosophy

- **Small & Frequent Releases**: Avoid high-risk "big bang" deployments. Release features in small, decoupled, logical steps.
- **Automated Validation**: Rely on automated testing and pipeline validation rather than manual smoke testing.
- **Zero Production Manual Adjustments**: Banish manual hotfixes or configuration edits on live production servers. Everything must flow through the automated pipeline.

---

## 2. CI/CD Pipeline Gates

Every commit pipeline must execute validation checks and fail fast on quality gate errors:
- **Linting & Formatting**: Enforce strict syntax formatting and type checks.
- **Automated Testing**: Run unit and integration tests.
- **Configuration Validation**: Verify environment variable structures against schemas before build execution.
- **Security Scanners**: Scan codebases for hardcoded credentials (e.g. TruffleHog) and vulnerability exposures.
- **Supply Chain Security**:
  - Check dependencies against security advisories (e.g. `npm audit`, `composer audit`).
  - Prefer official package registries; verify dependency integrity checksums (`package-lock.json`, `composer.lock`).
  - Maintain a Software Bill of Materials (SBOM) where appropriate.
  - Sign build artifacts and restrict access permissions to CI runners.

---

## 3. Infrastructure as Code (IaC)

- **Version-Controlled Infrastructure**: All environments must be provisioned and updated using version-controlled configurations (Terraform, Pulumi, Kubernetes manifests, or Docker Compose).
- **No Manual Adjustments**: Banish manual configuration changes in cloud provider web interfaces.
- **Environment Parity**: Dev, staging, and production environments should share identical service layouts, databases, and dependencies.

---

## 4. Production Secret Configuration

- **Zero Hardcoded Secrets**: Secrets, credentials, API keys, and connection strings must reside outside source code.
- **Secret Managers**: Retrieve sensitive variables using secure secrets management engines (e.g. HashiCorp Vault, AWS Secrets Manager, GCP Secret Manager).
- **Least Privilege Access**: monitoring systems, dashboards, CI/CD runners, and logging platforms must receive only the minimum permissions they require to operate.

---

## 5. Containerization Hygiene

Container configurations must be secure and lightweight:
- **Minimal Base Images**: Build from alpine or distroless images. Avoid shipping development tooling, package managers, or compilers into production images.
- **Non-Root Execution**: Configure Dockerfiles to run processes using a non-root system user.
- **Image Pinning**: Always pin base images and dependencies using explicit tags or digests (never use `latest`).
- **Minimize Image Layers**: Combine RUN instructions and clean up build cache flags in the same command layer.
- **Graceful Shutdown**: Support standard termination signals (`SIGTERM`) allowing application containers to complete active DB transactions and drain network requests before closing.

---

## 6. Database Deployment Safety

- **Backward-Compatible Schemas**: Database schema migrations must be backward-compatible with the immediately preceding application code version to support zero-downtime rolling deployments.
- **Safe Operations**:
  - Do not drop columns or rename tables during active deployments.
  - Separate database changes into two stages: first add the new columns/tables, then deploy code that reads/writes to them, and finally clean up legacy fields.
- **Reversible Migrations**: Write down-migration scripts for all schema changes, verifying that data can be restored safely if rollbacks trigger.

---

## 7. Zero-Downtime Strategies

Production rollouts should maintain application availability:
- **Rolling Deployments**: Gradually replace container instances, verifying target health check endpoints before routing traffic.
- **Blue/Green & Canary Releases**: Deploy releases to an isolated environment or route a small percentage of traffic (e.g. 5%) to a canary group, monitoring error metrics before routing full traffic.
- **Graceful Restarts**: Ensure servers and web servers drain active socket requests on restart.

---

## 8. Rollback Strategies

Prepare explicit rollback triggers and runbooks for every release:
- **Code Rollbacks**: Fast rollback of Docker tags or package versions.
- **Configuration Rollbacks**: Instantly revert environment variables without rebuilding the code.
- **Database Rollbacks**: Execute verified database down-migrations.
- **Feature Flag Rollbacks**: Disable unstable features via admin dashboard flags without touching server code.

---

## 9. Backups & Disaster Recovery (DR)

- **Verified Backups**: Maintain automated daily backups of databases and storage directories.
- **Restoration Verification**: Frequently test and verify data restoration from backups. A backup is useless unless it has been proven to restore correctly.
- **Disaster Recovery Runbooks**: Maintain actionable documents outlining procedures for database corruption, hosting region failures, primary network failures, and TLS certificate failures.

---

## 10. TLS, HTTPS, and Network Transport

- **HTTPS Only**: Enforce HTTPS for all production web services.
- **TLS Requirements**: Disable legacy TLS versions; enforce TLS 1.3 or modern TLS 1.2 profiles.
- **Secure Cookies**: Apply `Secure`, `HttpOnly`, and `SameSite=Strict` flags to session cookies.
- **HSTS Config**: Enforce HTTP Strict Transport Security (HSTS) headers.

---

## 11. Caching & Asset Management

- **Caching Control**: Clear and warm configuration caches, route caches, and view compilers during deployment to prevent stale configurations.
- **Asset Cache Busting**: Frontend script and style assets must be minified, compressed (Gzip/Brotli), and appended with unique version hashes (cache-busting).
- **Source Maps**: Do not deploy source maps to public servers unless explicitly required for client-side crash tracking.

---

## 12. Post-Deployment Verification

Immediately after deployment completes, verify:
1. Health and readiness endpoints (`/health/ready`) return `200 OK`.
2. Error rates do not spike in centralized trackers (Sentry).
3. Background queue workers and scheduled cron tasks execute successfully.
4. Critical user journeys function (user logins, page hydration).
5. Payment gateway operations are online.
6. Third-party API integrations resolve without timeouts.

---

## 13. Resource Allocation & Cost Awareness

SRE designs must optimize resource boundaries:
- **Scale Parameters**: Align compute instances, database connections, cache capacities, and memory configurations with measured production requirements.
- **Cost Metrics**: Review hosting costs before scaling capacity. Avoid keeping unused servers, disks, or database read replicas.

---

## 14. Required AI Deployment Review Questions

Before verifying that a feature is production-ready, ask:
1. *Can this code deploy safely without causing schema locks or downtime?*
2. *Is there an explicit, tested rollback path for code, database, and configurations?*
3. *Can operators monitor the health and throughput of this change in real-time?*
4. *Are errors, logs, and latency traces structured to allow immediate diagnosis?*
5. *Does the config change maintain credential security and protect secrets?*
6. *Have dependency upgrades been verified against licensing and security advisories?*

---

## Review Checklist

Verify the release against this production-readiness checklist:
- [ ] **CI pipeline passes**: Linting, formatting, compilation, and automated test steps complete successfully.
- [ ] **Tests pass**: No skipped or broken unit and integration tests are bypassed.
- [ ] **Security scans pass**: Code vulnerability and credential leaks scanners show zero warnings.
- [ ] **Secrets managed securely**: Secrets reside in secure variable injection systems, never committed to code repositories.
- [ ] **Configuration validated**: Configuration variables are checked for correctness against validation schemas.
- [ ] **Rollback documented**: Step-by-step code and database rollback paths are documented and verified.
- [ ] **Monitoring configured**: Real-time logging, P95/P99 latency metrics, error tracking, and alerting rules are set.
- [ ] **Health checks verified**: Liveness (`/health/live`) and readiness (`/health/ready`) APIs are active and functional.
- [ ] **Backups verified**: Automated backup routines are active and restoration commands are tested.
- [ ] **Dependency review completed**: Package dependencies are audited for licenses, security disclosures, and runtime compatibility.
- [ ] **Infrastructure reviewed**: All related IaC files are committed, reviewed, and match target staging environments.
- [ ] **Release notes prepared**: Clean changelogs detail new changes, deprecations, and configuration requirements.
- [ ] **Deployment automated where practical**: The deployment executes via automated scripts or webhook hooks.

---

## References
- Git and Collaboration: [core/12-git-and-collaboration-standard.md](core/12-git-and-collaboration-standard.md)
- Observability and Logging: [core/29-observability-and-operational-excellence-standard.md](core/29-observability-and-operational-excellence-standard.md)
- Infrastructure Local Setup: [environment/01-local-dev-standards.md](../environment/01-local-dev-standards.md)
- Docker Best Practices: [https://docs.docker.com/develop/develop-images/dockerfile_best-practices/](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
