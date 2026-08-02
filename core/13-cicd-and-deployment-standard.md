---
document_id: core-cicd-and-deployment-standard
title: CI/CD and Deployment Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-security-testing-and-threat-modeling
  - core-performance-engineering-standard
  - core-testing-engineering-standard
  - core-git-and-collaboration-standard
  - stacks-php-conventions
  - stacks-laravel-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# CI/CD and Deployment Engineering Standard

## Purpose & Inheritance
This document defines the core standards for continuous integration pipelines, container build configurations, deployment strategies, configuration management, and rollback procedures. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), and all preceding core engineering documents. It establishes strict deployment rules for human developers and AI coding agents.

---

## 1. CI/CD Philosophy

Deployment is not an operational afterthought; it is a **core phase of software engineering**. A professional software delivery pipeline must enforce:
- **Zero Manual Adjustments**: No developer or administrator may alter server states, database schemas, or configuration variables manually in production. Every change must be declared in code and executed via automation.
- **Repeatable Builds**: A build artifact generated in CI must behave identically across Staging and Production. Compile once; deploy anywhere.
- **Failing Early**: Build pipelines must run fast unit and lint validations first, stopping immediately on failures to provide fast feedback loops.
- **Continuous Post-Release Verification**: A deployment is not complete when the build logs show success. It is complete only when telemetry metrics confirm that the system is serving requests normally.

---

## 2. Software Delivery Pipeline

Our delivery pipeline consists of eight structured stages:

```text
[Code Change] ──> [Commit] ──> [Pull Request] ──> [Automated Checks (CI)]
                                                         │
[Monitor] <── [Deploy] <── [Build Image] <── [Review & Merge] <┘
```

1. **Code Change**: Developer edits code locally in a short-lived branch.
2. **Commit**: Commits are made using Conventional Commits. Pre-commit hooks verify that no secrets are committed.
3. **Pull Request**: Opening a PR triggers the CI pipeline.
4. **Automated Checks (CI)**: Pipeline runs lint checks, static analysis (SAST), and unit/integration tests.
5. **Review & Merge**: Team reviews the PR, checks security/architecture, and merges it to `main`.
6. **Build Image (Artifact Creation)**: CI compiles static assets, installs locked dependencies, builds a Docker image, and pushes it to a secure registry.
7. **Deploy**: The orchestrator deploys the image using a defined deployment strategy (e.g., Rolling or Blue-Green).
8. **Monitor**: Telemetry systems verify application health post-deploy.

---

## 3. Environment Strategy & Configuration Management

We enforce strict isolation between environments to prevent staging tests from polluting production resources.

### Environment Schema

| Environment | Purpose | Configuration Source | Data Handling | Access Bounds |
| :--- | :--- | :--- | :--- | :--- |
| **Development** | Local coding, manual debugging | Local `.env` file | Dummy data generated via factories | Local developer |
| **Testing** | CI pipeline execution | GitHub Actions Secrets | In-memory DB, reset between tests | Automated runner |
| **Staging** | Production-equivalent testing | Cloud Secret Manager / KMS | Sanitized, anonymized DB snapshot | Restricted read |
| **Production** | Live customer traffic | Cloud Secret Manager / KMS | Real user data, persistent backups | Strictly blocked |

### Configuration Management
- **Environment Variables**: Application configuration must be injected at runtime using environment variables. Never hardcode settings inside application code.
- **Feature Flags**: Wrap new or high-risk features in feature flags. This decouples code deployment from feature release, allowing features to be toggled off instantly without redeploying code.
- **Secrets Encryption**: Secrets must be stored in secure Key Management Services (KMS) or Secret Managers (e.g., AWS Secrets Manager, HashiCorp Vault). Never commit plain secrets to version control.

---

## 4. GitHub Actions Workflow Standards

We use GitHub Actions as our primary CI/CD runner.

### Workflow Principles
- **Fast Feedback Triggers**: Run pipelines on `pull_request` and `push` to target release branches (`main`, `master`).
- **Caching Dependencies**: Cache package folders (`vendor/`, `node_modules/`, `~/.pub-cache`) using runner cache keys to speed up pipeline execution.
- **Matrix Builds**: Run integration tests across supported database engines or PHP versions concurrently using build matrices:
  ```yaml
  strategy:
    matrix:
      php-version: [8.3, 8.4]
      db-type: [mysql, postgresql]
  ```
- **Job Concurrency Control**: Cancel in-progress pipelines when a developer pushes new commits to the same PR branch to save runner resources.

---

## 5. Docker & Container Best Practices

Containers ensure execution consistency across local development, staging, and production.

