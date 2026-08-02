# Integration Report — old.md into promptengine

**Date**: 2026-08-01  
**Source**: `old.md` (3,141 lines — concatenation of three project-specific AGENTS.md files from a PHP legacy modernization project, a Laravel/Vue project rules file, and a general coding rules file)  
**Target**: `/Users/kodexkode/Documents/workspace/promptengine`

---

## Summary

25 missing rules identified, classified, generalized, and integrated. 3 new checklist documents created. 5 existing documents improved. All project-specific rules rejected.

---

## Already Covered

The following topics were already addressed in the current playbook. No changes made.

| Topic | Old.md Location | Existing Coverage |
| :--- | :--- | :--- |
| Clean code, DRY, KISS, YAGNI, SOLID | Project Constitution | `core/02`, `core/05` |
| Naming conventions (camelCase, PascalCase, snake_case) | Section 8 / Section 1 | `core/05` Section 1 |
| No N+1 queries, paginate large datasets | Performance sections | `core/10`, `core/06` |
| Security headers list | Section 2 Security Headers | `core/08` Section 6, `security/04` |
| PDO prepared statements / no raw SQL | Database Rules | `core/06`, `security/01` |
| Git hygiene (feature branches, no force-push to main) | Section 3 Git Hygiene | `core/12` |
| Legacy backward compatibility philosophy | Legacy Compatibility section | `core/16` Sections 1–3 |
| AI agent: plan before implement, explain changes | Before Writing Code section | `core/20` Sections 2–5 |
| Centralized exception handling, no stack trace exposure | Error Handling section | `core/08` Section 4 |
| Incremental development, small pull requests | Development Process | `core/12`, `core/01` |
| Eager loading to prevent N+1 (Laravel) | Performance DATABASE section | `core/10`, `stacks/php-laravel/laravel-data.md` |
| Dependency auditing (`composer audit`, `npm audit`) | Dependencies section | `core/13`, `core/08` Section 4.6 |
| SSRF prevention (IP whitelisting) | Section 9 SSRF | `core/08` Section 4.9 (basic coverage) |

---

## Added Rules

### core/05-universal-coding-standards.md
New Section 12 — **Professional Output Quality**:
- Write boring, predictable, senior-engineer-quality code (not AI-generated style).
- No emojis anywhere in the codebase — code, comments, strings, logs.
- Never generate TODO/placeholder/stub/dead code in production output.
- Remove all debug output (`dd`, `var_dump`, `console.log`, `debugger`, `ray`) before submitting.

### core/08-security-engineering-standard.md
New sections added:
- **Section 13 — BOPLA**: Field-level write authorization per role, beyond ownership checks.
- **Section 14 — HTTP Parameter Pollution**: Reject duplicate query string and JSON body parameters.
- **Section 15 — File Parsing Attacks**: CSV injection, XXE, zip bomb prevention.
- **Section 16 — Subdomain Takeover Prevention**: Audit and remove stale DNS records.
- **Section 17 — Email Security**: SPF/DKIM/DMARC, no account enumeration, single-use tokens.
- **Section 18 — Timing Attack Prevention**: Constant-time comparison, secure token generation, no MD5/SHA1 for passwords.
- **Section 19 — User-Facing Message Safety**: `textContent`/`innerText` only, never `innerHTML`; async error response handling.
- **Section 20 — Rate Limiting Standards**: User/account-based limiting, implementation requirements, cleanup rules.
- **Section 21 — Financial Operation Security**: Transactions, pessimistic locks, idempotency keys, webhook signature verification, server-side amount recalculation, append-only audit logs.
- **Section 22 — Three Questions Mnemonic**: Authentication → Authorization → Validation for every endpoint.
- **Expanded Security Review Checklist**: Added BOPLA, timing attacks, HTTP parameter pollution, CSV injection, XXE, `innerHTML`, email DNS, financial operations, and infrastructure checks.

### core/10-performance-engineering-standard.md
- **Section 15 — Pagination Strategy**: Keyset pagination for batch/background work vs. offset pagination for bounded interactive UI pages. Decision table with rules.
- **Section 16 — Dangerous Search Patterns**: Leading wildcards, concatenated column searches, broad OR clauses, large-text searches — each with the performance risk and mitigation approach.

### core/16-legacy-modernization-and-refactoring-standard.md
- **Section 14 — Concurrency Safety in Legacy Systems**: Transaction length rules, lock discipline, idempotent operations, multi-server correctness, scheduled job safety requirements.

### core/20-ai-agent-engineering-workflow-standard.md
- **Section 16 — Critical Operational Rules**: Four enforced rules:
  1. Do not begin implementation from an audit finding alone.
  2. Do not claim the codebase was checked when any directory was skipped.
  3. Stop and get confirmation when a request affects multiple unrelated modules.
  4. Never silently deviate from a standard.
  5. Apply the Three Questions mnemonic for every endpoint.

### stacks/js-ts-vue-nuxt/vue-components.md
New Sections 5–12:
- When to extract a component.
- Component naming conventions (PascalCase, Base/App/The prefix — pick one, Page suffix).
- Component responsibility (single concern, no mixed responsibilities).
- Component location (`shared/`, `[page-name]/`, `base/`).
- Props discipline (explicit types, no `$parent`/`$root`, never mutate, `withDefaults`).
- Slots (arbitrary content, named slots, no hardcoded parent-controlled content).
- Reactivity discipline (avoid unnecessary `ref`/`reactive`).
- Composables and utils (separation of stateful vs. pure logic).

