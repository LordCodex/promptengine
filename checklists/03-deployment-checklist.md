---
document_id: checklists-deployment
title: Pre-Deployment Checklist
ecosystem: cross-cutting
dependencies:
  - core-security-engineering-standard
  - core-cicd-and-deployment-standard
  - core-git-and-collaboration-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Pre-Deployment Checklist

Run this checklist before every production deployment, regardless of change size. A single missed item can cause a production incident.

---

## Environment Configuration
- [ ] `APP_DEBUG=false` in the production `.env`.
- [ ] `APP_ENV=production` in the production `.env`.
- [ ] All secrets and credentials are in `.env` — nothing hardcoded in source files.
- [ ] `.env` is not publicly accessible from the web server.
- [ ] `.env.example` is updated with all new variable keys (values left empty).
- [ ] `storage/` and `bootstrap/cache/` are writable by the web server process.
- [ ] Uploaded files are stored outside `public/`.

---

## Security Verification
- [ ] Run `composer audit` — no known vulnerabilities.
- [ ] Run `npm audit` — no known vulnerabilities.
- [ ] Security headers are configured: `X-Frame-Options`, `X-Content-Type-Options`, `Strict-Transport-Security`, `Content-Security-Policy`.
- [ ] Error pages show no stack traces and no implementation details.
- [ ] All debug output removed: `dd()`, `dump()`, `var_dump()`, `print_r()`, `console.log()`, `debugger`, `ray()`.
- [ ] All unused routes, controllers, and middleware removed or disabled.
- [ ] Test and seeder routes are disabled in production.
- [ ] Headers that reveal the stack (`X-Powered-By`, `Server`) are removed or suppressed.

---

## Database & Migrations
- [ ] A full database backup has been taken before running migrations.
- [ ] Migrations have been reviewed for backward compatibility (no destructive drops or renames without a prior deprecation step).
- [ ] `php artisan migrate --force` (or equivalent) is run only after the backup is confirmed.
- [ ] Schema changes have been tested against a production-sized dataset for performance impact.

---

## Application Caches (Laravel)
- [ ] `php artisan config:cache`
- [ ] `php artisan route:cache`
- [ ] `php artisan view:cache`
- [ ] OPcache is enabled and warmed on the production server.

---

## Dependency & Build
- [ ] `composer install --no-dev --optimize-autoloader` (production Composer install).
- [ ] `npm ci` (production NPM install from lock file).
- [ ] Frontend assets have been built with `npm run build`.
- [ ] Cache-busting (version hashes) applied to static assets where relevant.

---

## Git & Code Review
- [ ] Every change was reviewed in a pull request before merging.
- [ ] The pull request diff has been self-reviewed for secrets, debug output, and leftover test code.
- [ ] Never committed directly to `main` or `master`.
- [ ] No sensitive data (passwords, API keys, database dumps) in any committed file.
- [ ] `*.log` files, `vendor/`, `node_modules/`, and `storage/` are in `.gitignore`.

---

## Monitoring & Alerting
- [ ] Application logging is configured and writing to the correct log destination.
- [ ] Alerts are configured for repeated auth failures, error spikes, and unusual activity.
- [ ] A rollback plan is documented and tested for this deployment.
- [ ] An incident response contact list is available if a post-deploy issue is detected.

---

## Post-Deployment Verification
- [ ] Health check endpoint (or smoke test) confirms the application is responding.
- [ ] Critical user journeys (login, primary feature, checkout/payment) have been manually verified.
- [ ] No new error spikes in application logs immediately after deployment.
- [ ] Queue workers are running and processing jobs (if applicable).
- [ ] Scheduled job registrations are active (if applicable).

---

## References
- CI/CD & Deployment Standard: [core/13-cicd-and-deployment-standard.md](../core/13-cicd-and-deployment-standard.md)
- Security Engineering Standard: [core/08-security-engineering-standard.md](../core/08-security-engineering-standard.md)
- Git & Collaboration Standard: [core/12-git-and-collaboration-standard.md](../core/12-git-and-collaboration-standard.md)
- Infrastructure & DevOps Standard: [core/14-infrastructure-and-devops-standard.md](../core/14-infrastructure-and-devops-standard.md)
