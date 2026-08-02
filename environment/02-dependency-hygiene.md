---
document_id: env-dependency-hygiene
title: Dependency Hygiene and Lockfile Management
ecosystem: cross-cutting
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Dependency Hygiene and Lockfile Management

## Purpose
This document establishes rules for importing external libraries, pinning dependencies, and managing package manifest lockfiles to maintain build stability and secure code pipelines.

## Scope
Applies to PHP Composer (`composer.json`), Node NPM (`package.json`), and Dart Pub (`pubspec.yaml`) configuration systems.

---

## Dependency Rules

### 1. Lockfile Integrity
- **Always commit lockfiles**: You must commit `composer.lock`, `package-lock.json`, and `pubspec.lock` to the git repository.
- **Do not modify lockfiles manually**: Only generate or update lockfiles through the package manager commands (`composer update`, `npm install`, `flutter pub upgrade`).
- **Single Dependency Updates**: When updating a library, target it specifically (e.g. `composer update vendor/package` rather than running a global `composer update`). This limits changes and reduces deployment risk.

### 2. Version Locking Constraints
To prevent unexpected breakages, define version ranges in manifest files using the following standards:

| Dependency Type | Version Constraint Pattern | Example |
| :--- | :--- | :--- |
| Core Frameworks | Strict patch or minor locking | Laravel: `"laravel/framework": "^10.0"` |
| Database Drivers | Exact version pinning | Redis: `"predis/predis": "2.2.0"` |
| Client UI Components | Caret locking (non-breaking upgrades) | Vue: `"vue": "^3.3.0"` |
| Native Device Bridges | Pinned minor updates | Flutter plugins: `path_provider: ^2.1.0` |

---

## Dependency Security and Auditing

- **Automated Security Audits**: Run auditing tools regularly to inspect dependencies for known vulnerabilities:
  - PHP: `composer audit`
  - JS/TS: `npm audit` or `pnpm audit`
  - Dart: `dart pub deps` (inspect tree layout)
- **Zero Vulnerability Standard**: CI/CD pipelines must fail build runs if vulnerabilities categorized as `high` or `critical` are found in committed configurations.

---

## Evaluation Checklist for New Libraries
Before importing a new third-party dependency, developers and agents must answer the following questions:
1. *Can this functionality be written natively in less than 50 lines of clean code?* If yes, write it custom.
2. *Is the repository actively maintained?* Check if the last commit was within the last 6 months.
3. *What is the license?* Verify compatibility (e.g., MIT, BSD, Apache. Avoid AGPL in commercial contexts).
4. *How many nested transitive dependencies does this import pull in?* Check with `npm ls` or equivalent command before committing.

---

## Common Mistakes & Anti-Patterns
- **Wildcard Ranges**: Using wildcard version targets (`"package": "*"`) which cause production builds to fetch untested releases.
- **Unused Packages**: Retaining unused packages in the `require`/`dependencies` block instead of moving them to `require-dev`/`devDependencies` or purging them completely.
- **Direct Code Modification**: Modifying dependency source code directly in local `vendor/` or `node_modules/` folders instead of extending classes or opening pull requests to upstream repositories.

---

## References
- Environment setup: [01-local-dev-standards.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/01-local-dev-standards.md)
- Automated pipelines: [03-ci-cd-pipelines.md](file:///Users/kodexkode/Documents/workspace/promptengine/environment/03-ci-cd-pipelines.md)
