# ADR Guidelines

These guidelines define when a code change warrants an Architecture Decision Record (ADR).

---

## When to Create an ADR

An ADR **must** be created or updated when you introduce, modify, or remove:
1. **Architecture Model**: Swapping paradigms (Monolith vs Microservices), altering communication flow.
2. **Security & Cryptography**: Authentication strategies (JWT vs Sessions), access tokens controls, cipher choices.
3. **Database & Persistence**: Adding a database engine (NoSQL/SQL), major schema changes, partitioning keys.
4. **Infrastructure & Deployments**: Container platforms, CI/CD runners, network topology.
5. **Caching & Queue Engines**: EViction rules, worker layout changes.
6. **Major Dependency Upgrades**: Adopting high-impact libraries or frameworks.
7. **Public API Structures**: Dynamic protocols (GraphQL vs REST), rate limiting.

---

## When NOT to Create an ADR

Do **not** create an ADR for:
- Bug fixes.
- Small refactoring edits.
- Visual styles or CSS changes.
- Standard minor feature additions.
- Formatting tasks.