### Multi-Stage Dockerfile Blueprint (Laravel Production)
```dockerfile
# Stage 1: Build PHP & Node dependencies
FROM php:8.3-fpm-alpine AS builder

WORKDIR /var/www/html

# Install dependencies and extensions
RUN apk add --no-cache git unzip libpng-dev libzip-dev zip
RUN docker-php-ext-install pdo_mysql gd zip

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer
COPY composer.json composer.lock ./
RUN composer install --no-dev --no-scripts --no-autoloader --prefer-dist

COPY . .
RUN composer dump-autoload --optimize --classmap-authoritative

# Stage 2: Final Production Runtime Image
FROM php:8.3-fpm-alpine

WORKDIR /var/www/html

RUN apk add --no-cache libpng libzip
RUN docker-php-ext-install pdo_mysql gd zip

COPY --from=builder /var/www/html /var/www/html

# Enforce Non-Root User Execution
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -D appuser
USER appuser

EXPOSE 9000
CMD ["php-fpm"]
```

### Docker Directives
- **Enforce Non-Root Users**: Never run container entrypoints as `root` in production. Always declare a low-privilege `USER` in the Dockerfile.
- **Pin Base Images**: Enforce explicit version tags on base images (e.g., `php:8.3-fpm-alpine`). Never build from the `latest` tag.
- **Lock Dependencies**: Copy `composer.lock` or `package-lock.json` and install dependencies using strict lock flags (`composer install --no-dev`, `npm ci`).

---

## 6. Deployment Strategies

### 1. Rolling Deployment
- **Method**: The orchestrator updates running container instances incrementally (one by one or in batches) until all nodes run the new version.
- **Benefits**: Requires zero additional hardware overhead; maintains constant throughput.
- **Risks**: Triggers a "mix-version state" where user requests hit both old and new containers. Database migrations must be backward-compatible.

### 2. Blue-Green Deployment
- **Method**: Provision an entirely separate server cluster running the new code version ("Green"). Run health checks, then switch the load balancer traffic pointer from the old cluster ("Blue") to the new cluster.
- **Benefits**: Zero downtime; provides instant rollback by switching the traffic pointer back to Blue.
- **Risks**: Requires double the server resource footprint during deployments; database changes must support both versions.

### 3. Canary Deployment
- **Method**: Deploy the new container version to a small subset of servers (e.g., 5% of traffic). Monitor error rates and user impact. If successful, roll out the change to the remaining servers.
- **Benefits**: Minimizes blast radius on critical failures.
- **Risks**: Complex routing requirements.

---

## 7. Database Migrations in CI/CD

Database migrations are the highest-risk phase of deployment. They must be managed with strict sequencing.

```text
Deploy Stage 1: Run backward-compatible schema changes (Expand)
       ↓
Deploy Stage 2: Deploy new application code container cluster
       ↓
Deploy Stage 3: Decommission deprecated columns (Contract)
```

### Migration Safety Rules
- **Automatic Execution**: Run database migrations automatically in staging and production immediately before the new container cluster boots.
- **No Downgrades in Production**: Never execute rollback migrations (`down()` scripts) on production databases containing live data. Doing so can cause catastrophic data loss.
- **Lock Mitigation**: Specify statement timeout limits on migration runs to prevent them from blocking database connection queues during lock contentions.

---

## 8. Zero-Downtime Releases

Zero-downtime releases require coordinating proxy load balancers and application container states.

### Implementation Checklist
1. **Health Check Endpoints**: Implement a `/health` endpoint that returns a `200 OK` status code only when the database and cache connection checks succeed.
2. **Graceful Shutdown**: Configure application servers to handle `SIGTERM` signals correctly. Let them finish active HTTP requests before exiting.
3. **Load Balancer Draining**: The load balancer must drain connections from old containers, stopping traffic before they are terminated.

---

## 9. Rollback Strategy

Every deployment must have a pre-configured rollback plan.

### Rollback Protocols
- **Application Rollback**: Revert to the previously tagged Docker image on the container registry. This takes seconds and does not require building new code.
- **Configuration Rollback**: Revert environment variable changes in the Secrets Manager and trigger a rollout.
- **Database Rollback Strategy**:
  - Since running down migrations in production is prohibited, rollbacks are handled by deploying a hotfix version of the code that supports the modified database schema.
  - Apply the "Expand and Contract" pattern to ensure that the database schema is compatible with both the old and new application versions.

---

## 10. Post-Deployment Monitoring

A deployment is incomplete until telemetry checks confirm that system behavior remains normal.

### Verification Controls
- **Error Rates Check**: Monitor error logs for spikes in `5xx` status codes or exception traces.
- **Latency Monitoring**: Verify that $P_{95}$ and $P_{99}$ response times do not exceed SLA targets.
- **CPU & Memory Profiling**: Monitor server resource consumption to identify memory leaks or CPU bottlenecks.

---

## 11. Security in CI/CD

Protecting the build pipeline prevents supply chain attacks and credential exposure.

