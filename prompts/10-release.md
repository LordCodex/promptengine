# 10. Release Readiness Prompt

---

## Purpose
Instructs the AI to perform a pre-release audit to ensure container configurations, secrets injection, deployment tasks, and change logs are production-ready.

## When to use
Use before merging a release branch to `main` or deploying to staging/production.

## Example
Pre-release checking a new Laravel microservice build.

---

## Copy-and-Paste Prompt

```markdown
We are preparing a new release for "{PROJECT_NAME}".

Please perform a production readiness audit:
1. **Docker Configurations**: Review Dockerfiles and compose files. Ensure pinned base images are used, multi-stage builds are optimized, and containers run under non-root users.
2. **Secrets Hygiene**: Scan environment variable templates (`.env.example` or compose files). Confirm no live credentials, API keys, or certificates are committed in source files.
3. **Database Migration Checks**: Check recent migration scripts for long-running locks or potential table locking behaviors during deployment.
4. **Onboarding & Setup scripts**: Verify that the local setup/diagnostics scripts execute cleanly.
5. **Changelog & Documentation**: Reconcile completed features against `docs/Progress.md` and `docs/Roadmap.md`. Update any deployment-specific parameters in `docs/Deployment.md`.
6. **Release Notes Draft**: Generate a concise release summary showing features added, bugs resolved, and operational requirements.
```

---

## Expected AI Behaviour
1. The AI audits configuration files, Dockerfiles, and compose files.
2. It outputs a checklist of security and build quality parameters.
3. It highlights deployment risks (e.g. un-cached dependencies, raw secret values).
4. It drafts a copy-and-paste ready Release Notes summary based on `docs/Progress.md`.

## Common Mistakes
- **Skipping lock checks**: The AI failing to inspect migration locks on large tables. Ensure migrations are explicitly listed.
- **Leaked environment parameters**: Missing variables in `.env.example`.