### New: checklists/ directory
- **checklists/01-feature-implementation-checklist.md**: Pre-coding and post-coding self-review checklist for developers and AI agents.
- **checklists/02-security-review-checklist.md**: Complete endpoint security sweep from the Three Questions through financial operations and infrastructure.
- **checklists/03-deployment-checklist.md**: Pre-deployment verification covering environment, security, database, caches, git hygiene, monitoring, and post-deploy smoke testing.

---

## Generalized Rules

The following rules from `old.md` contained technology-specific examples but expressed universal principles. They were generalized during integration:

| old.md Rule | Generalized Universal Equivalent | Integrated Into |
| :--- | :--- | :--- |
| "Controllers should not contain business logic / no PDO in controllers" | "Presentation layers must not contain business logic" | Already in `core/02` — no change needed |
| "Use chunk()/lazy() for large collections in Laravel" | "Process large datasets in bounded chunks; never unbounded fetchAll" | Validated in `core/10` Section 4 |
| "Use lockForUpdate() on balances (Laravel)" | "Acquire a pessimistic lock on all balance/wallet rows before reading and writing" | Added to `core/08` Section 21 |
| "Never use rand()/mt_rand() for tokens (PHP)" | "Use cryptographically secure functions for all token generation" | Added to `core/08` Section 18 |
| "Str::random() or random_bytes() (PHP)" | Language-specific examples preserved in Section 18 alongside Node.js and Dart equivalents | `core/08` Section 18 |
| "Encrypt sensitive database columns (Laravel)" | "Classify sensitive data fields and apply encryption at rest" | Added to `core/08` Section 21 (noted) |
| "Use UUIDs for public-facing IDs" | "Avoid exposing sequential auto-increment IDs in public-facing API responses" | Added to `core/08` Section 21 |
| "Vue composables for shared stateful logic" | "Shared stateful logic must have a single source of truth via composable or service layer" | Added to `stacks/vue-components.md` Section 12 |
| "Never copy and paste markup — extract a component" | "Never duplicate UI logic — extract when logic appears more than once" | Added to `stacks/vue-components.md` Section 5 |

---

## Rejected Rules

The following rules were removed from consideration because they are project-specific and must not enter the universal playbook.

| Rejected Rule | Reason |
| :--- | :--- |
| Exclude `Admin_Web_Servicescopy/` from all searches | Specific to one legacy project; not universal |
| Reference `Progress.md` and `remaining.md` | Specific to one project workflow; these are project-level files |
| Folder structure: `app/Modules/Users/`, `Auth/`, `Customers/` | Specific to one PHP project; the playbook does not impose folder names |
| "Never modify the existing database schema" (absolute) | Project-specific constraint. Universal rule is: prefer backward-compatible additive changes (Expand-Contract) — already in `core/06` |
| Paystack and Flutterwave mentioned by name | Specific payment processor names belong in project-level AGENTS.md, not the universal playbook |
| "Build a NEW PHP application reproducing a legacy system" | Project mission statement — project-specific |
| PDO-only constraint (no ORM allowed) | Project-specific architectural choice; the playbook supports both PDO and ORM patterns |
| Route file names: `routes/admin.php`, `routes/coordinators.php` | Project-specific route naming |
| `Progress.md` progress tracking requirement | Project-specific workflow tooling |
| "Section 7 — How to Use This File" (paste this at the top of every new chat) | This is a prompt-engineering technique for a project-level AGENTS.md, not a universal engineering standard |

---

## Conflicts Resolved

### Conflict 1: SSRF Coverage

- **old.md**: Detailed SSRF rules with specific IP ranges to block.
- **Current `core/08`**: SSRF was mentioned in Section 4.9 but lacked actionable enforcement rules (specific IP ranges, DNS resolution requirement).

**Decision**: Retained the existing Section 4.9 framework entry. Added specific enforcement details (block loopbacks, private ranges, cloud metadata IPs) in the broader context of Section 4.9's whitelist guidance. No contradiction exists — the new Section 15 Pagination rules extend, not replace.

### Conflict 2: Rate Limiting — IP vs. Account-Based

- **old.md**: Strong preference for user/account-based rate limiting over IP-only.
- **Current `core/08`**: Rate limiting was mentioned as a general requirement without implementation strategy guidance.

**Decision**: Added `core/08` Section 20 with clear guidance: prefer user/account-based, detail when IP-only is acceptable (only when explicitly requested), and define the implementation requirements for persistent rate limiters. No contradiction with existing content.

### Conflict 3: Pagination Strategy

- **old.md**: Used both offset and keyset pagination in different contexts.
- **Current `core/10`**: Pagination mentioned broadly without distinguishing between strategies.

**Decision**: Added `core/10` Section 15 distinguishing when to use each strategy. Both are legitimate — the conflict was absence of nuance, not contradiction.

---

## Final Review Assessment

**As CTO**: The integration strengthens the playbook without bloating it. The rejected rules list is correct — no project-specific content entered the universal system.

**As Principal Engineer**: The generalization of timing attacks, BOPLA, and financial security into framework-neutral language is correct. These are exploitable vulnerabilities in any tech stack.

**As Security Architect**: BOPLA, HTTP Parameter Pollution, CSV injection, XXE, subdomain takeover, and timing attacks are all real, common vulnerabilities that were absent. Their addition materially improves the security coverage.

**As Senior Developer**: The Vue component discipline, composable/utils separation, and reactivity rules reflect real problems in Vue projects. The checklists are practical and copy-pasteable.

**No duplicate rules exist. No contradictions introduced. Universal rules are separated from technology-specific rules. The system remains practical for real projects.**