### Security Controls
- **OIDC (OpenID Connect)**: Avoid storing permanent AWS or GCP access keys in GitHub Secrets. Use OIDC federation to allow runners to request temporary AWS/GCP credentials dynamically.
- **SCA Scans**: Configure automated Software Composition Analysis (SCA) to scan dependencies and Docker base images for open vulnerabilities.
- **Immutable Artifacts**: Tag Docker images with unique commit hashes (`SHA`) and enforce read-only permissions on production image tags to prevent them from being overwritten.

---

## 12. Legacy Deployments

When working with legacy environments that rely on manual deployments (e.g., FTP uploads or manual `git pull` commands on a VPS):
1. **Document the Manual Steps**: Record the current manual deployment process in a Markdown file.
2. **Automate Incrementally**: Replace manual steps step-by-step. Start by automating code checkouts, then database migrations, and finally asset compilation.
3. **Containerize**: Package the legacy application inside a Docker container to decouple it from host dependency configurations, allowing safe server upgrades.

---

## 13. Decision Matrices

Use these matrices to identify the correct deployment engineering decision based on project context.

### Matrix 1: Docker vs. Native Deployment
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Microservice architectures, multi-stack codebases, Kubernetes environments | **Docker Container** | Guarantees environmental consistency; packages dependencies. |
| Monolithic applications on simple, managed VPS platforms | **Native Deployment** | Low complexity; avoids container registry and orchestration overhead. |

### Matrix 2: Manual Deployment vs. Automated Deployment
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| All production-ready systems, staging run tracks, SaaS systems | **Automated** | Eliminates human execution errors and speeds up verification. |
| Local developer test environments, initial infrastructure bootstrap | **Manual** | Low initial setup overhead when automation pipelines do not exist. |

### Matrix 3: Continuous Deployment vs. Approval Deployment
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| High-trust SaaS features, low-impact web assets, strong automated test suites | **Continuous (Fully Auto)** | Shortens delivery cycles; merges are deployed instantly. |
| Financial systems, core authentication updates, compliance systems | **Approval Gate** | Adds a final human validation step to check risks and verify readiness. |

### Matrix 4: Rolling vs. Blue-Green Deployment
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Systems with strict budget limits, low-impact resources | **Rolling Deployment** | Fits within current hardware resource limits; updates nodes incrementally. |
| High-value web portals, zero-tolerance downtime systems | **Blue-Green** | Allows instant rollback by switching load balancer traffic pointers. |

### Matrix 5: Monolith Deployment vs. Multiple Services
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Simple application layouts, small engineering teams | **Monolith** | Simplifies deployment pipelines, database management, and routing. |
| High-scale distributed domains, multiple independent development teams | **Multiple Services** | Decouples deploy dependencies; enables independent scaling. |

### Matrix 6: Self-Hosted VPS vs. Managed Platform (PaaS / Serverless)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Cost-sensitive startup setups, low traffic systems, legacy systems | **Self-Hosted VPS** | Highly customizable; low upfront platform costs. |
| Fast growth SaaS systems, high-traffic APIs, dynamic scale demands | **Managed Platform** | Outsources server maintenance, scaling, and backups to the cloud provider. |

---

## 14. AI Deployment Rules

AI agents modifying pipelines or Docker files in this repository must follow these rules:

1. **Verify Environment Configurations**: Do not change environment variables in production setups without verifying staging and development setups first.
2. **Never Remove Pipeline Safety Checks**: Do not disable automated tests, static analysis steps, or security audits to bypass pipeline build failures.
3. **No Raw Secrets in Workflows**: Ensure all credential settings use GitHub Secrets variables or KMS integration. Never write hardcoded keys inside pipeline files.
4. **Use Multi-Stage Dockerfiles**: Enforce multi-stage Docker builds to keep image sizes small and minimize vulnerability scopes.
5. **No Blind Image Swaps**: Do not upgrade base container images without testing dependencies locally first.

---

## 15. Deployment Review Checklist

Use this checklist during code review to evaluate deployment pipeline safety.

### Before Deployment
- [ ] Have all automated checks (linting, static analysis, unit/integration tests) passed?
- [ ] Are database migrations backward-compatible (using the "Expand and Contract" pattern)?
- [ ] Is the rollback plan documented?
- [ ] Have secrets been configured securely in the cloud Secret Manager (no plain credentials in code)?

### During Deployment
- [ ] Is post-deployment logging active?
- [ ] Are slow SQL query logs and error telemetry being monitored?

### After Deployment
- [ ] Has the `/health` check returned `200 OK`?
- [ ] Have response latency and error rates been verified against SLA targets?

---

## References
- Database Migration Safety: [06-database-engineering-standard.md](06-database-engineering-standard.md)
- Testing & CI Integrations: [11-testing-engineering-standard.md](11-testing-engineering-standard.md)
- Git Branching & Merges: [12-git-and-collaboration-standard.md](12-git-and-collaboration-standard.md)
